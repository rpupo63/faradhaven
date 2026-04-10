import { useState } from 'react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible';
import { rollD20, type RollResult } from '@/lib/dice';
import { DND5E_SKILLS, DND5E_SAVING_THROWS } from '@/types/game';
import type { NormalizedCharacterSheet } from '@/types/game';
import { ChevronDown, ChevronRight } from 'lucide-react';
import { RaIcon } from '@/components/ui/RaIcon';

interface SkillsSectionProps {
  sheet: NormalizedCharacterSheet;
  onRollResult?: (result: RollResult) => void;
  className?: string;
}

export function SkillsSection({
  sheet,
  onRollResult,
  className,
}: SkillsSectionProps) {
  const [skillsOpen, setSkillsOpen] = useState(true);
  const [savesOpen, setSavesOpen] = useState(false);
  const [lastRoll, setLastRoll] = useState<RollResult | null>(null);

  const { modifiers, skill_proficiencies, saving_throw_proficiencies } = sheet;

  const handleRoll = async (modifier: number, label: string) => {
    const result = await rollD20(modifier, label);
    if (!result) return;
    const rollResult: RollResult = {
      type: 'ability_check',
      label,
      total: result.total,
      natural: result.natural,
      modifier: result.modifier,
      timestamp:
        // eslint-disable-next-line react-hooks/purity -- click handler, not render
        Date.now(),
    };
    setLastRoll(rollResult);
    onRollResult?.(rollResult);
  };

  const formatMod = (n: number) => (n >= 0 ? `+${n}` : `${n}`);

  return (
    <Card className={cn('arcane-border bg-card', className)}>
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center gap-2 text-base font-tome-subheading text-primary">
          <RaIcon name="book" className="text-sm" />
          Skills & Saving Throws
        </CardTitle>
        <p className="text-xs text-muted-foreground font-tome-marginalia">
          D&D 5e: d20 + ability modifier + proficiency (if proficient)
        </p>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Skills */}
        <Collapsible open={skillsOpen} onOpenChange={setSkillsOpen}>
          <CollapsibleTrigger asChild>
            <Button
              variant="ghost"
              size="sm"
              className="w-full justify-between font-tome-subheading"
            >
              <span>Skills</span>
              {skillsOpen ? (
                <ChevronDown className="h-4 w-4" />
              ) : (
                <ChevronRight className="h-4 w-4" />
              )}
            </Button>
          </CollapsibleTrigger>
          <CollapsibleContent>
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-2 pt-2">
              {DND5E_SKILLS.map((skill) => {
                const bonus = modifiers.skills[skill.id] ?? 0;
                const isProficient = skill_proficiencies.includes(skill.id);
                return (
                  <Button
                    key={skill.id}
                    variant="outline"
                    size="sm"
                    className={cn(
                      'h-auto py-2 px-3 justify-between font-tome-marginalia text-left',
                      isProficient && 'border-primary/50'
                    )}
                    onClick={() =>
                      handleRoll(bonus, `${skill.name} (${formatMod(bonus)})`)
                    }
                  >
                    <span className="truncate">
                      {skill.name}
                      {isProficient && (
                        <span className="text-primary ml-0.5" title="Proficient">
                          ●
                        </span>
                      )}
                    </span>
                    <span className="shrink-0 text-xs tabular-nums">
                      {formatMod(bonus)}
                    </span>
                  </Button>
                );
              })}
            </div>
          </CollapsibleContent>
        </Collapsible>

        {/* Saving Throws */}
        <Collapsible open={savesOpen} onOpenChange={setSavesOpen}>
          <CollapsibleTrigger asChild>
            <Button
              variant="ghost"
              size="sm"
              className="w-full justify-between font-tome-subheading"
            >
              <span className="flex items-center gap-2">
                <RaIcon name="shield" className="text-sm" />
                Saving Throws
              </span>
              {savesOpen ? (
                <ChevronDown className="h-4 w-4" />
              ) : (
                <ChevronRight className="h-4 w-4" />
              )}
            </Button>
          </CollapsibleTrigger>
          <CollapsibleContent>
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-2 pt-2">
              {DND5E_SAVING_THROWS.map((save) => {
                const bonus = modifiers.saving_throws[save.id] ?? 0;
                const isProficient = saving_throw_proficiencies.includes(save.id);
                return (
                  <Button
                    key={save.id}
                    variant="outline"
                    size="sm"
                    className={cn(
                      'h-auto py-2 px-3 justify-between font-tome-marginalia text-left',
                      isProficient && 'border-primary/50'
                    )}
                    onClick={() =>
                      handleRoll(bonus, `${save.name} Save (${formatMod(bonus)})`)
                    }
                  >
                    <span className="truncate">
                      {save.name}
                      {isProficient && (
                        <span className="text-primary ml-0.5" title="Proficient">
                          ●
                        </span>
                      )}
                    </span>
                    <span className="shrink-0 text-xs tabular-nums">
                      {formatMod(bonus)}
                    </span>
                  </Button>
                );
              })}
            </div>
          </CollapsibleContent>
        </Collapsible>

        {/* Last roll from this section */}
        {lastRoll && (
          <div
            className={cn(
              'rounded-lg border p-2 font-display text-center text-sm',
              lastRoll.natural === 20
                ? 'border-green-500/50 bg-green-500/10 text-green-700 dark:text-green-400'
                : lastRoll.natural === 1
                  ? 'border-red-500/50 bg-red-500/10 text-red-700 dark:text-red-400'
                  : 'border-border bg-muted/30'
            )}
          >
            {lastRoll.label}: {lastRoll.total}{' '}
            <span className="text-muted-foreground text-xs">
              (d20: {lastRoll.natural}
              {lastRoll.modifier >= 0 ? ' + ' : ' '}
              {lastRoll.modifier})
            </span>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
