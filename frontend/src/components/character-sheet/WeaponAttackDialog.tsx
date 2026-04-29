import { useState, useEffect } from 'react';
import { RaIcon } from '@/components/ui/RaIcon';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import { ApiCharacterWeapon, NormalizedCharacterSheet, ApiWeapon, ApiWeaponDamage } from '@/types/game';
import { rollD20, rollDiceNotation, dispatchClearDice, parseDiceNotation } from '@/lib/dice';
import { useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/context/AuthContext';
import { gainResource } from '@/lib/api/resources';
import { useToast } from '@/hooks/use-toast';

const MODIFIER_DISPLAY_NAMES: Record<string, string> = {
  piston_core: 'Piston Core',
  venom_coating: 'Venom Coating',
  shadow_coating: 'Shadow Coating',
  lethal_coating: 'Lethal Coating',
  silence_coating: 'Silence Coating',
};

interface WeaponAttackDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  selectedWeapon: ApiCharacterWeapon | null;
  sheet: NormalizedCharacterSheet;
  freeHands: number; // New prop for available hand slots
}

type WeaponAttackPhase = 'attack' | 'damage' | 'complete';

export function WeaponAttackDialog({ open, onOpenChange, selectedWeapon, sheet, freeHands }: WeaponAttackDialogProps) {
  const { modifiers } = sheet;
  const queryClient = useQueryClient();
  const { token } = useAuth();
  const { toast } = useToast();

  const [phase, setPhase] = useState<WeaponAttackPhase>('attack');
  const [attackRoll, setAttackRoll] = useState<{
    total: number;
    natural: number;
    modifier: number;
  } | null>(null);
  const [damageRolls, setDamageRolls] = useState<{
    damageType: string;
    dice: string;
    rolls: number[];
    modifier: number;
    total: number;
  }[]>([]);

  const [weaponMode, setWeaponMode] = useState<'one_handed' | 'two_handed' | 'transformed'>('one_handed');
  const [effectiveWeapon, setEffectiveWeapon] = useState<ApiWeapon | null>(null);
  const [twoHandedOverrideDamage, setTwoHandedOverrideDamage] = useState<string | null>(null);


  // Reset state when dialog opens with a new weapon
  useEffect(() => {
    if (open && selectedWeapon) {
      setTimeout(() => {
        setPhase('attack');
        setAttackRoll(null);
        setDamageRolls([]);
        setWeaponMode('one_handed'); // Reset to one-handed by default
        setTwoHandedOverrideDamage(null);
        setEffectiveWeapon(selectedWeapon.weapon); // Start with the base weapon
      }, 0);
    } else {
      dispatchClearDice();
    }

    // Ensure we clear dice if the dialog unmounts while open (e.g. navigation)
    return () => {
      dispatchClearDice();
    };
  }, [open, selectedWeapon]);

  // Effect to update effectiveWeapon when weaponMode changes
  useEffect(() => {
    if (!selectedWeapon) {
      setTimeout(() => setEffectiveWeapon(null), 0);
      return;
    }

    let tempEffectiveWeapon: ApiWeapon = { ...selectedWeapon.weapon };
    let currentProperties = [...(selectedWeapon.weapon.properties || [])];

    // Handle two-handed mode
    if (weaponMode === 'two_handed') {
      const weaponAllowsTwoHanded = selectedWeapon.weapon.properties?.some(p => p.toLowerCase() === 'versatile' || p.toLowerCase() === 'two-handed') ||
                                     selectedWeapon.weapon.two_handed_damage_dice;
      
      if (weaponAllowsTwoHanded) {
        // Use versatile damage dice if available, otherwise two-handed specific damage dice
        const dmgDice = selectedWeapon.weapon.two_handed_damage_dice || selectedWeapon.weapon.versatile_damage_dice;
        if (dmgDice) {
          tempEffectiveWeapon.damages = tempEffectiveWeapon.damages?.map((d, i) => i === 0 ? { ...d, damage_dice: dmgDice } : d) || [];
          setTimeout(() => setTwoHandedOverrideDamage(dmgDice), 0);
        }

        // Add two-handed specific properties
        if (selectedWeapon.weapon.two_handed_properties) {
          currentProperties = [...currentProperties, ...selectedWeapon.weapon.two_handed_properties];
        }
        if (!currentProperties.some(p => p.toLowerCase() === 'two-handed')) {
          currentProperties.push('Two-Handed');
        }
      }
    } else {
      setTimeout(() => setTwoHandedOverrideDamage(null), 0);
    }

    // Handle transformed mode
    if (weaponMode === 'transformed') {
      if (selectedWeapon.transformed_weapon_data) { // Full weapon swap
        tempEffectiveWeapon = { ...selectedWeapon.transformed_weapon_data };
        currentProperties = [...(selectedWeapon.transformed_weapon_data.properties || [])];
      } else if (selectedWeapon.weapon.transformed_self_damage_dice || selectedWeapon.weapon.transformed_self_properties) {
        // Self-modifying transformation (e.g., Saw Cleaver extending)
        if (selectedWeapon.weapon.transformed_self_damage_dice) {
          tempEffectiveWeapon.damages = tempEffectiveWeapon.damages?.map((d, i) => i === 0 ? { ...d, damage_dice: selectedWeapon.weapon.transformed_self_damage_dice! } : d) || [];
        }
        if (selectedWeapon.weapon.transformed_self_properties) {
          currentProperties = [...currentProperties, ...selectedWeapon.weapon.transformed_self_properties];
        }
      }
    }
    
    // Apply combined properties
    tempEffectiveWeapon.properties = Array.from(new Set(currentProperties));

    setTimeout(() => setEffectiveWeapon(tempEffectiveWeapon), 0);
  }, [weaponMode, selectedWeapon]);

    const handleAttackRoll = async () => {
      if (!effectiveWeapon) return;

      const weapon = effectiveWeapon;

      // Determine attack modifier based on weapon's AttackModifier field

      const attackAbility = (weapon.attack_modifier || 'Strength').toLowerCase() as keyof typeof modifiers;

      let baseModifier = 0;

  

      if (attackAbility && modifiers[attackAbility] !== undefined) {

        // Use the specified ability modifier + proficiency

        baseModifier = (modifiers[attackAbility] as number) + modifiers.proficiency;

      } else {

        // Fallback to legacy range-based detection

        const isRanged = weapon.range_type?.toLowerCase() === 'ranged';

        baseModifier = isRanged ? modifiers.ranged_attack : modifiers.melee_attack;

      }

  

      // Add weapon's fixed attack bonus if present (e.g. "+1 weapon")

      let weaponBonus = 0;

      // Check if attack_modifier is a number (e.g. "+1 weapon")

      const parsed = parseInt(weapon.attack_modifier || '', 10);

      if (!isNaN(parsed) && weapon.attack_modifier.startsWith('+')) weaponBonus = parsed; // Only if explicitly a bonus

  

      const totalModifier = baseModifier + weaponBonus;

      const result = await rollD20(totalModifier, `${weapon.name} Attack`);
      if (!result) return;

      setAttackRoll({
        total: result.total,
        natural: result.natural,
        modifier: totalModifier,
      });

      setPhase('damage');

    };

  const handleDamageRoll = async () => {
    if (!effectiveWeapon || !effectiveWeapon.damages) return;

    const weapon = effectiveWeapon;
    const isCrit = attackRoll?.natural === 20;
    const rolls: typeof damageRolls = [];

    // Determine ability modifier to add to damage
    const attackAbility = (weapon.attack_modifier || 'Strength').toLowerCase() as keyof typeof modifiers;
    let damageAbilityMod = 0;

    if (attackAbility && modifiers[attackAbility] !== undefined) {
      damageAbilityMod = modifiers[attackAbility] as number;
    } else {
      const isRanged = weapon.range_type?.toLowerCase() === 'ranged';
      const hasFinesse = weapon.properties?.some(p => p.toLowerCase().includes('finesse'));
      damageAbilityMod = isRanged ? modifiers.dexterity : modifiers.strength;
      if (hasFinesse) {
        damageAbilityMod = Math.max(modifiers.strength, modifiers.dexterity);
      }
    }

    // Roll base weapon damages
    for (const [idx, damage] of weapon.damages.entries()) {
      let diceNotation = twoHandedOverrideDamage && idx === 0 ? twoHandedOverrideDamage : damage.damage_dice;

      // Special case for Sanguinist Bite scaling
      if (weapon.name === 'Bite' && sheet.bite_damage_dice) {
        diceNotation = `${sheet.bite_damage_dice}d8`;
      }

      const abilityMod = idx === 0 ? damageAbilityMod : 0;
      const modStr = abilityMod !== 0 ? (abilityMod > 0 ? `+${abilityMod}` : `${abilityMod}`) : '';

      // On a crit, double the dice count so all dice roll in one animation.
      const parsed = parseDiceNotation(diceNotation);
      const critMultiplier = isCrit ? 2 : 1;
      const animNotation = `${parsed.count * critMultiplier}d${parsed.sides}`;
      const fullNotation = `${diceNotation}${modStr}`;

      const result = await rollDiceNotation(animNotation === `${parsed.count}d${parsed.sides}` ? fullNotation : animNotation, `${weapon.name} Damage`);
      if (!result) return;

      // If we animated with doubled dice (crit), the modifier was not included — add it back.
      const total = isCrit
        ? result.rolls.reduce((a, b) => a + b, 0) + abilityMod
        : result.total;

      rolls.push({
        damageType: damage.damage_type,
        dice: isCrit ? `${diceNotation} (x2 crit)` : diceNotation,
        rolls: result.rolls,
        modifier: abilityMod,
        total,
      });
    }

    // Roll modifier bonus damage (Piston Core, Coatings, etc.)
    for (const mod of selectedWeapon.active_modifiers?.filter(m => m.is_active) ?? []) {
      for (const bd of mod.bonus_damage ?? []) {
        const isFlat = bd.dice.match(/^[+-]?\d+$/);
        if (isFlat) {
          rolls.push({
            damageType: bd.damage_type || 'Modifier',
            dice: `${MODIFIER_DISPLAY_NAMES[mod.modifier_type] || mod.modifier_type}`,
            rolls: [parseInt(bd.dice, 10)],
            modifier: 0,
            total: parseInt(bd.dice, 10),
          });
        } else {
          const parsed = parseDiceNotation(bd.dice);
          const critMultiplier = isCrit ? 2 : 1;
          const animNotation = `${parsed.count * critMultiplier}d${parsed.sides}`;

          const modLabel = MODIFIER_DISPLAY_NAMES[mod.modifier_type] || mod.modifier_type;
          const result = await rollDiceNotation(
            animNotation,
            `${weapon.name} — ${modLabel}${bd.damage_type ? ` (${bd.damage_type})` : ''}`
          );
          if (!result) return;

          const total = result.rolls.reduce((a, b) => a + b, 0) + result.modifier;
          rolls.push({
            damageType: bd.damage_type || 'Modifier',
            dice: isCrit ? `${bd.dice} (x2 crit)` : bd.dice,
            rolls: result.rolls,
            modifier: result.modifier,
            total,
          });
        }
      }
    }

    setDamageRolls(rolls);
    setPhase('complete');

    // The Thirst: on a Bite hit you regain Ichor equal to proficiency bonus (tracked when damage is rolled).
    const weaponName = (effectiveWeapon.name || selectedWeapon.weapon.name || '').trim();
    const isBite =
      weaponName === 'Bite' ||
      selectedWeapon.character_weapon_id === 'virtual-bite';
    if (sheet.class.name === 'The Sanguinist' && isBite && token) {
      const prof = modifiers.proficiency;
      if (prof > 0) {
        try {
          await gainResource(sheet.character.id, 'max_blood_ichor', prof, token);
          await queryClient.invalidateQueries({ queryKey: ['character-sheet', sheet.character.id] });
          await queryClient.refetchQueries({ queryKey: ['character-sheet', sheet.character.id] });
        } catch (e) {
          console.error('Failed to apply Bite ichor regen:', e);
          toast({
            title: 'Could not update Blood Ichor',
            description: e instanceof Error ? e.message : 'Resource row may be missing — reload the character sheet.',
            variant: 'destructive',
          });
        }
      }
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent 
        className="max-w-sm"
        onPointerDownOutside={(e) => e.preventDefault()}
      >
        {selectedWeapon && (
          <div className="py-2">
            <DialogHeader>
              <DialogTitle className="text-center font-tome-subheading text-primary flex items-center justify-center gap-2">
                {selectedWeapon.is_primary && <RaIcon name="crown" className="text-sm text-primary" />}
                <RaIcon name="crossed-swords" className="text-base" />
                {effectiveWeapon?.name || selectedWeapon.custom_name || selectedWeapon.weapon.name}
              </DialogTitle>
              <DialogDescription className="sr-only">
                Perform an attack and roll damage with your selected weapon.
              </DialogDescription>
              {/* Show active modifiers */}
              {selectedWeapon.active_modifiers && selectedWeapon.active_modifiers.length > 0 && (
                <div className="flex justify-center gap-1 mt-2">
                  {selectedWeapon.active_modifiers.map((m, i) => (
                    <Badge
                      key={i}
                      variant={m.is_active ? "default" : "outline"}
                      className={cn(
                        "text-micro h-5 gap-1",
                        m.is_active ? "bg-primary/80" : "opacity-60"
                      )}
                    >
                      <RaIcon name="lightning-bolt" className="text-xs" />
                      {MODIFIER_DISPLAY_NAMES[m.modifier_type] || m.modifier_type}
                    </Badge>
                  ))}
                </div>
              )}

              {/* Weapon Properties */}
              {effectiveWeapon?.properties && effectiveWeapon.properties.length > 0 && (
                <div className="flex flex-wrap justify-center gap-1 mt-2">
                  {effectiveWeapon.properties.map((prop, i) => (
                    <Badge
                      key={i}
                      variant="outline"
                      size="tiny"
                      font="tomeMarginalia"
                      className="opacity-70"
                    >
                      {prop}
                    </Badge>
                  ))}
                </div>
              )}

              {/* Weapon Description & Secondary Effects */}
              {(effectiveWeapon?.description || effectiveWeapon?.secondary_effect) && (
                <div className="mt-3 px-4 space-y-2">
                  {effectiveWeapon.description && (
                    <p className="text-center text-xs text-muted-foreground font-tome-marginalia italic">
                      {effectiveWeapon.description}
                    </p>
                  )}
                  {effectiveWeapon.secondary_effect && (
                    <div className="text-center text-micro text-primary/90 font-tome-marginalia border-t border-primary/10 pt-2">
                      <span className="font-bold uppercase tracking-tighter mr-1 text-primary">Effect:</span>
                      {effectiveWeapon.secondary_effect}
                    </div>
                  )}
                </div>
              )}

              {/* Transformation Description */}
              {selectedWeapon.weapon.transformation_description && (
                 <div className="mt-3 px-4 space-y-2">
                   <div className="text-center text-micro text-yellow-400/90 font-tome-marginalia border-t border-yellow-400/10 pt-2">
                      <span className="font-bold uppercase tracking-tighter mr-1 text-yellow-400">Transform:</span>
                      {selectedWeapon.weapon.transformation_description}
                   </div>
                 </div>
              )}
            </DialogHeader>

                          {/* Weapon Mode Controls and Current Mode Display */}
                          {effectiveWeapon && (
                            <div className="flex flex-col items-center gap-2 mt-4">
                              <div className="text-sm font-tome-subheading text-muted-foreground">
                                Mode: <span className="text-primary capitalize">{weaponMode.replace('_', ' ')}</span>
                              </div>
                              <div className="flex flex-wrap justify-center gap-2">
                                {/* Two-Handed Button */}                  {(selectedWeapon.weapon.properties?.some(p => p.toLowerCase() === 'versatile' || p.toLowerCase() === 'two-handed') ||
                    selectedWeapon.weapon.two_handed_damage_dice) && (
                    <Button
                      variant={weaponMode === 'two_handed' ? 'default' : 'outline'}
                      size="sm"
                      onClick={() => setWeaponMode('two_handed')}
                      disabled={
                        weaponMode === 'two_handed' ||
                        (
                          // Check if switching to two-handed would consume an additional hand
                          (selectedWeapon.weapon.properties?.some(p => p.toLowerCase() === 'versatile') || 
                           (selectedWeapon.weapon.two_handed_damage_dice && !selectedWeapon.weapon.properties?.some(p => p.toLowerCase() === 'two-handed'))) &&
                          freeHands < 1 &&
                          weaponMode === 'one_handed'
                        )
                      }
                    >
                      Two-Handed
                    </Button>
                  )}

                  {/* Transform Button */}
                  {(selectedWeapon.weapon.transformed_to_weapon_id || selectedWeapon.weapon.properties?.some(p => p.toLowerCase() === 'transformative')) && (
                    <Button
                      variant={weaponMode === 'transformed' ? 'default' : 'outline'}
                      size="sm"
                      onClick={() => {
                        // For transformations to a different weapon (with transformed_to_weapon_id)
                        if (selectedWeapon.weapon.transformed_to_weapon_id && selectedWeapon.transformed_weapon_data) {
                          setWeaponMode('transformed');
                        } else if (selectedWeapon.weapon.properties?.some(p => p.toLowerCase() === 'transformative')) {
                            // This branch is for self-modifying transformations (like Saw Cleaver)
                            // where there's no separate transformed_weapon_data.
                            // The effectiveWeapon logic currently won't change unless new explicit
                            // fields for self-transformation are added to the Weapon model.
                            // For now, this just sets the mode to 'transformed'.
                            setWeaponMode('transformed');
                        }
                      }}
                      disabled={weaponMode === 'transformed'}
                    >
                      Transform
                    </Button>
                  )}

                  {/* One-Handed/Base Button (if applicable) */}
                  {(weaponMode !== 'one_handed' && (selectedWeapon.weapon.properties?.some(p => p.toLowerCase() === 'versatile') || selectedWeapon.weapon.two_handed_damage_dice || selectedWeapon.weapon.transformed_to_weapon_id)) && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setWeaponMode('one_handed')}
                    >
                      Base Mode
                    </Button>
                  )}
                </div>
              </div>
            )}

            <div className="mt-4 space-y-6">
              {/* PHASE 1: ATTACK ROLL */}
              <div className={cn(
                'rounded-lg border p-4 transition-all',
                phase === 'attack' ? 'border-primary ring-1 ring-primary/20 bg-primary/5' : 'opacity-50 border-border bg-muted/20'
              )}>
                <div className="flex justify-between items-center mb-3">
                  <span className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Attack Roll</span>
                  {attackRoll && (
                    <Badge
                      variant={
                        attackRoll.natural === 20
                          ? 'success'
                          : 'outline'
                      }
                      theme={
                        attackRoll.natural === 1
                          ? 'destructive-outline'
                          : undefined
                      }
                    >
                      {attackRoll.natural === 20 ? 'CRIT!' : attackRoll.natural === 1 ? 'Fumble' : 'Rolled'}
                    </Badge>
                  )}
                </div>
                
                {attackRoll ? (
                  <div className="text-center h-12 flex items-center justify-center">
                    {/* Blank result area as requested */}
                  </div>
                ) : (
                  <Button onClick={handleAttackRoll} className="w-full gap-2" variant="default">
                    <RaIcon name="archery-target" className="text-sm" />
                    Roll Attack
                  </Button>
                )}
              </div>

              {/* PHASE 2: DAMAGE ROLL */}
              <div className={cn(
                'rounded-lg border p-4 transition-all',
                phase === 'damage' ? 'border-primary ring-1 ring-primary/20 bg-primary/5' : 
                phase === 'complete' ? 'border-border bg-muted/20' : 'opacity-20 border-border bg-muted/10 pointer-events-none'
              )}>
                <div className="flex justify-between items-center mb-3">
                  <span className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Damage Roll</span>
                </div>

                {damageRolls.length > 0 ? (
                  <div className="space-y-3 h-12 flex items-center justify-center">
                    {/* Blank result area as requested */}
                  </div>
                ) : (
                  <Button 
                    onClick={handleDamageRoll} 
                    className="w-full gap-2" 
                    variant="default"
                    disabled={phase !== 'damage'}
                  >
                    <RaIcon name="fire" className="text-sm" />
                    Roll Damage
                  </Button>
                )}
              </div>
            </div>

            <div className="mt-6 flex gap-2">
              <Button
                variant="outline"
                className="flex-1"
                onClick={() => onOpenChange(false)}
              >
                {phase === 'complete' ? 'Close' : 'Cancel'}
              </Button>
              {phase === 'complete' && (
                <Button
                  variant="default"
                  className="flex-1"
                  onClick={() => {
                    setPhase('attack');
                    setAttackRoll(null);
                    setDamageRolls([]);
                  }}
                >
                  Attack Again
                </Button>
              )}
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}