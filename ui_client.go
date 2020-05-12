package main

import (
	"flag"
	//"tap/backend/local"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"
	"tap/rpc_client"
	"tap/server"
	"tap/ui"
	"time"
)

var (
	dir *string
)

func main() {
	logfile, err := os.OpenFile("/tmp/ui", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logfile)
	dir = flag.String("dir", "./", "")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "unix://"+server.UNIX_SOCK_FILE, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	rpcClient := server.NewPlayClient(conn)
	rpc_client.SetRpcClient(rpcClient)

	window := ui.NewWindow()
	defer func() {
		if err := recover(); err != nil {
			window.Close()
			fmt.Println(err)
		}
	}()

	if !rpc_client.ServerIsHealthLive() {
		panic("server is not running")
	}

	if !rpc_client.SetLocalProvider([]string{*dir, "./testDir"}) {
		panic("set localprovider fatal")
	}
	rpc_client.SetOutput(window.GetOutput())

	window.Init()

	window.Close()
}
