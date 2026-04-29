import { UserCircle, Users } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { MoneyPanel } from './MoneyPanel';
import type { NormalizedCharacterSheet } from '@/types/game';
import { cn } from '@/lib/utils';
import {
  getChangelingPersonaStorageKey,
  isChangelingSheet,
  shouldShowCharacterSpecificCard,
} from '@/lib/characterSpecificInfo';

export interface CharacterSpecificInfoCardProps {
  sheet: NormalizedCharacterSheet;
  onMoneyChange?: (newTotal: number) => void | Promise<void>;
  className?: string;
}

/**
 * Character-specific sheet details: party, Changeling persona (local), currency.
 * Visibility rules live in `@/lib/characterSpecificInfo`.
 */
export function CharacterSpecificInfoCard({
  sheet,
  onMoneyChange,
  className,
}: CharacterSpecificInfoCardProps) {
  if (!shouldShowCharacterSpecificCard(sheet)) {
    return null;
  }

  const characterId = sheet.character.id;
  const partyName = sheet.character.partyName?.trim();
  const changeling = isChangelingSheet(sheet);
  const showMoney = sheet.money !== undefined && sheet.money !== null;

  return (
    <Card className={cn('arcane-border bg-card', className)}>
      <CardHeader className="pb-2 pt-3 px-3 md:px-4">
        <CardTitle className="text-sm font-tome-subheading text-primary flex items-center gap-2">
          <Users className="h-4 w-4 shrink-0 opacity-90" />
          Character details
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-3 px-3 pb-3 pt-0 md:px-4 md:pb-4">
        {partyName ? (
          <div className="flex items-start gap-2 rounded-md border border-border/50 bg-muted/15 px-2.5 py-2">
            <Users className="h-4 w-4 shrink-0 text-primary mt-0.5" />
            <div className="min-w-0">
              <p className="text-micro font-tome-marginalia text-muted-foreground uppercase">Party</p>
              <p className="text-sm font-tome-marginalia text-foreground truncate" title={partyName}>
                {partyName}
              </p>
            </div>
          </div>
        ) : null}

        {changeling ? (
          <div className="flex items-center gap-2">
            <UserCircle className="h-4 w-4 text-primary shrink-0" />
            <span className="text-xs font-tome-marginalia text-muted-foreground uppercase shrink-0">
              Form
            </span>
            <Input
              placeholder="Current persona…"
              className="h-8 text-sm bg-transparent border-border/60 focus-visible:ring-1 font-tome-marginalia flex-1 min-w-0"
              defaultValue={typeof localStorage !== 'undefined' ? localStorage.getItem(getChangelingPersonaStorageKey(characterId)) || '' : ''}
              onChange={(e) => {
                try {
                  localStorage.setItem(getChangelingPersonaStorageKey(characterId), e.target.value);
                } catch {
                  /* ignore quota / private mode */
                }
              }}
            />
          </div>
        ) : null}

        {showMoney ? <MoneyPanel sheet={sheet} onMoneyChange={onMoneyChange} /> : null}
      </CardContent>
    </Card>
  );
}
