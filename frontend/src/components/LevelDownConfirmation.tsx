import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { AlertTriangle } from 'lucide-react';

interface LevelDownConfirmationProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  currentLevel: number;
  onConfirm: () => void;
  isPending?: boolean;
}

export function LevelDownConfirmation({
  open,
  onOpenChange,
  currentLevel,
  onConfirm,
  isPending,
}: LevelDownConfirmationProps) {
  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent className="arcane-border bg-card">
        <AlertDialogHeader>
          <AlertDialogTitle className="flex items-center gap-2 font-tome-subheading text-primary">
            <AlertTriangle className="h-5 w-5 text-destructive" />
            Revert to Level {currentLevel - 1}?
          </AlertDialogTitle>
          <AlertDialogDescription asChild>
            <div className="space-y-2 text-sm text-muted-foreground">
              <p className="font-tome-marginalia">
                This will restore your character to their exact state before
                leveling up to level {currentLevel}.
              </p>
              <p className="text-destructive font-medium">
                The following changes will be reverted:
              </p>
              <ul className="list-disc list-inside text-sm space-y-1">
                <li>Ability score improvements (if any)</li>
                <li>New spells learned</li>
                <li>New skill proficiencies</li>
                <li>Any features gained at level {currentLevel}</li>
              </ul>
            </div>
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={isPending}>Cancel</AlertDialogCancel>
          <AlertDialogAction
            onClick={onConfirm}
            disabled={isPending}
            className="bg-destructive hover:bg-destructive/90"
          >
            {isPending ? 'Reverting...' : 'Confirm Level Down'}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
