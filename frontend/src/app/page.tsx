'use client';
import AsyncLock from 'async-lock';
import { ClientReadableStream } from 'grpc-web';
import { useState } from 'react';

import JoinRoomForm from '@/app/components/JoinRoomForm';
import NewGameCard from '@/app/components/NewGameCard';
import NewRoomForm from '@/app/components/NewRoomForm';
import { StartNewGameConfig } from '@/app/types/newGame';
import { Player } from '@/app/types/player';
import {
  ConnectRequest,
  ConnectResponse,
  CreateRoomRequest,
  MessageType,
  VoteRequest,
  ShowVotesRequest,
  NewGameRequest,
} from '@/proto/planning_poker_pb';
import { PlanningPokerClient } from '@/proto/Planning_pokerServiceClientPb';
import VoteButton from './components/VoteButton';
import UserState from './components/UserState';

const client = new PlanningPokerClient('http://localhost:51000');
const FIBONACCHI = [0, 1, 2, 3, 5, 8, 13, 21];
const VOTE_RESET_NUMBER = -1;

export default function Home() {
  let lock = new AsyncLock();
  const [stream, setStream] = useState<ClientReadableStream<ConnectResponse> | null>(null);
  const [response, setResponse] = useState<ConnectResponse | null>(null);
  const [players, setPlayers] = useState<Map<string, Player>>(new Map());
  const [roomId, setRoomId] = useState<string>('');
  const [name, setName] = useState<string>('');
  const [votedNumber, setVotedNumber] = useState<number | null>(null);
  const [isShown, setIsShown] = useState<boolean>(false);

  const createNewRoom = (config: StartNewGameConfig) => {
    const req = new CreateRoomRequest
    req.setId(config.userName);
    req.setRoomname(config.room);

    if(stream !== null) {
      // TODO: 既存のストリームを終了させる処理をおこなう。新規部屋作成処理は続行する。
      console.error('stream is already exist');
      return;
    }

    const connection = client.createRoom(req);
    connection.on('data', async (res: ConnectResponse) => {
      await receiveMessage(res);
      setResponse(res);
    });

    setName(config.userName);
  };

  const joinRoom = (config: StartNewGameConfig) => {
    const req = new ConnectRequest();
    req.setId(config.userName);
    req.setRoomid(config.room);

    if(stream !== null) {
      // TODO: 既存のストリームを終了させる処理をおこなう。部屋参加処理は続行する。
      console.error('stream is already exist');
      return;
    }

    const connection = client.connect(req);
    connection.on('data', async (res: ConnectResponse) => {
      await receiveMessage(res);
      setResponse(res);
    });

    setName(config.userName);
    setRoomId(config.room);
  };

  const vote = (num: number) => {
    const req = new VoteRequest();

    if(num === votedNumber) {
      num = VOTE_RESET_NUMBER;
    }

    req.setId(name);
    req.setRoomid(roomId);
    req.setVote(num);

    client.vote(req, {}, (err, res) => {
      if (err) {
        console.error(err);
        return;
      }

      console.log(res);

      if(num === VOTE_RESET_NUMBER) {
        setVotedNumber(null);
        return;
      }
      setVotedNumber(num);
    });
  };

  const showVotes = () => {
    const req = new ShowVotesRequest();
    req.setRoomid(roomId);

    client.showVotes(req, {}, (err, res) => {
      if (err) {
        console.error(err);
        return;
      }

      console.log(res);
    });
  };

  const startNewGame = () => {
    const req = new NewGameRequest();
    req.setRoomid(roomId);

    client.newGame(req, {}, (err, res) => {
      if (err) {
        console.error(err);
        return;
      }

      console.log(res);
    });
  };

  const receiveMessage = async (res: ConnectResponse) => {
    console.log('receive message', res);
    await lock.acquire('receiveMessage', () => {
      switch (res.getType()) {
        case MessageType.JOIN:
          players.set(res.getMessage(), {name: res.getMessage(), isVoted: false, vote: null});
          setPlayers(players);
          break;
        case MessageType.VOTE:
          const player = players.get(res.getMessage())
          if (player) {
            player.isVoted = true;
          } else {
            players.set(res.getMessage(), {name: res.getMessage(), isVoted: true, vote: null});
          }
          setPlayers(players);
          break;
        case MessageType.SHOW_VOTES:
          const votes: Map<string, number> = JSON.parse(res.getMessage());
          for (const [key, value] of Object.entries(votes)) {
            const player = players.get(key);
            if (player) {
              player.vote = value;
              setPlayers(players);
            }
            setIsShown(true);
          }
          break;
        case MessageType.LEAVE:
          players.delete(res.getMessage());
          setPlayers(players);
          break;
        case MessageType.NEW_GAME:
          const newPlayers = new Map<string, Player>();
          for(let k in players.keys()) {
            newPlayers.set(k, {name: k, isVoted: false, vote: null})
          }
          setPlayers(newPlayers);
          setIsShown(false);
          break;
        case MessageType.CREATE_ROOM:
          setRoomId(res.getMessage());
          break;
        case MessageType.STATUS:
          const playerStatuses: Map<string, boolean> = JSON.parse(res.getMessage());
          for (const [key, value] of Object.entries(playerStatuses)) {
            players.set(key, {name: key, isVoted: value, vote: null});
            setPlayers(players);
          }
          break;
        case MessageType.RESET_VOTE:
          const resettingPlayer = players.get(res.getMessage())
          if (resettingPlayer) {
            resettingPlayer.isVoted = false;
          } else {
            players.set(res.getMessage(), {name: res.getMessage(), isVoted: false, vote: null});
          }
          setPlayers(players);
          break;
      }
      console.log(players);
    });
  }

  return (
    <main className='md:flex p-4 md:p-8'>
      <div className='md:w-1/3 px-6'>
        <NewGame createNewRoom={createNewRoom} joinRoom={joinRoom} />
      </div>
      <div className='mt-6 md:mt-0 md:w-2/3'>
        <div className='mb-4'>
          <h2 className='text-2xl font-bold leading-9 text-gray-900'>ルームID: {roomId}</h2>
          <p>ユーザ名: {name}</p>
        </div>
        <div className='mb-4 flex flex-wrap'>
          {
            Array.from(players.values()).map((player) => (
              <UserState key={player.name} player={player} />
            ))
          }
        </div>
        <div className='my-4'>
          <ResultAndNewGame
            isShown={isShown}
            showVotes={showVotes}
            startNewGame={startNewGame}
          />
        </div>
        <div>
          {
            FIBONACCHI.map((num) => (
              <VoteButton key={num} num={num} isSelected={num===votedNumber} vote={vote}/>
            ))
          }
        </div>
      </div>
    </main>
  )
}

type NewGameProps = {
  createNewRoom: (config: StartNewGameConfig) => void;
  joinRoom: (config: StartNewGameConfig) => void;
}

const NewGame: React.FC<NewGameProps> = ({ createNewRoom, joinRoom }) => {
  return (
    <div className='flex flex-col'>
      <NewGameCard title='新規ルームを作成' form={<NewRoomForm connectRoom={createNewRoom} />} />
      <NewGameCard title='既存のルームへ参加' form={<JoinRoomForm connectRoom={joinRoom} />} />
    </div>
  );
}

type ResultAndNewGameProps = {
  isShown: boolean;
  showVotes: () => void;
  startNewGame: () => void;
}

const ResultAndNewGame: React.FC<ResultAndNewGameProps> = ({ isShown, showVotes, startNewGame }) => {
  const name = isShown ? '次のポーカーへ' : '結果を見る';

  const onClick: React.MouseEventHandler<HTMLButtonElement> = (event) => {
    event.preventDefault()
    if (isShown) {
      startNewGame();
    } else {
      showVotes();
    }
  }
  
  return (
    <button
      className='flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600'
      onClick={onClick}
    >
      {name}
    </button>
  );
};