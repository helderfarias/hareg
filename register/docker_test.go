package register

import (
	dockerapi "github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var reg *DockerRegister
var bindings, networks map[dockerapi.Port][]dockerapi.PortBinding

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func setup() {
	reg = NewDockerRegister("unix://localhost.sock")
	reg.network = &NetworkFake{}
	bindings = make(map[dockerapi.Port][]dockerapi.PortBinding)
	networks = make(map[dockerapi.Port][]dockerapi.PortBinding)
}

func TestGetNetworkSettingsByNetwork(t *testing.T) {
	networks[dockerapi.Port("100")] = []dockerapi.PortBinding{dockerapi.PortBinding{HostIP: "0.0.0.0", HostPort: "100"}, dockerapi.PortBinding{HostIP: "0.0.0.0", HostPort: "100"}}

	ip, port := reg.getNetworkSettings(bindings, networks)

	assert.Equal(t, "0.0.0.0", ip)
	assert.Equal(t, "100", port)
	assert.NotNil(t, reg.network)
}

func TestGetNetworkSettingsByBindinds(t *testing.T) {
	bindings[dockerapi.Port("100")] = []dockerapi.PortBinding{dockerapi.PortBinding{HostIP: "0.0.0.0", HostPort: "100"}, dockerapi.PortBinding{HostIP: "0.0.0.0", HostPort: "100"}}

	ip, port := reg.getNetworkSettings(bindings, networks)

	assert.Equal(t, "0.0.0.0", ip)
	assert.Equal(t, "100", port)
	assert.NotNil(t, reg.network)
}

type NetworkFake struct {
}

func (this *NetworkFake) ResolverIP(ip string) string {
	return "0.0.0.0"
}
