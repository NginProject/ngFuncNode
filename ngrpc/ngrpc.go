// Package ngrpc:
// ngrpc is a simple ngFuncNode's library to call some ngind jsonrpc

package ngrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
)

// LocalNgindRPCClient is one ngind intance in local machine
type LocalNgindRPCClient struct {
	Schema string
	Host   string
	Port   int
}

// StrListParamsReq is the jsonrpc request with string list params
type StrListParamsReq struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	ID      int      `json:"id"`
}

// New a LocalNgindRpc
func New(schema, host string, port int) *LocalNgindRPCClient {
	if schema != "https" {
		schema = "http"
	}
	if host == "" {
		host = "127.0.0.1"
	}
	if port <= 0 || port > 65535 {
		port = 52521
	}
	return &LocalNgindRPCClient{
		Schema: schema,
		Host:   host,
		Port:   port,
	}

}

// GetCoinbase will get string-type coinbase from local ngind
func (n *LocalNgindRPCClient) GetCoinbase() (string, error) {
	url := n.Schema + "://" + n.Host + ":" + strconv.Itoa(n.Port)
	var client = &http.Client{
		Timeout: time.Second * 6,
	}
	reqJSON := &StrListParamsReq{
		Jsonrpc: "2.0",
		Method:  "ngin_coinbase",
		Params:  []string{},
		ID:      0,
	}
	reqBody, err := json.Marshal(reqJSON)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if err != nil {
		res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)
	addr, err := jsonparser.GetString(body, "result")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return addr, nil
}

// GetBalance will get the account's balance from local ngind
func (n *LocalNgindRPCClient) GetBalance(account string) *big.Int {
	url := n.Schema + "://" + n.Host + ":" + strconv.Itoa(n.Port)
	var client = &http.Client{
		Timeout: time.Second * 6,
	}
	reqJSON := &StrListParamsReq{
		Jsonrpc: "2.0",
		Method:  "ngin_getBalance",
		Params:  []string{account, "latest"},
		ID:      0,
	}
	reqBody, _ := json.Marshal(reqJSON)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)
	json, _, _, _ := jsonparser.Get(body, "result")
	i := new(big.Int)
	i.SetString(string(json)[2:], 16)
	return i
}

// GetBlockNum will get the latest block number for local ngind
func (n *LocalNgindRPCClient) GetBlockNum() uint64 {
	url := n.Schema + "://" + n.Host + ":" + strconv.Itoa(n.Port)
	var client = &http.Client{
		Timeout: time.Second * 6,
	}
	reqJSON := &StrListParamsReq{
		Jsonrpc: "2.0",
		Method:  "ngin_blockNumber",
		Params:  []string{},
		ID:      0,
	}
	reqBody, _ := json.Marshal(reqJSON)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)
	json, _, _, _ := jsonparser.Get(body, "result")
	i, _ := strconv.ParseUint(string(json)[2:], 16, 32)
	return i
}

// GetENode returns "enode://[key]@[ip]:[port]"
func (n *LocalNgindRPCClient) GetENode(ip string, port int) (string, error) {
	url := n.Schema + "://" + n.Host + ":" + strconv.Itoa(n.Port)
	var client = &http.Client{
		Timeout: time.Second * 6,
	}
	reqJSON := &StrListParamsReq{
		Jsonrpc: "2.0",
		Method:  "admin_nodeInfo",
		Params:  []string{},
		ID:      0,
	}
	reqBody, err := json.Marshal(reqJSON)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if err != nil {
		res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)
	id, err := jsonparser.GetString(body, "result", "id")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Remember to run with `--rpcapi admin,ngin,net,web3`")
		return "", err
	}
	enode := "enode://" + id + "@" + ip + ":" + strconv.Itoa(port)
	return enode, nil
}

// Balance2Chan send balance to *big.Int channel
func (n *LocalNgindRPCClient) Balance2Chan(addr string, balanceChan chan *big.Int) {
	for {
		balanceChan <- n.GetBalance(addr)
		time.Sleep(6 * time.Second)
	}
}

// BlockNum2Chan send block number to uint64 channel
func (n *LocalNgindRPCClient) BlockNum2Chan(blockNumChan chan uint64) {
	for {
		blockNumChan <- n.GetBlockNum()
		time.Sleep(6 * time.Second)
	}
}
