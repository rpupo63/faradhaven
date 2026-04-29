import { HelpCircle } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { cn } from '@/lib/utils';
import type { NormalizedCharacterSheet } from '@/types/game';
import {
  normalizeSpellCastingComponent,
  SPELL_CASTING_ABBREV,
  SPELL_CASTING_COMPONENT_LABELS,
  SPELL_CASTING_GUIDANCE,
} from '@/lib/spellCastingStyle';

export interface ClassSpellcastingStyleProps {
  sheet: NormalizedCharacterSheet;
  className?: string;
}

export function ClassSpellcastingStyle({ sheet, className }: ClassSpellcastingStyleProps) {
  const raw = sheet.class.spell_casting_component;
  const kind = normalizeSpellCastingComponent(raw);
  if (!kind) return null;

  const label = SPELL_CASTING_COMPONENT_LABELS[kind];
  const abbr = SPELL_CASTING_ABBREV[kind];
  const guide = SPELL_CASTING_GUIDANCE[kind];
  const flavor = sheet.class.spell_casting_description?.trim();

  return (
    <Card className={cn('arcane-border bg-card/80', className)}>
      <CardHeader className="pb-2 pt-3 px-3 md:px-4">
        <CardTitle className="text-sm font-tome-subheading text-primary flex items-center gap-2">
          <span>Spellcasting style</span>
          <Badge variant="secondary" className="font-tome-marginalia text-xs">
            {abbr} — {label}
          </Badge>
          <TooltipProvider>
            <Tooltip delayDuration={200}>
              <TooltipTrigger asChild>
                <button
                  type="button"
                  className="inline-flex text-muted-foreground hover:text-primary transition-colors"
                  aria-label="Spellcasting style details"
                >
                  <HelpCircle className="h-4 w-4" />
                </button>
              </TooltipTrigger>
              <TooltipContent
                side="top"
                className="max-w-sm p-3 text-left text-popover-foreground"
              >
                <div className="space-y-2 text-xs leading-snug">
                  {flavor ? (
                    <div>
                      <p className="font-tome-subheading text-primary text-[11px] uppercase tracking-wide">
                        How this class casts
                      </p>
                      <p className="mt-0.5 text-foreground/95">{flavor}</p>
                    </div>
                  ) : null}
                  <div>
                    <p className="font-tome-subheading text-primary text-[11px] uppercase tracking-wide">
                      Typical advantages
                    </p>
                    <ul className="mt-0.5 list-disc pl-4 space-y-1">
                      {guide.advantages.map((line) => (
                        <li key={line}>{line}</li>
                      ))}
                    </ul>
                  </div>
                  <div>
                    <p className="font-tome-subheading text-primary text-[11px] uppercase tracking-wide">
                      Typical risks
                    </p>
                    <ul className="mt-0.5 list-disc pl-4 space-y-1">
                      {guide.risks.map((line) => (
                        <li key={line}>{line}</li>
                      ))}
                    </ul>
                  </div>
                </div>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </CardTitle>
      </CardHeader>
      {flavor ? (
        <CardContent className="px-3 pb-3 pt-0 md:px-4 md:pb-4">
          <p className="text-xs font-tome-marginalia text-muted-foreground leading-snug line-clamp-3 md:line-clamp-none">
            {flavor}
          </p>
        </CardContent>
      ) : null}
    </Card>
  );
}
