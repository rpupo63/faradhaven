import { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { setCharacterParty } from '@/lib/api';
import { ApiCharacter, ApiCharacterSheet, ApiParty } from '@/types/game/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { User, UserMinus, UserPlus, Users } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { createPartyAndJoinWithCharacter, getPartiesByUserId, addCharacterToParty } from '@/lib/api/party';

interface PartyMembersProps {
  characterId: string;
  token?: string;
  characterSheet?: ApiCharacterSheet;
}

export function PartyMembers({ characterId, token, characterSheet }: PartyMembersProps) {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [showCreatePartyDialog, setShowCreatePartyDialog] = useState(false);
  const [newPartyName, setNewPartyName] = useState('');
  const [showJoinPartyDialog, setShowJoinPartyDialog] = useState(false);
  const [selectedPartyId, setSelectedPartyId] = useState<string>('');

  const currentCharacter = characterSheet?.character;
  const party = currentCharacter?.party;
  const members: ApiCharacter[] = party?.members || [];
  const isPartyOwner = party?.owner_id === currentCharacter?.user_id;

  const { data: userParties, isLoading: isLoadingParties } = useQuery({
    queryKey: ['user-parties', currentCharacter?.user_id],
    queryFn: () => currentCharacter?.user_id ? getPartiesByUserId(currentCharacter.user_id, token) : Promise.resolve([]),
    enabled: !!currentCharacter?.user_id && showJoinPartyDialog,
  });

  const createPartyMutation = useMutation({
    mutationFn: async () => {
      if (!token) throw new Error("Authentication token not available.");
      if (!characterId) throw new Error("Character ID not available.");
      if (!newPartyName.trim()) throw new Error("Party name cannot be empty.");

      await createPartyAndJoinWithCharacter(newPartyName, characterId, token);
      return true;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      setShowCreatePartyDialog(false);
      setNewPartyName('');
    },
    onError: (err) => {
      console.error("Failed to create party:", err);
      alert("Failed to create party: " + err.message);
    }
  });

  const joinPartyMutation = useMutation({
    mutationFn: async (partyId: string) => {
      if (!token) throw new Error("Authentication token not available.");
      if (!characterId) throw new Error("Character ID not available.");

      await addCharacterToParty(partyId, characterId, token);
      return true;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      setShowJoinPartyDialog(false);
      setSelectedPartyId('');
    },
    onError: (err) => {
      console.error("Failed to join party:", err);
      alert("Failed to join party: " + err.message);
    }
  });

  const leavePartyMutation = useMutation({
    mutationFn: (partyId: string) => setCharacterParty(characterId, null, token),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      // Optionally navigate away or show a message
    },
    onError: (err) => {
      console.error("Failed to leave party:", err);
      alert("Failed to leave party: " + err.message);
    }
  });

  if (!party) {
    return (
      <>
        <Alert>
          <AlertTitle>Not in a Party</AlertTitle>
          <AlertDescription>
            This character is not currently part of any party.
          </AlertDescription>
          <div className="mt-4 flex gap-2">
            <Button onClick={() => setShowCreatePartyDialog(true)}>
              <Users className="h-4 w-4 mr-2" /> Create New Party
            </Button>
            <Button variant="outline" onClick={() => setShowJoinPartyDialog(true)}>
              <UserPlus className="h-4 w-4 mr-2" /> Join Existing Party
            </Button>
          </div>
        </Alert>

        <Dialog open={showCreatePartyDialog} onOpenChange={setShowCreatePartyDialog}>
          <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>Create New Party</DialogTitle>
              <DialogDescription>
                Enter a name for your new party. Your current character will automatically join it.
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-1 gap-2 sm:grid-cols-4 sm:items-center sm:gap-4">
                <Label htmlFor="partyName" className="text-left sm:text-right">
                  Party Name
                </Label>
                <Input
                  id="partyName"
                  value={newPartyName}
                  onChange={(e) => setNewPartyName(e.target.value)}
                  className="sm:col-span-3 min-w-0"
                  disabled={createPartyMutation.isPending}
                />
              </div>
            </div>
            <DialogFooter>
              <Button
                type="submit"
                onClick={() => createPartyMutation.mutate()}
                disabled={!newPartyName.trim() || createPartyMutation.isPending}
              >
                {createPartyMutation.isPending ? 'Creating...' : 'Create Party'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        <Dialog open={showJoinPartyDialog} onOpenChange={setShowJoinPartyDialog}>
          <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>Join Existing Party</DialogTitle>
              <DialogDescription>
                Select one of your existing parties to join with this character.
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-1 gap-2">
                <Label htmlFor="partySelect">Select Party</Label>
                {isLoadingParties ? (
                  <p className="text-sm text-muted-foreground italic">Loading your parties...</p>
                ) : userParties && userParties.length > 0 ? (
                  <Select value={selectedPartyId} onValueChange={setSelectedPartyId}>
                    <SelectTrigger id="partySelect">
                      <SelectValue placeholder="Select a party..." />
                    </SelectTrigger>
                    <SelectContent>
                      {userParties.map((p) => (
                        <SelectItem key={p.id} value={p.id}>
                          {p.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                ) : (
                  <p className="text-sm text-muted-foreground italic">No other parties found. You can create a new one instead.</p>
                )}
              </div>
            </div>
            <DialogFooter>
              <Button
                type="submit"
                onClick={() => joinPartyMutation.mutate(selectedPartyId)}
                disabled={!selectedPartyId || joinPartyMutation.isPending}
              >
                {joinPartyMutation.isPending ? 'Joining...' : 'Join Party'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </>
    );
  }

  return (
    <div className="space-y-4">
      <h3 className="text-2xl font-bold font-display text-primary glow-text">
        {party.name} <span className="text-muted-foreground text-lg">({members.length} Members)</span>
      </h3>
      <p className="text-muted-foreground">
        Current members of your adventuring party.
      </p>

      <div className="flex gap-2">
        <Button
          variant="outline"
          onClick={() => leavePartyMutation.mutate(party.id)}
          disabled={leavePartyMutation.isPending}
        >
          <UserMinus className="h-4 w-4 mr-2" /> Leave Party
        </Button>
        {isPartyOwner && (
          <>


          </>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {members.map((member) => (
          <Card key={member.id} className="group hover:border-primary transition-colors">
            <CardHeader className="flex flex-row items-center space-x-4 p-4">
              <div className="w-12 h-12 rounded-full overflow-hidden border border-primary/20 bg-primary/5 shrink-0">
                {member.image_url ? (
                  <img src={member.image_url} alt={member.name} className="w-full h-full object-cover" />
                ) : (
                  <User className="w-8 h-8 text-primary/40" />
                )}
              </div>
              <div className="flex-grow min-w-0">
                <CardTitle className="text-lg leading-tight truncate">{member.name}</CardTitle>
                <CardDescription className="text-xs font-tome-marginalia">
                  Lvl {member.level} {member.race?.name} {member.class?.name}
                </CardDescription>
              </div>
            </CardHeader>
            <CardContent className="px-4 pb-4 flex justify-between items-center">

                {isPartyOwner && member.id !== characterId && ( // Can't remove self through this interface
                    <Button
                        variant="destructive"
                        size="sm"
                        // TODO: Implement remove member mutation
                        // onClick={() => removeMemberMutation.mutate(member.id)}
                    >
                        Remove
                    </Button>
                )}
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}