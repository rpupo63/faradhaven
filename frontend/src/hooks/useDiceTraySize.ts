import { useCallback, useEffect, useState } from 'react';

const MAX_TRAY = 800;
const MIN_TRAY = 280;
const VIEWPORT_FRACTION = 0.92;

function readViewport(): { w: number; h: number } {
  const vv = typeof window !== 'undefined' ? window.visualViewport : null;
  return {
    w: vv?.width ?? window.innerWidth,
    h: vv?.height ?? window.innerHeight,
  };
}

export function computeTrayPixelSize(): number {
  const { w, h } = readViewport();
  const edge = Math.min(w, h) * VIEWPORT_FRACTION;
  return Math.min(MAX_TRAY, Math.max(MIN_TRAY, Math.floor(edge)));
}

/** Dice visual scale: smaller tray → slightly smaller dice so they stay in view. */
export function computeDiceScale(trayPx: number): number {
  if (trayPx < 400) return 12;
  if (trayPx < 520) return 16;
  if (trayPx < 640) return 20;
  return 24;
}

export function useDiceTraySize(): number {
  const [trayPx, setTrayPx] = useState(() =>
    typeof window !== 'undefined' ? computeTrayPixelSize() : 520
  );

  const measure = useCallback(() => {
    setTrayPx(computeTrayPixelSize());
  }, []);

  useEffect(() => {
    measure();
    window.addEventListener('resize', measure);
    const vv = window.visualViewport;
    vv?.addEventListener('resize', measure);
    vv?.addEventListener('scroll', measure);
    return () => {
      window.removeEventListener('resize', measure);
      vv?.removeEventListener('resize', measure);
      vv?.removeEventListener('scroll', measure);
    };
  }, [measure]);

  return trayPx;
}
