import { useCallback, useEffect, useState } from 'react';
import { MIN_MAIN_BOUNDS_PX, useMainContentBounds } from '@/context/MainContentBoundsContext';

/** Cap tray so dice stay a comfortable size on large / ultrawide screens. */
const MAX_TRAY = 440;
const MIN_TRAY = 280;
/** Fraction of the shorter viewport/main edge — 0.92 made trays ~800px and dice huge. */
const VIEWPORT_FRACTION = 0.52;

function readViewport(): { w: number; h: number } {
  const vv = typeof window !== 'undefined' ? window.visualViewport : null;
  return {
    w: vv?.width ?? window.innerWidth,
    h: vv?.height ?? window.innerHeight,
  };
}

/** When `content` is set, tray is sized from the main column; otherwise from the visual viewport. */
export function computeTrayPixelSize(content?: { w: number; h: number } | null): number {
  let w: number;
  let h: number;
  if (
    content != null &&
    content.w >= MIN_MAIN_BOUNDS_PX &&
    content.h >= MIN_MAIN_BOUNDS_PX
  ) {
    w = content.w;
    h = content.h;
  } else {
    const v = readViewport();
    w = v.w;
    h = v.h;
  }
  const edge = Math.min(w, h) * VIEWPORT_FRACTION;
  return Math.min(MAX_TRAY, Math.max(MIN_TRAY, Math.floor(edge)));
}

/**
 * Die mesh scale for @3d-dice/dice-box — higher = larger dice in the scene.
 * Tray size (#dice-box pixels) is unchanged; this only scales the 3D die meshes inside it.
 */
export function computeDiceScale(trayPx: number): number {
  const base =
    trayPx <= 300 ? 15 : trayPx <= 360 ? 16 : trayPx <= 400 ? 17 : trayPx <= 440 ? 18 : 19;
  // Half of legacy base, +1 for slightly larger dice (tray px unchanged).
  return Math.max(4, Math.min(12, Math.round(base / 2) + 1));
}

export function useDiceTraySize(): number {
  const mainBounds = useMainContentBounds();
  const mainW = mainBounds?.width ?? 0;
  const mainH = mainBounds?.height ?? 0;

  const [trayPx, setTrayPx] = useState(() =>
    typeof window !== 'undefined' ? computeTrayPixelSize() : MAX_TRAY
  );

  const measure = useCallback(() => {
    const content =
      mainW >= MIN_MAIN_BOUNDS_PX && mainH >= MIN_MAIN_BOUNDS_PX
        ? { w: mainW, h: mainH }
        : null;
    setTrayPx(computeTrayPixelSize(content));
  }, [mainW, mainH]);

  useEffect(() => {
    const raf = requestAnimationFrame(() => measure());
    window.addEventListener('resize', measure);
    const vv = window.visualViewport;
    vv?.addEventListener('resize', measure);
    vv?.addEventListener('scroll', measure);
    return () => {
      cancelAnimationFrame(raf);
      window.removeEventListener('resize', measure);
      vv?.removeEventListener('resize', measure);
      vv?.removeEventListener('scroll', measure);
    };
  }, [measure]);

  return trayPx;
}
