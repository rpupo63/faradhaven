import { useState, useMemo, useEffect } from 'react';
import { Character, DND5E_SKILLS } from '@/types/game';
import {
  getSkillChoiceFromClass,
  getSkillFocusFromClass,
} from '@/lib/classLevelData';
import type { ApiClass } from '@/types/game';
import { Button } from './ui/button';
import { TomeInput } from './ui/TomeInput';
import { Label } from './ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './ui/select';
import { Slider } from './ui/slider';
import { Checkbox } from './ui/checkbox';
import { X } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useQuery } from '@tanstack/react-query';
import { getClasses, getRaces } from '@/lib/api';

function arraysEqual<T>(a: T[], b: T[]): boolean {
  if (a.length !== b.length) return false;
  for (let i = 0; i < a.length; i++) {
    if (a[i] !== b[i]) return false;
  }
  return true;
}

interface CharacterFormProps {
  character?: Character;
  classes?: ApiClass[];
  onSubmit: (data: {
    name: string;
    race: string;
    race_id?: string;
    class_id: string;
    class_name: string;
    level: number;
    skill_proficiencies?: string[];
  }) => void;
  onCancel: () => void;
}

/** Single-page character creation form with all fields visible */
export function CharacterForm({
  character,
  classes: classesProp,
  onSubmit,
  onCancel,
}: CharacterFormProps) {
  const { data: classesFromApi } = useQuery({
    queryKey: ['classes'],
    queryFn: () => getClasses(),
    staleTime: 60_000,
    retry: false,
  });
  const { data: racesFromApi } = useQuery({
    queryKey: ['races'],
    queryFn: () => getRaces(),
    staleTime: 60_000,
    retry: false,
  });
  const classes = useMemo(() => classesProp ?? classesFromApi ?? [], [classesProp, classesFromApi]);
  const races = useMemo(() => racesFromApi ?? [], [racesFromApi]);

  const [name, setName] = useState(character?.name || '');
  const [selectedRace, setSelectedRace] = useState<{ id: string; name: string } | null>(() => {
    if (character?.race_id && races.length) {
      const r = races.find((x) => x.id === character.race_id);
      return r ? { id: r.id, name: r.name } : null;
    }
    if (character?.race && races.length) {
      const r = races.find((x) => x.name === character.race);
      return r ? { id: r.id, name: r.name } : { id: '', name: character.race };
    }
    return races[0] ? { id: races[0].id, name: races[0].name } : null;
  });
  const [selectedClass, setSelectedClass] = useState<ApiClass | null>(() => {
    if (character?.class_id && classes.length) {
      return classes.find((c) => c.id === character.class_id) ?? null;
    }
    if (character?.class && classes.length) {
      return classes.find((c) => c.name === character.class) ?? null;
    }
    return classes[0] ?? null;
  });
  const [skillProficiencies, setSkillProficiencies] = useState<string[]>(() => {
    if (character?.skill_proficiencies?.length) return character.skill_proficiencies;
    const c = selectedClass ?? classes[0];
    return c ? getSkillFocusFromClass(c) : [];
  });
  const [level, setLevel] = useState(character?.level || 1);

  const skillChoiceCount = selectedClass?.skill_choice_count ?? 2;
  const skillChoiceOptions = useMemo(
    () => (selectedClass ? getSkillChoiceFromClass(selectedClass) : []),
    [selectedClass]
  );

  useEffect(() => {
    if (!selectedClass) return;
    const options = getSkillChoiceFromClass(selectedClass);
    const defaults = getSkillFocusFromClass(selectedClass);
    setTimeout(() => {
      setSkillProficiencies(prevSkillProficiencies => {
        const valid = prevSkillProficiencies.filter((id) => options.includes(id));
        if (valid.length === skillChoiceCount) return valid;
        const toAdd = defaults.filter((id) => options.includes(id) && !valid.includes(id));
        const newProficiencies = [...valid, ...toAdd].slice(0, skillChoiceCount);

        // Only return new state if it's different to prevent unnecessary re-renders
        if (arraysEqual(prevSkillProficiencies, newProficiencies)) {
          return prevSkillProficiencies; // No change
        }
        return newProficiencies;
      });
    }, 0);
  }, [selectedClass, skillChoiceCount]);

  useEffect(() => {
    if (classes.length && !selectedClass) {
      setTimeout(() => setSelectedClass(classes[0]), 0);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [classes]);

  useEffect(() => {
    if (races.length && !selectedRace) {
      setTimeout(() => setSelectedRace({ id: races[0].id, name: races[0].name }), 0);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [races]);

  const toggleSkill = (skillId: string) => {
    setSkillProficiencies((prev) => {
      if (prev.includes(skillId)) {
        return prev.filter((id) => id !== skillId);
      }
      if (prev.length >= skillChoiceCount) {
        return [skillId, ...prev.slice(0, -1)];
      }
      return [...prev, skillId];
    });
  };

  const canSubmit = () => {
    return (
      name.trim().length > 0 &&
      !!selectedClass &&
      !!selectedRace &&
      skillProficiencies.length === skillChoiceCount
    );
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (canSubmit() && selectedClass && selectedRace) {
      onSubmit({
        name: name.trim(),
        race: selectedRace.name,
        race_id: selectedRace.id || undefined,
        class_id: selectedClass.id,
        class_name: selectedClass.name,
        level,
        skill_proficiencies: skillProficiencies,
      });
    }
  };

  return (
    <div className="w-full max-w-2xl mx-auto">
      <div className="vellum-modal-content bg-card p-6 md:p-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-8 pb-4 border-b border-ink/20">
          <h2 className="font-tome-heading text-xl text-primary">
            {character ? 'Edit Character' : 'Create Character'}
          </h2>
          <Button variant="ghost" size="icon" onClick={onCancel} type="button">
            <X className="w-4 h-4" />
          </Button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-8">
          {/* Basic Info Section */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Name */}
            <div className="space-y-2">
              <Label htmlFor="name" className="font-tome-subheading text-base">
                Character Name
              </Label>
              <TomeInput
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Enter name..."
                autoFocus
              />
            </div>

            {/* Level */}
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <Label className="font-tome-subheading text-base">Level</Label>
                <span className="font-tome-marginalia text-primary text-lg font-semibold">
                  {level}
                </span>
              </div>
              <Slider
                value={[level]}
                onValueChange={([v]) => setLevel(v)}
                min={1}
                max={20}
                step={1}
                className="[&_[role=slider]]:bg-primary"
              />
            </div>
          </div>

          {/* Race & Class Section */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Race */}
            <div className="space-y-2">
              <Label className="font-tome-subheading text-base">Race</Label>
              <Select
                value={selectedRace?.id ?? ''}
                onValueChange={(id) => {
                  const r = races.find((x) => x.id === id);
                  if (r) setSelectedRace({ id: r.id, name: r.name });
                }}
              >
                <SelectTrigger className="tome-input-underline border-0 border-b-2 rounded-none bg-transparent font-tome-marginalia">
                  <SelectValue placeholder="Select race..." />
                </SelectTrigger>
                <SelectContent>
                  {races.map((r) => (
                    <SelectItem key={r.id} value={r.id}>
                      {r.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {races.length === 0 && (
                <p className="text-xs text-muted-foreground">
                  Races load from the server. Ensure the backend is running and seeded.
                </p>
              )}
            </div>

            {/* Class */}
            <div className="space-y-2">
              <Label className="font-tome-subheading text-base">Class</Label>
              <Select
                value={selectedClass?.id ?? ''}
                onValueChange={(id) => {
                  const c = classes.find((x) => x.id === id);
                  if (c) setSelectedClass(c);
                }}
              >
                <SelectTrigger className="tome-input-underline border-0 border-b-2 rounded-none bg-transparent font-tome-marginalia">
                  <SelectValue placeholder="Select class..." />
                </SelectTrigger>
                <SelectContent>
                  {classes.map((c) => (
                    <SelectItem key={c.id} value={c.id}>
                      {c.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {classes.length === 0 && (
                <p className="text-xs text-muted-foreground">
                  Classes load from the server. Ensure the backend is running and seeded.
                </p>
              )}
            </div>
          </div>

          {/* Skill Proficiencies Section */}
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <Label className="font-tome-subheading text-base">
                Skill Proficiencies
              </Label>
              <span className="text-sm text-muted-foreground font-tome-marginalia">
                {skillProficiencies.length} / {skillChoiceCount} selected
              </span>
            </div>
            <p className="text-sm text-muted-foreground font-tome-marginalia">
              Choose {skillChoiceCount} skills from your class options
            </p>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-2 max-h-64 overflow-y-auto p-1">
              {skillChoiceOptions.map((skillId) => {
                const skill = DND5E_SKILLS.find((s) => s.id === skillId);
                const isSelected = skillProficiencies.includes(skillId);
                return (
                  <label
                    key={skillId}
                    className={cn(
                      'flex items-center gap-2 rounded border px-3 py-2 cursor-pointer font-tome-marginalia text-sm transition-colors',
                      isSelected
                        ? 'border-primary bg-primary/10'
                        : 'border-border hover:border-primary/50 hover:bg-muted/50'
                    )}
                  >
                    <Checkbox
                      checked={isSelected}
                      onCheckedChange={() => toggleSkill(skillId)}
                    />
                    <span>{skill?.name ?? skillId}</span>
                  </label>
                );
              })}
            </div>
            {skillChoiceOptions.length === 0 && (
              <p className="text-sm text-muted-foreground italic">
                Select a class to see available skills
              </p>
            )}
          </div>

          {/* Form Actions */}
          <div className="flex gap-4 pt-6 border-t border-ink/20">
            <Button
              type="button"
              variant="outline"
              onClick={onCancel}
              className="flex-1"
            >
              Cancel
            </Button>
            <Button type="submit" className="flex-1" disabled={!canSubmit()}>
              {character ? 'Save Changes' : 'Create Character'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
