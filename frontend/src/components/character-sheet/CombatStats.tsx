import { Pencil } from 'lucide-react';
import { RaIcon } from '@/components/ui/RaIcon';
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

  const showSneakAttack = (class_level?.sneak_attack_dice ?? 0) > 0;

  return (
    <div className="flex flex-col gap-2">
      {/* Primary Stats: HP, AC, Hit Dice */}
      <div className="grid grid-cols-3 gap-2">
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
          <Pencil className="absolute top-1 right-1 h-2 w-2 text-muted-foreground/30" />
          <CardContent className="p-2 text-center h-full flex flex-col justify-center">
            <div className="flex items-center justify-center gap-1 mb-0.5">
            <RaIcon name="health" className="text-xs text-primary" />
            <span className="text-[10px] font-tome-marginalia text-muted-foreground uppercase leading-none">HP</span>
            </div>
            <p className="font-display text-xl text-primary leading-tight">
              {currentHP}<span className="text-xs text-muted-foreground">/{maxHP}</span>
            </p>
            {tempHP > 0 && (
              <p className="text-[10px] text-blue-500 font-tome-marginalia leading-none">+{tempHP}</p>
            )}
          </CardContent>
        </Card>

        {/* AC Card */}
        <Card className="arcane-border bg-card">
          <CardContent className="p-2 text-center h-full flex flex-col justify-center">
            <div className="flex items-center justify-center gap-1 mb-0.5">
              <RaIcon name="shield" className="text-xs text-primary" />
              <span className="text-[10px] font-tome-marginalia text-muted-foreground uppercase leading-none">AC</span>
            </div>
            <p className="font-display text-xl text-primary leading-tight">{ac}</p>
            {(equippedArmorItem || hasShieldEquipped || hasUnarmoredDefense) && (
              <p className="text-[8px] text-muted-foreground font-tome-marginalia leading-none truncate mt-0.5 px-1">
                {equippedArmorItem?.name || (hasUnarmoredDefense ? 'Unarmored' : 'Base')}
                {hasShieldEquipped && ' + Shield'}
              </p>
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
          <Pencil className="absolute top-1 right-1 h-2 w-2 text-muted-foreground/30" />
          <CardContent className="p-2 text-center h-full flex flex-col justify-center">
            <div className="flex items-center justify-center gap-1 mb-0.5">
              <RaIcon name="perspective-dice-six" className="text-xs text-primary" />
              <span className="text-[10px] font-tome-marginalia text-muted-foreground uppercase leading-none">Hit Dice</span>
            </div>
            <p className="font-display text-xl text-primary leading-tight">
              {hitDiceRemaining}<span className="text-xs text-muted-foreground">/{hitDiceTotal}</span>
            </p>
            <p className="text-[10px] text-muted-foreground font-tome-marginalia leading-none">d{hitDie}</p>
          </CardContent>
        </Card>
      </div>

      {/* Secondary Stats: Speed, Init, Prof, Save DC, Sneak Attack */}
      <div className={cn(
        "grid gap-2",
        showSneakAttack ? "grid-cols-3 md:grid-cols-5" : "grid-cols-2 md:grid-cols-4"
      )}>
        {/* Speed */}
        {speed != null && (
          <Card className="arcane-border bg-card">
            <CardContent className="p-1.5 text-center h-full flex flex-col justify-center">
              <div className="flex items-center justify-center gap-1 mb-0.5">
                <RaIcon name="shoe-prints" className="text-xs text-primary" />
                <span className="text-[10px] font-tome-marginalia text-muted-foreground uppercase leading-none">Speed</span>
              </div>
              <p className="font-display text-lg text-primary leading-tight">{speed}<span className="text-xs">ft</span></p>
            </CardContent>
          </Card>
        )}

        {/* Initiative */}
        <Card
          className="arcane-border bg-card cursor-pointer hover:bg-primary/5 transition-colors"
          onClick={() => onRoll?.('Initiative', modifiers.initiative)}
        >
          <CardContent className="p-1.5 text-center h-full flex flex-col justify-center">
            <div className="flex items-center justify-center gap-1 mb-0.5">
              <RaIcon name="lightning-bolt" className="text-xs text-primary" />
              <span className="text-[10px] font-tome-marginalia text-muted-foreground uppercase leading-none">Init</span>
            </div>
            <p className="font-display text-lg text-primary leading-tight">{formatMod(modifiers.initiative)}</p>
          </CardContent>
        </Card>

        {/* Proficiency Bonus */}
        <Card className="arcane-border bg-card">
          <CardContent className="p-1.5 text-center h-full flex flex-col justify-center">
            <div className="flex items-center justify-center gap-1 mb-0.5">
              <RaIcon name="trophy" className="text-xs text-primary" />
              <span className="text-[10px] font-tome-marginalia text-muted-foreground uppercase leading-none">Prof</span>
            </div>
            <p className="font-display text-lg text-primary leading-tight">{formatMod(pb)}</p>
          </CardContent>
        </Card>

        {/* Save DC */}
        {saveDC > 0 && (
          <Card className="arcane-border bg-card">
            <CardContent className="p-1.5 text-center h-full flex flex-col justify-center">
              <div className="flex items-center justify-center gap-1 mb-0.5">
                <RaIcon name="bolt-shield" className="text-xs text-primary" />
                <span className="text-[10px] font-tome-marginalia text-muted-foreground uppercase leading-none truncate">Save DC</span>
              </div>
              <p className="font-display text-lg text-primary leading-tight">{saveDC}</p>
            </CardContent>
          </Card>
        )}

        {/* Sneak Attack Dice */}
        {showSneakAttack && (
          <Card className="arcane-border bg-card">
            <CardContent className="p-1.5 text-center h-full flex flex-col justify-center">
              <div className="flex items-center justify-center gap-1 mb-0.5">
                <span className="text-[10px] font-tome-marginalia text-muted-foreground uppercase leading-none truncate">Sneak</span>
              </div>
              <p className="text-lg font-display text-primary leading-tight">{class_level!.sneak_attack_dice}d6</p>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
