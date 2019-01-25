// Main src for ngFuncNode
// Working for post data to governor server
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/NginProject/ngFuncNode/ip"
	"github.com/NginProject/ngFuncNode/ngrpc"
	"github.com/akamensky/argparse"
	"github.com/buger/jsonparser"
	"github.com/jbrodriguez/mlog"
)

// Submit local masternode config
func Submit(config *Config, data []byte) {
	url := config.Server
	var client = &http.Client{
		Timeout: time.Second * 6,
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		mlog.Warning(err.Error())
	}

	req.Header = map[string][]string{
		"Content-Type": {"application/json"},
	}

	res, err := client.Do(req)
	if err != nil {
		mlog.Warning(err.Error())
	}

	if err != nil {
		res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)
	status, err := jsonparser.GetInt(body, "status")
	if err != nil {
		mlog.Warning(err.Error())
	}
	detail, err := jsonparser.GetString(body, "detail")
	if err != nil {
		mlog.Warning(err.Error())
	}
	if status == 0 {
		mlog.Info(detail)
	}

}

// TODO: EncryptJSON

func main() {
	mlog.Start(mlog.LevelInfo, "")

	parser := argparse.NewParser("ngFuncNode", "Multi-functional Node Client on Ngin Network")

	config := GetConfig()

	fmt.Println(`Make sure your ngind with rpc is running, if not, run this command: ./ngind --rpc`)
	fmt.Println(`Make sure your 52520 port is open`)

	var ip_pstr *string = parser.String("a", "address", &argparse.Options{Required: true, Help: "Your external IP address (reqiured)."})
	var port_pint *int = parser.Int("p", "port", &argparse.Options{Default: 52520, Required: false, Help: "Your external IP port (reqiured)."})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	net_ip := net.ParseIP(*ip_pstr)
	if net_ip == nil {
		mlog.Warning("IP is wrong")
		os.Exit(0)
	}
	if !ip.IsPublicIP(net_ip) {
		mlog.Warning("Your IP is not a legel public address")
		os.Exit(0)
	}

	if *port_pint > 65535 {
		fmt.Println("Wrong Port Number")
		os.Exit(0)
	}

	addr, err := ngrpc.GetCoinbase()
	if err != nil {
		mlog.Warning("Cannot get coinbase from ngind, open new coinbase with command: ./ngind account new")
		os.Exit(0)
	}

	enode, err := ngrpc.GetENode(*ip_pstr, *port_pint)
	if err != nil {
		mlog.Warning("Cannot get enode address from ngind")
		os.Exit(0)
	}

	data := []byte(`{"ip": "", "address": "", "enode": "", "balance": "" }`)
	if data, err = jsonparser.Set(data, []byte(`"`+*ip_pstr+`"`), "ip"); err != nil {
		mlog.Warning(err.Error())
	}
	if data, err = jsonparser.Set(data, []byte(`"`+addr+`"`), "address"); err != nil {
		mlog.Warning(err.Error())
	}
	if data, err = jsonparser.Set(data, []byte(`"`+enode+`"`), "enode"); err != nil {
		mlog.Warning(err.Error())
	}
	balance := &big.Int{}
	var blockNum uint64
	balanceChan := make(chan *big.Int)
	blockNumChan := make(chan uint64)

	go ngrpc.GetBlockNum(blockNumChan)
	go ngrpc.GetBalance(addr, balanceChan)

	for {
		balance = <-balanceChan
		blockNum = <-blockNumChan
		era := int(blockNum / 100000)
		threshold := big.NewInt(int64(10 * int(math.Pow10(18)) * era * 2000))
		if balance.Cmp(threshold) < 0 {
			gap := threshold.Sub(threshold, balance)
			fmt.Println("Balance doesnt reach the threshold. More ", gap, " wei NG needed")
			os.Exit(0)
		}
		data, err = jsonparser.Set(data, []byte(`"`+balance.String()+`"`), "balance")
		//fmt.Println(string(data))
		Submit(config, data)
	}
}
