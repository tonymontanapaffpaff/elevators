package main

import (
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"github.com/tonymontanapaffpaff/elevators/client"
	"github.com/tonymontanapaffpaff/elevators/server"
)

var (
	host = "127.0.0.1"
	port = ":1234"
)

func main() {
	rand.Seed(time.Now().Unix())
	as := os.Args[1]
	switch as {
	case "server":
		runServer(port, os.Args)
		break
	case "client":
		runClient(host+port, os.Args)
		break
	}
}

func runServer(serverAddress string, args []string) {
	floorCount, err := strconv.Atoi(args[2])
	if err != nil || floorCount <= 0 {
		panic("missing valid floor count argument")
	}

	elevatorCount, err := strconv.Atoi(args[3])
	if err != nil || elevatorCount <= 0 {
		panic("missing valid elevator count argument")
	}

	// start a server
	srv := server.New(floorCount, elevatorCount)
	err = rpc.Register(srv)
	// create a TCP listener that will listen on `Port`
	listener, _ := net.Listen("tcp", serverAddress)
	// close the listener whenever we stop
	defer listener.Close()
	// wait for incoming connections
	rpc.Accept(listener)
}

func runClient(serverAddress string, args []string) {
	// start a client
	cl := client.New()
	cl.Start(serverAddress)
	// close client whenever we stop
	defer cl.End()
	name := args[2]     // client name
	schedule := args[3] // client schedule in a specified format (example: 5:3_9:2),
	// where underscore separates trips, colon separates floor number and time staying
	cl.AddWorker(name, schedule)
}
