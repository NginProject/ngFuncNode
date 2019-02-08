/*
Package Main
	Main src for ngFuncNode
	Working for post data to governor server

	TODO: EncryptJSON
*/
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net"
	"net/http"
	"os"
	"runtime"

	"strconv"
	"time"

	"github.com/NginProject/ngFuncNode/ngrpc"
	"github.com/NginProject/ngFuncNode/utils"
	"github.com/akamensky/argparse"
	"github.com/buger/jsonparser"
	gosocketio "github.com/mtfelian/golang-socketio"
	"github.com/mtfelian/golang-socketio/transport"
	// "github.com/rs/zerolog"
	// zlog "github.com/rs/zerolog/log"
)

// Submit local masternode config via http
func Submit(config *Config, data []byte) {
	url := config.Host
	var client = &http.Client{
		Timeout: time.Second * 6,
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
	}

	req.Header = map[string][]string{
		"Content-Type": {"application/json"},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}

	if err != nil {
		res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)
	status, err := jsonparser.GetInt(body, "status")
	if err != nil {
		log.Println(err.Error())
	}
	detail, err := jsonparser.GetString(body, "detail")
	if err != nil {
		log.Println(err.Error())
	}
	if status == 0 {
		log.Println(detail)
	}
}

type wsMsg struct {
	IP      string `json:"ip"`
	Account string `json:"account"`
	ENode   string `json:"enode"`
	Block   uint64 `json:"block"`
	// Time    int64  `json:"time"`
}

// Reply when server detect err and return it
type Reply struct {
	Status int    `json:"status,omitempty"`
	Detail string `json:"detail,omitempty"`
}

// NoData is an empty struct
type NoData struct{}

func main() {
	// output := zerolog.ConsoleWriter{
	// 	Out:        os.Stderr,
	// 	TimeFormat: time.RFC3339,
	// }
	// zlog.Logger = zlog.Output(output)

	parser := argparse.NewParser("ngFuncNode", `
	Multi-functional Node Client on Ngin Network. 
	Run "ngind --rpcapi admin,ngin,net,web3" before start ngFuncNode
	`)

	config := GetConfig()

	var ipFlag *string = parser.String("a", "address", &argparse.Options{Required: true, Help: "Your external IP address (reqiured)"})
	var portFlag *int = parser.Int("p", "port", &argparse.Options{Default: 52520, Required: false, Help: "Your external IP port"})
	var ngindFlag *bool = parser.Flag("n", "ngind", &argparse.Options{Default: false, Required: false, Help: "Whether auto-start ngind"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	// parse ip
	netIP := net.ParseIP(*ipFlag)
	if netIP == nil {
		log.Println("IP is wrong")
		os.Exit(0)
	}
	if !utils.IsPublicIP(netIP) {
		log.Println("Your IP is not a legel public address")
		os.Exit(0)
	}

	// parse port
	if *portFlag <= 0 || *portFlag > 65535 {
		log.Println("Wrong Port Number")
		os.Exit(0)
	}

	// parse ngind
	if *ngindFlag {
		if runtime.GOOS == "windows" {

		}
		if runtime.GOOS == "linux" {

		}
		if runtime.GOOS == "darwin" {

		}
	}

	n := ngrpc.New("", "", 0) // default

	account, err := n.GetCoinbase()
	if err != nil {
		log.Println(`Cannot get coinbase from ngind, open new coinbase with command: "./ngind account new"`)
		os.Exit(0)
	}

	enode, err := n.GetENode(*ipFlag, *portFlag)

	if err != nil {
		log.Println("Cannot get enode from ngind")
		os.Exit(0)
	}

	balance := &big.Int{}
	var blockNum uint64
	balanceChan := make(chan *big.Int)
	blockNumChan := make(chan uint64)

	go n.BlockNum2Chan(blockNumChan)
	go n.Balance2Chan(account, balanceChan)

	c, err := gosocketio.Dial(
		gosocketio.AddrWebsocket(config.Host, config.Port, false),
		&transport.WebsocketTransport{
			PingInterval:   6 * time.Second,
			PingTimeout:    60 * time.Second,
			ReceiveTimeout: 60 * time.Second,
			SendTimeout:    60 * time.Second,
			BufferSize:     1024 * 32,
		},
	)

	if err := c.On(gosocketio.OnConnection, func(channel *gosocketio.Channel, args interface{}) {
		log.Println("Connected")
		log.Println("New client connected, client id is ", c.Id())
		balance = <-balanceChan
		blockNum = <-blockNumChan
		era := int(blockNum / 100000)
		threshold := big.NewInt(int64(10 * int(math.Pow10(18)) * era * 2000))
		if balance.Cmp(threshold) < 0 {
			gap := threshold.Sub(threshold, balance)
			log.Println("Balance doesnt reach the threshold. More " + gap.String() + " wei NG needed")
			os.Exit(0)
		}
		msg := &wsMsg{
			IP:      *ipFlag,
			Account: account,
			ENode:   enode,
			Block:   blockNum,
		}
		c.Emit("register", msg)
	}); err != nil {
		log.Fatal(err)
	}

	if err := c.On(gosocketio.OnError, func(c *gosocketio.Channel, args *Reply) {
		log.Println("OnError")
		errReply := args
		log.Println("Error received: \nCode:" + strconv.Itoa(errReply.Status) + " - " + errReply.Detail)
	}); err != nil {
		log.Fatal(err)
	}

	if err := c.On("after_register", func(c *gosocketio.Channel, args *Reply) {
		log.Println("after_register")
		if args.Status == 1 {
			log.Println(args.Detail)
		} else {
			log.Fatal("register failed")
		}
	}); err != nil {
		log.Fatal(err)
	}

	if err := c.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel, args interface{}) {
		log.Println("Disconnected")
		// try reconnect
		// os.Exit(0)
	}); err != nil {
		log.Fatal(err)
	}

	if err := c.On("reward", func(c *gosocketio.Channel, args *Reply) {
		log.Println("New reward received: " + args.Detail)
	}); err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(10)
	}
}
