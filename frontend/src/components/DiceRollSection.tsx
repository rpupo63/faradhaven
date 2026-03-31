import { useState } from 'react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  rollD20,
  type RollResult,
  type RollType,
} from '@/lib/dice';
import type { NormalizedCharacterSheet } from '@/types/game';
import { Dices, Sparkles, Swords, Crosshair, Activity } from 'lucide-react';

const ABILITIES = [
  { id: 'strength', label: 'Strength', short: 'STR' },
  { id: 'dexterity', label: 'Dexterity', short: 'DEX' },
  { id: 'constitution', label: 'Constitution', short: 'CON' },
  { id: 'intelligence', label: 'Intelligence', short: 'INT' },
  { id: 'wisdom', label: 'Wisdom', short: 'WIS' },
  { id: 'charisma', label: 'Charisma', short: 'CHA' },
] as const;

type AbilityId = (typeof ABILITIES)[number]['id'];

interface DiceRollSectionProps {
  sheet: NormalizedCharacterSheet;
  className?: string;
}

export function DiceRollSection({ sheet, className }: DiceRollSectionProps) {
  const [lastRoll, setLastRoll] = useState<RollResult | null>(null);
  const [abilityForCheck, setAbilityForCheck] = useState<AbilityId>('strength');

  const { modifiers } = sheet;

  const handleRoll = async (type: RollType, modifier: number, label: string) => {
    const result = await rollD20(modifier);
    if (!result) return;
    setLastRoll({
      type,
      label,
      total: result.total,
      natural: result.natural,
      modifier: result.modifier,
      timestamp: Date.now(),
    });
  };

  const getAbilityMod = (id: AbilityId): number => {
    switch (id) {
      case 'strength':
        return modifiers.strength;
      case 'dexterity':
        return modifiers.dexterity;
      case 'constitution':
        return modifiers.constitution;
      case 'intelligence':
        return modifiers.intelligence;
      case 'wisdom':
        return modifiers.wisdom;
      case 'charisma':
        return modifiers.charisma;
      default:
        return 0;
    }
  };

  const abilityMod = getAbilityMod(abilityForCheck);
  const abilityLabel = ABILITIES.find((a) => a.id === abilityForCheck)?.label ?? abilityForCheck;

  return (
    <Card className={cn('arcane-border bg-card', className)}>
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center gap-2 text-base font-tome-subheading text-primary">
          <Dices className="h-4 w-4" />
          Dice Rolls
        </CardTitle>
        <p className="text-sm text-muted-foreground font-tome-marginalia">
          Roll d20 + modifier for attacks and ability checks
        </p>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Action buttons */}
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-2">
          <Button
            variant="outline"
            size="sm"
            className="gap-2 font-tome-marginalia h-auto py-2"
            onClick={() =>
              handleRoll(
                'spell_attack',
                modifiers.spell_attack,
                `Spell Attack (+${modifiers.spell_attack})`
              )
            }
          >
            <Sparkles className="h-4 w-4 shrink-0" />
            <span className="text-left">
              Spell Attack
              <span className="block text-xs text-muted-foreground">
                +{modifiers.spell_attack}
              </span>
            </span>
          </Button>
          <Button
            variant="outline"
            size="sm"
            className="gap-2 font-tome-marginalia h-auto py-2"
            onClick={() =>
              handleRoll(
                'melee_attack',
                modifiers.melee_attack,
                `Melee Attack (+${modifiers.melee_attack})`
              )
            }
          >
            <Swords className="h-4 w-4 shrink-0" />
            <span className="text-left">
              Melee Attack
              <span className="block text-xs text-muted-foreground">
                +{modifiers.melee_attack}
              </span>
            </span>
          </Button>
          <Button
            variant="outline"
            size="sm"
            className="gap-2 font-tome-marginalia h-auto py-2"
            onClick={() =>
              handleRoll(
                'ranged_attack',
                modifiers.ranged_attack,
                `Ranged Attack (+${modifiers.ranged_attack})`
              )
            }
          >
            <Crosshair className="h-4 w-4 shrink-0" />
            <span className="text-left">
              Ranged Attack
              <span className="block text-xs text-muted-foreground">
                +{modifiers.ranged_attack}
              </span>
            </span>
          </Button>
          <div className="flex flex-col gap-2">
            <Select
              value={abilityForCheck}
              onValueChange={(v) => setAbilityForCheck(v as AbilityId)}
            >
              <SelectTrigger className="h-8 text-xs font-tome-marginalia">
                <SelectValue placeholder="Ability" />
              </SelectTrigger>
              <SelectContent>
                {ABILITIES.map((a) => (
                  <SelectItem key={a.id} value={a.id}>
                    {a.short} ({modifiers[a.id] >= 0 ? '+' : ''}
                    {modifiers[a.id]})
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Button
              variant="outline"
              size="sm"
              className="gap-2 font-tome-marginalia h-auto py-2"
              onClick={() =>
                handleRoll(
                  'ability_check',
                  abilityMod,
                  `${abilityLabel} Check (+${abilityMod >= 0 ? '' : ''}${abilityMod})`
                )
              }
            >
              <Activity className="h-4 w-4 shrink-0" />
              <span className="text-left">
                Ability Check
                <span className="block text-xs text-muted-foreground">
                  {abilityLabel} +{abilityMod >= 0 ? '' : ''}{abilityMod}
                </span>
              </span>
            </Button>
          </div>
        </div>

        {/* Last roll result */}
        {lastRoll && (
          <div
            className={cn(
              'rounded-lg border p-3 font-display text-center transition-colors',
              lastRoll.natural === 20
                ? 'border-green-500/50 bg-green-500/10 text-green-700 dark:text-green-400'
                : lastRoll.natural === 1
                  ? 'border-red-500/50 bg-red-500/10 text-red-700 dark:text-red-400'
                  : 'border-border bg-muted/30'
            )}
          >
            <p className="text-xs text-muted-foreground font-tome-marginalia mb-1">
              {lastRoll.label}
            </p>
            <p className="text-2xl text-primary">
              {lastRoll.total}
              <span className="text-sm font-normal text-muted-foreground ml-1">
                (d20: {lastRoll.natural}
                {lastRoll.modifier >= 0 ? ' + ' : ' '}
                {lastRoll.modifier})
              </span>
            </p>
            {(lastRoll.natural === 20 || lastRoll.natural === 1) && (
              <p className="text-xs font-tome-marginalia mt-1">
                {lastRoll.natural === 20 ? 'Critical hit!' : 'Critical miss!'}
              </p>
            )}
          </div>
        )}

        {/* Save DC reminder */}
        <p className="text-xs text-muted-foreground font-tome-marginalia border-t border-border pt-3">
          Spell Save DC: {sheet.save_dc} (target rolls d20 + their save modifier)
        </p>
      </CardContent>
    </Card>
  );
}
