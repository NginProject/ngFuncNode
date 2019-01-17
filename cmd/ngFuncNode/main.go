// Main src for ngFuncNode
// Working for post data to governor server
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"os"

	"github.com/NginProject/ngFuncNode/ip"
	"github.com/NginProject/ngFuncNode/ngrpc"
	"github.com/buger/jsonparser"
)

// Submit local masternode config
func Submit(config *Config, data []byte) {
	url := config.Server
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		fmt.Println(err)
	}

	req.Header = map[string][]string{
		"Content-Type": {"application/json"},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	status, err := jsonparser.GetString(body, "status")
	if err != nil {
		fmt.Println(err)
	}
	detail, err := jsonparser.GetString(body, "detail")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(status)
	fmt.Println(detail)

}

// TODO: Encrypt the post form
func Encrypt(text string) (string, error) {
	key := []byte{0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
		0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
	}
	var iv = key[:aes.BlockSize]
	encrypted := make([]byte, len(text))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypter.XORKeyStream(encrypted, []byte(text))
	return hex.EncodeToString(encrypted), nil
}

func main() {
	config := GetConfig()
	fmt.Println(`Make sure your ngind with rpc is running, if not, run this command: ./ngind --rpc`)
	fmt.Println(`Make sure your 52520 port is open`)
	ip_str := ip.GetPublicIP()
	addr, err := ngrpc.GetCoinbase()
	if err != nil {
		fmt.Println("Cannot get coinbase from ngind, open new coinbase with command: ./ngind account new")
		os.Exit(0)
	}

	data := []byte(`{"ip": "", "address": "", "balance": "" }`)
	data, err = jsonparser.Set(data, []byte(`"`+ip_str+`"`), "ip")
	data, err = jsonparser.Set(data, []byte(`"`+addr+`"`), "address")
	fmt.Println(err)
	fmt.Println(string(data))
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
			fmt.Println("Balance doesnt reach the threshold. More ", gap, " NG needed")
			os.Exit(0)
		}
		data, err = jsonparser.Set(data, []byte(`"`+balance.String()+`"`), "balance")
		fmt.Println(string(data))
		Submit(config, data)
	}
}
