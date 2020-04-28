package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/straightdave/raft/pb"
	"google.golang.org/grpc"
)

var (
	fAddr = flag.String("host", "127.0.0.1:8765", "Server Address")
)

func main() {
	flag.Parse()

	cid := 0

	reader := bufio.NewReader(os.Stdin)

	cc, err := grpc.Dial(*fAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("failed to dial server %s: %v\n", *fAddr, err)
		return
	}

	c := pb.NewRaftClient(cc)

	for {
		fmt.Printf("%s %d> ", *fAddr, cid)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		parts := strings.Split(strings.TrimSpace(input), " ")
		resp, err := c.Command(context.Background(), &pb.CommandRequest{
			Cid: strconv.Itoa(cid),
			Entry: &pb.CommandEntry{
				Command: parts[0],
				Args:    parts[1:],
			},
		})
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(resp.GetResult())
		cid++
	}
}
