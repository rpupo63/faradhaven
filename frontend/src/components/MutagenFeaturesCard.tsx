import React, { useState, useEffect, useRef } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { Loader2, Flame, Zap } from 'lucide-react';
import { cn } from '@/lib/utils';
import { toast } from 'sonner';
import { castMutagen, rollTable, type RollTableResponse } from '@/lib/api/mechanics';
import { spendResource, gainResource } from '@/lib/api/resources';
import type { NormalizedCharacterSheet } from '@/types/game';

interface MutagenFeaturesCardProps {
  sheet: NormalizedCharacterSheet;
  token: string;
}

export const MutagenFeaturesCard: React.FC<MutagenFeaturesCardProps> = ({ sheet, token }) => {
  const queryClient = useQueryClient();
  const resources = sheet.class_resources ?? [];

  const madnessDCResource = resources.find(r => r.key === 'madness_base_dc');
  const resistanceResource = resources.find(r => r.key === 'madness_resistance_uses');
  const feralBonusResource = resources.find(r => r.key === 'feral_bonus');

  const madnessDC = madnessDCResource?.current_value ?? 10;
  const resistanceUses = resistanceResource?.current_value ?? 0;
  const maxResistanceUses = resistanceResource?.max_value ?? 0;
  const feralBonus = feralBonusResource?.current_value ?? 0;

  // Local feral mode state with 1-minute countdown
  const [isFeral, setIsFeral] = useState(false);
  const [feralSecondsLeft, setFeralSecondsLeft] = useState(0);
  const feralTimerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const [lastFeralRoll, setLastFeralRoll] = useState<RollTableResponse | null>(null);

  useEffect(() => {
    if (isFeral && feralSecondsLeft > 0) {
      feralTimerRef.current = setInterval(() => {
        setFeralSecondsLeft(s => {
          if (s <= 1) {
            setIsFeral(false);
            return 0;
          }
          return s - 1;
        });
      }, 1000);
    }
    return () => {
      if (feralTimerRef.current) clearInterval(feralTimerRef.current);
    };
  }, [isFeral, feralSecondsLeft]);

  const toggleFeral = () => {
    if (isFeral) {
      setIsFeral(false);
      setFeralSecondsLeft(0);
      if (feralTimerRef.current) clearInterval(feralTimerRef.current);
    } else {
      setIsFeral(true);
      setFeralSecondsLeft(60);
    }
  };

  const feralMinutes = Math.floor(feralSecondsLeft / 60);
  const feralSeconds = feralSecondsLeft % 60;

  const castMutation = useMutation({
    mutationFn: () => castMutagen(sheet.character.id, token),
    onSuccess: (data) => {
      toast.success(`Mutation cast! DC is now ${data.current_dc} (casts: ${data.madness_cast_count})`);
      queryClient.invalidateQueries({ queryKey: ['character-sheet', sheet.character.id] });
    },
    onError: (error: Error) => toast.error(error.message || 'Failed to cast mutation'),
  });

  const resetDCMutation = useMutation({
    mutationFn: async () => {
      const diff = 10 - madnessDC;
      if (diff > 0) return gainResource(sheet.character.id, 'madness_base_dc', diff, token);
      if (diff < 0) return spendResource(sheet.character.id, 'madness_base_dc', Math.abs(diff), token);
    },
    onSuccess: () => {
      toast.success('Madness DC reset to 10');
      queryClient.invalidateQueries({ queryKey: ['character-sheet', sheet.character.id] });
    },
    onError: (error: Error) => toast.error(error.message || 'Failed to reset DC'),
  });

  const rollFeralMutation = useMutation({
    mutationFn: () => rollTable(sheet.character.id, 'mutagen_feral', token),
    onSuccess: (data) => {
      setLastFeralRoll(data);
      queryClient.invalidateQueries({ queryKey: ['character-sheet', sheet.character.id] });
    },
    onError: (error: Error) => toast.error(error.message || 'Failed to roll feral table'),
  });

  const spendResistanceMutation = useMutation({
    mutationFn: () => spendResource(sheet.character.id, 'madness_resistance_uses', 1, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character-sheet', sheet.character.id] }),
    onError: (error: Error) => toast.error(error.message || 'Failed to spend resistance'),
  });

  const gainResistanceMutation = useMutation({
    mutationFn: () => gainResource(sheet.character.id, 'madness_resistance_uses', 1, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character-sheet', sheet.character.id] }),
    onError: (error: Error) => toast.error(error.message || 'Failed to restore resistance'),
  });

  return (
    <div className="space-y-4">
      {/* Madness DC Tracker */}
      <div>
        <h3 className="text-md font-semibold mb-2">Madness DC</h3>
        <div className="flex items-center justify-between">
          <div className="text-center">
            <p className="text-3xl font-display text-red-400">{madnessDC}</p>
            <p className="text-[10px] text-muted-foreground uppercase">Current DC</p>
          </div>
          <div className="flex flex-col gap-2">
            <Button size="sm" onClick={() => castMutation.mutate()} disabled={castMutation.isPending}>
              {castMutation.isPending && <Loader2 className="mr-1 h-3 w-3 animate-spin" />}
              Cast Ability
            </Button>
            <Button
              size="sm"
              variant="outline"
              onClick={() => resetDCMutation.mutate()}
              disabled={resetDCMutation.isPending || madnessDC === 10}
            >
              {resetDCMutation.isPending && <Loader2 className="mr-1 h-3 w-3 animate-spin" />}
              Reset DC
            </Button>
          </div>
        </div>
      </div>

      <Separator />

      {/* Feral Mode Toggle */}
      <div>
        <h3 className="text-md font-semibold mb-2">Feral Mode</h3>
        <div className="flex items-center justify-between">
          <div>
            {isFeral ? (
              <div>
                <Badge variant="destructive" className="mb-1">FERAL ACTIVE</Badge>
                <p className="text-sm text-orange-400 font-mono">
                  {feralMinutes}:{String(feralSeconds).padStart(2, '0')} remaining
                </p>
                <p className="text-xs text-muted-foreground">+{feralBonus} bonus damage</p>
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">Inactive</p>
            )}
          </div>
          <Button size="sm" variant={isFeral ? 'destructive' : 'outline'} onClick={toggleFeral}>
            <Flame className="mr-1 h-3 w-3" />
            {isFeral ? 'End Feral' : 'Enter Feral'}
          </Button>
        </div>
      </div>

      <Separator />

      {/* Roll Feral Table */}
      <div>
        <h3 className="text-md font-semibold mb-2">Feral Table Roll</h3>
        <Button className="w-full" onClick={() => rollFeralMutation.mutate()} disabled={rollFeralMutation.isPending}>
          {rollFeralMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          <Zap className="mr-1 h-4 w-4" />
          Roll Feral Table
        </Button>
        {lastFeralRoll && (
          <div className="mt-2 p-3 rounded bg-secondary text-sm">
            <p className="font-bold text-primary">Roll: {lastFeralRoll.roll}</p>
            {lastFeralRoll.effect && (
              <>
                <p className="font-semibold text-orange-400 mt-1">{lastFeralRoll.effect.name}</p>
                <p className="text-muted-foreground text-xs mt-1">{lastFeralRoll.effect.description}</p>
              </>
            )}
          </div>
        )}
      </div>

      {/* Madness Resistance */}
      {resistanceResource && (
        <>
          <Separator />
          <div>
            <h3 className="text-md font-semibold mb-1">Madness Resistance</h3>
            <p className="text-xs text-muted-foreground mb-3">
              Auto-succeed a failed Madness save (resets on long rest)
            </p>
            <div className="flex items-center justify-between">
              <div className="flex gap-1.5">
                {Array.from({ length: maxResistanceUses }).map((_, i) => (
                  <div
                    key={i}
                    className={cn(
                      'w-5 h-5 rounded-full border-2',
                      i < resistanceUses
                        ? 'bg-primary border-primary'
                        : 'bg-transparent border-muted-foreground/40'
                    )}
                  />
                ))}
              </div>
              <div className="flex gap-1">
                <Button
                  size="xs"
                  variant="outline"
                  className="h-7"
                  onClick={() => spendResistanceMutation.mutate()}
                  disabled={spendResistanceMutation.isPending || resistanceUses <= 0}
                >
                  Use
                </Button>
                <Button
                  size="xs"
                  variant="ghost"
                  className="h-7"
                  onClick={() => gainResistanceMutation.mutate()}
                  disabled={gainResistanceMutation.isPending || resistanceUses >= maxResistanceUses}
                >
                  +1
                </Button>
              </div>
            </div>
            <p className="text-xs text-muted-foreground mt-1">{resistanceUses}/{maxResistanceUses} remaining</p>
          </div>
        </>
      )}
    </div>
  );
};
