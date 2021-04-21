package client

import (
	"fmt"
	"net/rpc"
	"strconv"
	"strings"

	"github.com/tonymontanapaffpaff/elevators/server"
)

type Client struct {
	rpcClient *rpc.Client
}

func New() *Client {
	c := Client{}
	return &c
}

func (c *Client) Start(addr string) {
	// establish the connection of the RPC server
	rpcClient, err := rpc.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	c.rpcClient = rpcClient
}

func (c *Client) AddWorker(name string, schedule string) {
	// parse schedule string
	pairs := strings.Split(schedule, "_")
	schedulePairs := make([]server.WorkerSchedulePair, len(pairs))
	for i, pair := range pairs {
		s := strings.Split(pair, ":")
		schedulePairs[i].Floor, _ = strconv.Atoi(s[0])
		schedulePairs[i].Seconds, _ = strconv.Atoi(s[1])
	}
	// add worker to server
	workerResponse := server.WorkerResponse{}
	err := c.rpcClient.Call("Server.AddWorker",
		&server.WorkerRequest{Name: name, Schedule: schedulePairs}, &workerResponse)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(workerResponse.Message, schedule)
}

func (c *Client) End() {
	c.rpcClient.Close()
}
