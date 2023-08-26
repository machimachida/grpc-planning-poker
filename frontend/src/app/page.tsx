'use client';
import AsyncLock from 'async-lock';
import {ClientReadableStream} from 'grpc-web';
import {useState} from 'react';
import {useForm} from 'react-hook-form';

import {ConnectRequest, ConnectResponse, CreateRoomRequest, MessageType} from "@/proto/planning_poker_pb";
import {PlanningPokerClient} from '@/proto/Planning_pokerServiceClientPb';

const client = new PlanningPokerClient("http://localhost:51000");

type StartNewGameConfig = {
  room: string; // room名かroomIDを入れる(暫定的)
  userName: string;
}

type Player = {
  name: string;
  isVoted: boolean;
  vote: number | null;
}

export default function Home() {
  let lock = new AsyncLock();
  const [stream, setStream] = useState<ClientReadableStream<ConnectResponse> | null>(null);
  const [response, setResponse] = useState<ConnectResponse | null>(null);
  const [players, setPlayers] = useState<Map<string, Player>>(new Map());
  const [roomId, setRoomId] = useState<string>("");

  const createNewRoom = (config: StartNewGameConfig) => {
    const req = new CreateRoomRequest
    req.setId(config.userName);
    req.setRoomname(config.room);

    if(stream !== null) {
      // TODO: 既存のストリームを終了させる処理をおこなう。新規部屋作成処理は続行する。
      console.error("stream is already exist");
      return;
    }

    const connection = client.createRoom(req);
    connection.on("data", async (res: ConnectResponse) => {
      await receiveMessage(res);
      setResponse(res);
    });
  };

  const joinRoom = (config: StartNewGameConfig) => {
    const req = new ConnectRequest();
    req.setId(config.userName);
    req.setRoomid(config.room);

    if(stream !== null) {
      // TODO: 既存のストリームを終了させる処理をおこなう。部屋参加処理は続行する。
      console.error("stream is already exist");
      return;
    }

    const connection = client.connect(req);
    connection.on("data", async (res: ConnectResponse) => {
      await receiveMessage(res);
      setResponse(res);
    });
  };

  const receiveMessage = async (res: ConnectResponse) => {
    console.log("receive message", res);
    await lock.acquire("receiveMessage", () => {
      switch (res.getType()) {
        case MessageType.JOIN:
          players.set(res.getMessage(), {name: res.getMessage(), isVoted: false, vote: null});
          setPlayers(players);
          break;
        case MessageType.VOTE:
          const player = players.get(res.getId())
          if (player) {
            player.isVoted = true;
            setPlayers(players);
          } else {
            players.set(res.getId(), {name: res.getId(), isVoted: true, vote: null});
            setPlayers(players);
          }
          break;
        case MessageType.SHOW_VOTES:
          const votes: Map<string, number> = JSON.parse(res.getMessage());
          for (const [key, value] of Object.entries(votes)) {
            const player = players.get(key);
            if (player) {
              player.vote = value;
              setPlayers(players);
            }
          }
          break;
        case MessageType.LEAVE:
          players.delete(res.getId());
          setPlayers(players);
          break;
        case MessageType.NEW_GAME:
          for(let k in players.keys()) {
            players.set(k, {name: k, isVoted: false, vote: null})
          }
          setPlayers(players);
          break;
        case MessageType.CREATE_ROOM:
          setRoomId(res.getMessage());
          break;
        case MessageType.STATUS:
          const playerStatuses: Map<string, boolean> = JSON.parse(res.getMessage());
          const initPlayers = new Map<string, Player>();
          console.log(playerStatuses);
          for (const [key, value] of Object.entries(playerStatuses)) {
            initPlayers.set(key, {name: key, isVoted: value, vote: null});
            setPlayers(initPlayers);
          }
          break;
        case MessageType.RESET_VOTE:
          for(let k in players.keys()) {
            players.set(k, {name: k, isVoted: false, vote: null})
          }
          setPlayers(players);
            break;
      }
      console.log(players);
    });
  }


  return (
    <main className="flex min-h-screen items-center justify-between p-24">
      <div className='w-1/3'>
        <NewGame createNewRoom={createNewRoom} joinRoom={joinRoom} />
      </div>
      <div className="z-10 max-w-5xl w-full items-center justify-between font-mono text-sm">
        <div>
          <p>roomID: {roomId}</p>
        </div>
        {response && (
          <>
            <p>レスポンス情報</p>
            <p>{response.getId()}</p>
            <p>{response.getType().toString()}</p>
            <p>{response.getMessage()}</p>
          </>
        )}
        <div>
          <p>プレイヤー情報</p>
          {
            Array.from(players.values()).map((player) => {
              return (
                <div key={player.name}>
                  <p>{player.name}</p>
                  <p>{player.isVoted ? "投票済み" : "未投票"}</p>
                  <p>{player.vote}</p>
                </div>
              );
            })
          }
        </div>
      </div>
    </main>
  )
}

/* ゲーム管理用コンポーネント */

type NewGameProps = {
  createNewRoom: (config: StartNewGameConfig) => void;
  joinRoom: (config: StartNewGameConfig) => void;
}

const NewGame: React.FC<NewGameProps> = ({ createNewRoom, joinRoom }) => {
  return (
    <div className="flex min-h-full flex-1 flex-col px-6 py-12 lg:px-8">
      <NewGameCard title='新規ルームを作成' form={<NewRoomForm connectRoom={createNewRoom} />} />
      <NewGameCard title='既存のルームへ参加' form={<JoinRoomForm connectRoom={joinRoom} />} />
    </div>
  );
}

type NewRoomFormData = {
  roomName: string;
  userName: string;
}

type NewRoomFormProps = {
  connectRoom: (config: StartNewGameConfig) => void;
}

const NewRoomForm: React.FC<NewRoomFormProps> = ({ connectRoom }) => {
  const { register, handleSubmit } = useForm<NewRoomFormData>();

  const onSubmit = (data: NewRoomFormData) => {
    connectRoom({room: data.roomName, userName: data.userName})
  };

  return (
    <form className="space-y-6" action="#" method="POST" onSubmit={handleSubmit(onSubmit)}>
      <div>
        <label htmlFor='roomName' className="block text-sm font-medium leading-6 text-gray-900">
          ルーム名
        </label>
        <div className="mt-2">
          <input
            id='roomName'
            required
            className="block w-full rounded-md border-0 px-3 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
            {...register('roomName')} />
        </div>
      </div>
      <div>
        <label htmlFor='userName' className="block text-sm font-medium leading-6 text-gray-900">
          ユーザ名
        </label>
        <div className="mt-2">
          <input
            id="userName"
            required
            className="block w-full rounded-md border-0 px-3 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
            {...register('userName')}
          />
        </div>
      </div>
      <button
        type='submit'
        className="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
      >
        ルームを作成
      </button>
    </form>
  );
}

type JoinRoomFormData = {
  roomId: string;
  userName: string;
}

type JoinRoomFormProps = {
  connectRoom: (config: StartNewGameConfig) => void;
}

const JoinRoomForm: React.FC<JoinRoomFormProps> = ({ connectRoom }) => {
  const { register, handleSubmit } = useForm<JoinRoomFormData>();

  const onSubmit = (data: JoinRoomFormData) => {
    connectRoom({room: data.roomId, userName: data.userName});
  };

  return (
    <form className="space-y-6" action="#" method="POST" onSubmit={handleSubmit(onSubmit)}>
      <div>
        <label htmlFor='roomId' className="block text-sm font-medium leading-6 text-gray-900">
          ルームID
        </label>
        <div className="mt-2">
          <input
            id='roomId'
            required
            className="block w-full rounded-md border-0 px-3 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
            {...register('roomId')} />
        </div>
      </div>
      <div>
        <label htmlFor='userName' className="block text-sm font-medium leading-6 text-gray-900">
          ユーザ名
        </label>
        <div className="mt-2">
          <input
            id="userName"
            required
            className="block w-full rounded-md border-0 px-3 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
            {...register('userName')}
          />
        </div>
      </div>
      <button
        type='submit'
        className="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
      >
        ルームへ参加
      </button>
    </form>
  );
}

type NewGameCardProps = {
  title: string;
  form: React.ReactNode;
}

const NewGameCard: React.FC<NewGameCardProps> = ({title, form}) => {
  return (
    <>
      <div className="sm:mx-auto sm:w-full sm:max-w-sm">
        <h2 className="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">
          {title}
        </h2>
      </div>

      <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
        {form}
      </div>
    </>
  );
}
