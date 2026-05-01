import { useState, useCallback } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  getCharacterSpellbook,
  castSpell,
  getSpeedDial,
  saveSpeedDialSlot,
  clearSpeedDialSlot,
  type CharacterSpellbookScope,
  getComponents,
} from '@/lib/api';
import { ApiComponent, ApiSpell, ApiCharacterComponent, CastSpellResponse, ApiCharacterSheet, ApiSavedSpell } from '@/types/game';
import { getSpellComponentCount, hasAllComponents } from '@/lib/spellUtils';
import { toast } from 'sonner';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { SpellListPagination } from '@/components/SpellListPagination';
import { PreparedSpellCard } from '@/components/PreparedSpells';
import { buildCastToast } from '@/lib/toastUtils';
import { CharacterSpellForge } from '@/components/CharacterSpellForge';
import { CharacterSpellForgeV2 } from '@/components/CharacterSpellForgeV2';
import { RaIcon } from '@/components/ui/RaIcon';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Label } from '@/components/ui/label';
import { Card, CardContent } from '@/components/ui/card';
import { X, Bookmark, Save } from 'lucide-react';
import {
  CHARACTER_SPELLBOOK_PAGE_SIZE,
  CharacterSpellbookScopeFilter,
  EMPTY_BY_SCOPE,
  SCOPE_LABELS,
  useSpellbookListState,
} from '@/components/spellbook';
import { spellRequiresConcentration } from '@/lib/spellMechanics';

interface CharacterSpellbookProps {
  /** Class + race spell pool (from sheet or merged). Omit only if unknown; forge treats undefined as unrestricted. */
  availableComponents?: ApiComponent[];
  userId?: string;
  characterId?: string;
  token?: string;
  timerDuration?: number;
  components: ApiCharacterComponent[];
  sheet?: ApiCharacterSheet;
  isPowderMage?: boolean;
  speedDialSlots?: number;
  /** Controlled Spells / Spell Forge sub-tab (Spellbook page only). */
  spellbookSubTab?: 'spells' | 'forge' | 'forge2';
  onSpellbookSubTabChange?: (tab: 'spells' | 'forge' | 'forge2') => void;
}

export function CharacterSpellbook({
  availableComponents,
  userId,
  characterId,
  token,
  timerDuration,
  components,
  sheet,
  isPowderMage = false,
  speedDialSlots = 0,
  spellbookSubTab: spellbookSubTabProp,
  onSpellbookSubTabChange,
}: CharacterSpellbookProps) {
  const [internalSubTab, setInternalSubTab] = useState<'spells' | 'forge' | 'forge2'>('spells');
  const activeTab =
    spellbookSubTabProp !== undefined && onSpellbookSubTabChange
      ? spellbookSubTabProp
      : internalSubTab;
  const setActiveTab =
    spellbookSubTabProp !== undefined && onSpellbookSubTabChange
      ? onSpellbookSubTabChange
      : setInternalSubTab;
  const [editingSpell, setEditingSpell] = useState<ApiSpell | null>(null);
  const [spellScope, setSpellScope] = useState<CharacterSpellbookScope>('mine_or_castable');

  const { pointsFilter, setPointsFilter, currentPage: spellPage, setCurrentPage: setSpellPage } =
    useSpellbookListState();

  const queryClient = useQueryClient();

  const castMutation = useMutation({
    mutationFn: (spell: ApiSpell) =>
      castSpell(characterId!, getSpellComponentCount(spell), token, spell.id),
    onSuccess: (data: CastSpellResponse, spell: ApiSpell) => {
      const isMadnessWarning =
        sheet?.class?.name === 'The Mutagen' &&
        data.madness_cast_count !== undefined &&
        data.madness_cast_count > (sheet?.class_level?.proficiency_bonus || 2) +
          Math.floor(((sheet?.character?.constitution ?? 10) - 10) / 2);
      const message = buildCastToast(sheet, spell.name, data);
      if (isMadnessWarning) {
        toast.warning(message);
      } else {
        toast.success(message);
      }
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      queryClient.invalidateQueries({ queryKey: ['characterSpellbook', characterId] });
      queryClient.invalidateQueries({ queryKey: ['speed-dial', characterId] });
    },
    onError: (err: Error) => {
      toast.error(err.message || 'Failed to cast spell');
    },
  });

  const { data: allComponents = [] } = useQuery({
    queryKey: ['components'],
    queryFn: getComponents,
    staleTime: 60_000 * 10,
  });

  const { data: spellbookPage, isLoading: isLoadingSpells, error: spellsError } = useQuery({
    queryKey: ['characterSpellbook', characterId, spellScope, spellPage, pointsFilter],
    queryFn: () =>
      getCharacterSpellbook(characterId!, token, {
        scope: spellScope,
        page: spellPage,
        limit: CHARACTER_SPELLBOOK_PAGE_SIZE,
        level: pointsFilter === 'all' ? undefined : pointsFilter,
      }),
    enabled: !!characterId && !!token && activeTab === 'spells',
    placeholderData: (previousData) => previousData,
  });

  const { data: speedDialEntries = [] } = useQuery({
    queryKey: ['speed-dial', characterId],
    queryFn: () => getSpeedDial(characterId!, token),
    enabled: !!characterId && !!token && isPowderMage,
    staleTime: 30_000,
  });

  const saveSpeedDialMutation = useMutation({
    mutationFn: ({ slotIndex, name, componentIds }: { slotIndex: number; name: string; componentIds: string[] }) =>
      saveSpeedDialSlot(characterId!, slotIndex, { name, component_ids: componentIds }, token),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['speed-dial', characterId] });
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      toast.success('Speed Dial saved');
    },
    onError: (e: Error) => toast.error(e.message || 'Failed to save Speed Dial'),
  });

  const clearSpeedDialMutation = useMutation({
    mutationFn: (slotIndex: number) => clearSpeedDialSlot(characterId!, slotIndex, token),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['speed-dial', characterId] });
      toast.success('Speed Dial slot cleared');
    },
    onError: (e: Error) => toast.error(e.message || 'Failed to clear slot'),
  });

  const handleSaveToSpeedDial = useCallback((spell: ApiSpell, slotIndex: number) => {
    if (!characterId || !token) return;
    const componentIds = spell.components?.map(c => c.id) || [];
    if (componentIds.length === 0) {
      toast.error('This spell has no components');
      return;
    }
    if (componentIds.length > 3) {
      toast.error('Speed Dial holds at most 3 components');
      return;
    }
    saveSpeedDialMutation.mutate({
      slotIndex,
      name: spell.name,
      componentIds,
    });
  }, [characterId, token, saveSpeedDialMutation]);

  const pageSpells = spellbookPage?.spells ?? [];
  const totalPages = Math.max(
    1,
    Math.ceil((spellbookPage?.total_count ?? 0) / CHARACTER_SPELLBOOK_PAGE_SIZE),
  );

  const renderSpeedDialManagement = () => {
    if (!isPowderMage || speedDialSlots <= 0) return null;

    return (
      <Card className="arcane-border bg-card/50 mb-6">
        <CardContent className="p-4 space-y-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Bookmark className="h-4 w-4 text-primary" />
              <Label className="text-sm font-tome-heading text-primary">Speed Dial Slots</Label>
            </div>
            <Badge variant="outline" className="text-tiny font-normal">
              Max 3 components
            </Badge>
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            {Array.from({ length: speedDialSlots }).map((_, i) => {
              const entry = speedDialEntries.find((e: ApiSavedSpell) => e.slot_index === i);
              return (
                <div key={i} className="flex flex-col gap-2 p-3 rounded-md border border-border bg-background/50 relative group">
                  <div className="flex items-start justify-between gap-2">
                    <span className="text-xs font-medium truncate flex-1">
                      {entry ? entry.name : `Slot ${i + 1} (Empty)`}
                    </span>
                    {entry && (
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-6 w-6 text-muted-foreground hover:text-destructive -mr-1 -mt-1"
                        onClick={() => clearSpeedDialMutation.mutate(i)}
                        title="Clear Slot"
                        disabled={clearSpeedDialMutation.isPending}
                      >
                        <X className="h-3 w-3" />
                      </Button>
                    )}
                  </div>
                  {entry && entry.component_ids && (
                    <div className="flex gap-1 mt-auto flex-wrap">
                      {entry.component_ids.map((id: string, idx: number) => {
                        const comp = allComponents.find(c => c.id === id);
                        return (
                          <Badge key={`${id}-${idx}`} variant="secondary" className="text-[9px] px-1 h-3.5 flex gap-0.5 items-center">
                            {comp?.symbol}
                          </Badge>
                        );
                      })}
                    </div>
                  )}
                  {!entry && (
                    <div className="text-[10px] text-muted-foreground italic mt-auto">
                      Assign from spell card below
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </CardContent>
      </Card>
    );
  };

  const renderSpellsTab = () => (
    <div className="space-y-4">
      {renderSpeedDialManagement()}
      <CharacterSpellbookScopeFilter
        value={spellScope}
        onChange={(scope) => {
          setSpellScope(scope);
          setSpellPage(1);
        }}
      />
      <SpellListPagination
        initialSpells={pageSpells}
        emptyMessage={EMPTY_BY_SCOPE[spellScope]}
        renderSpellItem={(spell: ApiSpell) => (
          <PreparedSpellCard
            key={spell.id}
            spell={spell}
            onEdit={() => setEditingSpell(spell)}
            sheet={sheet}
            speedDialSpells={isPowderMage ? speedDialEntries : undefined}
            onCast={() => castMutation.mutate(spell)}
            isCasting={castMutation.isPending && castMutation.variables?.id === spell.id}
            isDisabled={!hasAllComponents(spell, availableComponents, components)}
            isPowderMage={isPowderMage}
            speedDialSlots={speedDialSlots}
            onSaveToSpeedDial={(slotIndex) => handleSaveToSpeedDial(spell, slotIndex)}
          />
        )}
        totalPages={totalPages}
        currentPage={spellPage}
        setCurrentPage={setSpellPage}
        isLoading={isLoadingSpells && !spellbookPage}
        error={spellsError}
        headerTitle={SCOPE_LABELS[spellScope]}
        pointsFilter={pointsFilter}
        setPointsFilter={setPointsFilter}
      />
    </div>
  );

  if (editingSpell) {
    const spellForForge = {
      ...editingSpell,
      concentration: spellRequiresConcentration(editingSpell.concentration),
    };
    return (
      <CharacterSpellForge
        availableComponents={availableComponents}
        userId={userId}
        characterId={characterId}
        token={token}
        timerDuration={timerDuration}
        components={components}
        spellToEdit={spellForForge}
        onClose={() => setEditingSpell(null)}
        isPowderMage={isPowderMage}
        speedDialSlots={speedDialSlots}
      />
    );
  }

  return (
    <Tabs
      value={activeTab}
      onValueChange={(v) => setActiveTab(v as 'spells' | 'forge' | 'forge2')}
      className="w-full min-w-0"
    >
      <TabsList className="grid h-auto w-full min-w-0 grid-cols-3 gap-0">
        <TabsTrigger value="spells" className="gap-1 sm:gap-2 px-1.5 sm:px-3 py-2.5 min-w-0 text-xs sm:text-sm">
          <RaIcon name="book" className="text-sm shrink-0" />
          <span className="truncate leading-tight text-center">Spells</span>
        </TabsTrigger>
        <TabsTrigger
          value="forge"
          className="gap-1 sm:gap-2 px-1.5 sm:px-3 py-2.5 min-w-0 text-xs sm:text-sm"
          disabled={!userId}
        >
          <RaIcon name="crystal-wand" className="text-sm shrink-0" />
          <span className="truncate leading-tight text-center">
            <span className="sm:hidden">Forge</span>
            <span className="hidden sm:inline">Spell Forge</span>
          </span>
        </TabsTrigger>
        <TabsTrigger
          value="forge2"
          className="gap-1 sm:gap-2 px-1.5 sm:px-3 py-2.5 min-w-0 text-xs sm:text-sm"
          disabled={!userId}
        >
          <RaIcon name="crystal-wand" className="text-sm shrink-0" />
          <span className="truncate leading-tight text-center">
            <span className="sm:hidden">Forge 2</span>
            <span className="hidden sm:inline">Spell Forge 2</span>
          </span>
        </TabsTrigger>
      </TabsList>
      <TabsContent value="spells" className="mt-4">
        {renderSpellsTab()}
      </TabsContent>
      <TabsContent value="forge" className="mt-4">
        <CharacterSpellForge
          availableComponents={availableComponents}
          userId={userId}
          characterId={characterId}
          token={token}
          timerDuration={timerDuration}
          components={components}
          isPowderMage={isPowderMage}
          speedDialSlots={speedDialSlots}
        />
      </TabsContent>
      <TabsContent value="forge2" className="mt-4">
        <CharacterSpellForgeV2
          availableComponents={availableComponents}
          userId={userId}
          characterId={characterId}
          token={token}
          timerDuration={timerDuration}
          components={components}
          currentStability={sheet?.class_resources?.find((r) => r.key === 'max_stability')?.current_value}
          maxStability={
            sheet?.class_resources?.find((r) => r.key === 'max_stability')?.max_value ??
            sheet?.class_resources?.find((r) => r.key === 'max_stability')?.value
          }
          maxBlueprintSlots={
            sheet?.class_resources?.find((r) => r.key === 'speed_dial_slots')?.max_value ??
            sheet?.class_resources?.find((r) => r.key === 'speed_dial_slots')?.value
          }
          isPowderMage={isPowderMage}
          speedDialSlots={speedDialSlots}
        />
      </TabsContent>
    </Tabs>
  );
}
