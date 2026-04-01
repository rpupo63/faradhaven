import { useState, useMemo, useCallback, useEffect, useRef } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getComponents, createSpell, updateSpell, synthesizeSpell } from '@/lib/api';
import { ElementTable, ComponentKeyboard } from '@/components/arcanum';
import { Button } from '@/components/ui/button';
import { LoadingButton } from '@/components/ui/loading-button';
import { LoadingQuill } from '@/components/LoadingQuill';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Atom, Sparkles, X, Save, Timer, RotateCcw, Keyboard, Layers, AlertCircle, ChevronUp, ChevronDown } from 'lucide-react';
import type { ApiComponent, ApiCharacterComponent, ApiSpell, SpellSynthesis } from '@/types/game';
import { DAMAGE_TYPES } from '@/types/game/state';
import { isValidSpellDuration, isValidSpellDamageDicePair, STANDARD_SPELL_DIE_SIZES } from '@/lib/spellMechanics';
import { LogicConnector, logicVariantFromName } from '@/components/spell-logic/LogicConnector';
import { splitSpellSequenceByLogica, spellChainHasLogica } from '@/lib/spellLogicPhases';

interface CharacterSpellForgeProps {
  availableComponents?: ApiComponent[];
  userId?: string;
  characterId?: string;
  token?: string;
  timerDuration?: number;
  /** Character's component inventory (counts shown for non-pool components) */
  components?: ApiCharacterComponent[];
  currentStability?: number;
  maxStability?: number;
  maxBlueprintSlots?: number;
  spellToEdit?: ApiSpell;
  onSaveComplete?: () => void;
  onClose?: () => void;
}

const ABILITIES = ["STR", "DEX", "CON", "INT", "WIS", "CHA"] as const;

export function CharacterSpellForge({
  availableComponents,
  userId,
  characterId,
  token,
  timerDuration,
  components,
  currentStability,
  maxStability,
  maxBlueprintSlots,
  spellToEdit,
  onSaveComplete,
  onClose,
}: CharacterSpellForgeProps) {
  const queryClient = useQueryClient();
  const isEditMode = !!spellToEdit;

  /** Ordered spell chain (narrative order; duplicate components allowed). */
  const [componentSequence, setComponentSequence] = useState<ApiComponent[]>([]);

  // Spell metadata state
  const [spellName, setSpellName] = useState('');
  const [spellDescription, setSpellDescription] = useState('');
  const [spellType, setSpellType] = useState('');
  /** Empty string = not set in UI; number = feet (0 = self-centered). */
  const [rangeFeet, setRangeFeet] = useState<number | ''>('');
  const [duration, setDuration] = useState('');
  const [concentration, setConcentration] = useState(false);
  const [saveAttr, setSaveAttr] = useState('');
  /** Empty = not set; number = dice count (e.g. 2 for 2d6). */
  const [damageDiceCount, setDamageDiceCount] = useState<number | ''>('');
  /** Empty = not set; die faces (6 = d6). */
  const [damageDieSize, setDamageDieSize] = useState<number | ''>('');
  const [damageType, setDamageType] = useState('');
  const [addModifier, setAddModifier] = useState(false);

  // Override tracking: fields the user has manually edited
  const [overrides, setOverrides] = useState<Set<string>>(() => new Set());

  // Synthesis state
  const [synthesis, setSynthesis] = useState<SpellSynthesis | null>(null);
  const [isSynthesizing, setIsSynthesizing] = useState(false);
  const synthesisTimerRef = useRef<NodeJS.Timeout | null>(null);

const COUNTDOWN_SECONDS = 3; // 3 seconds countdown

  // Timer state
  const [timeLeft, setTimeLeft] = useState<number>(timerDuration || 0);
  const [isTimerRunning, setIsTimerRunning] = useState(false);
  const [countdownValue, setCountdownValue] = useState(COUNTDOWN_SECONDS);
  const [isCountingDown, setIsCountingDown] = useState(false);
  const timerRef = useRef<NodeJS.Timeout | null>(null);

  // View mode state
  const [viewMode, setViewMode] = useState<'table' | 'keyboard'>('table');

  const selectedComponents = componentSequence;

  const clearForm = useCallback(() => {
    setComponentSequence([]);
    setSpellName('');
    setSpellDescription('');
    setSpellType('');
    setRangeFeet('');
    setDuration('');
    setConcentration(false);
    setSaveAttr('');
    setDamageDiceCount('');
    setDamageDieSize('');
    setDamageType('');
    setAddModifier(false);
    setOverrides(new Set());
    setSynthesis(null);
  }, []);

  // Timer effects
  useEffect(() => {
    if (timerDuration) setTimeLeft(timerDuration);
  }, [timerDuration]);

  useEffect(() => {
    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
      if (synthesisTimerRef.current) clearTimeout(synthesisTimerRef.current);
    };
  }, []);

  const { data: allComponents, isLoading, error } = useQuery({
    queryKey: ['components'],
    queryFn: getComponents,
    staleTime: 60_000 * 10,
  });

  // Populate form for edit mode
  useEffect(() => {
    if (isEditMode && allComponents && spellToEdit) {
      setSpellName(spellToEdit.name);
      setSpellDescription(spellToEdit.description || '');

      const byId = new Map((allComponents ?? []).map((c) => [c.id, c]));
      const spellComponents = spellToEdit.components?.length ? spellToEdit.components : [];
      const seq = spellComponents.map((c) => byId.get(c.id) ?? c);
      setComponentSequence(seq);

      setSpellType(spellToEdit.type || '');
      {
        const r = spellToEdit.range;
        if (typeof r === 'number' && !Number.isNaN(r)) setRangeFeet(r);
        else setRangeFeet('');
      }
      setDuration(spellToEdit.duration || '');
      setConcentration(spellToEdit.concentration || false);
      setSaveAttr(spellToEdit.save_attr || '');
      {
        const c = spellToEdit.damage_dice_count;
        const s = spellToEdit.damage_die_size;
        if (typeof c === 'number' && !Number.isNaN(c)) setDamageDiceCount(c);
        else setDamageDiceCount('');
        if (typeof s === 'number' && !Number.isNaN(s)) setDamageDieSize(s);
        else setDamageDieSize('');
      }
      setDamageType(spellToEdit.damage_type || '');
      setAddModifier(spellToEdit.add_modifier || false);
      // In edit mode, assume all fields are overridden
      setOverrides(new Set(['type', 'range', 'duration', 'concentration', 'damageDice', 'damageType']));
    } else if (!isEditMode) {
      clearForm();
    }
  }, [spellToEdit, isEditMode, allComponents, clearForm]);

  // Debounced synthesis call when components change
  useEffect(() => {
    if (synthesisTimerRef.current) clearTimeout(synthesisTimerRef.current);

    const componentIds = selectedComponents.map(c => c.id);
    if (componentIds.length === 0) {
      setSynthesis(null);
      return;
    }

    synthesisTimerRef.current = setTimeout(async () => {
      setIsSynthesizing(true);
      try {
        const result = await synthesizeSpell(componentIds, token);
        setSynthesis(result);

        // Auto-fill non-overridden fields
        if (!overrides.has('type') && result.suggested_type) {
          setSpellType(result.suggested_type);
        }
        if (!overrides.has('range') && result.suggested_range != null) {
          setRangeFeet(result.suggested_range);
        }
        if (!overrides.has('duration') && result.suggested_duration) {
          setDuration(result.suggested_duration || '');
        }
        if (!overrides.has('concentration')) {
          setConcentration(result.suggested_concentration);
        }
        if (!overrides.has('damageDice') && result.suggested_damage_dice_count != null && result.suggested_damage_die_size != null) {
          setDamageDiceCount(result.suggested_damage_dice_count);
          setDamageDieSize(result.suggested_damage_die_size);
        }
        if (!overrides.has('damageType') && result.suggested_damage_type) {
          setDamageType(result.suggested_damage_type);
        }
      } catch (err) { // Added err parameter to catch block
        // Silently fail synthesis — user can still fill manually
        console.error("Spell synthesis failed:", err); // Added console.error
      } finally {
        setIsSynthesizing(false);
      }
    }, 300);
  }, [selectedComponents, token, overrides]); // Added overrides to dependency array

  // Mark field as manually overridden
  const markOverride = (field: string) => {
    setOverrides(prev => new Set(prev).add(field));
  };

  const handleStartTimer = () => {
    if (!timerDuration) return;
    if (timerRef.current) clearInterval(timerRef.current);
    
    setIsCountingDown(true);
    setCountdownValue(COUNTDOWN_SECONDS);
    setIsTimerRunning(false); // Ensure main timer is not running yet

    let currentCountdown = COUNTDOWN_SECONDS;
    timerRef.current = setInterval(() => {
      currentCountdown -= 1;
      setCountdownValue(currentCountdown);
      if (currentCountdown <= 0) {
        if (timerRef.current) clearInterval(timerRef.current);
        setIsCountingDown(false);
        // Start main timer
        const startTime = Date.now();
        const durationMs = timerDuration * 1000;
        setTimeLeft(timerDuration);
        setIsTimerRunning(true);
        timerRef.current = setInterval(() => {
          const elapsed = Date.now() - startTime;
          const remaining = Math.max(0, durationMs - elapsed);
          setTimeLeft(remaining / 1000);
          if (remaining <= 0) {
            if (timerRef.current) clearInterval(timerRef.current);
            setIsTimerRunning(false);
          }
        }, 50);
      }
    }, 1000); // Decrement every second for countdown
  };

  const handleResetTimer = () => {
    if (timerRef.current) clearInterval(timerRef.current);
    setIsTimerRunning(false);
    setIsCountingDown(false);
    setCountdownValue(COUNTDOWN_SECONDS);
    setTimeLeft(timerDuration || 0);
  };

  const level = useMemo(() => synthesis?.level || Math.max(1, selectedComponents.length), [synthesis, selectedComponents]);
  /** Piston Brawler stability spend: same as spell level (component count), not summed tiers. */
  const stabilityCost = selectedComponents.length;

  const getSpellPayload = () => ({
    user_id: userId,
    character_id: characterId,
    name: spellName.trim(),
    description: spellDescription.trim(),
    component_ids: selectedComponents.map(c => c.id),
    slot_level: level,
    type: spellType || 'Utility',
    range: rangeFeet === '' ? undefined : rangeFeet,
    duration: duration.trim() || undefined,
    concentration: concentration,
    save_attr: saveAttr || undefined,
    damage_dice_count: damageDiceCount === '' ? undefined : damageDiceCount,
    damage_die_size: damageDieSize === '' ? undefined : damageDieSize,
    damage_type: damageType || undefined,
    add_modifier: addModifier,
  });

  const mutationOptions = {
    onSuccess: (..._args: unknown[]) => {
      queryClient.invalidateQueries({ queryKey: ['spells'] });
      if (characterId) {
        queryClient.invalidateQueries({ queryKey: ['character-spells', characterId] });
        queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      }
      if (onSaveComplete) onSaveComplete();
    },
  };

  const createSpellMutation = useMutation({
    mutationFn: (spellData: ReturnType<typeof getSpellPayload>) => {
      if (!token) throw new Error('Login required');
      return createSpell(spellData, token);
    },
    ...mutationOptions,
    onSuccess: (...args) => {
      mutationOptions.onSuccess(...args);
      clearForm();
    },
  });

  const updateSpellMutation = useMutation({
    mutationFn: (spellData: ReturnType<typeof getSpellPayload>) => {
      if (!token) throw new Error('Login required');
      if (!spellToEdit) throw new Error('No spell to update');
      return updateSpell({ ...spellData, id: spellToEdit.id }, token);
    },
    ...mutationOptions,
  });

  const handleSave = () => {
    if (!userId) throw new Error('Login required');
    if (!spellName.trim()) throw new Error('Spell name required');
    if (selectedComponents.length === 0) throw new Error('Select at least one component');

    if (currentStability !== undefined && currentStability < stabilityCost) {
      throw new Error('Insufficient Stability Charges');
    }

    const payload = getSpellPayload();
    if (isEditMode) {
      updateSpellMutation.mutate(payload);
    } else {
      createSpellMutation.mutate(payload);
    }
  };

  const availableComponentIds = useMemo(() => {
    if (!availableComponents) return undefined;
    return new Set(availableComponents.map((c) => c.id));
  }, [availableComponents]);

  // Append to ordered sequence (remove / reorder via the crucible strip)
  const handleComponentSelect = useCallback((component: ApiComponent) => {
    if (availableComponentIds && !availableComponentIds.has(component.id)) return;
    setComponentSequence((prev) => [...prev, component]);
  }, [availableComponentIds]);

  const removeSequenceAt = useCallback((index: number) => {
    setComponentSequence((prev) => prev.filter((_, i) => i !== index));
  }, []);

  const moveSequence = useCallback((index: number, delta: -1 | 1) => {
    const to = index + delta;
    setComponentSequence((prev) => {
      if (to < 0 || to >= prev.length) return prev;
      const next = [...prev];
      [next[index], next[to]] = [next[to], next[index]];
      return next;
    });
  }, []);

  const selectedComponentIds = useMemo(() => new Set(selectedComponents.map(c => c.id)), [selectedComponents]);

  const handleClearComponents = () => {
    setComponentSequence([]);
    setOverrides(new Set());
    setSynthesis(null);
  };

  if (isLoading) {
    return (
      <div className="min-h-[40vh] flex flex-col items-center justify-center">
        <LoadingQuill label="Preparing the spell forge..." />
      </div>
    );
  }

  if (error || !allComponents || allComponents.length === 0) {
    return (
      <div className="min-h-[40vh] flex flex-col items-center justify-center gap-4">
        <Atom className="h-16 w-16 text-muted-foreground" />
        <h2 className="text-xl font-tome-heading text-primary">
          {error ? 'The Forge is Cold' : 'Empty Forge'}
        </h2>
        <p className="text-muted-foreground font-tome-marginalia text-center max-w-md">
          {error ? 'Unable to access the component registry.' : 'No spell components have been catalogued yet.'}
        </p>
      </div>
    );
  }

  const availableCount = availableComponentIds?.size ?? allComponents.length;
  const totalCount = allComponents.length;
  const isDurationValid = duration.trim() === '' || isValidSpellDuration(duration.trim());
  const diceCountNum = damageDiceCount === '' ? undefined : damageDiceCount;
  const diceSizeNum = damageDieSize === '' ? undefined : damageDieSize;
  const isDamageDiceValid = isValidSpellDamageDicePair(diceCountNum, diceSizeNum);
  const hasValidationErrors =
    (synthesis && synthesis.validation_errors && synthesis.validation_errors.length > 0) ||
    !isDamageDiceValid ||
    !isDurationValid;
  const canSave =
    userId &&
    token &&
    spellName.trim() &&
    selectedComponents.length > 0 &&
    !hasValidationErrors &&
    isDamageDiceValid &&
    isDurationValid;

  const isPistonBrawler = currentStability !== undefined;
  const hasEnoughStability = isPistonBrawler ? (currentStability ?? 0) >= stabilityCost : true;
  const mutationInProgress = createSpellMutation.isPending || updateSpellMutation.isPending;
  const currentMutationError = isEditMode ? updateSpellMutation.error : createSpellMutation.error;

  return (
    <div className="w-full min-w-0">
      <div className="grid gap-6 lg:grid-cols-[1fr_340px] min-w-0">
        {/* Left: Component Table */}
        <div className="min-w-0">
          <div className="flex flex-col gap-4 mb-4">
            <div className="flex flex-wrap items-center justify-between gap-2 sm:gap-4 min-w-0">
              <div className="flex items-center bg-muted/30 p-1 rounded-md border border-border shrink-0">
                <Button variant={viewMode === 'table' ? 'secondary' : 'ghost'} size="sm" onClick={() => setViewMode('table')} className="gap-2 h-8">
                  <Atom className="h-4 w-4" /> Table
                </Button>
                <Button variant={viewMode === 'keyboard' ? 'secondary' : 'ghost'} size="sm" onClick={() => setViewMode('keyboard')} className="gap-2 h-8">
                  <Keyboard className="h-4 w-4" /> Keyboard
                </Button>
              </div>
              {availableComponentIds && (
                <div className="text-sm text-muted-foreground font-tome-marginalia text-right min-w-0 w-full sm:w-auto">
                  <span className="text-primary font-semibold">{availableCount}</span> of{' '}
                  <span className="text-muted-foreground">{totalCount}</span> available
                </div>
              )}
            </div>
          </div>

          {viewMode === 'table' ? (
            <ElementTable
              components={allComponents}
              availableComponentIds={availableComponentIds}
              selectedComponentIds={selectedComponentIds}
              onComponentClick={handleComponentSelect}
              disableDetailPopup={true}
              characterComponents={components}
            />
          ) : (
            <ComponentKeyboard
              availableComponents={availableComponents || allComponents}
              selectedComponents={selectedComponents}
              onComponentSelect={handleComponentSelect}
              isTimerActive={isTimerRunning}
            />
          )}
        </div>

        {/* Right: Spell Creation Panel */}
        <div className="space-y-4 min-w-0">
          {/* Timer Card */}
          {timerDuration !== undefined && (
            <Card className="arcane-border bg-card sticky top-4 z-10">
              <CardHeader className="pb-3">
                <CardTitle className="flex items-center gap-2 text-lg font-tome-heading text-primary">
                  <Timer className="h-5 w-5" /> Casting Timer
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex flex-col gap-4">
                  <div className="text-center">
                    {isCountingDown ? (
                      <span className="text-5xl font-mono font-bold text-yellow-500 animate-pulse">
                        {countdownValue > 0 ? countdownValue : 'GO!'}
                      </span>
                    ) : (
                      <span className={`text-5xl font-mono font-bold ${timeLeft <= 0.5 && isTimerRunning ? 'text-red-500 animate-pulse' : timeLeft === 0 ? 'text-red-500' : 'text-primary'}`}>
                        {timeLeft.toFixed(2)}s
                      </span>
                    )}
                  </div>
                  <div className="flex gap-2">
                    {!isTimerRunning && !isCountingDown ? (
                      <Button onClick={handleStartTimer} className="w-full bg-primary hover:bg-primary/90 text-primary-foreground font-bold">Start Casting</Button>
                    ) : (
                      <Button onClick={handleResetTimer} variant="destructive" className="w-full">Stop / Reset</Button>
                    )}
                    {!isTimerRunning && !isCountingDown && timeLeft !== timerDuration && (
                      <Button onClick={handleResetTimer} variant="outline" size="icon" title="Reset Timer"><RotateCcw className="h-4 w-4" /></Button>
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          <Card className={`arcane-border bg-card ${timerDuration === undefined ? 'sticky top-4' : ''}`}>
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center justify-between gap-2 text-lg font-tome-heading text-primary">
                <div className="flex items-center gap-2">
                  <Sparkles className="h-5 w-5" />
                  <span>{isEditMode ? 'Edit Spell' : 'Spell Crucible'}</span>
                </div>
                {onClose && (
                  <Button variant="ghost" size="icon" onClick={onClose} className="h-8 w-8">
                    <X className="h-4 w-4" />
                  </Button>
                )}
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Validation Errors */}
              {hasValidationErrors && (
                <div className="p-2 rounded-md bg-red-500/10 border border-red-500/30 space-y-1">
                  {synthesis?.validation_errors?.map((err, i) => (
                    <div key={i} className="flex items-center gap-2 text-xs text-red-500">
                      <AlertCircle className="h-3 w-3 flex-shrink-0" />
                      <span>{err}</span>
                    </div>
                  ))}
                  {!isDamageDiceValid && (
                    <div className="flex items-center gap-2 text-xs text-red-500">
                      <AlertCircle className="h-3 w-3 flex-shrink-0" />
                      <span>Set both dice count and die size (d4–d100), or leave both empty.</span>
                    </div>
                  )}
                </div>
              )}

              {/* Synthesis indicator */}
              {isSynthesizing && (
                <div className="text-center text-xs text-muted-foreground animate-pulse">
                  Synthesizing...
                </div>
              )}

              {/* Piston Brawler resources */}
              {isPistonBrawler && (
                <div className="grid grid-cols-2 gap-4 pb-2 border-b border-border/50">
                  <div className="text-center">
                    <p className="text-[10px] font-tome-marginalia text-muted-foreground uppercase">Stability (components)</p>
                    <p className={`font-display text-2xl ${stabilityCost > (currentStability ?? 0) ? 'text-red-500' : 'text-primary'}`}>{stabilityCost}</p>
                  </div>
                  <div className="text-center">
                    <p className="text-[10px] font-tome-marginalia text-muted-foreground uppercase">Blueprint Slots</p>
                    <p className="font-display text-2xl text-primary">{maxBlueprintSlots !== undefined ? maxBlueprintSlots : '\u221e'}</p>
                  </div>
                </div>
              )}

              {/* Ordered component sequence (narrative / logic flow) */}
              <div className="space-y-2 p-2 rounded-md border border-border bg-muted/20">
                <div className="flex items-center justify-between gap-2 flex-wrap">
                  <div className="flex items-center gap-2 min-w-0">
                    <Label className="text-sm font-tome-marginalia text-muted-foreground shrink-0">
                      Sequence ({selectedComponents.length})
                    </Label>
                    {selectedComponents.length > 0 && (
                      <Badge
                        variant="secondary"
                        className="text-[9px] font-normal font-tome-marginalia shrink-0 border border-border/60"
                      >
                        {spellChainHasLogica(selectedComponents) ? 'Multi-phase' : 'Single phase · default'}
                      </Badge>
                    )}
                  </div>
                  {selectedComponents.length > 0 && (
                    <Button variant="ghost" size="sm" onClick={handleClearComponents} className="h-6 px-2 gap-1 text-xs text-muted-foreground hover:text-destructive shrink-0">
                      <Layers className="h-3 w-3" /> Clear
                    </Button>
                  )}
                </div>
                <div className="text-[10px] text-muted-foreground font-tome-marginalia space-y-1.5">
                  <p>
                    <strong className="text-foreground/90">Default — single phase:</strong> add components in cast order. Table clicks append; arrows reorder; click a chip to remove. No If/Then/Therefore needed for normal spells.
                  </p>
                  <p>
                    <strong className="text-foreground/90">Optional — multi-phase:</strong> place <strong>If</strong>, <strong>Then</strong>, or <strong>Therefore</strong> between narrative beats (e.g. puddle, then freeze). Everything before the first link is phase 1, after it phase 2, and so on. Removing all logic links returns you to the default single-phase strip.
                  </p>
                </div>

                {selectedComponents.length === 0 ? (
                  <p className="text-xs text-muted-foreground italic text-center py-2">
                    Click components on the left to add them
                  </p>
                ) : spellChainHasLogica(selectedComponents) ? (
                  <div className="flex flex-col gap-3">
                    <p className="text-[9px] text-muted-foreground font-tome-marginalia rounded-md bg-muted/40 border border-border/60 px-2 py-1.5">
                      <span className="font-medium text-foreground/85">Multi-phase layout.</span>{' '}
                      Phases are split by your logic links. Delete those connectors to collapse back to the default single-phase view.
                    </p>
                    {splitSpellSequenceByLogica(selectedComponents).map((seg, segKey) => {
                      if (seg.kind === 'logic') {
                        const { comp, index: idx } = seg.item;
                        const logic = logicVariantFromName(comp.name);
                        return (
                          <div
                            key={`logic-${idx}-${segKey}`}
                            className="flex flex-col items-center gap-1 py-0.5"
                          >
                            <div className="h-px w-full max-w-[12rem] bg-gradient-to-r from-transparent via-border to-transparent" aria-hidden />
                            <button
                              type="button"
                              className="inline-flex items-center gap-1 px-2 py-1 rounded-md border border-dashed border-primary/40 bg-muted/30 hover:bg-muted/50"
                              onClick={() => removeSequenceAt(idx)}
                              title={`Remove ${comp.name} (logic link)`}
                            >
                              {logic ? (
                                <LogicConnector variant={logic} />
                              ) : (
                                <span className="text-xs font-tome-marginalia text-primary">{comp.name}</span>
                              )}
                              <div className="flex flex-col border-l border-border/60 pl-1 ml-0.5">
                                <Button
                                  type="button"
                                  variant="ghost"
                                  size="icon"
                                  className="h-5 w-5"
                                  disabled={idx === 0}
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    moveSequence(idx, -1);
                                  }}
                                  aria-label="Move link earlier"
                                >
                                  <ChevronUp className="h-3 w-3" />
                                </Button>
                                <Button
                                  type="button"
                                  variant="ghost"
                                  size="icon"
                                  className="h-5 w-5"
                                  disabled={idx >= selectedComponents.length - 1}
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    moveSequence(idx, 1);
                                  }}
                                  aria-label="Move link later"
                                >
                                  <ChevronDown className="h-3 w-3" />
                                </Button>
                              </div>
                              <X className="h-3 w-3 text-muted-foreground" />
                            </button>
                          </div>
                        );
                      }
                      return (
                        <div
                          key={`phase-${seg.phaseNumber}-${segKey}`}
                          className="rounded-lg border border-primary/25 bg-primary/[0.06] p-2 space-y-2"
                        >
                          <div className="flex items-center justify-between gap-2">
                            <span className="text-[10px] font-tome-marginalia uppercase tracking-wide text-primary/90">
                              Phase {seg.phaseNumber}
                            </span>
                            <span className="text-[9px] text-muted-foreground font-tome-marginalia">
                              Order within this phase matters
                            </span>
                          </div>
                          <div className="flex flex-wrap gap-2 items-stretch">
                            {seg.items.map(({ comp, index: idx }) => (
                              <div
                                key={`${comp.id}-${idx}`}
                                className="inline-flex items-center gap-1 rounded-md border border-border bg-background/90 px-1 py-0.5"
                              >
                                <div className="flex flex-col gap-0 border-r border-border/60 pr-1">
                                  <Button
                                    type="button"
                                    variant="ghost"
                                    size="icon"
                                    className="h-5 w-5"
                                    disabled={idx === 0}
                                    onClick={() => moveSequence(idx, -1)}
                                    aria-label="Move earlier"
                                  >
                                    <ChevronUp className="h-3 w-3" />
                                  </Button>
                                  <Button
                                    type="button"
                                    variant="ghost"
                                    size="icon"
                                    className="h-5 w-5"
                                    disabled={idx >= selectedComponents.length - 1}
                                    onClick={() => moveSequence(idx, 1)}
                                    aria-label="Move later"
                                  >
                                    <ChevronDown className="h-3 w-3" />
                                  </Button>
                                </div>
                                <button
                                  type="button"
                                  className="inline-flex items-center gap-1 px-2 py-0.5 text-xs rounded-full bg-primary/15 text-primary border border-primary/25 hover:bg-primary/25"
                                  onClick={() => removeSequenceAt(idx)}
                                  title={`${comp.name} (Tier ${comp.tier}) — remove`}
                                >
                                  <span className="font-mono font-bold">{comp.symbol}</span>
                                  <span className="hidden sm:inline">{comp.name}</span>
                                  <X className="h-3 w-3" />
                                </button>
                              </div>
                            ))}
                          </div>
                        </div>
                      );
                    })}
                  </div>
                ) : (
                  <div className="rounded-lg border border-border/90 bg-background/50 p-2 space-y-2">
                    <div className="flex items-center justify-between gap-2 flex-wrap">
                      <span className="text-[10px] font-tome-marginalia uppercase tracking-wide text-muted-foreground">
                        Single phase
                      </span>
                      <span className="text-[9px] text-muted-foreground font-tome-marginalia">
                        Standard — one ordered chain
                      </span>
                    </div>
                    <div className="flex flex-wrap gap-2 items-stretch">
                      {selectedComponents.map((comp, idx) => (
                        <div
                          key={`${comp.id}-${idx}`}
                          className="inline-flex items-center gap-1 rounded-md border border-border bg-background/80 px-1 py-0.5"
                        >
                          <div className="flex flex-col gap-0 border-r border-border/60 pr-1">
                            <Button
                              type="button"
                              variant="ghost"
                              size="icon"
                              className="h-5 w-5"
                              disabled={idx === 0}
                              onClick={() => moveSequence(idx, -1)}
                              aria-label="Move earlier"
                            >
                              <ChevronUp className="h-3 w-3" />
                            </Button>
                            <Button
                              type="button"
                              variant="ghost"
                              size="icon"
                              className="h-5 w-5"
                              disabled={idx >= selectedComponents.length - 1}
                              onClick={() => moveSequence(idx, 1)}
                              aria-label="Move later"
                            >
                              <ChevronDown className="h-3 w-3" />
                            </Button>
                          </div>
                          <button
                            type="button"
                            className="inline-flex items-center gap-1 px-2 py-0.5 text-xs rounded-full bg-primary/15 text-primary border border-primary/25 hover:bg-primary/25"
                            onClick={() => removeSequenceAt(idx)}
                            title={`${comp.name} (Tier ${comp.tier}) — remove`}
                          >
                            <span className="font-mono font-bold">{comp.symbol}</span>
                            <span className="hidden sm:inline">{comp.name}</span>
                            <X className="h-3 w-3" />
                          </button>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>

              {/* Spell Name */}
              <div>
                <Label className="text-sm font-tome-marginalia text-muted-foreground mb-2 block">Spell Name</Label>
                <Input value={spellName} onChange={(e) => setSpellName(e.target.value)} placeholder="Enter spell name..." className="bg-background" />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label className="text-sm font-tome-marginalia text-muted-foreground mb-2 block">Level</Label>
                  <div className="bg-muted/30 border border-border rounded-md px-3 py-2 text-sm font-mono font-bold text-primary">{level}</div>
                </div>
                <div>
                  <Label className="text-sm font-tome-marginalia text-muted-foreground mb-2 block">
                    Type {synthesis?.suggested_type && !overrides.has('type') && <Badge variant="secondary" size="tiny" className="ml-1">auto</Badge>}
                  </Label>
                  <Select value={spellType} onValueChange={(v) => { setSpellType(v); markOverride('type'); }}>
                    <SelectTrigger className="bg-background"><SelectValue placeholder="Type" /></SelectTrigger>
                    <SelectContent>
                      <SelectItem value="Attack">Attack</SelectItem>
                      <SelectItem value="Save">Save</SelectItem>
                      <SelectItem value="Effect">Effect</SelectItem>
                      <SelectItem value="Healing">Healing</SelectItem>
                      <SelectItem value="Utility">Utility</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label className="text-sm font-tome-marginalia text-muted-foreground mb-2 block">
                    Range (feet) {synthesis?.suggested_range != null && !overrides.has('range') && <Badge variant="secondary" size="tiny" className="ml-1">auto</Badge>}
                  </Label>
                  <Input
                    type="number"
                    min={0}
                    step={1}
                    value={rangeFeet === '' ? '' : rangeFeet}
                    onChange={(e) => {
                      const raw = e.target.value;
                      if (raw === '') {
                        setRangeFeet('');
                        markOverride('range');
                        return;
                      }
                      const n = parseInt(raw, 10);
                      if (!Number.isNaN(n) && n >= 0) {
                        setRangeFeet(n);
                        markOverride('range');
                      }
                    }}
                    placeholder="0 = self-centered"
                    className="bg-background"
                  />
                </div>
                <div>
                  <Label className="text-sm font-tome-marginalia text-muted-foreground mb-2 block">
                    Duration {synthesis?.suggested_duration && !overrides.has('duration') && <Badge variant="secondary" size="tiny" className="ml-1">auto</Badge>}
                  </Label>
                  <Input
                    value={duration}
                    onChange={(e) => {
                      setDuration(e.target.value);
                      markOverride('duration');
                    }}
                    placeholder="e.g. 1 min, instantaneous, concentration"
                    className={`bg-background ${duration.trim() !== '' && !isDurationValid ? 'border-red-500 ring-1 ring-red-500' : ''}`}
                  />
                  {duration.trim() !== '' && !isDurationValid && (
                    <p className="text-[10px] text-red-500 mt-1 font-tome-marginalia">
                      Use a timed form (1 min, 2 hours), rounds (1 round), or concentration / instantaneous / until dispelled / special / permanent.
                    </p>
                  )}
                </div>
              </div>

              <div className="space-y-4 pt-2 border-t border-border/50">
                {spellType === 'Save' && (
                  <div>
                    <Label className="text-sm font-tome-marginalia text-muted-foreground mb-2 block">Saving Throw Attribute</Label>
                    <Select value={saveAttr} onValueChange={setSaveAttr}>
                      <SelectTrigger className="bg-background"><SelectValue placeholder="Select Attribute" /></SelectTrigger>
                      <SelectContent>
                        {ABILITIES.map(a => (<SelectItem key={a} value={a}>{a}</SelectItem>))}
                      </SelectContent>
                    </Select>
                  </div>
                )}

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label className="text-sm font-tome-marginalia text-muted-foreground mb-2 block">
                      Dice count {synthesis?.suggested_damage_dice_count != null && synthesis?.suggested_damage_die_size != null && !overrides.has('damageDice') && <Badge variant="secondary" size="tiny" className="ml-1">auto</Badge>}
                    </Label>
                    <Input
                      type="number"
                      min={1}
                      max={99}
                      value={damageDiceCount === '' ? '' : damageDiceCount}
                      onChange={e => {
                        const v = e.target.value;
                        if (v === '') {
                          setDamageDiceCount('');
                        } else {
                          const n = parseInt(v, 10);
                          if (!Number.isNaN(n)) setDamageDiceCount(n);
                        }
                        markOverride('damageDice');
                      }}
                      placeholder="e.g. 2"
                      className={`bg-background ${!isDamageDiceValid ? 'border-red-500 ring-1 ring-red-500' : ''}`}
                    />
                  </div>
                  <div>
                    <Label className="text-sm font-tome-marginalia text-muted-foreground mb-2 block">Die size</Label>
                    <Select
                      value={damageDieSize === '' ? '__none__' : String(damageDieSize)}
                      onValueChange={v => {
                        setDamageDieSize(v === '__none__' ? '' : parseInt(v, 10));
                        markOverride('damageDice');
                      }}
                    >
                      <SelectTrigger className={`bg-background ${!isDamageDiceValid ? 'border-red-500 ring-1 ring-red-500' : ''}`}>
                        <SelectValue placeholder="Die" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="__none__">—</SelectItem>
                        {STANDARD_SPELL_DIE_SIZES.map(sz => (
                          <SelectItem key={sz} value={String(sz)}>
                            d{sz}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                <div>
                  <Label className="text-sm font-tome-marginalia text-muted-foreground mb-2 block">
                    Damage Type {synthesis?.suggested_damage_type && !overrides.has('damageType') && <Badge variant="secondary" size="tiny" className="ml-1">auto</Badge>}
                  </Label>
                  <Select value={damageType} onValueChange={(v) => { setDamageType(v); markOverride('damageType'); }}>
                    <SelectTrigger className="bg-background"><SelectValue placeholder="Type" /></SelectTrigger>
                    <SelectContent>
                      {DAMAGE_TYPES.map(t => (<SelectItem key={t} value={t}>{t}</SelectItem>))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="flex flex-col gap-3">
                  <div className="flex items-center space-x-2">
                    <Checkbox id="addMod" checked={addModifier} onCheckedChange={(c) => setAddModifier(c === true)} />
                    <Label htmlFor="addMod" className="cursor-pointer">Add Spellcasting Modifier?</Label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Checkbox id="conc" checked={concentration} onCheckedChange={(c) => { setConcentration(c === true); markOverride('concentration'); }} />
                    <Label htmlFor="conc" className="cursor-pointer">
                      Requires Concentration?
                      {synthesis?.suggested_concentration && !overrides.has('concentration') && <Badge variant="secondary" size="tiny" className="ml-2">auto</Badge>}
                    </Label>
                  </div>
                </div>
              </div>

              {/* Description */}
              <div>
                <Label className="text-sm font-tome-marginalia text-muted-foreground mb-2 block">Description</Label>
                <textarea
                  value={spellDescription}
                  onChange={(e) => setSpellDescription(e.target.value)}
                  placeholder="Describe what your spell does..."
                  rows={4}
                  className="w-full px-3 py-2 text-sm rounded-md border border-border bg-background resize-none focus:outline-none focus:ring-2 focus:ring-primary/50"
                />
              </div>

              {/* Save Button */}
              <LoadingButton
                onClick={handleSave}
                isLoading={mutationInProgress}
                disabled={!canSave || !hasEnoughStability || (maxBlueprintSlots !== undefined && maxBlueprintSlots < 1)}
                className="w-full gap-2"
                loadingText={isEditMode ? "Updating..." : "Forging..."}
              >
                <Save className="h-4 w-4" />
                {isEditMode ? 'Update Spell' : 'Forge Spell'}
              </LoadingButton>

              {currentMutationError && (
                <p className="text-sm text-red-500 text-center">
                  {currentMutationError.message || 'Failed to save spell'}
                </p>
              )}
              {createSpellMutation.isSuccess && (
                <p className="text-sm text-green-500 text-center">Spell created successfully!</p>
              )}
              {updateSpellMutation.isSuccess && (
                <p className="text-sm text-green-500 text-center">Spell updated successfully!</p>
              )}

              {(!userId || !token) && (
                <p className="text-xs text-amber-500 text-center font-tome-marginalia">Login required to save spells</p>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}