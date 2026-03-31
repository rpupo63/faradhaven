import { useState, useEffect } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { useAuth, DICE_THEME_COLOR_DEFAULT, DICE_FONT_COLOR_DEFAULT, DICE_THEME_DEFAULT } from '@/context/AuthContext';
import { updateCharacter } from '@/lib/api';
import { useToast } from '@/hooks/use-toast';
import type { DicePrefs } from '@/lib/api';

// Preset themes shipped with @3d-dice/dice-box
const THEMES = [
  { id: 'default', label: 'Classic' },
  { id: 'rust', label: 'Rust' },
  { id: 'diceOfRolling', label: 'Dice of Rolling' },
  { id: 'gemstone', label: 'Gemstone' },
] as const;

interface DiceCustomizerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function DiceCustomizer({ open, onOpenChange }: DiceCustomizerProps) {
  const {
    user,
    token,
    activeCharacterId,
    characterDicePrefs,
    updateUserDicePreferences,
    setDicePrefsOverride,
  } = useAuth();

  const characterId = activeCharacterId ?? undefined;
  const { toast } = useToast();

  // User-level form state
  const [userTheme, setUserTheme] = useState(user?.dice_theme ?? DICE_THEME_DEFAULT);
  const [userThemeColor, setUserThemeColor] = useState(user?.dice_theme_color ?? DICE_THEME_COLOR_DEFAULT);
  const [userFontColor, setUserFontColor] = useState(user?.dice_font_color ?? DICE_FONT_COLOR_DEFAULT);

  // Character override form state
  const [charOverrideEnabled, setCharOverrideEnabled] = useState(
    !!(characterDicePrefs?.dice_theme || characterDicePrefs?.dice_theme_color || characterDicePrefs?.dice_font_color)
  );
  const [charTheme, setCharTheme] = useState(characterDicePrefs?.dice_theme ?? DICE_THEME_DEFAULT);
  const [charThemeColor, setCharThemeColor] = useState(characterDicePrefs?.dice_theme_color ?? DICE_THEME_COLOR_DEFAULT);
  const [charFontColor, setCharFontColor] = useState(characterDicePrefs?.dice_font_color ?? DICE_FONT_COLOR_DEFAULT);

  const [saving, setSaving] = useState(false);

  // Sync form state when user/characterDicePrefs change (e.g. after a save)
  useEffect(() => {
    setUserTheme(user?.dice_theme ?? DICE_THEME_DEFAULT);
    setUserThemeColor(user?.dice_theme_color ?? DICE_THEME_COLOR_DEFAULT);
    setUserFontColor(user?.dice_font_color ?? DICE_FONT_COLOR_DEFAULT);
  }, [user?.dice_theme, user?.dice_theme_color, user?.dice_font_color]);

  useEffect(() => {
    const hasOverride = !!(
      characterDicePrefs?.dice_theme ||
      characterDicePrefs?.dice_theme_color ||
      characterDicePrefs?.dice_font_color
    );
    setCharOverrideEnabled(hasOverride);
    setCharTheme(characterDicePrefs?.dice_theme ?? DICE_THEME_DEFAULT);
    setCharThemeColor(characterDicePrefs?.dice_theme_color ?? DICE_THEME_COLOR_DEFAULT);
    setCharFontColor(characterDicePrefs?.dice_font_color ?? DICE_FONT_COLOR_DEFAULT);
  }, [characterDicePrefs]);

  // Preview: push current form values into the live dice config while dialog is open
  useEffect(() => {
    if (!open) return;
    const preview: DicePrefs = charOverrideEnabled && characterId
      ? { dice_theme: charTheme, dice_theme_color: charThemeColor, dice_font_color: charFontColor }
      : { dice_theme: userTheme, dice_theme_color: userThemeColor, dice_font_color: userFontColor };
    setDicePrefsOverride(preview);
    return () => setDicePrefsOverride(null);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, charOverrideEnabled, characterId, charTheme, charThemeColor, charFontColor, userTheme, userThemeColor, userFontColor]);

  const handleSave = async () => {
    setSaving(true);
    try {
      // Always save user-level prefs
      await updateUserDicePreferences({
        dice_theme: userTheme,
        dice_theme_color: userThemeColor,
        dice_font_color: userFontColor,
      });

      // Save character overrides if applicable
      if (characterId && token) {
        if (charOverrideEnabled) {
          await updateCharacter(
            characterId,
            { dice_theme: charTheme, dice_theme_color: charThemeColor, dice_font_color: charFontColor },
            token
          );
        } else {
          // Clear character overrides
          await updateCharacter(
            characterId,
            // Use the backend's clear flags
            { clear_dice_theme: true, clear_dice_colors: true } as never,
            token
          );
        }
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
          {/* ── User Defaults ── */}
          <section className="space-y-4">
            <h3 className="text-xs font-tome-marginalia uppercase tracking-widest text-muted-foreground">
              Your Defaults
            </h3>
            <DiceAppearanceControls
              theme={userTheme}
              themeColor={userThemeColor}
              fontColor={userFontColor}
              onThemeChange={setUserTheme}
              onThemeColorChange={setUserThemeColor}
              onFontColorChange={setUserFontColor}
            />
          </section>

          {/* ── Character Override ── */}
          {characterId && (
            <>
              <Separator />
              <section className="space-y-4">
                <div className="flex items-center justify-between">
                  <h3 className="text-xs font-tome-marginalia uppercase tracking-widest text-muted-foreground">
                    This Character
                  </h3>
                  <button
                    type="button"
                    onClick={() => setCharOverrideEnabled(v => !v)}
                    className={`relative inline-flex h-5 w-9 items-center rounded-full transition-colors focus:outline-none ${
                      charOverrideEnabled ? 'bg-primary' : 'bg-muted'
                    }`}
                  >
                    <span
                      className={`inline-block h-3.5 w-3.5 transform rounded-full bg-white shadow transition-transform ${
                        charOverrideEnabled ? 'translate-x-4' : 'translate-x-1'
                      }`}
                    />
                  </button>
                </div>
                {charOverrideEnabled && (
                  <DiceAppearanceControls
                    theme={charTheme}
                    themeColor={charThemeColor}
                    fontColor={charFontColor}
                    onThemeChange={setCharTheme}
                    onThemeColorChange={setCharThemeColor}
                    onFontColorChange={setCharFontColor}
                  />
                )}
                {!charOverrideEnabled && (
                  <p className="text-xs text-muted-foreground">
                    Using your default dice appearance.
                  </p>
                )}
              </section>
            </>
          )}

          {/* Preview indicator */}
          <p className="text-[11px] text-muted-foreground/60 font-tome-marginalia">
            Changes preview live — roll a die to see them in action.
          </p>

          <div className="flex justify-end gap-2 pt-1">
            <Button variant="ghost" onClick={() => onOpenChange(false)} disabled={saving}>
              Cancel
            </Button>
            <Button onClick={handleSave} disabled={saving}>
              {saving ? 'Saving…' : 'Save'}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}

// ── Shared appearance controls ──────────────────────────────────────────────

interface DiceAppearanceControlsProps {
  theme: string;
  themeColor: string;
  fontColor: string;
  onThemeChange: (v: string) => void;
  onThemeColorChange: (v: string) => void;
  onFontColorChange: (v: string) => void;
}

function DiceAppearanceControls({
  theme,
  themeColor,
  fontColor,
  onThemeChange,
  onThemeColorChange,
  onFontColorChange,
}: DiceAppearanceControlsProps) {
  return (
    <div className="space-y-4">
      {/* Theme preset selector */}
      <div className="space-y-2">
        <Label className="text-xs text-muted-foreground">Preset Theme</Label>
        <div className="grid grid-cols-2 gap-2">
          {THEMES.map(t => (
            <button
              key={t.id}
              type="button"
              onClick={() => onThemeChange(t.id)}
              className={`px-3 py-2 rounded-lg border text-xs font-tome-marginalia transition-colors ${
                theme === t.id
                  ? 'border-primary bg-primary/15 text-primary'
                  : 'border-border hover:border-primary/40 text-muted-foreground hover:text-foreground'
              }`}
            >
              {t.label}
            </button>
          ))}
        </div>
      </div>

      {/* Die body color */}
      <div className="flex items-center gap-3">
        <div className="flex-1 space-y-1">
          <Label className="text-xs text-muted-foreground">Die Color</Label>
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

        {/* Number color */}
        <div className="flex-1 space-y-1">
          <Label className="text-xs text-muted-foreground">Number Color</Label>
          <div className="flex items-center gap-2">
            <input
              type="color"
              value={fontColor}
              onChange={e => onFontColorChange(e.target.value)}
              className="h-8 w-8 cursor-pointer rounded border border-border bg-transparent p-0.5"
            />
            <input
              type="text"
              value={fontColor}
              onChange={e => {
                const v = e.target.value;
                if (/^#[0-9A-Fa-f]{0,6}$/.test(v)) onFontColorChange(v);
              }}
              className="flex-1 h-8 rounded border border-border bg-background px-2 text-xs font-mono"
              maxLength={7}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
