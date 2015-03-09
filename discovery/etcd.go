package discovery

import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/helderfarias/hareg/model"
	"log"
	"strings"
)

type EtcdDiscovery struct {
	etcdapi *etcd.Client
	url     string
}

func NewEtcdDiscovery(host string) *EtcdDiscovery {
	etcdClient := etcd.NewClient([]string{host})
	return &EtcdDiscovery{etcdapi: etcdClient, url: host}
}

func (this *EtcdDiscovery) AddService(srv model.Service) error {
	log.Println("Container added", srv)

	domain := "services/" + srv.Domain + "/domain"
	endpoint := srv.Endpoint
	_, errDomain := this.etcdapi.Set(domain, endpoint, 0)
	if errDomain != nil {
		log.Println("etcd: failed to register domain:", errDomain)
		return errDomain
	}

	backend := "services/" + srv.Domain + "/backend/" + srv.ContainerID
	server := srv.ExposedIP + ":" + srv.ExposedPort
	_, errBackend := this.etcdapi.Set(backend, server, 0)
	if errBackend != nil {
		log.Println("etcd: failed to register backend:", errBackend)
		return errBackend
	}

	return nil
}

func (this *EtcdDiscovery) RemoveService(srv model.Service) error {
	log.Println("Container removed", srv.ContainerID)

	backend := "services/" + srv.Domain + "/backend/" + srv.ContainerID
	_, errBackend := this.etcdapi.Delete(backend, true)
	if errBackend != nil {
		log.Println("etcd: failed to remove backend:", errBackend)
		return errBackend
	}

	resp, err := this.etcdapi.Get("services/"+srv.Domain+"/backend", true, true)
	if err != nil {
		log.Println("etc: failed get")
		return err
	}

	if resp.Node.Nodes.Len() == 0 {
		_, err := this.etcdapi.Delete("services/"+srv.Domain, true)
		if err != nil {
			log.Println("etc: failed delete domain")
			return err
		}
	}

	return nil
}

func makeEntries(hosts string) []string {
	pairs := strings.Split(hosts, ",")

	entries := make([]string, 0)
	for _, value := range pairs {
		entries = append(entries, value)
	}

	return entries
}
