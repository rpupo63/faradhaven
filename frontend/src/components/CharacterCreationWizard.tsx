import { useState, useEffect } from 'react';
import { ApiClass, ApiRace, ApiCreationOptions, CreateCharacterRequest } from '@/types/game';
import { getCreationOptions, createCharacter } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';
import { Card, CardContent, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { LoadingButton } from '@/components/ui/loading-button';
import { LoadingQuill } from '@/components/LoadingQuill';
import { ChevronRight } from 'lucide-react';
import { cn } from '@/lib/utils';

// Sub-components
import { RaceSelection } from './character-creation/RaceSelection';
import { ClassSelection } from './character-creation/ClassSelection';
import { AbilityScoreSelection } from './character-creation/AbilityScoreSelection';
import { DetailsSelection } from './character-creation/DetailsSelection';
import { ReviewSummary } from './character-creation/ReviewSummary';

interface CharacterCreationWizardProps {
  onComplete: (characterId: string) => void;
  onCancel: () => void;
}

const STEPS = ['Race', 'Class', 'Abilities', 'Details', 'Review'];

export function CharacterCreationWizard({ onComplete, onCancel }: CharacterCreationWizardProps) {
  const { token, userId } = useAuth();
  const [step, setStep] = useState(0);
  const [loading, setLoading] = useState(true);
  const [isCreating, setIsCreating] = useState(false);
  const [options, setOptions] = useState<ApiCreationOptions | null>(null);

  // Form State
  const [name, setName] = useState('');
  const [selectedRace, setSelectedRace] = useState<ApiRace | null>(null);
  const [selectedLineageId, setSelectedLineageId] = useState<string | null>(null); // maps to TraitOption ID
  const [selectedClass, setSelectedClass] = useState<ApiClass | null>(null);
  const [abilityMethod, setAbilityMethod] = useState<'standard' | 'pointbuy' | 'roll'>('standard');
  const [abilities, setAbilities] = useState({ strength: 10, dexterity: 10, constitution: 10, intelligence: 10, wisdom: 10, charisma: 10 });
  const [skillProficiencies, setSkillProficiencies] = useState<string[]>([]);
  const [selectedLanguages, setSelectedLanguages] = useState<string[]>([]);
  const [equipmentChoices, setEquipmentChoices] = useState<Record<string, string>>({}); // choice_id -> option_id

  useEffect(() => {
    if (!token) return;

    getCreationOptions(token)
      .then((data) => {
        setOptions(data);
        setLoading(false);
      })
      .catch(err => {
        console.error("Failed to load options", err);
        setLoading(false);
      });
  }, [token]);

  const handleCreate = async () => {
    if (!selectedRace || !selectedClass || !userId || !token) {
      if (!userId) alert("User ID not found. Please log in again.");
      return;
    }

    // Calculate final scores with bonuses
    const getBonus = (ability: string) => {
      let bonus = 0;
      if (selectedRace?.ability_score_bonuses?.[ability]) bonus += selectedRace.ability_score_bonuses[ability];
      if (selectedLineageId && selectedRace?.traits) {
        for (const t of selectedRace.traits) {
          const opt = t.options?.find(o => o.id === selectedLineageId);
          if (opt?.ability_score_bonuses?.[ability]) bonus += opt.ability_score_bonuses[ability];
        }
      }
      return bonus;
    };

    const payload: CreateCharacterRequest = {
      user_id: userId,
      name,
      race_id: selectedRace.id,
      lineage_id: selectedLineageId || null,
      class_id: selectedClass.id,
      level: 1,
      spellbook: [],
      strength: abilities.strength + getBonus('strength'),
      dexterity: abilities.dexterity + getBonus('dexterity'),
      constitution: abilities.constitution + getBonus('constitution'),
      intelligence: abilities.intelligence + getBonus('intelligence'),
      wisdom: abilities.wisdom + getBonus('wisdom'),
      charisma: abilities.charisma + getBonus('charisma'),
      current_spell_points: 0,
      money: 5000, // Starting gold: 50 gp (5000 cp)
      skill_proficiencies: skillProficiencies,
      languages: [...(selectedRace?.languages ?? []), ...selectedLanguages],
      equipment_choices: Object.values(equipmentChoices).filter(Boolean)
    };

    try {
      setIsCreating(true);
      const data = await createCharacter(payload, token);
      onComplete(data.id);
    } catch (e: unknown) {
      console.error(e);
      alert((e as Error).message || "Failed to create character");
    } finally {
      setIsCreating(false);
    }
  };

  const handleNext = () => {
    if (step === STEPS.length - 1) {
      handleCreate();
    } else {
      setStep(s => Math.min(STEPS.length - 1, s + 1));
    }
  };
  const handleBack = () => setStep(s => Math.max(0, s - 1));

  if (loading) {
    return <LoadingQuill label="Loading character options..." />;
  }

  if (!options) return <div className="p-12 text-center text-destructive font-bold">Failed to load character options. Please check your connection and login status.</div>;

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Progress */}
      <div className="flex justify-between mb-8">
        {STEPS.map((s, i) => (
          <div key={s} className={cn("flex flex-col items-center gap-2", i <= step ? "text-primary" : "text-muted-foreground")}>
            <div className={cn("w-8 h-8 rounded-full flex items-center justify-center border-2 font-display", 
              i <= step ? "border-primary bg-primary/10" : "border-muted-foreground")}>
              {i + 1}
            </div>
            <span className="text-xs font-tome-marginalia uppercase">{s}</span>
          </div>
        ))}
      </div>

      <Card className="arcane-border bg-card min-h-[600px] flex flex-col">
        <CardHeader>
          <CardTitle className="font-display text-3xl">{STEPS[step]}</CardTitle>
        </CardHeader>
        <CardContent className="flex-1">
          {step === 0 && (
            <RaceSelection 
              options={options} 
              selected={selectedRace} 
              onSelect={(race) => { setSelectedRace(race); setSelectedLanguages([]); }} 
              lineageId={selectedLineageId}
              onSelectLineage={setSelectedLineageId}
            />
          )}
          {step === 1 && (
            <ClassSelection 
              options={options} 
              selected={selectedClass} 
              onSelect={(cls) => {
                setSelectedClass(cls);
                setSkillProficiencies([]);
                setEquipmentChoices({});
              }} 
            />
          )}
          {step === 2 && (
            <AbilityScoreSelection 
              method={abilityMethod} 
              setMethod={setAbilityMethod}
              scores={abilities}
              setScores={setAbilities}
              race={selectedRace}
              lineageId={selectedLineageId}
            />
          )}
          {step === 3 && selectedClass && (
            <DetailsSelection 
              cls={selectedClass}
              race={selectedRace}
              name={name}
              setName={setName}
              skills={skillProficiencies}
              setSkills={setSkillProficiencies}
              selectedLanguages={selectedLanguages}
              setSelectedLanguages={setSelectedLanguages}
              equipment={equipmentChoices}
              setEquipment={setEquipmentChoices}
            />
          )}
          {step === 4 && (
            <ReviewSummary 
              name={name}
              race={selectedRace}
              cls={selectedClass}
              abilities={abilities}
              languages={[...(selectedRace?.languages ?? []), ...selectedLanguages]}
            />
          )}
        </CardContent>
        <CardFooter className="justify-between border-t border-border p-6">
          <Button variant="ghost" onClick={step === 0 ? onCancel : handleBack} disabled={isCreating}>
            {step === 0 ? 'Cancel' : 'Back'}
          </Button>
          {step === STEPS.length - 1 ? (
            <LoadingButton
              onClick={handleNext}
              isLoading={isCreating}
              disabled={!isValid(step, {selectedRace, selectedClass, name, skillProficiencies})}
              loadingText="Creating..."
            >
              Create Character <ChevronRight className="ml-2 h-4 w-4" />
            </LoadingButton>
          ) : (
            <Button onClick={handleNext} disabled={!isValid(step, {selectedRace, selectedClass, name, skillProficiencies})}>
              Next <ChevronRight className="ml-2 h-4 w-4" />
            </Button>
          )}
        </CardFooter>
      </Card>
    </div>
  );
}

interface ValidationData {
  selectedRace: ApiRace | null;
  selectedClass: ApiClass | null;
  name: string;
  skillProficiencies: string[];
}

function isValid(step: number, data: ValidationData) {
  if (step === 0) return !!data.selectedRace;
  if (step === 1) return !!data.selectedClass;
  if (step === 3) {
    if (!data.name) return false;
    const needed = data.selectedClass?.skill_choice_count ?? 2;
    return data.skillProficiencies.length >= needed;
  }
  return true;
}