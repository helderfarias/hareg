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

var etcdHost = flag.String("etcd", util.GetOpt("ETCD_HOST", ""), "Address for the Etcd")
var dockerHost = flag.String("docker", util.GetOpt("DOCKER_HOST", "unix:///var/run/docker.sock"), "Address for the Docker daemon")
var dockerCertPath = flag.String("certpath", util.GetOpt("DOCKER_CERT_PATH", ""), "Docker Cert path")

func main() {
	flag.Parse()

	log.Println("Starting service register...")

	disc := discovery.NewEtcdDiscovery(*etcdHost)

	reg := selectRegister()

	log.Println("Listening for registers...")
	reg.RunAndWatch(disc)

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

func selectRegister() *register.DockerRegister {
	if *dockerCertPath != "" {
		cert := *dockerCertPath + "/cert.pem"
		key := *dockerCertPath + "/key.pem"
		ca := *dockerCertPath + "/ca.pem"
		return register.NewDockerTLSRegister(*dockerHost, cert, key, ca)
	}

	return register.NewDockerRegister(*dockerHost)
}
