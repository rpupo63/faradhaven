import { useMemo, useState } from 'react';
import { ScrollText } from 'lucide-react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { TomePage } from '@/components/TomePage';
import { LoreMarkdown } from '@/components/LoreMarkdown';
import { getLoreEntries } from '@/content/lore/loreNotes';
import { cn } from '@/lib/utils';

export default function LorePage() {
  const entries = useMemo(() => getLoreEntries(), []);
  const [tab, setTab] = useState(entries[0]?.id ?? '');

  if (entries.length === 0) {
    return (
      <div className="space-y-4 animate-ink-bleed">
        <h1 className="flex items-center gap-3 font-tome-heading text-3xl text-primary">
          <ScrollText className="w-8 h-8 shrink-0" />
          Lore
        </h1>
        <p className="text-muted-foreground font-tome-marginalia">
          No lore notes found. Add markdown files under <code className="text-foreground/80">lore-notes/</code> in the
          repository root.
        </p>
      </div>
    );
  }

  return (
    <div className="w-full space-y-8 animate-ink-bleed">
      <div>
        <h1 className="flex items-center gap-3 font-tome-heading text-3xl text-primary">
          <ScrollText className="w-8 h-8 shrink-0" />
          Lore
        </h1>
        <p className="mt-1 text-muted-foreground font-tome-marginalia">
          One entry per note in <code className="text-foreground/80">lore-notes/</code>—edit those files to update the site.
        </p>
      </div>

      <Tabs value={tab} onValueChange={setTab} className="w-full">
        <TabsList
          className={cn(
            'mb-6 grid h-auto w-full gap-1 p-1 arcane-border bg-card',
            entries.length === 1 && 'grid-cols-1',
            entries.length === 2 && 'grid-cols-2',
            entries.length >= 3 && 'grid-cols-2 sm:grid-cols-3 lg:grid-cols-4'
          )}
        >
          {entries.map((e) => (
            <TabsTrigger
              key={e.id}
              value={e.id}
              className="min-h-10 font-tome-subheading uppercase tracking-wide text-xs sm:text-sm whitespace-normal text-center leading-tight px-2 py-2"
            >
              {e.label}
            </TabsTrigger>
          ))}
        </TabsList>

        {entries.map((e) => (
          <TabsContent key={e.id} value={e.id} className="mt-0">
            <TomePage prose deckled className="text-foreground/95">
              <LoreMarkdown source={e.markdown} />
            </TomePage>
          </TabsContent>
        ))}
      </Tabs>
    </div>
  );
}
