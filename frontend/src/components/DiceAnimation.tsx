import { useEffect, useRef, useState } from 'react';
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

  // Initialize dice-box once on mount and register with DiceManager.
  useEffect(() => {
    let cancelled = false;

    // Pre-load audio
    const audio = new Audio('/dice-rolling.mp3');
    audio.load();
    audioRef.current = audio;

    import('@3d-dice/dice-box').then(({ default: DiceBox }) => {
      if (cancelled) return;

      const tray = computeTrayPixelSize();
      const scale = computeDiceScale(tray);
      const narrow = tray < 520;

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
        startingHeight: narrow ? 4 : 5,
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
    }).catch(console.error);

    return () => {
      cancelled = true;
    };
    // Intentionally run only on mount — prefs changes handled via updateConfig below.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Apply preference changes live via updateConfig (no re-init needed).
  useEffect(() => {
    if (!boxRef.current) return;
    boxRef.current.updateConfig({
      theme: activeDicePrefs.dice_theme,
      themeColor: activeDicePrefs.dice_theme_color,
      fontColor: activeDicePrefs.dice_font_color,
    });
  }, [activeDicePrefs.dice_theme, activeDicePrefs.dice_theme_color, activeDicePrefs.dice_font_color]);

  // Keep canvas / physics bounds aligned when the viewport or tray size changes.
  useEffect(() => {
    const box = boxRef.current;
    if (!box) return;
    const scale = computeDiceScale(trayPx);
    const narrow = trayPx < 520;
    void box.updateConfig({
      scale,
      spinForce: narrow ? 4 : 6,
      throwForce: narrow ? 2 : 3,
      startingHeight: narrow ? 4 : 5,
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
      dismissTimer.current = setTimeout(() => setResult(null), 6000);
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

  return (
    <>
      {/* Dim backdrop while something is rolling or result is showing */}
      <AnimatePresence>
        {result && (
          <motion.div
            key="backdrop"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 z-[9998] bg-black/60 backdrop-blur-sm pointer-events-auto"
          />
        )}
      </AnimatePresence>

      {/*
        Result panel — Centered on screen.
        Z-index 9999 so it's behind the dice (Z-index 10000).
      */}
      <div
        className="pointer-events-none fixed inset-0 flex items-center justify-center z-[9999]"
      >
        <AnimatePresence>
          {result && (
            <motion.div
              key={result.id}
              initial={{ opacity: 0, scale: 0.9, y: 20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.9, y: -20 }}
              transition={{ type: 'spring', stiffness: 300, damping: 25 }}
              className="bg-background/95 border-2 border-primary/30 rounded-3xl p-4 sm:p-8 shadow-[0_0_50px_rgba(0,0,0,0.5)] backdrop-blur-xl flex flex-col items-center gap-4 min-w-0 max-w-[calc(100vw-2rem)] mx-4"
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

      {/*
        Bounded dice tray — centered on screen.
        dice-box creates a canvas sized to this element; physics walls
        are placed at its edges, so dice stay within this area.
      */}
      <div
        id="dice-box"
        style={{
          position: 'fixed',
          left: '50%',
          top: '50%',
          transform: 'translate(-50%, -50%)',
          width: trayPx,
          height: trayPx,
          zIndex: 10000,
          pointerEvents: 'none',
        }}
      />
    </>
  );
}
