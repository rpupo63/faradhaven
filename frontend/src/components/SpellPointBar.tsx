import { cn } from '@/lib/utils';
import { Progress } from '@/components/ui/progress';
import { RaIcon } from '@/components/ui/RaIcon';

interface SpellPointBarProps {
  current: number;
  max: number;
  className?: string;
  showLabel?: boolean;
}

/**
 * Displays current / max spell points (component-based resource pool).
 */
export function SpellPointBar({
  current,
  max,
  className,
  showLabel = true,
}: SpellPointBarProps) {
  const pct = max > 0 ? Math.min(100, (current / max) * 100) : 0;

  return (
    <div className={cn('space-y-2', className)}>
      {showLabel && (
        <div className="flex items-center justify-between text-sm">
          <span className="flex items-center gap-1.5 font-tome-subheading text-primary">
            <RaIcon name="aura" className="text-sm text-primary" />
            Spell Points
          </span>
          <span className="font-tome-marginalia text-muted-foreground">
            {current} / {max}
          </span>
        </div>
      )}
      <Progress value={pct} className="h-3 arcane-border" />
    </div>
  );
}
