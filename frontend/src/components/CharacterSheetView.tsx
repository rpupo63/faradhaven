import { useState, useEffect, useRef } from 'react';
import { cn } from '@/lib/utils';
import { type NormalizedCharacterSheet, type ApiCharacterWeapon } from '@/types/game';
import { ChevronUp, ChevronDown, Sun, Moon } from 'lucide-react';
import { RaIcon } from '@/components/ui/RaIcon';
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
import { ActiveEffectsSection } from './character-sheet/ActiveEffectsSection';
import { EquipmentSection } from './character-sheet/EquipmentSection';
import { FeaturesSection } from './character-sheet/FeaturesSection';
import { ActiveAbilitiesSection } from './character-sheet/ActiveAbilitiesSection';
import { MoneyPanel } from './character-sheet/MoneyPanel';
import { WeaponAttackDialog } from './character-sheet/WeaponAttackDialog';
import { DieOptions } from './character-sheet/DieOptions';
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
  /** Switch to Spellbook and focus Spell Forge (e.g. Powder Mage Speed Dial). */
  onGoToSpellForge?: () => void;
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
  onGoToSpellForge,
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
    <div className={cn('space-y-4 relative pb-0', className)}>
      {/* Header */}
      {!hideHeader && (
        <div className="border-b border-border pb-2">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="font-display text-xl md:text-2xl text-primary glow-text">
                {character.name}
              </h2>
              <p className="text-muted-foreground font-tome-marginalia text-xs md:text-sm">
                Level {character.level} {character.lineageName || character.raceName} {character.className}
                {character.archetypeName && ` - ${character.archetypeName}`}
              </p>
            </div>
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

        {/* Mobile Action Buttons */}
        <div className="flex gap-2 px-1">
             <Button
                variant="outline"
                size="sm"
                onClick={() => setIsLootModalOpen(true)}
                title="Generate Loot"
                className="flex-1 gap-2 h-10 px-2 shrink-0"
              >
                <RaIcon name="gold-bar" className="text-sm" />
                <span>Loot</span>
              </Button>
              <div className="flex gap-2 shrink-0">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={onLevelDown}
                  disabled={character.level <= 1 || !onLevelDown}
                  title="Level Down"
                  className="gap-1 h-10 px-3"
                >
                  <ChevronDown className="h-4 w-4" />
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={onLevelUp}
                  disabled={character.level >= 20 || !onLevelUp}
                  title="Level Up"
                  className="gap-1 h-10 px-3"
                >
                  <ChevronUp className="h-4 w-4" />
                </Button>
              </div>
        </div>

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
                <ClassResourceDisplay
                  resources={sheet.class_resources}
                  characterId={character.id}
                  token={token}
                  sheet={sheet}
                  onGoToSpellForge={onGoToSpellForge}
                />
              )}
              {token && <ComponentInventorySection sheet={sheet} token={token} />}
              {activeEffects && activeEffects.length > 0 && (
                <ActiveEffectsSection characterId={sheet.character.id} effects={activeEffects} />
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

      {/* Main 4-column layout (desktop): Abilities | Skills | Middle | Features */}
      <div className="hidden md:grid grid-cols-1 gap-4 min-w-0 md:grid-cols-[90px_180px_1fr] lg:grid-cols-[100px_200px_1fr_280px]">
        
        {/* FAR LEFT COLUMN: Ability Scores */}
        <div className="order-2 md:order-none space-y-4 min-w-0">
          {/* Action Buttons */}
          <div className="flex flex-col gap-2">
             <Button
                variant="outline"
                size="sm"
                onClick={() => setIsLootModalOpen(true)}
                title="Generate Loot"
                className="w-full gap-2 h-10 px-2 shrink-0"
              >
                <RaIcon name="gold-bar" className="text-sm" />
                <span>Loot</span>
              </Button>
              <div className="grid grid-cols-2 gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={onLevelDown}
                  disabled={character.level <= 1 || !onLevelDown}
                  title="Level Down"
                  className="gap-1 h-10 px-2 shrink-0"
                >
                  <ChevronDown className="h-4 w-4" />
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={onLevelUp}
                  disabled={character.level >= 20 || !onLevelUp}
                  title="Level Up"
                  className="gap-1 h-10 px-2 shrink-0"
                >
                  <ChevronUp className="h-4 w-4" />
                </Button>
              </div>
          </div>

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

          {/* Class Resources */}
          {sheet.class_resources && sheet.class_resources.length > 0 && token && (
            <ClassResourceDisplay
              resources={sheet.class_resources}
              characterId={character.id}
              token={token}
              sheet={sheet}
              onGoToSpellForge={onGoToSpellForge}
            />
          )}

          {/* Component Inventory */}
          {token && (
            <ComponentInventorySection sheet={sheet} token={token} />
          )}

          {/* Active Effects */}
          {activeEffects && activeEffects.length > 0 && (
            <ActiveEffectsSection
              characterId={sheet.character.id}
              effects={activeEffects}
            />
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