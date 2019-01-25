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

type GetCoinbaseJson struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      int      `json:"id"`
}

func GetCoinbase() (string, error) {
	url := "http://127.0.0.1:52521"
	var client = &http.Client{
		Timeout: time.Second * 6,
	}
	req_json := &GetCoinbaseJson{
		Jsonrpc: "2.0",
		Method:  "ngin_coinbase",
		Params:  []string{},
		Id:      0,
	}
	req_body, err := json.Marshal(req_json)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(req_body))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	addr, err := jsonparser.GetString(body, "result")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return addr, nil
}

func GetBalance(addr string, balanceChan chan *big.Int) {
	for {
		url := "http://127.0.0.1:52521"
		var client = &http.Client{
			Timeout: time.Second * 6,
		}
		req_json := &GetCoinbaseJson{
			Jsonrpc: "2.0",
			Method:  "ngin_getBalance",
			Params:  []string{addr, "latest"},
			Id:      0,
		}
		req_body, _ := json.Marshal(req_json)
		req, _ := http.NewRequest("POST", url, bytes.NewReader(req_body))
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		json, _, _, _ := jsonparser.Get(body, "result")
		i := new(big.Int)
		i.SetString(string(json)[2:], 16)
		balanceChan <- i
		time.Sleep(6 * time.Second)
	}
}

func GetBlockNum(blockNumChan chan uint64) {
	for {
		url := "http://127.0.0.1:52521"
		var client = &http.Client{
			Timeout: time.Second * 6,
		}
		req_json := &GetCoinbaseJson{
			Jsonrpc: "2.0",
			Method:  "ngin_blockNumber",
			Params:  []string{},
			Id:      0,
		}
		req_body, _ := json.Marshal(req_json)
		req, _ := http.NewRequest("POST", url, bytes.NewReader(req_body))
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		json, _, _, _ := jsonparser.Get(body, "result")
		i, _ := strconv.ParseUint(string(json)[2:], 16, 32)
		blockNumChan <- i
		time.Sleep(6 * time.Second)
	}
}

func GetENode(ip string, port int) (string, error) {
	url := "http://127.0.0.1:52521"
	var client = &http.Client{
		Timeout: time.Second * 6,
	}
	req_json := &GetCoinbaseJson{
		Jsonrpc: "2.0",
		Method:  "admin_nodeInfo",
		Params:  []string{},
		Id:      0,
	}
	req_body, err := json.Marshal(req_json)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(req_body))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	id, err := jsonparser.GetString(body, "result", "id")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Remember to run with `--rpcapi admin ...`")
		return "", err
	}
	enode := "enode://" + id + "@" + ip + ":" + string(port)
	return enode, nil
}
