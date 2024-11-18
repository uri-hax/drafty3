// hooks/useWindowWidth.ts
import { useState, useEffect } from 'react';

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