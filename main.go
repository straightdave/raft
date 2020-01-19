package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	pb "github.com/straightdave/raft/pb"
	grpc "google.golang.org/grpc"
)

var (
	fPort         = flag.Uint("port", 8765, "local port to listen (TCP)")
	fOtherServers = flag.String("servers", "", "the initial server list separated by commas")
)

func main() {
	flag.Parse()

	otherServers := strings.Split(*fOtherServers, ",")
	terminated := make(chan struct{})

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
	pb.RegisterRaftServer(gSvr, NewServerServiceImpl(otherServers))

	go func() {
		log.Printf("Serving TCP connections at :%d", *fPort)
		if err := gSvr.Serve(lis); err != nil {
			close(terminated)
			log.Fatal(err)
		}
	}()

	<-terminated
	log.Printf("Terminated")
}
