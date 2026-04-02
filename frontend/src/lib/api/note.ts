import { getBaseUrl, handleResponse, apiFetch } from './base';
import { SharedNote } from '@/types/note';

export async function getAllNotes(token: string): Promise<SharedNote[]> {
  const res = await apiFetch(`${getBaseUrl()}/api/notes`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return handleResponse<SharedNote[]>(res, 'Failed to fetch notes');
}

export async function createNote(
  note: {
    title: string;
    description: string;
    userId: string;
    username: string;
    partyId?: string | null;
    episodeTag?: string;
  },
  token: string,
  pdf?: File | null
): Promise<SharedNote> {
  const formData = new FormData();
  formData.append('title', note.title);
  formData.append('description', note.description);
  formData.append('userId', note.userId);
  formData.append('username', note.username);
  if (note.partyId) formData.append('partyId', note.partyId);
  if (note.episodeTag) formData.append('episodeTag', note.episodeTag);
  if (pdf) formData.append('pdf', pdf);

  const res = await apiFetch(`${getBaseUrl()}/api/notes`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
    body: formData,
  });
  return handleResponse<SharedNote>(res, 'Failed to create note');
}
