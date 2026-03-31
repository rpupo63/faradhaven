import { ApiRace, ApiClass } from '@/types/game';

interface ReviewSummaryProps {
  name: string;
  race: ApiRace | null;
  cls: ApiClass | null;
  abilities: Record<string, number>;
  languages?: string[];
}

export function ReviewSummary({ name, race, cls, abilities, languages }: ReviewSummaryProps) {
  return (
    <div className="space-y-6">
      <div className="text-center">
        <h2 className="font-display text-4xl text-primary mb-1">{name}</h2>
        <p className="text-xl text-muted-foreground font-tome-marginalia">
          Level 1 {race?.name} {cls?.name}
        </p>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-6 gap-2 text-center border-y py-4">
        {Object.entries(abilities).map(([key, val]) => (
          <div key={key} className="p-2">
            <div className="text-xs uppercase text-muted-foreground font-bold">{key.slice(0,3)}</div>
            <div className="font-display text-2xl">{val}</div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="space-y-2">
          <h3 className="font-display text-lg border-b">Race Features</h3>
          <ul className="list-disc list-inside text-sm text-muted-foreground">
            {race?.traits?.map(t => (
              <li key={t.id}>{t.name}</li>
            ))}
          </ul>
        </div>
        <div className="space-y-2">
          <h3 className="font-display text-lg border-b">Class Features</h3>
          <ul className="list-disc list-inside text-sm text-muted-foreground">
            {cls?.levels?.find(l => l.level === 1)?.level_features?.map(f => (
              <li key={f.id}>{f.name}</li>
            ))}
          </ul>
        </div>
      </div>
      
      {languages && languages.length > 0 && (
        <div className="space-y-2">
          <h3 className="font-display text-lg border-b">Languages</h3>
          <p className="text-sm text-muted-foreground">{languages.join(', ')}</p>
        </div>
      )}

      <div className="bg-primary/5 p-4 rounded text-center text-sm text-muted-foreground">
        Ready to embark on your adventure?
      </div>
    </div>
  );
}
