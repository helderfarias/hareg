package main

import (
	"flag"
	"github.com/helderfarias/hareg/discovery"
	"github.com/helderfarias/hareg/register"
	"github.com/helderfarias/hareg/util"
	"log"
	"os"
	"os/signal"
)

var etcdHost = flag.String("etcd", util.GetOpt("ETCD_HOST", "http://localhost:4001"), "Address for the Etcd")
var dockerHost = flag.String("docker", util.GetOpt("DOCKER_HOST", "unix:///var/run/docker.sock"), "Address for the Docker daemon")

func main() {
	flag.Parse()

	log.Println("Starting service register...")
	disc := discovery.NewEtcdDiscovery(*etcdHost)

	reg := register.NewDockerRegister(*dockerHost)

	log.Println("Listening for registers...")
	reg.RunAndWatch(disc)

	log.Println("Waiting for signals...")
	StartUp()
}

func StartUp() {
	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt)
	for {
		select {
		case <-s:
			log.Println("signal received, stopping")
			os.Exit(0)
		}
	}
}
