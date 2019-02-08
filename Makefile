build: config
	go build ./cmd/ngFuncNode

config: fmt
	go run ./cmd/config2go > ./cmd/ngFuncNode/config.go

fmt:
	go fmt ./...

target: build
