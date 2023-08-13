package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/machimachida/grpc-planning-poker/pb"
)

const (
	AVERAGE = "average"
)

func main() {
	name := flag.String("name", "SatoTaro", "name of user")
	waitSecond := flag.Int("wait", 600, "wait second")
	flag.Parse()

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	client := pb.NewPlanningPokerClient(conn)
	stream, cancel := connect(client, *name)
	go listenServerMessage(stream)
	go disconnectAfterWaitSecond(cancel, *waitSecond)

	println("Start planning poker!")

	for {
		var n int32
		_, err = fmt.Scan(&n)
		if err != nil {
			log.Println("failed to scan input. Please input integer.", err)
			continue
		}

		switch n {
		case -1:
			println("Disconnect...")
			disconnect(cancel)
		case -2:
			showVotes(client, *name)
			continue
		case -3:
			newGame(client)
			continue
		default:
			if n > 0 {
				vote(client, *name, n)
				println("Waiting others' votes...")
				continue
			}

			println("Please input natural number, -1(disconnect), -2(show votes) or -3(new game).")
		}
	}
}

func listenServerMessage(stream pb.PlanningPoker_ConnectClient) {
	for {
		message, err := stream.Recv()

		// 異常系
		if err != nil {
			// contextがキャンセルされた場合は、正常系として扱う
			if stream.Context().Err() == context.Canceled {
				log.Println("server is disconnected.")
				os.Exit(0)
			}

			log.Println("server is disconnected with error.", err)
			os.Exit(0)
		}
		if message == nil {
			log.Println("server sent nil message.")
			continue
		}

		// 正常系
		printlnBroadcastMessage(message)
	}
}

func disconnectAfterWaitSecond(cancel context.CancelFunc, waitSecond int) {
	time.Sleep(time.Duration(waitSecond) * time.Second)
	disconnect(cancel)
	os.Exit(0)
}

func connect(client pb.PlanningPokerClient, id string) (pb.PlanningPoker_ConnectClient, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := client.Connect(ctx, &pb.ConnectRequest{Id: id})
	if err != nil {
		panic(err)
	}
	return stream, cancel
}

func disconnect(cancel context.CancelFunc) {
	cancel()
}

func vote(client pb.PlanningPokerClient, id string, vote int32) {
	recv, err := client.Vote(context.Background(), &pb.VoteRequest{Id: id, Vote: vote})
	if err != nil {
		log.Println("failed to vote.", err)
	}
	println(recv.Message)
}

func showVotes(client pb.PlanningPokerClient, id string) {
	recv, err := client.ShowVotes(context.Background(), &pb.ShowVotesRequest{Id: id})
	if err != nil {
		log.Println("failed to show votes.", err)
	}
	println(recv.Message)
}

func newGame(client pb.PlanningPokerClient) {
	recv, err := client.NewGame(context.Background(), &pb.NewGameRequest{})
	if err != nil {
		log.Println("failed to new game.", err)
	}
	println(recv.Message)
}

func printlnBroadcastMessage(message *pb.ConnectResponse) {
	switch message.Type {
	case pb.MessageType_JOIN:
		println(color.HiGreenString(message.Message + " joined"))
	case pb.MessageType_VOTE:
		println(color.HiGreenString(message.Message + " voted"))
	case pb.MessageType_SHOW_VOTES:
		// message.MessageがJSON形式の文字列になっているので、パースして表示する
		var votes map[string]float32
		err := json.Unmarshal([]byte(message.Message), &votes)
		if err != nil {
			log.Println("failed to unmarshal message.", err)
		}
		println(color.HiGreenString("result"))

		avg := votes[AVERAGE]
		delete(votes, AVERAGE)
		for k, v := range votes {
			println(color.HiGreenString(fmt.Sprintf("%s: %.2f", k, v)))
		}
		if len(votes) > 0 {
			println(color.HiGreenString(fmt.Sprintf("average: %.2f", avg)))
		}
	case pb.MessageType_LEAVE:
		println(color.HiGreenString(message.Message + " left"))
	}
}
