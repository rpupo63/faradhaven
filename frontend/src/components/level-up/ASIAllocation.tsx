import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ApiLevelUpPreview, AbilityId } from '@/types/game';

interface ASIAllocationProps {
  preview: ApiLevelUpPreview;
  asiAllocation: Record<string, number>;
  onASIChange: (ability: string, delta: number) => void;
  asiPointsRemaining: number;
}

const ABILITIES: { id: AbilityId; name: string }[] = [
  { id: 'strength', name: 'Strength' },
  { id: 'dexterity', name: 'Dexterity' },
  { id: 'constitution', name: 'Constitution' },
  { id: 'intelligence', name: 'Intelligence' },
  { id: 'wisdom', name: 'Wisdom' },
  { id: 'charisma', name: 'Charisma' },
];

export function ASIAllocation({ preview, asiAllocation, onASIChange, asiPointsRemaining }: ASIAllocationProps) {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-tome-subheading text-lg text-primary">Ability Score Improvement</h3>
        <Badge variant={asiPointsRemaining > 0 ? 'secondary' : 'default'}>
          {asiPointsRemaining} points remaining
        </Badge>
      </div>
      <p className="text-sm text-muted-foreground font-tome-marginalia">
        Allocate up to +2 to one ability or +1 to two abilities
      </p>
      <div className="space-y-3">
        {ABILITIES.map((ability) => (
          <div key={ability.id} className="flex items-center justify-between p-3 rounded border border-border/50 bg-muted/30">
            <span className="font-tome-marginalia text-foreground">{ability.name}</span>
            <div className="flex items-center gap-3">
              <Button
                variant="outline"
                size="icon"
                className="h-8 w-8"
                onClick={() => onASIChange(ability.id, -1)}
                disabled={(asiAllocation[ability.id] ?? 0) <= 0}
              >
                -
              </Button>
              <span className="w-8 text-center font-display text-lg text-primary">
                +{asiAllocation[ability.id] ?? 0}
              </span>
              <Button
                variant="outline"
                size="icon"
                className="h-8 w-8"
                onClick={() => onASIChange(ability.id, 1)}
                disabled={
                  asiPointsRemaining <= 0 ||
                  (asiAllocation[ability.id] ?? 0) >= 2
                }
              >
                +
              </Button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
