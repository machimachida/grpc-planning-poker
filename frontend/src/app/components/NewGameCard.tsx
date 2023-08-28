type NewGameCardProps = {
  title: string;
  form: React.ReactNode;
}

export default function NewGameCard({title, form}: NewGameCardProps) {
  return (
    <>
      <div className='mx-auto w-full'>
        <h2 className='text-center text-2xl font-bold leading-9 text-gray-900'>
          {title}
        </h2>
      </div>

      <div className='my-10'>
        {form}
      </div>
    </>
  );
}
