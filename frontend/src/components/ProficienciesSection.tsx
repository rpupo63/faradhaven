import { cn } from '@/lib/utils';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { NormalizedCharacterSheet } from '@/types/game';
import { ScrollText } from 'lucide-react';

interface ProficienciesSectionProps {
  sheet: NormalizedCharacterSheet;
  className?: string;
}

/**
 * Displays class proficiencies (weapons, armor), tools, and skill focus.
 * Shown when class data includes proficiencies (e.g. from API).
 */
export function ProficienciesSection({ sheet, className }: ProficienciesSectionProps) {
  const { class: cls } = sheet;
  const hasProficiencies = cls.proficiencies || (cls.tools && cls.tools.length > 0) || (cls.skill_focus && cls.skill_focus.length > 0);

  if (!hasProficiencies) return null;

  return (
    <Card className={cn('arcane-border bg-card', className)}>
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center gap-2 text-base font-tome-subheading text-primary">
          <ScrollText className="h-4 w-4" />
          Proficiencies
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        {cls.proficiencies && (
          <div>
            <p className="text-xs text-muted-foreground font-tome-marginalia mb-1">Weapons & Armor</p>
            <p className="text-sm font-tome-marginalia">{cls.proficiencies}</p>
          </div>
        )}
        {cls.tools && cls.tools.length > 0 && (
          <div>
            <p className="text-xs text-muted-foreground font-tome-marginalia mb-1">Tools</p>
            <p className="text-sm font-tome-marginalia">{cls.tools.join(', ')}</p>
          </div>
        )}
        {cls.skill_focus && cls.skill_focus.length > 0 && (
          <div>
            <p className="text-xs text-muted-foreground font-tome-marginalia mb-1">Skill Focus</p>
            <p className="text-sm font-tome-marginalia">{cls.skill_focus.join(', ')}</p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
