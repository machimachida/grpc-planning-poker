package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/machimachida/grpc-planning-poker/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	name := flag.String("name", "SatoTaro", "name of user")
	waitSecond := flag.Int("wait", 600, "wait second")
	flag.Parse()

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewPlanningPokerClient(conn)
	stream := connect(client, *name)
	go func(*pb.PlanningPokerClient) {
		for {
			if message, err := stream.Recv(); err != nil {
				break
			} else {
				println(color.HiGreenString("broadcast: " + message.Message))
			}
		}
	}(&client)

	go func() {
		time.Sleep(time.Duration(*waitSecond) * time.Second)
		disconnect(client, *name)
		os.Exit(0)
	}()

	fmt.Println("Start planning poker!")
	for {
		var n int32
		fmt.Scan(&n)
		if n == -1 {
			fmt.Println("Disconnect...")
			break
		}
		if n == 0 {
			showVotes(client, *name)
			continue
		}
		vote(client, *name, n)
		fmt.Println("Waiting others' votes...")
	}

	disconnect(client, *name)
}

func connect(client pb.PlanningPokerClient, id string) pb.PlanningPoker_ConnectClient {
	stream, err := client.Connect(context.Background(), &pb.ConnectRequest{Id: id})
	if err != nil {
		panic(err)
	}
	return stream
}

func disconnect(client pb.PlanningPokerClient, id string) {
	recv, err := client.Disconnect(context.Background(), &pb.DisconnectRequest{Id: id})
	if err != nil {
		panic(err)
	}
	println(recv.Message)
}

func vote(client pb.PlanningPokerClient, id string, vote int32) {
	recv, err := client.Vote(context.Background(), &pb.VoteRequest{Id: id, Vote: vote})
	if err != nil {
		panic(err)
	}
	println(recv.Message)
}

func showVotes(client pb.PlanningPokerClient, id string) {
	recv, err := client.ShowVotes(context.Background(), &pb.ShowVotesRequest{Id: id})
	if err != nil {
		panic(err)
	}
	println(recv.Message)
}
