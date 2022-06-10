package main

import (
	client "github.com/dthung1602/instant-proxy-api-client"
	"net"
)

func main() {
	server := client.NewFakeServer(
		3000,
		[]net.IP{
			net.ParseIP("123.1.1.1"),
		},
	)
	server.ServeForever()
}
