import { Heart, Users, Swords } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { ApiLevelUpPreview, ApiArchetype, ApiWeapon } from '@/types/game';

interface LevelUpConfirmProps {
  preview: ApiLevelUpPreview;
  hpGain: number;
  hpChoice: 'average' | 'roll';
  hpRollResult: number | null;
  selectedArchetype: ApiArchetype | null;
  selectedWeapon: ApiWeapon | null;
  asiAllocation: Record<string, number>;
}

const MODIFIER_DISPLAY_NAMES: Record<string, string> = {
  piston_core: 'Piston Core',
  venom_coating: 'Venom Coating',
  shadow_coating: 'Shadow Coating',
  lethal_coating: 'Lethal Coating',
};

export function LevelUpConfirm({
  preview,
  hpGain,
  hpChoice,
  hpRollResult,
  selectedArchetype,
  selectedWeapon,
  asiAllocation,
}: LevelUpConfirmProps) {
  return (
    <div className="space-y-4">
      <h3 className="font-tome-subheading text-lg text-primary">Confirm Level Up</h3>
      <div className="space-y-3">
        <p className="text-muted-foreground font-tome-marginalia">
          You are about to level up to <strong className="text-primary">Level {preview.next_level}</strong>
        </p>

        {/* HP Gain Summary */}
        <div className="p-3 rounded border border-green-500/30 bg-green-500/10">
          <p className="font-semibold mb-2 text-foreground flex items-center gap-2">
            <Heart className="h-4 w-4 text-green-400" />
            Hit Point Increase
          </p>
          <div className="flex items-center gap-2">
            <Badge variant="element-heal">
              +{hpGain} HP
            </Badge>
            <span className="text-sm text-muted-foreground">
              ({hpChoice === 'average' ? 'Average' : `Rolled: ${hpRollResult}`} + CON)
            </span>
          </div>
          <p className="text-sm text-muted-foreground mt-1">
            New Max HP: <strong className="text-primary">{preview.current_max_hp + hpGain}</strong>
          </p>
        </div>

        {selectedArchetype && (
          <div className="p-3 rounded border border-primary/30 bg-primary/10">
            <p className="font-semibold mb-2 text-foreground flex items-center gap-2">
              <Users className="h-4 w-4 text-primary" />
              Archetype Selected
            </p>
            <Badge variant="default">{selectedArchetype.name}</Badge>
          </div>
        )}

        {selectedWeapon && preview.weapon_selection_info && (
          <div className="p-3 rounded border border-amber-500/30 bg-amber-500/10">
            <p className="font-semibold mb-2 text-foreground flex items-center gap-2">
              <Swords className="h-4 w-4 text-amber-400" />
              Signature Weapon
            </p>
            <Badge variant="element-lightning">{selectedWeapon.name}</Badge>
            <p className="text-sm text-muted-foreground mt-1">
              {MODIFIER_DISPLAY_NAMES[preview.weapon_selection_info.modifier_type] || preview.weapon_selection_info.modifier_type} will be applied
            </p>
          </div>
        )}

        {Object.keys(asiAllocation).filter(k => asiAllocation[k] > 0).length > 0 && (
          <div className="p-3 rounded border border-border/50 bg-muted/30">
            <p className="font-semibold mb-2 text-foreground">Ability Score Improvements:</p>
            <div className="flex flex-wrap gap-2">
              {Object.entries(asiAllocation)
                .filter(([, v]) => v > 0)
                .map(([ability, points]) => (
                  <Badge key={ability} variant="secondary">
                    {ability}: +{points}
                  </Badge>
                ))}
            </div>
          </div>
        )}
        {preview.class_level.level_features && preview.class_level.level_features.length > 0 && (
          <div className="p-3 rounded border border-border/50 bg-muted/30">
            <p className="font-semibold mb-2 text-foreground">Features Gained:</p>
            <div className="flex flex-wrap gap-2">
              {preview.class_level.level_features.map((f) => (
                <Badge key={f.id} variant="outline">{f.name}</Badge>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
