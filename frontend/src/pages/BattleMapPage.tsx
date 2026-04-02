import React, { useState, useEffect, useCallback } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { useMap } from '@/hooks/useMap';
import { BattleMapCanvas } from '@/components/battle-map/BattleMapCanvas';
import { DMSidebar } from '@/components/battle-map/DMSidebar';
import { RoomCodeEntry } from '@/components/battle-map/RoomCodeEntry';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardHeader, CardTitle, CardContent, CardDescription, CardFooter } from '@/components/ui/card';
import { getUserMaps, createMap, createMapElement, deleteMap, deleteToken } from '@/lib/api/map'; // Import createMapElement and deleteToken
import { GameMap, CreateMapRequest, MapElementType, CreateMapElementRequest } from '@/types/map'; // Import MapElementType, CreateMapElementRequest
import { useToast } from '@/hooks/use-toast';
import { LoadingQuill } from '@/components/LoadingQuill';
import { PlusCircle, Map as MapIcon, LogIn, RefreshCw, Trash2 } from 'lucide-react';

const DND_WORDS = [
  'dragon', 'goblin', 'dungeon', 'wizard', 'paladin', 'rogue', 'cleric', 'bard',
  'ranger', 'fighter', 'sorcerer', 'druid', 'warlock', 'monk', 'barbarian', 'tavern',
  'shield', 'sword', 'axe', 'bow', 'staff', 'wand', 'scroll', 'potion', 'flask',
  'mithral', 'adamantine', 'basilisk', 'beholder', 'lich', 'vampire', 'werewolf',
  'skeleton', 'ghoul', 'wraith', 'specter', 'banshee', 'dryad', 'pixie', 'gnome',
  'dwarf', 'elf', 'halfling', 'tiefling', 'aasimar', 'orc', 'kobold', 'hobgoblin',
  'troll', 'giant', 'cyclops', 'minotaur', 'harpy', 'medusa', 'sphinx', 'chimera',
  'hydra', 'wyvern', 'griffon', 'pegasus', 'unicorn', 'phoenix', 'crypt', 'forge',
  'citadel', 'tower', 'ruins', 'temple', 'shrine', 'altar', 'rune', 'glyph',
  'arcane', 'divine', 'infernal', 'eldritch', 'shadow', 'flame', 'frost', 'thunder',
  'radiant', 'necrotic', 'poison', 'chaos', 'vortex', 'abyss', 'portal', 'nexus',
  'curse', 'ritual', 'sigil', 'tome', 'grimoire', 'caster', 'invoke', 'conjure',
  'spectre', 'golem', 'gargoyle', 'manticore', 'basilisk', 'kraken', 'behemoth',
  'champion', 'herald', 'sage', 'oracle', 'sentinel', 'arbiter', 'wanderer',
];

function generateRoomCode(existingCodes?: string[]): string {
  const shuffled = [...DND_WORDS].sort(() => Math.random() - 0.5);
  for (const word of shuffled) {
    const code = word.toUpperCase();
    if (!existingCodes || !existingCodes.includes(code)) return code;
  }
  // All words taken — append a random number
  const word = shuffled[0].toUpperCase();
  return `${word}${Math.floor(Math.random() * 900) + 100}`;
}

export default function BattleMapPage({ initialMapId, initialRoomCode }: { initialMapId?: string; initialRoomCode?: string; }) {
  const { mapId: paramMapId, roomCode: paramRoomCode } = useParams();
  const mapId = initialMapId || paramMapId;
  const roomCode = initialRoomCode || paramRoomCode;
  const { token, user } = useAuth();
  
  // If we have mapId or roomCode, we are in View Mode
  if (mapId) {
      return <MapView mapIdOrCode={mapId} isRoomCode={false} token={token || undefined} userId={user?.id} />;
  }
  if (roomCode) {
      return <MapView mapIdOrCode={roomCode} isRoomCode={true} token={token || undefined} userId={user?.id} />;
  }

  // Otherwise, Landing Mode
  return <BattleMapLanding token={token} userId={user?.id} />;
}

import { ResizablePanelGroup, ResizablePanel, ResizableHandle } from '@/components/ui/resizable';

// --- Map View Component ---

const MapView = ({ mapIdOrCode, isRoomCode, token, userId }: { mapIdOrCode: string, isRoomCode: boolean, token?: string, userId?: string }) => {
    const { mapData, loading, error, moveTokenOptimistic, updateTokenProperty, refresh } = useMap(mapIdOrCode, isRoomCode, token);
    const { toast } = useToast();

    const [placementModeActive, setPlacementModeActive] = useState(false);
    const [placementElementType, setPlacementElementType] = useState<MapElementType | null>(null);
    const [placementElementProperties, setPlacementElementProperties] = useState<Partial<CreateMapElementRequest>>({});

    const handleStartPlacement = useCallback((elementType: MapElementType, properties: Partial<CreateMapElementRequest>) => {
        setPlacementModeActive(true);
        setPlacementElementType(elementType);
        setPlacementElementProperties(properties);
        toast({ title: `Placement mode active: ${elementType}. Click on the map to place.`, duration: 3000 });
    }, [toast]);

    const handleStopPlacement = useCallback(() => {
        setPlacementModeActive(false);
        setPlacementElementType(null);
        setPlacementElementProperties({});
        toast({ title: 'Placement mode deactivated.', duration: 3000 });
    }, [toast]);

    const handleCanvasClickForPlacement = useCallback(async (gridCoords: { gridX: number, gridY: number }[]) => {
        if (!placementModeActive || !placementElementType || !mapData?.id || !token) {
            return;
        }

        const placedCount = { success: 0, failed: 0 };

        await Promise.all(gridCoords.map(async ({ gridX, gridY }) => {
            try {
                const req: CreateMapElementRequest = {
                    type: placementElementType,
                    grid_x: gridX,
                    grid_y: gridY,
                    ...placementElementProperties,
                };
                await createMapElement(mapData.id, req, token);
                placedCount.success++;
            } catch (err) {
                console.error(`Failed to place map element at ${gridX},${gridY}:`, err);
                placedCount.failed++;
            }
        }));

        if (placedCount.success > 0) {
            toast({ title: `${placedCount.success} ${placementElementType}(s) placed successfully!` });
        }
        if (placedCount.failed > 0) {
            toast({ title: `${placedCount.failed} ${placementElementType}(s) failed to place.`, variant: 'destructive' });
        }
        refresh(); // Refresh map to show new elements
    }, [placementModeActive, placementElementType, placementElementProperties, mapData, token, refresh, toast]);

    const handleDeleteToken = useCallback(async (tokenId: string) => {
        if (!mapData?.id || !token) return;
        try {
            await deleteToken(mapData.id, tokenId, token);
            toast({ title: 'Token removed' });
            refresh();
        } catch (err) {
            toast({ title: 'Failed to remove token', variant: 'destructive' });
        }
    }, [mapData?.id, token, refresh, toast]);


    if (loading) {
        return (
            <div className="flex justify-center items-center h-[calc(100vh-4rem)]">
                <LoadingQuill label="Loading Map..." />
            </div>
        );
    }

    if (error || !mapData) {
        return (
            <div className="flex justify-center items-center h-[calc(100vh-4rem)] text-destructive">
                <Card className="bg-card border-destructive">
                    <CardHeader>
                        <CardTitle>Error Loading Map</CardTitle>
                        <CardDescription>{error || "Map not found"}</CardDescription>
                    </CardHeader>
                </Card>
            </div>
        );
    }

    const isDM = mapData.owner_id === userId;

    return (
        <ResizablePanelGroup direction="horizontal" className="flex h-[calc(100vh-4rem)]">
            <ResizablePanel defaultSize={isDM ? 75 : 100} minSize={30}>
                <div className="flex-1 bg-background relative h-full">
                    {/* Room Code Overlay for Players */}
                    <div className="absolute top-4 right-4 bg-card/80 p-2 rounded text-xs text-muted-foreground z-10">
                        Room: <span className="font-mono text-foreground">{mapData.room_code}</span>
                    </div>

                    <BattleMapCanvas
                        mapData={mapData}
                        currentUserId={userId}
                        userId={userId ?? ''}
                        authToken={token ?? ''}
                        onMoveToken={moveTokenOptimistic}
                        onUpdateToken={updateTokenProperty}
                        onDeleteToken={handleDeleteToken}
                        isDM={isDM}
                        onCanvasClickForPlacement={handleCanvasClickForPlacement}
                        placementModeActive={placementModeActive}
                    />
                </div>
            </ResizablePanel>

            {isDM && token && userId && (
                <>
                    <ResizableHandle withHandle />
                    <ResizablePanel defaultSize={25} minSize={20} maxSize={50}>
                        <DMSidebar
                            mapData={mapData}
                            token={token}
                            userId={userId}
                            refreshMap={refresh}
                            onStartPlacement={handleStartPlacement}
                            onStopPlacement={handleStopPlacement}
                            placementModeActive={placementModeActive}
                        />
                    </ResizablePanel>
                </>
            )}
        </ResizablePanelGroup>
    );
};

// --- Landing Component ---

const BattleMapLanding = ({ token, userId }: { token: string | null, userId?: string }) => {
    const [userMaps, setUserMaps] = useState<GameMap[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const { toast } = useToast();
    const navigate = useNavigate();

    // Create Map State
    const [newMapName, setNewMapName] = useState('');
    const [newMapBackgroundUrl, setNewMapBackgroundUrl] = useState('');
    const [newMapCode, setNewMapCode] = useState(() => generateRoomCode());
    const [newMapRows, setNewMapRows] = useState(20);
    const [newMapCols, setNewMapCols] = useState(20);

    const loadMaps = useCallback(async () => {
        if (!token || !userId) return;
        setIsLoading(true);
        try {
            const maps = await getUserMaps(userId, token);
            setUserMaps(maps);
        } catch (err) {
            console.error(err);
        } finally {
            setIsLoading(false);
        }
    }, [token, userId]);

    useEffect(() => {
        if (token && userId) {
            loadMaps();
        }
    }, [token, userId, loadMaps]);

    const handleDeleteMap = async (mapId: string, mapName: string, e: React.MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();
        if (!token) return;
        if (!confirm(`Delete "${mapName}"? This cannot be undone.`)) return;
        try {
            await deleteMap(mapId, token);
            setUserMaps(prev => prev.filter(m => m.id !== mapId));
            toast({ title: 'Map deleted' });
        } catch (err) {
            toast({ title: 'Failed to delete map', variant: 'destructive' });
        }
    };

    const handleCreateMap = async () => {
        if (!token) return;
        if (!newMapName || !newMapCode) {
            toast({ title: "Validation Error", description: "Name and Room Code are required", variant: "destructive" });
            return;
        }

        try {
            const req: CreateMapRequest = {
                name: newMapName,
                room_code: newMapCode,
                background_url: newMapBackgroundUrl,
                grid_rows: newMapRows,
                grid_cols: newMapCols,
                tile_size: 50
            };
            const map = await createMap(req, token);
            toast({ title: "Map Created!" });
            navigate(`/dm-tools?tab=battle-map&mapId=${map.id}`);
        } catch (err) {
            // Room code taken — auto-generate a new one with a number suffix
            const existingCodes = userMaps.map(m => m.room_code);
            const baseWord = newMapCode.replace(/\d+$/, '');
            const newCode = `${baseWord}${Math.floor(Math.random() * 900) + 100}`;
            setNewMapCode(newCode);
            toast({ title: "Room code taken", description: `Generated new code: ${newCode}. Try creating again.`, variant: "destructive" });
        }
    };

    return (
        <div className="w-full space-y-8">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                {/* Join Section */}
                <div className="space-y-6">
                    <h2 className="text-xl font-semibold text-foreground">Join a Session</h2>
                    <RoomCodeEntry />
                </div>

                {/* Create Section (Only if logged in) */}
                {token ? (
                    <div className="space-y-6">
                         <h2 className="text-xl font-semibold text-foreground">Create New Map</h2>
                         <Card className="bg-card border-border">
                            <CardHeader>
                                <CardTitle className="text-lg">Map Details</CardTitle>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                <div className="space-y-2">
                                    <Label>Map Name</Label>
                                    <Input value={newMapName} onChange={e => setNewMapName(e.target.value)} placeholder="e.g. Goblin Ambush" className="bg-input border-input"/>
                                </div>
                                <div className="space-y-2">
                                    <Label>Background Image URL</Label>
                                    <Input value={newMapBackgroundUrl} onChange={e => setNewMapBackgroundUrl(e.target.value)} placeholder="https://..." className="bg-input border-input"/>
                                </div>
                                <div className="space-y-2">
                                    <Label>Room Code (Unique)</Label>
                                    <div className="flex gap-2">
                                        <Input value={newMapCode} onChange={e => setNewMapCode(e.target.value.toUpperCase())} placeholder="e.g. GOBLIN" className="bg-input border-input font-mono"/>
                                        <Button
                                            type="button"
                                            variant="outline"
                                            size="icon"
                                            onClick={() => setNewMapCode(generateRoomCode(userMaps.map(m => m.room_code)))}
                                            title="Generate new room code"
                                        >
                                            <RefreshCw className="w-4 h-4" />
                                        </Button>
                                    </div>
                                </div>
                                <div className="grid grid-cols-2 gap-4">
                                    <div className="space-y-2">
                                        <Label>Rows</Label>
                                        <Input type="number" value={newMapRows} onChange={e => setNewMapRows(Number(e.target.value))} className="bg-input border-input"/>
                                    </div>
                                    <div className="space-y-2">
                                        <Label>Cols</Label>
                                        <Input type="number" value={newMapCols} onChange={e => setNewMapCols(Number(e.target.value))} className="bg-input border-input"/>
                                    </div>
                                </div>
                            </CardContent>
                            <CardFooter>
                                <Button onClick={handleCreateMap} className="w-full bg-primary hover:bg-primary/90">
                                    <PlusCircle className="w-4 h-4 mr-2" /> Create Map
                                </Button>
                            </CardFooter>
                         </Card>
                    </div>
                ) : (
                    <div className="flex items-center justify-center h-full">
                        <div className="text-center space-y-4">
                            <p className="text-muted-foreground">Log in to create and manage maps.</p>
                            <Link to="/login">
                                <Button variant="outline"><LogIn className="mr-2 h-4 w-4"/> Login</Button>
                            </Link>
                        </div>
                    </div>
                )}
            </div>

            {/* My Maps List */}
            {token && (
                <div className="mt-12">
                    <h2 className="text-xl font-semibold text-foreground mb-6">My Maps</h2>
                    {isLoading ? (
                        <LoadingQuill label="Loading your maps..." />
                    ) : userMaps.length > 0 ? (
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                            {userMaps.map(map => {
                                const isOwner = map.owner_id === userId;
                                return (
                                    <Link to={`/battle-map/${map.id}`} key={map.id}>
                                        <Card className="bg-card border-border hover:border-primary/50 transition-colors cursor-pointer group">
                                            <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-2">
                                                <div>
                                                    <CardTitle className="text-primary group-hover:text-primary/90">{map.name}</CardTitle>
                                                    <CardDescription>Code: {map.room_code}</CardDescription>
                                                </div>
                                                {isOwner ? (
                                                    <Button
                                                        variant="ghost"
                                                        size="icon"
                                                        className="text-muted-foreground hover:text-destructive shrink-0"
                                                        onClick={(e) => handleDeleteMap(map.id, map.name, e)}
                                                    >
                                                        <Trash2 className="w-4 h-4" />
                                                    </Button>
                                                ) : (
                                                    <span className="text-xs text-muted-foreground border border-border rounded px-1.5 py-0.5 shrink-0">Joined</span>
                                                )}
                                            </CardHeader>
                                            <CardContent>
                                                <div className="text-sm text-muted-foreground">
                                                    Size: {map.grid_cols}x{map.grid_rows}
                                                </div>
                                                <div className="text-sm text-muted-foreground">
                                                    Tokens: {map.tokens?.length || 0}
                                                </div>
                                            </CardContent>
                                        </Card>
                                    </Link>
                                );
                            })}
                        </div>
                    ) : (
                        <p className="text-muted-foreground italic">No maps created yet.</p>
                    )}
                </div>
            )}
        </div>
    );
};
