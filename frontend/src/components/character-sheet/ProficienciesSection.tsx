import { useState } from 'react';
import { Shield, BookOpen, Zap, Dices } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { NormalizedCharacterSheet, DND5E_SKILLS, DND5E_SAVING_THROWS } from '@/types/game';

interface ProficienciesSectionProps {
  sheet: NormalizedCharacterSheet;
  onRoll: (label: string, modifier: number) => void;
}

const formatMod = (n: number) => (n >= 0 ? `+${n}` : `${n}`);

export function ProficienciesSection({ sheet, onRoll }: ProficienciesSectionProps) {
  const {
    modifiers,
    skill_proficiencies,
    saving_throw_proficiencies,
  } = sheet;

  return (
    <div className="space-y-4">
      {/* Initiative */}
      <Card className="arcane-border bg-card">
        <CardContent className="p-0">
          <button
            onClick={() => onRoll('Initiative', modifiers.initiative)}
            className={cn(
              'w-full flex items-center justify-between py-2.5 px-3 rounded transition-colors',
              'hover:bg-primary/20 cursor-pointer'
            )}
          >
            <span className="flex items-center gap-2 text-sm font-tome-subheading text-primary">
              <Zap className="h-4 w-4" />
              Initiative
            </span>
            <span className="font-display text-xl text-primary tabular-nums">
              {formatMod(modifiers.initiative)}
            </span>
          </button>
        </CardContent>
      </Card>

      {/* Saving Throws */}
      <Card className="arcane-border bg-card">
        <CardHeader className="pb-2 px-3 pt-3">
          <CardTitle className="flex items-center gap-2 text-sm font-tome-subheading text-primary">
            <Shield className="h-4 w-4" />
            Saving Throws
          </CardTitle>
        </CardHeader>
        <CardContent className="px-3 pb-3">
          <div className="space-y-1">
            {DND5E_SAVING_THROWS.map((save) => {
              const bonus = modifiers.saving_throws[save.id] ?? 0;
              const isProficient = saving_throw_proficiencies.includes(save.id);
              return (
                <button
                  key={save.id}
                  onClick={() => onRoll(`${save.name} Save`, bonus)}
                  className={cn(
                    'w-full flex items-center justify-between py-1 px-2 rounded text-sm transition-colors',
                    'hover:bg-primary/20 cursor-pointer',
                    isProficient && 'bg-primary/10'
                  )}
                >
                  <span className="flex items-center gap-1.5 font-tome-marginalia">
                    {isProficient && (
                      <span className="w-1.5 h-1.5 rounded-full bg-primary" title="Proficient" />
                    )}
                    {!isProficient && <span className="w-1.5 h-1.5" />}
                    {save.name}
                  </span>
                  <span className="font-display text-primary tabular-nums">
                    {formatMod(bonus)}
                  </span>
                </button>
              );
            })}
          </div>
        </CardContent>
      </Card>

      {/* Skills List */}
      <Card className="arcane-border bg-card">
        <CardHeader className="pb-2 px-3 pt-3">
          <CardTitle className="flex items-center gap-2 text-sm font-tome-subheading text-primary">
            <BookOpen className="h-4 w-4" />
            Skills
          </CardTitle>
        </CardHeader>
        <CardContent className="px-3 pb-3">
          <div className="space-y-1">
            {DND5E_SKILLS.map((skill) => {
              const bonus = modifiers.skills[skill.id] ?? 0;
              const isProficient = skill_proficiencies.includes(skill.id);
              return (
                <button
                  key={skill.id}
                  onClick={() => onRoll(skill.name, bonus)}
                  className={cn(
                    'w-full flex items-center justify-between py-1 px-2 rounded text-sm transition-colors',
                    'hover:bg-primary/20 cursor-pointer',
                    isProficient && 'bg-primary/10'
                  )}
                >
                  <span className="flex items-center gap-1.5 font-tome-marginalia">
                    {isProficient && (
                      <span className="w-1.5 h-1.5 rounded-full bg-primary" title="Proficient" />
                    )}
                    {!isProficient && <span className="w-1.5 h-1.5" />}
                    <span className="text-xs text-muted-foreground uppercase tracking-wide">
                      {skill.ability.slice(0, 3)}
                    </span>
                    {skill.name}
                  </span>
                  <span className="font-display text-primary tabular-nums">
                    {formatMod(bonus)}
                  </span>
                </button>
              );
            })}
          </div>
        </CardContent>
      </Card>


    </div>
  );
}
