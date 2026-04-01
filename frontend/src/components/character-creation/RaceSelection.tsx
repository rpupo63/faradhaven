import { ScrollArea } from '@/components/ui/scroll-area';
import { Label } from '@/components/ui/label';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { ApiRace, ApiCreationOptions } from '@/types/game';
import { cn } from '@/lib/utils';

interface RaceSelectionProps {
  options: ApiCreationOptions;
  selected: ApiRace | null;
  onSelect: (r: ApiRace) => void;
  lineageId: string | null;
  onSelectLineage: (id: string) => void;
}

export function RaceSelection({ options, selected, onSelect, lineageId, onSelectLineage }: RaceSelectionProps) {
  const races = options?.races ?? [];
  
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 h-full">
      <ScrollArea className="md:col-span-1 border-r pr-4 h-[500px]">
        <div className="space-y-2">
          {races.map((race) => (
            <div
              key={race.id}
              className={cn(
                "p-3 rounded-lg border cursor-pointer transition-colors text-left text-sm",
                selected?.id === race.id
                  ? "bg-primary/10 border-primary text-primary"
                  : "bg-card border-border hover:bg-muted"
              )}
              onClick={() => { onSelect(race); onSelectLineage(''); }}
            >
              <div className="font-display font-bold">{race.name}</div>
              <div className="text-xs text-muted-foreground truncate">{race.size} {race.creature_type}</div>
            </div>
          ))}
        </div>
      </ScrollArea>
      
      <div className="md:col-span-2 pl-2">
        {selected ? (
          <ScrollArea className="h-[500px] pr-4">
            <h2 className="font-display text-3xl text-primary mb-2">{selected.name}</h2>
            {selected.ability_score_bonuses && (
              <div className="flex gap-2 mb-4">
                {Object.entries(selected.ability_score_bonuses).map(([stat, bonus]) => (
                  <span key={stat} className="bg-primary/20 text-primary px-2 py-1 rounded text-xs font-mono uppercase">
                    {stat.slice(0,3)} +{bonus}
                  </span>
                ))}
              </div>
            )}
            <p className="text-muted-foreground font-tome-text mb-4 whitespace-pre-wrap">{selected.description}</p>
            
            <h3 className="font-display text-lg mb-2 border-b">Traits</h3>
            <div className="space-y-4">
              {selected.traits?.map(trait => (
                <div key={trait.id} className="text-sm">
                  <div className="font-bold text-foreground">{trait.name}</div>
                  <div className="text-muted-foreground">{trait.description}</div>
                  
                  {/* Lineage Selection (if trait has options and is likely a lineage choice) */}
                  {trait.options && trait.options.length > 0 && (trait.name.includes("Lineage") || trait.name.includes("Ancestry") || trait.name.includes("Legacy")) && (
                    <div className="mt-3 ml-2 border-l-2 border-primary/30 pl-3">
                      <Label className="mb-2 block text-xs uppercase text-primary">Select Lineage</Label>
                      <RadioGroup value={lineageId || ''} onValueChange={onSelectLineage}>
                        {trait.options.map(opt => (
                          <div key={opt.id} className="flex items-start space-x-2 mb-2">
                            <RadioGroupItem value={opt.id} id={opt.id} className="mt-1" />
                            <div className="grid gap-1.5 leading-none">
                              <Label htmlFor={opt.id} className="cursor-pointer font-bold">{opt.name}</Label>
                              <p className="text-xs text-muted-foreground">{opt.description}</p>
                              {opt.ability_score_bonuses && (
                                <div className="flex gap-1 mt-1">
                                  {Object.entries(opt.ability_score_bonuses).map(([s, b]) => (
                                    <span key={s} className="text-micro bg-muted px-1 rounded uppercase">
                                      {s.slice(0,3)} +{b}
                                    </span>
                                  ))}
                                </div>
                              )}
                            </div>
                          </div>
                        ))}
                      </RadioGroup>
                    </div>
                  )}
                </div>
              ))}
            </div>
          </ScrollArea>
        ) : (
          <div className="flex h-full items-center justify-center text-muted-foreground">
            Select a race to view details
          </div>
        )}
      </div>
    </div>
  );
}
