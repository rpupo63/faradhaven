import { ApiRaceWithTraits, ApiTrait, ApiTraitOption } from '@/types/game';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { cn } from '@/lib/utils';

interface RaceBookProps {
  raceData: ApiRaceWithTraits;
  className?: string;
}

function FeatureBadge({ label, tooltip }: { label: string; tooltip: string }) {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <span className="cursor-help rounded bg-primary/10 px-1.5 py-0.5 hover:bg-primary/20 transition-colors">
          {label}
        </span>
      </TooltipTrigger>
      <TooltipContent>
        <p className="max-w-xs">{tooltip}</p>
      </TooltipContent>
    </Tooltip>
  );
}

function TraitItem({ trait }: { trait: ApiTrait }) {
  const hasOptions = trait.options && trait.options.length > 0;
  const hasBadges = trait.action_type || trait.uses_per_rest || trait.reset_condition || trait.level_req > 1;
  
  return (
    <div className="space-y-1">
      <h4 className="font-tome-subheading text-primary">{trait.name}</h4>
      {hasBadges && (
        <TooltipProvider delayDuration={200}>
          <div className="flex flex-wrap gap-2 text-xs text-muted-foreground font-tome-marginalia">
            {trait.action_type && (
              <FeatureBadge
                label={trait.action_type}
                tooltip={`Action Type: This trait requires ${trait.action_type === 'Passive' ? 'no action — it is always active' : `using your ${trait.action_type.toLowerCase()} to activate`}.`}
              />
            )}
            {trait.uses_per_rest && (
              <FeatureBadge
                label={`${trait.uses_per_rest}/rest`}
                tooltip={`Uses Per Rest: You can use this trait ${trait.uses_per_rest} time${Number(trait.uses_per_rest) > 1 ? 's' : ''} before needing to rest.`}
              />
            )}
            {trait.reset_condition && (
              <FeatureBadge
                label={trait.reset_condition}
                tooltip={`Reset Condition: This trait's uses are restored after a ${trait.reset_condition.toLowerCase()}.`}
              />
            )}
            {trait.level_req > 1 && (
              <FeatureBadge
                label={`Level ${trait.level_req}+`}
                tooltip={`Level Requirement: You must be at least level ${trait.level_req} to use this trait.`}
              />
            )}
          </div>
        </TooltipProvider>
      )}
      {trait.description && (
        <p className="font-tome-marginalia text-foreground/90 pl-2 border-l-2 border-primary/30">
          {trait.description}
        </p>
      )}
      {hasOptions && (
        <div className="mt-2 pl-2 border-l-2 border-primary/20 space-y-1">
          <span className="text-xs font-tome-subheading text-muted-foreground">
            Options:
          </span>
          {trait.options!.map((opt: ApiTraitOption) => (
            <div key={opt.id} className="text-sm">
              <span className="font-medium text-primary">{opt.name}</span>
              {opt.description && (
                <span className="text-muted-foreground font-tome-marginalia">
                  {' '}
                  — {opt.description}
                </span>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export function RaceBook({ raceData, className }: RaceBookProps) {
  const { name, description, creature_type, size, base_speed, photo_url, ability_score_bonuses, languages, bonus_language_count, traits } =
    raceData;
  const sortedTraits = [...(traits ?? [])].sort((a, b) => a.level_req - b.level_req);
  const hasBonuses = ability_score_bonuses && Object.keys(ability_score_bonuses).length > 0;

  return (
    <div className={cn('space-y-8', className)}>
      {/* Race header – D&D Beyond style */}
      <Card className="arcane-border overflow-hidden">
        {/* Race Photo Banner */}
        {photo_url && (
          <div className="relative w-full h-80 md:h-[28rem] lg:h-[32rem] overflow-hidden">
            <img
              src={photo_url}
              alt={`${name} race artwork`}
              className="w-full h-full object-cover scale-125 -translate-y-[10%] object-top"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-background via-background/40 to-transparent" />
          </div>
        )}
        <CardHeader className={cn('pb-2', photo_url && 'relative -mt-16 z-10')}>
          <CardTitle className="font-tome-heading text-2xl md:text-3xl text-primary glow-text">
            {name}
          </CardTitle>
          <div className="flex flex-wrap gap-4 text-sm text-muted-foreground font-tome-marginalia">
            {creature_type && (
              <span>
                <strong className="text-foreground">Type:</strong> {creature_type}
              </span>
            )}
            {size && (
              <span>
                <strong className="text-foreground">Size:</strong> {size}
              </span>
            )}
            {base_speed != null && (
              <span>
                <strong className="text-foreground">Speed:</strong> {base_speed} ft.
              </span>
            )}
          </div>
          {/* Ability Score Bonuses */}
          {hasBonuses && (
            <div className="flex flex-wrap gap-2 mt-3">
              {Object.entries(ability_score_bonuses!).map(([stat, bonus]) => (
                <span key={stat} className="bg-primary/15 text-primary px-2.5 py-1 rounded text-xs font-mono uppercase font-semibold">
                  {stat.slice(0,3)} +{bonus}
                </span>
              ))}
            </div>
          )}
          {/* Languages */}
          {languages && languages.length > 0 && (
            <div className="flex flex-wrap items-center gap-2 mt-2 text-sm text-muted-foreground font-tome-marginalia">
              <strong className="text-foreground">Languages:</strong>
              {languages.map(lang => (
                <span key={lang} className="bg-muted px-2 py-0.5 rounded text-xs">
                  {lang}
                </span>
              ))}
              {(bonus_language_count ?? 0) > 0 && (
                <span className="text-xs italic">
                  + {bonus_language_count} of your choice
                </span>
              )}
            </div>
          )}
        </CardHeader>
      </Card>

      {/* Description */}
      {description && (
        <Card className="arcane-border overflow-hidden">
          <CardHeader className="pb-2">
            <CardTitle className="font-tome-subheading text-lg text-primary">
              Description
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="tome-prose font-tome-marginalia text-foreground/95 whitespace-pre-wrap">
              {description}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Racial Traits – expandable */}
      <Card className="arcane-border overflow-hidden">
        <CardHeader className="pb-2">
          <CardTitle className="font-tome-subheading text-lg text-primary">
            Racial Traits
          </CardTitle>
          <p className="text-sm text-muted-foreground font-tome-marginalia">
            Abilities and features granted by this race
          </p>
        </CardHeader>
        <CardContent>
          {sortedTraits.length > 0 ? (
            <Accordion type="multiple" defaultValue={['1']} className="w-full">
              {sortedTraits.map((trait) => (
                <AccordionItem
                  key={trait.id}
                  value={trait.id}
                  className="border-border/60"
                >
                  <AccordionTrigger className="font-tome-subheading text-base hover:no-underline py-3">
                    <span className="flex items-center gap-2 text-left">
                      <span>{trait.name}</span>
                      {trait.level_req > 1 && (
                        <span className="text-xs text-muted-foreground">
                          (Level {trait.level_req})
                        </span>
                      )}
                    </span>
                  </AccordionTrigger>
                  <AccordionContent>
                    <TraitItem trait={trait} />
                  </AccordionContent>
                </AccordionItem>
              ))}
            </Accordion>
          ) : (
            <p className="text-muted-foreground font-tome-marginalia italic">
              No traits defined for this race.
            </p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
