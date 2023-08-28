import { useForm } from 'react-hook-form';
import { StartNewGameConfig } from '@/app/types/newGame';

type JoinRoomFormData = {
  roomId: string;
  userName: string;
}

type JoinRoomFormProps = {
  connectRoom: (config: StartNewGameConfig) => void;
}

export default function JoinRoomForm({ connectRoom }: JoinRoomFormProps) {
  const { register, handleSubmit } = useForm<JoinRoomFormData>();

  const onSubmit = (data: JoinRoomFormData) => {
    connectRoom({room: data.roomId, userName: data.userName});
  };

  return (
    <form className='space-y-6' action='#' method='POST' onSubmit={handleSubmit(onSubmit)}>
      <div>
        <label htmlFor='roomId' className='block text-sm font-medium leading-6 text-gray-900'>
          ルームID
        </label>
        <div className='mt-2'>
          <input
            id='roomId'
            required
            className='block w-full rounded-md border-0 px-3 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6'
            {...register('roomId')} />
        </div>
      </div>
      <div>
        <label htmlFor='userName' className='block text-sm font-medium leading-6 text-gray-900'>
          ユーザ名
        </label>
        <div className='mt-2'>
          <input
            id='userName'
            required
            className='block w-full rounded-md border-0 px-3 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6'
            {...register('userName')}
          />
        </div>
      </div>
      <button
        type='submit'
        className='flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600'
      >
        ルームへ参加
      </button>
    </form>
  );
}
