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

	"google.golang.org/grpc"

	"github.com/straightdave/raft"
	"github.com/straightdave/raft/kv"
	"github.com/straightdave/raft/pb"
)

var (
	fPort  = flag.Uint("raft.port", 8765, "local port to listen (TCP)")
	fPeers = flag.String("raft.peers", "", "the initial server list separated by commas")
)

func main() {
	flag.Parse()

	var peers []string
	raw := strings.TrimSpace(*fPeers)
	for _, p := range strings.Split(raw, ",") {
		pp := strings.TrimSpace(p)
		if pp != "" {
			peers = append(peers, pp)
		}
	}

	log.Printf("Starts with %d peer(s): %v\n", len(peers), peers)

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
	pb.RegisterRaftServer(
		gSvr,
		raft.NewRaftServer(
			*fPort,
			raft.WithPeers(peers),
			raft.WithExecutor(kv.NewExecutor()),
		),
	)

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
