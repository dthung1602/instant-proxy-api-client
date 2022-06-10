
test:
	go test ./...

fake-server:
	cd resources && go run fake_instant_proxy_server_up.go
