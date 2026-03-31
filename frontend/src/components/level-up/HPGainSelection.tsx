import { useMemo } from 'react';
import { Heart, Dices } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ApiLevelUpPreview } from '@/types/game';
import { cn } from '@/lib/utils';
import { rollHitDie } from '@/lib/dice';

interface HPGainSelectionProps {
  preview: ApiLevelUpPreview;
  hpChoice: 'average' | 'roll';
  setHpChoice: (choice: 'average' | 'roll') => void;
  hpRollResult: number | null;
  setHpRollResult: (result: number) => void;
}

export function HPGainSelection({ preview, hpChoice, setHpChoice, hpRollResult, setHpRollResult }: HPGainSelectionProps) {
  const handleRollHP = async () => {
    const roll = await rollHitDie(preview.hit_die, 'Hit Point Increase');
    if (roll === null) return;
    setHpRollResult(roll);
  };

  const hpGain = useMemo(() => {
    if (hpChoice === 'average') {
      return preview.hp_gain_average;
    }
    if (hpRollResult !== null) {
      const total = hpRollResult + preview.con_mod;
      return Math.max(1, total);
    }
    return 0;
  }, [preview, hpChoice, hpRollResult]);

  return (
    <div className="space-y-4">
      <h3 className="font-tome-subheading text-lg text-primary flex items-center gap-2">
        <Heart className="h-5 w-5" />
        Hit Point Increase
      </h3>
      <p className="text-sm text-muted-foreground font-tome-marginalia">
        Choose how to determine your HP gain for this level
      </p>

      <div className="grid gap-4 md:grid-cols-2">
        {/* Average Option */}
        <button
          type="button"
          onClick={() => {
            setHpChoice('average');
          }}
          className={cn(
            'p-4 rounded border text-left transition-all',
            hpChoice === 'average'
              ? 'border-primary bg-primary/10 ring-2 ring-primary'
              : 'border-border/50 bg-muted/30 hover:border-primary/50'
          )}
        >
          <p className="font-tome-subheading text-foreground mb-1">Take Average</p>
          <p className="text-sm text-muted-foreground mb-3">Safe, consistent HP gain</p>
          <div className="flex items-center gap-2">
            <Badge variant="secondary" className="text-lg px-3 py-1">
              +{preview.hp_gain_average} HP
            </Badge>
            <span className="text-xs text-muted-foreground">
              (d{preview.hit_die}/2 + 1 + {preview.con_mod >= 0 ? '+' : ''}{preview.con_mod})
            </span>
          </div>
        </button>

        {/* Roll Option */}
        <button
          type="button"
          onClick={() => setHpChoice('roll')}
          className={cn(
            'p-4 rounded border text-left transition-all',
            hpChoice === 'roll'
              ? 'border-primary bg-primary/10 ring-2 ring-primary'
              : 'border-border/50 bg-muted/30 hover:border-primary/50'
          )}
        >
          <p className="font-tome-subheading text-foreground mb-1 flex items-center gap-2">
            <Dices className="h-4 w-4" />
            Roll for HP
          </p>
          <p className="text-sm text-muted-foreground mb-3">Risk it for potentially more HP</p>
          <div className="flex items-center gap-2">
            <Badge variant="outline" className="text-lg px-3 py-1">
              d{preview.hit_die} + {preview.con_mod >= 0 ? '+' : ''}{preview.con_mod}
            </Badge>
            <span className="text-xs text-muted-foreground">
              ({Math.max(1, 1 + preview.con_mod)}-{preview.hit_die + preview.con_mod} HP)
            </span>
          </div>
        </button>
      </div>

      {/* Roll UI */}
      {hpChoice === 'roll' && (
        <div className="p-4 rounded border border-border/50 bg-muted/30 space-y-4">
          <div className="flex items-center justify-between">
            <span className="font-tome-marginalia text-foreground">
              Roll your d{preview.hit_die}
            </span>
            <Button onClick={handleRollHP} variant="default" size="sm" className="gap-2">
              <Dices className="h-4 w-4" />
              Roll
            </Button>
          </div>
          {hpRollResult !== null && (
            <div className="text-center space-y-2">
              <div className="flex items-center justify-center gap-3">
                <Badge variant="outline" className="text-xl px-4 py-2">
                  d{preview.hit_die}: {hpRollResult}
                </Badge>
                <span className="text-muted-foreground">+</span>
                <Badge variant="outline" className="text-xl px-4 py-2">
                  CON: {preview.con_mod >= 0 ? '+' : ''}{preview.con_mod}
                </Badge>
                <span className="text-muted-foreground">=</span>
                <Badge variant="element-heal" size="xl-compact">
                  +{hpGain} HP
                </Badge>
              </div>
              <p className="text-sm text-muted-foreground">
                New Max HP: {preview.current_max_hp} + {hpGain} = <strong className="text-primary">{preview.current_max_hp + hpGain}</strong>
              </p>
            </div>
          )}
          {hpRollResult === null && (
            <p className="text-center text-sm text-muted-foreground">
              Click Roll to determine your HP gain
            </p>
          )}
        </div>
      )}

      {/* Preview for average */}
      {hpChoice === 'average' && (
        <div className="p-4 rounded border border-green-500/30 bg-green-500/10 text-center">
          <p className="text-sm text-muted-foreground mb-1">New Max HP</p>
          <p className="font-display text-2xl text-green-400">
            {preview.current_max_hp} + {preview.hp_gain_average} = <strong>{preview.current_max_hp + preview.hp_gain_average}</strong>
          </p>
        </div>
      )}
    </div>
  );
}
