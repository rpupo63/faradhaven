import { BestiaryEntry } from '@/types/game';
import { Button } from './ui/button';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
  TooltipProvider,
} from './ui/tooltip';
import { Pencil, Trash2, Skull, Shield, Heart, Swords } from 'lucide-react';
import { cn } from '@/lib/utils';

interface BestiaryCardProps {
  entry: BestiaryEntry;
  onEdit: () => void;
  onDelete: () => void;
}

function modifierFromScore(score: number) {
  const mod = Math.floor((score - 10) / 2);
  return mod >= 0 ? `+${mod}` : `${mod}`;
}

/** Field Journal entry: left = charcoal-sketched illustration, right = ink-stain stat block */
export function BestiaryCard({ entry, onEdit, onDelete }: BestiaryCardProps) {
  return (
    <TooltipProvider delayDuration={200}>
      <article
        className={cn(
          'grid grid-cols-1 lg:grid-cols-[1fr_1.2fr] gap-0 overflow-hidden',
          'arcane-border rounded-lg',
          'bg-card/80'
        )}
      >
        {/* Left: large "charcoal-sketched" monster illustration */}
        <div className="relative aspect-square lg:aspect-auto lg:min-h-[280px] bg-muted/30 flex items-center justify-center overflow-hidden">
          {entry.imageUrl ? (
            <img
              src={entry.imageUrl}
              alt={entry.name}
              className="w-full h-full object-cover charcoal-sketch"
            />
          ) : (
            <div className="flex flex-col items-center justify-center gap-2 text-charcoal/60 p-6">
              <Skull className="w-16 h-16" strokeWidth={1} />
              <span className="font-tome-marginalia text-sm">No sketch</span>
            </div>
          )}
        </div>

        {/* Right: ink-stain stat block – ruler-and-leaky-pen borders */}
        <div className="ink-stain-block p-4 lg:p-5 flex flex-col">
          <div className="flex items-start justify-between gap-2 mb-3">
            <div className="min-w-0">
              <h3 className="font-tome-heading text-lg text-primary truncate">
                {entry.name}
              </h3>
              <p className="text-xs text-muted-foreground font-tome-marginalia mt-0.5">
                {entry.size} {entry.type}, {entry.alignment}
              </p>
            </div>
            <div className="flex gap-1 shrink-0">
              <Button variant="ghost" size="icon" onClick={onEdit} className="h-8 w-8">
                <Pencil className="w-4 h-4" />
              </Button>
              <Button
                variant="ghost"
                size="icon"
                onClick={onDelete}
                className="h-8 w-8 text-destructive hover:text-destructive"
              >
                <Trash2 className="w-4 h-4" />
              </Button>
            </div>
          </div>

          {/* Core stats row */}
          <div className="grid grid-cols-3 gap-2 text-center mb-3">
            <div className="py-2 px-1">
              <Shield className="w-4 h-4 mx-auto text-primary mb-0.5" />
              <p className="text-micro text-muted-foreground font-tome-marginalia">AC</p>
              <p className="font-tome-subheading text-base text-foreground">{entry.armorClass}</p>
            </div>
            <div className="py-2 px-1">
              <Heart className="w-4 h-4 mx-auto text-dried-blood mb-0.5" />
              <p className="text-micro text-muted-foreground font-tome-marginalia">HP</p>
              <p className="font-tome-subheading text-base text-foreground">{entry.hitPoints}</p>
              <p className="text-micro text-muted-foreground">({entry.hitDice})</p>
            </div>
            <div className="py-2 px-1">
              <Skull className="w-4 h-4 mx-auto text-faded-gold mb-0.5" />
              <p className="text-micro text-muted-foreground font-tome-marginalia">CR</p>
              <p className="font-tome-subheading text-base text-foreground">{entry.challengeRating}</p>
            </div>
          </div>

          <p className="text-sm mb-2">
            <span className="text-muted-foreground font-tome-marginalia">Speed </span>
            <span className="text-foreground">{entry.speed}</span>
          </p>

          {/* Ability scores – compact */}
          <div className="grid grid-cols-6 gap-1 text-center text-xs mb-3">
            {[
              { label: 'STR', value: entry.strength },
              { label: 'DEX', value: entry.dexterity },
              { label: 'CON', value: entry.constitution },
              { label: 'INT', value: entry.intelligence },
              { label: 'WIS', value: entry.wisdom },
              { label: 'CHA', value: entry.charisma },
            ].map(({ label, value }) => (
              <div key={label} className="py-1 border-b border-ink/20">
                <p className="text-muted-foreground font-tome-marginalia">{label}</p>
                <p className="font-tome-subheading text-foreground">{value}</p>
                <p className="text-primary text-micro">{modifierFromScore(value)}</p>
              </div>
            ))}
          </div>

          {/* Attacks – hover tooltip = yellowed sticky note */}
          {entry.attacks.length > 0 && (
            <div className="space-y-1.5 mt-auto">
              <div className="flex items-center gap-2 text-sm">
                <Swords className="w-4 h-4 text-dried-blood" />
                <span className="font-tome-subheading text-primary">Attacks</span>
              </div>
              <div className="space-y-1">
                {entry.attacks.map((attack) => (
                  <Tooltip key={attack.id}>
                    <TooltipTrigger asChild>
                      <div className="py-1.5 px-2 rounded border border-ink/20 bg-parchment/50 text-sm cursor-help">
                        <span className="font-medium text-foreground">{attack.name}</span>
                        <span className="text-muted-foreground text-xs ml-2">
                          +{attack.attackBonus} to hit · {attack.damageDice} {attack.damageType}
                        </span>
                      </div>
                    </TooltipTrigger>
                    <TooltipContent className="tooltip-parchment" side="left">
                      <p className="font-tome-marginalia text-xs">
                        <strong>{attack.name}</strong> — +{attack.attackBonus} to hit, {attack.reach || '5 ft.'} reach, {attack.damageDice} {attack.damageType} damage.
                        {attack.description && ` ${attack.description}`}
                      </p>
                    </TooltipContent>
                  </Tooltip>
                ))}
              </div>
            </div>
          )}

          {/* Abilities – hover for sticky-note tooltip */}
          {entry.abilities && (
            <div className="mt-2 pt-2 border-t border-ink/20">
              {entry.abilities.split(/[.;]\s*/).filter(Boolean).slice(0, 3).map((ability: string, i: number) => {
                const label = ability.slice(0, 30).trim() + (ability.length > 30 ? '…' : '');
                return (
                  <Tooltip key={i}>
                    <TooltipTrigger asChild>
                      <span className="text-xs text-foreground border-b border-dashed border-ink/40 cursor-help font-tome-marginalia">
                        {label}
                      </span>
                    </TooltipTrigger>
                    <TooltipContent className="tooltip-parchment" side="left">
                      <p className="font-tome-marginalia text-xs">{ability}</p>
                    </TooltipContent>
                  </Tooltip>
                );
              })}
            </div>
          )}

          {entry.description && (
            <p className="text-xs text-muted-foreground font-tome-marginalia mt-2 pt-2 border-t border-ink/10 italic line-clamp-2">
              {entry.description}
            </p>
          )}
        </div>
      </article>
    </TooltipProvider>
  );
}
