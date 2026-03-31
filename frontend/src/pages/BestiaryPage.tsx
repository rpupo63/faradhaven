import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Skull, Plus } from 'lucide-react';
import { useGame } from '@/context/GameContext';
import { useAuth } from '@/context/AuthContext';
import { BestiaryCard } from '@/components/BestiaryCard';
import { BestiaryForm } from '@/components/BestiaryForm';
import { BestiaryEntry } from '@/types/game';

export function BestiaryTabContent() {
  const { state, addBestiaryEntry, updateBestiaryEntry, deleteBestiaryEntry } = useGame();
  const { isAuthenticated } = useAuth();
  const [showForm, setShowForm] = useState(false);
  const [editingEntry, setEditingEntry] = useState<BestiaryEntry | null>(null);

  const handleCreate = (data: Omit<BestiaryEntry, 'id' | 'createdAt'>) => {
    addBestiaryEntry(data);
    setShowForm(false);
  };

  const handleUpdate = (data: Omit<BestiaryEntry, 'id' | 'createdAt'>) => {
    if (editingEntry) {
      updateBestiaryEntry(editingEntry.id, data);
      setEditingEntry(null);
    }
  };

  const handleEdit = (entry: BestiaryEntry) => {
    setEditingEntry(entry);
    setShowForm(false);
  };

  const handleDelete = (id: string) => {
    if (confirm('Are you sure you want to delete this creature?')) {
      deleteBestiaryEntry(id);
    }
  };

  return (
    <div className="w-full space-y-8">
      {/* Header – Field Journal title */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-3 rounded-full border-2 border-faded-gold/50 bg-primary/10">
            <Skull className="w-6 h-6 text-primary" />
          </div>
          <div>
            <h1 className="font-tome-heading text-2xl text-primary glow-text">Bestiary</h1>
            <p className="text-muted-foreground text-sm font-tome-marginalia">
              {state.bestiary.length} creature{state.bestiary.length !== 1 ? 's' : ''} catalogued
            </p>
          </div>
        </div>

        {isAuthenticated && !showForm && !editingEntry && (
          <Button onClick={() => setShowForm(true)} variant="quill" size="lg">
            <Plus className="w-4 h-4" />
            Add Creature
          </Button>
        )}
      </div>

      {/* Create/Edit Form */}
      {(showForm || editingEntry) && (
        <BestiaryForm
          entry={editingEntry || undefined}
          onSubmit={editingEntry ? handleUpdate : handleCreate}
          onCancel={() => {
            setShowForm(false);
            setEditingEntry(null);
          }}
        />
      )}

      {/* Bestiary Grid */}
      {state.bestiary.length > 0 ? (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {state.bestiary.map((entry) => (
            <BestiaryCard
              key={entry.id}
              entry={entry}
              onEdit={() => handleEdit(entry)}
              onDelete={() => handleDelete(entry.id)}
            />
          ))}
        </div>
      ) : !showForm && (
        <div className="arcane-border rounded-xl p-12 text-center">
          <Skull className="w-16 h-16 mx-auto mb-6 text-muted-foreground" />
          <h2 className="font-tome-heading text-2xl text-primary mb-2">No Creatures Yet</h2>
          <p className="text-muted-foreground font-tome-marginalia mb-6">
            {isAuthenticated
              ? "Begin cataloguing the monsters your adventurers will face"
              : "No creatures have been catalogued in the bestiary yet."}
          </p>
          {isAuthenticated && (
            <Button onClick={() => setShowForm(true)} variant="seal" size="lg">
              <Plus className="w-4 h-4" />
              Add Your First Creature
            </Button>
          )}
        </div>
      )}
    </div>
  );
}
