import { ApiClassWithLevels, ApiClassLevel, ApiLevelFeature, ApiComponent } from '@/types/game';
import { ApiClassResourceDefinition } from '@/types/game/class-resources'; // New import
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
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import { Sparkles } from 'lucide-react';

// Categories that represent core spell substance (highlighted)
const coreCategories = ['Forma', 'Essentia'];

// Category display names
const categoryLabels: Record<string, string> = {
  Forma: 'Forma (Shape)',
  Scopus: 'Scopus (Targeting)',
  Essentia: 'Essentia (Domain)',
  Actio: 'Actio (Action)',
  Magnitudo: 'Magnitudo (Scale)',
  Logica: 'Logica (Flow)',
};

const categoryOrder = ['Forma', 'Scopus', 'Essentia', 'Actio', 'Magnitudo', 'Logica'];

/** Helper to determine the main resource type for a class based on its resource definitions. */
function getResourceType(
  resourceDefs: ApiClassResourceDefinition[] | undefined
): string | undefined {
  if (!resourceDefs || resourceDefs.length === 0) {
    return undefined;
  }

  const keys = new Set(resourceDefs.map((def) => def.key));

  if (keys.has('concurrency_limit') && keys.has('yield_die')) return 'components';
  if (keys.has('timer_duration') && keys.has('speed_dial_slots')) return 'timer';
  if (keys.has('madness_base_dc') || keys.has('feral_bonus')) return 'madness';
  if (keys.has('echo_slots') || keys.has('madness_die')) return 'echo_slots'; // Renamed to echo_slots to encompass Lorewright's multiple slots
  if (keys.has('max_stability')) return 'stability';
  if (keys.has('max_blood_ichor') && keys.has('bite_damage_dice')) return 'blood_ichor';
  if (keys.has('shadow_points')) return 'shadow_points';
  if (keys.has('spell_points')) return 'spell_points';

  return undefined;
}

/** Hoverable badge that shows feature description on hover */
function FeatureBadge({ name, description }: { name: string; description?: string }) {
  if (!description) {
    return (
      <Badge
        variant="outline"
        className="cursor-default font-tome-marginalia text-xs border-primary/30 bg-primary/5 hover:bg-primary/10 transition-colors"
      >
        {name}
      </Badge>
    );
  }

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <Badge
          variant="outline"
          className="cursor-help font-tome-marginalia text-xs border-primary/30 bg-primary/5 hover:bg-primary/15 transition-colors"
        >
          {name}
        </Badge>
      </TooltipTrigger>
      <TooltipContent side="top" className="max-w-sm">
        <p className="font-tome-marginalia text-sm">{description}</p>
      </TooltipContent>
    </Tooltip>
  );
}

function formatAbility(ability: string): string {
  const map: Record<string, string> = {
    strength: 'Strength',
    dexterity: 'Dexterity',
    constitution: 'Constitution',
    intelligence: 'Intelligence',
    wisdom: 'Wisdom',
    charisma: 'Charisma',
  };
  return map[ability?.toLowerCase()] ?? ability;
}

/** Component badge with hover tooltip showing description */
function ComponentBadge({ component }: { component: ApiComponent }) {
  const isCore = coreCategories.includes(component.category);
  const badge = (
    <Badge
      variant={isCore ? 'default' : 'outline'}
      className={cn(
        'text-sm font-tome-marginalia transition-colors',
        isCore
          ? 'bg-primary/20 text-primary border-primary/40 hover:bg-primary/30'
          : 'border-muted-foreground/40 hover:bg-muted/50',
        component.description && 'cursor-help'
      )}
    >
      <Sparkles className="h-3 w-3 mr-1.5" />
      {component.name}
    </Badge>
  );

  if (!component.description) {
    return badge;
  }

  return (
    <Tooltip>
      <TooltipTrigger asChild>{badge}</TooltipTrigger>
      <TooltipContent side="top" className="max-w-sm">
        <p className="font-tome-marginalia text-sm whitespace-pre-wrap">{component.description}</p>
      </TooltipContent>
    </Tooltip>
  );
}

/** Render feature badges for the table cell with hover tooltips */
function FeatureBadges({ features }: { features?: ApiLevelFeature[] }) {
  if (!features || features.length === 0) {
    return <span className="text-muted-foreground">—</span>;
  }

  return (
    <TooltipProvider delayDuration={200}>
      <div className="flex flex-wrap gap-1.5">
        {features.map((feature) => (
          <FeatureBadge
            key={feature.id}
            name={feature.name}
            description={feature.description}
          />
        ))}
      </div>
    </TooltipProvider>
  );
}


interface ClassBookProps {
  classData: ApiClassWithLevels;
  className?: string;
}

export function ClassBook({ classData, className }: ClassBookProps) {
  const { name, hit_die, primary_ability, photo_url, levels, components, resource_definitions } = classData;
  const sortedLevels = [...(levels ?? [])].sort((a, b) => a.level - b.level);
  
  // Group components by category
  const groupedComponents = components?.reduce((acc, comp) => {
    (acc[comp.category] = acc[comp.category] || []).push(comp);
    return acc;
  }, {} as Record<string, ApiComponent[]>) || {};

  const hasComponents = components && components.length > 0;

  const resourceType = getResourceType(resource_definitions);

  return (
    <div className={cn('space-y-8', className)}>
      {/* Class header – D&D Beyond style */}
      <Card className="arcane-border overflow-hidden">
        {/* Class Photo Banner */}
        {photo_url && (
          <div className="relative w-full h-80 md:h-[28rem] lg:h-[32rem] overflow-hidden">
            <img
              src={photo_url}
              alt={`${name} class artwork`}
              className="w-full h-full object-cover scale-125 object-top"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-background via-background/40 to-transparent" />
          </div>
        )}
        <CardHeader className={cn("pb-2", photo_url && "relative -mt-16 z-10")}>
          <CardTitle className="font-tome-heading text-2xl md:text-3xl text-primary glow-text">
            {name}
          </CardTitle>
          {classData.description && (
            <p className="text-base text-muted-foreground font-tome-marginalia mt-2 max-w-2xl">
              {classData.description}
            </p>
          )}
          <div className="flex flex-wrap gap-4 text-sm text-muted-foreground font-tome-marginalia mt-2">
            <span>
              <strong className="text-foreground">Hit Die:</strong> {hit_die ? `d${hit_die}` : 'N/A'}
            </span>
            <span>
              <strong className="text-foreground">Primary Ability:</strong>{' '}
              {primary_ability ? formatAbility(primary_ability) : 'N/A'}
            </span>
          </div>
        </CardHeader>
      </Card>

      {/* Class Components – Full list available from level 1 */}
      {hasComponents && (
        <Card className="arcane-border overflow-hidden">
          <CardHeader className="pb-2">
            <CardTitle className="font-tome-subheading text-lg text-primary flex items-center gap-2">
              <Sparkles className="h-5 w-5" />
              Spell Components
            </CardTitle>
            <p className="text-sm text-muted-foreground font-tome-marginalia">
              Magical components available to this class for crafting spells — unlocked from level 1
            </p>
          </CardHeader>
          <CardContent className="space-y-4">
            <TooltipProvider delayDuration={200}>
              {/* Components grouped by category */}
              {categoryOrder.map((cat) => {
                const comps = groupedComponents[cat];
                if (!comps || comps.length === 0) return null;
                return (
                  <div key={cat} className="space-y-2">
                    <h4 className="text-sm font-tome-subheading text-muted-foreground uppercase tracking-wider flex items-center gap-2">
                      <Sparkles className="h-3.5 w-3.5" />
                      {categoryLabels[cat] || cat}
                    </h4>
                    <div className="flex flex-wrap gap-2">
                      {comps.map((comp) => (
                        <ComponentBadge key={comp.id} component={comp} />
                      ))}
                    </div>
                  </div>
                );
              })}
            </TooltipProvider>
          </CardContent>
        </Card>
      )}

      {/* Madness Table (Lorewright) */}
      {classData.madness_table && (
        <Card className="arcane-border overflow-hidden">
          <CardHeader className="pb-2">
            <CardTitle className="font-tome-subheading text-lg text-primary">
              Madness Effects
            </CardTitle>
            <p className="text-sm text-muted-foreground font-tome-marginalia">
              Roll a d50 on this table when the mind fractures
            </p>
          </CardHeader>
          <CardContent className="p-0">
            <div className="overflow-x-auto max-h-[500px]">
              <table className="w-full tome-table">
                <thead>
                  <tr>
                    <th className="w-16 text-center">d50</th>
                    <th>Effect</th>
                  </tr>
                </thead>
                <tbody>
                  {Object.entries(classData.madness_table)
                    .sort((a, b) => Number(a[0]) - Number(b[0]))
                    .map(([roll, effect]) => (
                      <tr key={roll}>
                        <td className="text-center font-tome-subheading">{roll}</td>
                        <td className="font-tome-marginalia">{effect}</td>
                      </tr>
                    ))}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Level progression table – D&D PHB style */}
      <Card className="arcane-border overflow-hidden">
        <CardHeader className="pb-2">
          <CardTitle className="font-tome-subheading text-lg text-primary">
            Level Progression
          </CardTitle>
          <p className="text-sm text-muted-foreground font-tome-marginalia">
            Features gained at each level
          </p>
        </CardHeader>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full tome-table">
              <thead>
                <tr>
                  <th>Level</th>
                  <th>Prof. Bonus</th>
                  {resourceType === 'components' && (
                    <>
                      <th>Concurrency</th>
                      <th>Yield</th>
                    </>
                  )}
                  {resourceType === 'timer' && (
                    <>
                      <th>Timer</th>
                      <th>Speed Dial</th>
                      <th>Max Length</th>
                    </>
                  )}
                  {resourceType === 'madness' && (
                    <>
                      <th>Feral Bonus</th>
                    </>
                  )}
                  {resourceType === 'echo_slots' && (
                    <>
                      <th>Harvest Slots</th>
                      <th>Madness Die</th>
                    </>
                  )}
                  {resourceType === 'stability' && (
                    <>
                      <th>Max Stability</th>
                      <th>Max Spell Level</th>
                    </>
                  )}
                  {resourceType === 'blood_ichor' && (
                    <>
                      <th>Ichor Points</th>
                      <th>Bite Damage</th>
                    </>
                  )}
                  {resourceType === 'shadow_points' && (
                    <th>Shadow Points</th>
                  )}
                  {sortedLevels.some((l) => (l.sneak_attack_dice ?? 0) > 0) && (
                    <th>Sneak Strike</th>
                  )}
                  {sortedLevels.some(
                    (l) => l.cantrips_known !== undefined || l.spells_known !== undefined
                  ) && (
                    <>
                      <th>Cantrips</th>
                      <th>Spells</th>
                    </>
                  )}
                  {(resourceType === 'spell_points' || !resourceType) && (
                    <th>Spell Points</th>
                  )}
                  <th>Features</th>
                </tr>
              </thead>
              <tbody>
                {sortedLevels.map((cl) => (
                  <tr key={cl.id}>
                    <td className="font-tome-subheading">{cl.level}</td>
                    <td>+{cl.proficiency_bonus}</td>
                    {resourceType === 'components' && (
                      <>
                        <td>{cl.resources?.concurrency_limit}</td>
                        <td>d{cl.resources?.yield_die}</td>
                      </>
                    )}
                    {resourceType === 'timer' && (
                      <>
                        <td>{cl.resources?.timer_duration}s</td>
                        <td>{cl.resources?.speed_dial_slots}</td>
                        <td>{cl.max_spell_level ? cl.max_spell_level : '∞'}</td>
                      </>
                    )}
                    {resourceType === 'madness' && (
                      <>
                        <td>+{cl.resources?.feral_bonus}</td>
                      </>
                    )}
                    {resourceType === 'echo_slots' && (
                      <>
                        <td>{cl.resources?.echo_slots || 0}</td>
                        <td>d{cl.resources?.madness_die}</td>
                      </>
                    )}
                    {resourceType === 'stability' && (
                      <>
                        <td>{cl.resources?.max_stability}</td>
                        <td>{cl.max_spell_level || '—'}</td>
                      </>
                    )}
                    {resourceType === 'blood_ichor' && (
                      <>
                        <td>{cl.resources?.max_blood_ichor || '—'}</td>
                        <td>{cl.resources?.bite_damage_dice ? `${cl.resources?.bite_damage_dice}d6` : '—'}</td>
                      </>
                    )}
                    {resourceType === 'shadow_points' && (
                      <td>{cl.resources?.shadow_points || '—'}</td>
                    )}
                    {sortedLevels.some((l) => (l.sneak_attack_dice ?? 0) > 0) && (
                      <td>{cl.sneak_attack_dice ? `${cl.sneak_attack_dice}d6` : '—'}</td>
                    )}
                    {sortedLevels.some(
                      (l) => l.cantrips_known !== undefined || l.spells_known !== undefined
                    ) && (
                      <>
                        <td>{cl.cantrips_known ?? '—'}</td>
                        <td>{cl.spells_known ?? '—'}</td>
                      </>
                    )}
                    {(resourceType === 'spell_points' || !resourceType) && (
                      <td>{cl.max_spell_points}</td>
                    )}
                    <td className="max-w-[320px]">
                      <FeatureBadges features={cl.level_features} />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      {/* Level-up descriptions – expandable by level */}
      <Card className="arcane-border overflow-hidden">
        <CardHeader className="pb-2">
          <CardTitle className="font-tome-subheading text-lg text-primary">
            Level Up — Feature Descriptions
          </CardTitle>
          <p className="text-sm text-muted-foreground font-tome-marginalia">
            Detailed descriptions of features gained at each level
          </p>
        </CardHeader>
        <CardContent>
          <Accordion type="multiple" defaultValue={['1']} className="w-full">
            {sortedLevels.map((cl) => (
              <LevelAccordionItem key={cl.id} level={cl} />
            ))}
          </Accordion>
        </CardContent>
      </Card>
    </div>
  );
}

function LevelAccordionItem({ level }: { level: ApiClassLevel }) {
  const hasFeatures = level.level_features && level.level_features.length > 0;
  const asi = level.ability_score_improvement ?? 0;
  const badge =
    asi > 0 ? (
      <span className="ml-2 rounded bg-primary/15 px-1.5 py-0.5 text-xs font-tome-marginalia text-primary">
        ASI
      </span>
    ) : null;

  // Get feature names for the trigger
  const featureNames = hasFeatures
    ? level.level_features!.map(f => f.name).join(', ')
    : null;

  return (
    <AccordionItem value={level.level.toString()} className="border-border/60">
      <AccordionTrigger className="font-tome-subheading text-base hover:no-underline py-3">
        <span className="flex items-center gap-2 text-left">
          <span className="shrink-0">Level {level.level}</span>
          {badge}
          {featureNames && (
            <span className="text-sm font-tome-marginalia text-muted-foreground">
              — {featureNames}
            </span>
          )}
        </span>
      </AccordionTrigger>
      <AccordionContent>
        {hasFeatures ? (
          <div className="space-y-4">
            {level.level_features!.map((feature) => (
              <FeatureItem key={feature.id} feature={feature} />
            ))}
          </div>
        ) : (
          <p className="text-muted-foreground font-tome-marginalia italic">
            No additional features at this level.
          </p>
        )}
      </AccordionContent>
    </AccordionItem>
  );
}

function FeatureItem({ feature }: { feature: ApiLevelFeature }) {
  return (
    <div className="space-y-2">
      <div className="flex items-start gap-2">
        <TooltipProvider delayDuration={200}>
          <Tooltip>
            <TooltipTrigger asChild>
              <h4 className="font-tome-subheading text-primary cursor-help inline-flex items-center gap-1.5 rounded px-1.5 py-0.5 bg-primary/5 hover:bg-primary/10 transition-colors">
                {feature.name}
              </h4>
            </TooltipTrigger>
            {feature.description && (
              <TooltipContent side="right" className="max-w-md">
                <p className="font-tome-marginalia text-sm whitespace-pre-wrap">{feature.description}</p>
              </TooltipContent>
            )}
          </Tooltip>
        </TooltipProvider>
      </div>
      {feature.description && (
        <p className="font-tome-marginalia text-foreground/90 pl-3 border-l-2 border-primary/30 whitespace-pre-wrap">
          {feature.description}
        </p>
      )}
    </div>
  );
}

