// @generated by protoc-gen-connect-es v1.0.0 with parameter "target=ts"
// @generated from file proto/v1/planning_poker.proto (package proto.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { ConnectRequest, ConnectResponse, CreateRoomRequest, NewGameRequest, NewGameResponse, ShowVotesRequest, ShowVotesResponse, VoteRequest, VoteResponse } from "./planning_poker_pb.ts";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service proto.v1.PlanningPokerService
 */
export const PlanningPokerService = {
  typeName: "proto.v1.PlanningPokerService",
  methods: {
    /**
     * @generated from rpc proto.v1.PlanningPokerService.CreateRoom
     */
    createRoom: {
      name: "CreateRoom",
      I: CreateRoomRequest,
      O: ConnectResponse,
      kind: MethodKind.ServerStreaming,
    },
    /**
     * @generated from rpc proto.v1.PlanningPokerService.Connect
     */
    connect: {
      name: "Connect",
      I: ConnectRequest,
      O: ConnectResponse,
      kind: MethodKind.ServerStreaming,
    },
    /**
     * @generated from rpc proto.v1.PlanningPokerService.Vote
     */
    vote: {
      name: "Vote",
      I: VoteRequest,
      O: VoteResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc proto.v1.PlanningPokerService.ShowVotes
     */
    showVotes: {
      name: "ShowVotes",
      I: ShowVotesRequest,
      O: ShowVotesResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc proto.v1.PlanningPokerService.NewGame
     */
    newGame: {
      name: "NewGame",
      I: NewGameRequest,
      O: NewGameResponse,
      kind: MethodKind.Unary,
    },
  }
} as const;

