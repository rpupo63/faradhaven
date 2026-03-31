import { useQuery } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { History, Star, Sparkles, BookOpen } from 'lucide-react';
import { getLevelHistory, type LevelUpHistoryEntry } from '@/lib/api';
import { DND5E_SKILLS } from '@/types/game';

interface LevelHistoryProps {
  characterId: string;
  token?: string;
}

export function LevelHistory({ characterId, token }: LevelHistoryProps) {
  const { data: history, isLoading } = useQuery({
    queryKey: ['levelHistory', characterId],
    queryFn: () => getLevelHistory(characterId, token),
  });

  if (isLoading) {
    return (
      <Card className="arcane-border">
        <CardContent className="p-6 text-center text-muted-foreground font-tome-marginalia">
          Loading history...
        </CardContent>
      </Card>
    );
  }

  if (!history || history.length === 0) {
    return (
      <Card className="arcane-border">
        <CardContent className="p-6 text-center text-muted-foreground">
          <History className="h-8 w-8 mx-auto mb-2 opacity-50" />
          <p className="font-tome-marginalia">No level-up history recorded</p>
          <p className="text-sm">Level up your character to see changes tracked here</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="arcane-border bg-card">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 font-tome-subheading text-primary">
          <History className="h-5 w-5" />
          Level History
        </CardTitle>
      </CardHeader>
      <CardContent>
        <ScrollArea className="h-[400px] pr-4">
          <div className="space-y-4">
            {history.map((entry) => (
              <LevelHistoryCard key={entry.id} entry={entry} />
            ))}
          </div>
        </ScrollArea>
      </CardContent>
    </Card>
  );
}

function LevelHistoryCard({ entry }: { entry: LevelUpHistoryEntry }) {
  const hasASI = entry.asi_allocation && Object.keys(entry.asi_allocation).length > 0;
  const hasSkills = entry.skill_selections && entry.skill_selections.length > 0;
  const hasSpells = entry.spells_learned && entry.spells_learned.length > 0;
  const hasFeatures = entry.features_gained && entry.features_gained.length > 0;

  return (
    <div className="p-4 rounded-lg border border-border/50 bg-muted/30">
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          <Star className="h-4 w-4 text-primary" />
          <span className="font-tome-subheading text-foreground">Level {entry.level}</span>
        </div>
        <span className="text-xs text-muted-foreground font-tome-marginalia">
          {new Date(entry.created_at).toLocaleDateString()}
        </span>
      </div>

      <div className="space-y-2 text-sm">
        {hasASI && (
          <div className="flex items-start gap-2">
            <Sparkles className="h-4 w-4 text-primary mt-0.5" />
            <div>
              <span className="text-muted-foreground">ASI: </span>
              {Object.entries(entry.asi_allocation!)
                .filter(([, v]) => v > 0)
                .map(([ability, points]) => (
                  <Badge key={ability} variant="secondary" className="mr-1 text-xs">
                    {ability} +{points}
                  </Badge>
                ))}
            </div>
          </div>
        )}

        {hasSkills && (
          <div className="flex items-start gap-2">
            <BookOpen className="h-4 w-4 text-primary mt-0.5" />
            <div>
              <span className="text-muted-foreground">Skills: </span>
              {entry.skill_selections!.map((skillId) => {
                const skill = DND5E_SKILLS.find((s) => s.id === skillId);
                return (
                  <Badge key={skillId} variant="outline" className="mr-1 text-xs">
                    {skill?.name ?? skillId}
                  </Badge>
                );
              })}
            </div>
          </div>
        )}

        {hasSpells && (
          <div className="flex items-start gap-2">
            <Sparkles className="h-4 w-4 text-primary mt-0.5" />
            <div>
              <span className="text-muted-foreground">Spells learned: </span>
              <span className="text-foreground">{entry.spells_learned!.length}</span>
            </div>
          </div>
        )}

        {hasFeatures && (
          <div className="flex items-start gap-2">
            <Star className="h-4 w-4 text-primary mt-0.5" />
            <div>
              <span className="text-muted-foreground">Features: </span>
              {entry.features_gained!.map((feature) => (
                <Badge key={feature} variant="outline" className="mr-1 text-xs">
                  {feature}
                </Badge>
              ))}
            </div>
          </div>
        )}

        {!hasASI && !hasSkills && !hasSpells && !hasFeatures && (
          <p className="text-muted-foreground italic font-tome-marginalia">
            No choices made at this level
          </p>
        )}
      </div>
    </div>
  );
}
