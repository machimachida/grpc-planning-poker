'use client';
import {useEffect, useState, useRef} from 'react';

type ExternalDisplayProps = {
  buttonRef: React.RefObject<HTMLButtonElement>;
  outsideClickCore: () => void;
};

const ExternalDisplay: React.FC<ExternalDisplayProps> = ({
  buttonRef,
  outsideClickCore,
}) => {
  const divRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const outSideClick = (event: MouseEvent) => {
      if (divRef.current?.contains(event.target as Node)) {
        return;
      }
      if (buttonRef.current?.contains(event.target as Node)) {
        return;
      }
      outsideClickCore();
    };
    document.addEventListener('click', outSideClick);
    return () => {
      document.removeEventListener('click', outSideClick);
    };
  }, [open]);

  return (
    <div
      ref={divRef}
      className='w-64 h-64 bg-gray-200'
    ></div>
  );
}

export default ExternalDisplay;