type VoteButtonProps = {
  num: number;
  isSelected: boolean;
  vote: (num: number) => void;
}

export default function VoteButton({num, isSelected, vote}: VoteButtonProps) {
  const bgColor = isSelected ? 'bg-indigo-600' : 'bg-gray-400';

  const onClick: React.MouseEventHandler<HTMLButtonElement> = (event) => {
    event.preventDefault();
    vote(num);
  }

  return (
    <button
      className={`${bgColor} shadow-sm rounded-lg font-semibold w-24 h-36 m-2 text-4xl text-white hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600`}
      onClick={onClick}
    >{num}</button>
  );
}
