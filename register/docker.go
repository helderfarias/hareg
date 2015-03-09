package register

import (
	dockerapi "github.com/fsouza/go-dockerclient"
	"github.com/helderfarias/hareg/discovery"
	"github.com/helderfarias/hareg/model"
	"github.com/helderfarias/hareg/util"
	"log"
	"regexp"
	"strings"
	"sync"
)

type DockerRegister struct {
	docker  *dockerapi.Client
	network util.Network
	sync.RWMutex
}

var servicesCache map[string]model.Service

func init() {
	servicesCache = make(map[string]model.Service)
}

func NewDockerRegister(endpoint string) *DockerRegister {
	docker, err := dockerapi.NewClient(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	return &DockerRegister{
		docker:  docker,
		network: &util.InetAddress{},
	}
}

func NewDockerTLSRegister(endpoint, cert, key, ca string) *DockerRegister {
	docker, err := dockerapi.NewTLSClient(endpoint, cert, key, ca)
	if err != nil {
		log.Fatal(err)
	}

	return &DockerRegister{docker: docker, network: &util.InetAddress{}}
}

func (this *DockerRegister) RunAndWatch(disc *discovery.EtcdDiscovery) {
	containers, err := this.docker.ListContainers(dockerapi.ListContainersOptions{})
	if err != nil {
		log.Fatalf("Unable to register running containers: %v", err)
	}

	for _, listing := range containers {
		this.registerContainer(listing.ID, disc)
	}

	events := make(chan *dockerapi.APIEvents)
	this.docker.AddEventListener(events)
	go this.captureEvents(events, disc)
}

func (this *DockerRegister) captureEvents(events chan *dockerapi.APIEvents, disc *discovery.EtcdDiscovery) {
	for msg := range events {
		var running = regexp.MustCompile("start|^Up.*$")
		var stopping = regexp.MustCompile("die")

		switch {
		case running.MatchString(msg.Status):
			this.registerContainer(msg.ID[:12], disc)
		case stopping.MatchString(msg.Status):
			this.removeContainer(msg.ID[:12], disc)
		}
	}
}

func (this *DockerRegister) removeContainer(containerId string, disc *discovery.EtcdDiscovery) {
	service := servicesCache[containerId]

	if service.ContainerID != "" {
		this.Lock()
		err := disc.RemoveService(service)

		if err == nil {
			delete(servicesCache, service.ContainerID)
		}
		this.Unlock()
	}
}

func (this *DockerRegister) registerContainer(containerId string, disc *discovery.EtcdDiscovery) {
	container, err := this.docker.InspectContainer(containerId)
	if err != nil {
		log.Println("unable to inspect container:", containerId[:12], err)
		return
	}

	domain, endpoint := this.getServiceDomain(container)
	if len(domain) == 0 || len(endpoint) == 0 {
		return
	}

	ip, port := this.getNetworkSettings(container.HostConfig.PortBindings, container.NetworkSettings.Ports)

	service := model.Service{}
	service.ExposedPort = this.mapperPort(port, container)
	service.ExposedIP = this.mapperIP(ip, container)
	service.ContainerID = container.ID[0:12]
	service.ContainerHostName = container.Name[1:]
	service.Domain = domain
	service.Endpoint = endpoint

	this.Lock()
	errService := disc.AddService(service)

	if errService == nil {
		servicesCache[service.ContainerID] = service
	}
	this.Unlock()
}

func (this *DockerRegister) getNetworkSettings(bindings, networks map[dockerapi.Port][]dockerapi.PortBinding) (string, string) {
	var hostIP, hostPort string

	for port, published := range bindings {
		if len(published) > 0 {
			hostIP = published[0].HostIP
			hostPort = published[0].HostPort
		} else {
			hostPort = string(port)
		}
	}

	for port, published := range networks {
		if len(published) > 0 {
			hostIP = published[0].HostIP
			hostPort = published[0].HostPort
		} else {
			hostPort = string(port)
		}
	}

	return hostIP, hostPort
}

func (this *DockerRegister) mapperPort(port string, container *dockerapi.Container) string {
	portDefault := port

	if port == "" {
		for exposed := range container.Config.ExposedPorts {
			if exposed.Port() != "" {
				portDefault = exposed.Port()
			}
		}
	}

	if strings.Contains(portDefault, "/") {
		return strings.Split(portDefault, "/")[0]
	}

	return portDefault
}

func (this *DockerRegister) mapperIP(ip string, container *dockerapi.Container) string {
	if ip == "" && container.HostConfig.NetworkMode == "bridge" {
		return container.NetworkSettings.IPAddress
	}

	return this.network.ResolverIP(ip)
}

func (this *DockerRegister) getServiceDomain(container *dockerapi.Container) (domain, endpoint string) {
	for _, env := range container.Config.Env {
		if strings.HasPrefix(env, "SERVICE_DOMAIN") {
			pairs := strings.Split(env, "=")[1]
			domain := strings.Split(pairs, ":")[0]
			endpoint := strings.Split(pairs, ":")[1]
			return domain, endpoint
		}
	}

	return "", ""
}
