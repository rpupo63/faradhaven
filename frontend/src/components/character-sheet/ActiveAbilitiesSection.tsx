import { useState } from 'react';
import { Zap, ChevronDown, ChevronUp, Sparkles } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { NormalizedCharacterSheet, DND5E_SKILLS } from '@/types/game';
import { ApiTrait, ApiLevelFeature, ApiClassResource } from '@/types/game/api';
import { cn } from '@/lib/utils';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { callTraitAbility, callFeatureAbility, restoreTraitAbility } from '@/lib/api/abilities';

interface ActiveAbilitiesSectionProps {
  sheet: NormalizedCharacterSheet;
  characterId: string;
  token: string;
}

function actionTypeBadgeVariant(actionType: string): 'default' | 'secondary' | 'outline' {
  switch (actionType) {
    case 'Action': return 'default';
    case 'Bonus Action': return 'secondary';
    case 'Reaction': return 'outline';
    default: return 'outline';
  }
}

const TRANSFORMATION_TRAITS = ['Shifting', 'Celestial Revelation', 'Large Form', 'Draconic Flight'];

function ResourceCostLabel({
  costs,
  classResources,
}: {
  costs: Array<{ key: string; amount: number }>;
  classResources?: ApiClassResource[];
}) {
  if (!costs.length) return null;
  const labels = costs.map((c) => {
    const res = classResources?.find((r) => r.key === c.key);
    const name = res?.display_name ?? c.key;
    return `${c.amount} ${name}`;
  });
  return (
    <span className="text-xs text-muted-foreground">Costs: {labels.join(', ')}</span>
  );
}

/** Can the character afford this feature's resource costs? */
function canAffordCosts(
  costs: Array<{ key: string; amount: number }>,
  classResources?: ApiClassResource[]
): boolean {
  if (!costs.length) return true;
  for (const cost of costs) {
    const res = classResources?.find((r) => r.key === cost.key);
    if (!res) return false;
    const current = res.current_value ?? res.value;
    if (current < cost.amount) return false;
  }
  return true;
}

interface TraitAbilityCardProps {
  trait: ApiTrait;
  currentUses?: number;
  maxUses?: number;
  characterId: string;
  token: string;
}

function TraitAbilityCard({ trait, currentUses, maxUses, characterId, token }: TraitAbilityCardProps) {
  const queryClient = useQueryClient();
  const [expanded, setExpanded] = useState(false);
  // local state for transformation visual feedback
  const [isActive, setIsActive] = useState(false);

  const hasUsesTracking = !!trait.uses_per_rest;
  const isTransformation = TRANSFORMATION_TRAITS.includes(trait.name);
  const isPassiveTrigger = trait.action_type === 'Passive';

  const mutation = useMutation({
    mutationFn: () => callTraitAbility(characterId, trait.id, token),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      if (isTransformation) setIsActive(true);
    },
    onError: (err) => console.error('Failed to use trait:', err),
  });

  const restoreMutation = useMutation({
    mutationFn: () => restoreTraitAbility(characterId, trait.id, token),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      if (isTransformation) setIsActive(false);
    },
    onError: (err) => console.error('Failed to restore trait:', err),
  });

  const disabled = mutation.isPending || restoreMutation.isPending || (hasUsesTracking && currentUses === 0);

  return (
    <div className={cn(
      "rounded-md border p-3 space-y-2 transition-all duration-300",
      isActive ? "border-primary/50 bg-primary/5 shadow-[0_0_15px_rgba(var(--primary),0.1)]" : "border-border bg-card/50"
    )}>
      <div className="flex items-start justify-between gap-2">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 flex-wrap">
            <span className="text-sm font-tome-subheading text-primary truncate">{trait.name}</span>
            <Badge variant={actionTypeBadgeVariant(trait.action_type)} className="text-[10px] shrink-0">
              {trait.action_type}
            </Badge>
            {isActive && <Sparkles className="h-3 w-3 text-primary animate-pulse" />}
          </div>
          
          {hasUsesTracking && currentUses !== undefined && maxUses !== undefined && maxUses > 1 && (
            <div className="flex gap-1 mt-1.5">
              {Array.from({ length: maxUses }).map((_, i) => (
                <div 
                  key={i} 
                  className={cn(
                    "h-1.5 w-1.5 rounded-full transition-colors",
                    i < currentUses ? "bg-primary" : "bg-muted border border-primary/20"
                  )} 
                />
              ))}
            </div>
          )}
          
          {hasUsesTracking && currentUses !== undefined && (maxUses === undefined || maxUses <= 1) && (
            <span className="text-[10px] text-muted-foreground uppercase font-tome-marginalia mt-1 block">
              {currentUses > 0 ? 'Available' : 'Spent'}
            </span>
          )}
        </div>

        <div className="flex items-center gap-1 shrink-0">
          {maxUses === 1 ? (
            <Button
              size="sm"
              variant={currentUses! > 0 ? "outline" : "ghost"}
              className={cn(
                "h-7 px-2 text-[10px] font-display",
                currentUses! > 0 ? "border-primary/30 text-primary hover:bg-primary/10" : "text-muted-foreground"
              )}
              onClick={() => currentUses! > 0 ? mutation.mutate() : restoreMutation.mutate()}
              disabled={mutation.isPending || restoreMutation.isPending}
            >
              {mutation.isPending || restoreMutation.isPending ? '...' : 
               currentUses! > 0 ? (isTransformation ? 'ACTIVATE' : (isPassiveTrigger ? 'MARK USED' : 'USE')) : 'RESTORE'}
            </Button>
          ) : (
            <div className="flex gap-1">
               {currentUses !== undefined && currentUses < (maxUses || 0) && (
                 <Button
                    size="sm"
                    variant="ghost"
                    className="h-7 px-2 text-[10px] font-display text-muted-foreground hover:text-primary"
                    onClick={() => restoreMutation.mutate()}
                    disabled={restoreMutation.isPending}
                 >
                   +
                 </Button>
               )}
               <Button
                size="sm"
                variant="outline"
                className="h-7 px-2 text-xs"
                onClick={() => mutation.mutate()}
                disabled={disabled}
              >
                {mutation.isPending ? '...' : (isPassiveTrigger ? 'Mark Used' : 'Use')}
              </Button>
            </div>
          )}
          <Button
            size="sm"
            variant="ghost"
            className="h-7 w-7 p-0"
            onClick={() => setExpanded(!expanded)}
          >
            {expanded ? <ChevronUp className="h-3 w-3" /> : <ChevronDown className="h-3 w-3" />}
          </Button>
        </div>
      </div>
      {expanded && (
        <div className="space-y-2 animate-in fade-in duration-200">
          <p className="text-xs text-muted-foreground leading-relaxed">
            {trait.description}
          </p>
          
          {trait.name === 'Severed from Dreams' && (
            <div className="pt-2 border-t border-border/50">
              <p className="text-[10px] text-muted-foreground uppercase font-tome-marginalia mb-1">Select Daily Proficiency:</p>
              <Select
                defaultValue={localStorage.getItem(`kalashtar_skill_${characterId}`) || ''}
                onValueChange={(val) => {
                  localStorage.setItem(`kalashtar_skill_${characterId}`, val);
                  queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
                }}
              >
                <SelectTrigger className="h-8 text-xs bg-background/50">
                  <SelectValue placeholder="Choose a skill..." />
                </SelectTrigger>
                <SelectContent>
                  {DND5E_SKILLS.map(skill => (
                    <SelectItem key={skill.id} value={skill.id} className="text-xs">
                      {skill.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

interface FeatureAbilityCardProps {
  feature: ApiLevelFeature;
  classResources?: ApiClassResource[];
  characterId: string;
  token: string;
}

function FeatureAbilityCard({ feature, classResources, characterId, token }: FeatureAbilityCardProps) {
  const queryClient = useQueryClient();
  const [expanded, setExpanded] = useState(false);

  const costs = feature.resource_costs ?? [];
  const affordable = canAffordCosts(costs, classResources);

  const mutation = useMutation({
    mutationFn: () => callFeatureAbility(characterId, feature.id, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] }),
    onError: (err) => console.error('Failed to use feature:', err),
  });

  const disabled = mutation.isPending || !affordable;

  return (
    <div className="rounded-md border border-border bg-card/50 p-3 space-y-2">
      <div className="flex items-start justify-between gap-2">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 flex-wrap">
            <span className="text-sm font-tome-subheading text-primary truncate">{feature.name}</span>
            {feature.action_type && (
              <Badge variant={actionTypeBadgeVariant(feature.action_type)} className="text-[10px] shrink-0">
                {feature.action_type}
              </Badge>
            )}
          </div>
          {costs.length > 0 && (
            <div className="mt-1">
              <ResourceCostLabel costs={costs} classResources={classResources} />
            </div>
          )}
        </div>
        <div className="flex items-center gap-1 shrink-0">
          <Button
            size="sm"
            variant="outline"
            className={cn('h-7 px-2 text-xs', !affordable && 'opacity-50')}
            onClick={() => mutation.mutate()}
            disabled={disabled}
          >
            {mutation.isPending ? '...' : 'Use'}
          </Button>
          <Button
            size="sm"
            variant="ghost"
            className="h-7 w-7 p-0"
            onClick={() => setExpanded(!expanded)}
          >
            {expanded ? <ChevronUp className="h-3 w-3" /> : <ChevronDown className="h-3 w-3" />}
          </Button>
        </div>
      </div>
      {expanded && (
        <p className="text-xs text-muted-foreground leading-relaxed">{feature.description}</p>
      )}
    </div>
  );
}

export function ActiveAbilitiesSection({ sheet, characterId, token }: ActiveAbilitiesSectionProps) {
  const traitUseStates = sheet.trait_use_states ?? {};
  const traitMaxUses = sheet.trait_max_uses ?? {};
  const classResources = sheet.class_resources;

  // Active race traits: include Passive traits if they have uses_per_rest (tracking)
  const activeTraits = (sheet.race_traits ?? []).filter(
    (t): t is ApiTrait => (!!t.action_type && t.action_type !== 'Passive') || !!t.uses_per_rest
  );

  // Active class features: have action_type set
  const activeFeatures = (sheet.class_level.level_features ?? []).filter(
    (f): f is ApiLevelFeature => !!f.action_type
  );

  if (activeTraits.length === 0 && activeFeatures.length === 0) return null;

  return (
    <Card className="arcane-border">
      <CardHeader className="pb-3">
        <CardTitle className="flex items-center gap-2 text-sm font-tome-subheading text-primary">
          <Zap className="h-4 w-4" />
          Active Abilities
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {activeTraits.length > 0 && (
          <div className="space-y-2">
            <p className="text-xs font-tome-marginalia text-muted-foreground uppercase tracking-wider">
              Race Abilities
            </p>
            {activeTraits.map((trait) => (
              <TraitAbilityCard
                key={trait.id}
                trait={trait}
                currentUses={traitUseStates[trait.id]}
                maxUses={traitMaxUses[trait.id]}
                characterId={characterId}
                token={token}
              />
            ))}
          </div>
        )}

        {activeFeatures.length > 0 && (
          <div className="space-y-2">
            <p className="text-xs font-tome-marginalia text-muted-foreground uppercase tracking-wider">
              Class Abilities
            </p>
            {activeFeatures.map((feature) => (
              <FeatureAbilityCard
                key={feature.id}
                feature={feature}
                classResources={classResources}
                characterId={characterId}
                token={token}
              />
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
