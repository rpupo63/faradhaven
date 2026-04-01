import { useState, useRef, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useGame } from '@/context/GameContext';
import { useAuth } from '@/context/AuthContext';
import { CharacterSheetView } from '@/components/CharacterSheetView';
import { CharacterSpellbook } from '@/components/CharacterSpellbook';
import { BackstoryEditor } from '@/components/BackstoryEditor';
import { LoadingQuill } from '@/components/LoadingQuill';
import {
  computeCharacterSheetFromApi,
  normalizeComputedSheet,
  normalizeApiSheet,
  getSkillFocusFromClass,
  getSavingThrowsFromClass,
} from '@/lib/classLevelData';
import {
  getCharacterSheet,
  getClassWithLevels,
  updateCharacter as updateCharacterApi,
  levelDownCharacter,
  updateHP,
  setTempHP,
  spendHitDice,
  shortRest,
  longRest,
  uploadCharacterImage,
} from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogTitle } from '@/components/ui/dialog';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { ArrowLeft, ChevronUp, ChevronDown, History, FileText, Wand2, BookOpen, Sparkles, Camera, User, FlaskRound, Coins, Type, Users, Recycle, Droplets } from 'lucide-react';
import { useQuery, useQueryClient, useMutation } from '@tanstack/react-query';
import { LevelUpWizard } from '@/components/LevelUpWizard';
import { LevelHistory } from '@/components/LevelHistory';
import { LevelDownConfirmation } from '@/components/LevelDownConfirmation';
import { HarvestModal } from '@/components/HarvestModal';
import { LootModal } from '@/components/LootModal';
import { ScavengeModal } from '@/components/ScavengeModal';
import { PartyBestiary } from '@/components/PartyBestiary'; // NEW: Import PartyBestiary
import { PartyMembers } from '@/components/PartyMembers';   // NEW: Import PartyMembers
import { dispatchClearDice } from '@/lib/dice';
type CharacterTab = 'sheet' | 'spellbook' | 'backstory' | 'bestiary' | 'party'; // NEW: Added bestiary and party

export default function CharacterSheetPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { state, updateCharacter, setActiveCharacter } = useGame();

  // Tab state
  const [activeTab, setActiveTab] = useState<CharacterTab>('sheet');

  // Level-up/down dialog state
  const [showLevelUpWizard, setShowLevelUpWizard] = useState(false);
  const [showLevelHistory, setShowLevelHistory] = useState(false);
  const [showLevelDownConfirm, setShowLevelDownConfirm] = useState(false);
  const [showHarvestModal, setShowHarvestModal] = useState(false); // New state for HarvestModal
  const [showLootModal, setShowLootModal] = useState(false);
  const [showScavengeModal, setShowScavengeModal] = useState(false);

  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isUploading, setIsUploading] = useState(false);

  const { token, setDicePrefsOverride } = useAuth();

  const character = state.characters.find((c) => c.id === id);

  // Ensure this character is set as active when viewing their sheet
  // This allows Spellbook and SpellForge to work correctly via useGame
  useEffect(() => {
    if (character && state.activeCharacterId !== character.id) {
      setActiveCharacter(character.id);
    }
  }, [character, state.activeCharacterId, setActiveCharacter]);

  const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !id || !token) return;

    try {
      setIsUploading(true);
      await uploadCharacterImage(id, file, token);
      queryClient.invalidateQueries({ queryKey: ['character-sheet', id] });
    } catch (error) {
      console.error('Failed to upload image:', error);
      alert('Failed to upload image. Please try again.');
    } finally {
      setIsUploading(false);
    }
  };

  // Try API first if we have a token
  const { data: apiSheet, isLoading: isLoadingSheet } = useQuery({
    queryKey: ['character-sheet', id],
    queryFn: () => getCharacterSheet(id!, token ?? undefined),
    enabled: !!id && !!token,
    staleTime: 30_000,
    retry: false, // Don't retry on 404/401 — fall back to local
  });

  // For local character with class_id: fetch class data to compute sheet
  const { data: classWithLevels, isLoading: isLoadingClass } = useQuery({
    queryKey: ['class', character?.class_id],
    queryFn: () => getClassWithLevels(character!.class_id!, token ?? undefined),
    enabled: !!character?.class_id && !apiSheet,
    staleTime: 60_000,
    retry: false,
  });

  const sheet = (() => {
    if (apiSheet) return normalizeApiSheet(apiSheet);
    if (character && classWithLevels) {
      const levelData = classWithLevels.levels?.find((l) => l.level === Math.max(1, character.level));
      if (!levelData) return null;
      const computed = computeCharacterSheetFromApi(
        {
          ...character,
          race: character.race,
          class: character.class,
          strength: character.strength ?? 10,
          dexterity: character.dexterity ?? 10,
          constitution: character.constitution ?? 10,
          intelligence: character.intelligence ?? 10,
          wisdom: character.wisdom ?? 10,
          charisma: character.charisma ?? 10,
          current_spell_points: character.current_spell_points,
        },
        classWithLevels,
        levelData
      );
      return normalizeComputedSheet(
        computed,
        getSkillFocusFromClass(classWithLevels),
        getSavingThrowsFromClass(classWithLevels)
      );
    }
    if (character && !character.class_id) {
      // Legacy character without class_id - cannot compute from backend
      return null;
    }
    return null;
  })();

  // Apply this character's dice overrides (if any) while the sheet is mounted.
  useEffect(() => {
    if (!apiSheet?.character) return;
    const { dice_theme, dice_theme_color, dice_font_color } = apiSheet.character;
    if (dice_theme || dice_theme_color || dice_font_color) {
      setDicePrefsOverride({ dice_theme, dice_theme_color, dice_font_color });
    }
    return () => setDicePrefsOverride(null);
  }, [
    apiSheet?.character?.dice_theme,
    apiSheet?.character?.dice_theme_color,
    apiSheet?.character?.dice_font_color,
    setDicePrefsOverride,
  ]);

  // HP Management handlers (only for API-backed characters)
  const handleHPChange = async (delta: number, source?: string) => {
    if (!apiSheet || !token) return;
    await updateHP(id!, delta, source, token);
    queryClient.invalidateQueries({ queryKey: ['character-sheet', id] });
  };

  const handleTempHPChange = async (value: number) => {
    if (!apiSheet || !token) return;
    await setTempHP(id!, value, token);
    queryClient.invalidateQueries({ queryKey: ['character-sheet', id] });
  };

  const handleUseHitDice = async (rolls: number[]) => {
    if (!apiSheet || !token) return;
    await spendHitDice(id!, rolls, token);
    queryClient.invalidateQueries({ queryKey: ['character-sheet', id] });
  };

  const handleShortRest = async () => {
    if (!apiSheet || !token) return;
    await shortRest(id!, token);
    queryClient.invalidateQueries({ queryKey: ['character-sheet', id] });
  };

  const handleLongRest = async () => {
    if (!apiSheet || !token) return;
    await longRest(id!, token);
    queryClient.invalidateQueries({ queryKey: ['character-sheet', id] });
  };

  const handleMoneyChange = async (newTotal: number) => {
    if (!apiSheet || !token) return;
    await updateCharacterApi(id!, { money: newTotal }, token);
    queryClient.invalidateQueries({ queryKey: ['character-sheet', id] });
  };

  const handleNotesChange = async (notes: string) => {
    if (!apiSheet || !token) return;
    await updateCharacterApi(id!, { notes }, token);
    queryClient.invalidateQueries({ queryKey: ['character-sheet', id] });
  };

  // Level-down mutation
  const levelDownMutation = useMutation({
    mutationFn: () => levelDownCharacter(id!, token ?? undefined),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['character-sheet', id] });
      queryClient.invalidateQueries({ queryKey: ['levelHistory', id] });
      setShowLevelDownConfirm(false);
    },
  });

  const handleLevelUpComplete = () => {
    setShowLevelUpWizard(false);
    queryClient.invalidateQueries({ queryKey: ['character-sheet', id] });
  };

  // Show loading state while fetching character data
  const isLoading = isLoadingSheet || isLoadingClass;

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[50vh]">
        <LoadingQuill label="Loading character..." />
      </div>
    );
  }

  if ((!character && !apiSheet) || !sheet) {
    return (
      <div className="space-y-4">
        <Button variant="ghost" onClick={() => navigate('/')} className="gap-2">
          <ArrowLeft className="h-4 w-4" />
          Back to Characters
        </Button>
        <p className="text-muted-foreground">
          {character && !character.class_id
            ? 'Character sheet requires class data from the backend. Edit the character to assign a class.'
            : 'Character not found.'}
        </p>
      </div>
    );
  }

  const currentLevel = apiSheet?.character?.level ?? character?.level ?? 1;

  return (
    <div className="space-y-6">
      {/* Unified Header Row */}
      <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as CharacterTab)} className="w-full">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-4">

          {/* Top Row (Mobile): Back Button + Short Info */}
          <div className="flex items-center gap-4 w-full md:w-auto">
            <Button variant="ghost" onClick={() => navigate('/')} className="gap-2 shrink-0">
              <ArrowLeft className="h-4 w-4" />
              <span className="hidden sm:inline">Back</span>
            </Button>

            {/* Mobile Character Summary + Avatar */}
            <div className="flex items-center gap-3 md:hidden overflow-hidden">
              <div
                className="relative group cursor-pointer shrink-0"
                onClick={() => fileInputRef.current?.click()}
              >
                <div className="w-10 h-10 rounded-full border border-primary/20 bg-primary/5 flex items-center justify-center overflow-hidden">
                  {sheet.character.image_url ? (
                    <img src={sheet.character.image_url} alt={sheet.character.name} className="w-full h-full object-cover" />
                  ) : (
                    <User className="w-6 h-6 text-primary/40" />
                  )}
                </div>
                {isUploading && (
                  <div className="absolute inset-0 bg-background/60 flex items-center justify-center rounded-full">
                    <div className="w-4 h-4 border-2 border-primary border-t-transparent animate-spin rounded-full" />
                  </div>
                )}
              </div>
              <div className="min-w-0">
                <h2 className="font-display text-lg text-primary glow-text leading-tight truncate">
                  {sheet.character.name}
                </h2>
                <p className="text-xs text-muted-foreground font-tome-marginalia">
                  Lvl {sheet.character.level} {sheet.character.className}
                </p>
              </div>
            </div>
          </div>

          <input
            type="file"
            ref={fileInputRef}
            className="hidden"
            accept="image/*"
            onChange={handleImageUpload}
          />

          {/* Tab Navigation: Bottom on mobile, Center on Desktop */}
          <div className="order-last md:order-none w-full md:w-auto min-w-0">
            <TabsList className="grid grid-cols-5 w-full md:w-auto min-w-0 h-12 md:h-10 gap-0">
              <TabsTrigger value="sheet" className="gap-1 sm:gap-2 px-1.5 sm:px-3 h-full min-w-0">
                <FileText className="h-4 w-4 shrink-0" />
                <span className="hidden sm:inline truncate">Sheet</span>
              </TabsTrigger>

              <TabsTrigger value="spellbook" className="gap-1 sm:gap-2 px-1.5 sm:px-3 h-full min-w-0">
                <Wand2 className="h-4 w-4 shrink-0" />
                <span className="hidden sm:inline truncate">Spellbook</span>
              </TabsTrigger>
              <TabsTrigger value="backstory" className="gap-1 sm:gap-2 px-1.5 sm:px-3 h-full min-w-0">
                <History className="h-4 w-4 shrink-0" />
                <span className="hidden sm:inline truncate">Story</span>
              </TabsTrigger>
              <TabsTrigger value="bestiary" className="gap-1 sm:gap-2 px-1.5 sm:px-3 h-full min-w-0">
                <Sparkles className="h-4 w-4 shrink-0" />
                <span className="hidden sm:inline truncate">Bestiary</span>
              </TabsTrigger>
              <TabsTrigger value="party" className="gap-1 sm:gap-2 px-1.5 sm:px-3 h-full min-w-0">
                <Users className="h-4 w-4 shrink-0" />
                <span className="hidden sm:inline truncate">Party</span>
              </TabsTrigger>
            </TabsList>
          </div>

          {/* Right: Level Controls + Desktop Info */}
          <div className="flex flex-col md:flex-row items-end md:items-center gap-2 md:gap-4 w-full md:w-auto min-w-0">
            {/* Level Controls */}
            <div className="flex flex-col items-stretch gap-1.5 w-full md:w-auto min-w-0 md:items-end">
              <span className="text-xs text-muted-foreground md:hidden uppercase tracking-wider font-semibold text-right">
                Actions
              </span>
              <div className="flex flex-wrap items-center justify-end gap-1.5 sm:gap-2 min-w-0">
              {sheet.character.className === "The Lorewright" && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setShowHarvestModal(true)}
                  title="Perform Harvest (Visceral Psychometry)"
                  className="gap-1 h-8 px-2 shrink-0"
                >
                  <FlaskRound className="h-4 w-4" />
                  <span className="hidden lg:inline">Harvest</span>
                </Button>
              )}
              {(sheet.character.className === "The Lorewright" ||
                sheet.character.className === "The Ironwright") && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setShowScavengeModal(true)}
                  title="Scavenge Components"
                  className="gap-1 h-8 px-2 shrink-0"
                >
                  <Recycle className="h-4 w-4" />
                  <span className="hidden lg:inline">Scavenge</span>
                </Button>
              )}
              {sheet.character.className === "The Sanguinist" && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setShowScavengeModal(true)}
                  title="Sanguine Extraction"
                  className="gap-1 h-8 px-2 shrink-0"
                >
                  <Droplets className="h-4 w-4" />
                  <span className="hidden lg:inline">Extract</span>
                </Button>
              )}
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowLootModal(true)}
                title="Generate Loot"
                className="gap-1 h-8 px-2 shrink-0"
              >
                <Coins className="h-4 w-4" />
                <span className="hidden lg:inline">Loot</span>
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => {
                  if (!apiSheet || !token) {
                    alert('Please log in to level your character');
                    return;
                  }
                  if (currentLevel <= 1) {
                    alert('Character is already at level 1');
                    return;
                  }
                  setShowLevelDownConfirm(true);
                }}
                disabled={currentLevel <= 1}
                title="Level Down"
                className="gap-1 h-8 px-2 shrink-0"
              >
                <ChevronDown className="h-4 w-4" />
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => {
                  if (!apiSheet || !token) {
                    alert('Please log in to level your character');
                    return;
                  }
                  if (currentLevel >= 20) {
                    alert('Character is already at max level (20)');
                    return;
                  }
                  setShowLevelUpWizard(true);
                }}
                disabled={currentLevel >= 20}
                title="Level Up"
                className="gap-1 h-8 px-2 shrink-0"
              >
                <ChevronUp className="h-4 w-4" />
              </Button>

              </div>
            </div>

            {/* Desktop Character Info */}
            <div className="flex items-center gap-4 md:flex">
              <div className="text-right">
                <h2 className="font-display text-xl text-primary glow-text">
                  {sheet.character.name}
                </h2>
                <p className="text-sm text-muted-foreground font-tome-marginalia">
                  Level {sheet.character.level} {sheet.character.raceName} {sheet.character.className}
                  {sheet.character.archetypeName && ` - ${sheet.character.archetypeName}`}
                </p>
              </div>

              <div
                className="relative group cursor-pointer"
                onClick={() => fileInputRef.current?.click()}
              >
                <div className="w-14 h-14 rounded-lg arcane-border overflow-hidden bg-primary/5 flex items-center justify-center transition-all
  group-hover:ring-2 group-hover:ring-primary/50">
                  {sheet.character.image_url ? (
                    <img src={sheet.character.image_url} alt={sheet.character.name} className="w-full h-full object-cover" />
                  ) : (
                    <User className="w-8 h-8 text-primary/40" />
                  )}
                  <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center
  justify-center">
                    <Camera className="w-6 h-6 text-white" />
                  </div>
                </div>
                {isUploading && (
                  <div className="absolute inset-0 bg-background/60 flex items-center justify-center rounded-lg">
                    <div className="w-6 h-6 border-2 border-primary border-t-transparent animate-spin rounded-full" />
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>

        <TabsContent value="sheet" className="mt-6">
          <CharacterSheetView
            sheet={sheet}
            onShortRest={apiSheet && token ? handleShortRest : undefined}
            onLongRest={apiSheet && token ? handleLongRest : undefined}
            onHPChange={apiSheet && token ? handleHPChange : undefined}
            onTempHPChange={apiSheet && token ? handleTempHPChange : undefined}
            onUseHitDice={apiSheet && token ? handleUseHitDice : undefined}
            onMoneyChange={apiSheet && token ? handleMoneyChange : undefined}
            onNotesChange={apiSheet && token ? handleNotesChange : undefined}
            hideHeader
          />
        </TabsContent>



        <TabsContent value="spellbook" className="mt-6">
          <CharacterSpellbook
            availableComponents={sheet.available_components || []}
            userId={apiSheet?.character?.user_id}
            characterId={id}
            token={token ?? undefined}
            timerDuration={sheet.class_resources?.find(r => r.key === 'timer_duration')?.value}
            components={sheet.components || []}
            sheet={apiSheet}
          />
        </TabsContent>

        <TabsContent value="backstory" className="mt-6">
          <BackstoryEditor
            characterId={id!}
            token={token ?? undefined}
            characterName={sheet.character.name}
          />
        </TabsContent>

        {/* NEW: Bestiary Tab */}
        <TabsContent value="bestiary" className="mt-6">
          {id && token && <PartyBestiary characterId={id} token={token} />}
        </TabsContent>

        {/* NEW: Party Tab */}
        <TabsContent value="party" className="mt-6">
          {id && token && <PartyMembers characterId={id} token={token} characterSheet={apiSheet} />}
        </TabsContent>
      </Tabs>

      {/* Level-up Wizard Dialog */}
      <Dialog 
        open={showLevelUpWizard} 
        onOpenChange={(open) => {
          setShowLevelUpWizard(open);
          if (!open) dispatchClearDice();
        }}
      >
        <DialogContent className="max-w-2xl p-0 border-0 bg-transparent" noPadding>
          <DialogTitle className="sr-only">Level Up Wizard</DialogTitle>
          <DialogDescription className="sr-only">Interface for leveling up your character and choosing new features.</DialogDescription>
          <LevelUpWizard
            characterId={id!}
            token={token ?? undefined}
            onComplete={handleLevelUpComplete}
            onCancel={() => {
              setShowLevelUpWizard(false);
              dispatchClearDice();
            }}
          />
        </DialogContent>
      </Dialog>



      {/* Level-down Confirmation */}
      <LevelDownConfirmation
        open={showLevelDownConfirm}
        onOpenChange={(open) => {
          setShowLevelDownConfirm(open);
          if (!open) dispatchClearDice();
        }}
        currentLevel={currentLevel}
        onConfirm={() => levelDownMutation.mutate()}
        isPending={levelDownMutation.isPending}
      />

      {/* Harvest Modal */}
      {id && token && (
        <HarvestModal
          isOpen={showHarvestModal}
          onClose={() => {
            setShowHarvestModal(false);
            dispatchClearDice();
          }}
          characterId={id}
          token={token}
          onHarvestSuccess={() => queryClient.invalidateQueries({ queryKey: ['character-sheet', id] })}
        />
      )}

      {/* Loot Modal */}
      {id && token && (
        <LootModal
          isOpen={showLootModal}
          onClose={() => {
            setShowLootModal(false);
            dispatchClearDice();
          }}
          characterId={id}
          token={token}
        />
      )}

      {/* Scavenge / Extract Modal */}
      {id && token && sheet && (
        <ScavengeModal
          isOpen={showScavengeModal}
          onClose={() => {
            setShowScavengeModal(false);
            dispatchClearDice();
          }}
          characterId={id}
          token={token}
          className={sheet.character.className}
          yieldDie={sheet.class_resources?.find(r => r.key === 'yield_die')?.value}
          characterLevel={currentLevel}
          onSuccess={() => queryClient.invalidateQueries({ queryKey: ['character-sheet', id] })}
        />
      )}
    </div>
  );
}