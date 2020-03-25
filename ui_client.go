package main

import (
	"flag"
	//"tap/backend/local"
	//"fmt"
	"google.golang.org/grpc"
	"log"
	"os"
	"tap/server"
	"tap/ui"
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

	conn, err := grpc.Dial("unix://"+server.UNIX_SOCK_FILE, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	rpcClient := server.NewPlayServerClient(conn)
	window := ui.NewWindow(rpcClient)

	if !window.ServerIsHealthLive() {
		panic("server is not running")
	}

	if !window.SetLocalProvider([]string{*dir}) {
		panic("set localprovider fatal")
	}

	window.Init()
	window.Close()
}
