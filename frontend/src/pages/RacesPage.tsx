import { Link, useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Users, ArrowLeft } from 'lucide-react';
import { LoadingQuill } from '@/components/LoadingQuill';
import { cn } from '@/lib/utils';
import { useQuery } from '@tanstack/react-query';
import { getRaces, getRaceWithTraits } from '@/lib/api';
import { ApiRace } from '@/types/game';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { RaceBook } from '@/components/RaceBook';

const getToken = () => typeof window !== 'undefined' ? localStorage.getItem('token') : null;

export function RaceCard({
  raceData,
  asLink = true,
}: {
  raceData: ApiRace;
  /** When false, renders the same card without navigation (e.g. modal preview). */
  asLink?: boolean;
}) {
  const card = (
      <Card className="arcane-border h-full hover:bg-primary/5 transition-all hover:shadow-lg hover:shadow-primary/10 cursor-pointer group overflow-hidden">
        {/* Race Photo */}
        {raceData.photo_url && (
          <div className="relative w-full h-40 overflow-hidden">
            <img
              src={raceData.photo_url}
              alt={`${raceData.name} race artwork`}
              className="w-full h-full object-cover scale-125 -translate-y-[10%] object-top transition-transform group-hover:scale-[1.35]"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-background/90 via-background/20 to-transparent" />
          </div>
        )}
        <CardHeader className={cn('pb-3', raceData.photo_url && 'pt-4')}>
          <div className="flex items-start justify-between gap-2 mb-2">
            <h3 className="font-tome-heading text-xl text-primary group-hover:glow-text transition-colors">
              {raceData.name}
            </h3>
            {raceData.base_speed != null && (
              <Badge variant="outline" font="tomeMarginalia" className="shrink-0">
                {raceData.base_speed} ft.
              </Badge>
            )}
          </div>
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            {raceData.size && (
              <span className="font-tome-marginalia">{raceData.size}</span>
            )}
            {raceData.creature_type && (
              <>
                <span>•</span>
                <span className="font-tome-marginalia">{raceData.creature_type}</span>
              </>
            )}
          </div>
        </CardHeader>
        <CardContent>
          {/* Ability Score Bonuses */}
          {raceData.ability_score_bonuses && Object.keys(raceData.ability_score_bonuses).length > 0 && (
            <div className="flex flex-wrap gap-1.5 mb-3">
              {Object.entries(raceData.ability_score_bonuses).map(([stat, bonus]) => (
                <span key={stat} className="bg-primary/15 text-primary px-2 py-0.5 rounded text-fine font-mono uppercase font-medium">
                  {stat.slice(0, 3)} +{bonus}
                </span>
              ))}
            </div>
          )}
          {raceData.description && (
            <p className="text-sm text-muted-foreground font-tome-marginalia line-clamp-3 mb-4">
              {raceData.description}
            </p>
          )}
          <div className="pt-2 border-t border-border/50">
            <div className="flex items-center gap-2 text-xs text-primary font-tome-marginalia group-hover:gap-3 transition-all">
              <span>View Details</span>
              <ArrowLeft className="h-3 w-3 rotate-180" />
            </div>
          </div>
        </CardContent>
      </Card>
  );

  const wrapClass = 'block h-full';
  if (asLink) {
    return (
      <Link to={`/game-rules/races/${raceData.id}`} className={wrapClass}>
        {card}
      </Link>
    );
  }
  return <div className={wrapClass}>{card}</div>;
}

export function RacesTabContent({ raceId }: { raceId?: string }) {
  const navigate = useNavigate();
  const token = getToken();

  const { data: races, isLoading: loadingList } = useQuery({
    queryKey: ['races'],
    queryFn: () => getRaces(token ?? undefined),
    staleTime: 60_000,
    retry: false,
  });

  const { data: raceWithTraits, isLoading: loadingDetail } = useQuery({
    queryKey: ['race', raceId],
    queryFn: () => getRaceWithTraits(raceId!, token ?? undefined),
    enabled: !!raceId,
    staleTime: 60_000,
    retry: false,
  });

  if (raceId) {
    return (
      <div className="w-full space-y-6">
        <Button
          variant="ghost"
          onClick={() => navigate('/game-rules/races')}
          className="gap-2"
        >
          <ArrowLeft className="h-4 w-4" />
          Back to Races
        </Button>
        {loadingDetail ? (
          <LoadingQuill label="Loading race details..." />
        ) : raceWithTraits ? (
          <RaceBook raceData={raceWithTraits} />
        ) : (
          <div className="arcane-border rounded-xl p-12 text-center">
            <Users className="w-16 h-16 mx-auto mb-6 text-muted-foreground" />
            <h2 className="font-tome-heading text-2xl text-primary mb-2">
              Race Not Found
            </h2>
            <p className="text-muted-foreground font-tome-marginalia mb-6">
              This race may not exist or you may need to log in to view it.
            </p>
            <Button onClick={() => navigate('/game-rules/races')} variant="seal">
              Back to Races
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
          <Users className="w-6 h-6 text-primary" />
        </div>
        <div>
          <h1 className="font-tome-heading text-3xl text-primary glow-text">
            Race Compendium
          </h1>
          <p className="text-muted-foreground text-sm font-tome-marginalia mt-1">
            Explore the races of Faradhaven — each with unique traits and
            abilities
          </p>
        </div>
      </div>

      {/* Races Grid */}
      {loadingList ? (
        <LoadingQuill label="Loading races..." />
      ) : races && races.length > 0 ? (
        <>
          <div className="flex items-center justify-between">
            <p className="text-sm text-muted-foreground font-tome-marginalia">
              {races.length} {races.length === 1 ? 'race' : 'races'} available
            </p>
          </div>
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {races.filter(r => r).map((r) => (
              <RaceCard key={r.id} raceData={r} />
            ))}
          </div>
        </>
      ) : (
        <div className="arcane-border rounded-xl p-12 text-center">
          <Users className="w-16 h-16 mx-auto mb-6 text-muted-foreground" />
          <h2 className="font-tome-heading text-2xl text-primary mb-2">
            No Races Available
          </h2>
          <p className="text-muted-foreground font-tome-marginalia mb-6">
            Races are loaded from the server. Ensure the backend is running and
            the database has been seeded with Faradhaven races.
          </p>
        </div>
      )}
    </div>
  );
}
