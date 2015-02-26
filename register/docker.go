package register

import (
	dockerapi "github.com/fsouza/go-dockerclient"
	"github.com/helderfarias/hareg/discovery"
	"github.com/helderfarias/hareg/model"
	"log"
	"regexp"
)

type DockerRegister struct {
	docker *dockerapi.Client
}

func NewDockerRegister(dockerHost string) *DockerRegister {
	docker, err := dockerapi.NewClient(dockerHost)
	if err != nil {
		log.Fatal(err)
	}

	return &DockerRegister{docker: docker}
}

func (this *DockerRegister) RunAndWatch(disc *discovery.EtcdDiscovery) {
	log.Println("Starting")

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
		log.Printf("Event: %s %s", msg.ID, msg.Status)

		container, err := this.docker.InspectContainer(msg.ID)
		if err != nil {
			log.Printf("Unable to inspect container %s, skipping", msg.ID)
			return
		}

		name := container.Config.Hostname + "." + container.Config.Domainname + "."
		var running = regexp.MustCompile("start|^Up.*$").MatchString(msg.Status)
		var stopping = regexp.MustCompile("die").MatchString(msg.Status)

		switch {
		case running:
			log.Printf("Adding record for %v", name)
		case stopping:
			log.Printf("Removing record for %v", name)
		}
	}
}

func (this *DockerRegister) registerContainer(containerId string) {
	container, err := this.docker.InspectContainer(containerId)
	if err != nil {
		log.Println("unable to inspect container:", containerId[:12], err)
		return
	}

	service := model.Service{ExposedIP: container.NetworkSettings.IPAddress}
	log.Println(service)

	for port, published := range container.HostConfig.PortBindings {
		log.Println(port, published)
	}

	for port, published := range container.NetworkSettings.Ports {
		log.Println(port, published)
	}
}
