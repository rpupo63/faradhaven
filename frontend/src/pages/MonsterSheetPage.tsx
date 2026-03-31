import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/context/AuthContext';
import { getMonsterById, updateMonster, deleteMonster } from '@/lib/api';
import { LoadingQuill } from '@/components/LoadingQuill';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { ArrowLeft, FileText, BookOpen, Trash2, Edit, User as UserIcon, Loader2 } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from '@/components/ui/dialog';

type MonsterTab = 'sheet' | 'notes'; // Added notes tab for editing original prompt

export default function MonsterSheetPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { token } = useAuth();

  const [activeTab, setActiveTab] = useState<MonsterTab>('sheet');
  const [showEditNotesDialog, setShowEditNotesDialog] = useState(false);
  const [editedNotes, setEditedNotes] = useState('');

  const { data: monster, isLoading, error } = useQuery({
    queryKey: ['monster-sheet', id],
    queryFn: () => getMonsterById(id!, token ?? undefined),
    enabled: !!id && !!token,
  });

  const updateMonsterNotesMutation = useMutation({
    mutationFn: (notes: string) => updateMonster(id!, { notes }, token ?? undefined),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['monster-sheet', id] });
      setShowEditNotesDialog(false);
    },
    onError: (err) => {
      console.error('Failed to update monster notes:', err);
      alert('Failed to update monster notes. Please try again.');
    },
  });

  const deleteMonsterMutation = useMutation({
    mutationFn: () => deleteMonster(id!, token ?? undefined),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['monsters'] }); // Invalidate monster list
      navigate('/dm-tools?tab=monster-builder'); // Go back to dashboard
    },
    onError: (err) => {
      console.error('Failed to delete monster:', err);
      alert('Failed to delete monster. Please try again.');
    },
  });

  const handleEditNotes = () => {
    setEditedNotes(monster?.notes || '');
    setShowEditNotesDialog(true);
  };

  const handleSaveNotes = () => {
    updateMonsterNotesMutation.mutate(editedNotes);
  };

  const handleDeleteMonster = () => {
    if (window.confirm(`Are you sure you want to delete ${monster?.name || 'this monster'}? This action cannot be undone.`)) {
      deleteMonsterMutation.mutate();
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[50vh]">
        <LoadingQuill label="Loading monster sheet..." />
      </div>
    );
  }

  if (error) {
    return <div className="text-destructive">Error loading monster: {error.message}</div>;
  }

  if (!monster) {
    return (
      <div className="space-y-4">
        <Button variant="ghost" onClick={() => navigate('/dm-tools?tab=monster-builder')} className="gap-2">
          <ArrowLeft className="h-4 w-4" />
          Back to Monster Builder
        </Button>
        <p className="text-muted-foreground">Monster not found.</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as MonsterTab)} className="w-full">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-4">
          <div className="flex items-center gap-4 w-full md:w-auto">
            <Button variant="ghost" onClick={() => navigate('/dm-tools?tab=monster-builder')} className="gap-2 shrink-0">
              <ArrowLeft className="h-4 w-4" />
              <span className="hidden sm:inline">Back</span>
            </Button>
          </div>

          {/* Tab Navigation */}
          <div className="order-last md:order-none w-full md:w-auto overflow-x-auto">
            <TabsList className="grid grid-cols-2 w-full md:w-auto min-w-[200px]">
              <TabsTrigger value="sheet" className="gap-2 px-3">
                <FileText className="h-4 w-4" />
                <span className="hidden sm:inline">Sheet</span>
              </TabsTrigger>
              <TabsTrigger value="notes" className="gap-2 px-3">
                <BookOpen className="h-4 w-4" />
                <span className="hidden sm:inline">Notes</span>
              </TabsTrigger>
            </TabsList>
          </div>

          {/* Right: Actions and Info */}
          <div className="flex flex-col md:flex-row items-end md:items-center gap-2 md:gap-4 w-full md:w-auto">
            <div className="flex items-center justify-end gap-2 w-full md:w-auto">
              <Button
                variant="outline"
                size="sm"
                onClick={handleEditNotes}
                title="Edit Notes"
                className="gap-1 h-8 px-2"
              >
                <Edit className="h-4 w-4" />
                <span className="hidden lg:inline">Edit Notes</span>
              </Button>
              <Button
                variant="destructive"
                size="sm"
                onClick={handleDeleteMonster}
                title="Delete Monster"
                className="gap-1 h-8 px-2"
                disabled={deleteMonsterMutation.isPending}
              >
                <Trash2 className="h-4 w-4" />
                <span className="hidden lg:inline">Delete</span>
              </Button>
            </div>

            <div className="flex items-center gap-4 md:flex">
              <div className="text-right">
                <h2 className="font-display text-xl text-primary glow-text">
                  {monster.name}
                </h2>
                <p className="text-sm text-muted-foreground font-tome-marginalia">
                  CR {monster.challenge_rating} {monster.size} {monster.type}
                </p>
              </div>

              <div className="relative group">
                <div className="w-14 h-14 rounded-lg arcane-border overflow-hidden bg-primary/5 flex items-center justify-center">
                  {monster.image_url ? (
                    <img src={monster.image_url} alt={monster.name} className="w-full h-full object-cover" />
                  ) : (
                    <UserIcon className="w-8 h-8 text-primary/40" />
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>

        <TabsContent value="sheet" className="mt-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Column 1: Core Stats */}
            <Card>
              <CardHeader><CardTitle>Core Stats</CardTitle></CardHeader>
              <CardContent className="grid gap-2">
                <p><strong>Alignment:</strong> {monster.alignment}</p>
                <p><strong>AC:</strong> {monster.armor_class}</p>
                <p><strong>HP:</strong> {monster.hit_points} ({monster.hit_dice})</p>
                <p><strong>Speed:</strong> {monster.speed}</p>
                <p><strong>Proficiency Bonus:</strong> {monster.proficiency_bonus}</p>
              </CardContent>
            </Card>

            {/* Column 2: Abilities */}
            <Card>
              <CardHeader><CardTitle>Abilities</CardTitle></CardHeader>
              <CardContent className="grid grid-cols-2 gap-2">
                <p><strong>STR:</strong> {monster.strength}</p>
                <p><strong>DEX:</strong> {monster.dexterity}</p>
                <p><strong>CON:</strong> {monster.constitution}</p>
                <p><strong>INT:</strong> {monster.intelligence}</p>
                <p><strong>WIS:</strong> {monster.wisdom}</p>
                <p><strong>CHA:</strong> {monster.charisma}</p>
              </CardContent>
            </Card>

            {/* Column 3: Proficiencies & Immunities */}
            <Card>
              <CardHeader><CardTitle>Proficiencies & Defenses</CardTitle></CardHeader>
              <CardContent className="grid gap-2 text-sm">
                <p><strong>Saving Throws:</strong> {monster.saving_throws.join(', ') || 'None'}</p>
                <p><strong>Skills:</strong> {monster.skills.join(', ') || 'None'}</p>
                <p><strong>Senses:</strong> {monster.senses || 'None'}</p>
                <p><strong>Languages:</strong> {monster.languages || 'None'}</p>
                <p><strong>Damage Immunities:</strong> {monster.damage_immunities.join(', ') || 'None'}</p>
                <p><strong>Damage Resistances:</strong> {monster.damage_resitances.join(', ') || 'None'}</p>
                <p><strong>Condition Immunities:</strong> {monster.condition_immunities.join(', ') || 'None'}</p>
              </CardContent>
            </Card>

            {/* Attacks */}
            <Card className="lg:col-span-3">
              <CardHeader><CardTitle>Attacks</CardTitle></CardHeader>
              <CardContent className="grid gap-4">
                {monster.attacks && monster.attacks.length > 0 ? (
                  monster.attacks.map((attack) => (
                    <div key={attack.id} className="border-b border-border/50 pb-2 last:border-b-0">
                      <h3 className="font-medium text-lg">{attack.name}</h3>
                      <p className="text-sm text-muted-foreground">{attack.description}</p>
                    </div>
                  ))
                ) : (
                  <p className="text-muted-foreground">No attacks listed.</p>
                )}
              </CardContent>
            </Card>

            {/* Actions */}
            <Card className="lg:col-span-3">
              <CardHeader><CardTitle>Actions</CardTitle></CardHeader>
              <CardContent className="grid gap-4">
                {monster.actions && monster.actions.length > 0 ? (
                  monster.actions.map((action) => (
                    <div key={action.id} className="border-b border-border/50 pb-2 last:border-b-0">
                      <h3 className="font-medium text-lg">{action.name}</h3>
                      <p className="text-sm text-muted-foreground">{action.description}</p>
                    </div>
                  ))
                ) : (
                  <p className="text-muted-foreground">No actions listed.</p>
                )}
              </CardContent>
            </Card>

            {/* Special Traits */}
            <Card className="lg:col-span-3">
              <CardHeader><CardTitle>Special Traits</CardTitle></CardHeader>
              <CardContent className="grid gap-2">
                {monster.special_traits && monster.special_traits.length > 0 ? (
                  monster.special_traits.map((trait, index) => (
                    <p key={index} className="text-sm text-muted-foreground">{trait}</p>
                  ))
                ) : (
                  <p className="text-muted-foreground">No special traits listed.</p>
                )}
              </CardContent>
            </Card>

            {/* Legendary Actions */}
            {monster.legendary_actions && monster.legendary_actions.length > 0 && (
              <Card className="lg:col-span-3">
                <CardHeader><CardTitle>Legendary Actions</CardTitle></CardHeader>
                <CardContent className="grid gap-2">
                  {monster.legendary_actions.map((action, index) => (
                    <p key={index} className="text-sm text-muted-foreground">{action}</p>
                  ))}
                </CardContent>
              </Card>
            )}

            {/* Reactions */}
            {monster.reactions && monster.reactions.length > 0 && (
              <Card className="lg:col-span-3">
                <CardHeader><CardTitle>Reactions</CardTitle></CardHeader>
                <CardContent className="grid gap-2">
                  {monster.reactions.map((reaction, index) => (
                    <p key={index} className="text-sm text-muted-foreground">{reaction}</p>
                  ))}
                </CardContent>
              </Card>
            )}

            {/* Lair Actions */}
            {monster.lair_actions && monster.lair_actions.length > 0 && (
              <Card className="lg:col-span-3">
                <CardHeader><CardTitle>Lair Actions</CardTitle></CardHeader>
                <CardContent className="grid gap-2">
                  {monster.lair_actions.map((action, index) => (
                    <p key={index} className="text-sm text-muted-foreground">{action}</p>
                  ))}
                </CardContent>
              </Card>
            )}

          </div>
        </TabsContent>

        <TabsContent value="notes" className="mt-6">
          <Card>
            <CardHeader>
              <CardTitle>Original Notes / Description</CardTitle>
              <CardDescription>This is the original text you provided to the AI to generate the monster. You can edit it here for your own reference.</CardDescription>
            </CardHeader>
            <CardContent>
              <p className="whitespace-pre-wrap text-sm text-muted-foreground">{monster.notes}</p>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {/* Edit Notes Dialog */}
      <Dialog open={showEditNotesDialog} onOpenChange={setShowEditNotesDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Monster Notes</DialogTitle>
            <DialogDescription>
              Modify the notes section of your monster. This is for your personal reference.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <Textarea
              value={editedNotes}
              onChange={(e) => setEditedNotes(e.target.value)}
              rows={10}
              className="min-h-[120px]"
            />
          </div>
          <DialogFooter>
            <Button
              onClick={() => setShowEditNotesDialog(false)}
              variant="outline"
            >
              Cancel
            </Button>
            <Button
              onClick={handleSaveNotes}
              disabled={updateMonsterNotesMutation.isPending}
            >
              {updateMonsterNotesMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Save Changes
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
