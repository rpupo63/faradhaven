import { Zap } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { cn } from '@/lib/utils';
import type { NormalizedCharacterSheet } from '@/types/game';

interface PistonBrawlerStatusCardProps {
  sheet: NormalizedCharacterSheet;
  onClick?: () => void;
  className?: string;
}

export function PistonBrawlerStatusCard({ sheet, onClick, className }: PistonBrawlerStatusCardProps) {
  const resources = sheet.class_resources ?? [];
  const stabilityResource = resources.find(r => r.key === 'max_stability');
  const currentStability = stabilityResource?.current_value ?? stabilityResource?.value ?? 0;
  const maxStability = stabilityResource?.max_value ?? stabilityResource?.value ?? 0;

  const level = sheet.character.level;
  const archetype = sheet.character.archetypeName;
  const overdriveDice = level >= 11
    ? (archetype === 'The Vulcan' && level >= 15 ? 4 : 3)
    : 2;

  return (
    <Card
      className={cn('arcane-border bg-card', onClick && 'cursor-pointer hover:ring-1 hover:ring-primary/50 transition-all', className)}
      onClick={onClick}
    >
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center justify-between text-base font-tome-subheading text-primary">
          <div className="flex items-center gap-2">
            <Zap className="h-4 w-4" />
            Piston Core
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-around text-center">
          <div>
            <p className="text-2xl font-display text-primary">{currentStability}/{maxStability}</p>
            <p className="text-micro font-tome-marginalia text-muted-foreground uppercase">Stability</p>
          </div>
          <div>
            <p className="text-2xl font-display text-orange-400">{overdriveDice}d6</p>
            <p className="text-micro font-tome-marginalia text-muted-foreground uppercase">Overdrive</p>
          </div>
        </div>
        <p className="text-micro font-tome-marginalia text-muted-foreground text-center italic mt-2">
          Click to manage Piston Core
        </p>
      </CardContent>
    </Card>
  );
}
