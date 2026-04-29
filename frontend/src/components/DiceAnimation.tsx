import { useEffect, useLayoutEffect, useRef, useState } from 'react';
import { createPortal } from 'react-dom';
import { AnimatePresence, motion } from 'framer-motion';
import { useLocation } from 'react-router-dom';
import { DICE_CLEAR_EVENT } from '@/lib/dice';
import { DiceManager, type DiceResult } from '@/lib/dice-manager';
import { useAuth } from '@/context/AuthContext';
import {
  computeDiceScale,
  computeTrayPixelSize,
  useDiceTraySize,
} from '@/hooks/useDiceTraySize';

interface ResultDisplay extends DiceResult {
  id: string;
  notation: string;
}

/**
 * Portaled to body. Use 100dvh + w-screen so the flex center matches the *visual* viewport
 * (inset-0 / 100vh alone often sits low on mobile when the URL bar resizes the layout viewport).
 */
const DICE_LAYER_BASE =
  'fixed left-0 right-0 top-0 flex h-[100dvh] w-screen max-w-[100vw] min-h-0';

/** @3d-dice/dice-box: height where the toss starts (default 8). We had 4–5 which read as “too low”. */
const STARTING_HEIGHT = { narrow: 12, wide: 14 } as const;

export function DiceAnimation() {
  const trayPx = useDiceTraySize();
  const [diceBoxGeneration, setDiceBoxGeneration] = useState(0);
  const [result, setResult] = useState<ResultDisplay | null>(null);
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const location = useLocation();
  const dismissTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const { activeDicePrefs } = useAuth();

  // Keep a stable ref to the box instance so we can call updateConfig later.
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const boxRef = useRef<any | null>(null);

  // Pre-load audio (effect — no DOM dependency)
  useEffect(() => {
    const audio = new Audio('/dice-rolling.mp3');
    audio.load();
    audioRef.current = audio;
  }, []);

  // Init dice-box after #dice-box is in the portaled DOM (layout phase beats async paint races).
  useLayoutEffect(() => {
    let cancelled = false;

    import('@3d-dice/dice-box').then(({ default: DiceBox }) => {
      let didInit = false;
      const attach = () => {
        if (cancelled || didInit) return;
        if (!document.getElementById('dice-box')) {
          requestAnimationFrame(attach);
          return;
        }
        didInit = true;

      const tray = computeTrayPixelSize();
      const scale = computeDiceScale(tray);
      const narrow = tray < 400;

      const box = new DiceBox({
        container: '#dice-box',
        assetPath: '/assets/dice-box/',
        scale,
        gravity: 14,
        mass: 1,
        friction: 0.8,
        restitution: 0.1,
        angularDamping: 0.4,
        linearDamping: 0.4,
        spinForce: narrow ? 4 : 6,
        throwForce: narrow ? 2 : 3,
        startingHeight: narrow ? STARTING_HEIGHT.narrow : STARTING_HEIGHT.wide,
        lightIntensity: 1,
        theme: activeDicePrefs.dice_theme,
        themeColor: activeDicePrefs.dice_theme_color,
        fontColor: activeDicePrefs.dice_font_color,
      });

      box.init().then(() => {
        if (cancelled) return;
        boxRef.current = box;
        DiceManager.register(box);
        setDiceBoxGeneration((n) => n + 1);
      }).catch(console.error);
      };
      requestAnimationFrame(attach);
    }).catch(console.error);

    return () => {
      cancelled = true;
    };
    // Intentionally run only on mount — prefs changes handled via updateConfig below.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Apply preference changes live via updateConfig (no re-init needed).
  // Include diceBoxGeneration so prefs apply once the box finishes init (user prefs often load after first paint).
  useEffect(() => {
    if (!boxRef.current) return;
    boxRef.current.updateConfig({
      theme: activeDicePrefs.dice_theme,
      themeColor: activeDicePrefs.dice_theme_color,
      fontColor: activeDicePrefs.dice_font_color,
    });
  }, [
    activeDicePrefs.dice_theme,
    activeDicePrefs.dice_theme_color,
    activeDicePrefs.dice_font_color,
    diceBoxGeneration,
  ]);

  // Keep canvas / physics bounds aligned when the viewport or tray size changes.
  useEffect(() => {
    const box = boxRef.current;
    if (!box) return;
    const scale = computeDiceScale(trayPx);
    const narrow = trayPx < 400;
    void box.updateConfig({
      scale,
      spinForce: narrow ? 4 : 6,
      throwForce: narrow ? 2 : 3,
      startingHeight: narrow ? STARTING_HEIGHT.narrow : STARTING_HEIGHT.wide,
    });
    if (typeof box.resizeWorld === 'function') {
      box.resizeWorld();
    }
  }, [trayPx, diceBoxGeneration]);

  // Play sound the moment dice appear (roll start), not when they settle.
  useEffect(() => {
    DiceManager.setOnRollStart(() => {
      if (audioRef.current) {
        audioRef.current.currentTime = 0;
        audioRef.current.play().catch(() => {});
      }
    });
    return () => DiceManager.setOnRollStart(null);
  }, []);

  // Show result overlay after dice physically settle.
  useEffect(() => {
    DiceManager.setOnResult((r: DiceResult & { notation: string }) => {
      if (dismissTimer.current) clearTimeout(dismissTimer.current);
      setResult({ id: crypto.randomUUID(), ...r });
      // Results stay up until clicked or 6 seconds pass
      dismissTimer.current = setTimeout(() => {
        setResult(null);
        DiceManager.clear();
      }, 6000);
    });
    return () => DiceManager.setOnResult(null);
  }, []);

  // Clear on route change.
  useEffect(() => {
    setResult(null);
    DiceManager.clear();
  }, [location.pathname]);

  // DICE_CLEAR_EVENT — dispatched by dialogs/panels on close.
  // Click anywhere — dismisses the overlay AND removes physical dice.
  useEffect(() => {
    const handleClear = () => {
      setResult(null);
      if (dismissTimer.current) clearTimeout(dismissTimer.current);
      DiceManager.clear();
    };
    const handleClick = () => {
      // Only clear if a result is currently showing.
      // This prevents the click that triggers the roll from immediately clearing it.
      setResult((prev) => {
        if (prev) {
          DiceManager.clear();
          if (dismissTimer.current) clearTimeout(dismissTimer.current);
          return null;
        }
        return null;
      });
    };

    window.addEventListener(DICE_CLEAR_EVENT, handleClear);
    window.addEventListener('click', handleClick);
    return () => {
      window.removeEventListener(DICE_CLEAR_EVENT, handleClear);
      window.removeEventListener('click', handleClick);
    };
  }, []);

  /*
   * Portal to document.body so fixed layers use the real viewport (Layout uses transforms /
   * overflow / blur that break fixed + getBoundingClientRect combos).
   * @3d-dice/dice-box appends an inline <canvas> — force it to fill #dice-box or it stays ~300×150.
   */
  const layers = (
    <>
      {/* Dim backdrop while something is rolling or result is showing */}
      <AnimatePresence>
        {result && (
          <motion.div
            key="backdrop"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className={`${DICE_LAYER_BASE} z-[9998] bg-black/60 backdrop-blur-sm pointer-events-auto`}
          />
        )}
      </AnimatePresence>

      <div
        className={`${DICE_LAYER_BASE} z-[9999] pointer-events-none items-center justify-center`}
      >
        <AnimatePresence>
          {result && (
            <motion.div
              key={result.id}
              initial={{ opacity: 0, scale: 0.9, y: 20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.9, y: -20 }}
              transition={{ type: 'spring', stiffness: 300, damping: 25 }}
              className="bg-background/95 border-2 border-primary/30 rounded-3xl p-4 sm:p-8 shadow-[0_0_50px_rgba(0,0,0,0.5)] backdrop-blur-xl flex flex-col items-center gap-4 min-w-0 max-w-[calc(100%-2rem)] mx-4"
            >
              {result.label && (
                <span className="text-sm font-tome-marginalia text-muted-foreground uppercase tracking-[0.2em]">
                  {result.label}
                </span>
              )}

              {/* Total */}
              <div className="relative">
                <span className="text-5xl sm:text-7xl font-bold font-display text-primary tracking-tighter leading-none glow-text">
                  {result.total}
                </span>
              </div>

              {/* Breakdown */}
              <div className="flex flex-col items-center gap-1">
                <div className="flex items-center gap-2 text-muted-foreground font-tome-marginalia">
                  {result.rolls.length === 1 ? (
                    <span className="text-lg">
                      {result.rolls[0]}
                      {result.modifier !== undefined && result.modifier !== 0 && (
                        <> {result.modifier > 0 ? '+' : '-'} {Math.abs(result.modifier)}</>
                      )}
                    </span>
                  ) : (
                    <div className="flex flex-wrap justify-center gap-1 max-w-[200px]">
                      {result.rolls.map((v, i) => (
                        <span key={i} className="text-sm opacity-80">
                          {v}{i < result.rolls.length - 1 ? ' +' : ''}
                        </span>
                      ))}
                      {result.modifier !== undefined && result.modifier !== 0 && (
                        <span className="text-sm">
                          {result.modifier > 0 ? '+' : '-'} {Math.abs(result.modifier)}
                        </span>
                      )}
                    </div>
                  )}
                </div>
                <span className="text-micro uppercase font-tome-marginalia text-muted-foreground/60 tracking-widest">
                  {result.notation}
                </span>
              </div>

              {/* Individual die values for multi-roll when many dice are involved */}
              {result.rolls.length > 1 && result.rolls.length <= 10 && (
                <div className="flex flex-wrap justify-center gap-2 mt-2">
                  {result.rolls.map((v, i) => (
                    <span
                      key={i}
                      className="w-8 h-8 flex items-center justify-center rounded-lg border border-primary/20 bg-primary/5 font-display text-xs text-primary/80"
                    >
                      {v}
                    </span>
                  ))}
                </div>
              )}
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      <div
        className={`${DICE_LAYER_BASE} z-[10000] pointer-events-none items-center justify-center`}
      >
        <div
          id="dice-box"
          className="relative shrink-0 [&>canvas.dice-box-canvas]:block [&>canvas.dice-box-canvas]:h-full [&>canvas.dice-box-canvas]:w-full"
          style={{
            width: `${trayPx}px`,
            height: `${trayPx}px`,
          }}
        />
      </div>
    </>
  );

  if (typeof document === 'undefined') {
    return null;
  }

  return createPortal(layers, document.body);
}
