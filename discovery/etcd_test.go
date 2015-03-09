package discovery

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitialize(t *testing.T) {
	disc := NewEtcdDiscovery("localhost:4001")

	assert.NotNil(t, disc.etcdapi)
}

func TestCreateEntries(t *testing.T) {
	entries := makeEntries("10.10.10.2:4001,10.10.10.3:4001")

	assert.NotEmpty(t, entries)
	assert.Equal(t, []string{"10.10.10.2:4001", "10.10.10.3:4001"}, entries)
}
