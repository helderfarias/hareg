package discovery

import (
	"github.com/coreos/go-etcd/etcd"
)

type EtcdDiscovery struct {
	etcdapi *etcd.Client
}

func NewEtcdDiscovery(etcdHost string) *EtcdDiscovery {
	etcdClient := etcd.NewClient([]string{etcdHost})
	return &EtcdDiscovery{etcdapi: etcdClient}
}

func (this *EtcdDiscovery) AddService() {

}
