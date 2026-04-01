import type { ReactNode } from 'react';
import { cn } from '@/lib/utils';
import type { NormalizedCharacterSheet } from '@/types/game';

interface MobileSheetSummaryStripProps {
  sheet: NormalizedCharacterSheet;
  onHPTap?: () => void;
  onInitiativeTap?: () => void;
}

const formatMod = (n: number) => (n >= 0 ? `+${n}` : `${n}`);

/**
 * Sticky at-a-glance row for narrow viewports: HP, AC, speed, proficiency, initiative.
 */
export function MobileSheetSummaryStrip({
  sheet,
  onHPTap,
  onInitiativeTap,
}: MobileSheetSummaryStripProps) {
  const pb = sheet.class_level?.proficiency_bonus ?? 2;
  const init = sheet.modifiers.initiative;
  const speed = sheet.speed;

  const cell = (
    label: string,
    value: ReactNode,
    opts?: { onClick?: () => void; emphasize?: boolean }
  ) => (
    <button
      type="button"
      disabled={!opts?.onClick}
      onClick={opts?.onClick}
      className={cn(
        'flex flex-col items-center justify-center rounded-md border border-border/60 bg-card/90 px-1 py-1.5 text-center min-w-0',
        opts?.onClick && 'active:bg-primary/10 cursor-pointer',
        !opts?.onClick && 'cursor-default'
      )}
    >
      <span className="text-tiny font-tome-marginalia uppercase tracking-wider text-muted-foreground leading-none">
        {label}
      </span>
      <span
        className={cn(
          'text-sm font-display tabular-nums leading-tight mt-0.5',
          opts?.emphasize ? 'text-primary' : 'text-foreground'
        )}
      >
        {value}
      </span>
    </button>
  );

  return (
    <div
      className="sticky top-0 z-30 -mx-1 px-1 py-2 mb-2 border-b border-border/50 bg-background/95 backdrop-blur-sm md:hidden"
      role="region"
      aria-label="Quick stats"
    >
      <div className="grid grid-cols-5 gap-1">
        {cell(
          'HP',
          <>
            {sheet.current_hp}
            <span className="text-muted-foreground font-normal">/{sheet.max_hp}</span>
          </>,
          onHPTap ? { onClick: onHPTap, emphasize: true } : undefined
        )}
        {cell('AC', sheet.ac)}
        {cell('SPD', speed != null ? `${speed}` : '—')}
        {cell('Prof', formatMod(pb))}
        {cell(
          'Init',
          formatMod(init),
          onInitiativeTap
            ? {
                onClick: onInitiativeTap,
                emphasize: true,
              }
            : undefined
        )}
      </div>
    </div>
  );
}
