import type { CharacterSpellbookScope } from '@/lib/api';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { SCOPE_LABELS } from './characterSpellbookScopeConstants';

const SCOPES = Object.keys(SCOPE_LABELS) as CharacterSpellbookScope[];

export function CharacterSpellbookScopeFilter({
  value,
  onChange,
  id = 'spellbook-scope',
}: {
  value: CharacterSpellbookScope;
  onChange: (scope: CharacterSpellbookScope) => void;
  id?: string;
}) {
  return (
    <div className="space-y-2">
      <Label htmlFor={id} className="text-xs font-tome-marginalia text-muted-foreground">
        Show
      </Label>
      <Select value={value} onValueChange={(v) => onChange(v as CharacterSpellbookScope)}>
        <SelectTrigger id={id} className="w-full max-w-md bg-background">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {SCOPES.map((scope) => (
            <SelectItem key={scope} value={scope}>
              {SCOPE_LABELS[scope]}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
