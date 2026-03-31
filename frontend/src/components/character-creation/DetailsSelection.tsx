import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { Checkbox } from '@/components/ui/checkbox';
import { Separator } from '@/components/ui/separator';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Badge } from '@/components/ui/badge';
import { ApiClass, ApiRace } from '@/types/game';

// All available languages in Faradhaven (standard + exotic)
const ALL_LANGUAGES = [
  'Common', 'Dwarvish', 'Elvish', 'Giant', 'Gnomish', 'Goblin', 'Halfling', 'Orc',
  'Abyssal', 'Celestial', 'Draconic', 'Deep Speech', 'Infernal', 'Primordial', 'Quori', 'Sylvan', 'Undercommon',
];

interface DetailsSelectionProps {
  cls: ApiClass;
  race: ApiRace | null;
  name: string;
  setName: (n: string) => void;
  skills: string[];
  setSkills: (s: string[]) => void;
  selectedLanguages: string[];
  setSelectedLanguages: (l: string[]) => void;
  equipment: Record<string, string>;
  setEquipment: (e: Record<string, string> | ((prev: Record<string, string>) => Record<string, string>)) => void;
}

export function DetailsSelection({ cls, race, name, setName, skills, setSkills, selectedLanguages, setSelectedLanguages, equipment, setEquipment }: DetailsSelectionProps) {
  const toggleSkill = (skill: string) => {
    if (skills.includes(skill)) {
      setSkills(skills.filter(s => s !== skill));
    } else {
      if (skills.length < (cls.skill_choice_count || 2)) {
        setSkills([...skills, skill]);
      }
    }
  };

  const handleEquipChange = (choiceId: string, optionId: string) => {
    setEquipment((prev: Record<string, string>) => ({ ...prev, [choiceId]: optionId }));
  };

  return (
    <div className="space-y-8 h-full overflow-y-auto pr-4">
      <div className="space-y-2">
        <Label htmlFor="char-name" className="text-lg">Character Name</Label>
        <Input 
          id="char-name" 
          value={name} 
          onChange={(e) => setName(e.target.value)} 
          placeholder="Enter a name..." 
          className="text-lg font-display h-12"
        />
      </div>

      <Separator />

      <div className="space-y-4">
        <div className="flex justify-between items-center">
          <Label className="text-lg">Skill Proficiencies</Label>
          <span className="text-sm text-muted-foreground">
            Choose {skills.length} / {cls.skill_choice_count || 2}
          </span>
        </div>
        <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
          {Array.from(new Set(cls.skill_choice || []))
            .filter(skill => !cls.skill_focus?.some(f => f.toLowerCase() === skill.toLowerCase()))
            .map(skill => (
              <div key={skill} className="flex items-center space-x-2">
                <Checkbox 
                  id={`skill-${skill}`} 
                  checked={skills.includes(skill)}
                  onCheckedChange={() => toggleSkill(skill)}
                  disabled={!skills.includes(skill) && skills.length >= (cls.skill_choice_count || 2)}
                />
                <label 
                  htmlFor={`skill-${skill}`} 
                  className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
                >
                  {skill}
                </label>
              </div>
            ))}
        </div>
        {cls.skill_focus && cls.skill_focus.length > 0 && (
          <div className="text-sm text-muted-foreground">
            <span className="font-bold">Fixed Proficiencies:</span> {cls.skill_focus.join(", ")}
          </div>
        )}
      </div>

      {/* Language Selection */}
      {race && (
        <>
          <Separator />
          <div className="space-y-4">
            <Label className="text-lg">Languages</Label>
            
            {/* Known languages from race */}
            {race.languages && race.languages.length > 0 && (
              <div className="space-y-2">
                <div className="text-xs uppercase text-muted-foreground">Known from Race</div>
                <div className="flex flex-wrap gap-2">
                  {race.languages.map(lang => (
                    <Badge key={lang} variant="secondary" font="tomeMarginalia">
                      {lang}
                    </Badge>
                  ))}
                </div>
              </div>
            )}

            {/* Bonus language selection */}
            {(race.bonus_language_count ?? 0) > 0 && (
              <div className="space-y-2">
                <div className="flex justify-between items-center">
                  <div className="text-xs uppercase text-muted-foreground">Choose Additional Languages</div>
                  <span className="text-sm text-muted-foreground">
                    {selectedLanguages.length} / {race.bonus_language_count}
                  </span>
                </div>
                <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                  {ALL_LANGUAGES
                    .filter(lang => !(race.languages ?? []).includes(lang))
                    .map(lang => (
                      <div key={lang} className="flex items-center space-x-2">
                        <Checkbox
                          id={`lang-${lang}`}
                          checked={selectedLanguages.includes(lang)}
                          onCheckedChange={() => {
                            if (selectedLanguages.includes(lang)) {
                              setSelectedLanguages(selectedLanguages.filter(l => l !== lang));
                            } else if (selectedLanguages.length < (race.bonus_language_count ?? 0)) {
                              setSelectedLanguages([...selectedLanguages, lang]);
                            }
                          }}
                          disabled={!selectedLanguages.includes(lang) && selectedLanguages.length >= (race.bonus_language_count ?? 0)}
                        />
                        <label
                          htmlFor={`lang-${lang}`}
                          className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
                        >
                          {lang}
                        </label>
                      </div>
                    ))}
                </div>
              </div>
            )}

            {/* No bonus languages message */}
            {(race.bonus_language_count ?? 0) === 0 && (
              <p className="text-sm text-muted-foreground font-tome-marginalia italic">
                Your race does not grant additional language choices.
              </p>
            )}
          </div>
        </>
      )}

      <Separator />

      <div className="space-y-6">
        <Label className="text-lg">Starting Equipment</Label>
        
        {/* Automatic Equipment */}
        {cls.starting_equipment && cls.starting_equipment.length > 0 && (
          <div className="bg-muted/20 p-3 rounded border">
            <div className="text-xs uppercase text-muted-foreground mb-1">Included</div>
            <div className="text-sm">{cls.starting_equipment.join(", ")}</div>
          </div>
        )}

        {/* Choices */}
        {cls.equipment_choices?.map((choice, idx) => (
          <div key={choice.id} className="border p-4 rounded-lg space-y-3">
            <Label className="text-base text-primary">Choice {idx + 1}: {choice.instruction}</Label>
            <RadioGroup 
              value={equipment[choice.id] || ''} 
              onValueChange={(val) => handleEquipChange(choice.id, val)}
            >
              {choice.options.map(opt => (
                <div key={opt.id} className="flex items-center space-x-2">
                  <RadioGroupItem value={opt.id} id={opt.id} />
                  <Label htmlFor={opt.id} className="font-normal cursor-pointer">
                    {opt.description}
                  </Label>
                </div>
              ))}
            </RadioGroup>
          </div>
        ))}
      </div>
    </div>
  );
}
