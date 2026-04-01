import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { CharacterCard } from '@/components/CharacterCard';
import { CharacterCreationWizard } from '@/components/CharacterCreationWizard';
import { LoadingQuill } from '@/components/LoadingQuill';
import { RefreshIndicator } from '@/components/ui/refresh-indicator';
import { Button } from '@/components/ui/button';
import { Plus, Users } from 'lucide-react';
import { useQuery, useQueryClient, useMutation } from '@tanstack/react-query';
import { getCharactersByUserId, deleteCharacter } from '@/lib/api';
import { toast } from 'sonner';

export default function CharactersPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { token, userId, activeCharacterId, setActiveCharacter, isAuthenticated } = useAuth();
  const [showWizard, setShowWizard] = useState(false);

  // Fetch characters from API
  const { data: apiResponse, isLoading, isFetching } = useQuery({
    queryKey: ['characters', userId],
    queryFn: () => getCharactersByUserId(userId!, token ?? undefined),
    enabled: isAuthenticated && !!userId && !!token,
  });

  const apiCharacters = apiResponse?.characters || [];

  const deleteCharacterMutation = useMutation({
    mutationFn: (id: string) => deleteCharacter(id, token ?? undefined),
    onSuccess: (_, deletedId) => {
      queryClient.invalidateQueries({ queryKey: ['characters', userId] });
      // If we deleted the active character, clear it
      if (activeCharacterId === deletedId) {
        setActiveCharacter(null).catch(() => {});
      }
    },
  });

  const handleCharacterCreated = (characterId: string) => {
    setShowWizard(false);
    queryClient.invalidateQueries({ queryKey: ['characters', userId] });
    // Set newly created character as active
    setActiveCharacter(characterId).catch(() => {});
    navigate(`/character/${characterId}/sheet`);
  };

  const handleSelectCharacter = (id: string) => {
    navigate(`/character/${id}/sheet`);
  };

  const handleSetActiveCharacter = async (id: string) => {
    try {
      await setActiveCharacter(id);
      toast.success('Active character updated');
    } catch (error) {
      toast.error('Failed to set active character');
    }
  };

  if (showWizard) {
    return (
      <div className="w-full">
        <CharacterCreationWizard
          onComplete={handleCharacterCreated}
          onCancel={() => setShowWizard(false)}
        />
      </div>
    );
  }

  return (
    <div className="w-full space-y-12 relative">
      <RefreshIndicator isFetching={isFetching && !isLoading} />

      {/* Header */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="p-3 rounded-full bg-primary/20">
            <Users className="w-6 h-6 text-primary" />
          </div>
          <div>
            <h1 className="font-tome-heading text-2xl text-primary glow-text">Your Characters</h1>
            <p className="text-muted-foreground text-sm">
              {apiCharacters?.length || 0} character{(apiCharacters?.length !== 1) ? 's' : ''} created
            </p>
          </div>
        </div>

        <Button onClick={() => setShowWizard(true)} variant="quill" size="lg" className="w-full sm:w-auto">
          <Plus className="w-4 h-4" />
          New Character
        </Button>
      </div>

      {isLoading ? (
        <LoadingQuill label="Loading characters..." />
      ) : apiCharacters && apiCharacters.length > 0 ? (
        <div className="grid gap-4 md:grid-cols-2">
          {apiCharacters.map((char) => (
            <CharacterCard
              key={char.id}
              character={{
                id: char.id,
                name: char.name,
                race: char.race?.name || 'Unknown Race',
                class: char.class?.name || 'Unknown Class',
                level: char.level,
                spellbook: char.spellbook || [],
                createdAt: 0 // Placeholder
              }}
              isActive={activeCharacterId === char.id}
              onSelect={() => handleSelectCharacter(char.id)}
              onSetActive={() => handleSetActiveCharacter(char.id)}
              onDelete={() => {
                if (confirm('Are you sure you want to delete this character? This action cannot be undone.')) {
                  deleteCharacterMutation.mutate(char.id);
                }
              }}
            />
          ))}
        </div>
      ) : (
        <div className="arcane-border rounded-xl p-6 md:p-12 text-center">
          <Users className="w-16 h-16 mx-auto mb-6 text-muted-foreground" />
          <h2 className="font-display text-2xl text-primary mb-2">No Characters Yet</h2>
          <p className="text-muted-foreground mb-6">
            Begin your adventure by creating your first character
          </p>
          <Button onClick={() => setShowWizard(true)} variant="seal" size="lg" className="w-full sm:w-auto h-auto py-3 whitespace-normal">
            <Plus className="w-4 h-4" />
            Create Your First Character
          </Button>
        </div>
      )}
    </div>
  );
}
