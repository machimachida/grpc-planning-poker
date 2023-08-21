import * as jspb from 'google-protobuf'



export class CreateRoomRequest extends jspb.Message {
  getId(): string;
  setId(value: string): CreateRoomRequest;

  getRoomname(): string;
  setRoomname(value: string): CreateRoomRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateRoomRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateRoomRequest): CreateRoomRequest.AsObject;
  static serializeBinaryToWriter(message: CreateRoomRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateRoomRequest;
  static deserializeBinaryFromReader(message: CreateRoomRequest, reader: jspb.BinaryReader): CreateRoomRequest;
}

export namespace CreateRoomRequest {
  export type AsObject = {
    id: string,
    roomname: string,
  }
}

export class ConnectRequest extends jspb.Message {
  getId(): string;
  setId(value: string): ConnectRequest;

  getRoomid(): string;
  setRoomid(value: string): ConnectRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConnectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ConnectRequest): ConnectRequest.AsObject;
  static serializeBinaryToWriter(message: ConnectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConnectRequest;
  static deserializeBinaryFromReader(message: ConnectRequest, reader: jspb.BinaryReader): ConnectRequest;
}

export namespace ConnectRequest {
  export type AsObject = {
    id: string,
    roomid: string,
  }
}

export class ConnectResponse extends jspb.Message {
  getId(): string;
  setId(value: string): ConnectResponse;

  getType(): MessageType;
  setType(value: MessageType): ConnectResponse;

  getMessage(): string;
  setMessage(value: string): ConnectResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConnectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ConnectResponse): ConnectResponse.AsObject;
  static serializeBinaryToWriter(message: ConnectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConnectResponse;
  static deserializeBinaryFromReader(message: ConnectResponse, reader: jspb.BinaryReader): ConnectResponse;
}

export namespace ConnectResponse {
  export type AsObject = {
    id: string,
    type: MessageType,
    message: string,
  }
}

export class VoteRequest extends jspb.Message {
  getId(): string;
  setId(value: string): VoteRequest;

  getVote(): number;
  setVote(value: number): VoteRequest;

  getRoomid(): string;
  setRoomid(value: string): VoteRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VoteRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VoteRequest): VoteRequest.AsObject;
  static serializeBinaryToWriter(message: VoteRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VoteRequest;
  static deserializeBinaryFromReader(message: VoteRequest, reader: jspb.BinaryReader): VoteRequest;
}

export namespace VoteRequest {
  export type AsObject = {
    id: string,
    vote: number,
    roomid: string,
  }
}

export class VoteResponse extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): VoteResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VoteResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VoteResponse): VoteResponse.AsObject;
  static serializeBinaryToWriter(message: VoteResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VoteResponse;
  static deserializeBinaryFromReader(message: VoteResponse, reader: jspb.BinaryReader): VoteResponse;
}

export namespace VoteResponse {
  export type AsObject = {
    message: string,
  }
}

export class ShowVotesRequest extends jspb.Message {
  getId(): string;
  setId(value: string): ShowVotesRequest;

  getRoomid(): string;
  setRoomid(value: string): ShowVotesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ShowVotesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ShowVotesRequest): ShowVotesRequest.AsObject;
  static serializeBinaryToWriter(message: ShowVotesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ShowVotesRequest;
  static deserializeBinaryFromReader(message: ShowVotesRequest, reader: jspb.BinaryReader): ShowVotesRequest;
}

export namespace ShowVotesRequest {
  export type AsObject = {
    id: string,
    roomid: string,
  }
}

export class ShowVotesResponse extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): ShowVotesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ShowVotesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ShowVotesResponse): ShowVotesResponse.AsObject;
  static serializeBinaryToWriter(message: ShowVotesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ShowVotesResponse;
  static deserializeBinaryFromReader(message: ShowVotesResponse, reader: jspb.BinaryReader): ShowVotesResponse;
}

export namespace ShowVotesResponse {
  export type AsObject = {
    message: string,
  }
}

export class NewGameRequest extends jspb.Message {
  getId(): string;
  setId(value: string): NewGameRequest;

  getRoomid(): string;
  setRoomid(value: string): NewGameRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NewGameRequest.AsObject;
  static toObject(includeInstance: boolean, msg: NewGameRequest): NewGameRequest.AsObject;
  static serializeBinaryToWriter(message: NewGameRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NewGameRequest;
  static deserializeBinaryFromReader(message: NewGameRequest, reader: jspb.BinaryReader): NewGameRequest;
}

export namespace NewGameRequest {
  export type AsObject = {
    id: string,
    roomid: string,
  }
}

export class NewGameResponse extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): NewGameResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NewGameResponse.AsObject;
  static toObject(includeInstance: boolean, msg: NewGameResponse): NewGameResponse.AsObject;
  static serializeBinaryToWriter(message: NewGameResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NewGameResponse;
  static deserializeBinaryFromReader(message: NewGameResponse, reader: jspb.BinaryReader): NewGameResponse;
}

export namespace NewGameResponse {
  export type AsObject = {
    message: string,
  }
}

export enum MessageType { 
  JOIN = 0,
  VOTE = 1,
  SHOW_VOTES = 2,
  LEAVE = 3,
  NEW_GAME = 4,
  CREATE_ROOM = 5,
  STATUS = 6,
  RESET_VOTE = 7,
}
