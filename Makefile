build: config
	go build ./cmd/ngFuncNode

config: fmt
	go get ./cmd/config2go 
	config2go > ./cmd/ngFuncNode/config.go

fmt:
	go fmt ./...

target: build
