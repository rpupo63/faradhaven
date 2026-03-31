import { cn } from '@/lib/utils';

interface DiceQuickRollProps {
  onRoll: (sides: number) => void;
  className?: string;
}

const DICE_TYPES = [4, 6, 8, 10, 12, 20, 100];

export function DiceQuickRoll({ onRoll, className }: DiceQuickRollProps) {
  return (
    <div className={cn("grid grid-cols-2 gap-3", className)}>
      {DICE_TYPES.map((sides) => (
        <button
          key={sides}
          onClick={() => onRoll(sides)}
          aria-label={`Roll d${sides}`}
          className={cn(
            'flex items-center justify-center py-4 px-2 rounded border border-primary/20 bg-primary/5 transition-all duration-200',
            'hover:bg-primary/20 hover:border-primary/50 hover:scale-105 active:scale-95'
          )}
        >
          <span className="font-display text-primary text-2xl">d{sides}</span>
        </button>
      ))}
    </div>
  );
}
