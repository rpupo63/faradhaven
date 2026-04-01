import React, { useState, useEffect, useRef } from 'react';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import { toast } from 'sonner';
import type { NormalizedCharacterSheet } from '@/types/game';

interface PowderMageFeaturesCardProps {
  sheet: NormalizedCharacterSheet;
  token: string;
}

export const PowderMageFeaturesCard: React.FC<PowderMageFeaturesCardProps> = ({ sheet }) => {
  const resources = sheet.class_resources ?? [];
  const timerDuration = resources.find(r => r.key === 'timer_duration')?.value ?? 2;
  const speedDialSlots = resources.find(r => r.key === 'speed_dial_slots')?.value ?? 0;

  const [casting, setCasting] = useState(false);
  const [secondsLeft, setSecondsLeft] = useState(0);
  const [done, setDone] = useState(false);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    if (casting && secondsLeft > 0) {
      timerRef.current = setInterval(() => {
        setSecondsLeft(s => {
          if (s <= 1) {
            setCasting(false);
            setDone(true);
            toast.info('Casting window closed!');
            return 0;
          }
          return s - 1;
        });
      }, 1000);
    }
    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
    };
  }, [casting, secondsLeft]);

  const startCast = () => {
    setDone(false);
    setCasting(true);
    setSecondsLeft(timerDuration);
  };

  const reset = () => {
    setCasting(false);
    setDone(false);
    setSecondsLeft(0);
    if (timerRef.current) clearInterval(timerRef.current);
  };

  return (
    <div className="space-y-4">
      {/* Casting Timer */}
      <div>
        <h3 className="text-md font-semibold mb-2">Casting Timer</h3>
        <div className="flex items-center justify-between">
          <div className="text-center">
            <p className="text-3xl font-display text-orange-400">{timerDuration}s</p>
            <p className="text-micro text-muted-foreground uppercase">Window Duration</p>
          </div>
          <div className="flex flex-col gap-2 items-end">
            {casting && (
              <Badge variant="destructive">CASTING</Badge>
            )}
            {done && !casting && (
              <Badge variant="outline" className="text-muted-foreground">DONE</Badge>
            )}
            {casting && (
              <p className="text-sm text-orange-400 font-mono">{secondsLeft}s remaining</p>
            )}
            {casting || done ? (
              <Button size="sm" variant="outline" onClick={reset}>
                Reset
              </Button>
            ) : (
              <Button size="sm" onClick={startCast}>
                Start Cast
              </Button>
            )}
          </div>
        </div>
      </div>

      <Separator />

      {/* Speed Dial Slots */}
      <div>
        <h3 className="text-md font-semibold mb-2">Speed Dial Slots</h3>
        {speedDialSlots > 0 ? (
          <>
            <div className="flex gap-1.5 mb-2">
              {Array.from({ length: speedDialSlots }).map((_, i) => (
                <div
                  key={i}
                  className={cn(
                    'w-5 h-5 rounded-full border-2 bg-primary border-primary'
                  )}
                />
              ))}
            </div>
            <p className="text-xs text-muted-foreground">
              Pre-save up to {speedDialSlots} component string(s) for instant casting without using the timer.
            </p>
          </>
        ) : (
          <p className="text-xs text-muted-foreground italic">Unlocked at level 3</p>
        )}
      </div>
    </div>
  );
};
