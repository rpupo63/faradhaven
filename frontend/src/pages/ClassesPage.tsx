import { Link, useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { ArrowLeft } from 'lucide-react';
import { RaIcon } from '@/components/ui/RaIcon';
import { LoadingQuill } from '@/components/LoadingQuill';
import { cn } from '@/lib/utils';
import { useQuery } from '@tanstack/react-query';
import { getClasses, getClassWithLevels } from '@/lib/api';
import { ApiClass, ApiComponent, ApiLevelFeature } from '@/types/game';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { ClassBook } from '@/components/ClassBook';

const getToken = () => typeof window !== 'undefined' ? localStorage.getItem('token') : null;

function formatAbility(ability: string): string {
  const map: Record<string, string> = {
    strength: 'STR',
    dexterity: 'DEX',
    constitution: 'CON',
    intelligence: 'INT',
    wisdom: 'WIS',
    charisma: 'CHA',
  };
  return map[ability?.toLowerCase()] ?? ability.toUpperCase();
}

// Categories that represent core spell substance (highlighted)
const coreCategories = ['Forma', 'Essentia'];

/** Hoverable badge that shows component description on hover */
function ComponentBadge({ component }: { component: ApiComponent }) {
  const isCore = coreCategories.includes(component.category);
  const badge = (
    <Badge
      variant={isCore ? 'default' : 'outline'}
      className={cn(
        'text-xs font-tome-marginalia transition-colors',
        isCore
          ? 'bg-primary/20 text-primary border-primary/30 hover:bg-primary/30'
          : 'border-muted-foreground/30 hover:bg-muted/50',
        component.description && 'cursor-help'
      )}
    >
      {component.name}
    </Badge>
  );

  if (!component.description) {
    return badge;
  }

  return (
    <Tooltip>
      <TooltipTrigger asChild>{badge}</TooltipTrigger>
      <TooltipContent side="top" className="max-w-xs">
        <p className="font-tome-marginalia text-sm">{component.description}</p>
      </TooltipContent>
    </Tooltip>
  );
}

/** Hoverable badge that shows feature description on hover */
function FeatureBadge({ feature }: { feature: ApiLevelFeature }) {
  const badge = (
    <Badge
      variant="outline"
      className={cn(
        'text-xs font-tome-marginalia transition-colors border-faded-gold/50 bg-accent/20 text-accent-foreground hover:bg-accent/30',
        feature.description && 'cursor-help'
      )}
    >
      {feature.name}
    </Badge>
  );

  if (!feature.description) {
    return badge;
  }

  return (
    <Tooltip>
      <TooltipTrigger asChild>{badge}</TooltipTrigger>
      <TooltipContent side="top" className="max-w-sm">
        <p className="font-tome-marginalia text-sm whitespace-pre-wrap">{feature.description}</p>
      </TooltipContent>
    </Tooltip>
  );
}

function ClassCard({ classData }: { classData: ApiClass }) {
  const hasComponents = classData.components && classData.components.length > 0;
  // Sort components: core categories first, then others
  const sortedComponents = [...(classData.components ?? [])].sort((a, b) => {
    const aIsCore = coreCategories.includes(a.category);
    const bIsCore = coreCategories.includes(b.category);
    if (aIsCore && !bIsCore) return -1;
    if (!aIsCore && bIsCore) return 1;
    return 0;
  });

  // Get level 1 features from the class data
  const level1Features = classData.levels?.[0]?.level_features ?? [];
  const hasFeatures = level1Features.length > 0;

  return (
    <Link to={`/game-rules/classes/${classData.id}`} className="block h-full">
      <Card className="arcane-border h-full hover:bg-primary/5 transition-all hover:shadow-lg hover:shadow-primary/10 cursor-pointer group overflow-hidden">
        {/* Class Photo */}
        {classData.photo_url && (
          <div className="relative w-full h-40 overflow-hidden">
            <img
              src={classData.photo_url}
              alt={`${classData.name} class artwork`}
              className="w-full h-full object-cover scale-125 object-top transition-transform group-hover:scale-[1.35]"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-background/90 via-background/20 to-transparent" />
          </div>
        )}
        <CardHeader className={cn("pb-3", classData.photo_url && "pt-4")}>
          <div className="flex items-start justify-between gap-2 mb-2">
            <h3 className="font-tome-heading text-xl text-primary group-hover:glow-text transition-colors">
              {classData.name}
            </h3>
            <Badge variant="secondary" font="mono" size="default">
              d{classData.hit_die}
            </Badge>
          </div>
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <span className="font-tome-marginalia">Primary:</span>
            <Badge variant="secondary" className="font-mono text-xs">
              {formatAbility(classData.primary_ability)}
            </Badge>
          </div>
          {classData.description && (
            <p className="text-sm text-muted-foreground font-tome-marginalia mt-2 line-clamp-3">
              {classData.description}
            </p>
          )}
        </CardHeader>
        <CardContent className="space-y-4">
          <TooltipProvider delayDuration={200}>
            {/* Components Section with overflow ellipsis */}
            {hasComponents && (
              <div className="space-y-2">
                <div className="flex items-center gap-2 text-xs font-tome-subheading text-muted-foreground uppercase tracking-wider">
                  <RaIcon name="aura" className="text-xs" />
                  <span>Components</span>
                </div>
                <div className="flex flex-wrap gap-1.5 max-h-[4.5rem] overflow-hidden relative">
                  {sortedComponents.map((comp) => (
                    <ComponentBadge key={comp.id} component={comp} />
                  ))}
                  {/* Fade overlay for overflow */}
                  <div className="absolute bottom-0 left-0 right-0 h-4 bg-gradient-to-t from-card to-transparent pointer-events-none" />
                </div>
              </div>
            )}

            {/* Level 1 Features Section */}
            <div className="space-y-2">
              <div className="flex items-center gap-2 text-xs font-tome-subheading text-muted-foreground uppercase tracking-wider">
                <RaIcon name="scroll-unfurled" className="text-xs" />
                <span>Starting Features</span>
              </div>
              {hasFeatures ? (
                <div className="flex flex-wrap gap-1.5 max-h-[4.5rem] overflow-hidden relative">
                  {level1Features.map((feature) => (
                    <FeatureBadge key={feature.id} feature={feature} />
                  ))}
                  {/* Fade overlay for overflow */}
                  <div className="absolute bottom-0 left-0 right-0 h-4 bg-gradient-to-t from-card to-transparent pointer-events-none" />
                </div>
              ) : (
                <p className="text-xs text-muted-foreground/80 font-tome-marginalia italic">
                  View class details for features
                </p>
              )}
            </div>

            {/* Proficiencies Preview */}
            {classData.proficiencies && (
              <div className="space-y-1">
                <div className="flex items-center gap-2 text-xs font-tome-subheading text-muted-foreground uppercase tracking-wider">
                  <RaIcon name="shield" className="text-xs" />
                  <span>Proficiencies</span>
                </div>
                <p className="text-xs text-muted-foreground font-tome-marginalia line-clamp-2">
                  {classData.proficiencies}
                </p>
              </div>
            )}

            {/* Skills Preview */}
            {(classData.skill_focus && classData.skill_focus.length > 0) && (
              <div className="space-y-1">
                <div className="text-xs font-tome-subheading text-muted-foreground uppercase tracking-wider">
                  Skills
                </div>
                <div className="flex flex-wrap gap-1">
                  {classData.skill_focus.slice(0, 3).map((skill, idx) => (
                    <Badge
                      key={idx}
                      variant="outline"
                      className="text-xs font-tome-marginalia border-muted-foreground/20"
                    >
                      {skill}
                    </Badge>
                  ))}
                  {classData.skill_focus.length > 3 && (
                    <Badge
                      variant="outline"
                      className="text-xs font-tome-marginalia border-muted-foreground/20"
                    >
                      +{classData.skill_focus.length - 3}
                    </Badge>
                  )}
                </div>
              </div>
            )}
          </TooltipProvider>

          {/* View Details Link */}
          <div className="pt-2 border-t border-border/50">
            <div className="flex items-center gap-2 text-xs text-primary font-tome-marginalia group-hover:gap-3 transition-all">
              <span>View Details</span>
              <ArrowLeft className="h-3 w-3 rotate-180" />
            </div>
          </div>
        </CardContent>
      </Card>
    </Link>
  );
}

export function ClassesTabContent({ classId }: { classId?: string }) {
  const navigate = useNavigate();
  const token = getToken();

  const { data: classes, isLoading: loadingList } = useQuery({
    queryKey: ['classes'],
    queryFn: () => getClasses(token ?? undefined),
    staleTime: 60_000,
    retry: false,
  });

  const { data: classWithLevels, isLoading: loadingDetail } = useQuery({
    queryKey: ['class', classId],
    queryFn: () => getClassWithLevels(classId!, token ?? undefined),
    enabled: !!classId,
    staleTime: 60_000,
    retry: false,
  });

  if (classId) {
    return (
      <div className="w-full space-y-6">
        <Button
          variant="ghost"
          onClick={() => navigate('/game-rules/classes')}
          className="gap-2"
        >
          <ArrowLeft className="h-4 w-4" />
          Back to Classes
        </Button>
        {loadingDetail ? (
          <LoadingQuill label="Loading class details..." />
        ) : classWithLevels ? (
          <ClassBook classData={classWithLevels} />
        ) : (
          <div className="arcane-border rounded-xl p-12 text-center">
            <RaIcon name="book" className="text-6xl mx-auto mb-6 text-muted-foreground block" />
            <h2 className="font-tome-heading text-2xl text-primary mb-2">
              Class Not Found
            </h2>
            <p className="text-muted-foreground font-tome-marginalia mb-6">
              This class may not exist or you may need to log in to view it.
            </p>
            <Button onClick={() => navigate('/game-rules/classes')} variant="seal">
              Back to Classes
            </Button>
          </div>
        )}
      </div>
    );
  }

  return (
    <div className="w-full space-y-12">
      {/* Header */}
      <div className="flex items-center gap-4">
        <div className="p-3 rounded-full border-2 border-faded-gold/50 bg-primary/10">
          <RaIcon name="book" className="text-xl text-primary" />
        </div>
        <div>
          <h1 className="font-tome-heading text-3xl text-primary glow-text">
            Class Compendium
          </h1>
          <p className="text-muted-foreground text-sm font-tome-marginalia mt-1">
            Explore the unique classes of Faradhaven — each with their own components, abilities, and playstyle
          </p>
        </div>
      </div>

      {/* Classes Grid */}
      {loadingList ? (
        <LoadingQuill label="Loading classes..." />
      ) : classes && classes.length > 0 ? (
        <>
          <div className="flex items-center justify-between">
            <p className="text-sm text-muted-foreground font-tome-marginalia">
              {classes.length} {classes.length === 1 ? 'class' : 'classes'} available
            </p>
          </div>
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {classes.map((c) => (
              <ClassCard key={c.id} classData={c} />
            ))}
          </div>
        </>
      ) : (
        <div className="arcane-border rounded-xl p-12 text-center">
          <RaIcon name="book" className="text-6xl mx-auto mb-6 text-muted-foreground block" />
          <h2 className="font-tome-heading text-2xl text-primary mb-2">
            No Classes Available
          </h2>
          <p className="text-muted-foreground font-tome-marginalia mb-6">
            Classes are loaded from the server. Ensure you are logged in and the
            database has been seeded with Faradhaven classes.
          </p>
        </div>
      )}
    </div>
  );
}
