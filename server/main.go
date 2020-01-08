package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	pb "github.com/straightdave/raft/server/pb"
	grpc "google.golang.org/grpc"
)

var (
	fPort            = flag.Uint("port", 8801, "local port to listen")
	fOtherServerList = flag.String("servers", "", "server list separated by commas")

	otherServers []string
	terminated   chan struct{}
)

func init() {
	flag.Parse()
	otherServers = strings.Split(*fOtherServerList, ",")
	terminated = make(chan struct{})
}

func main() {
	ctx := context.Background()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		close(terminated)
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *fPort))
	if err != nil {
		log.Fatal(err)
	}

	gSvr := grpc.NewServer()
	pb.RegisterRaftServer(gSvr, NewServerServiceImpl(ctx, otherServers))

	go func() {
		log.Printf("Serving at :%d", *fPort)
		if err := gSvr.Serve(lis); err != nil {
			close(terminated)
			log.Fatal(err)
		}
	}()

	<-terminated
	log.Printf("Server terminated")
}
