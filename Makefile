build_nc:
	go build -o nats-pub cmd/pub/main.go
	go build -o nats-sub cmd/sub/main.go

install:
	go install

mod_init:
	go mod init github.com/ragul28/nats-go-sample
	go get -u

mod:
	go mod tidy
	go mod verify
