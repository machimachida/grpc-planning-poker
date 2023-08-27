type ExternalButtonProps = {
  onClick: React.MouseEventHandler<HTMLButtonElement>;
  buttonRef: React.RefObject<HTMLButtonElement>;
};

const ExternalButton: React.FC<ExternalButtonProps> = (
  {onClick, buttonRef}
) => (
  <button
    ref={buttonRef}
    className='my-4'
    onClick={onClick}
  >
    <span>Button</span>
  </button>
);

export default ExternalButton;