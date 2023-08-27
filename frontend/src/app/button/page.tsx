'use client';
import {useEffect, useState, useRef} from 'react';
import ExternalButton from '@/components/ExternalButton';
import ExternalDisplay from '@/components/ExternalDisplay';

export default function Page() {
  const [open, setOpen] = useState(false);

  const buttonRef = useRef<HTMLButtonElement>(null);
  //const divRef = useRef<HTMLDivElement>(null);

  const onClick: React.MouseEventHandler<HTMLButtonElement> = (_event) => {
    setOpen(!open);
  };

  //const outSideClick = (event: MouseEvent) => {
  //  if (buttonRef.current?.contains(event.target as Node)) {
  //    return;
  //  }
  //  if (divRef.current?.contains(event.target as Node)) {
  //    return;
  //  }
  //  setOpen(false);
  //};

  //useEffect(() => {
  //  if (open) {
  //    document.addEventListener('click', outSideClick);
  //  } else {
  //    document.removeEventListener('click', outSideClick);
  //  }
  //  return () => {
  //    document.removeEventListener('click', outSideClick);
  //  };
  //}, [open]);

  return (
    <div>
      <h1>Welcome to Button!</h1>
      <ExternalButton
        onClick={onClick}
        buttonRef={buttonRef}
      />
      {open && (<ExternalDisplay buttonRef={buttonRef} outsideClickCore={() => setOpen(false)} />)}
    </div>
  );
}
