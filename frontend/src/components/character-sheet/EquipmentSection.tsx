import { Swords, Target, Star, Shield, ShieldCheck, Hand } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { NormalizedCharacterSheet, ApiCharacterWeapon, ApiWeapon, ApiWeaponDamage } from '@/types/game';
import { updateEquipment } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';
import { useState } from 'react';
import { cn } from '@/lib/utils';

interface EquipmentSectionProps {
  sheet: NormalizedCharacterSheet;
  onWeaponClick: (weapon: ApiCharacterWeapon) => void;
  onGenerateLoot?: () => void;
  onEquipmentChange: () => void;
  isTwoHanded: (properties?: string[]) => boolean; // Passed from parent
  usedHands: number; // Passed from parent
}

const formatMod = (n: number) => (n >= 0 ? `+${n}` : `${n}`);

// Helper to calculate attack and damage modifiers for a weapon
function getWeaponModifiers(weapon: ApiWeapon, modifiers: NormalizedCharacterSheet['modifiers']) {
  // Attack Modifier Calculation
  const attackAbility = (weapon.attack_modifier || 'Strength').toLowerCase() as keyof typeof modifiers;
  let baseAttackModifier = 0;

  if (attackAbility && modifiers[attackAbility] !== undefined) {
    baseAttackModifier = (modifiers[attackAbility] as number) + modifiers.proficiency;
  } else {
    const isRanged = weapon.range_type?.toLowerCase() === 'ranged';
    baseAttackModifier = isRanged ? modifiers.ranged_attack : modifiers.melee_attack;
  }

  let weaponBonus = 0;
  const parsed = parseInt(weapon.attack_modifier || '', 10);
  if (!isNaN(parsed) && weapon.attack_modifier.startsWith('+')) {
    weaponBonus = parsed;
  }
  const totalAttackModifier = baseAttackModifier + weaponBonus;

  // Damage Ability Modifier Calculation (only the ability part, not proficiency)
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

  return { totalAttackModifier, damageAbilityMod };
}

export function EquipmentSection({ sheet, onWeaponClick, onGenerateLoot, onEquipmentChange, isTwoHanded, usedHands }: EquipmentSectionProps) {
  const { character, modifiers } = sheet;
  const { token } = useAuth();
  const [isLoading, setIsLoading] = useState<string | null>(null);
  const [equipError, setEquipError] = useState<string | null>(null);
  const [attackError, setAttackError] = useState<string | null>(null);

  const freeHands = 2 - usedHands;

  const handleEquipmentChange = async (
    itemId: string,
    isWeapon: boolean,
    equip: boolean,
    slot: 'armor' | 'shield' | 'weapon'
  ) => {
    if (!token) return;
    setIsLoading(itemId);
    setEquipError(null);
    try {
      await updateEquipment(character.id, { item_id: itemId, is_weapon: isWeapon, equip, slot }, token);
      onEquipmentChange();
    } catch (error) {
      const msg = error instanceof Error ? error.message : 'Failed to update equipment';
      setEquipError(msg);
    } finally {
      setIsLoading(null);
    }
  };
  
  const sortedWeapons = [...(sheet.inventory_weapons || [])].sort((a, b) => {
    if (a.is_equipped && !b.is_equipped) return -1;
    if (!a.is_equipped && b.is_equipped) return 1;
    return 0;
  });

  const armor = sheet.inventory_items?.filter(it => it.category === 'Armor').sort((a, b) => {
    const aEquipped = sheet.character.equipped_armor_id === a.id;
    const bEquipped = sheet.character.equipped_armor_id === b.id;
    if (aEquipped && !bEquipped) return -1;
    if (!aEquipped && bEquipped) return 1;
    return 0;
  });

  const shields = sheet.inventory_items?.filter(it => it.category === 'Shield').sort((a, b) => {
    const aEquipped = sheet.character.equipped_shield_id === a.id;
    const bEquipped = sheet.character.equipped_shield_id === b.id;
    if (aEquipped && !bEquipped) return -1;
    if (!aEquipped && bEquipped) return 1;
    return 0;
  });
  
  const otherItems = sheet.inventory_items?.filter(it => it.category !== 'Armor' && it.category !== 'Shield');

  return (
    <Card className="arcane-border bg-card">
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center justify-between text-base font-tome-subheading text-primary">
          <div className="flex items-center gap-2">
            <Swords className="h-4 w-4" />
            Equipment & Weapons
          </div>

        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Hand Slots Indicator */}
        <div className="flex flex-wrap items-center justify-between gap-2 text-xs border border-border/50 rounded p-2 bg-muted/10 min-w-0">
          <div className="flex items-center gap-1.5 text-muted-foreground min-w-0">
            <Hand className="h-3.5 w-3.5 shrink-0" />
            <span>Hand Slots</span>
          </div>
          <div className="flex items-center gap-1 shrink-0">
            {[0, 1].map(i => (
              <div
                key={i}
                className={cn(
                  'w-3 h-3 rounded-full border',
                  i < usedHands
                    ? 'bg-primary border-primary'
                    : 'bg-transparent border-muted-foreground/40'
                )}
              />
            ))}
            <span className="ml-1.5 text-muted-foreground">{freeHands}/2 free</span>
          </div>
        </div>

        {equipError && (
          <p className="text-xs text-destructive bg-destructive/10 rounded px-2 py-1">{equipError}</p>
        )}
        {attackError && (
          <p className="text-xs text-destructive bg-destructive/10 rounded px-2 py-1">{attackError}</p>
        )}

        {/* Weapons */}
        {sortedWeapons.length > 0 && (
          <div className="space-y-2">
            <p className="text-[10px] uppercase tracking-wider text-muted-foreground font-bold">Weapons</p>
            <div className="space-y-2">
              {sortedWeapons.map((cw) => {
                const isVirtual = cw.character_weapon_id === 'virtual-bite';
                const handCost = isTwoHanded(cw.weapon.properties) ? 2 : 1;
                const canEquip = cw.is_equipped || freeHands >= handCost;
                
                const { totalAttackModifier, damageAbilityMod } = getWeaponModifiers(cw.weapon, modifiers);

                return (
                  <div key={cw.character_weapon_id}
                    className={cn(`w-full min-w-0 text-left text-sm p-2 rounded border transition-colors group`,
                      cw.is_equipped ? 'border-primary/50 bg-primary/5' : 'border-border/50 bg-muted/10'
                    )}
                  >
                    <div className="flex flex-col gap-2 min-[480px]:flex-row min-[480px]:justify-between min-[480px]:items-start">
                      <button onClick={() => {
                        if (!cw.is_equipped) {
                          setAttackError('Weapon must be equipped to attack.');
                          setTimeout(() => setAttackError(null), 3000); // Clear error after 3 seconds
                          return;
                        }
                        onWeaponClick(cw);
                      }} className="min-w-0 flex-1 text-left">
                        <span className="font-bold text-primary group-hover:text-primary/90 flex flex-wrap items-center gap-x-1.5 gap-y-1">
                          {cw.is_equipped && <Star className="h-3 w-3 fill-primary shrink-0" />}
                          <Target className={`h-3 w-3 shrink-0 ${cw.is_equipped ? 'opacity-100' : 'opacity-0 group-hover:opacity-100'} transition-opacity`} />
                          <span className="break-words">{cw.custom_name || cw.weapon.name}</span>
                          <span className="text-xs text-muted-foreground font-normal w-full min-[480px]:w-auto min-[480px]:ml-2">
                            ({formatMod(totalAttackModifier)} Atk, {formatMod(damageAbilityMod)} Dmg)
                          </span>
                        </span>
                      </button>
                      <div className="flex flex-wrap items-center gap-1 shrink-0 justify-end">
                        {cw.weapon.properties?.includes('Transformative') && (
                          <Badge variant="outline" size="tiny" theme="warning">Transforms</Badge>
                        )}
                        {cw.weapon.properties?.includes('Versatile') && (
                          <Badge variant="outline" size="tiny">Versatile</Badge>
                        )}
                        {handCost === 2 && (
                          <Badge variant="outline" size="tiny">2H</Badge>
                        )}
                        {cw.is_primary && sheet.class.name === 'The Piston Brawler' && (
                          <Badge variant="outline" size="tiny" theme="info">⚙ Piston Core</Badge>
                        )}
                        {!isVirtual && (
                          <Button
                            size="xs"
                            variant="outline"
                            className="h-6"
                            onClick={() => handleEquipmentChange(cw.character_weapon_id, true, !cw.is_equipped, 'weapon')}
                            disabled={!!isLoading || (!cw.is_equipped && !canEquip)}
                            title={!cw.is_equipped && !canEquip ? 'Not enough free hands' : undefined}
                          >
                            {isLoading === cw.character_weapon_id ? '...' : (cw.is_equipped ? 'Unequip' : 'Equip')}
                          </Button>
                        )}
                      </div>
                    </div>
                    {/* Base damage */}
                    <div className="flex flex-wrap gap-2 mt-1">
                      {cw.weapon.damages?.map((d, i) => (
                        <Badge key={i} variant="secondary" size="sm">
                          {d.damage_dice} {d.damage_type}
                        </Badge>
                      ))}
                      {/* Modifier bonus damage */}
                      {cw.active_modifiers?.filter(m => m.is_active).map((m, i) =>
                        m.bonus_damage?.map((bd, j) => (
                          <Badge key={`mod-${i}-${j}`} variant="accent" size="sm">
                            {bd.dice} {bd.damage_type}
                          </Badge>
                        ))
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}

        {/* Armor */}
        {armor && armor.length > 0 && (
          <div className="space-y-2">
            <p className="text-[10px] uppercase tracking-wider text-muted-foreground font-bold">Armor</p>
            <div className="space-y-2">
              {armor.map((it) => {
                const isEquipped = sheet.character.equipped_armor_id === it.id;
                return (
                  <div key={it.id} className={cn("text-sm p-2 rounded border min-w-0", isEquipped ? 'border-primary/50 bg-primary/5' : 'border-border/50 bg-muted/10')}>
                    <div className="flex flex-col gap-2 min-[400px]:flex-row min-[400px]:justify-between min-[400px]:items-start">
                      <div className="min-w-0">
                        <span className="font-bold text-primary flex items-center gap-1.5 flex-wrap">
                          {isEquipped && <ShieldCheck className="h-3 w-3 text-primary shrink-0" />}
                          <span className="break-words">{it.name}</span>
                        </span>
                        {it.armor_type && <span className="text-[10px] text-muted-foreground">{it.armor_type} Armor</span>}
                      </div>
                      <Button size="xs" variant="outline" className="h-6 shrink-0 self-end min-[400px]:ml-2" onClick={() => handleEquipmentChange(it.id, false, !isEquipped, 'armor')} disabled={!!isLoading}>
                        {isLoading === it.id ? '...' : (isEquipped ? 'Unequip' : 'Equip')}
                      </Button>
                    </div>
                    <div className="flex flex-wrap gap-2 mt-1">
                      {it.base_ac != null && (
                        <Badge variant="secondary" size="sm">AC {it.base_ac}</Badge>
                      )}
                      {it.stealth_disadvantage && (
                        <Badge variant="warning" size="sm">Stealth Disadv.</Badge>
                      )}
                      {it.strength_requirement != null && (
                        <Badge variant="outline" size="sm">Str {it.strength_requirement}+</Badge>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}
        
        {/* Shields */}
        {shields && shields.length > 0 && (
          <div className="space-y-2">
            <p className="text-[10px] uppercase tracking-wider text-muted-foreground font-bold">Shield</p>
            <div className="space-y-2">
              {shields.map((it) => {
                const isEquipped = sheet.character.equipped_shield_id === it.id;
                const canEquipShield = isEquipped || freeHands >= 1;
                return (
                  <div key={it.id} className={cn("text-sm p-2 rounded border min-w-0", isEquipped ? 'border-primary/50 bg-primary/5' : 'border-border/50 bg-muted/10')}>
                    <div className="flex flex-col gap-2 min-[400px]:flex-row min-[400px]:justify-between min-[400px]:items-start">
                       <span className="font-bold text-primary flex items-center gap-1.5 min-w-0 break-words">
                          {isEquipped && <Shield className="h-3 w-3 text-primary shrink-0" />}
                          {it.name}
                        </span>
                      <Button
                        size="xs"
                        variant="outline"
                        className="h-6 shrink-0 self-end min-[400px]:ml-2"
                        onClick={() => handleEquipmentChange(it.id, false, !isEquipped, 'shield')}
                        disabled={!!isLoading || (!isEquipped && !canEquipShield)}
                        title={!isEquipped && !canEquipShield ? 'Not enough free hands' : undefined}
                      >
                        {isLoading === it.id ? '...' : (isEquipped ? 'Unequip' : 'Equip')}
                      </Button>
                    </div>
                     {it.effects && <p className="text-[10px] text-muted-foreground mt-1 italic">{it.effects}</p>}
                  </div>
                );
              })}
            </div>
          </div>
        )}

        {/* Other Items */}
        {otherItems && otherItems.length > 0 && (
          <div className="space-y-2">
            <p className="text-[10px] uppercase tracking-wider text-muted-foreground font-bold">Items & Gear</p>
            <div className="space-y-2">
              {otherItems.map((it) => (
                <div key={it.id} className="text-sm p-2 rounded border border-border/50 bg-muted/10 min-w-0">
                  <div className="flex flex-col gap-1 min-[400px]:flex-row min-[400px]:justify-between min-[400px]:items-start">
                    <span className="font-bold text-primary break-words min-w-0">{it.name}</span>
                    <span className="text-[10px] text-muted-foreground shrink-0">{it.category}</span>
                  </div>
                  {it.effects && (
                    <p className="text-[10px] text-muted-foreground mt-1 italic">{it.effects}</p>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Legacy/Generic Inventory */}
        {character.inventory && character.inventory.length > 0 && (
          <div className="space-y-2 border-t border-border/30 pt-2">
            <p className="text-[10px] uppercase tracking-wider text-muted-foreground font-bold">Other Gear</p>
            <ul className="text-xs font-tome-marginalia space-y-1">
              {character.inventory.map((item, idx) => (
                <li key={idx} className="flex items-center gap-2 text-muted-foreground">
                  <span className="w-1 h-1 rounded-full bg-primary/30" />
                  {item}
                </li>
              ))}
            </ul>
          </div>
        )}

        {(!character.inventory?.length && !sheet.inventory_weapons?.length && !sheet.inventory_items?.length) && (
          <p className="text-sm font-tome-marginalia text-muted-foreground italic text-center py-4">
            No equipment tracked yet
          </p>
        )}

        {/* Attack bonuses */}
        <div className="mt-3 grid grid-cols-1 min-[400px]:grid-cols-3 gap-2 text-center border-t border-border/30 pt-3">
          <div className="rounded border border-border p-2">
            <p className="text-xs text-muted-foreground">Melee</p>
            <p className="font-display text-primary">{formatMod(modifiers.melee_attack)}</p>
          </div>
          <div className="rounded border border-border p-2">
            <p className="text-xs text-muted-foreground">Ranged</p>
            <p className="font-display text-primary">{formatMod(modifiers.ranged_attack)}</p>
          </div>
          <div className="rounded border border-border p-2">
            <p className="text-xs text-muted-foreground">Spell</p>
            <p className="font-display text-primary">{formatMod(modifiers.spell_attack)}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}