import { ShieldAlert } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { cn } from '@/lib/utils';

interface NotorietyMeterProps {
  characterId: string;
  notoriety: number;
  sanguine_mp?: number;
  sanguine_br?: number;
  onClick?: () => void;
  className?: string;
}

export function NotorietyMeter({ notoriety, sanguine_mp, sanguine_br, onClick, className }: NotorietyMeterProps) {
  // Notoriety is from -20 to 20.
  // We want to map it to 0 to 100 for the progress bar.
  // formula: ((notoriety + 20) / 40) * 100
  const percentage = ((notoriety + 20) / 40) * 100;

  const getStatusColor = () => {
    if (notoriety <= -10) return 'text-blue-500'; // Very Good/Doctor
    if (notoriety < 0) return 'text-blue-300'; // Good
    if (notoriety === 0) return 'text-muted-foreground'; // Neutral
    if (notoriety < 10) return 'text-red-300'; // Notorious
    return 'text-red-500'; // Infamous/Predator
  };

  const getLabel = () => {
    if (notoriety <= -15) return 'Medical Prodigy';
    if (notoriety <= -5) return 'Benevolent';
    if (notoriety < 5) return 'Neutral';
    if (notoriety < 15) return 'Starving Predator';
    return 'Blood Rage';
  };

  return (
    <Card className={cn("arcane-border bg-card", onClick && "cursor-pointer hover:ring-1 hover:ring-primary/50 transition-all", className)} onClick={onClick}>
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center justify-between text-base font-tome-subheading text-primary">
          <div className="flex items-center gap-2">
            <ShieldAlert className="h-4 w-4" />
            Notoriety Meter
          </div>
          <span className={cn("text-xs font-tome-marginalia uppercase tracking-widest", getStatusColor())}>
            {getLabel()}
          </span>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="relative pt-2">
          <Progress 
            value={percentage} 
            className="h-3 arcane-border bg-blue-900/20"
          />
          <div className="absolute top-0 left-1/2 -translate-x-1/2 h-5 w-0.5 bg-border z-10" />
        </div>
        
        <div className="flex items-center justify-between gap-4 text-center">
          <div>
            <p className="text-lg font-display text-blue-400">{sanguine_mp ?? 0}</p>
            <p className="text-micro font-tome-marginalia text-muted-foreground uppercase">Medical Prodigy</p>
          </div>
          <div className="text-center">
            <p className="text-2xl font-display text-primary">{notoriety}</p>
            <p className="text-micro font-tome-marginalia text-muted-foreground uppercase">Net Notoriety</p>
          </div>
          <div>
            <p className="text-lg font-display text-red-400">{sanguine_br ?? 0}</p>
            <p className="text-micro font-tome-marginalia text-muted-foreground uppercase">Blood Rage</p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
