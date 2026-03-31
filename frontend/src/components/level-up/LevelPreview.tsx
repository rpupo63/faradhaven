import { ApiLevelUpPreview } from '@/types/game';

interface LevelPreviewProps {
  preview: ApiLevelUpPreview;
}

export function LevelPreview({ preview }: LevelPreviewProps) {
  return (
    <div className="space-y-4">
      <h3 className="font-tome-subheading text-lg text-primary">Level {preview.next_level} Preview</h3>
      <div className="grid grid-cols-2 gap-4">
        <div className="p-4 rounded border border-border/50 bg-muted/30">
          <p className="text-sm text-muted-foreground font-tome-marginalia">Proficiency Bonus</p>
          <p className="font-display text-xl text-primary">+{preview.class_level.proficiency_bonus}</p>
        </div>
        <div className="p-4 rounded border border-border/50 bg-muted/30">
          <p className="text-sm text-muted-foreground font-tome-marginalia">Max Spell Points</p>
          <p className="font-display text-xl text-primary">{preview.class_level.max_spell_points}</p>
        </div>
        {preview.asi_points_available > 0 && (
          <div className="p-4 rounded border border-primary/30 bg-primary/10">
            <p className="text-sm text-muted-foreground font-tome-marginalia">ASI Points</p>
            <p className="font-display text-xl text-primary">+{preview.asi_points_available}</p>
          </div>
        )}
        {preview.new_spells_allowed > 0 && (
          <div className="p-4 rounded border border-primary/30 bg-primary/10">
            <p className="text-sm text-muted-foreground font-tome-marginalia">New Spells</p>
            <p className="font-display text-xl text-primary">+{preview.new_spells_allowed}</p>
          </div>
        )}
      </div>
      {preview.class_level.level_features && preview.class_level.level_features.length > 0 && (
        <div className="space-y-2">
          <p className="font-tome-subheading text-primary">New Features</p>
          {preview.class_level.level_features.map((feature) => (
            <div key={feature.id} className="p-3 rounded border border-border/50 bg-muted/30">
              <p className="font-semibold text-foreground">{feature.name}</p>
              <p className="text-sm text-muted-foreground">{feature.description}</p>
            </div>
          ))}
        </div>
      )}
      {preview.requires_archetype_choice && (
        <div className="p-4 rounded border border-primary/30 bg-primary/10">
          <p className="text-sm text-muted-foreground font-tome-marginalia">Archetype Choice</p>
          <p className="font-display text-lg text-primary">Choose your path at this level</p>
        </div>
      )}
    </div>
  );
}
