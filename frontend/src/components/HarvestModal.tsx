import React, { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { toast } from '@/components/ui/use-toast';
import { Loader2 } from 'lucide-react';

// Define types for HarvestableAbility and HarvestableAbilitiesResponse
// These should ideally come from a shared types file if the frontend and backend share a codebase
interface HarvestableAbility {
  id: string;
  name: string;
  type: 'skill' | 'attack' | 'recipe';
  description: string;
}

interface HarvestableAbilitiesResponse {
  beast_id: string; // uuid
  abilities: HarvestableAbility[];
  fracture_dc: number;
  requires_save: boolean;
}

interface ConfirmHarvestRequest {
  character_id: string;
  beast_id: string;
  ability_id: string;
  ability_type: string;
  did_fail_save: boolean;
}

// API Calls (Assuming these exist in lib/api.ts or similar)
const getHarvestableAbilitiesApi = async (
  characterId: string,
  beastId: string,
  token: string
): Promise<HarvestableAbilitiesResponse> => {
  const response = await fetch(`/api/beasts/${beastId}/harvestable-abilities?characterID=${characterId}`, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
  });
  if (!response.ok) {
    throw new Error('Failed to fetch harvestable abilities');
  }
  return response.json();
};

const confirmHarvestApi = async (data: ConfirmHarvestRequest, token: string): Promise<unknown> => {
  const response = await fetch(`/api/characters/${data.character_id}/harvest`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  });
  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || 'Failed to confirm harvest');
  }
  return response.json();
};

interface HarvestModalProps {
  isOpen: boolean;
  onClose: () => void;
  characterId: string;
  token: string;
  onHarvestSuccess?: () => void;
}

export const HarvestModal: React.FC<HarvestModalProps> = ({
  isOpen,
  onClose,
  characterId,
  token,
  onHarvestSuccess,
}) => {
  const queryClient = useQueryClient();
  const [beastIdInput, setBeastIdInput] = useState('');
  const [selectedAbility, setSelectedAbility] = useState<HarvestableAbility | null>(null);
  const [didFailSave, setDidFailSave] = useState<'yes' | 'no' | null>(null);

  // Query to fetch harvestable abilities from a beast
  const {
    data: harvestData,
    isLoading: isLoadingHarvestData,
    refetch,
    isError: isErrorHarvestData,
    error: harvestError,
  } = useQuery<HarvestableAbilitiesResponse, Error>({
    queryKey: ['harvestableAbilities', beastIdInput, characterId],
    queryFn: () => getHarvestableAbilitiesApi(characterId, beastIdInput, token),
    enabled: !!beastIdInput && !!characterId && !!token,
    staleTime: 0,
    retry: false,
  });

  // Mutation to confirm the harvest
  const confirmHarvestMutation = useMutation<unknown, Error, ConfirmHarvestRequest>({
    mutationFn: (data) => confirmHarvestApi(data, token),
    onSuccess: (dataFromApi) => {
      toast({
        title: 'Harvest Successful!',
        description: `${selectedAbility?.name} harvested.`,
      });
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      onHarvestSuccess?.();
      onClose();
    },
    onError: (error) => {
      toast({
        title: 'Harvest Failed',
        description: error.message || 'Could not harvest the ability.',
        variant: 'destructive',
      });
    },
  });

  // Reset state when modal opens/closes
  useEffect(() => {
    if (!isOpen) {
      setTimeout(() => {
        setBeastIdInput('');
        setSelectedAbility(null);
        setDidFailSave(null);
      }, 0);
    }
  }, [isOpen]);

  const handleFetchAbilities = () => {
    if (beastIdInput) {
      refetch();
    }
  };

  const handleConfirmHarvest = () => {
    if (!selectedAbility || !harvestData || didFailSave === null) {
      toast({
        title: 'Missing Information',
        description: 'Please select an ability and indicate if you failed the save.',
        variant: 'destructive',
      });
      return;
    }

    confirmHarvestMutation.mutate({
      character_id: characterId,
      beast_id: harvestData.beast_id,
      ability_id: selectedAbility.id,
      ability_type: selectedAbility.type,
      did_fail_save: didFailSave === 'yes',
    });
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Perform Visceral Psychometry</DialogTitle>
          <DialogDescription>
            Consume the liver of a recently deceased beast to absorb its abilities.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="beastId" className="text-right">
              Beast ID
            </Label>
            <Input
              id="beastId"
              value={beastIdInput}
              onChange={(e) => setBeastIdInput(e.target.value)}
              placeholder="Enter Beast UUID"
              className="col-span-3"
            />
          </div>
          <Button onClick={handleFetchAbilities} disabled={!beastIdInput || isLoadingHarvestData}>
            {isLoadingHarvestData && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Fetch Abilities
          </Button>

          {isErrorHarvestData && <p className="text-red-500 text-sm">{harvestError?.message || 'Error fetching abilities.'}</p>}

          {harvestData && (
            <>
              <div className="grid gap-2">
                <Label>Harvestable Abilities:</Label>
                {harvestData.abilities.length === 0 ? (
                  <p className="text-muted-foreground text-sm">No abilities found for this beast.</p>
                ) : (
                  <Select
                    onValueChange={(value) => {
                      const ability = harvestData.abilities.find((a) => a.id === value);
                      setSelectedAbility(ability || null);
                    }}
                    value={selectedAbility?.id || ''}
                  >
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Select an ability to harvest" />
                    </SelectTrigger>
                    <SelectContent>
                      {harvestData.abilities.map((ability) => (
                        <SelectItem key={ability.id} value={ability.id}>
                          {ability.name} ({ability.type}) - {ability.description}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              </div>

              {harvestData.requires_save && (
                <div className="grid gap-2">
                  <Label>Fracture Save (DC {harvestData.fracture_dc}): Did you fail?</Label>
                  <RadioGroup value={didFailSave || ''} onValueChange={(value: 'yes' | 'no') => setDidFailSave(value)}>
                    <div className="flex items-center space-x-2">
                      <RadioGroupItem value="yes" id="fail-yes" />
                      <Label htmlFor="fail-yes">Yes</Label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <RadioGroupItem value="no" id="fail-no" />
                      <Label htmlFor="fail-no">No</Label>
                    </div>
                  </RadioGroup>
                </div>
              )}
            </>
          )}
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button
            onClick={handleConfirmHarvest}
            disabled={!selectedAbility || confirmHarvestMutation.isPending || (harvestData?.requires_save && didFailSave === null)}
          >
            {confirmHarvestMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Confirm Harvest
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
