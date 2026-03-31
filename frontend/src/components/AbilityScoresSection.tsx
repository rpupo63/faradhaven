import { cn } from '@/lib/utils';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { NormalizedCharacterSheet } from '@/types/game';
import { BarChart3 } from 'lucide-react';

const ABILITIES = [
  { id: 'strength', label: 'STR', full: 'Strength' },
  { id: 'dexterity', label: 'DEX', full: 'Dexterity' },
  { id: 'constitution', label: 'CON', full: 'Constitution' },
  { id: 'intelligence', label: 'INT', full: 'Intelligence' },
  { id: 'wisdom', label: 'WIS', full: 'Wisdom' },
  { id: 'charisma', label: 'CHA', full: 'Charisma' },
] as const;

interface AbilityScoresSectionProps {
  sheet: NormalizedCharacterSheet;
  className?: string;
}

export function AbilityScoresSection({ sheet, className }: AbilityScoresSectionProps) {
  const { character, class: cls, modifiers } = sheet;
  const primaryAbility = cls.primary_ability.toLowerCase();

  return (
    <Card className={cn('arcane-border bg-card', className)}>
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center gap-2 text-base font-tome-subheading text-primary">
          <BarChart3 className="h-4 w-4" />
          Ability Scores
        </CardTitle>
        <p className="text-xs text-muted-foreground font-tome-marginalia">
          Spellcasting: {cls.primary_ability}
        </p>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-6 gap-3">
          {ABILITIES.map(({ id, label, full }) => {
            const score =
              id === 'strength'
                ? character.strength
                : id === 'dexterity'
                  ? character.dexterity
                  : id === 'constitution'
                    ? character.constitution
                    : id === 'intelligence'
                      ? character.intelligence
                      : id === 'wisdom'
                        ? character.wisdom
                        : character.charisma;
            const mod = modifiers[id];
            const isPrimary = primaryAbility === id;

            return (
              <div
                key={id}
                className={cn(
                  'rounded-lg border p-2 text-center transition-colors',
                  isPrimary
                    ? 'border-primary/50 bg-primary/5'
                    : 'border-border bg-muted/20'
                )}
              >
                <p className="text-xs font-tome-marginalia text-muted-foreground uppercase tracking-wider">
                  {label}
                  {isPrimary && (
                    <span className="ml-1 text-primary" title="Spellcasting ability">
                      *
                    </span>
                  )}
                </p>
                <p className="font-display text-xl text-primary">{score}</p>
                <p
                  className={cn(
                    'text-sm font-tome-marginalia',
                    mod >= 0 ? 'text-primary' : 'text-muted-foreground'
                  )}
                >
                  {mod >= 0 ? '+' : ''}
                  {mod}
                </p>
              </div>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}
