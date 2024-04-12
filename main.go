package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	jsonrpc "github.com/ybbus/jsonrpc/v3"
)

var fsURL string

func main() {
	flag.StringVar(&fsURL, "fs-url", "", "")
	flag.Parse()

	fmt.Println(run(context.Background(), flag.Args()))
}

func run(ctx context.Context, urls []string) error {

	for {
		for _, url := range urls {
			client := NewSubpsaceClient(url)
			farmInfo, err := client.GetFarmerAppInfo(ctx)
			if err != nil {
				log.Println("request subspace node fail ", err)
				continue
			}
			if farmInfo.Syncing {
				sendFs(fsURL, fmt.Sprintf("%s 同步错误", url))
			}
		}

		time.Sleep(time.Minute)
	}
}

func sendFs(url string, msg string) error {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(fmt.Sprintf(`{"msg_type":"text","content":{"text":"%s"}}`, msg)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	return err
}

type SubspaceClient struct {
	client jsonrpc.RPCClient
}

func NewSubpsaceClient(url string) *SubspaceClient {
	rpcClient := jsonrpc.NewClient(url)
	return &SubspaceClient{
		rpcClient,
	}
}

type FarmAppInfo struct {
	GenesisHash       string         `json:"genesisHash"`
	DsnBootstrapNodes []string       `json:"dsnBootstrapNodes"`
	Syncing           bool           `json:"syncing"`
	FarmTimeout       map[string]int `json:"farmingTimeout"`
}

func (c *SubspaceClient) GetFarmerAppInfo(ctx context.Context) (FarmAppInfo, error) {
	result := FarmAppInfo{}
	resp, err := c.client.Call(ctx, "subspace_getFarmerAppInfo")
	if err != nil {
		return result, err
	}

	fmt.Println(resp.Result)
	err = resp.GetObject(&result)
	if err != nil {
		return result, err
	}
	return result, err
}
