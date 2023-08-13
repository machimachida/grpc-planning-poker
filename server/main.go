package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
	rm = RoomMap{rooms: make(map[string]*Room)}
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

var (
	ErrReservedUserName = errors.New("this name is reserved")
	ErrExistRoom        = errors.New("this room is already exist")
)

func (*Server) CreateRoom(req *pb.CreateRoomRequest, stream pb.PlanningPoker_CreateRoomServer) error {
	log.Println("CreateRoom function was invoked with a request from " + req.Id)

	if req.Id == AVERAGE {
		return ErrReservedUserName
	}

	rm.mu.Lock()
	for room := range rm.rooms {
		if room == req.Id {
			rm.mu.Unlock()
			return ErrExistRoom
		}
	}

	id := uuid.New().String()
	rm.rooms[id] = &Room{
		id:          id,
		name:        req.Id,
		connections: ConnectionMap{streams: make(map[string]pb.PlanningPoker_ConnectServer, 1)},
		voteMap:     &sync.Map{},
	}
	rm.mu.Unlock()

	log.Println("room created: " + id)
	err := stream.Send(&pb.ConnectResponse{Message: id, Type: pb.MessageType_CREATE_ROOM})
	if err != nil {
		log.Println("failed to send message.", err)
	}

	err = connectWithRoom(stream, id, req.Id)
	return err
}

func (*Server) Connect(req *pb.ConnectRequest, stream pb.PlanningPoker_ConnectServer) error {
	log.Println("Connect function was invoked with a streaming request from " + req.Id)

	if req.Id == AVERAGE {
		return ErrReservedUserName
	}

	err := connectWithRoom(stream, req.RoomId, req.Id)
	return err
}

func connectWithRoom(stream pb.PlanningPoker_ConnectServer, roomId, name string) error {
	r, ok := rm.rooms[roomId]
	if !ok {
		msg := fmt.Sprintf("room %s not found", roomId)
		err := errors.New(msg)
		log.Println(err)
		return err
	}

	err := r.connections.Connect(stream, name)
	if err != nil {
		log.Println("failed to connect", err)
		return err
	}

	// クライアントがルームに参加した際の、他ユーザの接続状況を通知する
	userVoteStatus := make(map[string]bool, len(r.connections.streams))
	for id := range r.connections.streams {
		if _, ok := r.voteMap.Load(id); ok {
			userVoteStatus[id] = true
		} else {
			userVoteStatus[id] = false
		}
	}
	b, err := json.Marshal(userVoteStatus)
	if err != nil {
		log.Println("failed to marshal user vote status.", err)
	} else {
		err := stream.Send(&pb.ConnectResponse{Message: string(b), Type: pb.MessageType_STATUS})
		if err != nil {
			log.Println("failed to send message.", err)
		}
	}

	// 参加したことを全ユーザに通知する
	r.connections.Broadcast(name, pb.MessageType_JOIN)

	for {
		select {
		case <-stream.Context().Done():
			log.Println(name + " is disconnected from " + roomId)
			if err := stream.Context().Err(); err != nil {
				log.Println(name, err)
			}
			r.connections.Disconnect(name)
			r.connections.Broadcast(name, pb.MessageType_LEAVE)

			// 参加者がいなくなったらルームを削除する
			if len(r.connections.streams) == 0 {
				rm.mu.Lock()
				log.Println("room " + roomId + " is closed")
				delete(rm.rooms, roomId)
				rm.mu.Unlock()
			}

			return nil
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (*Server) Vote(_ context.Context, req *pb.VoteRequest) (*pb.VoteResponse, error) {
	log.Printf("Vote function was invoked with a request from %s with vote \"%d\"\n", req.Id, req.Vote)

	r, ok := rm.rooms[req.RoomId]
	if !ok {
		err := fmt.Errorf("room %s not found", req.RoomId)
		log.Println(err)
		return nil, err
	}

	if req.Vote == -1 {
		r.connections.Broadcast(req.Id, pb.MessageType_RESET_VOTE)
		r.voteMap.Delete(req.Id)
	} else if req.Vote <= 0 {
		err := fmt.Errorf("invalid vote %d", req.Vote)
		log.Println(err)
		return &pb.VoteResponse{Message: "invalid vote"}, err
	} else {
		r.connections.Broadcast(req.Id, pb.MessageType_VOTE)
		r.voteMap.Store(req.Id, req.Vote)
	}

	return &pb.VoteResponse{Message: "voted"}, nil
}

func (*Server) ShowVotes(_ context.Context, req *pb.ShowVotesRequest) (*pb.ShowVotesResponse, error) {
	log.Println("ShowVotes function was invoked with a request from " + req.Id)

	r, ok := rm.rooms[req.RoomId]
	if !ok {
		err := fmt.Errorf("room %s not found", req.RoomId)
		log.Println(err)
		return nil, err
	}

	votes := make(map[string]float32, len(r.connections.streams)+1)
	var sum int32
	var l int
	r.voteMap.Range(func(key, value any) bool {
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
	r.connections.Broadcast(string(b), pb.MessageType_SHOW_VOTES)

	return &pb.ShowVotesResponse{Message: "accepted"}, nil
}

func (*Server) NewGame(_ context.Context, req *pb.NewGameRequest) (*pb.NewGameResponse, error) {
	log.Println("NewGame function was invoked with a request from " + req.Id)

	r, ok := rm.rooms[req.RoomId]
	if !ok {
		err := fmt.Errorf("room %s not found", req.RoomId)
		log.Println(err)
		return nil, err
	}

	var m sync.Map
	r.voteMap = &m
	log.Println("new game start in Room " + req.RoomId)
	r.connections.Broadcast("new game start", pb.MessageType_NEW_GAME)

	return &pb.NewGameResponse{Message: "accepted"}, nil
}

type RoomMap struct {
	mu    sync.Mutex
	rooms map[string]*Room
}

type Room struct {
	id          string
	name        string
	connections ConnectionMap
	voteMap     *sync.Map
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
