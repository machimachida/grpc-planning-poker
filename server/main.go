package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/machimachida/grpc-planning-poker/pb"
	"google.golang.org/grpc"
)

var (
	cm      = ConnectionMap{streams: make(map[string]pb.PlanningPoker_ConnectServer)}
	voteMap = sync.Map{}
)

type Server struct {
	pb.UnimplementedPlanningPokerServer
}

func (*Server) Connect(req *pb.ConnectRequest, stream pb.PlanningPoker_ConnectServer) error {
	fmt.Println("Connect function was invoked with a streaming request from " + req.Id)

	err := cm.Connect(stream, req.Id)
	if err != nil {
		return err
	}
	cm.Broadcast(req.Id + " is connected")

	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (*Server) Disconnect(ctx context.Context, req *pb.DisconnectRequest) (*pb.DisconnectResponse, error) {
	fmt.Println("Disconnect function was invoked with a request from " + req.Id)

	cm.Broadcast(req.Id + " is disconnected")

	cm.Disconnect(req.Id)
	voteMap.Delete(req.Id)
	return &pb.DisconnectResponse{Message: "disconnected"}, nil
}

func (*Server) Vote(ctx context.Context, req *pb.VoteRequest) (*pb.VoteResponse, error) {
	fmt.Printf("Vote function was invoked with a request from %s with vote \"%d\"\n", req.Id, req.Vote)

	cm.Broadcast(req.Id + " is voted")

	voteMap.Store(req.Id, req.Vote)
	return &pb.VoteResponse{Message: "voted"}, nil
}

func (*Server) ShowVotes(ctx context.Context, req *pb.ShowVotesRequest) (*pb.ShowVotesResponse, error) {
	fmt.Println("ShowVotes function was invoked with a request from " + req.Id)

	cm.Broadcast("show votes by " + req.Id)

	var sum int32
	var size int
	voteMap.Range(func(key, value interface{}) bool {
		cm.Broadcast(fmt.Sprintf("%s: %d", key, value))
		sum += value.(int32)
		size++
		return true
	})
	cm.Broadcast(fmt.Sprintf("average: %f\n", float32(sum)/float32(size)))

	voteMap = sync.Map{}
	cm.Broadcast("new game start")

	return &pb.ShowVotesResponse{Message: "accepted"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPlanningPokerServer(s, &Server{})
	fmt.Println("Server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Connectionの状態を保持する構造体。
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

func (cm *ConnectionMap) Broadcast(message string) {
	cm.mu.Lock()
	for _, stream := range cm.streams {
		stream.Send(&pb.ConnectResponse{Message: message})
	}
	cm.mu.Unlock()
}
