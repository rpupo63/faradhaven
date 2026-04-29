import { useState, useEffect } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import {
  useAuth,
  DICE_THEME_COLOR_DEFAULT,
  DICE_FONT_COLOR_DEFAULT,
  DICE_THEME_DEFAULT,
} from '@/context/AuthContext';
import { updateCharacter } from '@/lib/api';
import { useToast } from '@/hooks/use-toast';
import type { DicePrefs } from '@/lib/api';

interface DiceCustomizerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function DiceCustomizer({ open, onOpenChange }: DiceCustomizerProps) {
  const queryClient = useQueryClient();
  const { token, activeCharacterId, characterDicePrefs, setDicePreview, setCharacterDicePrefs } = useAuth();

  const characterId = activeCharacterId ?? undefined;
  const { toast } = useToast();

  const [charThemeColor, setCharThemeColor] = useState(
    characterDicePrefs?.dice_theme_color ?? DICE_THEME_COLOR_DEFAULT
  );

  const [saving, setSaving] = useState(false);

  useEffect(() => {
    setCharThemeColor(characterDicePrefs?.dice_theme_color ?? DICE_THEME_COLOR_DEFAULT);
  }, [characterDicePrefs]);

  // Preview while open for the active character only.
  useEffect(() => {
    if (!open || !characterId) return;
    const preview: DicePrefs = {
      dice_theme: DICE_THEME_DEFAULT,
      dice_theme_color: charThemeColor,
      dice_font_color: DICE_FONT_COLOR_DEFAULT,
    };
    setDicePreview(preview);
    return () => setDicePreview(null);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, characterId, charThemeColor]);

  const handleSave = async () => {
    setSaving(true);
    try {
      if (characterId && token) {
        await updateCharacter(
          characterId,
          {
            dice_theme: DICE_THEME_DEFAULT,
            dice_theme_color: charThemeColor,
          },
          token
        );

        setCharacterDicePrefs({
          dice_theme: DICE_THEME_DEFAULT,
          dice_theme_color: charThemeColor,
          dice_font_color: DICE_FONT_COLOR_DEFAULT,
        });
        void queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      }

      toast({ title: 'Dice preferences saved' });
      onOpenChange(false);
    } catch {
      toast({ title: 'Failed to save preferences', variant: 'destructive' });
    } finally {
      setSaving(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle className="font-display tracking-wide">Dice Appearance</DialogTitle>
        </DialogHeader>

        <div className="space-y-6 pt-2">
          {characterId && (
            <section className="space-y-4">
              <p className="text-sm text-muted-foreground">
                Dice appearance applies to your active character.
              </p>
              <CharacterDiceColors
                themeColor={charThemeColor}
                onThemeColorChange={setCharThemeColor}
              />
            </section>
          )}

          {!characterId && (
            <>
              <Separator />
              <p className="text-xs text-muted-foreground">
                Set an active character to customize that character&apos;s dice colors.
              </p>
            </>
          )}

          <p className="text-fine text-muted-foreground/60 font-tome-marginalia">
            Changes preview live — roll a die to see them in action.
          </p>

          <div className="flex justify-end gap-2 pt-1">
            <Button variant="ghost" onClick={() => onOpenChange(false)} disabled={saving}>
              Cancel
            </Button>
            <Button onClick={handleSave} disabled={saving || !characterId}>
              {saving ? 'Saving…' : 'Save'}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}

function CharacterDiceColors({
  themeColor,
  onThemeColorChange,
}: {
  themeColor: string;
  onThemeColorChange: (v: string) => void;
}) {
  return (
    <div className="space-y-1">
      <Label className="text-xs text-muted-foreground">Die color</Label>
      <div className="flex items-center gap-2">
        <input
          type="color"
          value={themeColor}
          onChange={e => onThemeColorChange(e.target.value)}
          className="h-8 w-8 cursor-pointer rounded border border-border bg-transparent p-0.5"
        />
        <input
          type="text"
          value={themeColor}
          onChange={e => {
            const v = e.target.value;
            if (/^#[0-9A-Fa-f]{0,6}$/.test(v)) onThemeColorChange(v);
          }}
          className="flex-1 h-8 rounded border border-border bg-background px-2 text-xs font-mono"
          maxLength={7}
        />
      </div>
    </div>
  );
}
