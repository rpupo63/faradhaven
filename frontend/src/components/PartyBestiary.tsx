import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getCharacterSheet, addBeastToParty } from '@/lib/api';
import { ApiBeast, ApiCharacter, ApiAttack } from '@/types/game/api';
import { BestiaryEntry } from '@/types/game/base';
import { LoadingQuill } from '@/components/LoadingQuill';
import { BeastCard } from '@/components/BeastCard';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { BugOff, Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { BestiaryForm } from '@/components/BestiaryForm';

interface PartyBestiaryProps {
  characterId: string;
  token?: string;
}

export function PartyBestiary({ characterId, token }: PartyBestiaryProps) {
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);

  const { data: characterSheet, isLoading, error } = useQuery({
    queryKey: ['character-sheet', characterId],
    queryFn: () => getCharacterSheet(characterId, token),
    enabled: !!characterId && !!token,
  });

  const addBeastMutation = useMutation({
    mutationFn: (data: Omit<BestiaryEntry, 'id' | 'createdAt'>) => {
      if (!characterSheet?.character?.party?.id || !token) {
        throw new Error('Party information is not available.');
      }
      
      const apiBeastData: Omit<ApiBeast, 'id'> = {
        user_id: characterSheet.character.user_id,
        name: data.name,
        image_url: data.imageUrl,
        size: data.size,
        type: data.type,
        alignment: data.alignment,
        armor_class: data.armorClass,
        hit_points: data.hitPoints,
        hit_dice: data.hitDice,
        speed: data.speed,
        strength: data.strength,
        dexterity: data.dexterity,
        constitution: data.constitution,
        intelligence: data.intelligence,
        wisdom: data.wisdom,
        charisma: data.charisma,
        challenge_rating: data.challengeRating,
        abilities: data.abilities ? [data.abilities] : [],
        actions: data.actions ? [data.actions] : [],
        legendary_actions: data.legendaryActions ? [data.legendaryActions] : [],
        description: data.description,
        attacks: data.attacks.map(a => ({
          name: a.name,
          attack_bonus: a.attackBonus,
          damage_type: a.damageType,
          damage_dice: a.damageDice,
          reach: a.reach,
          description: a.description,
        })) as ApiAttack[],
      };

      return addBeastToParty(characterSheet.character.party.id, apiBeastData, token);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      setShowForm(false);
    },
    onError: (err) => {
      alert(`Failed to add beast: ${err.message}`);
    },
  });

  if (isLoading) {
    return <LoadingQuill label="Loading party bestiary..." />;
  }

  if (error) {
    return (
      <Alert variant="destructive">
        <BugOff className="h-4 w-4" />
        <AlertTitle>Error loading bestiary</AlertTitle>
        <AlertDescription>
          {error.message || 'An unknown error occurred while fetching party bestiary data.'}
        </AlertDescription>
      </Alert>
    );
  }

  const party = characterSheet?.character?.party;
  const isPartyOwner = party?.owner_id === characterSheet?.character?.user_id;
  const identifiedBeasts: ApiBeast[] = party?.identified_beasts || [];

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-2xl font-bold font-display text-primary glow-text">Identified Bestiary</h3>
          <p className="text-muted-foreground">
            These are the beasts your party has identified and learned about.
          </p>
        </div>
        {isPartyOwner && !showForm && (
          <Button onClick={() => setShowForm(true)}>
            <Plus className="w-4 h-4 mr-2" />
            Add Beast
          </Button>
        )}
      </div>

      {showForm && (
        <BestiaryForm
          onSubmit={(data) => addBeastMutation.mutate(data)}
          onCancel={() => setShowForm(false)}
        />
      )}

      {!showForm && identifiedBeasts.length === 0 ? (
        <Alert>
          <AlertTitle>No Identified Beasts</AlertTitle>
          <AlertDescription>
            Your party has not yet identified any beasts. Explore the world and encounter new creatures!
          </AlertDescription>
        </Alert>
      ) : !showForm && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {identifiedBeasts.map((beast) => (
            <BeastCard key={beast.id} beast={beast} />
          ))}
        </div>
      )}
    </div>
  );
}