import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { PreparedSpells } from '@/components/PreparedSpells';
import { restSpellPreparationHint } from '@/lib/restSpellPreparation';

interface RestSpellPreparationDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  restKind: 'short' | 'long';
  characterId: string;
  token: string;
  characterName: string;
  /** e.g. "The Elixirist" from sheet.class.name */
  gameClassName: string | undefined;
}

export function RestSpellPreparationDialog({
  open,
  onOpenChange,
  restKind,
  characterId,
  token,
  characterName,
  gameClassName,
}: RestSpellPreparationDialogProps) {
  const titleRest = restKind === 'short' ? 'Short rest' : 'Long rest';
  const hint = restSpellPreparationHint(gameClassName, restKind);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto gap-4">
        <DialogHeader className="text-left space-y-2">
          <DialogTitle className="font-tome-heading text-primary">
            {titleRest}: adjust prepared spells
          </DialogTitle>
          <DialogDescription className="text-muted-foreground font-tome-marginalia">
            {hint}
          </DialogDescription>
        </DialogHeader>
        <div className="min-h-0">
          <PreparedSpells characterId={characterId} token={token} characterName={characterName} />
        </div>
      </DialogContent>
    </Dialog>
  );
}
