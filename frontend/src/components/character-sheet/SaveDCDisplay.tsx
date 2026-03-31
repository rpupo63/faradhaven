import { ShieldCheck } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';

interface SaveDCDisplayProps {
  saveDC: number;
}

export function SaveDCDisplay({ saveDC }: SaveDCDisplayProps) {
  if (!saveDC) return null;

  return (
    <Card className="arcane-border bg-card">
      <CardContent className="p-3">
        <div className="flex items-center justify-center gap-4 text-center">
          <div>
            <p className="text-xs font-tome-marginalia text-muted-foreground uppercase tracking-wider">Spell Save DC</p>
            <p className="font-display text-2xl text-primary">{saveDC}</p>
          </div>
          <ShieldCheck className="h-6 w-6 text-primary/50" />
        </div>
      </CardContent>
    </Card>
  );
}
