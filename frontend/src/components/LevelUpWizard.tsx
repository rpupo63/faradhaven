import { useState, useMemo, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { LoadingButton } from '@/components/ui/loading-button';
import { LoadingQuill } from '@/components/LoadingQuill';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { Sparkles, X, ChevronRight, ChevronLeft } from 'lucide-react';
import {
  getLevelUpPreview,
  levelUpCharacter,
  type LevelUpRequest,
} from '@/lib/api';
import { dispatchClearDice } from '@/lib/dice';

// Sub-components
import { LevelPreview } from './level-up/LevelPreview';
import { ArchetypeSelection } from './level-up/ArchetypeSelection';
import { WeaponSelection } from './level-up/WeaponSelection';
import { HPGainSelection } from './level-up/HPGainSelection';
import { ASIAllocation } from './level-up/ASIAllocation';
import { LevelUpConfirm } from './level-up/LevelUpConfirm';
import { ModularChoiceSelection } from './level-up/ModularChoiceSelection';

interface LevelUpWizardProps {
  characterId: string;
  token?: string;
  onComplete: () => void;
  onCancel: () => void;
}

type WizardStep = 'preview' | 'archetype' | 'weapon' | 'hp' | 'asi' | 'modular_choice' | 'confirm';
type HPChoice = 'average' | 'roll';

export function LevelUpWizard({
  characterId,
  token,
  onComplete,
  onCancel,
}: LevelUpWizardProps) {
  const queryClient = useQueryClient();
  const [step, setStep] = useState<WizardStep>('preview');

  // Clear dice when wizard unmounts
  useEffect(() => {
    return () => {
      dispatchClearDice();
    };
  }, []);

  // Choices state
  const [asiAllocation, setAsiAllocation] = useState<Record<string, number>>({});
  const [hpChoice, setHpChoice] = useState<HPChoice>('average');
  const [hpRollResult, setHpRollResult] = useState<number | null>(null);
  const [selectedArchetypeId, setSelectedArchetypeId] = useState<string | null>(null);
  const [selectedWeaponId, setSelectedWeaponId] = useState<string | null>(null);
  const [modularChoice, setModularChoice] = useState<{ mp: number; br: number } | null>(null);

  // Fetch preview data
  const { data: preview, isLoading } = useQuery({
    queryKey: ['levelUpPreview', characterId],
    queryFn: () => getLevelUpPreview(characterId, token),
  });

  // Level-up mutation
  const levelUpMutation = useMutation({
    mutationFn: (request: LevelUpRequest) =>
      levelUpCharacter(characterId, request, token),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['characterSheet', characterId] });
      queryClient.invalidateQueries({ queryKey: ['levelHistory', characterId] });
      onComplete();
    },
  });

  const totalASISpent = useMemo(
    () => Object.values(asiAllocation).reduce((sum, v) => sum + v, 0),
    [asiAllocation]
  );

  const asiPointsRemaining = (preview?.asi_points_available ?? 0) - totalASISpent;

  const handleASIChange = (ability: string, delta: number) => {
    setAsiAllocation((prev) => {
      const current = prev[ability] ?? 0;
      const newValue = Math.max(0, Math.min(2, current + delta));
      if (delta > 0 && asiPointsRemaining <= 0) return prev;
      return { ...prev, [ability]: newValue };
    });
  };

  const steps: WizardStep[] = useMemo(() => {
    const result: WizardStep[] = ['preview'];
    if (preview?.requires_archetype_choice) {
      result.push('archetype');
    }
    if (preview?.weapon_selection_info) {
      result.push('weapon');
    }
    if (preview?.class_level?.level_features?.some(f => f.name.startsWith('Modular Choice'))) {
      result.push('modular_choice');
    }
    if (preview && preview.next_level > 1) {
      result.push('hp');
    }
    if (preview?.asi_points_available && preview.asi_points_available > 0) {
      result.push('asi');
    }
    result.push('confirm');
    return result;
  }, [preview]);

  const hpGain = useMemo(() => {
    if (!preview) return 0;
    if (hpChoice === 'average') {
      return preview.hp_gain_average;
    }
    if (hpRollResult !== null) {
      const total = hpRollResult + preview.con_mod;
      return Math.max(1, total);
    }
    return 0;
  }, [preview, hpChoice, hpRollResult]);

  const currentStepIndex = steps.indexOf(step);
  const progress = ((currentStepIndex + 1) / steps.length) * 100;

  const goNext = () => {
    const nextIndex = currentStepIndex + 1;
    if (nextIndex < steps.length) {
      setStep(steps[nextIndex]);
    }
  };

  const goBack = () => {
    const prevIndex = currentStepIndex - 1;
    if (prevIndex >= 0) {
      setStep(steps[prevIndex]);
    }
  };

  const handleConfirm = () => {
    const request: LevelUpRequest = {};
    if (Object.keys(asiAllocation).length > 0) {
      request.asi_allocation = asiAllocation;
    }
    if (hpChoice === 'roll' && hpRollResult !== null) {
      request.hp_roll_result = hpRollResult;
    }
    if (selectedArchetypeId) {
      request.archetype_id = selectedArchetypeId;
    }
    if (selectedWeaponId) {
      request.primary_weapon_id = selectedWeaponId;
    }
    if (modularChoice) {
      request.mp_change = modularChoice.mp;
      request.br_change = modularChoice.br;
    }
    levelUpMutation.mutate(request);
  };

  const selectedArchetype = useMemo(() => {
    if (!selectedArchetypeId || !preview?.available_archetypes) return null;
    return preview.available_archetypes.find(a => a.id === selectedArchetypeId) || null;
  }, [selectedArchetypeId, preview]);

  const selectedWeapon = useMemo(() => {
    if (!selectedWeaponId || !preview?.weapon_selection_info?.eligible_weapons) return null;
    return preview.weapon_selection_info.eligible_weapons.find(w => w.id === selectedWeaponId) || null;
  }, [selectedWeaponId, preview]);

  if (isLoading) {
    return (
      <Card className="w-full max-w-2xl mx-auto arcane-border">
        <CardContent className="p-8">
          <LoadingQuill label="Loading level-up preview..." />
        </CardContent>
      </Card>
    );
  }

  if (!preview) {
    return (
      <Card className="w-full max-w-2xl mx-auto arcane-border">
        <CardContent className="p-8 text-center">
          <p className="text-destructive">Failed to load level-up data</p>
          <Button onClick={onCancel} className="mt-4">Close</Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="w-full max-w-2xl mx-auto arcane-border bg-card">
      <CardHeader className="border-b border-border/50">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2 font-tome-subheading text-primary">
            <Sparkles className="h-5 w-5" />
            Level Up to {preview.next_level}
          </CardTitle>
          <Button variant="ghost" size="icon" onClick={onCancel}>
            <X className="h-4 w-4" />
          </Button>
        </div>
        <Progress value={progress} className="mt-4" />
      </CardHeader>

      <CardContent className="p-6 space-y-6">
        {step === 'preview' && <LevelPreview preview={preview} />}
        
        {step === 'archetype' && (
          <ArchetypeSelection
            preview={preview}
            selectedArchetypeId={selectedArchetypeId}
            onSelect={setSelectedArchetypeId}
          />
        )}

        {step === 'weapon' && (
          <WeaponSelection
            preview={preview}
            selectedWeaponId={selectedWeaponId}
            onSelect={setSelectedWeaponId}
          />
        )}

        {step === 'modular_choice' && (
          <ModularChoiceSelection 
            preview={preview}
            onSelect={setModularChoice}
          />
        )}

        {step === 'hp' && (
          <HPGainSelection 
            preview={preview}
            hpChoice={hpChoice}
            setHpChoice={setHpChoice}
            hpRollResult={hpRollResult}
            setHpRollResult={setHpRollResult}
          />
        )}

        {step === 'asi' && (
          <ASIAllocation 
            preview={preview}
            asiAllocation={asiAllocation}
            onASIChange={handleASIChange}
            asiPointsRemaining={asiPointsRemaining}
          />
        )}

        {step === 'confirm' && (
          <LevelUpConfirm
            preview={preview}
            hpGain={hpGain}
            hpChoice={hpChoice}
            hpRollResult={hpRollResult}
            selectedArchetype={selectedArchetype}
            selectedWeapon={selectedWeapon}
            asiAllocation={asiAllocation}
          />
        )}

        <div className="flex justify-between pt-4 border-t border-border/50">
          <Button 
            variant="outline" 
            onClick={goBack} 
            disabled={currentStepIndex === 0 || levelUpMutation.isPending}
            className="gap-2"
          >
            <ChevronLeft className="h-4 w-4" />
            Back
          </Button>
          
          {step === 'confirm' ? (
            <LoadingButton
              onClick={handleConfirm}
              isLoading={levelUpMutation.isPending}
              className="gap-2 bg-primary-glow"
              loadingText="Processing..."
            >
              Confirm Level Up
            </LoadingButton>
          ) : (
            <Button
              onClick={goNext}
              disabled={
                (step === 'archetype' && !selectedArchetypeId) ||
                (step === 'weapon' && !selectedWeaponId) ||
                (step === 'modular_choice' && !modularChoice) ||
                (step === 'hp' && hpChoice === 'roll' && hpRollResult === null)
              }
              className="gap-2"
            >
              Next
              <ChevronRight className="h-4 w-4" />
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
}