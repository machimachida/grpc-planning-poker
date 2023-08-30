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

const FIBONACCHI = [1, 2, 3, 5, 8, 13, 21];
const VOTE_RESET_NUMBER = -1;
const AVERAGE_KEY_NAME = "average";

if(process.env.NEXT_PUBLIC_SERVER === undefined) {
  throw new Error('SERVER is not defined');
}
const client = new PlanningPokerClient(process.env.NEXT_PUBLIC_SERVER);


type Players = {
  [key: string]: Player;
}

export default function Home() {
  let lock = new AsyncLock();
  const [stream, setStream] = useState<ClientReadableStream<ConnectResponse> | null>(null);
  const [response, setResponse] = useState<ConnectResponse | null>(null); // これがないとなぜか再レンダリングしない。要調査
  const [players, setPlayers] = useState<Players>({});
  const [roomId, setRoomId] = useState<string>('');
  const [name, setName] = useState<string>('');
  const [votedNumber, setVotedNumber] = useState<number | null>(null);
  const [isShown, setIsShown] = useState<boolean>(false);
  const [average, setAverage] = useState<number>(0);

  const createNewRoom = (config: StartNewGameConfig) => {
    const req = new CreateRoomRequest
    req.setId(config.userName);
    req.setRoomname(config.room);

    if(stream !== null) {
      stream.cancel();
    }

    const connection = client.createRoom(req);
    connection.on('data', async (res: ConnectResponse) => {
      await receiveMessage(res);
      setResponse(res);
    });

    setStream(connection);
    setName(config.userName);
  };

  const joinRoom = (config: StartNewGameConfig) => {
    const req = new ConnectRequest();
    req.setId(config.userName);
    req.setRoomid(config.room);

    if(stream !== null) {
      stream.cancel();
    }

    const connection = client.connect(req);
    connection.on('data', async (res: ConnectResponse) => {
      await receiveMessage(res);
      setResponse(res);
    });

    setStream(connection);
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
        alert('投票できませんでした。');
        return;
      }

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
        alert('投票結果を表示できませんでした。');
        return;
      }
    });
  };

  const startNewGame = () => {
    const req = new NewGameRequest();
    req.setRoomid(roomId);

    client.newGame(req, {}, (err, _res) => {
      if (err) {
        console.error(err);
        alert('新しいゲームを開始できませんでした。');
        return;
      }
    });
  };

  const receiveMessage = async (res: ConnectResponse) => {
    console.log("receive message", res)
    await lock.acquire('receiveMessage', () => {
      switch (res.getType()) {
        case MessageType.JOIN:
          console.log("join", res.getMessage());
          players[res.getMessage()] = {name: res.getMessage(), isVoted: false, vote: null};
          setPlayers(players);
          break;
        case MessageType.VOTE:
          console.log("vote", res.getMessage());
          const player = players[res.getMessage()]
          if (player) {
            player.isVoted = true;
          } else {
            players[res.getMessage()] = {name: res.getMessage(), isVoted: false, vote: null};
          }
          setPlayers(players);
          break;
        case MessageType.SHOW_VOTES:
          console.log("show votes", res.getMessage());
          const votes: {[key: string]: number} = JSON.parse(res.getMessage());
          for (const [key, value] of Object.entries(votes)) {
            if (key === AVERAGE_KEY_NAME) {
              setAverage(value);
              continue;
            }

            const player = players[key];
            if (player) {
              player.vote = value;
              setPlayers(players);
            }
            setIsShown(true);
          }
          break;
        case MessageType.LEAVE:
          console.log("leave", res.getMessage());
          delete players[res.getMessage()];
          setPlayers(players);
          break;
        case MessageType.NEW_GAME:
          console.log("new game", res.getMessage());
          for(const [key, _value] of Object.entries(players)) {
            players[key] = {name: key, isVoted: false, vote: null}
          }
          setPlayers(players);
          setVotedNumber(null);
          setIsShown(false);
          break;
        case MessageType.CREATE_ROOM:
          console.log("create room", res.getMessage());
          setRoomId(res.getMessage());
          break;
        case MessageType.STATUS:
          console.log("status", res.getMessage());
          const playerStatuses: {[key: string]: boolean} = JSON.parse(res.getMessage());
          for (const [key, value] of Object.entries(playerStatuses)) {
            players[key] = {name: key, isVoted: value, vote: null};
          }
          setPlayers(players);
          console.log(players);
          break;
        case MessageType.RESET_VOTE:
          console.log("reset vote", res.getMessage());
          const resettingPlayer = players[res.getMessage()]
          if (resettingPlayer) {
            resettingPlayer.isVoted = false;
          } else {
            players[res.getMessage()] = {name: res.getMessage(), isVoted: false, vote: null};
          }
          setPlayers(players);
          break;
      }
      console.log("players", players)
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
            Object.values(players).map((player) => (
              <UserState key={player.name} player={player} />
            ))
          }
        </div>
        <div className='h-8'>
          {isShown && <p className='font-bold leading-9 text-gray-900'>平均: {average}</p>
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