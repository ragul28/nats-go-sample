# NATS Go sample

Simple NATS project.  

## NATS core PubSub

* Run nats in docker
```
docker compose up -d nats
```

* Testing nats pub sub
```sh
go run cmd/sub/main.go foo
```

```sh
go run cmd/pub/main.go foo bar
```