import { useState, useEffect } from 'react';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { dispatchClearDice } from '@/lib/dice';
import { Minus, Plus } from 'lucide-react';

interface DieOptionsProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onRoll: (count: number, sides: number) => void;
  /** At 0 HP: offer Faradhaven death save (5d20, best counts). */
  deathSaveOption?: {
    show: boolean;
    disabled: boolean;
    disabledReason?: string;
    busy?: boolean;
    /** Return true after dice finish so the dialog can close; false if cancelled. */
    onRollDeathSave: () => boolean | Promise<boolean>;
  };
}

export function DieOptions({ open, onOpenChange, onRoll, deathSaveOption }: DieOptionsProps) {
  const [count, setCount] = useState(1);

  useEffect(() => {
    if (!open) {
      dispatchClearDice();
      setCount(1);
    }
  }, [open]);

  const dice = [4, 6, 8, 10, 12, 20, 100];

  const handleDieRoll = (sides: number) => {
    onRoll(count, sides);
    onOpenChange(false); // Close the dialog after rolling
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-xs">
        <DialogHeader>
          <DialogTitle className="text-center text-primary font-display">Roll Any Die</DialogTitle>
          <DialogDescription className="sr-only">Select a die to roll.</DialogDescription>
        </DialogHeader>
        {deathSaveOption?.show && (
          <div className="px-4 pt-2 pb-1 border-b border-border/60">
            <p className="text-[10px] font-tome-marginalia text-muted-foreground uppercase tracking-wide mb-1.5">
              Dying (0 HP)
            </p>
            <Button
              type="button"
              variant="destructive"
              className="w-full h-11 font-tome-subheading text-sm"
              disabled={deathSaveOption.disabled || deathSaveOption.busy}
              title={deathSaveOption.disabled ? deathSaveOption.disabledReason : undefined}
              onClick={async () => {
                const finished = await deathSaveOption.onRollDeathSave();
                if (finished) onOpenChange(false);
              }}
            >
              {deathSaveOption.busy ? 'Rolling…' : 'Death save — 5d20, keep best'}
            </Button>
            <p className="text-[10px] text-muted-foreground font-tome-marginalia mt-1.5 leading-snug">
              Best die 11+ = success, 10 or below = failure (track on sheet).
            </p>
          </div>
        )}

        <div className="px-4 py-3 border-b border-border/60 flex items-center justify-between">
          <span className="text-xs font-tome-subheading text-muted-foreground uppercase tracking-wider">Quantity</span>
          <div className="flex items-center gap-3">
            <Button
              variant="outline"
              size="icon"
              className="h-8 w-8"
              onClick={() => setCount(Math.max(1, count - 1))}
              disabled={count <= 1}
            >
              <Minus className="h-3 w-3" />
            </Button>
            <span className="text-sm font-bold min-w-[1.2rem] text-center">{count}</span>
            <Button
              variant="outline"
              size="icon"
              className="h-8 w-8"
              onClick={() => setCount(Math.min(20, count + 1))}
              disabled={count >= 20}
            >
              <Plus className="h-3 w-3" />
            </Button>
          </div>
        </div>

        <div className="grid grid-cols-3 gap-2 p-4">
          {dice.map((sides) => (
            <Button
              key={sides}
              variant="outline"
              className="text-lg h-12"
              onClick={() => handleDieRoll(sides)}
            >
              D{sides}
            </Button>
          ))}
        </div>
      </DialogContent>
    </Dialog>
  );
}