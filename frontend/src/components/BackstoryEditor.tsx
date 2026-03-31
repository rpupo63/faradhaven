import { useCallback, useEffect, useState } from 'react';
import { useEditor, EditorContent } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import Placeholder from '@tiptap/extension-placeholder';
import Underline from '@tiptap/extension-underline';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { getBackstory, updateBackstory } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import {
  Bold,
  Italic,
  Underline as UnderlineIcon,
  List,
  ListOrdered,
  Quote,
  Redo,
  Undo,
  Save,
  Loader2,
  BookOpen,
  Palette,
  Check,
} from 'lucide-react';

interface BackstoryEditorProps {
  characterId: string;
  token?: string;
  characterName?: string;
}

const COLOR_OPTIONS = [
  { name: 'Black', value: '#000000' },
  { name: 'Dark Charcoal', value: '#1a1a1a' },
  { name: 'Ink', value: 'var(--ink)' },
  { name: 'Dried Blood', value: 'var(--dried-blood)' },
  { name: 'Faded Gold', value: 'var(--faded-gold)' },
  { name: 'Charcoal', value: 'var(--charcoal)' },
  { name: 'Fire', value: 'hsl(var(--element-fire))' },
  { name: 'Ice', value: 'hsl(var(--element-ice))' },
  { name: 'Lightning', value: 'hsl(var(--element-lightning))' },
  { name: 'Nature', value: 'hsl(var(--element-nature))' },
  { name: 'Dark', value: 'hsl(var(--element-dark))' },
  { name: 'Navy', value: '#000080' },
  { name: 'Midnight', value: '#191970' },
  { name: 'Maroon', value: '#800000' },
  { name: 'Forest', value: '#006400' },
  { name: 'Teal', value: '#008080' },
  { name: 'Indigo', value: '#4b0082' },
  { name: 'Purple', value: '#800080' },
  { name: 'Brown', value: '#8b4513' },
  { name: 'Sienna', value: '#a0522d' },
];

export function BackstoryEditor({ characterId, token, characterName }: BackstoryEditorProps) {
  const queryClient = useQueryClient();
  const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false);
  const [backstoryColor, setBackstoryColor] = useState('#000000');

  // Fetch current backstory
  const { data: backstoryData, isLoading } = useQuery({
    queryKey: ['backstory', characterId],
    queryFn: () => getBackstory(characterId, token),
    enabled: !!characterId && !!token,
    staleTime: 30_000,
  });

  // Save mutation
  const saveMutation = useMutation({
    mutationFn: (html: string) => updateBackstory(characterId, html, token, backstoryColor),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['backstory', characterId] });
      setHasUnsavedChanges(false);
    },
  });

  // Initialize TipTap editor
  const editor = useEditor({
    extensions: [
      StarterKit.configure({
        heading: {
          levels: [2, 3],
        },
      }),
      Underline,
      Placeholder.configure({
        placeholder: 'Write your character\'s backstory here...\n\nTell the tale of their origins, motivations, and the events that shaped who they are.',
      }),
    ],
    content: '',
    editorProps: {
      attributes: {
        class: 'prose prose-sm max-w-none min-h-[300px] p-4 focus:outline-none',
      },
    },
    onUpdate: () => {
      setHasUnsavedChanges(true);
    },
  });

          const resetUnsavedChanges = useCallback(() => {
            setTimeout(() => setHasUnsavedChanges(false), 0);
          }, []);  
        // Update editor content when backstory data loads
        useEffect(() => {
          if (editor && backstoryData?.backstory !== undefined) {
            const currentEditorContent = editor.getHTML();
            const newBackstoryContent = backstoryData.backstory || '';
  
            // If the editor content doesn't match the new backstory, update it
            if (currentEditorContent !== newBackstoryContent) {
              editor.commands.setContent(newBackstoryContent);
              resetUnsavedChanges(); // Call the memoized reset function
            }
            
                  if (backstoryData.backstory_hex_color) {
                    setTimeout(() => setBackstoryColor(backstoryData.backstory_hex_color), 0);
                  }          }
        }, [editor, backstoryData?.backstory, backstoryData?.backstory_hex_color, resetUnsavedChanges]);
  const handleSave = useCallback(() => {
    if (!editor) return;
    const html = editor.getHTML();
    // Don't save if content is just the empty paragraph TipTap creates
    const isEmpty = html === '<p></p>' || html === '';
    saveMutation.mutate(isEmpty ? '' : html);
  }, [editor, saveMutation]);

  const handleColorChange = (color: string) => {
    setBackstoryColor(color);
    setHasUnsavedChanges(true);
  };

  if (!token) {
    return (
      <Card className="arcane-border bg-card">
        <CardContent className="py-12 text-center">
          <BookOpen className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
          <h3 className="text-lg font-tome-heading text-primary mb-2">Login Required</h3>
          <p className="text-muted-foreground font-tome-marginalia">
            Please log in to write and save your character's backstory.
          </p>
        </CardContent>
      </Card>
    );
  }

  if (isLoading) {
    return (
      <Card className="arcane-border bg-card">
        <CardContent className="py-12 text-center">
          <Loader2 className="h-8 w-8 mx-auto text-primary animate-spin mb-4" />
          <p className="text-muted-foreground font-tome-marginalia">Loading backstory...</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="arcane-border bg-card">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2 text-lg font-tome-heading text-primary">
            <BookOpen className="h-5 w-5" />
            {characterName ? `${characterName}'s Backstory` : 'Backstory'}
          </CardTitle>
          <div className="flex items-center gap-2">
            {hasUnsavedChanges && (
              <span className="text-xs text-amber-500 font-tome-marginalia">
                Unsaved changes
              </span>
            )}
            <Button
              size="sm"
              onClick={handleSave}
              disabled={!hasUnsavedChanges || saveMutation.isPending}
              className="gap-2"
            >
              {saveMutation.isPending ? (
                <>
                  <Loader2 className="h-4 w-4 animate-spin" />
                  Saving...
                </>
              ) : (
                <>
                  <Save className="h-4 w-4" />
                  Save
                </>
              )}
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        {/* Toolbar */}
        <div className="flex flex-wrap gap-1 p-2 rounded-md border border-border bg-muted/30">
          <ToolbarButton
            onClick={() => editor?.chain().focus().toggleBold().run()}
            isActive={editor?.isActive('bold')}
            title="Bold"
          >
            <Bold className="h-4 w-4" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor?.chain().focus().toggleItalic().run()}
            isActive={editor?.isActive('italic')}
            title="Italic"
          >
            <Italic className="h-4 w-4" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor?.chain().focus().toggleUnderline().run()}
            isActive={editor?.isActive('underline')}
            title="Underline"
          >
            <UnderlineIcon className="h-4 w-4" />
          </ToolbarButton>

          <div className="w-px h-6 bg-border mx-1 self-center" />

          <ToolbarButton
            onClick={() => editor?.chain().focus().toggleHeading({ level: 2 }).run()}
            isActive={editor?.isActive('heading', { level: 2 })}
            title="Heading"
          >
            <span className="font-bold text-sm">H2</span>
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor?.chain().focus().toggleHeading({ level: 3 }).run()}
            isActive={editor?.isActive('heading', { level: 3 })}
            title="Subheading"
          >
            <span className="font-bold text-xs">H3</span>
          </ToolbarButton>

          <div className="w-px h-6 bg-border mx-1 self-center" />

          <ToolbarButton
            onClick={() => editor?.chain().focus().toggleBulletList().run()}
            isActive={editor?.isActive('bulletList')}
            title="Bullet List"
          >
            <List className="h-4 w-4" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor?.chain().focus().toggleOrderedList().run()}
            isActive={editor?.isActive('orderedList')}
            title="Numbered List"
          >
            <ListOrdered className="h-4 w-4" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor?.chain().focus().toggleBlockquote().run()}
            isActive={editor?.isActive('blockquote')}
            title="Quote"
          >
            <Quote className="h-4 w-4" />
          </ToolbarButton>

          <div className="w-px h-6 bg-border mx-1 self-center" />
          
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="ghost" size="sm" className="h-8 w-8 p-0" title="Text Color">
                <Palette className="h-4 w-4" style={{ color: backstoryColor }} />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-64 p-2">
              <div className="grid grid-cols-5 gap-1">
                {COLOR_OPTIONS.map((color) => (
                  <button
                    key={color.name}
                    className="w-8 h-8 rounded-full border border-border flex items-center justify-center hover:scale-110 transition-transform"
                    style={{ background: color.value }}
                    onClick={() => handleColorChange(color.value)}
                    title={color.name}
                  >
                    {backstoryColor === color.value && (
                      <Check className="h-4 w-4 text-white drop-shadow-md" />
                    )}
                  </button>
                ))}
              </div>
            </PopoverContent>
          </Popover>

          <div className="w-px h-6 bg-border mx-1 self-center" />

          <ToolbarButton
            onClick={() => editor?.chain().focus().undo().run()}
            disabled={!editor?.can().undo()}
            title="Undo"
          >
            <Undo className="h-4 w-4" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor?.chain().focus().redo().run()}
            disabled={!editor?.can().redo()}
            title="Redo"
          >
            <Redo className="h-4 w-4" />
          </ToolbarButton>
        </div>

        {/* Editor */}
        <div 
          className="rounded-md border border-border bg-background/50 min-h-[350px]"
          style={{ 
            color: backstoryColor,
            '--tw-prose-body': backstoryColor,
            '--tw-prose-headings': backstoryColor,
            '--tw-prose-lead': backstoryColor,
            '--tw-prose-links': backstoryColor,
            '--tw-prose-bold': backstoryColor,
            '--tw-prose-counters': backstoryColor,
            '--tw-prose-bullets': backstoryColor,
            '--tw-prose-hr': backstoryColor,
            '--tw-prose-quotes': backstoryColor,
            '--tw-prose-quote-borders': backstoryColor,
            '--tw-prose-captions': backstoryColor,
            '--tw-prose-code': backstoryColor,
            '--tw-prose-pre-code': backstoryColor,
            '--tw-prose-pre-bg': backstoryColor,
            '--tw-prose-th-borders': backstoryColor,
            '--tw-prose-td-borders': backstoryColor,
          } as React.CSSProperties}
        >
          <EditorContent editor={editor} />
        </div>

        {/* Status messages */}
        {saveMutation.isError && (
          <p className="text-sm text-red-500 text-center">
            {saveMutation.error?.message || 'Failed to save backstory'}
          </p>
        )}
        {saveMutation.isSuccess && !hasUnsavedChanges && (
          <p className="text-sm text-green-500 text-center">
            Backstory saved successfully!
          </p>
        )}
      </CardContent>
    </Card>
  );
}

interface ToolbarButtonProps {
  onClick: () => void;
  isActive?: boolean;
  disabled?: boolean;
  title?: string;
  children: React.ReactNode;
}

function ToolbarButton({ onClick, isActive, disabled, title, children }: ToolbarButtonProps) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      title={title}
      className={`p-2 rounded hover:bg-accent/50 transition-colors ${
        isActive ? 'bg-accent text-accent-foreground' : 'text-muted-foreground'
      } ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
    >
      {children}
    </button>
  );
}
