import { ScrollArea } from '@/components/ui/scroll-area';
import { Shield, Heart } from 'lucide-react';
import { ApiClass, ApiCreationOptions } from '@/types/game';
import { cn } from '@/lib/utils';

interface ClassSelectionProps {
  options: ApiCreationOptions;
  selected: ApiClass | null;
  onSelect: (c: ApiClass) => void;
}

export function ClassSelection({ options, selected, onSelect }: ClassSelectionProps) {
  const classes = options?.classes ?? [];
  
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 h-full">
      <ScrollArea className="md:col-span-1 border-r pr-4 h-[500px]">
        <div className="space-y-2">
          {classes.map((cls) => (
            <div
              key={cls.id}
              className={cn(
                "p-3 rounded-lg border cursor-pointer transition-colors text-left text-sm",
                selected?.id === cls.id
                  ? "bg-primary/10 border-primary text-primary"
                  : "bg-card border-border hover:bg-muted"
              )}
              onClick={() => onSelect(cls)}
            >
              <div className="font-display font-bold">{cls.name}</div>
              <div className="text-xs text-muted-foreground truncate">Hit Die: d{cls.hit_die}</div>
            </div>
          ))}
        </div>
      </ScrollArea>
      
      <div className="md:col-span-2 pl-2">
        {selected ? (
          <ScrollArea className="h-[500px] pr-4">
            <h2 className="font-display text-3xl text-primary mb-2">{selected.name}</h2>
            <div className="flex gap-4 mb-4 text-sm text-muted-foreground">
              <div className="flex items-center gap-1"><Heart className="w-4 h-4" /> d{selected.hit_die} Hit Die</div>
              <div className="flex items-center gap-1"><Shield className="w-4 h-4" /> {selected.saving_throws?.join(", ")} Saves</div>
            </div>
            
            <p className="text-muted-foreground font-tome-text mb-4 whitespace-pre-wrap">{selected.description}</p>
            
            <div className="grid grid-cols-2 gap-4 mb-6">
              <div className="bg-muted/20 p-3 rounded border">
                <h4 className="font-bold text-sm mb-1 text-primary">Proficiencies</h4>
                <p className="text-xs text-muted-foreground">{selected.proficiencies || "None"}</p>
              </div>
              <div className="bg-muted/20 p-3 rounded border">
                <h4 className="font-bold text-sm mb-1 text-primary">Skills</h4>
                <p className="text-xs text-muted-foreground">
                  Choose {selected.skill_choice_count} from: {selected.skill_choice?.join(", ")}
                </p>
              </div>
            </div>

            <h3 className="font-display text-lg mb-2 border-b">Class Features</h3>
            <div className="space-y-4">
              {selected.levels?.find(l => l.level === 1)?.level_features?.map(feature => (
                <div key={feature.id} className="text-sm">
                  <div className="font-bold text-foreground">{feature.name}</div>
                  <div className="text-muted-foreground">{feature.description}</div>
                </div>
              ))}
            </div>
          </ScrollArea>
        ) : (
          <div className="flex h-full items-center justify-center text-muted-foreground">
            Select a class to view details
          </div>
        )}
      </div>
    </div>
  );
}
