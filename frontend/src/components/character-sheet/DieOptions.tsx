import { useEffect } from 'react';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { dispatchClearDice } from '@/lib/dice';

interface DieOptionsProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onRoll: (sides: number) => void;
}

export function DieOptions({ open, onOpenChange, onRoll }: DieOptionsProps) {
  useEffect(() => {
    if (!open) {
      dispatchClearDice();
    }
  }, [open]);

  const dice = [4, 6, 8, 10, 12, 20, 100];

  const handleDieRoll = (sides: number) => {
    onRoll(sides);
    onOpenChange(false); // Close the dialog after rolling
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-xs">
        <DialogHeader>
          <DialogTitle className="text-center text-primary font-display">Roll Any Die</DialogTitle>
          <DialogDescription className="sr-only">Select a die to roll.</DialogDescription>
        </DialogHeader>
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