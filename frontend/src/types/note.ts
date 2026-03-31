export interface SharedNote {
  id: string;
  title: string;
  description: string;
  pdfUrl?: string | null;
  userId: string;
  username: string;
  partyId?: string | null;
  episodeTag?: string;
  createdAt: string;
}
