package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolverIPDefaullt(t *testing.T) {
	inet := new(InetAddress)

	ip := inet.ResolverIP("0.0.0.0")

	assert.NotEmpty(t, ip)
	assert.NotEqual(t, "127.0.0.1", ip)
}

func TestResolverIPIsEmpty(t *testing.T) {
	inet := new(InetAddress)

	ip := inet.ResolverIP("")

	assert.NotEmpty(t, ip)
	assert.NotEqual(t, "127.0.0.1", ip)
}
