import React, { useState, useEffect } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { Loader2 } from 'lucide-react';
import { toast } from 'sonner';
import {
  generateLootPreview,
  confirmLootPickup,
  getLootOptions,
  type LootSource,
  type LootResult,
  type LootRoomTheme,
  type LootLocation,
} from '@/lib/api/loot';
import type { LootPreviewResponse, LootPartyMember } from '@/types/game/api';
import { Badge } from '@/components/ui/badge';
import {
  getThemesForSource,
  deriveLootLevelBand
} from '@/lib/lootThemes';

interface LootModalProps {
  isOpen: boolean;
  onClose: () => void;
  characterId: string;
  characterLevel: number;
  token: string;
}

const SOURCE_OPTIONS: { value: LootSource; label: string }[] = [
  { value: 'common_enemy', label: 'Common Enemy' },
  { value: 'boss_enemy', label: 'Boss Enemy' },
  { value: 'room', label: 'Room' },
];

const LOCATION_OPTIONS: { value: LootLocation; label: string }[] = [
  { value: 'indoor', label: 'Indoor' },
  { value: 'underground', label: 'Underground' },
  { value: 'urban', label: 'Urban' },
  { value: 'slums', label: 'Slums' },
  { value: 'estate', label: 'Estate' },
  { value: 'street', label: 'Street' },
  { value: 'wilds', label: 'Wilds' },
];

function formatCopper(cp: number): string {
  if (cp >= 100) {
    const gp = Math.floor(cp / 100);
    const remainder = cp % 100;
    return remainder > 0 ? `${gp} gp ${remainder} cp` : `${gp} gp`;
  }
  if (cp >= 10) {
    const sp = Math.floor(cp / 10);
    const remainder = cp % 10;
    return remainder > 0 ? `${sp} sp ${remainder} cp` : `${sp} sp`;
  }
  return `${cp} cp`;
}

function normalizeLootLevel(level: number): number {
  if (!Number.isFinite(level)) {
    return 1;
  }
  return Math.min(20, Math.max(1, Math.round(level)));
}

export const LootModal: React.FC<LootModalProps> = ({ isOpen, onClose, characterId, characterLevel, token }) => {
  const queryClient = useQueryClient();
  const { data: lootOptions } = useQuery({
    queryKey: ['loot-options'],
    queryFn: () => getLootOptions(),
    staleTime: 60_000,
  });

  const sourceOptions: { value: LootSource; label: string }[] =
    lootOptions?.sources?.map((source) => ({
      value: source,
      label: source.replace('_', ' ').replace(/\b\w/g, (c) => c.toUpperCase()),
    })) ?? SOURCE_OPTIONS;
  const lootLevelOptions = lootOptions?.loot_levels ?? Array.from({ length: 20 }, (_, i) => i + 1);

  const [source, setSource] = useState<LootSource>('common_enemy');
  const [roomTheme, setRoomTheme] = useState<LootRoomTheme>('dungeon');
  const [location, setLocation] = useState<LootLocation>('indoor');
  const [lootLevel, setLootLevel] = useState<number>(() => normalizeLootLevel(characterLevel));
  const [result, setResult] = useState<LootResult | null>(null);
  const [previewSessionId, setPreviewSessionId] = useState<string | null>(null);
  const [partyMembers, setPartyMembers] = useState<LootPartyMember[]>([]);
  const [dropAssignments, setDropAssignments] = useState<Record<number, string>>({});
  const levelBand = deriveLootLevelBand(lootLevel || characterLevel);
  const themeOptions = getThemesForSource(source).filter((option) =>
    lootOptions?.themes?.length ? lootOptions.themes.includes(option.value) : true
  );
  const locationOptions: { value: LootLocation; label: string }[] =
    lootOptions?.locations?.map((loc) => ({
      value: loc,
      label: loc.charAt(0).toUpperCase() + loc.slice(1),
    })) ?? LOCATION_OPTIONS;
  const selectedTheme = themeOptions.find((option) => option.value === roomTheme) ?? themeOptions[0];

  const lootMutation = useMutation<LootPreviewResponse, Error>({
    mutationFn: () => generateLootPreview(characterId, source, roomTheme, location, lootLevel, token),
    onSuccess: (data) => {
      setPreviewSessionId(data.session_id);
      setResult(data.loot);
      setPartyMembers(data.party_members);
      const defaultRecipient = data.party_members[0]?.id ?? characterId;
      const nextAssignments: Record<number, string> = {};
      (data.loot.drops ?? []).forEach((_, idx) => {
        nextAssignments[idx] = defaultRecipient;
      });
      setDropAssignments(nextAssignments);
      toast.success('Loot rolled. Assign drops to party members, then confirm pickup.');
    },
    onError: (error) => {
      toast.error(error.message || 'Failed to generate loot');
    },
  });
  const confirmMutation = useMutation<LootResult, Error>({
    mutationFn: async () => {
      if (!previewSessionId || !result) {
        throw new Error('No pending loot session to confirm');
      }
      const assignments = (result.drops ?? []).map((_, idx) => ({
        drop_index: idx,
        character_id: dropAssignments[idx] ?? characterId,
      }));
      return confirmLootPickup(characterId, previewSessionId, assignments, token);
    },
    onSuccess: (data) => {
      setResult(data);
      setPreviewSessionId(null);
      setDropAssignments({});
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      toast.success('Loot distributed to party inventory.');
    },
    onError: (error) => {
      toast.error(error.message || 'Failed to confirm loot pickup');
    },
  });

  useEffect(() => {
    if (themeOptions.length > 0 && !themeOptions.some((theme) => theme.value === roomTheme)) {
      setRoomTheme(themeOptions[0].value);
    }
  }, [roomTheme, themeOptions]);

  useEffect(() => {
    if (!isOpen) {
      setTimeout(() => setResult(null), 0);
      setPreviewSessionId(null);
      setDropAssignments({});
      setPartyMembers([]);
    }
  }, [isOpen]);

  useEffect(() => {
    if (isOpen) {
      setLootLevel(normalizeLootLevel(characterLevel));
    }
  }, [isOpen, characterLevel]);

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Generate Loot</DialogTitle>
          <DialogDescription>
            Build a themed loot profile with source, room style, reward intensity, and level context.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label>1. Source</Label>
            <Select value={source} onValueChange={(v) => setSource(v as LootSource)}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {sourceOptions.map((o) => (
                  <SelectItem key={o.value} value={o.value}>
                    {o.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="grid gap-2">
            <Label>2. Theme</Label>
            <Select value={roomTheme} onValueChange={(v) => setRoomTheme(v as LootRoomTheme)}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {themeOptions.map((o) => (
                  <SelectItem key={o.value} value={o.value}>
                    {o.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {selectedTheme && <p className="text-xs text-muted-foreground">{selectedTheme.description}</p>}
          </div>

          <div className="grid gap-2">
            <Label>3. Location</Label>
            <Select value={location} onValueChange={(v) => setLocation(v as LootLocation)}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {locationOptions.map((o) => (
                  <SelectItem key={o.value} value={o.value}>
                    {o.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="grid gap-2">
            <Label>4. Loot Level</Label>
            <Select value={String(lootLevel)} onValueChange={(v) => setLootLevel(Number(v))}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {lootLevelOptions.map((lvl) => (
                  <SelectItem key={lvl} value={String(lvl)}>
                    Level {lvl}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="flex flex-wrap gap-2">
            <Badge variant="secondary">Level Band: {levelBand}</Badge>
            <Badge variant="secondary">Theme: {roomTheme}</Badge>
            <Badge variant="secondary">Location: {location}</Badge>
            <Badge variant="secondary">Loot Level: {lootLevel}</Badge>
          </div>

          {result && (
            <div className="rounded-md border p-3 space-y-2 text-sm">
              <p className="font-medium">
                Gold earned: {formatCopper(result.gold_earned)}
              </p>
              <p className="text-muted-foreground">
                Item drops rolled: {result.items_rolled ?? result.items.length} · Weapon drops rolled:{' '}
                {result.weapons_rolled ?? result.weapons.length}
              </p>
              {(typeof result.expected_budget === 'number' || typeof result.session_budget === 'number') && (
                <p className="text-muted-foreground">
                  Budget: {result.session_budget ?? 0} cp rolled from expected {result.expected_budget ?? 0} cp ·
                  Ending: {result.ending_budget ?? 0} cp
                  {result.debt_used ? ` · Debt used (${result.debt_amount ?? 0} cp)` : ''}
                </p>
              )}
              {(result.profile_notes?.length ?? 0) > 0 && (
                <ul className="list-disc list-inside text-muted-foreground">
                  {result.profile_notes.map((note, idx) => (
                    <li key={`${note}-${idx}`}>{note}</li>
                  ))}
                </ul>
              )}
              {result.drops && result.drops.length > 0 ? (
                <div>
                  <p className="font-medium">Discovered (in order)</p>
                  <ol className="list-decimal list-inside space-y-1">
                    {result.drops.map((drop, i) => (
                      <li key={`${drop.kind}-${drop.name}-${i}`} className="space-y-1">
                        <span className="capitalize text-muted-foreground">{drop.kind}:</span> {drop.name}{' '}
                        <span className="text-muted-foreground">({drop.rarity})</span>
                        {previewSessionId && (
                          <div className="pt-1">
                            <Label className="text-xs">Assign to</Label>
                            <Select
                              value={dropAssignments[i] ?? partyMembers[0]?.id ?? characterId}
                              onValueChange={(v) => setDropAssignments((prev) => ({ ...prev, [i]: v }))}
                            >
                              <SelectTrigger className="h-8">
                                <SelectValue placeholder="Select party member" />
                              </SelectTrigger>
                              <SelectContent>
                                {partyMembers.map((member) => (
                                  <SelectItem key={member.id} value={member.id}>
                                    {member.name}
                                  </SelectItem>
                                ))}
                              </SelectContent>
                            </Select>
                          </div>
                        )}
                      </li>
                    ))}
                  </ol>
                </div>
              ) : (
                <>
                  {result.items.length > 0 && (
                    <div>
                      <p className="font-medium">Items:</p>
                      <ul className="list-disc list-inside">
                        {result.items.map((item, i) => (
                          <li key={`${item.id}-${i}`}>
                            {item.name}{' '}
                            <span className="text-muted-foreground">({item.rarity})</span>
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}
                  {result.weapons.length > 0 && (
                    <div>
                      <p className="font-medium">Weapons:</p>
                      <ul className="list-disc list-inside">
                        {result.weapons.map((w, i) => (
                          <li key={`${w.id}-${i}`}>
                            {w.name}{' '}
                            <span className="text-muted-foreground">({w.rarity})</span>
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}
                </>
              )}
              {(!result.drops || result.drops.length === 0) &&
                result.items.length === 0 &&
                result.weapons.length === 0 && (
                  <p className="text-muted-foreground">No items or weapons found — just gold this time.</p>
                )}
            </div>
          )}
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Close
          </Button>
          {previewSessionId ? (
            <Button
              onClick={() => confirmMutation.mutate()}
              disabled={confirmMutation.isPending || (result?.drops?.length ?? 0) === 0}
            >
              {confirmMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Confirm Pickup
            </Button>
          ) : (
            <Button onClick={() => lootMutation.mutate()} disabled={lootMutation.isPending}>
              {lootMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Roll Loot
            </Button>
          )}
          {previewSessionId && (
            <Button
              variant="outline"
              onClick={() => {
                setPreviewSessionId(null);
                setDropAssignments({});
                setResult(null);
              }}
              disabled={confirmMutation.isPending}
            >
              Discard
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
