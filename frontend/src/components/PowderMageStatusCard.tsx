import { Timer } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { cn } from '@/lib/utils';
import type { NormalizedCharacterSheet } from '@/types/game';

interface PowderMageStatusCardProps {
  sheet: NormalizedCharacterSheet;
  onClick?: () => void;
  className?: string;
}

export function PowderMageStatusCard({ sheet, onClick, className }: PowderMageStatusCardProps) {
  const resources = sheet.class_resources ?? [];
  const timerDuration = resources.find(r => r.key === 'timer_duration')?.value ?? 2;
  const speedDialSlots = resources.find(r => r.key === 'speed_dial_slots')?.value ?? 0;

  return (
    <Card
      className={cn('arcane-border bg-card', onClick && 'cursor-pointer hover:ring-1 hover:ring-primary/50 transition-all', className)}
      onClick={onClick}
    >
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center justify-between text-base font-tome-subheading text-primary">
          <div className="flex items-center gap-2">
            <Timer className="h-4 w-4" />
            Casting Window
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-around text-center">
          <div>
            <p className="text-2xl font-display text-orange-400">{timerDuration}s</p>
            <p className="text-[10px] font-tome-marginalia text-muted-foreground uppercase">Casting Timer</p>
          </div>
          <div>
            <p className="text-2xl font-display text-primary">{speedDialSlots}</p>
            <p className="text-[10px] font-tome-marginalia text-muted-foreground uppercase">Speed Dial Slots</p>
          </div>
        </div>
        <p className="text-[10px] font-tome-marginalia text-muted-foreground text-center italic mt-2">
          Click to manage Continuous Ignition
        </p>
      </CardContent>
    </Card>
  );
}
