import { useForm } from 'react-hook-form';
import { StartNewGameConfig } from '@/app/types/newGame';

type NewRoomFormData = {
  roomName: string;
  userName: string;
}

type NewRoomFormProps = {
  connectRoom: (config: StartNewGameConfig) => void;
}

export default function NewRoomForm({ connectRoom }: NewRoomFormProps) {
  const { register, handleSubmit } = useForm<NewRoomFormData>();

  const onSubmit = (data: NewRoomFormData) => {
    connectRoom({room: data.roomName, userName: data.userName})
  };

  return (
    <form className='space-y-6' action='#' method='POST' onSubmit={handleSubmit(onSubmit)}>
      <div>
        <label htmlFor='roomName' className='block text-sm font-medium leading-6 text-gray-900'>
          ルーム名
        </label>
        <div className='mt-2'>
          <input
            id='roomName'
            required
            className='block w-full rounded-md border-0 px-3 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6'
            {...register('roomName')} />
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
        ルームを作成
      </button>
    </form>
  );
}