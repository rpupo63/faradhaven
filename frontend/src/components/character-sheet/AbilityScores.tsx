import { cn } from '@/lib/utils';
import { NormalizedCharacterSheet } from '@/types/game';

interface AbilityScoresProps {
  sheet: NormalizedCharacterSheet;
}

export const ABILITIES = [
  { id: 'strength', label: 'STR', full: 'Strength' },
  { id: 'dexterity', label: 'DEX', full: 'Dexterity' },
  { id: 'constitution', label: 'CON', full: 'Constitution' },
  { id: 'intelligence', label: 'INT', full: 'Intelligence' },
  { id: 'wisdom', label: 'WIS', full: 'Wisdom' },
  { id: 'charisma', label: 'CHA', full: 'Charisma' },
] as const;

const formatMod = (n: number) => (n >= 0 ? `+${n}` : `${n}`);

export function AbilityScores({ sheet }: AbilityScoresProps) {
  const { character, modifiers } = sheet;
  const primaryAbility = sheet.class.primary_ability.toLowerCase();

  const getAbilityScore = (id: string): number => {
    switch (id) {
      case 'strength': return character.strength;
      case 'dexterity': return character.dexterity;
      case 'constitution': return character.constitution;
      case 'intelligence': return character.intelligence;
      case 'wisdom': return character.wisdom;
      case 'charisma': return character.charisma;
      default: return 10;
    }
  };

  return (
    <div className="grid grid-cols-3 gap-1 md:grid-cols-1 md:gap-2">
      {ABILITIES.map(({ id, label }) => {
        const score = getAbilityScore(id);
        const mod = modifiers[id as keyof typeof modifiers] as number;
        const isPrimary = primaryAbility === id;

        return (
          <div
            key={id}
            className={cn(
              'rounded-md md:rounded-lg border px-1 py-1.5 text-center md:p-2',
              isPrimary
                ? 'border-primary/50 bg-primary/5'
                : 'border-border bg-muted/20'
            )}
          >
            <p className="text-micro font-tome-marginalia text-muted-foreground uppercase tracking-wide leading-none md:text-xs md:tracking-wider">
              {label}
              {isPrimary && <span className="ml-0.5 text-primary">*</span>}
            </p>
            <p className="font-display text-lg leading-tight text-primary md:text-xl">{score}</p>
            <p
              className={cn(
                'text-fine font-tome-marginalia leading-none md:text-sm',
                mod >= 0 ? 'text-primary' : 'text-muted-foreground'
              )}
            >
              {formatMod(mod)}
            </p>
          </div>
        );
      })}
    </div>
  );
}
