import { useState, useEffect, useRef } from 'react';
import { cn } from '@/lib/utils';
import { type NormalizedCharacterSheet, type ApiCharacterWeapon } from '@/types/game';
import { ChevronUp, ChevronDown, Sun, Moon, Skull } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { LootModal } from '@/components/LootModal';
import { rollD20, rollD, dispatchClearDice } from '@/lib/dice';
import { getActiveEffects } from '@/lib/api/mechanics';
import { getWeaponById } from '@/lib/api/character';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/context/AuthContext';

// New sub-components
import { AbilityScores } from './character-sheet/AbilityScores';
import { ProficienciesSection } from './character-sheet/ProficienciesSection';
import { CombatStats } from './character-sheet/CombatStats';
import { HPPanel } from './character-sheet/HPPanel';
import { HitDicePanel } from './character-sheet/HitDicePanel';
import { ClassResourceDisplay } from './character-sheet/ClassResourceDisplay';
import { ComponentInventorySection } from './character-sheet/ComponentInventorySection';
import { RacialResourceTracker } from './character-sheet/RacialResourceTracker';
import { HarvestBankSection } from './character-sheet/HarvestBankSection';
import { ActiveEffectsSection } from './character-sheet/ActiveEffectsSection';
import { ConstructsSection } from './character-sheet/ConstructsSection';
import { NotorietyMeter } from './character-sheet/NotorietyMeter';
import { EquipmentSection } from './character-sheet/EquipmentSection';
import { FeaturesSection } from './character-sheet/FeaturesSection';
import { ActiveAbilitiesSection } from './character-sheet/ActiveAbilitiesSection';
import { MoneyPanel } from './character-sheet/MoneyPanel';
import { WeaponAttackDialog } from './character-sheet/WeaponAttackDialog';
import { DieOptions } from './character-sheet/DieOptions';
import { SanguinistFeaturesCard } from './SanguinistFeaturesCard';
import { LorewrightFeaturesCard } from './LorewrightFeaturesCard';
import { MutagenStatusCard } from './MutagenStatusCard';
import { MutagenFeaturesCard } from './MutagenFeaturesCard';
import { PowderMageStatusCard } from './PowderMageStatusCard';
import { PowderMageFeaturesCard } from './PowderMageFeaturesCard';
import { PistonBrawlerStatusCard } from './PistonBrawlerStatusCard';
import { PistonBrawlerFeaturesCard } from './PistonBrawlerFeaturesCard';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import { UserCircle } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { MobileSheetSummaryStrip } from './character-sheet/MobileSheetSummaryStrip';
import { MobileSheetBottomBar } from './character-sheet/MobileSheetBottomBar';

interface CharacterSheetViewProps {
  sheet: NormalizedCharacterSheet;
  onShortRest?: () => void | Promise<void>;
  onLongRest?: () => void | Promise<void>;
  onHPChange?: (delta: number, source?: string) => void | Promise<void>;
  onTempHPChange?: (value: number) => void | Promise<void>;
  onUseHitDice?: (rolls: number[]) => void | Promise<void>;
  onLevelUp?: () => void;
  onLevelDown?: () => void;
  onMoneyChange?: (newTotal: number) => void | Promise<void>;
  onNotesChange?: (notes: string) => void | Promise<void>;
  hideHeader?: boolean;
  className?: string;
}

/**
 * Character Sheet: D&D-style layout with skills on left, abilities beside them,
 * combat stats in center, and class features on the right.
 */
export function CharacterSheetView({
  sheet,
  onShortRest,
  onLongRest,
  onHPChange,
  onTempHPChange,
  onUseHitDice,
  onLevelUp,
  onLevelDown,
  onMoneyChange,
  onNotesChange,
  hideHeader = false,
  className,
}: CharacterSheetViewProps) {
  const {
    character,
  } = sheet;

  const { token } = useAuth();
  const queryClient = useQueryClient();

  // Helper functions for hand calculations (moved from EquipmentSection)
  const isTwoHanded = (properties?: string[]): boolean => {
    return !!properties?.includes('Two-Handed');
  };

  const computeUsedHands = (sheet: NormalizedCharacterSheet): number => {
    let hands = 0;
    for (const cw of sheet.inventory_weapons ?? []) {
      if (cw.is_equipped && cw.character_weapon_id !== 'virtual-bite') {
        hands += isTwoHanded(cw.weapon.properties) ? 2 : 1;
      }
    }
    if (sheet.character.equipped_shield_id) {
      hands++;
    }
    return hands;
  };

  // Calculate hands for rendering
  const usedHands = computeUsedHands(sheet);
  const freeHands = 2 - usedHands;

  // State for LootModal
  const [isLootModalOpen, setIsLootModalOpen] = useState(false);
  const [deathPopupOpen, setDeathPopupOpen] = useState(false); // Moved here

  // Fetch active effects
  const { data: activeEffects } = useQuery({
    queryKey: ['active-effects', character.id],
    queryFn: () => getActiveEffects(character.id, token || undefined),
    enabled: !!character.id && !!token,
  });

  // Track previous HP to detect when it drops to 0 or below
  const prevHPRef = useRef(sheet.current_hp);

  useEffect(() => {
    // Play death sound and show popup if HP drops to 0 or below
    if (sheet.current_hp <= 0 && prevHPRef.current > 0) {
      const audio = new Audio('/death-sound.mp3');
      audio.play().catch(err => console.error('Error playing death sound:', err));
      setTimeout(() => setDeathPopupOpen(true), 0);
    }
    prevHPRef.current = sheet.current_hp;
  }, [sheet.current_hp]);

  // Expanded panel state
  const [expandedPanel, setExpandedPanel] = useState<'hp' | 'hitdice' | null>(null);
  const [classFeaturesOpen, setClassFeaturesOpen] = useState(false);
  const [dieOptionsOpen, setDieOptionsOpen] = useState(false);

  // Weapon attack state
  const [weaponAttackDialogOpen, setWeaponAttackDialogOpen] = useState(false);
  const [selectedWeapon, setSelectedWeapon] = useState<ApiCharacterWeapon | null>(null);

  const handleRoll = async (label: string, modifier: number) => {
    await rollD20(modifier, label);
  };

  const handleGenericRoll = async (sides: number) => {
    await rollD(sides, `d${sides} Roll`);
  };

  const handleWeaponClick = async (charWeapon: ApiCharacterWeapon) => {
    let weaponToSet = charWeapon;

    if (charWeapon.weapon.transformed_to_weapon_id && !charWeapon.transformed_weapon_data) {
      if (!token) {
        console.error("No authentication token available to fetch transformed weapon data.");
        setSelectedWeapon(charWeapon); // Still open dialog with base weapon
        setWeaponAttackDialogOpen(true);
        return;
      }
      try {
        const transformedWeapon = await getWeaponById(charWeapon.weapon.transformed_to_weapon_id, token);
        weaponToSet = {
          ...charWeapon,
          transformed_weapon_data: transformedWeapon,
        };
      } catch (error) {
        console.error("Failed to fetch transformed weapon data:", error);
        // Fallback to base weapon if transformed data can't be fetched
        weaponToSet = charWeapon;
      }
    }
    setSelectedWeapon(weaponToSet);
    setWeaponAttackDialogOpen(true);
  };

  const handleEquipmentChange = () => {
    queryClient.invalidateQueries({ queryKey: ['character-sheet', character.id] });
  };

  const handleMobileHPTap = () => {
    if (!onHPChange) {
      alert('Please log in to modify HP');
      return;
    }
    setExpandedPanel(expandedPanel === 'hp' ? null : 'hp');
  };

  return (
    <div className={cn('space-y-6 relative pb-24 md:pb-0', className)}>
      {/* Header */}
      {!hideHeader && (
        <div className="border-b border-border pb-4">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="font-display text-2xl text-primary glow-text">
                {character.name}
              </h2>
              <p className="text-muted-foreground font-tome-marginalia">
                Level {character.level} {character.lineageName || character.raceName} {character.className}
                {character.archetypeName && ` - ${character.archetypeName}`}
              </p>
            </div>
            {/* Level Controls */}
            {(onLevelUp || onLevelDown) && (
              <div className="flex items-center gap-2">
                {onLevelDown && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={onLevelDown}
                    disabled={character.level <= 1}
                    className="gap-1"
                  >
                    <ChevronDown className="h-4 w-4" />
                    Level Down
                  </Button>
                )}
                {onLevelUp && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={onLevelUp}
                    disabled={character.level >= 20}
                    className="gap-1"
                  >
                    <ChevronUp className="h-4 w-4" />
                    Level Up
                  </Button>
                )}
              </div>
            )}
          </div>
        </div>
      )}

      {/* Mobile: sticky quick stats + accordions + fixed bottom bar (see MobileSheetBottomBar) */}
      <div className="md:hidden space-y-3">
        <MobileSheetSummaryStrip
          sheet={sheet}
          onHPTap={onHPChange ? handleMobileHPTap : undefined}
          onInitiativeTap={() => void handleRoll('Initiative', sheet.modifiers.initiative)}
        />
        <Accordion
          type="multiple"
          defaultValue={['combat', 'abilities']}
          className="w-full rounded-lg border border-border/60 bg-card/40 px-1 text-sm leading-snug"
        >
          <AccordionItem value="combat" className="border-border/50 px-2">
            <AccordionTrigger className="py-2 text-sm font-tome-subheading text-primary hover:no-underline">
              Combat & vitals
            </AccordionTrigger>
            <AccordionContent className="space-y-3 text-sm leading-snug pb-3">
              <CombatStats
                sheet={sheet}
                expandedPanel={expandedPanel}
                setExpandedPanel={setExpandedPanel}
                onHPChange={onHPChange}
                onUseHitDice={onUseHitDice}
                onRoll={handleRoll}
              />
              {sheet.character.raceName.includes('Changeling') && (
                <Card className="arcane-border bg-card">
                  <CardContent className="py-2 px-3 flex items-center gap-2">
                    <UserCircle className="h-4 w-4 text-primary shrink-0" />
                    <span className="text-xs font-tome-marginalia text-muted-foreground uppercase shrink-0">Form:</span>
                    <Input
                      placeholder="Current Persona..."
                      className="h-7 text-xs bg-transparent border-none focus-visible:ring-0 px-0 font-tome-subheading"
                      defaultValue={localStorage.getItem(`persona_${character.id}`) || ''}
                      onChange={(e) => localStorage.setItem(`persona_${character.id}`, e.target.value)}
                    />
                  </CardContent>
                </Card>
              )}
              {sheet.money !== undefined && (
                <MoneyPanel sheet={sheet} onMoneyChange={onMoneyChange} />
              )}
            </AccordionContent>
          </AccordionItem>

          <AccordionItem value="abilities" className="border-border/50 px-2">
            <AccordionTrigger className="py-2 text-sm font-tome-subheading text-primary hover:no-underline">
              Ability scores
            </AccordionTrigger>
            <AccordionContent className="space-y-3 pb-3">
              <AbilityScores sheet={sheet} />
            </AccordionContent>
          </AccordionItem>

          <AccordionItem value="skills" className="border-border/50 px-2">
            <AccordionTrigger className="py-2 text-sm font-tome-subheading text-primary hover:no-underline">
              Saves & skills
            </AccordionTrigger>
            <AccordionContent className="pb-3">
              <ProficienciesSection sheet={sheet} onRoll={handleRoll} />
            </AccordionContent>
          </AccordionItem>

          <AccordionItem value="resources" className="border-border/50 px-2">
            <AccordionTrigger className="py-2 text-sm font-tome-subheading text-primary hover:no-underline">
              Resources & components
            </AccordionTrigger>
            <AccordionContent className="space-y-3 pb-3">
              {sheet.class_resources && sheet.class_resources.length > 0 && token && (
                <ClassResourceDisplay resources={sheet.class_resources} characterId={character.id} token={token} />
              )}
              {token && <ComponentInventorySection sheet={sheet} token={token} />}
              {token && (
                <RacialResourceTracker sheet={sheet} characterId={character.id} token={token} />
              )}
              {sheet.class.name === 'The Lorewright' && (
                <>
                  <HarvestBankSection sheet={sheet} />
                  <Card
                    className="arcane-border bg-card cursor-pointer hover:ring-1 hover:ring-primary/50 transition-all"
                    onClick={() => setClassFeaturesOpen(true)}
                  >
                    <CardContent className="flex items-center justify-between py-3 px-4">
                      <div className="flex items-center gap-2">
                        <Skull className="h-4 w-4 text-red-400" />
                        <span className="text-sm font-tome-subheading text-primary">Trauma & Madness</span>
                      </div>
                      <span className="text-lg font-display text-red-400">{sheet.trauma ?? 0}</span>
                    </CardContent>
                  </Card>
                </>
              )}
              {sheet.class.name === 'The Ironwright' && <ConstructsSection sheet={sheet} />}
              {activeEffects && activeEffects.length > 0 && (
                <ActiveEffectsSection characterId={sheet.character.id} effects={activeEffects} />
              )}
              {sheet.class.name === 'The Sanguinist' && (
                <NotorietyMeter
                  characterId={sheet.character.id}
                  notoriety={sheet.character.sanguine_notoriety ?? 0}
                  sanguine_mp={sheet.character.sanguine_mp ?? 0}
                  sanguine_br={sheet.character.sanguine_br ?? 0}
                  onClick={() => setClassFeaturesOpen(true)}
                />
              )}
              {sheet.class.name === 'The Mutagen' && (
                <MutagenStatusCard sheet={sheet} onClick={() => setClassFeaturesOpen(true)} />
              )}
              {sheet.class.name === 'The Powder Mage' && (
                <PowderMageStatusCard sheet={sheet} onClick={() => setClassFeaturesOpen(true)} />
              )}
              {sheet.class.name === 'The Piston Brawler' && (
                <PistonBrawlerStatusCard sheet={sheet} onClick={() => setClassFeaturesOpen(true)} />
              )}
            </AccordionContent>
          </AccordionItem>

          <AccordionItem value="equipment" className="border-border/50 px-2">
            <AccordionTrigger className="py-2 text-sm font-tome-subheading text-primary hover:no-underline">
              Equipment & weapons
            </AccordionTrigger>
            <AccordionContent className="pb-3">
              <EquipmentSection
                sheet={sheet}
                onWeaponClick={handleWeaponClick}
                onGenerateLoot={() => setIsLootModalOpen(true)}
                onEquipmentChange={handleEquipmentChange}
                isTwoHanded={isTwoHanded}
                usedHands={usedHands}
              />
            </AccordionContent>
          </AccordionItem>

          <AccordionItem value="features" className="border-border/50 px-2 border-b-0">
            <AccordionTrigger className="py-2 text-sm font-tome-subheading text-primary hover:no-underline">
              Features & notes
            </AccordionTrigger>
            <AccordionContent className="space-y-3 pb-3">
              {token && (
                <ActiveAbilitiesSection sheet={sheet} characterId={character.id} token={token} />
              )}
              <FeaturesSection sheet={sheet} onNotesChange={onNotesChange} />
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      </div>

      <MobileSheetBottomBar
        onShortRest={
          onShortRest
            ? async () => {
                await onShortRest();
                setExpandedPanel('hitdice');
              }
            : undefined
        }
        onLongRest={onLongRest}
        onOpenDice={() => setDieOptionsOpen(true)}
      />

      {/* Main 4-column layout (desktop): Abilities | Skills | Middle | Features */}
      <div className="hidden md:grid grid-cols-1 gap-6 min-w-0 md:grid-cols-[90px_180px_1fr] lg:grid-cols-[100px_200px_1fr_280px]">
        
        {/* FAR LEFT COLUMN: Ability Scores */}
        <div className="order-2 md:order-none space-y-4 min-w-0">
          <AbilityScores sheet={sheet} />
          <Button
            variant="outline"
            className="hidden md:flex w-full h-auto min-h-12 whitespace-normal text-wrap"
            onClick={() => setDieOptionsOpen(true)}
          >
            Roll Any Die
          </Button>
        </div>

        {/* SECOND COLUMN: Saving Throws, Skills, Languages */}
        <div className="order-3 md:order-none space-y-4 min-w-0">
          <ProficienciesSection 
            sheet={sheet} 
            onRoll={handleRoll} 
          />
        </div>

        {/* MIDDLE COLUMN: HP, Stats, Equipment, Weapons, Languages */}
        <div className="space-y-4 order-1 md:order-none min-w-0">
          <CombatStats
            sheet={sheet}
            expandedPanel={expandedPanel}
            setExpandedPanel={setExpandedPanel}
            onHPChange={onHPChange}
            onUseHitDice={onUseHitDice}
            onRoll={handleRoll}
          />

          {/* Changeling Persona */}
          {sheet.character.raceName.includes('Changeling') && (
            <Card className="arcane-border bg-card">
              <CardContent className="py-2 px-3 flex items-center gap-2">
                <UserCircle className="h-4 w-4 text-primary shrink-0" />
                <span className="text-xs font-tome-marginalia text-muted-foreground uppercase shrink-0">Form:</span>
                <Input 
                  placeholder="Current Persona..." 
                  className="h-7 text-xs bg-transparent border-none focus-visible:ring-0 px-0 font-tome-subheading"
                  defaultValue={localStorage.getItem(`persona_${character.id}`) || ''}
                  onChange={(e) => localStorage.setItem(`persona_${character.id}`, e.target.value)}
                />
              </CardContent>
            </Card>
          )}

          {/* Money Tracking */}          {sheet.money !== undefined && (
            <MoneyPanel sheet={sheet} onMoneyChange={onMoneyChange} />
          )}

          {/* Generic Class Resources */}
          {sheet.class_resources && sheet.class_resources.length > 0 && token && (
            <ClassResourceDisplay resources={sheet.class_resources} characterId={character.id} token={token} />
          )}

          {/* Component Inventory */}
          {token && (
            <ComponentInventorySection sheet={sheet} token={token} />
          )}

          {/* Racial Resources */}
          {token && (
            <RacialResourceTracker sheet={sheet} characterId={character.id} token={token} />
          )}

          {/* Harvest Bank (Lorewright only) */}
          {sheet.class.name === 'The Lorewright' && (
            <>
              <HarvestBankSection sheet={sheet} />
              <Card 
                className="arcane-border bg-card cursor-pointer hover:ring-1 hover:ring-primary/50 transition-all"
                onClick={() => setClassFeaturesOpen(true)}
              >
                <CardContent className="flex items-center justify-between py-3 px-4">
                  <div className="flex items-center gap-2">
                    <Skull className="h-4 w-4 text-red-400" />
                    <span className="text-sm font-tome-subheading text-primary">Trauma & Madness</span>
                  </div>
                  <span className="text-lg font-display text-red-400">{sheet.trauma ?? 0}</span>
                </CardContent>
              </Card>
            </>
          )}

          {/* Constructs Section (Ironwright only) */}
          {sheet.class.name === 'The Ironwright' && (
            <ConstructsSection sheet={sheet} />
          )}

          {/* Active Effects */}
          {activeEffects && activeEffects.length > 0 && (
            <ActiveEffectsSection 
              characterId={sheet.character.id} 
              effects={activeEffects} 
            />
          )}

          {/* Notoriety Meter (Sanguinist only) - click to open moral seesaw controls */}
          {sheet.class.name === 'The Sanguinist' && (
            <NotorietyMeter
              characterId={sheet.character.id}
              notoriety={sheet.character.sanguine_notoriety ?? 0}
              sanguine_mp={sheet.character.sanguine_mp ?? 0}
              sanguine_br={sheet.character.sanguine_br ?? 0}
              onClick={() => setClassFeaturesOpen(true)}
            />
          )}

          {/* Mutagen State (Mutagen only) - click to open feral mechanics */}
          {sheet.class.name === 'The Mutagen' && (
            <MutagenStatusCard sheet={sheet} onClick={() => setClassFeaturesOpen(true)} />
          )}

          {/* Powder Mage — Casting Window */}
          {sheet.class.name === 'The Powder Mage' && (
            <PowderMageStatusCard sheet={sheet} onClick={() => setClassFeaturesOpen(true)} />
          )}

          {/* Piston Brawler — Piston Core */}
          {sheet.class.name === 'The Piston Brawler' && (
            <PistonBrawlerStatusCard sheet={sheet} onClick={() => setClassFeaturesOpen(true)} />
          )}

          {/* Detailed Equipment Section */}
          <EquipmentSection 
            sheet={sheet} 
            onWeaponClick={handleWeaponClick} 
            onGenerateLoot={() => setIsLootModalOpen(true)} 
            onEquipmentChange={handleEquipmentChange}
            isTwoHanded={isTwoHanded} // Pass isTwoHanded helper
            usedHands={usedHands}     // Pass usedHands calculation
          />
        </div>

        {/* RIGHT COLUMN: Class Features & Racial Traits */}
        <div className="order-4 md:order-none space-y-4 min-w-0">
          <div className="hidden md:grid grid-cols-2 gap-2">
            {onShortRest && (
              <Button
                variant="outline"
                onClick={async () => {
                  await onShortRest();
                  setExpandedPanel('hitdice');
                }}
                className="w-full gap-2 font-tome-marginalia h-12 px-2"
              >
                <Moon className="h-5 w-5 shrink-0" />
                <div className="text-left overflow-hidden">
                  <span className="block text-sm">Short Rest</span>
                  <span className="text-micro text-muted-foreground truncate">Resources</span>
                </div>
              </Button>
            )}
            {onLongRest && (
              <Button
                variant="outline"
                onClick={onLongRest}
                className="w-full gap-2 font-tome-marginalia h-12 px-2"
              >
                <Sun className="h-5 w-5 shrink-0" />
                <div className="text-left overflow-hidden">
                  <span className="block text-sm">Long Rest</span>
                  <span className="text-micro text-muted-foreground truncate">Full restore</span>
                </div>
              </Button>
            )}
          </div>
          {token && (
            <ActiveAbilitiesSection
              sheet={sheet}
              characterId={character.id}
              token={token}
            />
          )}
          <FeaturesSection sheet={sheet} onNotesChange={onNotesChange} />
        </div>
      </div>


      {/* Weapon Attack Dialog */}
      <WeaponAttackDialog 
        open={weaponAttackDialogOpen} 
        onOpenChange={setWeaponAttackDialogOpen}
        selectedWeapon={selectedWeapon}
        sheet={sheet}
        freeHands={freeHands} // Pass freeHands calculation
      />

      {/* Die Options Dialog */}
      <DieOptions
        open={dieOptionsOpen}
        onOpenChange={setDieOptionsOpen}
        onRoll={handleGenericRoll}
      />

      {/* HP Management Dialog */}
      <Dialog 
        open={expandedPanel === 'hp'} 
        onOpenChange={(open) => {
          if (!open) {
            setExpandedPanel(null);
            dispatchClearDice();
          }
        }}
      >
        <DialogContent className="max-w-xs">
          <DialogHeader>
            <DialogTitle className="text-center text-primary font-display">Manage HP</DialogTitle>
            <DialogDescription className="sr-only">Update your current and temporary health points.</DialogDescription>
          </DialogHeader>
          {onHPChange && (
            <HPPanel 
              currentHP={sheet.current_hp}
              maxHP={sheet.max_hp}
              tempHP={sheet.temp_hp}
              onHPChange={onHPChange} 
              onTempHPChange={onTempHPChange}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Hit Dice Dialog */}
      <Dialog 
        open={expandedPanel === 'hitdice'} 
        onOpenChange={(open) => {
          if (!open) {
            setExpandedPanel(null);
            dispatchClearDice();
          }
        }}
      >
        <DialogContent className="max-w-xs">
          <DialogHeader>
            <DialogTitle className="text-center text-primary font-display">Short Rest Healing</DialogTitle>
            <DialogDescription className="sr-only">Spend hit dice to recover health points.</DialogDescription>
          </DialogHeader>
          {onUseHitDice && <HitDicePanel sheet={sheet} onUseHitDice={onUseHitDice} />}
        </DialogContent>
      </Dialog>

      {/* Death Popup Dialog */}
      <Dialog 
        open={deathPopupOpen} 
        onOpenChange={(open) => {
          setDeathPopupOpen(open);
          if (!open) dispatchClearDice();
        }}
      >
        <DialogContent className="max-w-md bg-transparent border-none shadow-none flex flex-col items-center justify-center p-0" noPadding>
          <DialogHeader className="sr-only">
            <DialogTitle>Character Death</DialogTitle>
            <DialogDescription>Your character has reached 0 hit points.</DialogDescription>
          </DialogHeader>
          <img 
            src="/fading.gif" 
            alt="Fading away" 
            className="w-full h-auto rounded-lg"
          />
          <h2 className="text-3xl font-display text-white glow-text mt-4">YOU DIED</h2>
        </DialogContent>
      </Dialog>

      {/* Class Features Dialog (Sanguinist moral seesaw, Lorewright tools, etc.) */}
      {token && (
        <Dialog 
          open={classFeaturesOpen} 
          onOpenChange={(open) => {
            setClassFeaturesOpen(open);
            if (!open) dispatchClearDice();
          }}
        >
          <DialogContent className="max-w-md">
            <DialogHeader>
              <DialogTitle className="text-center text-primary font-display">
                {sheet.class.name === 'The Sanguinist' && 'The Moral Seesaw'}
                {sheet.class.name === 'The Lorewright' && 'Lorewright Tools'}
                {sheet.class.name === 'The Mutagen' && 'Feral Mechanics'}
                {sheet.class.name === 'The Powder Mage' && 'Continuous Ignition'}
                {sheet.class.name === 'The Piston Brawler' && 'Piston Core'}
              </DialogTitle>
              <DialogDescription className="sr-only">
                Class-specific feature controls.
              </DialogDescription>
            </DialogHeader>
            {sheet.class.name === 'The Sanguinist' && (
              <SanguinistFeaturesCard sheet={sheet} token={token} />
            )}
            {sheet.class.name === 'The Lorewright' && (
              <LorewrightFeaturesCard sheet={sheet} token={token} />
            )}
            {sheet.class.name === 'The Mutagen' && (
              <MutagenFeaturesCard sheet={sheet} token={token} />
            )}
            {sheet.class.name === 'The Powder Mage' && (
              <PowderMageFeaturesCard sheet={sheet} token={token} />
            )}
            {sheet.class.name === 'The Piston Brawler' && (
              <PistonBrawlerFeaturesCard sheet={sheet} token={token} />
            )}
          </DialogContent>
        </Dialog>
      )}

      {/* Loot Modal */}
      {token && (
        <LootModal
          isOpen={isLootModalOpen}
          onClose={() => setIsLootModalOpen(false)}
          characterId={character.id}
          token={token}
        />
      )}
    </div>
  );
}