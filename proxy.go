package client

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Proxy struct {
	IP   net.IP
	Port uint16
}

func (proxy *Proxy) String() string {
	return fmt.Sprintf("%s:%d", proxy.IP.String(), proxy.Port)
}

func MakeProxy(str string) (*Proxy, error) {
	str = strings.Trim(str, " \n\r\t")
	parts := strings.Split(str, ":")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid proxy string '%s'", str)
	}

	ip := net.ParseIP(parts[0])
	port, portErr := strconv.Atoi(parts[1])

	if port < 0 || port > 65535 {
		portErr = errors.New("port out of range")
	}

	if ip == nil || portErr != nil {
		return nil, fmt.Errorf("invalid proxy string '%s'", str)
	}

	return &Proxy{ip, uint16(port)}, nil
}

func MakeProxies(strings []string) ([]*Proxy, error) {
	proxies := make([]*Proxy, len(strings))
	for i, line := range strings {
		proxy, parseErr := MakeProxy(line)
		if parseErr != nil {
			return nil, parseErr
		}
		proxies[i] = proxy
	}
	return proxies, nil
}
