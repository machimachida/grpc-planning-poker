/**
 * @fileoverview gRPC-Web generated client stub for planning_poker
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.4.2
// 	protoc              v4.23.3
// source: proto/planning_poker.proto


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as proto_planning_poker_pb from './planning_poker_pb';


export class PlanningPokerClient {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: any; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; }) {
    if (!options) options = {};
    if (!credentials) credentials = {};
    options['format'] = 'text';

    this.client_ = new grpcWeb.GrpcWebClientBase(options);
    this.hostname_ = hostname.replace(/\/+$/, '');
    this.credentials_ = credentials;
    this.options_ = options;
  }

  methodDescriptorCreateRoom = new grpcWeb.MethodDescriptor(
    '/planning_poker.PlanningPoker/CreateRoom',
    grpcWeb.MethodType.SERVER_STREAMING,
    proto_planning_poker_pb.CreateRoomRequest,
    proto_planning_poker_pb.ConnectResponse,
    (request: proto_planning_poker_pb.CreateRoomRequest) => {
      return request.serializeBinary();
    },
    proto_planning_poker_pb.ConnectResponse.deserializeBinary
  );

  createRoom(
    request: proto_planning_poker_pb.CreateRoomRequest,
    metadata?: grpcWeb.Metadata): grpcWeb.ClientReadableStream<proto_planning_poker_pb.ConnectResponse> {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/planning_poker.PlanningPoker/CreateRoom',
      request,
      metadata || {},
      this.methodDescriptorCreateRoom);
  }

  methodDescriptorConnect = new grpcWeb.MethodDescriptor(
    '/planning_poker.PlanningPoker/Connect',
    grpcWeb.MethodType.SERVER_STREAMING,
    proto_planning_poker_pb.ConnectRequest,
    proto_planning_poker_pb.ConnectResponse,
    (request: proto_planning_poker_pb.ConnectRequest) => {
      return request.serializeBinary();
    },
    proto_planning_poker_pb.ConnectResponse.deserializeBinary
  );

  connect(
    request: proto_planning_poker_pb.ConnectRequest,
    metadata?: grpcWeb.Metadata): grpcWeb.ClientReadableStream<proto_planning_poker_pb.ConnectResponse> {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/planning_poker.PlanningPoker/Connect',
      request,
      metadata || {},
      this.methodDescriptorConnect);
  }

  methodDescriptorVote = new grpcWeb.MethodDescriptor(
    '/planning_poker.PlanningPoker/Vote',
    grpcWeb.MethodType.UNARY,
    proto_planning_poker_pb.VoteRequest,
    proto_planning_poker_pb.VoteResponse,
    (request: proto_planning_poker_pb.VoteRequest) => {
      return request.serializeBinary();
    },
    proto_planning_poker_pb.VoteResponse.deserializeBinary
  );

  vote(
    request: proto_planning_poker_pb.VoteRequest,
    metadata: grpcWeb.Metadata | null): Promise<proto_planning_poker_pb.VoteResponse>;

  vote(
    request: proto_planning_poker_pb.VoteRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: proto_planning_poker_pb.VoteResponse) => void): grpcWeb.ClientReadableStream<proto_planning_poker_pb.VoteResponse>;

  vote(
    request: proto_planning_poker_pb.VoteRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: proto_planning_poker_pb.VoteResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/planning_poker.PlanningPoker/Vote',
        request,
        metadata || {},
        this.methodDescriptorVote,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/planning_poker.PlanningPoker/Vote',
    request,
    metadata || {},
    this.methodDescriptorVote);
  }

  methodDescriptorShowVotes = new grpcWeb.MethodDescriptor(
    '/planning_poker.PlanningPoker/ShowVotes',
    grpcWeb.MethodType.UNARY,
    proto_planning_poker_pb.ShowVotesRequest,
    proto_planning_poker_pb.ShowVotesResponse,
    (request: proto_planning_poker_pb.ShowVotesRequest) => {
      return request.serializeBinary();
    },
    proto_planning_poker_pb.ShowVotesResponse.deserializeBinary
  );

  showVotes(
    request: proto_planning_poker_pb.ShowVotesRequest,
    metadata: grpcWeb.Metadata | null): Promise<proto_planning_poker_pb.ShowVotesResponse>;

  showVotes(
    request: proto_planning_poker_pb.ShowVotesRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: proto_planning_poker_pb.ShowVotesResponse) => void): grpcWeb.ClientReadableStream<proto_planning_poker_pb.ShowVotesResponse>;

  showVotes(
    request: proto_planning_poker_pb.ShowVotesRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: proto_planning_poker_pb.ShowVotesResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/planning_poker.PlanningPoker/ShowVotes',
        request,
        metadata || {},
        this.methodDescriptorShowVotes,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/planning_poker.PlanningPoker/ShowVotes',
    request,
    metadata || {},
    this.methodDescriptorShowVotes);
  }

  methodDescriptorNewGame = new grpcWeb.MethodDescriptor(
    '/planning_poker.PlanningPoker/NewGame',
    grpcWeb.MethodType.UNARY,
    proto_planning_poker_pb.NewGameRequest,
    proto_planning_poker_pb.NewGameResponse,
    (request: proto_planning_poker_pb.NewGameRequest) => {
      return request.serializeBinary();
    },
    proto_planning_poker_pb.NewGameResponse.deserializeBinary
  );

  newGame(
    request: proto_planning_poker_pb.NewGameRequest,
    metadata: grpcWeb.Metadata | null): Promise<proto_planning_poker_pb.NewGameResponse>;

  newGame(
    request: proto_planning_poker_pb.NewGameRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: proto_planning_poker_pb.NewGameResponse) => void): grpcWeb.ClientReadableStream<proto_planning_poker_pb.NewGameResponse>;

  newGame(
    request: proto_planning_poker_pb.NewGameRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: proto_planning_poker_pb.NewGameResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/planning_poker.PlanningPoker/NewGame',
        request,
        metadata || {},
        this.methodDescriptorNewGame,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/planning_poker.PlanningPoker/NewGame',
    request,
    metadata || {},
    this.methodDescriptorNewGame);
  }

}

