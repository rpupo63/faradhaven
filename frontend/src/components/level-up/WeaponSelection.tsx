import { Swords, Star } from 'lucide-react';
import { ApiLevelUpPreview } from '@/types/game';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';

interface WeaponSelectionProps {
  preview: ApiLevelUpPreview;
  selectedWeaponId: string | null;
  onSelect: (id: string) => void;
}

const MODIFIER_DISPLAY_NAMES: Record<string, string> = {
  piston_core: 'Piston Core',
  venom_coating: 'Venom Coating',
  shadow_coating: 'Shadow Coating',
  lethal_coating: 'Lethal Coating',
};

export function WeaponSelection({ preview, selectedWeaponId, onSelect }: WeaponSelectionProps) {
  const weaponInfo = preview.weapon_selection_info;
  if (!weaponInfo) return null;

  const modifierName = MODIFIER_DISPLAY_NAMES[weaponInfo.modifier_type] || weaponInfo.modifier_type;

  return (
    <div className="space-y-4">
      <h3 className="font-tome-subheading text-lg text-primary flex items-center gap-2">
        <Swords className="h-5 w-5" />
        Select Your Signature Weapon
      </h3>
      <p className="text-sm text-muted-foreground font-tome-marginalia">
        {weaponInfo.description}
      </p>
      <div className="p-3 rounded border border-primary/30 bg-primary/5">
        <p className="text-xs text-muted-foreground">
          <Star className="h-3 w-3 inline mr-1" />
          The <span className="font-semibold text-primary">{modifierName}</span> will be permanently bonded to your chosen weapon.
        </p>
      </div>
      <div className="space-y-2 max-h-80 overflow-y-auto">
        {weaponInfo.eligible_weapons.map((weapon) => (
          <button
            key={weapon.id}
            type="button"
            onClick={() => onSelect(weapon.id)}
            className={cn(
              'w-full p-4 rounded border text-left transition-all',
              selectedWeaponId === weapon.id
                ? 'border-primary bg-primary/10 ring-2 ring-primary'
                : 'border-border/50 bg-muted/30 hover:border-primary/50'
            )}
          >
            <div className="flex justify-between items-start mb-2">
              <p className="font-tome-subheading text-foreground">{weapon.name}</p>
              <span className="text-micro text-muted-foreground">{weapon.range_type}</span>
            </div>
            <div className="flex flex-wrap gap-2 mb-2">
              {weapon.damages?.map((d, i) => (
<Badge key={i} variant="secondary" size="sm">
                  {d.damage_dice} {d.damage_type}
                </Badge>
              ))}
              {weapon.properties?.map((prop, i) => (
                <Badge key={`prop-${i}`} variant="outline" className="text-micro h-4 py-0">
                  {prop}
                </Badge>
              ))}
            </div>
            {weapon.description && (
              <p className="text-xs text-muted-foreground line-clamp-2">{weapon.description}</p>
            )}
          </button>
        ))}
      </div>
    </div>
  );
}
