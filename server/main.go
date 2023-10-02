package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/rs/cors"

	pokerv1 "github.com/machimachida/grpc-planning-poker/gen/proto/v1"
	"github.com/machimachida/grpc-planning-poker/gen/proto/v1/pokerv1connect"
)

const (
	AVERAGE = "average"
)

var (
	rm = RoomMap{rooms: make(map[string]*Room)}
)

var (
	ErrReservedUserName = errors.New("this name is reserved")
	ErrExistRoom        = errors.New("this room is already exist")
)

type pokerServer struct{}

func (s *pokerServer) CreateRoom(ctx context.Context, req *connect.Request[pokerv1.CreateRoomRequest], stream *connect.ServerStream[pokerv1.ConnectResponse]) error {
	log.Println("CreateRoom function was invoked with a request from " + req.Msg.Id)

	if req.Msg.Id == AVERAGE {
		return connect.NewError(
			connect.CodeAlreadyExists,
			ErrReservedUserName,
		)
	}

	rm.mu.Lock()
	for room := range rm.rooms {
		if room == req.Msg.Id {
			rm.mu.Unlock()
			return connect.NewError(
				connect.CodeAlreadyExists,
				ErrExistRoom,
			)
		}
	}

	id := req.Msg.RoomId
	rm.rooms[id] = &Room{
		id:          id,
		connections: ConnectionMap{streams: make(map[string]StreamState, 1)},
		voteMap:     &sync.Map{},
	}
	rm.mu.Unlock()

	log.Println("room created: " + id)

	err := stream.Send(&pokerv1.ConnectResponse{
		Id:      id,
		Type:    pokerv1.MessageType_MESSAGE_TYPE_CREATE_ROOM,
		Message: id,
	})
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	err = connectWithRoom(ctx, stream, id, req.Msg.Id)
	if err != nil {
		return connect.NewError(
			connect.CodeInternal,
			err,
		)
	}
	return nil
}

func (s *pokerServer) Connect(ctx context.Context, req *connect.Request[pokerv1.ConnectRequest], stream *connect.ServerStream[pokerv1.ConnectResponse]) error {
	log.Println("Connect function was invoked with a request from " + req.Msg.Id)

	if req.Msg.Id == AVERAGE {
		return connect.NewError(
			connect.CodeAlreadyExists,
			ErrReservedUserName,
		)
	}

	rm.mu.Lock()
	for room := range rm.rooms {
		if room == req.Msg.RoomId {
			rm.mu.Unlock()
			err := connectWithRoom(ctx, stream, room, req.Msg.Id)
			if err != nil {
				return connect.NewError(
					connect.CodeInternal,
					err,
				)
			}
			return nil
		}
	}
	rm.mu.Unlock()

	return connect.NewError(
		connect.CodeNotFound,
		errors.New("room not found"),
	)
}

func connectWithRoom(ctx context.Context, stream *connect.ServerStream[pokerv1.ConnectResponse], roomId, name string) error {
	r, ok := rm.rooms[roomId]
	if !ok {
		msg := fmt.Sprintf("room %s not found", roomId)
		err := errors.New(msg)
		log.Println(err)
		return err
	}

	err := r.connections.Connect(ctx, stream, name)
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
		err := stream.Send(&pokerv1.ConnectResponse{
			Message: string(b),
			Type:    pokerv1.MessageType_MESSAGE_TYPE_STATUS,
		})
		if err != nil {
			log.Println("failed to send message.", err)
		}
	}

	// 参加したことを全ユーザに通知する
	r.connections.Broadcast(name, pokerv1.MessageType_MESSAGE_TYPE_JOIN)
	r.currentUsedAt = time.Now()

	for {
		select {
		case <-ctx.Done():
			log.Println(name + " is disconnected from " + roomId)
			if err := ctx.Err(); err != nil {
				log.Println(name, err)
			}
			r.connections.Disconnect(name)
			r.connections.Broadcast(name, pokerv1.MessageType_MESSAGE_TYPE_LEAVE)

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

func (s *pokerServer) Vote(_ context.Context, req *connect.Request[pokerv1.VoteRequest]) (*connect.Response[pokerv1.VoteResponse], error) {
	log.Printf("Vote function was invoked with a request from %s with vote \"%d\"\n", req.Msg.Id, req.Msg.Vote)

	r, ok := rm.rooms[req.Msg.RoomId]
	if !ok {
		err := fmt.Errorf("room %s not found", req.Msg.RoomId)
		log.Println(err)
		return nil, connect.NewError(
			connect.CodeNotFound,
			err,
		)
	}

	if req.Msg.Vote == -1 {
		r.connections.Broadcast(req.Msg.Id, pokerv1.MessageType_MESSAGE_TYPE_RESET_VOTE)
		r.voteMap.Delete(req.Msg.Id)
	} else if req.Msg.Vote <= 0 {
		err := fmt.Errorf("invalid vote %d", req.Msg.Vote)
		log.Println(err)
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			err,
		)
	} else {
		r.connections.Broadcast(req.Msg.Id, pokerv1.MessageType_MESSAGE_TYPE_VOTE)
		r.voteMap.Store(req.Msg.Id, req.Msg.Vote)
	}
	r.currentUsedAt = time.Now()

	return connect.NewResponse(&pokerv1.VoteResponse{
		Message: "voted",
	}), nil
}

func (s *pokerServer) ShowVotes(_ context.Context, req *connect.Request[pokerv1.ShowVotesRequest]) (*connect.Response[pokerv1.ShowVotesResponse], error) {
	log.Println("ShowVotes function was invoked with a request from " + req.Msg.Id)

	r, ok := rm.rooms[req.Msg.RoomId]
	if !ok {
		err := fmt.Errorf("room %s not found", req.Msg.RoomId)
		log.Println(err)
		return nil, connect.NewError(
			connect.CodeNotFound,
			err,
		)
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
		return connect.NewResponse(&pokerv1.ShowVotesResponse{
			Message: "no votes",
		}), nil
	}

	votes[AVERAGE] = float32(sum) / float32(l)
	b, err := json.Marshal(votes)
	if err != nil {
		log.Println("failed to marshal votes.", err)
	}
	r.connections.Broadcast(string(b), pokerv1.MessageType_MESSAGE_TYPE_SHOW_VOTES)
	r.currentUsedAt = time.Now()

	return connect.NewResponse(&pokerv1.ShowVotesResponse{
		Message: "accepted",
	}), nil
}

func (s *pokerServer) NewGame(_ context.Context, req *connect.Request[pokerv1.NewGameRequest]) (*connect.Response[pokerv1.NewGameResponse], error) {
	log.Println("NewGame function was invoked with a request from " + req.Msg.Id)

	r, ok := rm.rooms[req.Msg.RoomId]
	if !ok {
		err := fmt.Errorf("room %s not found", req.Msg.RoomId)
		log.Println(err)
		return nil, connect.NewError(
			connect.CodeNotFound,
			err,
		)
	}

	var m sync.Map
	r.voteMap = &m
	log.Println("new game start in Room " + req.Msg.RoomId)
	r.connections.Broadcast("new game start", pokerv1.MessageType_MESSAGE_TYPE_NEW_GAME)
	r.currentUsedAt = time.Now()

	return connect.NewResponse(&pokerv1.NewGameResponse{
		Message: "accepted",
	}), nil
}

func main() {
	// 1時間に1回、使われていないルームがあるか確認する
	// 6時間使われていないルームは削除する
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			rm.mu.Lock()
			for id, r := range rm.rooms {
				if time.Since(r.currentUsedAt) > 6*time.Hour {
					log.Println("room " + id + " is closed because it is not used for a long time")
					delete(rm.rooms, id)
				}
			}
			rm.mu.Unlock()
		}
	}()

	corsHandler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{
			"Accept-Encoding",
			"Content-Encoding",
			"Content-Type",
			"Content-Type",
			"Connect-Protocol-Version",
			"Connect-Timeout-Ms",
			"Connect-Accept-Encoding",  // Unused in web browsers, but added for future-proofing
			"Connect-Content-Encoding", // Unused in web browsers, but added for future-proofing
			"Grpc-Timeout",             // Used for gRPC-web
			"X-Grpc-Web",               // Used for gRPC-web
			"X-User-Agent",             // Used for gRPC-web
		},
		ExposedHeaders: []string{
			"Content-Encoding",         // Unused in web browsers, but added for future-proofing
			"Connect-Content-Encoding", // Unused in web browsers, but added for future-proofing
			"Grpc-Status",              // Required for gRPC-web
			"Grpc-Message",             // Required for gRPC-web
		},
	})

	server := &pokerServer{}
	mux := http.NewServeMux()
	mux.Handle(pokerv1connect.NewPlanningPokerServiceHandler(server))
	handler := corsHandler.Handler(mux)

	log.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}

type RoomMap struct {
	mu    sync.Mutex
	rooms map[string]*Room
}

type Room struct {
	id            string
	connections   ConnectionMap
	voteMap       *sync.Map
	currentUsedAt time.Time
}

// ConnectionMap Connectionの状態を保持する構造体。
// この構造体は、Connect関数で生成され、Disconnect関数で削除される。
// streamsは、クライアントのIDをキーとして、クライアントとの接続を保持する。
type ConnectionMap struct {
	mu      sync.Mutex
	streams map[string]StreamState
}

type StreamState struct {
	ctx    context.Context
	stream *connect.ServerStream[pokerv1.ConnectResponse]
}

var ErrAlreadyConnected = errors.New("already connected")

func (cm *ConnectionMap) Connect(ctx context.Context, stream *connect.ServerStream[pokerv1.ConnectResponse], name string) error {
	cm.mu.Lock()
	if _, ok := cm.streams[name]; ok {
		cm.mu.Unlock()
		return ErrAlreadyConnected
	}
	cm.streams[name] = StreamState{
		ctx:    ctx,
		stream: stream,
	}
	cm.mu.Unlock()
	return nil
}

func (cm *ConnectionMap) Disconnect(name string) {
	cm.mu.Lock()
	cm.streams[name].ctx.Done()
	delete(cm.streams, name)
	cm.mu.Unlock()
}

func (cm *ConnectionMap) Broadcast(message string, mt pokerv1.MessageType) {
	cm.mu.Lock()
	for id, state := range cm.streams {
		err := state.stream.Send(&pokerv1.ConnectResponse{
			Type:    mt,
			Message: message,
		})
		if err != nil {
			log.Println("failed to send message to "+id, err)
		}
	}
	cm.mu.Unlock()
}
