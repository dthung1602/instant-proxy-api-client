package client

import (
	"net"
	"reflect"
	"testing"
)

// ------------------------------------
//   Test make proxy
// ------------------------------------

func TestMakeProxy(t *testing.T) {
	proxyString := "1.0.123.255:12345"
	expected := &Proxy{
		IP:   net.ParseIP("1.0.123.255"),
		Port: 12345,
	}

	proxyAdrr, err := MakeProxy(proxyString)

	if err != nil {
		t.Errorf("found error %v", err)
	}

	if !reflect.DeepEqual(*proxyAdrr, *expected) {
		t.Errorf("proxy adrr is not equal")
	}
}

func TestMakeProxyWithPortError(t *testing.T) {
	proxyString := "1.0.123.255:-1"
	_, err := MakeProxy(proxyString)

	if err == nil {
		t.Errorf("error should be returned")
	}

	proxyString = "1.0.123.255:65536"
	_, err = MakeProxy(proxyString)

	if err == nil {
		t.Errorf("error should be returned")
	}
}

func TestMakeProxies(t *testing.T) {
	proxyStrings := []string{
		"1.0.123.255:12345",
		"212.123.84.94:1250",
	}
	expected := []*Proxy{
		{
			IP:   net.ParseIP("1.0.123.255"),
			Port: 12345,
		},
		{
			IP:   net.ParseIP("212.123.84.94"),
			Port: 1250,
		},
	}

	proxies, err := MakeProxies(proxyStrings)

	if err != nil {
		t.Errorf("error not expected: %v", err)
		return
	}

	if !reflect.DeepEqual(proxies, expected) {
		t.Error("proxies are not equal")
		return
	}
}
