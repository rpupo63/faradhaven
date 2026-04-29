import { Moon, Sun } from 'lucide-react';
import { RaIcon } from '@/components/ui/RaIcon';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

interface MobileSheetBottomBarProps {
  onShortRest?: () => void | Promise<void>;
  onLongRest?: () => void | Promise<void>;
  onOpenDice: () => void;
  className?: string;
}

/**
 * Fixed bottom actions on small screens: short rest, long rest, roll any die.
 */
export function MobileSheetBottomBar({
  onShortRest,
  onLongRest,
  onOpenDice,
  className,
}: MobileSheetBottomBarProps) {
  return (
    <div
      className={cn(
        'md:hidden fixed bottom-0 left-0 right-0 z-40 border-t border-border/60 bg-card/95 backdrop-blur-md px-2 pt-2 pb-[calc(0.5rem+env(safe-area-inset-bottom))] shadow-[0_-4px_20px_rgba(0,0,0,0.08)]',
        className
      )}
      role="toolbar"
      aria-label="Sheet quick actions"
    >
      <div className="flex items-stretch gap-2 max-w-lg mx-auto">
        {onShortRest && (
          <Button
            variant="outline"
            size="sm"
            className="h-11 flex-1 min-w-0 flex-col gap-0 py-1.5 px-1"
            onClick={async () => {
              await onShortRest();
            }}
          >
            <Moon className="h-4 w-4 shrink-0" />
            <span className="text-micro font-tome-marginalia leading-tight">Short rest</span>
          </Button>
        )}
        {onLongRest && (
          <Button
            variant="outline"
            size="sm"
            className="h-11 flex-1 min-w-0 flex-col gap-0 py-1.5 px-1"
            onClick={() => void onLongRest()}
          >
            <Sun className="h-4 w-4 shrink-0" />
            <span className="text-micro font-tome-marginalia leading-tight">Long rest</span>
          </Button>
        )}
        <Button
          variant="secondary"
          size="sm"
          className="h-11 flex-1 min-w-0 flex-col gap-0 py-1.5 px-1"
          onClick={onOpenDice}
        >
          <RaIcon name="perspective-dice-six" className="text-sm shrink-0" />
          <span className="text-micro font-tome-marginalia leading-tight text-center">
            Roll any die
          </span>
        </Button>
      </div>
    </div>
  );
}
