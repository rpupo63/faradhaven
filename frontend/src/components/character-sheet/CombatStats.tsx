import { Heart, Dices, Shield, Footprints, Pencil, Zap, ShieldCheck, GraduationCap } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Card, CardContent } from '@/components/ui/card';
import { NormalizedCharacterSheet } from '@/types/game';

interface CombatStatsProps {
  sheet: NormalizedCharacterSheet;
  expandedPanel: 'hp' | 'hitdice' | null;
  setExpandedPanel: (panel: 'hp' | 'hitdice' | null) => void;
  onHPChange?: (delta: number) => void | Promise<void>;
  onUseHitDice?: (rolls: number[]) => void | Promise<void>;
  onRoll?: (label: string, modifier: number) => void;
}

const formatMod = (n: number) => (n >= 0 ? `+${n}` : `${n}`);

export function CombatStats({
  sheet,
  expandedPanel,
  setExpandedPanel,
  onHPChange,
  onUseHitDice,
  onRoll,
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
    modifiers,
    save_dc: saveDC,
    class_level,
  } = sheet;

  const pb = class_level?.proficiency_bonus ?? 2;

  const hasUnarmoredDefense = class_level?.level_features?.some(f => f.name === 'Unarmored Defense');
  const hasArmorEquipped = !!sheet.character.equipped_armor_id;
  const hasShieldEquipped = !!sheet.character.equipped_shield_id;
  const equippedArmorItem = hasArmorEquipped
    ? sheet.inventory_items?.find(it => it.id === sheet.character.equipped_armor_id)
    : undefined;

  return (
    <div className="grid grid-cols-3 gap-1 md:gap-3">
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
        <Pencil className="absolute top-1 right-1 h-2.5 w-2.5 text-muted-foreground/30 md:top-2 md:right-2 md:h-3 md:w-3" />
        <CardContent className="p-1.5 md:p-4 text-center h-full flex flex-col justify-center min-h-[60px] md:min-h-0">
          <div className="flex items-center justify-center gap-1 mb-0.5 md:mb-1">
            <Heart className="h-3 w-3 md:h-4 md:w-4 text-red-500" />
            <span className="text-[10px] md:text-xs font-tome-marginalia text-muted-foreground uppercase leading-none">HP</span>
          </div>
          <p className="font-display text-lg md:text-2xl text-primary leading-tight">
            {currentHP}<span className="text-xs md:text-sm text-muted-foreground">/{maxHP}</span>
          </p>
          {tempHP > 0 && (
            <p className="text-[10px] md:text-xs text-blue-500 font-tome-marginalia leading-none">+{tempHP}</p>
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
        <Pencil className="absolute top-1 right-1 h-2.5 w-2.5 text-muted-foreground/30 md:top-2 md:right-2 md:h-3 md:w-3" />
        <CardContent className="p-1.5 md:p-4 text-center h-full flex flex-col justify-center min-h-[60px] md:min-h-0">
          <div className="flex items-center justify-center gap-1 mb-0.5 md:mb-1">
            <Dices className="h-3 w-3 md:h-4 md:w-4 text-primary" />
            <span className="text-[10px] md:text-xs font-tome-marginalia text-muted-foreground uppercase leading-none">Hit Dice</span>
          </div>
          <p className="font-display text-lg md:text-2xl text-primary leading-tight">
            {hitDiceRemaining}<span className="text-xs md:text-sm text-muted-foreground">/{hitDiceTotal}</span>
          </p>
          <p className="text-[10px] md:text-xs text-muted-foreground font-tome-marginalia leading-none">d{hitDie}</p>
        </CardContent>
      </Card>

      {/* AC Card */}
      <Card className="arcane-border bg-card">
        <CardContent className="p-1.5 md:p-4 text-center h-full flex flex-col justify-center min-h-[60px] md:min-h-0">
          <div className="flex items-center justify-center gap-1 mb-0.5 md:mb-1">
            <Shield className="h-3 w-3 md:h-4 md:w-4 text-primary" />
            <span className="text-[10px] md:text-xs font-tome-marginalia text-muted-foreground uppercase leading-none">AC</span>
          </div>
          <p className="font-display text-lg md:text-2xl text-primary leading-tight">{ac}</p>
          <p className="hidden md:block text-micro text-muted-foreground font-tome-marginalia leading-tight mt-1">
            {equippedArmorItem?.name || (hasUnarmoredDefense ? 'Unarmored' : 'Base')}
            {hasShieldEquipped && ' + Shield'}
          </p>
        </CardContent>
      </Card>

      {/* Speed */}
      {speed != null && (
        <Card className="arcane-border bg-card">
          <CardContent className="p-1.5 md:p-4 text-center h-full flex flex-col justify-center min-h-[60px] md:min-h-0">
            <div className="flex items-center justify-center gap-1 mb-0.5 md:mb-1">
              <Footprints className="h-3 w-3 md:h-4 md:w-4 text-primary" />
              <span className="text-[10px] md:text-xs font-tome-marginalia text-muted-foreground uppercase leading-none">Speed</span>
            </div>
            <p className="font-display text-lg md:text-xl text-primary leading-tight">{speed}<span className="text-xs">ft</span></p>
          </CardContent>
        </Card>
      )}

      {/* Initiative */}
      <Card
        className="arcane-border bg-card cursor-pointer hover:bg-primary/5 transition-colors"
        onClick={() => onRoll?.('Initiative', modifiers.initiative)}
      >
        <CardContent className="p-1.5 md:p-4 text-center h-full flex flex-col justify-center min-h-[60px] md:min-h-0">
          <div className="flex items-center justify-center gap-1 mb-0.5 md:mb-1">
            <Zap className="h-3 w-3 md:h-4 md:w-4 text-primary" />
            <span className="text-[10px] md:text-xs font-tome-marginalia text-muted-foreground uppercase leading-none">Init</span>
          </div>
          <p className="font-display text-lg md:text-xl text-primary leading-tight">{formatMod(modifiers.initiative)}</p>
        </CardContent>
      </Card>

      {/* Proficiency Bonus */}
      <Card className="arcane-border bg-card">
        <CardContent className="p-1.5 md:p-4 text-center h-full flex flex-col justify-center min-h-[60px] md:min-h-0">
          <div className="flex items-center justify-center gap-1 mb-0.5 md:mb-1">
            <GraduationCap className="h-3 w-3 md:h-4 md:w-4 text-primary" />
            <span className="text-[10px] md:text-xs font-tome-marginalia text-muted-foreground uppercase leading-none">Prof</span>
          </div>
          <p className="font-display text-lg md:text-xl text-primary leading-tight">{formatMod(pb)}</p>
        </CardContent>
      </Card>

      {/* Save DC */}
      {saveDC > 0 && (
        <Card className="arcane-border bg-card">
          <CardContent className="p-1.5 md:p-4 text-center h-full flex flex-col justify-center min-h-[60px] md:min-h-0">
            <div className="flex items-center justify-center gap-1 mb-0.5 md:mb-1">
              <ShieldCheck className="h-3 w-3 md:h-4 md:w-4 text-primary" />
              <span className="text-[10px] md:text-xs font-tome-marginalia text-muted-foreground uppercase leading-none">Save DC</span>
            </div>
            <p className="font-display text-lg md:text-xl text-primary leading-tight">{saveDC}</p>
          </CardContent>
        </Card>
      )}

      {/* Sneak Attack Dice */}
      {(class_level?.sneak_attack_dice ?? 0) > 0 && (
        <Card className="arcane-border bg-card">
          <CardContent className="p-1.5 md:p-4 text-center h-full flex flex-col justify-center min-h-[60px] md:min-h-0">
            <div className="flex items-center justify-center gap-1 mb-0.5 md:mb-1">
              <span className="text-[10px] md:text-xs font-tome-marginalia text-muted-foreground uppercase leading-none truncate w-full">Sneak</span>
            </div>
            <p className="text-base md:text-xl font-display text-primary leading-tight">{class_level.sneak_attack_dice}d6</p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
