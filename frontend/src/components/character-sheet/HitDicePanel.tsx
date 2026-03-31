import { Heart } from 'lucide-react';
import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { NormalizedCharacterSheet } from '@/types/game';
import { rollHitDice, dispatchClearDice } from '@/lib/dice';

interface HitDicePanelProps {
  sheet: NormalizedCharacterSheet;
  onUseHitDice: (rolls: number[]) => void | Promise<void>;
}

export function HitDicePanel({ sheet, onUseHitDice }: HitDicePanelProps) {
  const {
    hit_dice_remaining: hitDiceRemaining,
    hit_die: hitDie,
    modifiers,
  } = sheet;

  const [hitDiceToUse, setHitDiceToUse] = useState(1);
  const [lastHitDiceRoll, setLastHitDiceRoll] = useState<{ die: number; total: number }[] | null>(null);

  // Clear dice when panel unmounts
  useEffect(() => {
    return () => {
      dispatchClearDice();
    };
  }, []);

  const handleUseHitDice = async () => {
    if (hitDiceToUse < 1 || hitDiceToUse > hitDiceRemaining) return;

    // Roll all hit dice as one animation, then apply CON mod to each.
    const diceRolls = await rollHitDice(hitDiceToUse, hitDie);
    if (!diceRolls) return;

    const results: { die: number; total: number }[] = diceRolls.map(die => ({
      die,
      total: Math.max(1, die + modifiers.constitution),
    }));

    setLastHitDiceRoll(results);
    onUseHitDice(results.map(r => r.total));
    setHitDiceToUse(1);
  };

  return (
    <div className="space-y-4 py-2">
      <p className="text-xs text-muted-foreground text-center font-tome-marginalia">
        Each die heals d{hitDie}{modifiers.constitution >= 0 ? '+' : ''}{modifiers.constitution} HP
      </p>

      {/* Dice Count Selector */}
      <div className="flex items-center justify-center gap-3">
        <Button
          variant="outline"
          size="sm"
          onClick={() => setHitDiceToUse(Math.max(1, hitDiceToUse - 1))}
          disabled={hitDiceToUse <= 1}
        >
          -
        </Button>
        <span className="font-display text-xl text-primary min-w-[2rem] text-center">
          {hitDiceToUse}
        </span>
        <Button
          variant="outline"
          size="sm"
          onClick={() => setHitDiceToUse(Math.min(hitDiceRemaining, hitDiceToUse + 1))}
          disabled={hitDiceToUse >= hitDiceRemaining}
        >
          +
        </Button>
      </div>

      {/* Roll Button */}
      <Button
        onClick={handleUseHitDice}
        disabled={hitDiceRemaining === 0}
        className="w-full gap-2"
        variant="default"
      >
        <Heart className="h-4 w-4" />
        Roll {hitDiceToUse} Hit {hitDiceToUse === 1 ? 'Die' : 'Dice'}
      </Button>

      {/* Last Roll Results */}
      {lastHitDiceRoll && lastHitDiceRoll.length > 0 && (
        <div className="pt-2 border-t border-border/50 text-center h-12 flex items-center justify-center">
          {/* Blank result area as requested */}
        </div>
      )}

      {hitDiceRemaining === 0 && (
        <p className="text-center text-xs text-muted-foreground">
          Take a long rest to recover hit dice
        </p>
      )}
    </div>
  );
}
