import React from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { Progress } from '@/components/ui/progress';
import { Loader2, Zap } from 'lucide-react';
import { cn } from '@/lib/utils';
import { toast } from 'sonner';
import { spendResource, gainResource } from '@/lib/api/resources';
import type { NormalizedCharacterSheet } from '@/types/game';

interface PistonBrawlerFeaturesCardProps {
  sheet: NormalizedCharacterSheet;
  token: string;
}

export const PistonBrawlerFeaturesCard: React.FC<PistonBrawlerFeaturesCardProps> = ({ sheet, token }) => {
  const queryClient = useQueryClient();
  const resources = sheet.class_resources ?? [];

  const stabilityResource = resources.find(r => r.key === 'max_stability');
  const currentStability = stabilityResource?.current_value ?? stabilityResource?.value ?? 0;
  const maxStability = stabilityResource?.max_value ?? stabilityResource?.value ?? 0;

  const speedDialSlots = resources.find(r => r.key === 'speed_dial_slots')?.value ?? 0;

  const level = sheet.character.level;
  const archetype = sheet.character.archetypeName;
  const overdriveDice = level >= 11
    ? (archetype === 'The Vulcan' && level >= 15 ? 4 : 3)
    : 2;

  const spendStabilityMutation = useMutation({
    mutationFn: () => spendResource(sheet.character.id, 'max_stability', 1, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character-sheet', sheet.character.id] }),
    onError: (error: Error) => toast.error(error.message || 'Failed to spend stability'),
  });

  const gainStabilityMutation = useMutation({
    mutationFn: () => gainResource(sheet.character.id, 'max_stability', 1, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character-sheet', sheet.character.id] }),
    onError: (error: Error) => toast.error(error.message || 'Failed to gain stability'),
  });

  const overdriveStrikeMutation = useMutation({
    mutationFn: () => spendResource(sheet.character.id, 'max_stability', 1, token),
    onSuccess: () => {
      toast.success(`Overdrive Strike! Roll ${overdriveDice}d6 force damage — add to your hit!`);
      queryClient.invalidateQueries({ queryKey: ['character-sheet', sheet.character.id] });
    },
    onError: (error: Error) => toast.error(error.message || 'Failed to trigger overdrive'),
  });

  const progressPct = maxStability > 0 ? (currentStability / maxStability) * 100 : 0;

  return (
    <div className="space-y-4">
      {/* Stability Charges */}
      <div>
        <h3 className="text-md font-semibold mb-2">Stability Charges</h3>
        <div className="flex items-center justify-between mb-2">
          <p className="text-3xl font-display text-primary">{currentStability}/{maxStability}</p>
          <div className="flex gap-1">
            <Button
              size="xs"
              variant="ghost"
              className="h-7"
              onClick={() => gainStabilityMutation.mutate()}
              disabled={gainStabilityMutation.isPending || currentStability >= maxStability}
            >
              {gainStabilityMutation.isPending && <Loader2 className="mr-1 h-3 w-3 animate-spin" />}
              +1
            </Button>
            <Button
              size="xs"
              variant="outline"
              className="h-7"
              onClick={() => spendStabilityMutation.mutate()}
              disabled={spendStabilityMutation.isPending || currentStability <= 0}
            >
              {spendStabilityMutation.isPending && <Loader2 className="mr-1 h-3 w-3 animate-spin" />}
              Use (−1)
            </Button>
          </div>
        </div>
        <Progress value={progressPct} className="h-2 mb-3" />
        <Button
          className="w-full"
          variant="destructive"
          onClick={() => overdriveStrikeMutation.mutate()}
          disabled={overdriveStrikeMutation.isPending || currentStability <= 0}
        >
          {overdriveStrikeMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          <Zap className="mr-1 h-4 w-4" />
          Overdrive Strike ({overdriveDice}d6)
        </Button>
      </div>

      <Separator />

      {/* Blueprint Slots */}
      <div>
        <h3 className="text-md font-semibold mb-2">Blueprint Slots</h3>
        {speedDialSlots > 0 ? (
          <>
            <div className="flex gap-1.5 mb-2">
              {Array.from({ length: speedDialSlots }).map((_, i) => (
                <div
                  key={i}
                  className={cn('w-5 h-5 rounded-full border-2 bg-primary border-primary')}
                />
              ))}
            </div>
            <p className="text-xs text-muted-foreground">
              During a short or long rest, pre-calculate a component combination for instant casting in combat.
            </p>
          </>
        ) : (
          <p className="text-xs text-muted-foreground italic">Blueprint Slots unlock at level 13</p>
        )}
      </div>
    </div>
  );
};
