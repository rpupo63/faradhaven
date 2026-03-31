import { useQuery } from '@tanstack/react-query';
import { getClasses } from '@/lib/api';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

const NONE_VALUE = '__none__';

export interface MonsterGenerationClassSelectProps {
  value: string;
  onValueChange: (className: string) => void;
  token?: string;
  disabled?: boolean;
  id?: string;
}

/**
 * Optional Faradhaven class picker for AI monster generation.
 * Uses the same /api/classes list as character creation; the backend matches by class name to seed JSON.
 */
export function MonsterGenerationClassSelect({
  value,
  onValueChange,
  token,
  disabled,
  id = 'monster-gen-class',
}: MonsterGenerationClassSelectProps) {
  const { data: classes = [], isLoading } = useQuery({
    queryKey: ['classes'],
    queryFn: () => getClasses(token),
    enabled: !!token,
  });

  const sorted = [...classes].sort((a, b) => a.name.localeCompare(b.name));

  return (
    <div className="grid gap-2">
      <Label htmlFor={id}>Enemy theme (optional)</Label>
      <Select
        value={value ? value : NONE_VALUE}
        onValueChange={(v) => onValueChange(v === NONE_VALUE ? '' : v)}
        disabled={disabled || isLoading || !token}
      >
        <SelectTrigger id={id} className="w-full">
          <SelectValue placeholder={isLoading ? 'Loading classes…' : 'Choose source'} />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value={NONE_VALUE}>Custom creature (from description only)</SelectItem>
          {sorted.map((c) => (
            <SelectItem key={c.id} value={c.name}>
              {c.name}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      <p className="text-xs text-muted-foreground">
        Pick a Faradhaven class to base the enemy on seeded class features and flavor; leave as custom for a normal monster from your description alone.
      </p>
    </div>
  );
}
