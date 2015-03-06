package register

import (
	dockerapi "github.com/fsouza/go-dockerclient"
	"github.com/helderfarias/hareg/discovery"
	"github.com/helderfarias/hareg/model"
	"github.com/helderfarias/hareg/util"
	"log"
	"regexp"
	"strings"
)

type DockerRegister struct {
	docker  *dockerapi.Client
	network util.Network
}

func NewDockerRegister(dockerHost string) *DockerRegister {
	docker, err := dockerapi.NewClient(dockerHost)
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
		this.registerContainer(listing.ID)
	}

	events := make(chan *dockerapi.APIEvents)
	this.docker.AddEventListener(events)
	go this.captureEvents(events)
}

func (this *DockerRegister) captureEvents(events chan *dockerapi.APIEvents) {
	for msg := range events {
		container, err := this.docker.InspectContainer(msg.ID)
		if err != nil {
			log.Printf("Unable to inspect container %s - %s, skipping", msg.ID, err)
			return
		}

		//name := container.Config.Hostname + "." + container.Config.Domainname + "."
		var running = regexp.MustCompile("start|^Up.*$")
		var stopping = regexp.MustCompile("die")

		switch {
		case running.MatchString(msg.Status):
			this.registerContainer(container.ID)
		case stopping.MatchString(msg.Status):
			this.removeContainer(container.ID)
		}
	}
}

func (this *DockerRegister) removeContainer(containerId string) {
	log.Println("Container removed", containerId)
}

func (this *DockerRegister) registerContainer(containerId string) {
	log.Println("Container added", containerId)

	container, err := this.docker.InspectContainer(containerId)
	if err != nil {
		log.Println("unable to inspect container:", containerId[:12], err)
		return
	}

	ip, port := this.getNetworkSettings(container.HostConfig.PortBindings, container.NetworkSettings.Ports)

	service := model.Service{}
	service.ExposedPort = this.mapperPort(port, container)
	service.ExposedIP = this.mapperIP(ip, container)
	service.ContainerID = container.ID[0:12]
	service.ContainerHostName = container.Name[1:]

	log.Println("Service ", service)
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
