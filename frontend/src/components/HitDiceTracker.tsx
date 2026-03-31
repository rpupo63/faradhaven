import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Dices, Heart } from 'lucide-react';
import { cn } from '@/lib/utils';
import { rollHitDice } from '@/lib/dice';

interface HitDiceTrackerProps {
  hitDiceTotal: number;
  hitDiceRemaining: number;
  hitDie: number; // e.g., 10 for d10
  conMod: number;
  onUseHitDice: (rolls: number[]) => void | Promise<void>;
  disabled?: boolean;
  className?: string;
}

export function HitDiceTracker({
  hitDiceTotal,
  hitDiceRemaining,
  hitDie,
  conMod,
  onUseHitDice,
  disabled = false,
  className,
}: HitDiceTrackerProps) {
  const [diceToUse, setDiceToUse] = useState(1);
  const [lastRolls, setLastRolls] = useState<{ die: number; total: number }[]>([]);

  const handleUseHitDice = async () => {
    if (diceToUse > hitDiceRemaining || diceToUse < 1) return;

    // Roll all dice as one animation, then apply CON mod to each.
    const diceRolls = await rollHitDice(diceToUse, hitDie);
    if (!diceRolls) return;

    const rolls: { die: number; total: number }[] = diceRolls.map(die => ({
      die,
      total: Math.max(1, die + conMod),
    }));

    setLastRolls(rolls);
    onUseHitDice(rolls.map(r => r.total));
  };

  // Generate visual dice indicators
  const diceIndicators = [];
  for (let i = 0; i < hitDiceTotal; i++) {
    const isAvailable = i < hitDiceRemaining;
    diceIndicators.push(
      <div
        key={i}
        className={cn(
          'w-6 h-6 rounded flex items-center justify-center text-xs font-bold border',
          isAvailable
            ? 'bg-primary/20 border-primary/50 text-primary'
            : 'bg-muted/30 border-muted/50 text-muted-foreground'
        )}
        title={isAvailable ? 'Available' : 'Used'}
      >
        d{hitDie}
      </div>
    );
  }

  return (
    <Card className={cn('arcane-border bg-card', className)}>
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center gap-2 text-base font-tome-subheading text-primary">
          <Dices className="h-4 w-4" />
          Hit Dice
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Hit Dice Display */}
        <div className="text-center">
          <div className="font-display text-2xl text-primary">
            {hitDiceRemaining}
            <span className="text-lg text-muted-foreground">/{hitDiceTotal}</span>
            <span className="text-lg text-muted-foreground ml-2">d{hitDie}</span>
          </div>
          <p className="text-xs text-muted-foreground mt-1">
            Each die heals d{hitDie}{conMod >= 0 ? '+' : ''}{conMod} HP
          </p>
        </div>

        {/* Visual Dice Indicators */}
        <div className="flex flex-wrap justify-center gap-1">
          {diceIndicators}
        </div>

        {/* Dice Usage Controls */}
        <div className="flex items-center justify-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setDiceToUse(Math.max(1, diceToUse - 1))}
            disabled={disabled || diceToUse <= 1}
          >
            -
          </Button>
          <Badge variant="secondary" size="min-w-sm">
            {diceToUse}
          </Badge>
          <Button
            variant="outline"
            size="sm"
            onClick={() => setDiceToUse(Math.min(hitDiceRemaining, diceToUse + 1))}
            disabled={disabled || diceToUse >= hitDiceRemaining}
          >
            +
          </Button>
        </div>

        {/* Use Hit Dice Button */}
        <Button
          onClick={handleUseHitDice}
          disabled={disabled || hitDiceRemaining === 0 || diceToUse < 1}
          className="w-full gap-2"
          variant="default"
        >
          <Heart className="h-4 w-4" />
          Use {diceToUse} Hit {diceToUse === 1 ? 'Die' : 'Dice'}
        </Button>

        {/* Last Roll Results */}
        {lastRolls.length > 0 && (
          <div className="pt-2 border-t border-border/50">
            <p className="text-xs text-muted-foreground mb-2">Last Roll</p>
            <div className="flex flex-wrap gap-2 justify-center">
              {lastRolls.map((roll, idx) => (
                <Badge key={idx} variant="outline" className="gap-1">
                  <span className="text-muted-foreground">d{hitDie}:</span>
                  <span>{roll.die}</span>
                  <span className="text-muted-foreground">
                    {conMod >= 0 ? '+' : ''}{conMod}
                  </span>
                  <span className="text-green-400">=</span>
                  <span className="text-green-400 font-bold">{roll.total}</span>
                </Badge>
              ))}
            </div>
            <p className="text-center text-sm text-green-400 mt-2 font-medium">
              Total: +{lastRolls.reduce((sum, r) => sum + r.total, 0)} HP
            </p>
          </div>
        )}

        {hitDiceRemaining === 0 && (
          <p className="text-center text-xs text-muted-foreground">
            Take a long rest to recover hit dice
          </p>
        )}
      </CardContent>
    </Card>
  );
}
