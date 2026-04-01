import { useEffect, useRef, useState } from 'react';
import { useSharedNotes } from '@/hooks/useSharedNotes';
import { useAuth } from '@/context/AuthContext';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Newspaper, Plus, User, Calendar, FileText, Lock, BookOpen } from 'lucide-react';
import { format } from 'date-fns';
import { cn } from '@/lib/utils';
import { SharedNote } from '@/types/note';
import { getPartiesByUserId } from '@/lib/api/party';
import { ApiParty } from '@/types/game/api';

export default function BulletinPage() {
  const { notes, isLoading, addNote, token } = useSharedNotes();
  const { user } = useAuth();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [newNote, setNewNote] = useState({ title: '', description: '', episodeTag: '' });
  const [selectedPartyId, setSelectedPartyId] = useState<string>('public');
  const [pdfFile, setPdfFile] = useState<File | null>(null);
  const [parties, setParties] = useState<ApiParty[]>([]);
  const pdfInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (user && token) {
      getPartiesByUserId(user.id, token).then(setParties).catch(() => {});
    }
  }, [user, token]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newNote.title || !newNote.description || !user) return;
    const partyId = selectedPartyId === 'public' ? null : selectedPartyId;
    addNote(newNote.title, newNote.description, user.id, user.name, partyId, newNote.episodeTag || undefined, pdfFile);
    setNewNote({ title: '', description: '', episodeTag: '' });
    setSelectedPartyId('public');
    setPdfFile(null);
    if (pdfInputRef.current) pdfInputRef.current.value = '';
    setIsDialogOpen(false);
  };

  return (
    <div className="space-y-8 animate-ink-bleed">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <h1 className="flex items-center gap-3">
            <Newspaper className="w-8 h-8 text-primary" />
            The Faradhaven Bulletin
          </h1>
          <p className="text-muted-foreground font-tome-marginalia mt-1">
            Freshly clipped news and rumors from across the realm.
          </p>
        </div>

        {user && (
          <Dialog
            open={isDialogOpen}
            onOpenChange={(open) => {
              setIsDialogOpen(open);
              if (!open) {
                setPdfFile(null);
                if (pdfInputRef.current) pdfInputRef.current.value = '';
              }
            }}
          >
            <DialogTrigger asChild>
              <Button className="hand-drawn-border bg-primary/10 hover:bg-primary/20 text-primary border-primary/40 gap-2">
                <Plus className="w-4 h-4" />
                Post Clipping
              </Button>
            </DialogTrigger>
            <DialogContent className="vellum-modal-content max-w-md">
              <DialogHeader>
                <DialogTitle className="font-tome-heading text-2xl text-primary">Draft a Clipping</DialogTitle>
              </DialogHeader>
              <form onSubmit={handleSubmit} className="space-y-4 py-4">
                <div className="space-y-2">
                  <label className="font-tome-subheading text-sm uppercase tracking-wider">Headline</label>
                  <Input
                    value={newNote.title}
                    onChange={(e) => setNewNote({ ...newNote, title: e.target.value })}
                    placeholder="Extra! Extra! Read all about it..."
                    className="tome-input-underline"
                    maxLength={100}
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label className="font-tome-subheading text-sm uppercase tracking-wider">The Full Story</label>
                  <Textarea
                    value={newNote.description}
                    onChange={(e) => setNewNote({ ...newNote, description: e.target.value })}
                    placeholder="Detail the events, rumors, or calls for adventure..."
                    className="min-h-[150px] bg-white/30 border-2 border-faded-gold/30 focus:border-primary/50 transition-colors font-serif italic"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label className="font-tome-subheading text-sm uppercase tracking-wider">Episode / Session (optional)</label>
                  <Input
                    value={newNote.episodeTag}
                    onChange={(e) => setNewNote({ ...newNote, episodeTag: e.target.value })}
                    placeholder="e.g. Episode 7 or The Sunken Vault"
                    className="tome-input-underline"
                    maxLength={80}
                  />
                </div>
                <div className="space-y-2">
                  <label className="font-tome-subheading text-sm uppercase tracking-wider">Visibility</label>
                  <Select value={selectedPartyId} onValueChange={setSelectedPartyId}>
                    <SelectTrigger className="tome-input-underline">
                      <SelectValue placeholder="Public (everyone)" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="public">Public (everyone)</SelectItem>
                      {parties.map((p) => (
                        <SelectItem key={p.id} value={p.id}>
                          {p.name} (party only)
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <label className="font-tome-subheading text-sm uppercase tracking-wider">
                    Attachment (PDF, optional)
                  </label>
                  <Input
                    ref={pdfInputRef}
                    type="file"
                    accept="application/pdf"
                    className="tome-input-underline cursor-pointer file:mr-3 file:rounded file:border-0 file:bg-primary/10 file:px-2 file:py-1 file:text-sm"
                    onChange={(e) => setPdfFile(e.target.files?.[0] ?? null)}
                  />
                  {pdfFile && (
                    <div className="flex items-center justify-between gap-2 text-xs text-muted-foreground">
                      <span className="truncate font-tome-marginalia">{pdfFile.name}</span>
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        className="shrink-0 h-7 text-xs"
                        onClick={() => {
                          setPdfFile(null);
                          if (pdfInputRef.current) pdfInputRef.current.value = '';
                        }}
                      >
                        Clear
                      </Button>
                    </div>
                  )}
                </div>
                <DialogFooter>
                  <Button type="submit" className="w-full bg-primary hover:bg-primary/90 text-white">
                    Publish to the Bulletin
                  </Button>
                </DialogFooter>
              </form>
            </DialogContent>
          </Dialog>
        )}
      </div>

      {isLoading ? (
        <div className="flex justify-center py-20">
          <div className="animate-quill-write text-primary font-tome-marginalia text-lg">
            Scribing bulletin...
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {notes.map((note, index) => (
            <NewsClipping key={note.id} note={note} index={index} />
          ))}
        </div>
      )}
    </div>
  );
}

function NewsClipping({ note, index }: { note: SharedNote; index: number }) {
  // Random rotation for that "clipped and pinned" look
  const rotation = (index % 3 - 1) * 1.5;
  const isLong = note.description.length > 200;
  const snippet = isLong ? note.description.slice(0, 200) + '...' : note.description;
  
  return (
    <Dialog>
      <DialogTrigger asChild>
        <div 
          className={cn(
            "relative p-6 bg-[#fdfaf3] text-[#2c2421] shadow-xl transition-all hover:scale-[1.02] hover:z-10 group cursor-pointer",
            "before:absolute before:inset-0 before:bg-[url('https://www.transparenttextures.com/patterns/pinstriped-suit.png')] before:opacity-5 before:pointer-events-none"
          )}
          style={{ 
            transform: `rotate(${rotation}deg)`,
            border: '1px solid rgba(44, 36, 33, 0.2)',
            boxShadow: '4px 4px 15px rgba(0,0,0,0.1), -1px -1px 2px rgba(255,255,255,0.5)'
          }}
        >
          {/* Tape effect at the top */}
          <div className="absolute -top-3 left-1/2 -translate-x-1/2 w-16 h-8 bg-primary/10 border border-primary/5 rotate-1 opacity-60 pointer-events-none" />

          <div className="flex flex-col h-full space-y-4">
            <div className="border-b-2 border-primary/20 pb-2">
              <h2 className="font-tome-heading text-2xl uppercase tracking-tighter leading-tight group-hover:text-primary transition-colors">
                {note.title}
              </h2>
              <div className="flex flex-wrap gap-1.5 mt-1.5">
                {note.episodeTag && (
                  <span className="inline-flex items-center gap-1 text-[0.65rem] font-tome-marginalia px-1.5 py-0.5 rounded bg-primary/10 text-primary border border-primary/20">
                    <BookOpen className="w-2.5 h-2.5" />
                    {note.episodeTag}
                  </span>
                )}
                {note.partyId && (
                  <span className="inline-flex items-center gap-1 text-[0.65rem] font-tome-marginalia px-1.5 py-0.5 rounded bg-amber-100 text-amber-800 border border-amber-300">
                    <Lock className="w-2.5 h-2.5" />
                    Party only
                  </span>
                )}
              </div>
            </div>
            
            <div className="flex-1">
              <p className="font-serif leading-relaxed text-sm first-letter:text-4xl first-letter:font-bold first-letter:mr-1 first-letter:float-left first-letter:leading-[0.8]">
                {snippet}
              </p>
              {isLong && (
                <span className="text-primary text-xs font-tome-marginalia mt-2 block hover:underline">
                  Read full clipping...
                </span>
              )}
              {note.pdfUrl && (
                <a
                  href={note.pdfUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  onClick={(e) => e.stopPropagation()}
                  className="mt-3 inline-flex items-center gap-1.5 text-primary text-xs font-tome-marginalia hover:underline relative z-20"
                >
                  <FileText className="w-3.5 h-3.5 shrink-0" />
                  View PDF
                </a>
              )}
            </div>

            <div className="pt-4 mt-auto border-t border-primary/10 flex flex-col gap-2">
              <div className="flex items-center justify-between text-[0.7rem] font-tome-marginalia opacity-70 italic">
                <div className="flex items-center gap-1">
                  <User className="w-3 h-3" />
                  <span>Reported by: {note.username}</span>
                </div>
                <div className="flex items-center gap-1">
                  <Calendar className="w-3 h-3" />
                  <span>{format(new Date(note.createdAt), 'MMM d, yyyy')}</span>
                </div>
              </div>
            </div>
          </div>
          
          {/* Ragged edge effect at bottom */}
          <div className="absolute -bottom-1 left-0 right-0 h-2 bg-[#fdfaf3] deckled-edge" />
        </div>
      </DialogTrigger>
      <DialogContent className="vellum-modal-content max-w-2xl">
        <DialogHeader>
          <DialogTitle className="font-tome-heading text-3xl text-primary border-b-2 border-primary/20 pb-2 mb-2">
            {note.title}
          </DialogTitle>
          <div className="flex flex-wrap gap-2 mb-2">
            {note.episodeTag && (
              <span className="inline-flex items-center gap-1 text-xs font-tome-marginalia px-2 py-0.5 rounded bg-primary/10 text-primary border border-primary/20">
                <BookOpen className="w-3 h-3" />
                {note.episodeTag}
              </span>
            )}
            {note.partyId && (
              <span className="inline-flex items-center gap-1 text-xs font-tome-marginalia px-2 py-0.5 rounded bg-amber-100 text-amber-800 border border-amber-300">
                <Lock className="w-3 h-3" />
                Party only
              </span>
            )}
          </div>
        </DialogHeader>
        <div className="space-y-6">
          <p className="font-serif leading-relaxed text-lg first-letter:text-5xl first-letter:font-bold first-letter:mr-2 first-letter:float-left first-letter:leading-[0.8] whitespace-pre-wrap">
            {note.description}
          </p>
          {note.pdfUrl && (
            <a
              href={note.pdfUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 text-primary font-tome-marginalia hover:underline"
            >
              <FileText className="w-4 h-4" />
              Open attached PDF
            </a>
          )}
          <div className="pt-6 border-t border-primary/10 flex items-center justify-between text-sm font-tome-marginalia opacity-70 italic">
            <div className="flex items-center gap-2">
              <User className="w-4 h-4" />
              <span>Reported by: {note.username}</span>
            </div>
            <div className="flex items-center gap-2">
              <Calendar className="w-4 h-4" />
              <span>{format(new Date(note.createdAt), 'MMMM d, yyyy')}</span>
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
