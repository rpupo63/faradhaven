import React, { useState } from 'react';
import { GameMap, CreateTokenRequest, MapElementType, CreateMapElementRequest, MapToken } from '@/types/map'; // Import MapElementType
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { addToken, deleteMap, setInitiative, clearInitiative } from '@/lib/api/map';
import { useToast } from '@/hooks/use-toast';
import { useNavigate } from 'react-router-dom';
import { EntityCombobox, TokenEntity } from './EntityCombobox';
import { Checkbox } from '@/components/ui/checkbox'; // Assuming you have a Checkbox component
import { DiceManager } from '@/lib/dice-manager';
import { getAllCharacters } from '@/lib/api/character';
import { getMonstersByUser } from '@/lib/api/monster';

interface DMSidebarProps {
  mapData: GameMap;
  token: string;
  userId: string;
  refreshMap: () => void;
  onStartPlacement: (elementType: MapElementType, properties: Partial<CreateMapElementRequest>) => void; // New prop
  onStopPlacement: () => void; // New prop to stop placement mode
  placementModeActive: boolean; // New prop to indicate if placement mode is active
}

interface RollDetail {
  natural: number;
  modifier: number;
  total: number;
  dexterity: number;
}

export const DMSidebar: React.FC<DMSidebarProps> = ({ mapData, token, userId, refreshMap, onStartPlacement, onStopPlacement, placementModeActive }) => {
  const { toast } = useToast();
  const navigate = useNavigate();
  const [selectedEntity, setSelectedEntity] = useState<TokenEntity | null>(null);
  const [newTokenName, setNewTokenName] = useState('');
  const [newTokenColor, setNewTokenColor] = useState('#ff0000');

  // Initiative state
  const [rollDetails, setRollDetails] = useState<Record<string, RollDetail>>({});
  const [isRolling, setIsRolling] = useState(false);

  const initiativeOrder: MapToken[] = [...mapData.tokens]
    .filter(t => t.initiative_order != null)
    .sort((a, b) => (a.initiative_order ?? 0) - (b.initiative_order ?? 0));

  // State for Map Elements
  const [selectedElementType, setSelectedElementType] = useState<MapElementType | null>(null);
  const [elementColor, setElementColor] = useState('#6B8E2380');
  const [elementTrapName, setElementTrapName] = useState('');
  const [elementTrapDescription, setElementTrapDescription] = useState('');
  const [elementTrapVisible, setElementTrapVisible] = useState(false);
  const [elementElevationLevel, setElementElevationLevel] = useState(0);

  const handleEntitySelect = (entity: TokenEntity) => {
    setSelectedEntity(entity);
    setNewTokenName(entity.name);
  };

  const handleAddToken = async () => {
    if (!selectedEntity) {
      toast({ title: 'Select a character, monster, or party first', variant: 'destructive' });
      return;
    }

    if (selectedEntity.entity_type === 'party') {
      const members = selectedEntity.party_members ?? [];
      if (members.length === 0) {
        toast({ title: 'Party has no members', variant: 'destructive' });
        return;
      }
      const results = await Promise.allSettled(
        members.map((member) =>
          addToken(mapData.id, {
            name: member.name,
            image_url: member.image_url || '',
            token_type: 'pc',
            grid_x: 0,
            grid_y: 0,
            size: 1,
            color: newTokenColor,
            visible: true,
            character_id: member.id,
            assigned_user_id: member.owner_user_id,
          }, token)
        )
      );
      const failed = results.filter((r) => r.status === 'rejected').length;
      const succeeded = results.length - failed;
      if (succeeded > 0) toast({ title: `${succeeded} party member${succeeded !== 1 ? 's' : ''} added to map` });
      if (failed > 0) toast({ title: `${failed} token${failed !== 1 ? 's' : ''} failed to add`, variant: 'destructive' });
      setSelectedEntity(null);
      setNewTokenName('');
      refreshMap();
      return;
    }

    try {
      const req: CreateTokenRequest = {
        name: newTokenName || selectedEntity.name,
        image_url: selectedEntity.image_url || '',
        token_type: selectedEntity.entity_type === 'character' ? 'pc' : 'npc',
        grid_x: 0,
        grid_y: 0,
        size: 1,
        color: newTokenColor,
        visible: true,
        ...(selectedEntity.entity_type === 'character'
          ? { character_id: selectedEntity.id, assigned_user_id: selectedEntity.owner_user_id }
          : { monster_id: selectedEntity.id }),
      };
      await addToken(mapData.id, req, token);
      toast({ title: 'Token added' });
      setSelectedEntity(null);
      setNewTokenName('');
      refreshMap();
    } catch (err) {
      toast({ title: 'Failed to add token', variant: 'destructive' });
    }
  };

  const handleRollInitiative = async () => {
    if (mapData.tokens.length === 0) {
      toast({ title: 'No tokens on map', variant: 'destructive' });
      return;
    }
    setIsRolling(true);
    try {
      const charIds = mapData.tokens.filter(t => t.character_id).map(t => t.character_id!);
      const monsterIds = mapData.tokens.filter(t => t.monster_id).map(t => t.monster_id!);

      const [charsRes, monstersRes] = await Promise.all([
        charIds.length > 0 ? getAllCharacters(token, { limit: 1000 }) : Promise.resolve({ characters: [] }),
        monsterIds.length > 0 ? getMonstersByUser(userId, token) : Promise.resolve([]),
      ]);

      const charDataMap: Record<string, { mod: number, score: number }> = {};
      for (const c of charsRes.characters) {
        const dexMod = Math.floor((c.dexterity - 10) / 2);
        let bonus = 0;
        // Jack of All Trades (Bard Level 2+)
        if (c.class?.name === 'The Bard' && c.level >= 2) {
            const prof = Math.floor((c.level - 1) / 4) + 2;
            bonus += Math.floor(prof / 2);
        }
        charDataMap[c.id] = { mod: dexMod + bonus, score: c.dexterity };
      }

      const monsterDataMap: Record<string, { mod: number, score: number }> = {};
      for (const m of monstersRes) {
        monsterDataMap[m.id] = { mod: Math.floor((m.dexterity - 10) / 2), score: m.dexterity };
      }

      // Roll initiative for all tokens synchronously — showing one 3D die per
      // token in sequence would be unusable with many tokens on the map.
      const results: { mapToken: MapToken; detail: RollDetail }[] = mapData.tokens.map(mapToken => {
        const data = mapToken.character_id != null
          ? charDataMap[mapToken.character_id]
          : mapToken.monster_id != null
          ? monsterDataMap[mapToken.monster_id]
          : null;
        
        const modifier = data?.mod ?? 0;
        const dexterity = data?.score ?? 10;

        const { rolls } = DiceManager.rollSync('1d20');
        const natural = rolls[0];
        return { mapToken, detail: { natural, modifier, total: natural + modifier, dexterity } };
      });

      // Tie-breaker: total first, then dexterity score
      results.sort((a, b) => {
        if (b.detail.total !== a.detail.total) {
          return b.detail.total - a.detail.total;
        }
        return b.detail.dexterity - a.detail.dexterity;
      });

      const newRollDetails: Record<string, RollDetail> = {};
      for (const r of results) {
        newRollDetails[r.mapToken.id] = r.detail;
      }
      setRollDetails(newRollDetails);

      const entries = results.map((r, i) => ({ token_id: r.mapToken.id, order: i + 1 }));
      await setInitiative(mapData.id, entries, token);
      refreshMap();
      toast({ title: 'Initiative rolled!' });
    } catch (err) {
      console.error(err);
      toast({ title: 'Failed to roll initiative', variant: 'destructive' });
    } finally {
      setIsRolling(false);
    }
  };

  const handleClearInitiative = async () => {
    try {
      await clearInitiative(mapData.id, token);
      setRollDetails({});
      refreshMap();
      toast({ title: 'Initiative cleared' });
    } catch (err) {
      toast({ title: 'Failed to clear initiative', variant: 'destructive' });
    }
  };

  const handleDeleteMap = async () => {
    if (confirm('Are you sure you want to delete this map? This cannot be undone.')) {
        try {
            await deleteMap(mapData.id, token);
            toast({ title: 'Map deleted' });
            navigate('/battle-map');
        } catch (err) {
            toast({ title: 'Failed to delete map', variant: 'destructive' });
        }
    }
  }

  const handleElementTypeChange = (value: MapElementType) => {
    setSelectedElementType(value);
    // Reset properties when type changes
    setElementColor('#6B8E2380');
    setElementTrapName('');
    setElementTrapDescription('');
    setElementTrapVisible(false);
    setElementElevationLevel(0);
  };

  const buildElementProperties = (): Partial<CreateMapElementRequest> => {
    switch (selectedElementType) {
      case MapElementType.DifficultTerrain:
        return { difficult_terrain_color: elementColor };
      case MapElementType.Trap:
        return { trap_name: elementTrapName, trap_description: elementTrapDescription, trap_visible_to_players: elementTrapVisible };
      case MapElementType.Elevation:
        return { elevation_level: elementElevationLevel };
      case MapElementType.Wall:
        return {};
      default:
        return {};
    }
  };

  const handleStartPlacement = () => {
    if (selectedElementType) {
      onStartPlacement(selectedElementType, buildElementProperties());
    }
  };

  return (
    <div className="w-64 bg-sidebar border-l border-sidebar-border h-full text-sidebar-foreground flex flex-col">
      <div className="p-4 space-y-6 flex-grow overflow-y-auto">
      <div>
        <h3 className="text-lg font-bold text-sidebar-primary mb-2">DM Controls</h3>
        <p className="text-xs text-sidebar-foreground mb-4">Room Code: <span className="font-mono text-sidebar-foreground select-all">{mapData.room_code}</span></p>
      </div>

      <div className="space-y-3 border-t border-sidebar-border pt-4">
        <h4 className="font-semibold text-sm">Add Token</h4>
        <div className="space-y-2">
            <Label>Character / Monster</Label>
            <EntityCombobox
              value={selectedEntity}
              onSelect={handleEntitySelect}
              userId={userId}
              authToken={token}
            />
        </div>
        {selectedEntity?.entity_type !== 'party' && (
          <div className="space-y-2">
            <Label>Name</Label>
            <Input
                value={newTokenName}
                onChange={e => setNewTokenName(e.target.value)}
                placeholder="Token Name"
                className="bg-input border-input"
            />
          </div>
        )}
        <div className="space-y-2">
            <Label>Color</Label>
             <Select value={newTokenColor} onValueChange={setNewTokenColor}>
              <SelectTrigger className="w-full bg-input border-input">
                <SelectValue placeholder="Select Color" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="#ff0000">Red</SelectItem>
                <SelectItem value="#00ff00">Green</SelectItem>
                <SelectItem value="#0000ff">Blue</SelectItem>
                <SelectItem value="#ffff00">Yellow</SelectItem>
                <SelectItem value="#800080">Purple</SelectItem>
              </SelectContent>
            </Select>
        </div>
        <Button
          onClick={handleAddToken}
          disabled={!selectedEntity}
          className="w-full bg-primary hover:bg-primary/90"
        >
          {selectedEntity?.entity_type === 'party'
            ? `Spawn Party (${selectedEntity.party_members?.length ?? 0})`
            : 'Spawn Token'}
        </Button>
      </div>

      {/* Map Elements Section */}
      <div className="space-y-3 border-t border-sidebar-border pt-4">
        <h4 className="font-semibold text-sm">Map Elements</h4>
        <div className="space-y-2">
          <Label>Element Type</Label>
          <Select value={selectedElementType || ''} onValueChange={handleElementTypeChange}>
            <SelectTrigger className="w-full bg-input border-input">
              <SelectValue placeholder="Select Element Type" />
            </SelectTrigger>
            <SelectContent>
              {Object.values(MapElementType).map(type => (
                <SelectItem key={type} value={type}>
                  {type.charAt(0).toUpperCase() + type.slice(1).replace(/([A-Z])/g, ' $1')}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {selectedElementType === MapElementType.DifficultTerrain && (
          <div className="space-y-2">
            <Label>Color (Hex or RGBA)</Label>
            <Input
              type="text"
              value={elementColor}
              onChange={e => setElementColor(e.target.value)}
              placeholder="#RRGGBB or rgba(R,G,B,A)"
              className="bg-input border-input"
            />
          </div>
        )}

        {selectedElementType === MapElementType.Trap && (
          <>
            <div className="space-y-2">
              <Label>Trap Name</Label>
              <Input
                type="text"
                value={elementTrapName}
                onChange={e => setElementTrapName(e.target.value)}
                placeholder="e.g. Pit Trap"
                className="bg-input border-input"
              />
            </div>
            <div className="space-y-2">
              <Label>Description</Label>
              <Input
                type="text"
                value={elementTrapDescription}
                onChange={e => setElementTrapDescription(e.target.value)}
                placeholder="e.g. 10ft deep, 1d6 fall damage"
                className="bg-input border-input"
              />
            </div>
            <div className="flex items-center space-x-2">
              <Checkbox
                id="trap-visible"
                checked={elementTrapVisible}
                onCheckedChange={checked => setElementTrapVisible(checked === true)}
              />
              <Label htmlFor="trap-visible">Visible to Players</Label>
            </div>
          </>
        )}

        {selectedElementType === MapElementType.Elevation && (
          <div className="space-y-2">
            <Label>Elevation Level</Label>
            <Input
              type="number"
              value={elementElevationLevel}
              onChange={e => setElementElevationLevel(Number(e.target.value))}
              className="bg-input border-input"
            />
          </div>
        )}

        {placementModeActive ? (
            <Button
                onClick={onStopPlacement}
                className="w-full bg-red-500 hover:bg-red-600"
            >
                Stop Placing {selectedElementType ? selectedElementType.charAt(0).toUpperCase() + selectedElementType.slice(1).replace(/([A-Z])/g, ' $1') : ''}
            </Button>
        ) : (
            <Button
                onClick={handleStartPlacement}
                disabled={!selectedElementType}
                className="w-full bg-blue-500 hover:bg-blue-600"
            >
                Start Placing {selectedElementType ? selectedElementType.charAt(0).toUpperCase() + selectedElementType.slice(1).replace(/([A-Z])/g, ' $1') : 'Elements'}
            </Button>
        )}
      </div>

      {/* Initiative Section */}
      <div className="space-y-3 border-t border-sidebar-border pt-4">
        <h4 className="font-semibold text-sm">Initiative</h4>
        <div className="flex gap-2">
          <Button
            onClick={handleRollInitiative}
            disabled={isRolling || mapData.tokens.length === 0}
            className="flex-1 bg-amber-600 hover:bg-amber-700 text-white"
          >
            {isRolling ? 'Rolling...' : 'Roll Initiative'}
          </Button>
          {initiativeOrder.length > 0 && (
            <Button
              variant="outline"
              onClick={handleClearInitiative}
              className="text-destructive border-destructive hover:bg-destructive/10"
            >
              Clear
            </Button>
          )}
        </div>
        {initiativeOrder.length > 0 && (
          <div className="space-y-1">
            {initiativeOrder.map((t, i) => {
              const detail = rollDetails[t.id];
              return (
                <div key={t.id} className="flex items-center gap-2 text-sm py-1 px-2 rounded bg-sidebar-accent/50">
                  <span className="w-5 h-5 rounded-full flex items-center justify-center text-xs font-bold bg-primary text-primary-foreground shrink-0">
                    {i + 1}
                  </span>
                  <div className="w-3 h-3 rounded-full shrink-0 border border-border" style={{ backgroundColor: t.color }} />
                  <span className="flex-1 font-medium truncate text-xs">{t.name}</span>
                  {detail && (
                    <span className="text-muted-foreground text-xs font-mono shrink-0">
                      {detail.natural}{detail.modifier !== 0 ? (detail.modifier > 0 ? `+${detail.modifier}` : detail.modifier) : ''}={detail.total}
                    </span>
                  )}
                </div>
              );
            })}
          </div>
        )}
        {initiativeOrder.length === 0 && mapData.tokens.length === 0 && (
          <p className="text-xs text-muted-foreground italic">Add tokens to the map first</p>
        )}
      </div>

      <div className="space-y-3 border-t border-sidebar-border pt-4">
        <h4 className="font-semibold text-sm text-destructive">Danger Zone</h4>
        <Button variant="destructive" onClick={handleDeleteMap} className="w-full">
            Delete Map
        </Button>
      </div>
      </div>
    </div>
  );
};
