// frontend/src/lib/api/madness.ts
import { getBaseUrl } from '@/lib/api/base';

// Define types mirroring backend services/madness.go
export interface MadnessRollResult {
  character_id: string; // uuid
  roll: number;
  madness_die: number;
  effect?: string;
}

export const rollMadness = async (characterId: string, token: string): Promise<MadnessRollResult> => {
  const response = await fetch(`${getBaseUrl()}/api/characters/${characterId}/madness/roll`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || 'Failed to roll madness die.');
  }

  return response.json();
};