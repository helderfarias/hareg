package util

import (
	"errors"
	"net"
)

type Network interface {
	ResolverIP(ip string) string
}

type InetAddress struct {
}

func (this *InetAddress) ResolverIP(ip string) string {
	if ip != "0.0.0.0" && ip != "" {
		return ip
	}

	local, err := lookup()
	if err != nil {
		return ""
	}

	return local.String()
}

func lookup() (net.IP, error) {
	tt, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, t := range tt {
		aa, err := t.Addrs()
		if err != nil {
			return nil, err
		}

		for _, a := range aa {
			ipnet, ok := a.(*net.IPNet)
			if !ok {
				continue
			}

			v4 := ipnet.IP.To4()
			if v4 == nil || v4[0] == 127 { // loopback address
				continue
			}

			return v4, nil
		}
	}

	return nil, errors.New("cannot find local IP address")
}
