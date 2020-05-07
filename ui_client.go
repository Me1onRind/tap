package main

import (
	"flag"
	//"tap/backend/local"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"
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
	//logfile.Write([]byte("123"))
	log.SetOutput(logfile)
	//log.Println("123")
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
	window := ui.NewWindow(rpcClient)

	if !window.ServerIsHealthLive() {
		panic("server is not running")
	}

	if !window.SetLocalProvider([]string{*dir, "./testDir"}) {
		panic("set localprovider fatal")
	}

	defer func() {
		if err := recover(); err != nil {
			window.Close()
			fmt.Println(err)
		}
	}()
	window.Init()

	window.Close()
}
