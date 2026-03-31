import { Users } from 'lucide-react';
import { ApiLevelUpPreview } from '@/types/game';
import { cn } from '@/lib/utils';

interface ArchetypeSelectionProps {
  preview: ApiLevelUpPreview;
  selectedArchetypeId: string | null;
  onSelect: (id: string) => void;
}

export function ArchetypeSelection({ preview, selectedArchetypeId, onSelect }: ArchetypeSelectionProps) {
  if (!preview.available_archetypes) return null;

  return (
    <div className="space-y-4">
      <h3 className="font-tome-subheading text-lg text-primary flex items-center gap-2">
        <Users className="h-5 w-5" />
        Choose Your Path
      </h3>
      <p className="text-sm text-muted-foreground font-tome-marginalia">
        Select your archetype. This choice will determine your specialized abilities.
      </p>
      <div className="space-y-3">
        {preview.available_archetypes.map((archetype) => (
          <button
            key={archetype.id}
            type="button"
            onClick={() => onSelect(archetype.id)}
            className={cn(
              'w-full p-4 rounded border text-left transition-all',
              selectedArchetypeId === archetype.id
                ? 'border-primary bg-primary/10 ring-2 ring-primary'
                : 'border-border/50 bg-muted/30 hover:border-primary/50'
            )}
          >
            <p className="font-tome-subheading text-foreground mb-2">{archetype.name}</p>
            <p className="text-sm text-muted-foreground">{archetype.description}</p>
          </button>
        ))}
      </div>
    </div>
  );
}
