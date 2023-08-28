import { Player } from '@/app/types/player';

type UserStateProps = {
  player: Player;
}

export default function UserState({player}: UserStateProps) {
  const bgColor = player.isVoted && player.vote === null ? 'bg-indigo-600' : 'bg-gray-400';
  return (
    <div
      key={player.name}
      className='p-2'
    >
      <div className={`${bgColor} rounded-lg w-24 h-36 flex justify-center items-center`}>
        {
          player.vote && (
            <p className='font-semibold text-4xl text-white'>{player.vote}</p>
          )
        }
      </div>
      <p className='line-clamp-2'>{player.name}</p>
    </div>
  )
}
