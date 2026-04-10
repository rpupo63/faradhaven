import { Hash, Layers, TrendingUp, Activity } from 'lucide-react';
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
import { ClassResourceExtraSections } from './classResourceSectionRegistry';
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

        {sheet && (
          <ClassResourceExtraSections
            characterId={characterId}
            token={token}
            sheet={sheet}
            onGoToSpellForge={onGoToSpellForge}
          />
        )}
      </CardContent>
    </Card>
  );
}
