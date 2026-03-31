import { Trash2, AlertTriangle, Skull } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ActiveEffect, removeEffect } from '@/lib/api';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { useAuth } from '@/context/AuthContext';
import { cn } from '@/lib/utils';

interface ActiveEffectsSectionProps {
  characterId: string;
  effects: ActiveEffect[];
}

export function ActiveEffectsSection({ characterId, effects }: ActiveEffectsSectionProps) {
  const { token } = useAuth();
  const queryClient = useQueryClient();

  const removeMutation = useMutation({
    mutationFn: (effectId: string) => removeEffect(characterId, effectId, token || undefined),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['active-effects', characterId] });
      toast.success('Effect removed');
    },
    onError: (err: Error) => {
      toast.error(err.message || 'Failed to remove effect');
    }
  });

  if (!effects || effects.length === 0) return null;

  return (
    <Card className="arcane-border bg-card border-l-4 border-l-orange-500/50">
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center gap-2 text-base font-tome-subheading text-orange-600 dark:text-orange-400">
          <AlertTriangle className="h-4 w-4" />
          Active Effects & Conditions
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-2">
        {effects.map((effect) => (
          <div 
            key={effect.id} 
            className="flex flex-col gap-1 p-2 rounded border bg-muted/20 border-border/50 animate-in fade-in slide-in-from-left-2 duration-300"
          >
            <div className="flex justify-between items-start">
              <div className="flex items-center gap-2">
                <span className="font-bold text-primary text-sm flex items-center gap-1.5">
                  {effect.category.includes('Feral') || effect.category.includes('Madness') ? (
                    <Skull className="h-3 w-3 text-red-500" />
                  ) : null}
                  {effect.name}
                </span>
                <Badge variant="outline" size="sm">
                  {effect.source}
                </Badge>
              </div>
              <Button
                variant="ghost"
                size="icon"
                className="h-6 w-6 text-muted-foreground hover:text-red-500 hover:bg-red-500/10"
                onClick={() => removeMutation.mutate(effect.id)}
                disabled={removeMutation.isPending}
              >
                <Trash2 className="h-3 w-3" />
              </Button>
            </div>
            
            <p className="text-xs font-tome-marginalia text-foreground/90">
              {effect.description}
            </p>
            
            {effect.mechanics && (
              <div className="text-[10px] font-tome-marginalia text-muted-foreground italic border-t border-border/30 pt-1 mt-1">
                {effect.mechanics}
              </div>
            )}
          </div>
        ))}
      </CardContent>
    </Card>
  );
}
