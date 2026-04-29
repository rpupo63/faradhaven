import { useState, useMemo, useCallback, useEffect, useRef, type DragEvent } from 'react';
import { createPortal } from 'react-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  getComponents,
  createSpell,
  updateSpell,
  synthesizeSpell,
  getCharacterSpells,
  previewSpellAIOpinion,
  type CreateSpellRequest,
} from '@/lib/api';
import { ElementTable, ComponentKeyboard, ComponentRpgGlyph, ComponentRpgChip } from '@/components/arcanum';
import { COMPONENT_DRAG_MIME } from '@/components/arcanum/ElementTable';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { LoadingButton } from '@/components/ui/loading-button';
import { LoadingQuill } from '@/components/LoadingQuill';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { X, Save, Timer, RotateCcw, Keyboard, Layers, AlertCircle, ChevronUp, ChevronDown, Loader2, Wand2 } from 'lucide-react';
import { RaIcon } from '@/components/ui/RaIcon';
import type { ApiComponent, ApiCharacterComponent, ApiSpell, SpellSynthesis } from '@/types/game';
import { toast } from 'sonner';
import { DAMAGE_TYPES } from '@/types/game/state';
import {
  isValidSpellDuration,
  isValidSpellDamageDicePair,
  STANDARD_SPELL_DIE_SIZES,
  formatSpellRangeFeet,
  formatSpellDamageDice,
} from '@/lib/spellMechanics';
import { LogicConnector, logicVariantFromName } from '@/components/spell-logic/LogicConnector';
import {
  splitSpellSequenceByLogica,
  spellChainHasLogica,
  bucketSequenceByComponentId,
  bucketIndexedPhase,
  deriveSpellPhaseWindows,
  toggleComponentWithinPhase,
  clampPhaseIndex,
  validateFormaScopusPerPhase,
} from '@/lib/spellLogicPhases';
import { findMatchingPreparedSpell } from '@/lib/powderMageSpellMatch';
import { resolveForgeVessel } from '@/components/spell-forge-vessels';
import {
  formatSynthesisConcentrationHint,
  formatSynthesisDamageTypeHint,
  formatSynthesisDiceHint,
  formatSynthesisDurationHint,
  formatSynthesisRangeHint,
  formatSynthesisTypeHint,
} from '@/lib/spellForgeSynthesisHints';
import {
  SpellForgeSynthesisHintBadges,
  SpellForgeSynthesisLabel,
} from '@/components/spell-forge/SpellForgeSynthesisLabel';

interface CharacterSpellForgeV2Props {
  availableComponents?: ApiComponent[];
  userId?: string;
  characterId?: string;
  token?: string;
  /** When true, after the casting timer ends on the keyboard, resolve the spell name or prompt to create. */
  isPowderMage?: boolean;
  /** Powder Mage: number of Speed Dial slots (from class resources). */
  speedDialSlots?: number;
  timerDuration?: number;
  /** Character's component inventory (counts shown for non-pool components) */
  components?: ApiCharacterComponent[];
  currentStability?: number;
  maxStability?: number;
  maxBlueprintSlots?: number;
  spellToEdit?: ApiSpell;
  onSaveComplete?: () => void;
  onClose?: () => void;
  /** Character class display name (e.g. `sheet.class.name`) for themed forge vessel. */
  gameClassName?: string;
}

const ABILITIES = ["STR", "DEX", "CON", "INT", "WIS", "CHA"] as const;

const EMPTY_PREPARED_SPELLS: ApiSpell[] = [];

/** Only send damage_type when it matches server enums (synthesis/UI can briefly hold unknown strings). */
function damageTypeForSpellApi(raw: string): string | undefined {
  const t = raw.trim();
  if (!t) return undefined;
  return DAMAGE_TYPES.includes(t as (typeof DAMAGE_TYPES)[number]) ? t : undefined;
}

export function CharacterSpellForgeV2({
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
  isPowderMage = false,
  speedDialSlots = 0,
  gameClassName,
}: CharacterSpellForgeV2Props) {
  const queryClient = useQueryClient();
  const isEditMode = !!spellToEdit;

  /** Ordered spell chain (narrative order; duplicate components allowed). */
  const [componentSequence, setComponentSequence] = useState<ApiComponent[]>([]);
  const selectedComponents = componentSequence;
  const [activePhaseIndex, setActivePhaseIndex] = useState(0);

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

  /**
   * Manual field flags: when a key is present, spell synthesis will not auto-update that bucket
   * (range, duration, damage dice, type, damage type, concentration, or spell type).
   */
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

  /** Drag-and-drop + merge feedback (Spell Forge 2) */
  const [dragOverZone, setDragOverZone] = useState<string | null>(null);
  /** Viewport center of the drop zone — burst renders there (not anchored only to the Sequence box). */
  const [mergeBurst, setMergeBurst] = useState<{ key: number; cx: number; cy: number } | null>(null);
  const [dragTrail, setDragTrail] = useState<{ component: ApiComponent; x: number; y: number } | null>(null);

  /** Last component id added for class vessel entry animation */
  const [lastAddedComponentId, setLastAddedComponentId] = useState<string | null>(null);
  const prevSequenceLenRef = useRef(0);
  const vesselAnimInitRef = useRef(false);

  /** Powder Mage: timer finished → matched prepared spell or prompt to forge */
  const [lastCastResolution, setLastCastResolution] = useState<
    { kind: 'matched'; spell: ApiSpell } | { kind: 'unmatched' } | null
  >(null);
  const prevTimerRunningRef = useRef(false);
  const spellNameInputRef = useRef<HTMLInputElement>(null);
  const castConfirmNameRef = useRef<HTMLInputElement>(null);

  const { data: characterSpellsData } = useQuery({
    queryKey: ['character-spells', characterId, 'powder-forge'],
    queryFn: () => getCharacterSpells(characterId!, token, { page: 1, limit: 500 }),
    enabled: !!isPowderMage && !!characterId && !!token,
    staleTime: 30_000,
  });
  const preparedSpellsForMatch = characterSpellsData?.spells ?? EMPTY_PREPARED_SPELLS;

  const clearForm = useCallback(() => {
    setComponentSequence([]);
    setActivePhaseIndex(0);
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
    setLastCastResolution(null);
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

  // Powder Mage: when the main casting timer expires naturally, match prepared spells or prompt creation
  useEffect(() => {
    if (!isPowderMage || !timerDuration) {
      prevTimerRunningRef.current = isTimerRunning;
      return;
    }
    const finishedNaturally =
      prevTimerRunningRef.current &&
      !isTimerRunning &&
      !isCountingDown &&
      timeLeft <= 0.05 &&
      componentSequence.length > 0;

    if (finishedNaturally) {
      const matched = findMatchingPreparedSpell(componentSequence, preparedSpellsForMatch);
      if (matched) setLastCastResolution({ kind: 'matched', spell: matched });
      else setLastCastResolution({ kind: 'unmatched' });
    }
    prevTimerRunningRef.current = isTimerRunning;
  }, [
    isPowderMage,
    timerDuration,
    isTimerRunning,
    isCountingDown,
    timeLeft,
    componentSequence,
    preparedSpellsForMatch,
  ]);

  useEffect(() => {
    if (lastCastResolution?.kind === 'unmatched') {
      const t = window.setTimeout(() => castConfirmNameRef.current?.focus(), 100);
      return () => window.clearTimeout(t);
    }
  }, [lastCastResolution]);

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

  // Debounced synthesis call when forge-relevant inputs change
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
        const result = await synthesizeSpell(
          componentIds,
          {
            damageType: damageType || undefined,
            range: rangeFeet === '' ? undefined : rangeFeet,
          },
          token,
        );
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
  }, [selectedComponents, token, overrides, damageType, rangeFeet]);

  /** Mark a synthesis bucket as user-edited so auto-suggestion does not overwrite the displayed value. */
  const markOverride = (field: string) => {
    setOverrides((prev) => new Set(prev).add(field));
  };

  const handleStartTimer = () => {
    if (!timerDuration) return;
    if (timerRef.current) clearInterval(timerRef.current);
    if (synthesisTimerRef.current) {
      clearTimeout(synthesisTimerRef.current);
      synthesisTimerRef.current = null;
    }
    setComponentSequence([]);
    setActivePhaseIndex(0);
    setSynthesis(null);
    setIsSynthesizing(false);
    setOverrides(new Set());
    setLastCastResolution(null);

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
    setLastCastResolution(null);
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
    damage_type: damageTypeForSpellApi(damageType),
    add_modifier: addModifier,
  });

  const [aiDescriptionPending, setAiDescriptionPending] = useState(false);

  const handleGenerateDescription = async () => {
    if (!userId || !token) {
      toast.error('Login required');
      return;
    }
    if (selectedComponents.length === 0) {
      toast.error('Add at least one component first');
      return;
    }
    setAiDescriptionPending(true);
    try {
      const opinion = await previewSpellAIOpinion(getSpellPayload() as CreateSpellRequest, token);
      const rec = opinion.recommended_description?.trim();
      if (rec) {
        setSpellDescription(rec);
        toast.success('Description filled from AI review');
      } else {
        toast.message('No description suggestion', {
          description: opinion.description_opinion || opinion.overall_verdict || 'Try adjusting components or mechanics.',
        });
      }
    } catch (e) {
      toast.error(e instanceof Error ? e.message : 'Failed to generate description');
    } finally {
      setAiDescriptionPending(false);
    }
  };

  const mutationOptions = {
    onSuccess: (..._args: unknown[]) => {
      queryClient.invalidateQueries({ queryKey: ['spells'] });
      if (characterId) {
        queryClient.invalidateQueries({ queryKey: ['character-spells', characterId] });
        queryClient.invalidateQueries({ queryKey: ['characterSpellbook', characterId] });
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
    // Empty array means pool not loaded — treat like undefined (no restriction). Non-empty = allowlist.
    if (!availableComponents?.length) return undefined;
    return new Set(availableComponents.map((c) => c.id));
  }, [availableComponents]);

  const phaseWindows = useMemo(
    () => deriveSpellPhaseWindows(selectedComponents),
    [selectedComponents],
  );
  const resolvedActivePhaseIndex = useMemo(
    () => clampPhaseIndex(activePhaseIndex, phaseWindows),
    [activePhaseIndex, phaseWindows],
  );
  const activePhaseWindow = phaseWindows[resolvedActivePhaseIndex];
  const activePhaseComponents = useMemo(() => {
    if (!activePhaseWindow) return [];
    return activePhaseWindow.indices.map((idx) => selectedComponents[idx]).filter(Boolean);
  }, [activePhaseWindow, selectedComponents]);

  // Toggle within active phase: reclick removes from active phase, otherwise add to active phase end.
  const handleComponentSelect = useCallback((component: ApiComponent) => {
    if (availableComponentIds && !availableComponentIds.has(component.id)) return;
    setComponentSequence((prev) => {
      const result = toggleComponentWithinPhase(prev, component, activePhaseIndex);
      return result.nextSequence;
    });
    if (component.category === 'Logica') {
      setActivePhaseIndex((prev) => prev + 1);
    }
  }, [availableComponentIds, activePhaseIndex]);

  const applyComponentAtPhase = useCallback(
    (component: ApiComponent, phaseIdx: number) => {
      if (availableComponentIds && !availableComponentIds.has(component.id)) return;
      setComponentSequence((prev) => {
        const result = toggleComponentWithinPhase(prev, component, phaseIdx);
        return result.nextSequence;
      });
      if (component.category === 'Logica') {
        setActivePhaseIndex(phaseIdx + 1);
      } else {
        setActivePhaseIndex(phaseIdx);
      }
    },
    [availableComponentIds],
  );

  const resolveComponentFromDrag = useCallback(
    (dt: DataTransfer): ApiComponent | null => {
      const id = dt.getData(COMPONENT_DRAG_MIME) || dt.getData('text/plain').trim();
      if (!id || !allComponents?.length) return null;
      return allComponents.find((c) => c.id === id) ?? null;
    },
    [allComponents],
  );

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

  const selectedComponentIds = useMemo(
    () => new Set(activePhaseComponents.map((c) => c.id)),
    [activePhaseComponents],
  );
  const selectedAnywhereComponentIds = useMemo(
    () => new Set(selectedComponents.map((c) => c.id)),
    [selectedComponents],
  );

  const handleClearComponents = () => {
    setComponentSequence([]);
    setActivePhaseIndex(0);
    setOverrides(new Set());
    setSynthesis(null);
  };

  useEffect(() => {
    if (activePhaseIndex !== resolvedActivePhaseIndex) {
      setActivePhaseIndex(resolvedActivePhaseIndex);
    }
  }, [activePhaseIndex, resolvedActivePhaseIndex]);

  useEffect(() => {
    const clearDragUi = () => {
      setDragOverZone(null);
      setDragTrail(null);
    };
    window.addEventListener('dragend', clearDragUi);
    return () => window.removeEventListener('dragend', clearDragUi);
  }, []);

  /** HTML5 `drag` often reports clientX/Y as 0 (Chromium); `dragover` has reliable cursor coords. */
  const isDraggingForgeIcon = dragTrail !== null;
  useEffect(() => {
    if (!isDraggingForgeIcon) return;
    const onDragOver = (e: globalThis.DragEvent) => {
      e.preventDefault();
      setDragTrail((prev) =>
        prev ? { component: prev.component, x: e.clientX, y: e.clientY } : null,
      );
    };
    document.addEventListener('dragover', onDragOver, true);
    return () => document.removeEventListener('dragover', onDragOver, true);
  }, [isDraggingForgeIcon]);

  useEffect(() => {
    const len = selectedComponents.length;
    const prev = prevSequenceLenRef.current;
    if (!vesselAnimInitRef.current) {
      vesselAnimInitRef.current = true;
      prevSequenceLenRef.current = len;
      return;
    }
    if (len > prev) {
      const last = selectedComponents[len - 1];
      setLastAddedComponentId(last?.id ?? null);
      const t = window.setTimeout(() => setLastAddedComponentId(null), 900);
      prevSequenceLenRef.current = len;
      return () => window.clearTimeout(t);
    }
    if (len < prev) {
      setLastAddedComponentId(null);
    }
    prevSequenceLenRef.current = len;
  }, [selectedComponents]);

  const ForgeVessel = useMemo(() => resolveForgeVessel(gameClassName), [gameClassName]);

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
        <RaIcon name="aura" className="text-6xl text-muted-foreground" />
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
  const localPhaseErrors = validateFormaScopusPerPhase(selectedComponents);
  const hasValidationErrors =
    localPhaseErrors.length > 0 ||
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

  const showCastConfirmDialog =
    isPowderMage && timerDuration !== undefined && lastCastResolution !== null;
  const castConfirmUnmatched = lastCastResolution?.kind === 'unmatched';

  const phaseDropProps = (zoneId: string, phaseIdx: number) => ({
    onDragEnter: (e: DragEvent) => {
      e.preventDefault();
      e.dataTransfer.dropEffect = 'copy';
      setDragOverZone(zoneId);
    },
    onDragOver: (e: DragEvent) => {
      e.preventDefault();
      e.dataTransfer.dropEffect = 'copy';
      setDragOverZone(zoneId);
    },
    onDrop: (e: DragEvent) => {
      e.preventDefault();
      setDragOverZone(null);
      const comp = resolveComponentFromDrag(e.dataTransfer);
      if (!comp) return;
      const el = e.currentTarget as HTMLElement;
      const r = el.getBoundingClientRect();
      setMergeBurst((prev) => ({
        key: (prev?.key ?? 0) + 1,
        cx: r.left + r.width / 2,
        cy: r.top + r.height / 2,
      }));
      applyComponentAtPhase(comp, phaseIdx);
    },
  });

  const forgeOverlays =
    typeof document !== 'undefined'
      ? createPortal(
          <div className="pointer-events-none">
            <AnimatePresence>
              {dragTrail && (
                <motion.div
                  className="fixed z-[80] -translate-x-1/2 -translate-y-1/2"
                  style={{ left: dragTrail.x, top: dragTrail.y }}
                  initial={{ scale: 0.85, opacity: 0.92 }}
                  animate={{ rotate: [-5, 5, -5], scale: 1 }}
                  transition={{
                    rotate: { repeat: Infinity, duration: 1.15, ease: 'easeInOut' },
                    scale: { duration: 0.12 },
                  }}
                >
                  <div className="flex h-14 w-14 items-center justify-center rounded-xl border-2 border-primary/90 bg-background/95 shadow-2xl ring-2 ring-primary/40">
                    <ComponentRpgGlyph
                      component={dragTrail.component}
                      iconClassName="text-2xl text-primary"
                      fallbackClassName="font-mono text-lg font-bold text-primary"
                    />
                  </div>
                </motion.div>
              )}
            </AnimatePresence>
            <AnimatePresence mode="wait">
              {mergeBurst && (
                <motion.div
                  key={mergeBurst.key}
                  className="fixed z-[75] flex -translate-x-1/2 -translate-y-1/2 items-center justify-center"
                  style={{ left: mergeBurst.cx, top: mergeBurst.cy }}
                  initial={{ opacity: 0.95 }}
                  animate={{ opacity: 0 }}
                  transition={{ duration: 0.35 }}
                >
                  <motion.div
                    className="h-28 w-28 rounded-full bg-gradient-to-br from-amber-400/90 via-primary/60 to-transparent blur-md"
                    initial={{ scale: 0.15, opacity: 1 }}
                    animate={{ scale: 2.8, opacity: 0 }}
                    transition={{ duration: 0.48, ease: 'easeOut' }}
                  />
                </motion.div>
              )}
            </AnimatePresence>
          </div>,
          document.body,
        )
      : null;

  return (
    <div className="relative w-full min-w-0">
      {forgeOverlays}

      <div className="grid gap-6 lg:grid-cols-[minmax(0,1.2fr)_minmax(260px,360px)] min-w-0">
        {/* Primary: spell being forged — sticky while page scrolls */}
        <div className="flex flex-col gap-4 min-w-0 lg:sticky lg:top-4 lg:self-start lg:max-h-[80dvh] lg:min-h-0">
          {/* Timer Card */}
          {timerDuration !== undefined && (
            <Card className="arcane-border bg-card shrink-0 sticky top-4 z-10 lg:static">
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


          <Card
            className={`arcane-border bg-card flex min-h-0 flex-1 flex-col overflow-hidden lg:min-h-0 ${
              timerDuration === undefined ? 'max-lg:sticky max-lg:top-4' : ''
            }`}
          >
            <CardHeader className="shrink-0 pb-3">
              <CardTitle className="flex items-center justify-between gap-2 text-lg font-tome-heading text-primary">
                <div className="flex flex-wrap items-center gap-2">
                  <RaIcon name="aura" className="text-base" />
                  <span>{isEditMode ? 'Edit Spell' : 'Spell Crucible'}</span>
                  {!isEditMode && (
                    <Badge variant="outline" className="text-tiny font-tome-marginalia border-primary/40">
                      Forge 2
                    </Badge>
                  )}
                </div>
                {onClose && (
                  <Button variant="ghost" size="icon" onClick={onClose} className="h-8 w-8">
                    <X className="h-4 w-4" />
                  </Button>
                )}
              </CardTitle>
            </CardHeader>
            <CardContent className="min-h-0 flex-1 space-y-4 overflow-y-auto overscroll-contain">
              {/* Drop onto the vessel preview or the sequence box below — same phase as active. */}
              <div
                className={cn(
                  'relative rounded-xl transition-shadow',
                  dragOverZone === 'forge-vessel' &&
                    'ring-2 ring-primary/45 ring-offset-2 ring-offset-background',
                )}
                {...phaseDropProps('forge-vessel', resolvedActivePhaseIndex)}
              >
                <ForgeVessel
                  sequence={selectedComponents}
                  lastAddedComponentId={lastAddedComponentId}
                />
              </div>

              {/* Validation Errors */}
              {hasValidationErrors && (
                <div className="p-2 rounded-md bg-red-500/10 border border-red-500/30 space-y-1">
                  {localPhaseErrors.map((err, i) => (
                    <div key={`local-${i}`} className="flex items-center gap-2 text-xs text-red-500">
                      <AlertCircle className="h-3 w-3 flex-shrink-0" />
                      <span>{err}</span>
                    </div>
                  ))}
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
                    <p className="text-micro font-tome-marginalia text-muted-foreground uppercase">Stability (components)</p>
                    <p className={`font-display text-2xl ${stabilityCost > (currentStability ?? 0) ? 'text-red-500' : 'text-primary'}`}>{stabilityCost}</p>
                  </div>
                  <div className="text-center">
                    <p className="text-micro font-tome-marginalia text-muted-foreground uppercase">Blueprint Slots</p>
                    <p className="font-display text-2xl text-primary">{maxBlueprintSlots !== undefined ? maxBlueprintSlots : '\u221e'}</p>
                  </div>
                </div>
              )}

              {/* Ordered component sequence (narrative / logic flow) */}
              <div className="relative space-y-2 overflow-hidden rounded-md border border-border bg-muted/20 p-2">
                <div className="flex items-center justify-between gap-2 flex-wrap">
                  <div className="flex items-center gap-2 min-w-0">
                    <Label className="text-sm font-tome-marginalia text-muted-foreground shrink-0">
                      Sequence ({selectedComponents.length})
                    </Label>
                    {selectedComponents.length > 0 && (
                      <TooltipProvider delayDuration={150}>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <Badge
                              variant="secondary"
                              className="text-tiny font-normal font-tome-marginalia shrink-0 border border-border/60 cursor-help"
                            >
                              {spellChainHasLogica(selectedComponents) ? 'Multi-phase' : 'Single phase · buckets'}
                            </Badge>
                          </TooltipTrigger>
                          <TooltipContent className="max-w-xs text-xs font-tome-marginalia">
                            {spellChainHasLogica(selectedComponents)
                              ? 'Optional multi-phase mode: If / Then / Therefore links split the sequence into phases. Components bucket within each phase.'
                              : 'Default single-phase mode: components are bucketed by symbol and order is ignored.'}
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                    )}
                  </div>
                  {selectedComponents.length > 0 && (
                    <Button variant="ghost" size="sm" onClick={handleClearComponents} className="h-6 px-2 gap-1 text-xs text-muted-foreground hover:text-destructive shrink-0">
                      <Layers className="h-3 w-3" /> Clear
                    </Button>
                  )}
                </div>
                {selectedComponents.length === 0 ? (
                  <div
                    className={cn(
                      'relative z-10 flex min-h-[7rem] flex-col items-center justify-center rounded-md border-2 border-dashed px-2 py-6 transition-colors',
                      dragOverZone === 'seq-empty'
                        ? 'border-primary bg-primary/10'
                        : 'border-border/60 bg-muted/30',
                    )}
                    {...phaseDropProps('seq-empty', 0)}
                  >
                    <p className="text-xs text-muted-foreground italic text-center px-2">
                      Drag component icons here or click palette tiles to add them
                    </p>
                  </div>
                ) : spellChainHasLogica(selectedComponents) ? (
                  <div className="flex flex-col gap-3">
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
                          className={cn(
                            'relative z-10 rounded-lg space-y-2 border p-2 transition-shadow',
                            resolvedActivePhaseIndex + 1 === seg.phaseNumber
                              ? 'border-primary bg-primary/[0.12]'
                              : 'border-primary/25 bg-primary/[0.06]',
                            dragOverZone === `phase-${seg.phaseNumber}` &&
                              'ring-2 ring-primary/55 ring-offset-2 ring-offset-background',
                          )}
                          onClick={() => setActivePhaseIndex(seg.phaseNumber - 1)}
                          {...phaseDropProps(`phase-${seg.phaseNumber}`, seg.phaseNumber - 1)}
                        >
                          <div className="flex items-center justify-between gap-2">
                            <span className="text-micro font-tome-marginalia uppercase tracking-wide text-primary/90">
                              Phase {seg.phaseNumber}
                            </span>
                            <TooltipProvider delayDuration={150}>
                              <Tooltip>
                                <TooltipTrigger asChild>
                                  <span className="text-tiny text-muted-foreground font-tome-marginalia cursor-help">
                                    {resolvedActivePhaseIndex + 1 === seg.phaseNumber
                                      ? 'Active phase'
                                      : 'Bucketed — order ignored here'}
                                  </span>
                                </TooltipTrigger>
                                <TooltipContent className="max-w-xs text-xs font-tome-marginalia">
                                  {resolvedActivePhaseIndex + 1 === seg.phaseNumber
                                    ? `Active phase for clicks. Reclicking a component toggles it in Phase ${seg.phaseNumber}.`
                                    : 'Within a phase, duplicate components are bucketed and order is ignored.'}
                                </TooltipContent>
                              </Tooltip>
                            </TooltipProvider>
                          </div>
                          <div className="flex flex-wrap gap-2 items-stretch">
                            {bucketIndexedPhase(seg.items).map((bucket) => {
                              const n = bucket.indices.length;
                              const removeOne = () =>
                                removeSequenceAt(bucket.indices[bucket.indices.length - 1]);
                              return (
                                <button
                                  key={bucket.comp.id}
                                  type="button"
                                  className="inline-flex items-center gap-1 px-2 py-0.5 text-xs rounded-full bg-primary/15 text-primary border border-primary/25 hover:bg-primary/25"
                                  onClick={removeOne}
                                  title={`${bucket.comp.name} (Tier ${bucket.comp.tier}) — remove one`}
                                >
                                  <ComponentRpgGlyph
                                    component={bucket.comp}
                                    iconClassName="text-sm"
                                    fallbackClassName="font-mono font-bold text-xs"
                                  />
                                  <span className="hidden sm:inline">{bucket.comp.name}</span>
                                  {n > 1 && (
                                    <Badge variant="outline" className="h-4 px-1 text-[10px] font-mono">
                                      ×{n}
                                    </Badge>
                                  )}
                                  <X className="h-3 w-3" />
                                </button>
                              );
                            })}
                          </div>
                        </div>
                      );
                    })}
                  </div>
                ) : (
                  <div
                    className={cn(
                      'relative z-10 space-y-2 rounded-lg border border-border/90 bg-background/50 p-2 transition-shadow',
                      dragOverZone === 'seq-single' &&
                        'ring-2 ring-primary/50 ring-offset-2 ring-offset-background',
                    )}
                    {...phaseDropProps('seq-single', 0)}
                  >
                    <div className="flex items-center justify-between gap-2 flex-wrap">
                      <span className="text-micro font-tome-marginalia uppercase tracking-wide text-muted-foreground">
                        Single phase
                      </span>
                      <TooltipProvider delayDuration={150}>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <span className="text-tiny text-muted-foreground font-tome-marginalia cursor-help">
                              {resolvedActivePhaseIndex === 0 ? 'Active phase · bucketed — order ignored' : 'Bucketed — order ignored'}
                            </span>
                          </TooltipTrigger>
                          <TooltipContent className="max-w-xs text-xs font-tome-marginalia">
                            Default single-phase mode: clicks toggle components in this one phase, and component order is ignored.
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                    </div>
                    <div className="flex flex-wrap gap-2 items-stretch">
                      {bucketSequenceByComponentId(selectedComponents).map((bucket) => {
                        const n = bucket.indices.length;
                        const removeOne = () =>
                          removeSequenceAt(bucket.indices[bucket.indices.length - 1]);
                        return (
                          <button
                            key={bucket.comp.id}
                            type="button"
                            className="inline-flex items-center gap-1 px-2 py-0.5 text-xs rounded-full bg-primary/15 text-primary border border-primary/25 hover:bg-primary/25"
                            onClick={removeOne}
                            title={`${bucket.comp.name} (Tier ${bucket.comp.tier}) — remove one`}
                          >
                            <ComponentRpgGlyph
                              component={bucket.comp}
                              iconClassName="text-sm"
                              fallbackClassName="font-mono font-bold text-xs"
                            />
                            <span className="hidden sm:inline">{bucket.comp.name}</span>
                            {n > 1 && (
                              <Badge variant="outline" className="h-4 px-1 text-[10px] font-mono">
                                ×{n}
                              </Badge>
                            )}
                            <X className="h-3 w-3" />
                          </button>
                        );
                      })}
                    </div>
                  </div>
                )}
              </div>

              {/* Spell Name */}
              <div>
                <Label className="text-sm font-tome-marginalia text-muted-foreground mb-2 block">Spell Name</Label>
                <Input
                  ref={spellNameInputRef}
                  value={spellName}
                  onChange={(e) => setSpellName(e.target.value)}
                  placeholder="Enter spell name..."
                  className="bg-background"
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label className="text-sm font-tome-marginalia text-muted-foreground mb-2 block">Level</Label>
                  <div className="bg-muted/30 border border-border rounded-md px-3 py-2 text-sm font-mono font-bold text-primary">{level}</div>
                </div>
                <div>
                  <SpellForgeSynthesisLabel
                    showAutoBadge={Boolean(synthesis?.suggested_type && !overrides.has('type'))}
                    recommendedSummary={formatSynthesisTypeHint(synthesis)}
                    hasManualOverride={overrides.has('type')}
                  >
                    Type
                  </SpellForgeSynthesisLabel>
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
                  <SpellForgeSynthesisLabel
                    showAutoBadge={synthesis?.suggested_range != null && !overrides.has('range')}
                    recommendedSummary={formatSynthesisRangeHint(synthesis)}
                    hasManualOverride={overrides.has('range')}
                  >
                    Range (feet)
                  </SpellForgeSynthesisLabel>
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
                  <SpellForgeSynthesisLabel
                    showAutoBadge={Boolean(synthesis?.suggested_duration && !overrides.has('duration'))}
                    recommendedSummary={formatSynthesisDurationHint(synthesis)}
                    hasManualOverride={overrides.has('duration')}
                  >
                    Duration
                  </SpellForgeSynthesisLabel>
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
                    <p className="text-micro text-red-500 mt-1 font-tome-marginalia">
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
                    <SpellForgeSynthesisLabel
                      showAutoBadge={
                        synthesis?.suggested_damage_dice_count != null &&
                        synthesis?.suggested_damage_die_size != null &&
                        !overrides.has('damageDice')
                      }
                      recommendedSummary={formatSynthesisDiceHint(synthesis)}
                      hasManualOverride={overrides.has('damageDice')}
                    >
                      Dice count
                    </SpellForgeSynthesisLabel>
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
                    <SpellForgeSynthesisLabel
                      showAutoBadge={false}
                      recommendedSummary={formatSynthesisDiceHint(synthesis)}
                      hasManualOverride={overrides.has('damageDice')}
                    >
                      Die size
                    </SpellForgeSynthesisLabel>
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
                  <SpellForgeSynthesisLabel
                    showAutoBadge={Boolean(synthesis?.suggested_damage_type && !overrides.has('damageType'))}
                    recommendedSummary={formatSynthesisDamageTypeHint(synthesis)}
                    hasManualOverride={overrides.has('damageType')}
                  >
                    Damage Type
                  </SpellForgeSynthesisLabel>
                  <Select
                    value={damageType === '' ? '__none__' : damageType}
                    onValueChange={(v) => {
                      setDamageType(v === '__none__' ? '' : v);
                      markOverride('damageType');
                    }}
                  >
                    <SelectTrigger className="bg-background"><SelectValue placeholder="None (no damage)" /></SelectTrigger>
                    <SelectContent>
                      <SelectItem value="__none__">—</SelectItem>
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
                    <label htmlFor="conc" className="cursor-pointer text-sm font-tome-marginalia inline-flex flex-wrap items-center gap-1.5">
                      Requires Concentration?
                      <SpellForgeSynthesisHintBadges
                        showAutoBadge={Boolean(synthesis?.suggested_concentration && !overrides.has('concentration'))}
                        recommendedSummary={formatSynthesisConcentrationHint(synthesis)}
                        hasManualOverride={overrides.has('concentration')}
                      />
                    </label>
                  </div>
                </div>
              </div>

              {/* Description */}
              <div>
                <div className="mb-2 flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                  <Label className="text-sm font-tome-marginalia text-muted-foreground block sm:mb-0">
                    Description
                  </Label>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    className="shrink-0 gap-1.5 font-tome-marginalia"
                    disabled={
                      aiDescriptionPending ||
                      !userId ||
                      !token ||
                      selectedComponents.length === 0
                    }
                    onClick={handleGenerateDescription}
                    title="Same AI review as GM spell tools — fills recommended description when available"
                  >
                    {aiDescriptionPending ? (
                      <Loader2 className="h-3.5 w-3.5 animate-spin" aria-hidden />
                    ) : (
                      <Wand2 className="h-3.5 w-3.5" aria-hidden />
                    )}
                    Generate description
                  </Button>
                </div>
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

        {/* Components palette (secondary column) — card chrome matches Spell Crucible */}
        <Card className="arcane-border bg-card min-w-0 flex flex-col lg:max-h-[80dvh] lg:min-h-0 overflow-hidden">
          <CardHeader className="shrink-0 pb-3">
            <CardTitle className="flex items-center gap-2 text-lg font-tome-heading text-primary">
              <RaIcon name="crystal-wand" className="text-base" />
              <span>Components</span>
            </CardTitle>
          </CardHeader>
          <CardContent className="min-h-0 flex-1 space-y-4 overflow-y-auto overscroll-contain pt-0">
            <div className="flex flex-wrap items-center justify-between gap-2 sm:gap-4 min-w-0">
              <div className="flex items-center bg-muted/30 p-1 rounded-md border border-border shrink-0">
                <Button variant={viewMode === 'table' ? 'secondary' : 'ghost'} size="sm" onClick={() => setViewMode('table')} className="gap-2 h-8">
                  <RaIcon name="aura" className="text-sm" /> Table
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

            {viewMode === 'table' ? (
              <ElementTable
                components={allComponents}
                availableComponentIds={availableComponentIds}
                selectedComponentIds={selectedComponentIds}
                selectedAnywhereComponentIds={selectedAnywhereComponentIds}
                onComponentClick={handleComponentSelect}
                disableDetailPopup={true}
                characterComponents={components}
                enableComponentDrag
                onComponentDragStart={(comp, x, y) => setDragTrail({ component: comp, x, y })}
                onComponentDragEnd={() => setDragTrail(null)}
              />
            ) : (
              <ComponentKeyboard
                availableComponents={availableComponents || allComponents}
                selectedComponents={activePhaseComponents}
                selectedAnywhereComponentIds={selectedAnywhereComponentIds}
                onComponentSelect={handleComponentSelect}
                isTimerActive={isTimerRunning}
              />
            )}
          </CardContent>
        </Card>
      </div>

      <Dialog
        open={showCastConfirmDialog}
        onOpenChange={(open) => {
          if (!open) setLastCastResolution(null);
        }}
      >
        <DialogContent
          tome
          hideCloseButton={castConfirmUnmatched}
          className="max-w-2xl"
          onPointerDownOutside={(e) => {
            if (castConfirmUnmatched) e.preventDefault();
          }}
          onEscapeKeyDown={(e) => {
            if (castConfirmUnmatched) e.preventDefault();
          }}
        >
          {lastCastResolution?.kind === 'matched' && (
            <>
              <DialogHeader>
                <DialogTitle className="font-tome-heading text-primary">Cast complete</DialogTitle>
                <DialogDescription>
                  This component sequence matches a spell you already have prepared.
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-4 py-1">
                <p className="text-xl font-tome-heading text-primary">{lastCastResolution.spell.name}</p>
                <dl className="grid grid-cols-1 gap-3 text-sm sm:grid-cols-2">
                  <div>
                    <dt className="text-micro font-tome-marginalia uppercase text-muted-foreground">Level</dt>
                    <dd className="font-mono font-semibold text-foreground">
                      {lastCastResolution.spell.slot_level ?? lastCastResolution.spell.level}
                    </dd>
                  </div>
                  <div>
                    <dt className="text-micro font-tome-marginalia uppercase text-muted-foreground">Type</dt>
                    <dd>{lastCastResolution.spell.type ?? '—'}</dd>
                  </div>
                  <div>
                    <dt className="text-micro font-tome-marginalia uppercase text-muted-foreground">Range</dt>
                    <dd>{formatSpellRangeFeet(lastCastResolution.spell.range) || '—'}</dd>
                  </div>
                  <div>
                    <dt className="text-micro font-tome-marginalia uppercase text-muted-foreground">Duration</dt>
                    <dd>{lastCastResolution.spell.duration?.trim() ? lastCastResolution.spell.duration : '—'}</dd>
                  </div>
                </dl>
                {lastCastResolution.spell.damage_dice_count != null &&
                  lastCastResolution.spell.damage_die_size != null && (
                    <p className="text-sm">
                      <span className="text-muted-foreground font-tome-marginalia">Damage </span>
                      <span className="font-mono font-semibold">
                        {formatSpellDamageDice(
                          lastCastResolution.spell.damage_dice_count,
                          lastCastResolution.spell.damage_die_size,
                        )}
                      </span>
                      {lastCastResolution.spell.damage_type && (
                        <span className="text-muted-foreground"> {lastCastResolution.spell.damage_type}</span>
                      )}
                    </p>
                  )}
                {(lastCastResolution.spell.concentration || lastCastResolution.spell.save_attr) && (
                  <p className="text-xs text-muted-foreground font-tome-marginalia">
                    {lastCastResolution.spell.concentration && <span>Concentration. </span>}
                    {lastCastResolution.spell.save_attr && (
                      <span>Save {lastCastResolution.spell.save_attr}</span>
                    )}
                  </p>
                )}
                {lastCastResolution.spell.description?.trim() && (
                  <div className="rounded-md border border-border/60 bg-muted/20 p-3 text-sm text-foreground/90 max-h-40 overflow-y-auto">
                    {lastCastResolution.spell.description}
                  </div>
                )}
              </div>
              <DialogFooter>
                <Button type="button" onClick={() => setLastCastResolution(null)}>
                  OK
                </Button>
              </DialogFooter>
            </>
          )}

          {lastCastResolution?.kind === 'unmatched' && (
            <>
              <DialogHeader>
                <DialogTitle className="font-tome-heading text-primary">Confirm new spell</DialogTitle>
                <DialogDescription>
                  No prepared spell matches this sequence. Name it and set mechanics, then forge it—or cancel to
                  edit the sequence in the crucible.
                </DialogDescription>
              </DialogHeader>

              <div className="max-h-[min(60vh,28rem)] space-y-4 overflow-y-auto pr-1">
                <div className="rounded-md border border-border/60 bg-muted/20 p-2">
                  <p className="text-micro font-tome-marginalia uppercase text-muted-foreground mb-2">Sequence</p>
                  <div className="flex flex-wrap gap-2">
                    {selectedComponents.map((comp, idx) => (
                      <Badge key={`${comp.id}-${idx}`} variant="secondary" className="inline-flex items-center gap-1 font-mono">
                        <ComponentRpgGlyph
                          component={comp}
                          iconClassName="text-sm"
                          fallbackClassName="font-bold"
                        />
                        <span>{comp.name}</span>
                      </Badge>
                    ))}
                  </div>
                </div>

                <div>
                  <Label className="text-sm font-tome-marginalia text-muted-foreground mb-2 block">Spell Name</Label>
                  <Input
                    ref={castConfirmNameRef}
                    value={spellName}
                    onChange={(e) => setSpellName(e.target.value)}
                    placeholder="Enter spell name..."
                    className="bg-background"
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label className="text-sm font-tome-marginalia text-muted-foreground mb-2 block">Level</Label>
                    <div className="bg-muted/30 border border-border rounded-md px-3 py-2 text-sm font-mono font-bold text-primary">
                      {level}
                    </div>
                  </div>
                  <div>
                    <SpellForgeSynthesisLabel
                      showAutoBadge={Boolean(synthesis?.suggested_type && !overrides.has('type'))}
                      recommendedSummary={formatSynthesisTypeHint(synthesis)}
                      hasManualOverride={overrides.has('type')}
                    >
                      Type
                    </SpellForgeSynthesisLabel>
                    <Select
                      value={spellType}
                      onValueChange={(v) => {
                        setSpellType(v);
                        markOverride('type');
                      }}
                    >
                      <SelectTrigger className="bg-background">
                        <SelectValue placeholder="Type" />
                      </SelectTrigger>
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
                    <SpellForgeSynthesisLabel
                      showAutoBadge={synthesis?.suggested_range != null && !overrides.has('range')}
                      recommendedSummary={formatSynthesisRangeHint(synthesis)}
                      hasManualOverride={overrides.has('range')}
                    >
                      Range (feet)
                    </SpellForgeSynthesisLabel>
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
                    <SpellForgeSynthesisLabel
                      showAutoBadge={Boolean(synthesis?.suggested_duration && !overrides.has('duration'))}
                      recommendedSummary={formatSynthesisDurationHint(synthesis)}
                      hasManualOverride={overrides.has('duration')}
                    >
                      Duration
                    </SpellForgeSynthesisLabel>
                    <Input
                      value={duration}
                      onChange={(e) => {
                        setDuration(e.target.value);
                        markOverride('duration');
                      }}
                      placeholder="e.g. 1 min, instantaneous"
                      className={`bg-background ${duration.trim() !== '' && !isDurationValid ? 'border-red-500 ring-1 ring-red-500' : ''}`}
                    />
                  </div>
                </div>

                {spellType === 'Save' && (
                  <div>
                    <Label className="text-sm font-tome-marginalia text-muted-foreground mb-2 block">
                      Saving Throw
                    </Label>
                    <Select value={saveAttr} onValueChange={setSaveAttr}>
                      <SelectTrigger className="bg-background">
                        <SelectValue placeholder="Attribute" />
                      </SelectTrigger>
                      <SelectContent>
                        {ABILITIES.map((a) => (
                          <SelectItem key={a} value={a}>
                            {a}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                )}

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <SpellForgeSynthesisLabel
                      showAutoBadge={
                        synthesis?.suggested_damage_dice_count != null &&
                        synthesis?.suggested_damage_die_size != null &&
                        !overrides.has('damageDice')
                      }
                      recommendedSummary={formatSynthesisDiceHint(synthesis)}
                      hasManualOverride={overrides.has('damageDice')}
                    >
                      Dice count
                    </SpellForgeSynthesisLabel>
                    <Input
                      type="number"
                      min={1}
                      max={99}
                      value={damageDiceCount === '' ? '' : damageDiceCount}
                      onChange={(e) => {
                        const v = e.target.value;
                        if (v === '') setDamageDiceCount('');
                        else {
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
                    <SpellForgeSynthesisLabel
                      showAutoBadge={false}
                      recommendedSummary={formatSynthesisDiceHint(synthesis)}
                      hasManualOverride={overrides.has('damageDice')}
                    >
                      Die size
                    </SpellForgeSynthesisLabel>
                    <Select
                      value={damageDieSize === '' ? '__none__' : String(damageDieSize)}
                      onValueChange={(v) => {
                        setDamageDieSize(v === '__none__' ? '' : parseInt(v, 10));
                        markOverride('damageDice');
                      }}
                    >
                      <SelectTrigger
                        className={`bg-background ${!isDamageDiceValid ? 'border-red-500 ring-1 ring-red-500' : ''}`}
                      >
                        <SelectValue placeholder="Die" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="__none__">—</SelectItem>
                        {STANDARD_SPELL_DIE_SIZES.map((sz) => (
                          <SelectItem key={sz} value={String(sz)}>
                            d{sz}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                <div>
                  <SpellForgeSynthesisLabel
                    showAutoBadge={Boolean(synthesis?.suggested_damage_type && !overrides.has('damageType'))}
                    recommendedSummary={formatSynthesisDamageTypeHint(synthesis)}
                    hasManualOverride={overrides.has('damageType')}
                  >
                    Damage type
                  </SpellForgeSynthesisLabel>
                  <Select
                    value={damageType === '' ? '__none__' : damageType}
                    onValueChange={(v) => {
                      setDamageType(v === '__none__' ? '' : v);
                      markOverride('damageType');
                    }}
                  >
                    <SelectTrigger className="bg-background">
                      <SelectValue placeholder="None (no damage)" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="__none__">—</SelectItem>
                      {DAMAGE_TYPES.map((t) => (
                        <SelectItem key={t} value={t}>
                          {t}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="flex flex-col gap-3">
                  <div className="flex items-center space-x-2">
                    <Checkbox id="castConfirmAddMod" checked={addModifier} onCheckedChange={(c) => setAddModifier(c === true)} />
                    <Label htmlFor="castConfirmAddMod" className="cursor-pointer">
                      Add spellcasting modifier to damage
                    </Label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id="castConfirmConc"
                      checked={concentration}
                      onCheckedChange={(c) => {
                        setConcentration(c === true);
                        markOverride('concentration');
                      }}
                    />
                    <label htmlFor="castConfirmConc" className="cursor-pointer text-sm font-tome-marginalia inline-flex flex-wrap items-center gap-1.5">
                      Requires concentration
                      <SpellForgeSynthesisHintBadges
                        showAutoBadge={Boolean(synthesis?.suggested_concentration && !overrides.has('concentration'))}
                        recommendedSummary={formatSynthesisConcentrationHint(synthesis)}
                        hasManualOverride={overrides.has('concentration')}
                      />
                    </label>
                  </div>
                </div>

                <div>
                  <div className="mb-2 flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                    <Label className="text-sm font-tome-marginalia text-muted-foreground block sm:mb-0">
                      Description
                    </Label>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      className="shrink-0 gap-1.5 font-tome-marginalia"
                      disabled={
                        aiDescriptionPending ||
                        !userId ||
                        !token ||
                        selectedComponents.length === 0
                      }
                      onClick={handleGenerateDescription}
                      title="Same AI review as GM spell tools — fills recommended description when available"
                    >
                      {aiDescriptionPending ? (
                        <Loader2 className="h-3.5 w-3.5 animate-spin" aria-hidden />
                      ) : (
                        <Wand2 className="h-3.5 w-3.5" aria-hidden />
                      )}
                      Generate description
                    </Button>
                  </div>
                  <textarea
                    value={spellDescription}
                    onChange={(e) => setSpellDescription(e.target.value)}
                    placeholder="Describe what your spell does..."
                    rows={3}
                    className="w-full px-3 py-2 text-sm rounded-md border border-border bg-background resize-none focus:outline-none focus:ring-2 focus:ring-primary/50"
                  />
                </div>
              </div>

              <DialogFooter className="gap-2 sm:gap-0">
                <Button type="button" variant="outline" onClick={() => setLastCastResolution(null)}>
                  Cancel
                </Button>
                <LoadingButton
                  type="button"
                  onClick={() => {
                    try {
                      handleSave();
                    } catch (err) {
                      toast.error(err instanceof Error ? err.message : 'Cannot forge spell');
                    }
                  }}
                  isLoading={mutationInProgress}
                  disabled={!canSave || !hasEnoughStability || (maxBlueprintSlots !== undefined && maxBlueprintSlots < 1) || isEditMode}
                  loadingText="Forging..."
                  className="gap-2"
                >
                  <Save className="h-4 w-4" />
                  Forge spell
                </LoadingButton>
              </DialogFooter>
            </>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}