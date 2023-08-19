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

var (
	roomId string
)

func main() {
	name := flag.String("name", "SatoTaro", "name of user")
	waitSecond := flag.Int("wait", 600, "wait second")
	isCreatingRoom := flag.Bool("create", false, "create room")
	joinRoomId := flag.String("join", "", "join room id")
	flag.Parse()

	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	var stream pb.PlanningPoker_ConnectClient
	var cancel context.CancelFunc
	if *isCreatingRoom {
		fmt.Print("Please input room name: ")
		var roomName string
		for {
			_, err = fmt.Scan(&roomName)
			if err == nil {
				break
			}
			log.Println("failed to scan input.", err)
			println("Please input room name.")
		}
		stream, cancel = createRoom(client, *name, roomName)
	} else {
		stream, cancel = connect(client, *name, *joinRoomId)
		roomId = *joinRoomId
	}
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
		case -2:
			println("Disconnect...")
			disconnect(cancel)
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

func listenServerMessage(stream pb.PlanningPoker_ConnectClient) {
	for {
		message, err := stream.Recv()

		// 異常系
		if err != nil {
			// TODO: ルームが見つからなかった場合やルームが作成された場合もここへ遷移する(Context.Errなし)。要対応

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

func createRoom(client pb.PlanningPokerClient, name, roomName string) (pb.PlanningPoker_ConnectClient, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := client.CreateRoom(ctx, &pb.CreateRoomRequest{
		Id:       name,
		RoomName: roomName,
	})
	if err != nil {
		panic(err)
	}
	return stream, cancel
}

func connect(client pb.PlanningPokerClient, id, roomId string) (pb.PlanningPoker_ConnectClient, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := client.Connect(ctx, &pb.ConnectRequest{Id: id, RoomId: roomId})
	if err != nil {
		panic(err)
	}
	return stream, cancel
}

func disconnect(cancel context.CancelFunc) {
	cancel()
}

func vote(client pb.PlanningPokerClient, id, roomId string, vote int32) {
	recv, err := client.Vote(context.Background(), &pb.VoteRequest{Id: id, Vote: vote, RoomId: roomId})
	if err != nil {
		log.Println("failed to vote.", err)
		return
	}
	if recv == nil {
		log.Println("failed to vote. recv is nil.")
		return
	}
	println(recv.Message)
}

func showVotes(client pb.PlanningPokerClient, id, roomId string) {
	recv, err := client.ShowVotes(context.Background(), &pb.ShowVotesRequest{Id: id, RoomId: roomId})
	if err != nil {
		log.Println("failed to show votes.", err)
		return
	}
	if recv == nil {
		log.Println("failed to show vote. recv is nil.")
		return
	}
	println(recv.Message)
}

func newGame(client pb.PlanningPokerClient, id, roomId string) {
	recv, err := client.NewGame(context.Background(), &pb.NewGameRequest{Id: id, RoomId: roomId})
	if err != nil {
		log.Println("failed to new game.", err)
	}
	if recv == nil {
		log.Println("failed to show vote. recv is nil.")
		return
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
	case pb.MessageType_CREATE_ROOM:
		println(color.CyanString("room " + message.Message + " created"))
		roomId = message.Message
	case pb.MessageType_STATUS:
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
	case pb.MessageType_NEW_GAME:
		println(color.YellowString(message.Message))
	case pb.MessageType_RESET_VOTE:
		println(color.HiGreenString(message.Message + " reset their vote"))
	}
}
