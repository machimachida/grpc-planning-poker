package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/fatih/color"

	pokerv1 "github.com/machimachida/grpc-planning-poker/gen/proto/v1"
	"github.com/machimachida/grpc-planning-poker/gen/proto/v1/pokerv1connect"
)

const (
	AVERAGE = "average"
)

var (
	roomId string
)

func main() {
	name := flag.String("name", "Taro", "name of user")
	waitSecond := flag.Int("wait", 600, "wait second")
	isCreatingRoom := flag.Bool("create", false, "create room")
	joinRoomId := flag.String("join", "", "join room id")
	flag.Parse()

	client := pokerv1connect.NewPlanningPokerServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)

	ctx, cancel := context.WithCancel(context.Background())
	var (
		stream *connect.ServerStreamForClient[pokerv1.ConnectResponse]
		err    error
	)
	if *isCreatingRoom {
		fmt.Print("Please input room name: ")
		var roomName string
		for {
			_, err := fmt.Scan(&roomName)
			if err == nil {
				break
			}
			log.Println("failed to scan input.", err)
			println("Please input room name.")
		}
		stream, err = client.CreateRoom(ctx, connect.NewRequest(&pokerv1.CreateRoomRequest{
			Id:       *name,
			RoomName: roomName,
		}))
	} else {
		stream, err = client.Connect(ctx, connect.NewRequest(&pokerv1.ConnectRequest{
			Id:     *name,
			RoomId: *joinRoomId,
		}))
		roomId = *joinRoomId
	}
	if err != nil {
		log.Fatal("failed to create or join room.", err)
	}
	go listenServerMessage(stream)
	go disconnectAfterWaitSecond(cancel, *waitSecond)

	println("Start planning poker!")

	for {
		var n int32
		_, err := fmt.Scan(&n)
		if err != nil {
			log.Println("failed to scan input. Please input integer.", err)
			continue
		}

		switch n {
		case -2:
			println("Disconnect...")
			cancel()
		case -3:
			showVotes(client, *name, roomId)
			continue
		case -4:
			newGame(client, *name, roomId)
			continue
		default:
			if n > 0 || n == -1 {
				vote(client, *name, roomId, n)
				continue
			}

			println("Please input natural number, -1(reset your vote), -2(disconnect), -3(show votes) or -4(new game).")
		}
	}
}

func listenServerMessage(stream *connect.ServerStreamForClient[pokerv1.ConnectResponse]) {
	for {
		if stream.Receive() {
			message := stream.Msg()
			if message == nil {
				log.Println("server sent nil message.")
				return
			}
			printlnBroadcastMessage(message)
		}

		err := stream.Err()
		if err == nil {
			continue
		}

		if errors.Is(err, context.Canceled) {
			log.Println("server is disconnected.")
			os.Exit(0)
		}

		log.Println("failed to receive message.", err)
		os.Exit(1)
	}
}

func disconnectAfterWaitSecond(cancel context.CancelFunc, waitSecond int) {
	time.Sleep(time.Duration(waitSecond) * time.Second)
	cancel()
	os.Exit(0)
}

func vote(client pokerv1connect.PlanningPokerServiceClient, id, roomId string, vote int32) {
	res, err := client.Vote(context.Background(), connect.NewRequest(&pokerv1.VoteRequest{Id: id, Vote: vote, RoomId: roomId}))
	if err != nil {
		log.Println("failed to vote.", err)
		return
	}
	if res == nil {
		log.Println("failed to vote. res is nil.")
		return
	}
	println(res.Msg.Message)
}

func showVotes(client pokerv1connect.PlanningPokerServiceClient, id, roomId string) {
	res, err := client.ShowVotes(context.Background(), connect.NewRequest(&pokerv1.ShowVotesRequest{Id: id, RoomId: roomId}))
	if err != nil {
		log.Println("failed to show votes.", err)
		return
	}
	if res == nil {
		log.Println("failed to show vote. res is nil.")
		return
	}
	println(res.Msg.Message)
}

func newGame(client pokerv1connect.PlanningPokerServiceClient, id, roomId string) {
	res, err := client.NewGame(context.Background(), connect.NewRequest(&pokerv1.NewGameRequest{Id: id, RoomId: roomId}))
	if err != nil {
		log.Println("failed to new game.", err)
	}
	if res == nil {
		log.Println("failed to show vote. res is nil.")
		return
	}
	println(res.Msg.Message)
}

func printlnBroadcastMessage(message *pokerv1.ConnectResponse) {
	switch message.Type {
	case pokerv1.MessageType_MESSAGE_TYPE_JOIN:
		println(color.HiGreenString(message.Message + " joined"))
	case pokerv1.MessageType_MESSAGE_TYPE_VOTE:
		println(color.HiGreenString(message.Message + " voted"))
	case pokerv1.MessageType_MESSAGE_TYPE_SHOW_VOTES:
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
	case pokerv1.MessageType_MESSAGE_TYPE_LEAVE:
		println(color.HiGreenString(message.Message + " left"))
	case pokerv1.MessageType_MESSAGE_TYPE_CREATE_ROOM:
		println(color.CyanString("room " + message.Message + " created"))
		roomId = message.Message
	case pokerv1.MessageType_MESSAGE_TYPE_STATUS:
		// ユーザの状態情報を持つmessage.MessageがJSON形式の文字列になっているので、パースして表示する
		var status map[string]bool
		err := json.Unmarshal([]byte(message.Message), &status)
		if err != nil {
			log.Println("failed to unmarshal message.", err)
		}
		println(color.CyanString("room status"))
		for k, v := range status {
			println(color.CyanString(fmt.Sprintf("%s: %t", k, v)))
		}
	case pokerv1.MessageType_MESSAGE_TYPE_NEW_GAME:
		println(color.YellowString(message.Message))
	case pokerv1.MessageType_MESSAGE_TYPE_RESET_VOTE:
		println(color.HiGreenString(message.Message + " reset their vote"))
	}
}
