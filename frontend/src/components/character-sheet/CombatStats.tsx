import { Heart, Dices, Shield, Footprints, Pencil } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Card, CardContent } from '@/components/ui/card';
import { NormalizedCharacterSheet } from '@/types/game';

interface CombatStatsProps {
  sheet: NormalizedCharacterSheet;
  expandedPanel: 'hp' | 'hitdice' | null;
  setExpandedPanel: (panel: 'hp' | 'hitdice' | null) => void;
  onHPChange?: (delta: number) => void | Promise<void>;
  onUseHitDice?: (rolls: number[]) => void | Promise<void>;
}

export function CombatStats({
  sheet,
  expandedPanel,
  setExpandedPanel,
  onHPChange,
  onUseHitDice,
}: CombatStatsProps) {
  const {
    current_hp: currentHP,
    max_hp: maxHP,
    temp_hp: tempHP,
    hit_dice_total: hitDiceTotal,
    hit_dice_remaining: hitDiceRemaining,
    hit_die: hitDie,
    ac,
    speed,
  } = sheet;

  const hasUnarmoredDefense = sheet.class_level.level_features?.some(f => f.name === 'Unarmored Defense');
  const hasArmorEquipped = !!sheet.character.equipped_armor_id;
  const hasShieldEquipped = !!sheet.character.equipped_shield_id;
  const equippedArmorItem = hasArmorEquipped
    ? sheet.inventory_items?.find(it => it.id === sheet.character.equipped_armor_id)
    : undefined;

  return (
    <div className="space-y-4">
      {/* Combat Stats Row: HP, Hit Dice, AC */}
      <div className="grid gap-3 grid-cols-1 min-[400px]:grid-cols-3">
        {/* HP Card - Clickable */}
        <Card
          className={cn(
            'arcane-border bg-card transition-colors relative cursor-pointer hover:border-red-500/50',
            expandedPanel === 'hp' && 'border-red-500/50'
          )}
          onClick={() => {
            if (!onHPChange) {
              alert('Please log in to modify HP');
              return;
            }
            setExpandedPanel(expandedPanel === 'hp' ? null : 'hp');
          }}
        >
          <Pencil className="absolute top-2 right-2 h-3 w-3 text-muted-foreground/50" />
          <CardContent className="p-4 text-center">
            <div className="flex items-center justify-center gap-1.5 mb-1">
              <Heart className="h-4 w-4 text-red-500" />
              <span className="text-xs font-tome-marginalia text-muted-foreground uppercase">HP</span>
            </div>
            <p className="font-display text-2xl text-primary">
              {currentHP}<span className="text-sm text-muted-foreground">/{maxHP}</span>
            </p>
            {tempHP > 0 && (
              <p className="text-xs text-blue-500 font-tome-marginalia">+{tempHP} temp</p>
            )}
          </CardContent>
        </Card>

        {/* Hit Dice Card - Clickable */}
        <Card
          className={cn(
            'arcane-border bg-card transition-colors relative cursor-pointer hover:border-primary/50',
            expandedPanel === 'hitdice' && 'border-primary/50'
          )}
          onClick={() => {
            if (!onUseHitDice) {
              alert('Please log in to use hit dice');
              return;
            }
            setExpandedPanel(expandedPanel === 'hitdice' ? null : 'hitdice');
          }}
        >
          <Pencil className="absolute top-2 right-2 h-3 w-3 text-muted-foreground/50" />
          <CardContent className="p-4 text-center">
            <div className="flex items-center justify-center gap-1.5 mb-1">
              <Dices className="h-4 w-4 text-primary" />
              <span className="text-xs font-tome-marginalia text-muted-foreground uppercase">Hit Dice</span>
            </div>
            <p className="font-display text-2xl text-primary">
              {hitDiceRemaining}<span className="text-sm text-muted-foreground">/{hitDiceTotal}</span>
            </p>
            <p className="text-xs text-muted-foreground font-tome-marginalia">d{hitDie}</p>
          </CardContent>
        </Card>

        {/* AC Card */}
        <Card className="arcane-border bg-card">
          <CardContent className="p-4 text-center">
            <div className="flex items-center justify-center gap-1.5 mb-1">
              <Shield className="h-4 w-4 text-primary" />
              <span className="text-xs font-tome-marginalia text-muted-foreground uppercase">AC</span>
            </div>
            <p className="font-display text-2xl text-primary">{ac}</p>
            {equippedArmorItem ? (
              <p className="text-[10px] text-muted-foreground font-tome-marginalia leading-tight mt-1">
                {equippedArmorItem.name}
                {hasShieldEquipped && ' + Shield'}
              </p>
            ) : hasUnarmoredDefense ? (
              <p className="text-[10px] text-muted-foreground font-tome-marginalia leading-tight mt-1">
                Unarmored Defense
                {hasShieldEquipped && ' + Shield'}
              </p>
            ) : hasShieldEquipped ? (
              <p className="text-[10px] text-muted-foreground font-tome-marginalia leading-tight mt-1">
                Unarmored + Shield
              </p>
            ) : null}
          </CardContent>
        </Card>
      </div>

      {/* Sneak Attack Dice (Vapor Blade) */}
      {(sheet.class_level.sneak_attack_dice ?? 0) > 0 && (
        <Card className="arcane-border bg-card">
          <CardContent className="p-3 text-center">
            <p className="text-[10px] font-tome-marginalia text-muted-foreground uppercase tracking-wider">Sneak Strike</p>
            <p className="text-xl font-display text-primary">{sheet.class_level.sneak_attack_dice}d6</p>
          </CardContent>
        </Card>
      )}

      {/* Speed */}
      {speed != null && (
        <Card className="arcane-border bg-card">
          <CardContent className="p-3 text-center">
            <div className="flex items-center justify-center gap-1.5 mb-1">
              <Footprints className="h-4 w-4 text-primary" />
              <span className="text-xs font-tome-marginalia text-muted-foreground uppercase">Speed</span>
            </div>
            <p className="font-display text-xl text-primary">{speed}<span className="text-sm">ft</span></p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
