import NewGameCard from '@/app/components/NewGameCard';
import NewRoomForm from '@/app/components/NewRoomForm';
import JoinRoomForm from '@/app/components/JoinRoomForm';
import { StartNewGameConfig } from '@/app/types/newGame';

type NewGameProps = {
  createNewRoom: (config: StartNewGameConfig) => void;
  joinRoom: (config: StartNewGameConfig) => void;
}

export default function NewGame({ createNewRoom, joinRoom }: NewGameProps) {
  return (
    <div className='flex min-h-full flex-1 flex-col px-6 py-12 lg:px-8'>
      <NewGameCard title='新規ルームを作成' form={<NewRoomForm connectRoom={createNewRoom} />} />
      <NewGameCard title='既存のルームへ参加' form={<JoinRoomForm connectRoom={joinRoom} />} />
    </div>
  );
}