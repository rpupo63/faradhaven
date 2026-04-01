import { Brain } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { cn } from '@/lib/utils';
import type { NormalizedCharacterSheet } from '@/types/game';

interface MutagenStatusCardProps {
  sheet: NormalizedCharacterSheet;
  onClick?: () => void;
  className?: string;
}

export function MutagenStatusCard({ sheet, onClick, className }: MutagenStatusCardProps) {
  const resources = sheet.class_resources ?? [];
  const madnessDC = resources.find(r => r.key === 'madness_base_dc')?.current_value ?? 10;
  const feralBonus = resources.find(r => r.key === 'feral_bonus')?.current_value ?? 0;

  return (
    <Card
      className={cn('arcane-border bg-card', onClick && 'cursor-pointer hover:ring-1 hover:ring-primary/50 transition-all', className)}
      onClick={onClick}
    >
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center justify-between text-base font-tome-subheading text-primary">
          <div className="flex items-center gap-2">
            <Brain className="h-4 w-4" />
            Mutagen State
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-around text-center">
          <div>
            <p className="text-2xl font-display text-red-400">{madnessDC}</p>
            <p className="text-micro font-tome-marginalia text-muted-foreground uppercase">Madness DC</p>
          </div>
          <div>
            <p className="text-2xl font-display text-orange-400">+{feralBonus}</p>
            <p className="text-micro font-tome-marginalia text-muted-foreground uppercase">Feral Bonus</p>
          </div>
        </div>
        <p className="text-micro font-tome-marginalia text-muted-foreground text-center italic mt-2">
          Click to manage Feral Mechanics
        </p>
      </CardContent>
    </Card>
  );
}
