import React, { useMemo, useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/context/AuthContext';
import { getMonstersByUser, createMonster, previewMonster } from '@/lib/api';
import { MonsterGenerationClassSelect } from '@/components/MonsterGenerationClassSelect';
import { LoadingQuill } from '@/components/LoadingQuill';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { PlusCircle, Hammer, Loader2, Search } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import type { Monster } from '@/types/monster';
import type { MonsterGenerationContext } from '@/lib/api/monster';

const CR_VALUES = ['1/8', '1/4', '1/2', ...Array.from({ length: 20 }, (_, i) => String(i + 1))];
const templateOptions = [
  { id: 'bandit-captain', name: 'Bandit Captain', role: 'skirmisher', temperament: 'disciplined', encounter_goal: 'ambush' },
  { id: 'occult-caster', name: 'Occult Caster', role: 'caster', temperament: 'scheming', encounter_goal: 'control' },
  { id: 'boss-brute', name: 'Boss Brute', role: 'brute', temperament: 'aggressive', encounter_goal: 'frontline' },
];

export default function MonsterBuilderDashboardPage() {
  const enableMonsterAIV2 = import.meta.env.VITE_ENABLE_MONSTER_AI_V2 === 'true';
  const { user, token } = useAuth();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [showGenerateMonsterDialog, setShowGenerateMonsterDialog] = useState(false);
  const [showPreviewDialog, setShowPreviewDialog] = useState(false);
  const [description, setDescription] = useState('');
  const [challengeRating, setChallengeRating] = useState('1');
  const [faradhavenClassName, setFaradhavenClassName] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [filterCR, setFilterCR] = useState('all');
  const [sortBy, setSortBy] = useState<'newest' | 'name' | 'cr'>('newest');
  const [previewResult, setPreviewResult] = useState<Monster | null>(null);
  const [generationContext, setGenerationContext] = useState<MonsterGenerationContext>({
    role: '',
    environment: '',
    temperament: '',
    encounter_goal: '',
    party_level: undefined,
    template_id: '',
    class_theme_intensity: 'strong',
  });

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
      generation_context?: MonsterGenerationContext;
    }) => createMonster({ ...data, user_id: user!.id }, token ?? undefined),
    onSuccess: (newMonster) => {
      queryClient.invalidateQueries({ queryKey: ['monsters', user?.id] });
      navigate(`/monster/${newMonster.id}`);
    },
    onError: (err) => {
      console.error("Failed to generate monster:", err);
      alert("Failed to generate monster. Please try again."); // Basic error feedback
    },
  });

  const previewMonsterMutation = useMutation({
    mutationFn: (data: {
      description: string;
      challenge_rating: string;
      faradhaven_class_name?: string;
      generation_context?: MonsterGenerationContext;
    }) => previewMonster({ ...data, user_id: user!.id }, token ?? undefined),
    onSuccess: (preview) => {
      setPreviewResult(preview);
      setShowPreviewDialog(true);
    },
    onError: () => alert("Failed to generate preview. Please try again."),
  });

  const handleGeneratePreview = (e: React.FormEvent) => {
    e.preventDefault();
    if (!enableMonsterAIV2) {
      createMonsterMutation.mutate({
        description,
        challenge_rating: challengeRating,
        ...(faradhavenClassName.trim()
          ? { faradhaven_class_name: faradhavenClassName.trim() }
          : {}),
        generation_context: generationContext,
      });
      return;
    }
    previewMonsterMutation.mutate({
      description,
      challenge_rating: challengeRating,
      ...(faradhavenClassName.trim()
        ? { faradhaven_class_name: faradhavenClassName.trim() }
        : {}),
      generation_context: generationContext,
    });
  };

  const handleSavePreview = () => {
    if (!previewResult) return;
    createMonsterMutation.mutate({
      description,
      challenge_rating: challengeRating,
      ...(faradhavenClassName.trim()
        ? { faradhaven_class_name: faradhavenClassName.trim() }
        : {}),
      generation_context: generationContext,
    });
  };

  const filteredMonsters = useMemo(() => {
    const base = [...(monsters ?? [])].filter((monster) => {
      const matchesSearch = `${monster.name} ${monster.type}`.toLowerCase().includes(searchQuery.toLowerCase());
      const matchesCR = filterCR === 'all' ? true : monster.challenge_rating === filterCR;
      return matchesSearch && matchesCR;
    });
    if (sortBy === 'name') return base.sort((a, b) => a.name.localeCompare(b.name));
    if (sortBy === 'cr') return base.sort((a, b) => a.challenge_rating.localeCompare(b.challenge_rating));
    return base.sort((a, b) => (b.created_at || '').localeCompare(a.created_at || ''));
  }, [monsters, searchQuery, filterCR, sortBy]);

  const applyTemplate = (templateId: string) => {
    const tpl = templateOptions.find((t) => t.id === templateId);
    if (!tpl) return;
    setGenerationContext((prev) => ({
      ...prev,
      role: tpl.role,
      temperament: tpl.temperament,
      encounter_goal: tpl.encounter_goal,
      template_id: tpl.id,
    }));
    if (!description.trim()) {
      setDescription(`A ${tpl.name.toLowerCase()} patrolling a dangerous district.`);
    }
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

      <div className="grid grid-cols-1 md:grid-cols-4 gap-3">
        <div className="md:col-span-2 relative">
          <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input className="pl-8" placeholder="Search monsters..." value={searchQuery} onChange={(e) => setSearchQuery(e.target.value)} />
        </div>
        <Select value={filterCR} onValueChange={setFilterCR}>
          <SelectTrigger><SelectValue placeholder="Filter CR" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All CR</SelectItem>
            {CR_VALUES.map((cr) => <SelectItem key={cr} value={cr}>{cr}</SelectItem>)}
          </SelectContent>
        </Select>
        <Select value={sortBy} onValueChange={(v) => setSortBy(v as 'newest' | 'name' | 'cr')}>
          <SelectTrigger><SelectValue placeholder="Sort by" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="newest">Newest</SelectItem>
            <SelectItem value="name">Name</SelectItem>
            <SelectItem value="cr">Challenge Rating</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {filteredMonsters.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filteredMonsters.map((monster) => (
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
          <form onSubmit={handleGeneratePreview} className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="template">Template (optional)</Label>
              <Select value={generationContext.template_id || '__none__'} onValueChange={(value) => {
                if (value === '__none__') {
                  setGenerationContext((prev) => ({ ...prev, template_id: '' }));
                  return;
                }
                applyTemplate(value);
              }}>
                <SelectTrigger id="template"><SelectValue placeholder="Choose preset template" /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="__none__">No template</SelectItem>
                  {templateOptions.map((tpl) => <SelectItem key={tpl.id} value={tpl.id}>{tpl.name}</SelectItem>)}
                </SelectContent>
              </Select>
            </div>
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
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="grid gap-2">
                <Label htmlFor="cr">Challenge Rating</Label>
                <Select value={challengeRating} onValueChange={setChallengeRating}>
                  <SelectTrigger id="cr"><SelectValue /></SelectTrigger>
                  <SelectContent>
                    {CR_VALUES.map((cr) => <SelectItem key={cr} value={cr}>{cr}</SelectItem>)}
                  </SelectContent>
                </Select>
              </div>
              <div className="grid gap-2">
                <Label htmlFor="party-level">Party Level (optional)</Label>
                <Input
                  id="party-level"
                  type="number"
                  min={1}
                  max={20}
                  value={generationContext.party_level ?? ''}
                  onChange={(e) => {
                    const next = e.target.value ? Number(e.target.value) : undefined;
                    setGenerationContext((prev) => ({ ...prev, party_level: next }));
                  }}
                />
              </div>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
              <div className="grid gap-2">
                <Label>Role</Label>
                <Select value={generationContext.role || '__none__'} onValueChange={(v) => setGenerationContext((prev) => ({ ...prev, role: v === '__none__' ? '' : v }))}>
                  <SelectTrigger><SelectValue placeholder="Choose role" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="__none__">Any role</SelectItem>
                    <SelectItem value="brute">Brute</SelectItem>
                    <SelectItem value="skirmisher">Skirmisher</SelectItem>
                    <SelectItem value="caster">Caster</SelectItem>
                    <SelectItem value="controller">Controller</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="grid gap-2">
                <Label>Environment</Label>
                <Input value={generationContext.environment || ''} onChange={(e) => setGenerationContext((prev) => ({ ...prev, environment: e.target.value }))} placeholder="e.g. sewers, foundry, ruins" />
              </div>
              <div className="grid gap-2">
                <Label>Temperament</Label>
                <Input value={generationContext.temperament || ''} onChange={(e) => setGenerationContext((prev) => ({ ...prev, temperament: e.target.value }))} placeholder="e.g. disciplined, feral" />
              </div>
            </div>
            <div className="grid gap-2">
              <Label>Encounter Goal</Label>
              <Input value={generationContext.encounter_goal || ''} onChange={(e) => setGenerationContext((prev) => ({ ...prev, encounter_goal: e.target.value }))} placeholder="e.g. ambush, hold chokepoint" />
            </div>
            <div className="grid gap-2">
              <Label>Class theme intensity</Label>
              <Select value={generationContext.class_theme_intensity || 'strong'} onValueChange={(v) => setGenerationContext((prev) => ({ ...prev, class_theme_intensity: v as 'light' | 'strong' }))}>
                <SelectTrigger><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="light">Light</SelectItem>
                  <SelectItem value="strong">Strong</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <DialogFooter>
              <Button type="submit" disabled={previewMonsterMutation.isPending || createMonsterMutation.isPending}>
                {(previewMonsterMutation.isPending || createMonsterMutation.isPending) && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {enableMonsterAIV2 ? 'Generate Preview' : 'Generate Monster'}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      <Dialog open={showPreviewDialog} onOpenChange={setShowPreviewDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Preview Monster</DialogTitle>
            <DialogDescription>Review this generated monster before saving.</DialogDescription>
          </DialogHeader>
          {previewResult ? (
            <div className="space-y-3 text-sm">
              <p><strong>Name:</strong> {previewResult.name}</p>
              <p><strong>CR:</strong> {previewResult.challenge_rating}</p>
              <p><strong>Type:</strong> {previewResult.type}</p>
              <p><strong>AC / HP:</strong> {previewResult.armor_class} / {previewResult.hit_points}</p>
              <p className="text-muted-foreground line-clamp-4">{previewResult.visual_description}</p>
            </div>
          ) : (
            <p className="text-muted-foreground">No preview available.</p>
          )}
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowPreviewDialog(false)}>Keep Editing</Button>
            <Button onClick={handleSavePreview} disabled={createMonsterMutation.isPending}>
              {createMonsterMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Save Monster
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}