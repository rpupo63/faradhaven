import { useState, useEffect, useCallback } from 'react';
import { SharedNote } from '@/types/note';
import { getAllNotes, createNote } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';
import { useToast } from '@/hooks/use-toast';

export function useSharedNotes() {
  const [notes, setNotes] = useState<SharedNote[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const { token } = useAuth();
  const { toast } = useToast();

  const fetchNotes = useCallback(async () => {
    if (!token) return;
    setIsLoading(true);
    try {
      const data = await getAllNotes(token);
      setNotes(data);
    } catch (err) {
      console.error('Error fetching notes:', err);
      toast({
        title: 'Error',
        description: 'Failed to scribing bulletin notes. The archives might be sealed.',
        variant: 'destructive',
      });
    } finally {
      setIsLoading(false);
    }
  }, [toast]);

  useEffect(() => {
    if (token) fetchNotes();
  }, [fetchNotes, token]);

  const addNote = async (
    title: string,
    description: string,
    userId: string,
    username: string,
    partyId?: string | null,
    episodeTag?: string,
    pdf?: File | null
  ) => {
    if (!token) {
      toast({
        title: 'Unauthorized',
        description: 'You must be logged in to post to the bulletin.',
        variant: 'destructive',
      });
      return;
    }

    try {
      const newNote = await createNote(
        {
          title,
          description,
          userId,
          username,
          partyId,
          episodeTag,
        },
        token,
        pdf
      );
      setNotes((prev) => [newNote, ...prev]);
      toast({
        title: 'Success',
        description: 'Your clipping has been pinned to the bulletin.',
      });
    } catch (err) {
      console.error('Error adding note:', err);
      toast({
        title: 'Error',
        description: 'Failed to publish your clipping. The quill has run dry.',
        variant: 'destructive',
      });
    }
  };

  return { notes, isLoading, addNote, refreshNotes: fetchNotes, token };
}
