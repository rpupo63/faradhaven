import { useState, useEffect, useRef } from 'react';
import { Hash, Layers, TrendingUp, Activity, Loader2, Zap, ShieldAlert, ArrowUp, ArrowDown } from 'lucide-react';
import { RaIcon } from '@/components/ui/RaIcon';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { ApiClassResource } from '@/types/game/api';
import { cn } from '@/lib/utils';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { spendResource, gainResource } from '@/lib/api/resources';
import { castMutagen, rollTable, type RollTableResponse } from '@/lib/api/mechanics';
import { rollMadness } from '@/lib/api/madness';
import { updateSanguineNotoriety, forageComponents } from '@/lib/api/character';
import { PowderMageSpeedDialMenu } from './PowderMageSpeedDialMenu';
import { ConstructsSection } from './ConstructsSection';
import { HarvestBankSection } from './HarvestBankSection';
import { Input } from '@/components/ui/input';
import { toast } from 'sonner';
import type { NormalizedCharacterSheet } from '@/types/game';

interface ClassResourceDisplayProps {
  resources: ApiClassResource[];
  characterId: string;
  token: string;
  sheet?: NormalizedCharacterSheet;
  onGoToSpellForge?: () => void;
  onHarvest?: () => void;
  onScavenge?: () => void;
  scavengeLabel?: string;
  scavengeIcon?: 'recycle' | 'droplets';
}

/** Icon for a resource category */
function CategoryIcon({ category }: { category: string }) {
  switch (category) {
    case 'pool':
      return <RaIcon name="flask" className="text-sm text-primary" />;
    case 'die_size':
      return <RaIcon name="perspective-dice-six" className="text-sm text-primary" />;
    case 'limit':
      return <Hash className="h-4 w-4 text-primary" />;
    case 'slot_count':
      return <Layers className="h-4 w-4 text-primary" />;
    case 'modifier':
      return <TrendingUp className="h-4 w-4 text-primary" />;
    case 'state':
      return <Activity className="h-4 w-4 text-primary" />;
    default:
      return <Hash className="h-4 w-4 text-primary" />;
  }
}

/** Renders a single resource based on its category */
function ResourceItem({ resource, characterId, token }: { resource: ApiClassResource; characterId: string; token: string }) {
  const queryClient = useQueryClient();

  const spendMutation = useMutation({
    mutationFn: () => spendResource(characterId, resource.key, 1, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] }),
    onError: (error) => console.error('Failed to spend resource:', error),
  });

  const gainMutation = useMutation({
    mutationFn: () => gainResource(characterId, resource.key, 1, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] }),
    onError: (error) => console.error('Failed to gain resource:', error),
  });

  const isPending = spendMutation.isPending || gainMutation.isPending;

  const content = (() => {
    switch (resource.category) {
      case 'pool': {
        const current = resource.current_value ?? resource.value;
        const max = resource.max_value ?? resource.value;
        const pct = max > 0 ? (current / max) * 100 : 0;
        return (
          <div className="space-y-1.5">
            <div className="flex items-center justify-between">
              <span className="text-xl font-display text-primary">
                {current} / {max}
              </span>
            </div>
            <Progress value={pct} className="h-2" />
            {resource.is_trackable && (
              <div className="flex gap-1 justify-center">
                <Button
                  size="sm"
                  variant="outline"
                  className="h-6 w-8 p-0 text-xs"
                  onClick={() => gainMutation.mutate()}
                  disabled={isPending || current >= max}
                >
                  +1
                </Button>
                <Button
                  size="sm"
                  variant="outline"
                  className="h-6 w-8 p-0 text-xs"
                  onClick={() => spendMutation.mutate()}
                  disabled={isPending || current <= 0}
                >
                  -1
                </Button>
              </div>
            )}
          </div>
        );
      }

      case 'die_size':
        return (
          <span className="text-xl font-display text-primary">
            d{resource.value}
          </span>
        );

      case 'limit':
        return (
          <span className="text-xl font-display text-primary">
            {resource.value === 0 ? '∞' : resource.value}
          </span>
        );

      case 'slot_count': {
        const totalSlots = resource.value;
        const availableSlots = resource.is_trackable && resource.current_value != null
          ? resource.current_value
          : totalSlots;
        return (
          <div className="flex flex-wrap gap-1 items-center">
            {totalSlots > 0 ? (
              Array.from({ length: totalSlots }).map((_, i) => {
                const isAvailable = i < availableSlots;
                return resource.is_trackable ? (
                  <button
                    key={i}
                    onClick={() => isAvailable ? spendMutation.mutate() : gainMutation.mutate()}
                    disabled={isPending}
                    title={isAvailable ? 'Click to use slot' : 'Click to restore slot'}
                    className={cn(
                      'h-5 w-5 rounded-full flex items-center justify-center text-micro font-display transition-colors cursor-pointer',
                      isAvailable
                        ? 'bg-primary text-primary-foreground hover:bg-primary/80'
                        : 'border border-border text-muted-foreground hover:border-primary/50'
                    )}
                  >
                    {i + 1}
                  </button>
                ) : (
                  <Badge
                    key={i}
                    variant={isAvailable ? 'default' : 'outline'}
                    className="h-5 w-5 rounded-full p-0 flex items-center justify-center text-micro font-display"
                  >
                    {i + 1}
                  </Badge>
                );
              })
            ) : (
              <span className="text-xl font-display text-primary">0</span>
            )}
          </div>
        );
      }

      case 'modifier': {
        const displayValue = resource.value > 0 ? `+${resource.value}` : `${resource.value}`;
        const unit = resource.key === 'timer_duration' ? 's' : '';
        return (
          <span className="text-xl font-display text-primary">
            {displayValue}{unit}
          </span>
        );
      }

      case 'state': {
        const stateValue = resource.is_trackable && resource.current_value != null
          ? resource.current_value
          : resource.value;
        return (
          <div className="space-y-1.5">
            <span className="text-xl font-display text-primary">
              {stateValue}
            </span>
            {resource.is_trackable && (
              <div className="flex gap-1 justify-center">
                <Button
                  size="sm"
                  variant="outline"
                  className="h-6 w-8 p-0 text-xs"
                  onClick={() => gainMutation.mutate()}
                  disabled={isPending}
                >
                  +1
                </Button>
                <Button
                  size="sm"
                  variant="outline"
                  className="h-6 w-8 p-0 text-xs"
                  onClick={() => spendMutation.mutate()}
                  disabled={isPending}
                >
                  -1
                </Button>
              </div>
            )}
          </div>
        );
      }

      default:
        return (
          <span className="text-xl font-display text-primary">
            {resource.value}
          </span>
        );
    }
  })();

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <div className="flex flex-col items-center gap-1 min-w-[80px]">
            <p className="text-micro font-tome-marginalia text-muted-foreground uppercase text-center leading-tight">
              {resource.display_name}
            </p>
            {content}
          </div>
        </TooltipTrigger>
        {resource.description && (
          <TooltipContent side="top" className="max-w-xs">
            <p className="text-sm">{resource.description}</p>
          </TooltipContent>
        )}
      </Tooltip>
    </TooltipProvider>
  );
}

// ─── Piston Brawler ───────────────────────────────────────────────────────────

function PistonBrawlerActions({
  characterId,
  token,
  sheet,
}: {
  characterId: string;
  token: string;
  sheet: NormalizedCharacterSheet;
}) {
  const queryClient = useQueryClient();
  const resources = sheet.class_resources ?? [];
  const stabilityRes = resources.find((r) => r.key === 'max_stability');
  const currentStability = stabilityRes?.current_value ?? stabilityRes?.value ?? 0;

  const level = sheet.character.level;
  const archetype = sheet.character.archetypeName;
  const overdriveDice =
    level >= 11 ? (archetype === 'The Vulcan' && level >= 15 ? 4 : 3) : 2;

  const overdriveMutation = useMutation({
    mutationFn: () => spendResource(characterId, 'max_stability', 1, token),
    onSuccess: () => {
      toast.success(`Overdrive Strike! Roll ${overdriveDice}d6 force damage — add to your hit!`);
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
    },
    onError: (error: Error) => toast.error(error.message || 'Failed to trigger overdrive'),
  });

  return (
    <>
      <Separator />
      <Button
        className="w-full"
        variant="destructive"
        onClick={() => overdriveMutation.mutate()}
        disabled={overdriveMutation.isPending || currentStability <= 0}
      >
        {overdriveMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        <Zap className="mr-1 h-4 w-4" />
        Overdrive Strike ({overdriveDice}d6)
      </Button>
    </>
  );
}

// ─── Mutagen ─────────────────────────────────────────────────────────────────

function MutagenActions({
  characterId,
  token,
  sheet,
}: {
  characterId: string;
  token: string;
  sheet: NormalizedCharacterSheet;
}) {
  const queryClient = useQueryClient();
  const resources = sheet.class_resources ?? [];
  const madnessDCRes = resources.find((r) => r.key === 'madness_base_dc');
  const madnessDC = madnessDCRes?.current_value ?? 10;

  const [isFeral, setIsFeral] = useState(false);
  const [feralSecondsLeft, setFeralSecondsLeft] = useState(0);
  const feralTimerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const [lastFeralRoll, setLastFeralRoll] = useState<RollTableResponse | null>(null);

  useEffect(() => {
    if (isFeral && feralSecondsLeft > 0) {
      feralTimerRef.current = setInterval(() => {
        setFeralSecondsLeft((s) => {
          if (s <= 1) { setIsFeral(false); return 0; }
          return s - 1;
        });
      }, 1000);
    }
    return () => { if (feralTimerRef.current) clearInterval(feralTimerRef.current); };
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

  const feralMins = Math.floor(feralSecondsLeft / 60);
  const feralSecs = feralSecondsLeft % 60;

  const castMutation = useMutation({
    mutationFn: () => castMutagen(characterId, token),
    onSuccess: (data) => {
      toast.success(`Mutation cast! DC is now ${data.current_dc} (casts: ${data.madness_cast_count})`);
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
    },
    onError: (error: Error) => toast.error(error.message || 'Failed to cast mutation'),
  });

  const resetDCMutation = useMutation({
    mutationFn: async () => {
      const diff = 10 - madnessDC;
      if (diff > 0) return gainResource(characterId, 'madness_base_dc', diff, token);
      if (diff < 0) return spendResource(characterId, 'madness_base_dc', Math.abs(diff), token);
    },
    onSuccess: () => {
      toast.success('Madness DC reset to 10');
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
    },
    onError: (error: Error) => toast.error(error.message || 'Failed to reset DC'),
  });

  const rollFeralMutation = useMutation({
    mutationFn: () => rollTable(characterId, 'mutagen_feral', token),
    onSuccess: (data) => {
      setLastFeralRoll(data);
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
    },
    onError: (error: Error) => toast.error(error.message || 'Failed to roll feral table'),
  });

  return (
    <>
      <Separator />
      {/* Feral Mode */}
      <div className="flex items-center justify-between">
        <div>
          {isFeral ? (
            <div>
              <Badge variant="destructive" className="mb-1">FERAL ACTIVE</Badge>
              <p className="text-sm text-orange-400 font-mono">
                {feralMins}:{String(feralSecs).padStart(2, '0')} remaining
              </p>
            </div>
          ) : (
            <p className="text-sm text-muted-foreground">Feral Mode inactive</p>
          )}
        </div>
        <Button size="sm" variant={isFeral ? 'destructive' : 'outline'} onClick={toggleFeral}>
          <RaIcon name="fire" className="mr-1 text-xs" />
          {isFeral ? 'End Feral' : 'Enter Feral'}
        </Button>
      </div>
      {/* Actions */}
      <div className="flex gap-2">
        <Button
          className="flex-1"
          size="sm"
          onClick={() => castMutation.mutate()}
          disabled={castMutation.isPending}
        >
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
      {/* Roll Feral Table */}
      <Button
        className="w-full"
        variant="outline"
        size="sm"
        onClick={() => rollFeralMutation.mutate()}
        disabled={rollFeralMutation.isPending}
      >
        {rollFeralMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        <RaIcon name="lightning-bolt" className="mr-1 text-sm" />
        Roll Feral Table
      </Button>
      {lastFeralRoll && (
        <div className="p-3 rounded bg-secondary text-sm">
          <p className="font-bold text-primary">Roll: {lastFeralRoll.roll}</p>
          {lastFeralRoll.effect && (
            <>
              <p className="font-semibold text-orange-400 mt-1">{lastFeralRoll.effect.name}</p>
              <p className="text-muted-foreground text-xs mt-1">{lastFeralRoll.effect.description}</p>
            </>
          )}
        </div>
      )}
    </>
  );
}

// ─── Powder Mage ──────────────────────────────────────────────────────────────

function PowderMageActions({
  characterId,
  token,
  sheet,
  onGoToSpellForge,
}: {
  characterId: string;
  token: string;
  sheet: NormalizedCharacterSheet;
  onGoToSpellForge?: () => void;
}) {
  const resources = sheet.class_resources ?? [];
  const timerDuration = resources.find((r) => r.key === 'timer_duration')?.value ?? 2;

  const [casting, setCasting] = useState(false);
  const [secondsLeft, setSecondsLeft] = useState(0);
  const [done, setDone] = useState(false);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    if (casting && secondsLeft > 0) {
      timerRef.current = setInterval(() => {
        setSecondsLeft((s) => {
          if (s <= 1) {
            setCasting(false);
            setDone(true);
            toast.info('Casting window closed!');
            return 0;
          }
          return s - 1;
        });
      }, 1000);
    }
    return () => { if (timerRef.current) clearInterval(timerRef.current); };
  }, [casting, secondsLeft]);

  const startCast = () => { setDone(false); setCasting(true); setSecondsLeft(timerDuration); };
  const reset = () => {
    setCasting(false);
    setDone(false);
    setSecondsLeft(0);
    if (timerRef.current) clearInterval(timerRef.current);
  };

  return (
    <>
      <Separator />
      {/* Speed Dial */}
      <PowderMageSpeedDialMenu
        sheet={sheet}
        characterId={characterId}
        token={token}
        onGoToSpellForge={onGoToSpellForge}
      />
      <Separator />
      {/* Cast Timer */}
      <div className="flex items-center justify-between">
        <div className="text-center">
          <p className="text-2xl font-display text-orange-400">
            {casting ? secondsLeft : timerDuration}s
          </p>
          <p className="text-micro text-muted-foreground uppercase">
            {casting ? 'remaining' : 'window'}
          </p>
        </div>
        <div className="flex flex-col gap-1 items-end">
          {casting && <Badge variant="destructive">CASTING</Badge>}
          {done && !casting && <Badge variant="outline" className="text-muted-foreground">DONE</Badge>}
          {casting || done ? (
            <Button size="sm" variant="outline" onClick={reset}>Reset</Button>
          ) : (
            <Button size="sm" onClick={startCast}>Start Cast</Button>
          )}
        </div>
      </div>
    </>
  );
}

// ─── Ironwright ──────────────────────────────────────────────────────────────

function IronwrightSection({ sheet }: { sheet: NormalizedCharacterSheet }) {
  return <ConstructsSection sheet={sheet} />;
}

// ─── Lorewright ──────────────────────────────────────────────────────────────

function LorewrightSection({
  characterId,
  token,
  sheet,
}: {
  characterId: string;
  token: string;
  sheet: NormalizedCharacterSheet;
}) {
  const queryClient = useQueryClient();
  const [strainCR, setStrainCR] = useState('');
  const level = sheet.character.level;

  const scavengeMutation = useMutation({
    mutationFn: () => forageComponents(characterId, token),
    onSuccess: (data) => {
      const names = data.components.map((c: { name: string }) => c.name).join(', ');
      toast.success(`Scavenged ${data.yield} component(s): ${names}`);
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
    },
    onError: (error: Error) => toast.error(error.message || 'Failed to scavenge components'),
  });

  const madnessMutation = useMutation({
    mutationFn: () => rollMadness(characterId, token),
    onSuccess: (data) => {
      toast.info(`Madness roll: ${data.roll} on d${data.madness_die}. ${data.effect ? `Effect: ${data.effect}` : 'No effect.'}`);
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
    },
    onError: (error: Error) => toast.error(error.message || 'Failed to roll madness'),
  });

  return (
    <>
      <HarvestBankSection sheet={sheet} />
      <Separator />
      {/* Trauma */}
      <div className="flex items-center gap-2">
        <RaIcon name="skull" className="text-base text-red-500" />
        <span className="text-sm font-tome-subheading text-primary">Trauma:</span>
        <span className="text-lg font-display text-red-400">{sheet.trauma ?? 0}</span>
        <span className="text-xs text-muted-foreground">
          {sheet.trauma === 1 && '— Disadv. Charisma checks'}
          {sheet.trauma === 2 && '— Temporary flaw'}
          {(sheet.trauma ?? 0) >= 3 && '— Severe Disorientation'}
        </span>
      </div>
      {/* Psychic Strain DC */}
      <div className="space-y-1">
        <p className="text-xs text-muted-foreground">
          Psychic Strain: safe CR ≤ {(level / 3).toFixed(1)} — Wisdom save above that
        </p>
        <div className="flex gap-2 items-center">
          <Input
            type="number"
            placeholder="Creature CR"
            value={strainCR}
            onChange={(e) => setStrainCR(e.target.value)}
            className="h-8 text-sm w-28"
            step="0.25"
            min="0"
          />
          {(() => {
            const cr = parseFloat(strainCR);
            if (isNaN(cr) || strainCR === '') return null;
            const safeThreshold = level / 3;
            if (cr <= safeThreshold) return <p className="text-xs text-green-400">No save required</p>;
            return <p className="text-xs text-red-400">DC {10 + Math.ceil(cr - safeThreshold)} Wis save</p>;
          })()}
        </div>
      </div>
      {/* Actions */}
      <div className="flex gap-2">
        <Button
          className="flex-1"
          size="sm"
          variant="outline"
          onClick={() => madnessMutation.mutate()}
          disabled={madnessMutation.isPending}
        >
          {madnessMutation.isPending && <Loader2 className="mr-1 h-3 w-3 animate-spin" />}
          Roll Madness Die
        </Button>
        <Button
          className="flex-1"
          size="sm"
          variant="outline"
          onClick={() => scavengeMutation.mutate()}
          disabled={scavengeMutation.isPending}
        >
          {scavengeMutation.isPending && <Loader2 className="mr-1 h-3 w-3 animate-spin" />}
          Scavenge
        </Button>
      </div>
    </>
  );
}

// ─── Sanguinist ──────────────────────────────────────────────────────────────

function SanguinistSection({
  characterId,
  token,
  sheet,
}: {
  characterId: string;
  token: string;
  sheet: NormalizedCharacterSheet;
}) {
  const queryClient = useQueryClient();
  const sanguineMp = sheet.character.sanguine_mp ?? 0;
  const sanguineBr = sheet.character.sanguine_br ?? 0;
  const notoriety = sanguineMp - sanguineBr;
  const percentage = Math.min(100, Math.max(0, ((notoriety + 20) / 40) * 100));

  const ichorResource = sheet.class_resources?.find((r) => r.key === 'max_blood_ichor');
  const ichorBelowHalf =
    ichorResource &&
    ichorResource.max_value != null &&
    ichorResource.max_value > 0 &&
    (ichorResource.current_value ?? 0) < ichorResource.max_value * 0.5;

  const getLabel = () => {
    if (notoriety <= -15) return 'Medical Prodigy';
    if (notoriety <= -5) return 'Benevolent';
    if (notoriety < 5) return 'Neutral';
    if (notoriety < 15) return 'Starving Predator';
    return 'Blood Rage';
  };

  const getStatusColor = () => {
    if (notoriety <= -10) return 'text-blue-500';
    if (notoriety < 0) return 'text-blue-300';
    if (notoriety === 0) return 'text-muted-foreground';
    if (notoriety < 10) return 'text-red-300';
    return 'text-red-500';
  };

  const notorietyMutation = useMutation({
    mutationFn: ({ mpChange, brChange }: { mpChange: number; brChange: number }) =>
      updateSanguineNotoriety(characterId, { mpChange, brChange }, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] }),
    onError: (error: Error) => toast.error(error.message || 'Failed to update notoriety'),
  });

  const forageMutation = useMutation({
    mutationFn: () => forageComponents(characterId, token),
    onSuccess: (data) => {
      const names = data.components.map((c: { name: string }) => c.name).join(', ');
      toast.success(`Extracted ${data.yield} component(s): ${names}`);
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
    },
    onError: (error: Error) => toast.error(error.message || 'Failed to extract components'),
  });

  const isPending = notorietyMutation.isPending;

  return (
    <>
      <Separator />
      <div className="flex items-center justify-between">
        <span className="flex items-center gap-2 text-sm font-tome-subheading text-primary">
          <ShieldAlert className="h-4 w-4" />
          Sanguine Notoriety
        </span>
        <span className={cn('text-xs font-tome-marginalia uppercase tracking-widest', getStatusColor())}>
          {getLabel()}
        </span>
      </div>
      <div className="relative">
        <Progress value={percentage} className="h-2 bg-blue-900/20" />
        <div className="absolute top-0 left-1/2 -translate-x-1/2 h-4 w-0.5 bg-border z-10" />
      </div>
      <div className="flex items-center justify-between text-center">
        <div>
          <p className="text-lg font-display text-blue-400">{sanguineMp}</p>
          <p className="text-micro font-tome-marginalia text-muted-foreground uppercase">MP</p>
        </div>
        <div>
          <p className="text-xl font-display text-primary">{notoriety}</p>
          <p className="text-micro font-tome-marginalia text-muted-foreground uppercase">Net</p>
        </div>
        <div>
          <p className="text-lg font-display text-red-400">{sanguineBr}</p>
          <p className="text-micro font-tome-marginalia text-muted-foreground uppercase">BR</p>
        </div>
      </div>
      <div className="flex flex-wrap gap-1 justify-center">
        <Button size="sm" variant="outline" className="h-7 text-xs" onClick={() => notorietyMutation.mutate({ mpChange: 1, brChange: 0 })} disabled={isPending}>
          <ArrowUp className="h-3 w-3 mr-0.5" />+MP
        </Button>
        <Button size="sm" variant="outline" className="h-7 text-xs" onClick={() => notorietyMutation.mutate({ mpChange: -1, brChange: 0 })} disabled={isPending || sanguineMp <= 0}>
          <ArrowDown className="h-3 w-3 mr-0.5" />-MP
        </Button>
        <Button size="sm" variant="outline" className="h-7 text-xs" onClick={() => notorietyMutation.mutate({ mpChange: 0, brChange: 1 })} disabled={isPending}>
          <ArrowUp className="h-3 w-3 mr-0.5" />+BR
        </Button>
        <Button size="sm" variant="outline" className="h-7 text-xs" onClick={() => notorietyMutation.mutate({ mpChange: 0, brChange: -1 })} disabled={isPending || sanguineBr <= 0}>
          <ArrowDown className="h-3 w-3 mr-0.5" />-BR
        </Button>
      </div>
      {ichorBelowHalf && (
        <div className="rounded border border-yellow-500/50 bg-yellow-500/10 px-3 py-2 text-xs text-yellow-400">
          ⚠ Exhaustion risk: DC 13 Con save if you spend &gt;50% Ichor per turn
        </div>
      )}
      <Button
        className="w-full"
        size="sm"
        variant="outline"
        onClick={() => forageMutation.mutate()}
        disabled={forageMutation.isPending}
      >
        {forageMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        Extract Components
      </Button>
    </>
  );
}

// ─── Main export ─────────────────────────────────────────────────────────────

/**
 * ClassResourceDisplay renders all class resources for a character, plus any
 * class-specific interactive mechanics (Mutagen cast/feral, Piston overdrive,
 * Powder Mage speed dial + timer).
 */
export function ClassResourceDisplay({
  resources,
  characterId,
  token,
  sheet,
  onGoToSpellForge,
  onHarvest,
  onScavenge,
  scavengeLabel = 'Scavenge',
  scavengeIcon = 'recycle',
}: ClassResourceDisplayProps) {
  if (!resources || resources.length === 0) return null;

  const sorted = [...resources].sort((a, b) => a.display_order - b.display_order);
  const scavengeIconName = scavengeIcon === 'droplets' ? 'water-drop' : 'regeneration';
  const className = sheet?.class?.name;

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-tome-marginalia uppercase tracking-wider flex items-center gap-2">
          <CategoryIcon category={sorted[0]?.category ?? 'limit'} />
          Class Resources
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex flex-wrap gap-4 justify-center">
          {sorted.map((resource) => (
            <ResourceItem key={resource.key} resource={resource} characterId={characterId} token={token} />
          ))}
        </div>

        {(onHarvest || onScavenge) && (
          <div className="flex flex-wrap gap-2 justify-center border-t border-border pt-3">
            {onHarvest && (
              <Button variant="outline" size="sm" onClick={onHarvest} className="gap-1.5">
                <RaIcon name="flask" className="text-sm" />
                Harvest
              </Button>
            )}
            {onScavenge && (
              <Button variant="outline" size="sm" onClick={onScavenge} className="gap-1.5">
                <RaIcon name={scavengeIconName} className="text-sm" />
                {scavengeLabel}
              </Button>
            )}
          </div>
        )}

        {/* Class-specific interactive sections */}
        {sheet && className === 'The Piston Brawler' && (
          <PistonBrawlerActions characterId={characterId} token={token} sheet={sheet} />
        )}
        {sheet && className === 'The Mutagen' && (
          <MutagenActions characterId={characterId} token={token} sheet={sheet} />
        )}
        {sheet && className === 'The Powder Mage' && (
          <PowderMageActions
            characterId={characterId}
            token={token}
            sheet={sheet}
            onGoToSpellForge={onGoToSpellForge}
          />
        )}
        {sheet && className === 'The Ironwright' && (
          <IronwrightSection sheet={sheet} />
        )}
        {sheet && className === 'The Lorewright' && (
          <LorewrightSection characterId={characterId} token={token} sheet={sheet} />
        )}
        {sheet && className === 'The Sanguinist' && (
          <SanguinistSection characterId={characterId} token={token} sheet={sheet} />
        )}
      </CardContent>
    </Card>
  );
}
