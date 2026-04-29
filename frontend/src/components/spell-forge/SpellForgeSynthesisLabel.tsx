import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { Info } from 'lucide-react';
import { cn } from '@/lib/utils';
import { SPELL_SYNTHESIS_HINT_BODY } from '@/lib/spellForgeSynthesisHints';

export interface SpellForgeSynthesisHintProps {
  /** When synthesis fills this field automatically (user has not marked this bucket as manual). */
  showAutoBadge: boolean;
  /** Tooltip summary line after “Recommended:”; omit info icon when null/empty */
  recommendedSummary: string | null;
  /** User manually edited this bucket — forge keeps their value on component changes */
  hasManualOverride: boolean;
}

/** Auto badge + info tooltip (use beside a checkbox label or inside `SpellForgeSynthesisLabel`) */
export function SpellForgeSynthesisHintBadges({
  showAutoBadge,
  recommendedSummary,
  hasManualOverride,
}: SpellForgeSynthesisHintProps) {
  const showHint = Boolean(recommendedSummary && recommendedSummary.trim() !== '');

  return (
    <>
      {showAutoBadge && (
        <Badge variant="secondary" size="tiny">
          auto
        </Badge>
      )}
      {showHint && (
        <TooltipProvider delayDuration={200}>
          <Tooltip>
            <TooltipTrigger asChild>
              <button
                type="button"
                className="inline-flex shrink-0 rounded-sm text-muted-foreground hover:text-foreground focus:outline-none focus:ring-2 focus:ring-primary/40"
                aria-label="Synthesis recommendation details"
              >
                <Info className="h-3.5 w-3.5" />
              </button>
            </TooltipTrigger>
            <TooltipContent side="top" className="max-w-sm text-xs font-tome-marginalia">
              <p className="font-semibold text-foreground">Recommended</p>
              <p className="mt-1 text-popover-foreground">{recommendedSummary}</p>
              <p className="mt-2 text-muted-foreground">{SPELL_SYNTHESIS_HINT_BODY}</p>
              {hasManualOverride && (
                <p className="mt-2 border-t border-border/60 pt-2 text-amber-700 dark:text-amber-400">
                  You edited this field manually; the forge keeps your value and will not auto-replace it when the chain
                  changes.
                </p>
              )}
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      )}
    </>
  );
}

export interface SpellForgeSynthesisLabelProps extends SpellForgeSynthesisHintProps {
  className?: string;
  /** Visible label text */
  children: React.ReactNode;
}

/** Shared label row for Spell Forge / Spell Forge 2: optional auto badge + synthesis recommendation tooltip */
export function SpellForgeSynthesisLabel({
  className,
  children,
  showAutoBadge,
  recommendedSummary,
  hasManualOverride,
}: SpellForgeSynthesisLabelProps) {
  return (
    <Label
      className={cn('text-sm font-tome-marginalia text-muted-foreground mb-2 block', className)}
    >
      <span className="inline-flex flex-wrap items-center gap-1.5">
        {children}
        <SpellForgeSynthesisHintBadges
          showAutoBadge={showAutoBadge}
          recommendedSummary={recommendedSummary}
          hasManualOverride={hasManualOverride}
        />
      </span>
    </Label>
  );
}
