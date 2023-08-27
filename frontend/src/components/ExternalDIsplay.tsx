'use client';

type ExternalDisplayProps = {
  divRef: React.RefObject<HTMLDivElement>;
};

const ExternalDisplay: React.FC<ExternalDisplayProps> = ({divRef}) => {

  return (
    <div
      ref={divRef}
      className='w-64 h-64 bg-gray-200'
    ></div>
  );
}

export default ExternalDisplay;