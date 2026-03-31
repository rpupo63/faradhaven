
import { useState, useMemo } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { getCharacterSpellbook, castSpell } from '@/lib/api';
import { ApiComponent, ApiSpell, ApiCharacterComponent, CastSpellResponse, ApiCharacterSheet } from '@/types/game';
import { hasAllComponents } from '@/lib/spellUtils';
import { toast } from 'sonner';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { SpellListPagination } from '@/components/SpellListPagination';
import { PreparedSpellCard } from '@/components/PreparedSpells';
import { buildCastToast } from '@/lib/toastUtils';
import { CharacterSpellForge } from '@/components/CharacterSpellForge';
import { LoadingQuill } from './LoadingQuill';
import { Wand2, User, Sparkles } from 'lucide-react';


interface CharacterSpellbookProps {
  availableComponents: ApiComponent[];
  userId?: string;
  characterId?: string;
  token?: string;
  timerDuration?: number;
  components: ApiCharacterComponent[];
  sheet?: ApiCharacterSheet;
}

export function CharacterSpellbook({
  availableComponents,
  userId,
  characterId,
  token,
  timerDuration,
  components,
  sheet,
}: CharacterSpellbookProps) {
  const [activeTab, setActiveTab] = useState('available');
  const [editingSpell, setEditingSpell] = useState<ApiSpell | null>(null);

  const queryClient = useQueryClient();

  const castMutation = useMutation({
    mutationFn: (spell: ApiSpell) => castSpell(characterId!, spell.slot_level, token, spell.id),
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
    },
    onError: (err: Error) => {
      toast.error(err.message || 'Failed to cast spell');
    },
  });

  const { data: spellbook, isLoading: isLoadingSpells, error: spellsError } = useQuery({
    queryKey: ['characterSpellbook', characterId],
    queryFn: () => getCharacterSpellbook(characterId!, token),
    enabled: !!characterId && !!token,
  });

  const [pointsFilter, setPointsFilter] = useState<number | 'all'>('all');
  const [currentPage, setCurrentPage] = useState(1);

  const allAvailableSpells = useMemo(() => spellbook?.available_spells ?? [], [spellbook]);
  const myAvailableSpells = useMemo(() => spellbook?.my_spells ?? [], [spellbook]);

  const renderSpellList = (spells: ApiSpell[], emptyMessage: string, headerTitle: string, enableCast = false) => {
    if (isLoadingSpells) {
      return <LoadingQuill label="Loading spells..." />;
    }
    if (spellsError) {
      return <p className="text-destructive">Failed to load spells.</p>;
    }
    return (
      <SpellListPagination
        initialSpells={spells}
        emptyMessage={emptyMessage}
        renderSpellItem={(spell: ApiSpell) => (
          <PreparedSpellCard
            key={spell.id}
            spell={spell}
            onEdit={() => setEditingSpell(spell)}
            sheet={enableCast ? sheet : undefined}
            onCast={enableCast ? () => castMutation.mutate(spell) : undefined}
            isCasting={enableCast && castMutation.isPending && castMutation.variables?.id === spell.id}
            isDisabled={!hasAllComponents(spell, availableComponents, components)}
          />
        )}
        totalPages={1}
        currentPage={currentPage}
        setCurrentPage={setCurrentPage}
        isLoading={isLoadingSpells}
        error={spellsError}
        headerTitle={headerTitle}
        pointsFilter={pointsFilter}
        setPointsFilter={setPointsFilter}
      />
    );
  };

  if (editingSpell) {
    return (
      <CharacterSpellForge
        availableComponents={availableComponents}
        userId={userId}
        characterId={characterId}
        token={token}
        timerDuration={timerDuration}
        components={components}
        spellToEdit={editingSpell}
        onClose={() => setEditingSpell(null)}
      />
    )
  }

  return (
    <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
      <TabsList className="grid w-full grid-cols-3">
        <TabsTrigger value="available" className="gap-2">
          <Sparkles className="h-4 w-4" />
          Available Spells
        </TabsTrigger>
        <TabsTrigger value="my-spells" className="gap-2" disabled={!userId}>
          <User className="h-4 w-4" />
          My Spells
        </TabsTrigger>
        <TabsTrigger value="forge" className="gap-2" disabled={!userId}>
          <Wand2 className="h-4 w-4" />
          Spell Forge
        </TabsTrigger>
      </TabsList>
      <TabsContent value="available" className="mt-4">
        {renderSpellList(allAvailableSpells, "No spells are available with your current components.", "Available Spells", true)}
      </TabsContent>
      <TabsContent value="my-spells" className="mt-4">
        {renderSpellList(myAvailableSpells, "You have not created any spells yet.", "My Spells", true)}
      </TabsContent>
      <TabsContent value="forge" className="mt-4">
        <CharacterSpellForge
          availableComponents={availableComponents}
          userId={userId}
          characterId={characterId}
          token={token}
          timerDuration={timerDuration}
          components={components}
        />
      </TabsContent>
    </Tabs>
  );
}
