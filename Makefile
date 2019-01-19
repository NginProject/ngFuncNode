build-windows: config
	go build ./cmd/ngFuncNode
	mv ngFuncNode* bin/ 

build-linux: config
	go build ./cmd/ngFuncNode
	mv ngFuncNode* bin/

config: fmt
	go run ./cmd/config2go > ./cmd/ngFuncNode/config.go

fmt:
	go fmt ./...

target: build-windows build-linux
