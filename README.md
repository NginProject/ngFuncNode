# ngFuncNode

Multi-functional Node Client on Ngin Network

## Build

```bash
git clone github.com/NginProject/ngFuncNode
cd ngFuncNode

# For Linux
make
./ngFuncNode -a 1.1.1.1 # Replace 1.1.1.1 with your IP

# For Windows
go run ./cmd/config2go > ./cmd/ngFuncNode/config.go
go build ./cmd/ngFuncNode
ngFuncNode.exe -a 1.1.1.1 # Replace 1.1.1.1 with your IP
```

## Usage

How to run a ngFuncNode to gain reward?

1. run a mainnet ngind daemon `./ngind --rpc`
2. Make sure your 52520 port is exposed to internet
3. run the ngFuncNode with external IP address flag `./ngFuncNode --a 1.1.1.1`(change `1.1.1.1` with your IP)
4. When both `ngind` and `ngFuncNode` are running all the time. If they stop, delivering reward to your coinbase account would stop in next 2 or 3 blocks.

For Windows, just need to convert the `./ngind` to `ngind.exe`, the `./ngFuncNode` to `ngFuncNode.exe`