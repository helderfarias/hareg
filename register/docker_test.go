package register

import (
	dockerapi "github.com/fsouza/go-dockerclient"
	"testing"
)

func TestShouldBeBindingGetDefaultPort(t *testing.T) {
	reg := NewDockerRegister("unix://localhost.sock")
	bindings := make(map[dockerapi.Port][]dockerapi.PortBinding)
	networks := make(map[dockerapi.Port][]dockerapi.PortBinding)

	bindings[dockerapi.Port("200")] = []dockerapi.PortBinding{dockerapi.PortBinding{HostIP: "0.0.0.0", HostPort: "100/tcp"}}

	ip, port := reg.getDefaultPort(bindings, networks)

	if ip != "0.0.0.0" {
		t.Fatalf("Expected %s, but %s", "0.0.0.0", ip)
	}

	if port != "100" {
		t.Fatalf("Expected %s, but %s", "100", port)
	}
}
