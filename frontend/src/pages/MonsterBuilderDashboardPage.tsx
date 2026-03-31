import React, { useState } from 'react'; // Import useState
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'; // Import useMutation, useQueryClient
import { useAuth } from '@/context/AuthContext';
import { getMonstersByUser, createMonster } from '@/lib/api';
import { MonsterGenerationClassSelect } from '@/components/MonsterGenerationClassSelect';
import { LoadingQuill } from '@/components/LoadingQuill';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { PlusCircle, Hammer, Loader2 } from 'lucide-react'; // Import Loader2 for loading spinner
import { Link } from 'react-router-dom';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog'; // Import Dialog components
import { Label } from '@/components/ui/label'; // Import Label
import { Textarea } from '@/components/ui/textarea'; // Import Textarea
import { Input } from '@/components/ui/input'; // Import Input for CR

export default function MonsterBuilderDashboardPage() {
  const { user, token } = useAuth();
  const queryClient = useQueryClient(); // Initialize queryClient

  const [showGenerateMonsterDialog, setShowGenerateMonsterDialog] = useState(false); // State for dialog
  const [description, setDescription] = useState('');
  const [challengeRating, setChallengeRating] = useState('1'); // Default CR
  const [faradhavenClassName, setFaradhavenClassName] = useState('');

  const { data: monsters, isLoading, error } = useQuery({
    queryKey: ['monsters', user?.id],
    queryFn: () => getMonstersByUser(user!.id, token ?? undefined),
    enabled: !!user?.id && !!token,
  });

  const createMonsterMutation = useMutation({
    mutationFn: (data: {
      description: string;
      challenge_rating: string;
      faradhaven_class_name?: string;
    }) => createMonster({ ...data, user_id: user!.id }, token ?? undefined),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['monsters', user?.id] }); // Invalidate cache to refetch monsters
      setShowGenerateMonsterDialog(false); // Close dialog
      setDescription(''); // Clear form
      setChallengeRating('1'); // Reset CR
      setFaradhavenClassName('');
    },
    onError: (err) => {
      console.error("Failed to generate monster:", err);
      alert("Failed to generate monster. Please try again."); // Basic error feedback
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createMonsterMutation.mutate({
      description,
      challenge_rating: challengeRating,
      ...(faradhavenClassName.trim()
        ? { faradhaven_class_name: faradhavenClassName.trim() }
        : {}),
    });
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[50vh]">
        <LoadingQuill label="Loading monsters..." />
      </div>
    );
  }

  if (error) {
    return <div className="text-destructive">Error loading monsters: {error.message}</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-foreground">Your Monsters</h2>
        <Button size="sm" className="gap-2" onClick={() => setShowGenerateMonsterDialog(true)}>
          <PlusCircle className="w-4 h-4" />
          Generate New Monster
        </Button>
      </div>

      {monsters && monsters.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {monsters.map((monster) => (
            <Link to={`/monster/${monster.id}`} key={monster.id}> {/* Link to monster sheet */}
              <Card className="hover:shadow-lg transition-shadow">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-lg font-medium">{monster.name}</CardTitle>
                  <Hammer className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted-foreground">CR: {monster.challenge_rating}</p>
                  <p className="text-sm text-muted-foreground">Type: {monster.type}</p>
                  {monster.image_url && (
                    <img src={monster.image_url} alt={monster.name} className="mt-4 rounded-md object-cover h-32 w-full" />
                  )}
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      ) : (
        <p className="text-muted-foreground">No monsters generated yet. Click "Generate New Monster" to create your first one!</p>
      )}

      {/* Generate Monster Dialog */}
      <Dialog open={showGenerateMonsterDialog} onOpenChange={setShowGenerateMonsterDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Generate New Monster</DialogTitle>
            <DialogDescription>
              Describe the creature (or the scene), set a Challenge Rating, and optionally base the enemy on a Faradhaven
              class so the AI uses that class&apos;s seeded mechanics and flavor.
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleSubmit} className="grid gap-4 py-4">
            <MonsterGenerationClassSelect
              value={faradhavenClassName}
              onValueChange={setFaradhavenClassName}
              token={token ?? undefined}
              disabled={createMonsterMutation.isPending}
            />
            <div className="grid gap-2">
              <Label htmlFor="description">Description or scene</Label>
              <Textarea
                id="description"
                placeholder="e.g., A wyrmling with scorched scales — or, with a class selected: two guards patrolling a workshop."
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                required
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="cr" className="text-right">
                Challenge Rating
              </Label>
              <Input
                id="cr"
                type="text"
                placeholder="e.g., 1/2, 1, 5"
                value={challengeRating}
                onChange={(e) => setChallengeRating(e.target.value)}
                className="col-span-3"
                required
              />
            </div>
            <DialogFooter>
              <Button type="submit" disabled={createMonsterMutation.isPending}>
                {createMonsterMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Generate Monster
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}