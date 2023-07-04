package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/machimachida/grpc-planning-poker/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPlanningPokerServer
}

var cm = connectionMap{stream: make(map[string]pb.PlanningPoker_ConnectServer)}
var voteMap = sync.Map{}

func (*server) Connect(req *pb.ConnectRequest, stream pb.PlanningPoker_ConnectServer) error {
	fmt.Println("Connect function was invoked with a streaming request from " + req.Id)
	cm.connect(stream, req.Id)

	for {
		time.Sleep(1 * time.Second)
	}
}

func (*server) Disconnect(ctx context.Context, req *pb.DisconnectRequest) (*pb.DisconnectResponse, error) {
	fmt.Println("Disconnect function was invoked with a streaming request from " + req.Id)
	cm.broadcast(&pb.ConnectResponse{Message: req.Id + " is disconnected"})
	cm.disconnect(req.Id)
	voteMap.Delete(req.Id)
	return &pb.DisconnectResponse{Message: "disconnected"}, nil
}

func (*server) Vote(ctx context.Context, req *pb.VoteRequest) (*pb.VoteResponse, error) {
	fmt.Println(req.Vote)
	fmt.Println("Vote function was invoked with a streaming request from " + req.Id + " with " + string(req.Vote) + " vote")

	voteMap.Store(req.Id, req.Vote)

	cm.broadcast(&pb.ConnectResponse{Message: req.Id + " is voted"})
	return &pb.VoteResponse{Message: "voted"}, nil
}

func (*server) ShowVotes(ctx context.Context, req *pb.ShowVotesRequest) (*pb.ShowVotesResponse, error) {
	fmt.Println("ShowVotes function was invoked with a streaming request from " + req.Id)
	cm.broadcast(&pb.ConnectResponse{Message: "show votes"})

	var sum int32
	var size int
	voteMap.Range(func(key, value interface{}) bool {
		cm.broadcast(&pb.ConnectResponse{Message: fmt.Sprintf("%s: %d", key, value)})
		sum += value.(int32)
		size++
		return true
	})
	cm.broadcast(&pb.ConnectResponse{Message: fmt.Sprintf("average: %f", float32(sum)/float32(size))})

	voteMap = sync.Map{}
	cm.broadcast(&pb.ConnectResponse{Message: "new game start"})

	return &pb.ShowVotesResponse{Message: "accepted"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPlanningPokerServer(s, &server{})
	fmt.Println("Server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

type connectionMap struct {
	mu     sync.Mutex
	stream map[string]pb.PlanningPoker_ConnectServer
}

func (cm *connectionMap) connect(stream pb.PlanningPoker_ConnectServer, name string) {
	cm.mu.Lock()
	cm.stream[name] = stream
	cm.mu.Unlock()
}

func (cm *connectionMap) disconnect(name string) {
	cm.mu.Lock()
	delete(cm.stream, name)
	cm.mu.Unlock()
}

func (cm *connectionMap) broadcast(message *pb.ConnectResponse) {
	cm.mu.Lock()
	for _, stream := range cm.stream {
		stream.Send(message)
	}
	cm.mu.Unlock()
}

func (cm *connectionMap) send(message *pb.ConnectResponse, stream pb.PlanningPoker_ConnectServer) {
	cm.mu.Lock()
	stream.Send(message)
	cm.mu.Unlock()
}
