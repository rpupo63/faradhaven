import { Skull, HeartPulse } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { cn } from '@/lib/utils';

interface DeathSavesTrackerProps {
  successes: number;
  failures: number;
  stable: boolean;
  dead: boolean;
  canRoll: boolean;
  onRoll: () => void | Promise<void>;
  rolling?: boolean;
}

function Pips({
  count,
  filled,
  variant,
}: {
  count: number;
  filled: number;
  variant: 'success' | 'failure';
}) {
  const filledClass =
    variant === 'success'
      ? 'bg-emerald-600 border-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.35)]'
      : 'bg-red-700 border-red-600 shadow-[0_0_8px_rgba(185,28,28,0.35)]';
  const emptyClass =
    variant === 'success'
      ? 'border-emerald-700/50 bg-emerald-950/20'
      : 'border-red-900/50 bg-red-950/20';

  return (
    <div className="flex gap-1.5" role="img" aria-label={`${filled} of ${count} ${variant === 'success' ? 'successes' : 'failures'}`}>
      {Array.from({ length: count }, (_, i) => (
        <span
          key={i}
          className={cn(
            'h-3 w-3 rounded-full border-2 transition-colors',
            i < filled ? filledClass : emptyClass
          )}
        />
      ))}
    </div>
  );
}

export function DeathSavesTracker({
  successes,
  failures,
  stable,
  dead,
  canRoll,
  onRoll,
  rolling = false,
}: DeathSavesTrackerProps) {
  return (
    <Card className="arcane-border border-red-900/40 bg-red-950/15">
      <CardContent className="p-3 space-y-3">
        <div className="flex items-start justify-between gap-2">
          <div>
            <div className="flex items-center gap-2">
              <Skull className="h-4 w-4 text-red-400 shrink-0" aria-hidden />
              <h3 className="font-tome-subheading text-sm text-red-200">Death saving throws</h3>
            </div>
            <p className="text-[11px] text-muted-foreground font-tome-marginalia mt-1 leading-snug">
              Roll 5d20 and use the <span className="text-foreground/90">highest</span>. 11+ counts as a
              success; 10 or below is a failure. Three successes: you stabilize. Three failures: you die.
            </p>
          </div>
        </div>

        {dead ? (
          <p className="text-sm text-red-300 font-tome-subheading">This character has died.</p>
        ) : stable ? (
          <div className="flex items-center gap-2 text-sm text-emerald-200/95">
            <HeartPulse className="h-4 w-4 shrink-0 text-emerald-400" aria-hidden />
            <span className="font-tome-marginalia leading-snug">
              Stabilized at 0 HP — unconscious until healed.
            </span>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-1">
                <span className="text-[10px] font-tome-marginalia uppercase tracking-wide text-emerald-600/90">
                  Successes
                </span>
                <Pips count={3} filled={successes} variant="success" />
              </div>
              <div className="space-y-1">
                <span className="text-[10px] font-tome-marginalia uppercase tracking-wide text-red-400/90">
                  Failures
                </span>
                <Pips count={3} filled={failures} variant="failure" />
              </div>
            </div>
            {canRoll && (
              <Button
                type="button"
                variant="destructive"
                className="w-full font-tome-subheading"
                disabled={rolling}
                onClick={() => void onRoll()}
              >
                {rolling ? 'Rolling…' : 'Roll death save (5d20, best)'}
              </Button>
            )}
          </>
        )}
      </CardContent>
    </Card>
  );
}
