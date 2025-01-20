// src/hooks/useWindow.ts
import { useState, useEffect } from 'react';

/*
  Custom hook to get and track the current window width.
  This hook listens for window resize events and updates the width dynamically.
  returns The current width of the window.
*/

export default function useWindowWidth() {
  const [gridWidth, setGridWidth] = useState<number>(window.innerWidth);

  useEffect(() => {
    const handleResize = () => {
      setGridWidth(window.innerWidth);
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, []);

  return gridWidth;
}