import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/context/AuthContext';
import { getUncheckedSpells, gmUpdateSpell } from '@/lib/api/character';
import type { ApiSpell } from '@/types/game';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { isValidSpellDuration, SPELL_TYPES, SPELL_SAVE_ATTRIBUTES, STANDARD_SPELL_DIE_SIZES } from '@/lib/spellMechanics';
import { DAMAGE_TYPES } from '@/types/game/state';
import { CheckCircle, ChevronDown, ChevronUp, Shield, Wand2, Sparkles } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { Navigate } from 'react-router-dom';
import { cn } from '@/lib/utils';

const GM_EMAIL = 'rpupo63@gmail.com';

/** AI callouts: dark on light parchment, light on dark card */
const aiAccentLabel =
  'text-amber-950 dark:text-amber-200 font-medium';
const aiHintRow =
  'flex items-start gap-1.5 rounded-sm border border-amber-800/25 dark:border-amber-400/25 bg-amber-100/80 dark:bg-amber-950/50 px-2 py-1 text-[11px] text-amber-950 dark:text-amber-100';
const aiHintIcon = 'w-3 h-3 mt-0.5 shrink-0 text-amber-900 dark:text-amber-300';

interface SpellEditState {
  name: string;
  description: string;
  type: string;
  slot_level: number;
  range: number | '';
  duration: string;
  damage_dice_count: number | '';
  damage_die_size: number | '';
  damage_type: string;
  save_attr: string;
}

function parseSpellRangeForEdit(spell: ApiSpell): number | '' {
  const r = spell.range;
  if (typeof r === 'number' && !Number.isNaN(r)) return r;
  return '';
}

function rangeEditString(r: number | ''): string {
  return r === '' ? '' : String(r);
}

function spellToEditState(spell: ApiSpell): SpellEditState {
  return {
    name: spell.name,
    description: spell.description ?? '',
    type: spell.type ?? 'Utility',
    slot_level: spell.slot_level ?? spell.level,
    range: parseSpellRangeForEdit(spell),
    duration: spell.duration ?? '',
    damage_dice_count: typeof spell.damage_dice_count === 'number' && !Number.isNaN(spell.damage_dice_count) ? spell.damage_dice_count : '',
    damage_die_size: typeof spell.damage_die_size === 'number' && !Number.isNaN(spell.damage_die_size) ? spell.damage_die_size : '',
    damage_type: spell.damage_type ?? '',
    save_attr: spell.save_attr ?? '',
  };
}

function buildGmSpellUpdatePayload(edit: SpellEditState) {
  return {
    name: edit.name,
    description: edit.description,
    type: edit.type,
    slot_level: Number(edit.slot_level),
    range: edit.range === '' ? undefined : edit.range,
    duration: edit.duration.trim() || undefined,
    damage_dice_count: edit.damage_dice_count === '' ? undefined : edit.damage_dice_count,
    damage_die_size: edit.damage_die_size === '' ? undefined : edit.damage_die_size,
    damage_type: edit.damage_type.trim() || undefined,
    save_attr: edit.save_attr.trim() || undefined,
  };
}

// Shown below a field when the AI has a different recommendation
function AIHint({
  recommendation,
  opinion,
  currentValue,
  onApply,
}: {
  recommendation?: string | null;
  opinion?: string | null;
  currentValue: string;
  onApply?: (v: string) => void;
}) {
  const hasDiff = recommendation && recommendation !== currentValue;
  const hasOpinion = !!opinion;
  if (!hasDiff && !hasOpinion) return null;

  return (
    <div className="mt-1 space-y-0.5">
      {hasDiff && (
        <div className={aiHintRow}>
          <Sparkles className={aiHintIcon} />
          <span className="italic flex-1">{recommendation}</span>
          {onApply && (
            <button
              type="button"
              onClick={() => onApply(recommendation!)}
              className="text-primary font-medium hover:underline shrink-0 ml-1"
            >
              Apply
            </button>
          )}
        </div>
      )}
      {hasOpinion && (
        <p className="text-[11px] text-foreground/80 dark:text-foreground/75 italic border-l-2 border-amber-800/35 dark:border-amber-500/40 pl-3">
          {opinion}
        </p>
      )}
    </div>
  );
}

function SpellCard({ spell, token }: { spell: ApiSpell; token: string }) {
  const queryClient = useQueryClient();
  const { toast } = useToast();
  const [expanded, setExpanded] = useState(false);
  const [edit, setEdit] = useState<SpellEditState>(spellToEditState(spell));

  const set = (field: keyof SpellEditState) => (v: string | number) =>
    setEdit((s) => ({ ...s, [field]: v }));

  const approveMutation = useMutation({
    mutationFn: () =>
      gmUpdateSpell(spell.id, { ...buildGmSpellUpdatePayload(edit), checked: true }, token),
    onSuccess: () => {
      toast({ title: 'Spell approved', description: `"${edit.name}" marked as checked.` });
      queryClient.invalidateQueries({ queryKey: ['gm-unchecked-spells'] });
    },
    onError: () => toast({ title: 'Error', description: 'Failed to approve spell.', variant: 'destructive' }),
  });

  const saveMutation = useMutation({
    mutationFn: () =>
      gmUpdateSpell(spell.id, buildGmSpellUpdatePayload(edit), token),
    onSuccess: (updated) => {
      toast({ title: 'Spell saved', description: `"${updated.name}" updated.` });
      setEdit(spellToEditState(updated));
    },
    onError: () => toast({ title: 'Error', description: 'Failed to save spell.', variant: 'destructive' }),
  });

  const hasAI = !!spell.ai_overall_verdict;
  const busy = saveMutation.isPending || approveMutation.isPending;
  const durationOk = edit.duration.trim() === '' || isValidSpellDuration(edit.duration.trim());

  return (
    <div className="border border-faded-gold/40 rounded-lg bg-card/80 overflow-hidden">
      {/* Header row */}
      <button
        type="button"
        className="w-full flex items-center justify-between px-4 py-3 text-left hover:bg-primary/5 transition-colors"
        onClick={() => setExpanded((v) => !v)}
      >
        <div className="flex items-center gap-3 flex-wrap">
          <Wand2 className="w-4 h-4 text-primary shrink-0" />
          <span className="font-medium">{spell.name}</span>
          <Badge variant="outline" className="text-xs">Level {spell.slot_level ?? spell.level}</Badge>
          {spell.type && <Badge variant="secondary" className="text-xs">{spell.type}</Badge>}
          {hasAI && (
            <Badge
              variant="outline"
              className="text-[10px] gap-1 border-amber-800/45 dark:border-amber-400/40 text-amber-950 dark:text-amber-200 bg-amber-50/90 dark:bg-amber-950/35"
            >
              <Sparkles className="w-2.5 h-2.5 text-amber-900 dark:text-amber-300" /> AI reviewed
            </Badge>
          )}
        </div>
        <div className="flex items-center gap-2 ml-2 shrink-0">
          <span className="text-xs text-muted-foreground">
            {new Date(spell.created_at ?? '').toLocaleDateString()}
          </span>
          {expanded ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
        </div>
      </button>

      {expanded && (
        <div className="px-4 pb-4 border-t border-faded-gold/20 pt-4 space-y-4">

          {/* AI Verdict banner — only when AI data exists */}
          {hasAI && (
            <div className="rounded-md border border-amber-800/30 dark:border-amber-500/30 bg-amber-100/70 dark:bg-amber-950/45 px-3 py-2 text-xs text-amber-950 dark:text-amber-50">
              <span className="font-semibold uppercase tracking-wide mr-2 text-amber-900 dark:text-amber-200">
                AI Verdict:
              </span>
              {spell.ai_overall_verdict}
              {spell.ai_effect_opinion && (
                <span className="ml-2 text-amber-900/85 dark:text-amber-200/85">· {spell.ai_effect_opinion}</span>
              )}
            </div>
          )}

          {/* Fields grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-x-4 gap-y-3">

            {/* Name */}
            <div className="space-y-1">
              <label className={cn('text-xs uppercase tracking-wide', spell.ai_recommended_name ? aiAccentLabel : 'text-muted-foreground')}>
                Name
              </label>
              <Input value={edit.name} onChange={(e) => set('name')(e.target.value)} />
              <AIHint
                recommendation={spell.ai_recommended_name}
                currentValue={edit.name}
                onApply={set('name')}
              />
            </div>

            {/* Type */}
            <div className="space-y-1">
              <label className={cn('text-xs uppercase tracking-wide', spell.ai_recommended_type ? aiAccentLabel : 'text-muted-foreground')}>
                Type
              </label>
              <Select value={edit.type} onValueChange={(v) => set('type')(v)}>
                <SelectTrigger className="h-9">
                  <SelectValue placeholder="Type" />
                </SelectTrigger>
                <SelectContent>
                  {SPELL_TYPES.map((t) => (
                    <SelectItem key={t} value={t}>
                      {t}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <AIHint
                recommendation={spell.ai_recommended_type}
                currentValue={edit.type}
                onApply={set('type')}
              />
            </div>

            {/* Slot Level */}
            <div className="space-y-1">
              <label className="text-xs text-muted-foreground uppercase tracking-wide">Slot Level</label>
              <Input
                type="number"
                min={1}
                max={9}
                value={edit.slot_level}
                onChange={(e) => set('slot_level')(Number(e.target.value))}
              />
            </div>

            {/* Range */}
            <div className="space-y-1">
              <label className={cn('text-xs uppercase tracking-wide', spell.ai_recommended_range ? aiAccentLabel : 'text-muted-foreground')}>
                Range (feet)
              </label>
              <Input
                type="number"
                min={0}
                step={1}
                value={edit.range === '' ? '' : edit.range}
                onChange={(e) => {
                  const raw = e.target.value;
                  if (raw === '') set('range')('');
                  else {
                    const n = parseInt(raw, 10);
                    if (!Number.isNaN(n) && n >= 0) set('range')(n);
                  }
                }}
                placeholder="0 = self-centered"
              />
              <AIHint
                recommendation={spell.ai_recommended_range != null && spell.ai_recommended_range !== undefined ? String(spell.ai_recommended_range) : undefined}
                currentValue={rangeEditString(edit.range)}
                onApply={(v) => {
                  const n = parseInt(String(v).replace(/\D/g, ''), 10);
                  if (!Number.isNaN(n) && n >= 0) set('range')(n);
                }}
              />
            </div>

            {/* Duration */}
            <div className="space-y-1">
              <label className={cn('text-xs uppercase tracking-wide', spell.ai_recommended_duration ? aiAccentLabel : 'text-muted-foreground')}>
                Duration
              </label>
              <Input
                value={edit.duration}
                onChange={(e) => set('duration')(e.target.value)}
                placeholder="e.g. 1 min, instantaneous"
                className={cn(!durationOk && edit.duration.trim() !== '' && 'border-destructive')}
              />
              {!durationOk && edit.duration.trim() !== '' && (
                <p className="text-[11px] text-destructive">Duration does not match allowed patterns.</p>
              )}
              <AIHint
                recommendation={spell.ai_recommended_duration}
                currentValue={edit.duration}
                onApply={set('duration')}
              />
            </div>

            {/* Damage dice count */}
            <div className="space-y-1">
              <label className={cn('text-xs uppercase tracking-wide', spell.ai_recommended_damage_dice_count != null ? aiAccentLabel : 'text-muted-foreground')}>
                Damage dice count
              </label>
              <Input
                type="number"
                min={1}
                value={edit.damage_dice_count === '' ? '' : edit.damage_dice_count}
                onChange={(e) => {
                  const raw = e.target.value;
                  if (raw === '') set('damage_dice_count')('');
                  else {
                    const n = parseInt(raw, 10);
                    if (!Number.isNaN(n)) set('damage_dice_count')(n);
                  }
                }}
              />
              <AIHint
                recommendation={
                  spell.ai_recommended_damage_dice_count != null
                    ? String(spell.ai_recommended_damage_dice_count)
                    : undefined
                }
                opinion={spell.ai_damage_opinion}
                currentValue={edit.damage_dice_count === '' ? '' : String(edit.damage_dice_count)}
                onApply={(v) => {
                  const n = parseInt(v, 10);
                  if (!Number.isNaN(n)) set('damage_dice_count')(n);
                }}
              />
            </div>

            {/* Die size */}
            <div className="space-y-1">
              <label className={cn('text-xs uppercase tracking-wide', spell.ai_recommended_damage_die_size != null ? aiAccentLabel : 'text-muted-foreground')}>
                Die size
              </label>
              <Select
                value={edit.damage_die_size === '' ? '__none__' : String(edit.damage_die_size)}
                onValueChange={(v) => set('damage_die_size')(v === '__none__' ? '' : parseInt(v, 10))}
              >
                <SelectTrigger className="h-9">
                  <SelectValue placeholder="d6" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="__none__">—</SelectItem>
                  {STANDARD_SPELL_DIE_SIZES.map((sz) => (
                    <SelectItem key={sz} value={String(sz)}>
                      d{sz}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <AIHint
                recommendation={spell.ai_recommended_damage_die_size != null ? String(spell.ai_recommended_damage_die_size) : undefined}
                currentValue={edit.damage_die_size === '' ? '' : String(edit.damage_die_size)}
                onApply={(v) => {
                  const n = parseInt(v, 10);
                  if (!Number.isNaN(n)) set('damage_die_size')(n);
                }}
              />
            </div>

            {/* Damage type */}
            <div className="space-y-1">
              <label className={cn('text-xs uppercase tracking-wide', spell.ai_recommended_damage_type ? aiAccentLabel : 'text-muted-foreground')}>
                Damage type
              </label>
              <Select
                value={edit.damage_type.trim() === '' ? '__none__' : edit.damage_type}
                onValueChange={(v) => set('damage_type')(v === '__none__' ? '' : v)}
              >
                <SelectTrigger className="h-9">
                  <SelectValue placeholder="Type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="__none__">—</SelectItem>
                  {DAMAGE_TYPES.map((t) => (
                    <SelectItem key={t} value={t}>
                      {t}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <AIHint
                recommendation={spell.ai_recommended_damage_type ?? undefined}
                currentValue={edit.damage_type}
                onApply={set('damage_type')}
              />
            </div>

            {/* Save Attribute */}
            <div className="space-y-1">
              <label className={cn('text-xs uppercase tracking-wide', spell.ai_recommended_save_attr ? aiAccentLabel : 'text-muted-foreground')}>
                Save Attribute
              </label>
              <Select
                value={edit.save_attr.trim() === '' ? '__none__' : edit.save_attr}
                onValueChange={(v) => set('save_attr')(v === '__none__' ? '' : v)}
              >
                <SelectTrigger className="h-9">
                  <SelectValue placeholder="—" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="__none__">—</SelectItem>
                  {SPELL_SAVE_ATTRIBUTES.map((a) => (
                    <SelectItem key={a} value={a}>
                      {a}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <AIHint
                recommendation={spell.ai_recommended_save_attr ?? undefined}
                currentValue={edit.save_attr}
                onApply={set('save_attr')}
              />
            </div>
          </div>

          {/* Description (full width) */}
          <div className="space-y-1">
            <label className={cn('text-xs uppercase tracking-wide', spell.ai_recommended_description ? aiAccentLabel : 'text-muted-foreground')}>
              Description
            </label>
            <Textarea
              rows={3}
              value={edit.description}
              onChange={(e) => set('description')(e.target.value)}
            />
            <AIHint
              recommendation={spell.ai_recommended_description}
              opinion={spell.ai_description_opinion}
              currentValue={edit.description}
              onApply={set('description')}
            />
          </div>

          {/* Components */}
          {spell.components && spell.components.length > 0 && (
            <div className="space-y-1">
              <p className="text-xs text-muted-foreground uppercase tracking-wide">Components</p>
              <div className="flex flex-wrap gap-1">
                {spell.components.map((c) => (
                  <Badge key={c.id} variant="outline" className="text-xs">{c.name}</Badge>
                ))}
              </div>
            </div>
          )}

          {/* Actions */}
          <div className="flex gap-2 pt-2">
            <Button
              size="sm"
              variant="outline"
              onClick={() => saveMutation.mutate()}
              disabled={busy || !durationOk}
            >
              Save Changes
            </Button>
            <div className="flex-grow" />
            <Button
              size="sm"
              className="gap-2"
              onClick={() => approveMutation.mutate()}
              disabled={busy || !durationOk}
            >
              <CheckCircle className="w-4 h-4" />
              Approve & Check
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}

export default function GMSpellReviewPage() {
  const { user, token, isLoading } = useAuth();
  const isGM = user?.email === GM_EMAIL;

  const { data: spells, isLoading: spellsLoading } = useQuery({
    queryKey: ['gm-unchecked-spells'],
    queryFn: () => getUncheckedSpells(token!),
    enabled: !!token && isGM,
  });

  if (isLoading) {
    return <div className="text-muted-foreground p-8">Loading...</div>;
  }

  if (!isGM) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3 border-b border-faded-gold/40 pb-4">
        <Shield className="w-6 h-6 text-primary" />
        <div>
          <h1 className="text-3xl font-bold text-primary glow-text">GM Spell Review</h1>
          <p className="text-muted-foreground text-sm mt-1">
            Unchecked spells submitted by players — review and approve them here.
          </p>
        </div>
      </div>

      {spellsLoading && <p className="text-muted-foreground">Loading spells...</p>}

      {!spellsLoading && (!spells || spells.length === 0) && (
        <div className="text-center py-16 text-muted-foreground">
          <CheckCircle className="w-12 h-12 mx-auto mb-3 opacity-30" />
          <p className="text-lg font-medium">All caught up!</p>
          <p className="text-sm">No spells are awaiting review.</p>
        </div>
      )}

      {spells && spells.length > 0 && (
        <div className="space-y-3">
          <p className="text-sm text-muted-foreground">
            {spells.length} spell{spells.length !== 1 ? 's' : ''} pending review
          </p>
          {spells.map((spell) => (
            <SpellCard key={spell.id} spell={spell} token={token!} />
          ))}
        </div>
      )}
    </div>
  );
}
