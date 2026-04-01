import { Fragment, useState } from 'react';
import { Wrench, Plus, Minus, Trash2, Shield, Heart, Loader2, Dices, Recycle, Bug, ShieldAlert, Zap } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { NormalizedCharacterSheet } from '@/types/game';
import type { ApiMinion, ApiConstructTemplate, ApiCharacterComponent } from '@/types/game';
import { getMinions, getMinionTemplates, createConstruct, updateMinionHP, deleteMinion, forageComponents } from '@/lib/api';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/context/AuthContext';
import { toast } from 'sonner';

interface ConstructsSectionProps {
  sheet: NormalizedCharacterSheet;
}

export function ConstructsSection({ sheet }: ConstructsSectionProps) {
  const { token } = useAuth();
  const queryClient = useQueryClient();
  const [summonDialogOpen, setSummonDialogOpen] = useState(false);
  const [scavengeResult, setScavengeResult] = useState<number | null>(null);
  const [efficiencyDialog, setEfficiencyDialog] = useState<{ minionId: string; name: string } | null>(null);
  const [droneComponentPicker, setDroneComponentPicker] = useState(false);
  const [fortressDialog, setFortressDialog] = useState(false);
  const [kineticDialog, setKineticDialog] = useState(false);

  const characterId = sheet.character.id;
  // Read from generic class_resources instead of deprecated top-level fields
  const concurrencyLimit = sheet.class_resources?.find(r => r.key === 'concurrency_limit')?.value ?? 0;
  const yieldDie = sheet.class_resources?.find(r => r.key === 'yield_die')?.value ?? 4;
  const characterLevel = sheet.character.level;
  const archetypeName = sheet.character.archetypeName;
  const isSwarmEngineer = archetypeName === 'The Swarm Engineer';
  const isAutoForger = archetypeName === 'The Automaton Forger';
  const intMod = Math.max(1, Math.floor((sheet.character.intelligence - 10) / 2));
  const components = sheet.components || [];

  // Fetch active constructs
  const { data: constructMinions = [] } = useQuery({
    queryKey: ['minions', characterId, 'construct'],
    queryFn: () => getMinions(characterId, 'construct', token || undefined),
    enabled: !!characterId && !!token,
  });

  // Fetch drones
  const { data: droneMinions = [] } = useQuery({
    queryKey: ['minions', characterId, 'drone'],
    queryFn: () => getMinions(characterId, 'drone', token || undefined),
    enabled: !!characterId && !!token && isSwarmEngineer,
  });

  // Fetch templates
  const { data: templates } = useQuery({
    queryKey: ['construct-templates', characterId],
    queryFn: () => getMinionTemplates(characterId, token || undefined),
    enabled: !!characterId && !!token,
  });

  const activeConstructs = constructMinions.filter((m: ApiMinion) => m.is_active);
  const activeDrones = droneMinions.filter((m: ApiMinion) => m.is_active);

  // Check if character has a component by name
  const hasComponent = (name: string): boolean => {
    return components.some((c: ApiCharacterComponent) => c.component?.name === name && c.count >= 1);
  };

  // Check which required components are available
  const getComponentAvailability = (requiredComponents: string[]): { name: string; has: boolean }[] => {
    return requiredComponents.map((name) => ({
      name,
      has: hasComponent(name),
    }));
  };

  const allComponentsAvailable = (requiredComponents: string[]): boolean => {
    return requiredComponents.every((name) => hasComponent(name));
  };

  // Mutations
  const createMutation = useMutation({
    mutationFn: (templateKey: string) =>
      createConstruct(
        characterId,
        {
          minion_type: 'construct',
          template_key: templateKey,
        },
        token || undefined
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['minions', characterId] });
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      setSummonDialogOpen(false);
      toast.success('Construct assembled!');
    },
    onError: (err: Error) => toast.error(err.message || 'Failed to create construct'),
  });

  const createDroneMutation = useMutation({
    mutationFn: (componentId: string) =>
      createConstruct(
        characterId,
        {
          minion_type: 'drone',
          component_id: componentId,
        },
        token || undefined
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['minions', characterId] });
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      setDroneComponentPicker(false);
      toast.success('Drone deployed!');
    },
    onError: (err: Error) => toast.error(err.message || 'Failed to create drone'),
  });

  const hpMutation = useMutation({
    mutationFn: ({ minionId, delta }: { minionId: string; delta: number }) =>
      updateMinionHP(characterId, minionId, delta, token || undefined),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['minions', characterId] });
      // Check for construct reaching 0 HP → Efficiency prompt
      if (variables.delta < 0) {
        const minion = [...constructMinions, ...droneMinions].find((m) => m.id === variables.minionId);
        if (minion && (minion.current_hp + variables.delta) <= 0) {
          setEfficiencyDialog({ minionId: variables.minionId, name: minion.name });
        }
        // Fortress Protocol prompt (Automaton Forger L10+)
        if (isAutoForger && characterLevel >= 10 && minion?.minion_type === 'construct') {
          setFortressDialog(true);
        }
        // Kinetic Barrier prompt (L11+)
        if (characterLevel >= 11 && minion?.minion_type === 'construct') {
          setKineticDialog(true);
        }
      }
    },
    onError: (err: Error) => toast.error(err.message || 'Failed to update HP'),
  });

  const destroyMutation = useMutation({
    mutationFn: ({ minionId, reclaim }: { minionId: string; reclaim: boolean }) =>
      deleteMinion(characterId, minionId, reclaim, token || undefined),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['minions', characterId] });
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      if (data.components_reclaimed) {
        toast.success(`Construct destroyed. Reclaimed ${data.components_reclaimed} components.`);
      } else {
        toast.success('Construct destroyed.');
      }
    },
    onError: (err: Error) => toast.error(err.message || 'Failed to destroy construct'),
  });

  const scavengeMutation = useMutation({
    mutationFn: () => forageComponents(characterId, token || ''),
    onSuccess: (data) => {
      setScavengeResult(data.yield);
      const names = data.components.map((c) => c.name).join(', ');
      toast.success(`Scavenged ${data.yield} component(s): ${names}`);
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
    },
    onError: (err: Error) => toast.error(err.message || 'Failed to scavenge components'),
  });

  // Render a single minion card (used for both constructs and drones)
  const renderMinionCard = (minion: ApiMinion) => {
    const hpPct = minion.max_hp > 0 ? (minion.current_hp / minion.max_hp) * 100 : 0;
    return (
      <div className="border border-border rounded-md p-2 space-y-1">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="text-sm font-display text-primary">{minion.name}</span>
            {minion.template_key && (
              <Badge variant="secondary" size="sm" className="capitalize">
                {minion.template_key}
              </Badge>
            )}
          </div>
          <div className="flex items-center gap-1 text-xs text-muted-foreground">
            <Shield className="h-3 w-3" />
            {minion.armor_class}
          </div>
        </div>

        {/* HP Bar */}
        <div className="flex items-center gap-2">
          <Heart className="h-3 w-3 text-red-400" />
          <Progress value={hpPct} className="h-2 flex-1" />
          <span className="text-xs font-mono text-muted-foreground w-12 text-right">
            {minion.current_hp}/{minion.max_hp}
          </span>
        </div>

        {/* Controls */}
        <div className="flex items-center justify-between pt-1">
          <div className="flex items-center gap-1">
            <Button
              size="icon"
              variant="ghost"
              className="h-6 w-6"
              onClick={() => hpMutation.mutate({ minionId: minion.id, delta: -1 })}
              disabled={hpMutation.isPending}
            >
              <Minus className="h-3 w-3" />
            </Button>
            <Button
              size="icon"
              variant="ghost"
              className="h-6 w-6"
              onClick={() => hpMutation.mutate({ minionId: minion.id, delta: 1 })}
              disabled={hpMutation.isPending}
            >
              <Plus className="h-3 w-3" />
            </Button>
          </div>
          <div className="flex items-center gap-1">
            <Button
              size="sm"
              variant="ghost"
              className="h-6 px-2 text-xs text-green-400 hover:text-green-300"
              onClick={() => destroyMutation.mutate({ minionId: minion.id, reclaim: true })}
              disabled={destroyMutation.isPending}
              title="Destroy and reclaim components"
            >
              <Recycle className="h-3 w-3 mr-1" />
              Reclaim
            </Button>
            <Button
              size="icon"
              variant="ghost"
              className="h-6 w-6 text-red-400 hover:text-red-300"
              onClick={() => destroyMutation.mutate({ minionId: minion.id, reclaim: false })}
              disabled={destroyMutation.isPending}
              title="Destroy construct"
            >
              <Trash2 className="h-3 w-3" />
            </Button>
          </div>
        </div>
      </div>
    );
  };

  return (
    <>
      <Card className="arcane-border bg-card">
        <CardHeader className="pb-2">
          <CardTitle className="flex items-center justify-between text-base font-tome-subheading text-primary">
            <div className="flex items-center gap-2">
              <Wrench className="h-4 w-4" />
              Constructs
            </div>
            <Badge variant="outline" className="text-xs">
              {activeConstructs.length} / {concurrencyLimit}
            </Badge>
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {/* Active Constructs List */}
          {activeConstructs.length === 0 ? (
            <p className="text-xs text-muted-foreground italic font-tome-marginalia text-center py-2">
              No active constructs
            </p>
          ) : (
            <div className="space-y-2">
              {activeConstructs.map((minion) => (
                <Fragment key={minion.id}>{renderMinionCard(minion)}</Fragment>
              ))}
            </div>
          )}

          {/* Action Buttons */}
          <div className="flex gap-2 pt-1">
            <Button
              size="sm"
              variant="outline"
              className="flex-1 gap-1 border-primary/50"
              onClick={() => setSummonDialogOpen(true)}
              disabled={activeConstructs.length >= concurrencyLimit}
            >
              <Plus className="h-3 w-3" />
              Summon Construct
            </Button>
            <Button
              size="sm"
              variant="outline"
              className="gap-1 border-amber-500/50 text-amber-400 hover:bg-amber-900/20"
              onClick={() => scavengeMutation.mutate()}
              disabled={scavengeMutation.isPending}
            >
              {scavengeMutation.isPending ? (
                <Loader2 className="h-3 w-3 animate-spin" />
              ) : (
                <Dices className="h-3 w-3" />
              )}
              Scavenge (d{yieldDie})
            </Button>
          </div>

          {/* Scavenge Result */}
          {scavengeResult !== null && (
            <div className="text-center p-2 border border-amber-500/30 rounded bg-amber-900/10">
              <p className="text-xs text-muted-foreground font-tome-marginalia">Last Scavenge Yield</p>
              <p className="text-lg font-display text-amber-400">{scavengeResult} components</p>
            </div>
          )}

          {/* Drone Section (Swarm Engineer only) */}
          {isSwarmEngineer && (
            <div className="border-t border-border pt-3 mt-3">
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-2">
                  <Bug className="h-4 w-4 text-primary" />
                  <span className="text-sm font-tome-subheading text-primary">Drones</span>
                </div>
                <Badge variant="outline" className="text-xs">
                  {activeDrones.length} / {intMod}
                </Badge>
              </div>
              {activeDrones.length === 0 ? (
                <p className="text-xs text-muted-foreground italic text-center py-1">No active drones</p>
              ) : (
                <div className="space-y-2">
                  {activeDrones.map((minion) => (
                    <Fragment key={minion.id}>{renderMinionCard(minion)}</Fragment>
                  ))}
                </div>
              )}
              <Button
                size="sm"
                variant="outline"
                className="w-full mt-2 gap-1 border-primary/50"
                onClick={() => setDroneComponentPicker(true)}
                disabled={activeDrones.length >= intMod}
              >
                <Plus className="h-3 w-3" />
                Create Drone (1 component)
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Summon Construct Dialog */}
      <Dialog open={summonDialogOpen} onOpenChange={setSummonDialogOpen}>
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle className="text-center text-primary font-display">Summon Construct</DialogTitle>
            <DialogDescription className="text-center text-xs">
              Choose a construct template. Required components are deducted from your inventory.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-2">
            {templates &&
              Object.values(templates)
                .filter((t: ApiConstructTemplate) => t.key !== 'drone')
                .map((template: ApiConstructTemplate, templateIndex: number) => {
                  const required = template.required_components || [];
                  const availability = getComponentAvailability(required);
                  const canBuild = allComponentsAvailable(required);
                  const durationMinutes = required.length * 10;
                  return (
                    <button
                      key={template.key || `template-${templateIndex}`}
                      className="w-full text-left border border-border rounded-md p-3 hover:bg-accent/50 transition-colors disabled:opacity-50"
                      onClick={() => createMutation.mutate(template.key)}
                      disabled={createMutation.isPending || !canBuild}
                    >
                      <div className="flex items-center justify-between mb-1">
                        <span className="font-display text-sm text-primary">{template.name}</span>
                        <Badge variant="secondary" className="text-xs">{required.length} components</Badge>
                      </div>
                      <p className="text-xs text-muted-foreground mb-1">{template.description}</p>
                      <div className="flex gap-3 text-micro text-muted-foreground">
                        <span>AC {template.base_ac} | HP {template.base_hp} + 2×Lvl</span>
                        <span>Speed {template.base_speed}ft</span>
                        <span>{template.size}</span>
                      </div>
                      <div className="flex gap-3 text-micro text-muted-foreground mt-0.5">
                        <span>Duration: {durationMinutes} min</span>
                      </div>
                      {/* Required Components */}
                      {required.length > 0 && (
                        <div className="flex flex-wrap gap-1 mt-1">
                          {availability.map(({ name, has }, availIdx) => (
                            <Badge
                              key={`${name}-${availIdx}`}
                              variant={has ? 'secondary' : 'destructive'}
                              className="text-micro"
                            >
                              {name}
                            </Badge>
                          ))}
                        </div>
                      )}
                      {template.actions?.map((action, actionIdx) => (
                        <div key={`${action.name}-${actionIdx}`} className="text-micro text-muted-foreground mt-0.5">
                          {action.name}: {action.damage_dice || ''} {action.damage_type || ''} {action.range ? `(${action.range})` : ''}
                        </div>
                      ))}
                      {createMutation.isPending && (
                        <Loader2 className="h-3 w-3 animate-spin mt-1" />
                      )}
                    </button>
                  );
                })}
          </div>
        </DialogContent>
      </Dialog>

      {/* Drone Component Picker Dialog */}
      <Dialog open={droneComponentPicker} onOpenChange={setDroneComponentPicker}>
        <DialogContent className="max-w-xs">
          <DialogHeader>
            <DialogTitle className="text-center text-primary font-display">Create Drone</DialogTitle>
            <DialogDescription className="text-center text-xs">
              Choose which component to spend (any 1 component).
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-1 max-h-64 overflow-y-auto">
            {components
              .filter((c: ApiCharacterComponent) => c.count >= 1)
              .map((c: ApiCharacterComponent, compIdx: number) => (
                <button
                  key={c.component_id || `component-${compIdx}`}
                  className="w-full text-left border border-border rounded-md p-2 hover:bg-accent/50 transition-colors flex items-center justify-between"
                  onClick={() => createDroneMutation.mutate(c.component_id)}
                  disabled={createDroneMutation.isPending}
                >
                  <span className="text-sm">{c.component?.name || c.component_id}</span>
                  <Badge variant="outline" className="text-xs">×{c.count}</Badge>
                </button>
              ))}
          </div>
        </DialogContent>
      </Dialog>

      {/* Efficiency Prompt Dialog */}
      <Dialog open={!!efficiencyDialog} onOpenChange={() => setEfficiencyDialog(null)}>
        <DialogContent className="max-w-xs">
          <DialogHeader>
            <DialogTitle className="text-center text-primary font-display">Construct Destroyed</DialogTitle>
            <DialogDescription className="text-center text-xs">
              {efficiencyDialog?.name} has been reduced to 0 HP. Use <strong>Efficiency</strong> to reclaim 50% of components?
            </DialogDescription>
          </DialogHeader>
          <div className="flex gap-2">
            <Button
              className="flex-1"
              onClick={() => {
                if (efficiencyDialog) {
                  destroyMutation.mutate({ minionId: efficiencyDialog.minionId, reclaim: true });
                }
                setEfficiencyDialog(null);
              }}
            >
              <Recycle className="h-4 w-4 mr-1" />
              Reclaim
            </Button>
            <Button
              variant="outline"
              className="flex-1"
              onClick={() => {
                if (efficiencyDialog) {
                  destroyMutation.mutate({ minionId: efficiencyDialog.minionId, reclaim: false });
                }
                setEfficiencyDialog(null);
              }}
            >
              <Trash2 className="h-4 w-4 mr-1" />
              No
            </Button>
          </div>
        </DialogContent>
      </Dialog>

      {/* Fortress Protocol Prompt (Automaton Forger L10+) */}
      <Dialog open={fortressDialog} onOpenChange={setFortressDialog}>
        <DialogContent className="max-w-xs">
          <DialogHeader>
            <DialogTitle className="text-center text-primary font-display flex items-center justify-center gap-2">
              <ShieldAlert className="h-4 w-4" />
              Fortress Protocol
            </DialogTitle>
            <DialogDescription className="text-center text-xs">
              Spend 2 components to give your construct resistance to all damage from this source. (DM adjudicates)
            </DialogDescription>
          </DialogHeader>
          <div className="flex gap-2">
            <Button
              className="flex-1"
              onClick={() => {
                toast.info('Fortress Protocol activated — DM: apply resistance to all damage from this source.');
                setFortressDialog(false);
              }}
            >
              Activate
            </Button>
            <Button variant="outline" className="flex-1" onClick={() => setFortressDialog(false)}>
              Skip
            </Button>
          </div>
        </DialogContent>
      </Dialog>

      {/* Kinetic Barrier Prompt (L11+) */}
      <Dialog open={kineticDialog} onOpenChange={setKineticDialog}>
        <DialogContent className="max-w-xs">
          <DialogHeader>
            <DialogTitle className="text-center text-primary font-display flex items-center justify-center gap-2">
              <Zap className="h-4 w-4" />
              Kinetic Barrier
            </DialogTitle>
            <DialogDescription className="text-center text-xs">
              Spend Push + Zone components to give +5 AC as a reaction until the start of your next turn.
            </DialogDescription>
          </DialogHeader>
          <div className="flex gap-2">
            <Button
              className="flex-1"
              disabled={!hasComponent('Push') || !hasComponent('Zone')}
              onClick={() => {
                toast.info('Kinetic Barrier activated — DM: construct gains +5 AC reaction. Deduct Push + Zone components.');
                setKineticDialog(false);
              }}
            >
              Activate
              {(!hasComponent('Push') || !hasComponent('Zone')) && (
                <span className="text-micro ml-1">(missing)</span>
              )}
            </Button>
            <Button variant="outline" className="flex-1" onClick={() => setKineticDialog(false)}>
              Skip
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
