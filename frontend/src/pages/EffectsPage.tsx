import { useState, useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { getEffects } from '@/lib/api';
import { ApiEffect } from '@/types/game';
import { Badge, badgeVariants } from '@/components/ui/badge';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { useAuth } from '@/context/AuthContext';
import { Zap, Search, AlertCircle, Info, BookText, ChevronRight, Layers, Loader2 } from 'lucide-react';

function CategoryBadge({ category }: { category: string }) {
  const normalized = category.toLowerCase();
  let themeVariant: "muted-subtle-outline" | "element-ice-outline" | "element-dark-outline" | "element-nature-outline" = "muted-subtle-outline"; // Default

  if (normalized === 'condition') themeVariant = 'element-ice-outline';
  else if (normalized === 'madness') themeVariant = 'element-dark-outline';
  else if (normalized === 'class feature') themeVariant = 'element-nature-outline';

  return (
    <Badge
      variant="outline"
      font="tomeMarginalia"
      size="sm" // text-[10px] is now size="sm"
      theme={themeVariant}
      className="uppercase tracking-wider" // uppercase tracking-wider remains in className
    >
      {category}
    </Badge>
  );
}

interface Variation {
  label: string;
  description: string;
  mechanics: string;
}

interface EffectGroup {
  groupName: string;
  category: string;
  description: string;
  count: number;
  variations: Variation[];
}

type DisplayItem =
  | { type: 'single'; effect: ApiEffect }
  | { type: 'group'; group: EffectGroup };

const GROUP_CONFIGS = [
  { prefix: 'Feral Mutation', groupName: 'Feral Mutation Table', description: 'Roll on the Mutagen Feral Table to determine your specific mutation effect while in Feral Mode.' },
  { prefix: 'Steampunk Madness', groupName: 'Steampunk Madness Table', description: 'Roll on the Lorewright Madness Table to determine your steampunk-related madness effect.' },
] as const;

function parseExhaustionLevels(mechanics: string): Variation[] {
  return mechanics.split(/(?=\d+:)/).filter(Boolean).map(part => {
    const match = part.match(/^(\d+):\s*(.+?)\.?\s*$/);
    if (!match) return null;
    return { label: `Level ${match[1]}`, description: match[2], mechanics: match[2] };
  }).filter(Boolean) as Variation[];
}

function extractRollNumber(name: string): number {
  const match = name.match(/\(Roll (\d+)\)/);
  return match ? parseInt(match[1]) : 0;
}

function EffectCard({ effect }: { effect: ApiEffect }) {
  return (
    <Card className="arcane-border h-full hover:bg-primary/5 transition-all hover:shadow-lg hover:shadow-primary/10 group overflow-hidden flex flex-col">
      <CardHeader className="pb-3 pt-4 space-y-1">
        <div className="flex items-start justify-between gap-2">
          <h3 className="font-tome-heading text-lg text-primary group-hover:glow-text transition-colors leading-tight">
            {effect.name}
          </h3>
          <CategoryBadge category={effect.category} />
        </div>
      </CardHeader>
      <CardContent className="space-y-4 flex-grow flex flex-col">
        <div className="text-sm space-y-3 flex-grow">
          <p className="text-muted-foreground leading-relaxed italic border-l-2 border-primary/20 pl-3">
            {effect.description}
          </p>

          <div className="bg-primary/5 p-3 rounded-lg border border-primary/10 text-xs mt-2 space-y-2">
            <span className="font-bold text-primary flex items-center gap-1.5 uppercase tracking-widest text-[10px]">
              <Info className="w-3 h-3" /> Mechanics
            </span>
            <p className="text-foreground/90 font-tome-marginalia leading-relaxed">
              {effect.mechanics}
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

function GroupedEffectCard({ group, onClick }: { group: EffectGroup; onClick: () => void }) {
  return (
    <Card
      className="arcane-border h-full hover:bg-primary/5 transition-all hover:shadow-lg hover:shadow-primary/10 group overflow-hidden flex flex-col cursor-pointer"
      onClick={onClick}
    >
      <CardHeader className="pb-3 pt-4 space-y-1">
        <div className="flex items-start justify-between gap-2">
          <h3 className="font-tome-heading text-lg text-primary group-hover:glow-text transition-colors leading-tight">
            {group.groupName}
          </h3>
          <CategoryBadge category={group.category} />
        </div>
      </CardHeader>
      <CardContent className="space-y-4 flex-grow flex flex-col">
        <div className="text-sm space-y-3 flex-grow">
          <p className="text-muted-foreground leading-relaxed italic border-l-2 border-primary/20 pl-3">
            {group.description}
          </p>
        </div>
        <div className="flex items-center justify-between text-xs text-primary/70 pt-2 border-t border-primary/10">
          <span className="flex items-center gap-1.5 font-tome-marginalia">
            <Layers className="w-3.5 h-3.5" />
            {group.count} variation{group.count !== 1 ? 's' : ''}
          </span>
          <span className="flex items-center gap-1 font-tome-marginalia group-hover:text-primary transition-colors">
            View all <ChevronRight className="w-3.5 h-3.5" />
          </span>
        </div>
      </CardContent>
    </Card>
  );
}

function VariationsDialog({ group, open, onOpenChange }: { group: EffectGroup | null; open: boolean; onOpenChange: (open: boolean) => void }) {
  if (!group) return null;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent tome className="max-w-2xl max-h-[80vh] flex flex-col">
        <DialogHeader>
          <DialogTitle className="font-tome-heading text-2xl text-primary">
            {group.groupName}
          </DialogTitle>
          <DialogDescription className="italic">
            {group.description}
          </DialogDescription>
        </DialogHeader>
        <div className="overflow-y-auto space-y-3 pr-2 -mr-2">
          {group.variations.map((v, i) => (
            <div key={i} className="bg-primary/5 p-3 rounded-lg border border-primary/10 space-y-1">
              <span className="font-tome-heading text-sm text-primary">
                {v.label}
              </span>
              <p className="text-sm text-foreground/90 font-tome-marginalia leading-relaxed">
                {v.description}
              </p>
            </div>
          ))}
        </div>
      </DialogContent>
    </Dialog>
  );
}

export function EffectsTabContent() {
  const { token } = useAuth();

  const [search, setSearch] = useState('');
  const [category, setCategory] = useState('all');
  const [selectedGroup, setSelectedGroup] = useState<EffectGroup | null>(null);

  const { data: effects, isLoading } = useQuery({
    queryKey: ['effects'],
    queryFn: () => getEffects(token ?? undefined),
    staleTime: 60_000 * 10,
  });

  const categories = useMemo(() => {
    if (!effects) return [];
    const cats = new Set(effects.map(e => e.category));
    return Array.from(cats).sort();
  }, [effects]);

  const displayItems = useMemo((): DisplayItem[] => {
    if (!effects) return [];
    const searchLower = search.toLowerCase();

    const grouped = new Map<string, ApiEffect[]>();
    const singles: ApiEffect[] = [];

    for (const effect of effects) {
      const config = GROUP_CONFIGS.find(c => effect.name.startsWith(c.prefix + ' ('));
      if (config) {
        const existing = grouped.get(config.prefix) ?? [];
        existing.push(effect);
        grouped.set(config.prefix, existing);
      } else {
        singles.push(effect);
      }
    }

    const items: DisplayItem[] = [];

    for (const effect of singles) {
      const matchesSearch = !searchLower ||
        effect.name.toLowerCase().includes(searchLower) ||
        effect.description.toLowerCase().includes(searchLower) ||
        effect.mechanics.toLowerCase().includes(searchLower);
      const matchesCategory = category === 'all' || effect.category === category;
      if (!matchesSearch || !matchesCategory) continue;

      if (effect.name === 'Exhaustion') {
        const levels = parseExhaustionLevels(effect.mechanics);
        const group: EffectGroup = {
          groupName: 'Exhaustion',
          category: effect.category,
          description: effect.description,
          count: levels.length,
          variations: levels,
        };
        items.push({ type: 'group', group });
      } else {
        items.push({ type: 'single', effect });
      }
    }

    for (const [prefix, groupEffects] of grouped) {
      const config = GROUP_CONFIGS.find(c => c.prefix === prefix)!;
      const matchesCategory = category === 'all' || groupEffects.some(e => e.category === category);
      const matchesSearch = !searchLower ||
        config.groupName.toLowerCase().includes(searchLower) ||
        config.description.toLowerCase().includes(searchLower) ||
        groupEffects.some(e =>
          e.name.toLowerCase().includes(searchLower) ||
          e.description.toLowerCase().includes(searchLower) ||
          e.mechanics.toLowerCase().includes(searchLower)
        );
      if (!matchesSearch || !matchesCategory) continue;

      const sorted = [...groupEffects].sort((a, b) => extractRollNumber(a.name) - extractRollNumber(b.name));
      const group: EffectGroup = {
        groupName: config.groupName,
        category: sorted[0].category,
        description: config.description,
        count: sorted.length,
        variations: sorted.map(e => ({
          label: `Roll ${extractRollNumber(e.name)}`,
          description: e.description,
          mechanics: e.mechanics,
        })),
      };
      items.push({ type: 'group', group });
    }

    return items.sort((a, b) => {
      const nameA = a.type === 'single' ? a.effect.name : a.group.groupName;
      const nameB = b.type === 'single' ? b.effect.name : b.group.groupName;
      return nameA.localeCompare(nameB);
    });
  }, [effects, search, category]);

  return (
    <div className="w-full space-y-8 p-6">
      <div className="flex flex-col xl:flex-row xl:items-center justify-between gap-4">
        <div className="flex items-center gap-4">
          <div className="p-3 rounded-full border-2 border-faded-gold/50 bg-primary/10">
            <Zap className="w-6 h-6 text-primary" />
          </div>
          <div>
            <h1 className="font-tome-heading text-3xl text-primary glow-text">
              Status & Conditions
            </h1>
            <p className="text-muted-foreground text-sm font-tome-marginalia mt-1">
              Rules for madness, exhaustion, and magical afflictions
            </p>
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-3">
          <div className="relative w-full md:w-64">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
            <Input
              placeholder="Search effects..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9"
            />
          </div>
          <Select value={category} onValueChange={setCategory}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Category" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Categories</SelectItem>
              {categories.map(cat => (
                <SelectItem key={cat} value={cat}>{cat}</SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      <div className="arcane-border rounded-lg p-4 bg-primary/5 flex items-center gap-3">
        <BookText className="w-5 h-5 text-primary flex-shrink-0" />
        <p className="text-sm text-muted-foreground italic">
          "The mind is a fragile vessel; once cracked, the echoes of the void begin to seep in."
          <span className="ml-2 not-italic font-tome-marginalia text-primary/60">— Codex of the Rift</span>
        </p>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center py-16">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
        </div>
      ) : displayItems.length > 0 ? (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {displayItems.map((item) =>
            item.type === 'single' ? (
              <EffectCard key={item.effect.id} effect={item.effect} />
            ) : (
              <GroupedEffectCard
                key={item.group.groupName}
                group={item.group}
                onClick={() => setSelectedGroup(item.group)}
              />
            )
          )}
        </div>
      ) : (
        <div className="arcane-border rounded-xl p-12 text-center">
          <AlertCircle className="w-16 h-16 mx-auto mb-6 text-muted-foreground" />
          <h2 className="font-tome-heading text-2xl text-primary mb-2">
            No Effects Found
          </h2>
          <p className="text-muted-foreground font-tome-marginalia">
            Try adjusting your search criteria.
          </p>
        </div>
      )}

      <VariationsDialog
        group={selectedGroup}
        open={!!selectedGroup}
        onOpenChange={(open) => { if (!open) setSelectedGroup(null); }}
      />
    </div>
  );
}
