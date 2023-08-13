package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/machimachida/grpc-planning-poker/pb"
)

const (
	AVERAGE = "average"
)

var (
	cm      = ConnectionMap{streams: make(map[string]pb.PlanningPoker_ConnectServer)}
	voteMap = sync.Map{}
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPlanningPokerServer(s, &Server{})
	log.Println("Server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

type Server struct {
	pb.UnimplementedPlanningPokerServer
}

func (*Server) Connect(req *pb.ConnectRequest, stream pb.PlanningPoker_ConnectServer) error {
	log.Println("Connect function was invoked with a streaming request from " + req.Id)

	if req.Id == AVERAGE {
		return errors.New("this name is reserved")
	}

	err := cm.Connect(stream, req.Id)
	if err != nil {
		return err
	}
	cm.Broadcast(req.Id, pb.MessageType_JOIN)

	for {
		select {
		case <-stream.Context().Done():
			log.Println(req.Id + " is disconnected")
			if err := stream.Context().Err(); err != nil {
				log.Println(req.Id, err)
			}
			voteMap.Delete(req.Id)
			cm.Disconnect(req.Id)
			cm.Broadcast(req.Id, pb.MessageType_LEAVE)
			return nil
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (*Server) Vote(_ context.Context, req *pb.VoteRequest) (*pb.VoteResponse, error) {
	log.Printf("Vote function was invoked with a request from %s with vote \"%d\"\n", req.Id, req.Vote)

	cm.Broadcast(req.Id, pb.MessageType_VOTE)

	voteMap.Store(req.Id, req.Vote)
	return &pb.VoteResponse{Message: "voted"}, nil
}

func (*Server) ShowVotes(_ context.Context, req *pb.ShowVotesRequest) (*pb.ShowVotesResponse, error) {
	log.Println("ShowVotes function was invoked with a request from " + req.Id)

	votes := make(map[string]float32, len(cm.streams)+1)
	var sum int32
	var l int
	voteMap.Range(func(key, value any) bool {
		id, ok := key.(string)
		if !ok {
			return true
		}
		vote, ok := value.(int32)
		if !ok {
			return true
		}
		votes[id] = float32(vote)
		sum += vote
		l++
		return true
	})

	if sum == 0 {
		return &pb.ShowVotesResponse{Message: "no votes"}, nil
	}

	votes[AVERAGE] = float32(sum) / float32(l)
	b, err := json.Marshal(votes)
	if err != nil {
		log.Println("failed to marshal votes.", err)
	}
	cm.Broadcast(string(b), pb.MessageType_SHOW_VOTES)

	return &pb.ShowVotesResponse{Message: "accepted"}, nil
}

func (*Server) NewGame(_ context.Context, req *pb.NewGameRequest) (*pb.NewGameResponse, error) {
	log.Println("NewGame function was invoked with a request from " + req.Id)

	voteMap = sync.Map{}
	cm.Broadcast("new game start", pb.MessageType_NEW_GAME)

	return &pb.NewGameResponse{Message: "accepted"}, nil
}

// ConnectionMap Connectionの状態を保持する構造体。
// この構造体は、Connect関数で生成され、Disconnect関数で削除される。
// streamsは、クライアントのIDをキーとして、クライアントとの接続を保持する。
type ConnectionMap struct {
	mu      sync.Mutex
	streams map[string]pb.PlanningPoker_ConnectServer
}

var ErrAlreadyConnected = errors.New("already connected")

func (cm *ConnectionMap) Connect(stream pb.PlanningPoker_ConnectServer, name string) error {
	cm.mu.Lock()
	if _, ok := cm.streams[name]; ok {
		cm.mu.Unlock()
		return ErrAlreadyConnected
	}
	cm.streams[name] = stream
	cm.mu.Unlock()
	return nil
}

func (cm *ConnectionMap) Disconnect(name string) {
	cm.mu.Lock()
	cm.streams[name].Context().Done()
	delete(cm.streams, name)
	cm.mu.Unlock()
}

func (cm *ConnectionMap) Broadcast(message string, mt pb.MessageType) {
	cm.mu.Lock()
	for id, stream := range cm.streams {
		err := stream.Send(&pb.ConnectResponse{Message: message, Type: mt})
		if err != nil {
			log.Println("failed to send message to "+id, err)
		}
	}
	cm.mu.Unlock()
}
