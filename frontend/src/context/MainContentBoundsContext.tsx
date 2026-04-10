import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useLayoutEffect,
  useState,
  type RefObject,
  type ReactNode,
} from 'react';

export interface MainContentBounds {
  left: number;
  top: number;
  width: number;
  height: number;
}

/** Ignore pre-layout / broken flex passes so we fall back to visual viewport. */
export const MIN_MAIN_BOUNDS_PX = 64;

const MainContentBoundsContext = createContext<MainContentBounds | null>(null);

export function useMainContentBounds(): MainContentBounds | null {
  return useContext(MainContentBoundsContext);
}

function readVisualViewportAsBounds(): MainContentBounds {
  if (typeof window === 'undefined') {
    return { left: 0, top: 0, width: 0, height: 0 };
  }
  const vv = window.visualViewport;
  return {
    left: vv?.offsetLeft ?? 0,
    top: vv?.offsetTop ?? 0,
    width: vv?.width ?? window.innerWidth,
    height: vv?.height ?? window.innerHeight,
  };
}

/** Full visual viewport in screen coords — use when overlay rect is invalid or for portaled fixed layers. */
export function getViewportBoundsForDiceOverlay(): MainContentBounds {
  return readVisualViewportAsBounds();
}

function isUsableMainBounds(b: MainContentBounds | null): b is MainContentBounds {
  return (
    b != null &&
    b.width >= MIN_MAIN_BOUNDS_PX &&
    b.height >= MIN_MAIN_BOUNDS_PX
  );
}

/** Main column rect when available; otherwise visual viewport (for fixed overlay positioning). */
export function useDiceOverlayTargetRect(): MainContentBounds {
  const rawMain = useMainContentBounds();
  const mainBounds = isUsableMainBounds(rawMain) ? rawMain : null;
  const [vvRect, setVvRect] = useState<MainContentBounds>(() =>
    typeof window !== 'undefined' ? readVisualViewportAsBounds() : { left: 0, top: 0, width: 0, height: 0 }
  );

  useEffect(() => {
    if (mainBounds != null) return;
    const read = () => setVvRect(readVisualViewportAsBounds());
    read();
    window.addEventListener('resize', read);
    const vv = window.visualViewport;
    vv?.addEventListener('resize', read);
    vv?.addEventListener('scroll', read);
    return () => {
      window.removeEventListener('resize', read);
      vv?.removeEventListener('resize', read);
      vv?.removeEventListener('scroll', read);
    };
  }, [mainBounds]);

  return mainBounds ?? vvRect;
}

export function MainContentBoundsProvider({
  mainRef,
  children,
}: {
  mainRef: RefObject<HTMLElement | null>;
  children: ReactNode;
}) {
  const [bounds, setBounds] = useState<MainContentBounds | null>(null);

  const measure = useCallback(() => {
    const el = mainRef.current;
    if (!el) {
      setBounds(null);
      return;
    }
    const r = el.getBoundingClientRect();
    if (r.width < MIN_MAIN_BOUNDS_PX || r.height < MIN_MAIN_BOUNDS_PX) {
      setBounds(null);
      return;
    }
    setBounds({
      left: r.left,
      top: r.top,
      width: r.width,
      height: r.height,
    });
  }, [mainRef]);

  useLayoutEffect(() => {
    const el = mainRef.current;
    const ro = el ? new ResizeObserver(() => measure()) : null;
    if (el && ro) ro.observe(el);
    window.addEventListener('resize', measure);
    const vv = window.visualViewport;
    vv?.addEventListener('resize', measure);
    vv?.addEventListener('scroll', measure);

    // Read after flex layout; sync so first paint has real <main> size (rAF could run before flex resolves).
    // eslint-disable-next-line react-hooks/set-state-in-effect -- bounds synced from layout target
    measure();

    return () => {
      ro?.disconnect();
      window.removeEventListener('resize', measure);
      vv?.removeEventListener('resize', measure);
      vv?.removeEventListener('scroll', measure);
    };
  }, [measure, mainRef]);

  return (
    <MainContentBoundsContext.Provider value={bounds}>
      {children}
    </MainContentBoundsContext.Provider>
  );
}
