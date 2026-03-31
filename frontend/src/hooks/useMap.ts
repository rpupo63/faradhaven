import { useState, useEffect, useCallback, useRef } from 'react';
import { GameMap, UpdateTokenRequest, MapToken } from '@/types/map';
import { getMap, getMapByRoom, updateToken } from '@/lib/api/map';
import { useToast } from '@/hooks/use-toast';

export function useMap(mapIdOrCode: string, isRoomCode: boolean = false, token?: string) {
  const [mapData, setMapData] = useState<GameMap | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { toast } = useToast();
  const pollingRef = useRef<NodeJS.Timeout | null>(null);

  const fetchMap = useCallback(async () => {
    if (!mapIdOrCode) return;
    try {
      let data: GameMap;
      if (isRoomCode) {
        data = await getMapByRoom(mapIdOrCode, token);
      } else {
        data = await getMap(mapIdOrCode, token);
      }
      setMapData(prev => {
        if (JSON.stringify(prev) === JSON.stringify(data)) {
           return prev;
        }
        return data;
      });
      setError(null);
    } catch (err) {
      console.error(err);
      setError('Failed to load map');
    } finally {
      setLoading(false);
    }
  }, [mapIdOrCode, isRoomCode, token]);

  useEffect(() => {
    fetchMap();
  }, [fetchMap, mapIdOrCode, isRoomCode, token]);

  useEffect(() => {
    if (!mapData?.id || !token) return;

    const apiUrl = import.meta.env.VITE_API_URL || window.location.origin;
    const wsUrl = apiUrl.replace(/^http/, 'ws') + `/api/map/${mapData.id}/ws?token=${token}`;

    const ws = new WebSocket(wsUrl);

    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        if (message.type) {
          fetchMap();
        }
      } catch (e) {
        console.error('Error parsing map websocket message', e);
      }
    };

    return () => {
      ws.close();
    };
  }, [mapData?.id, token, fetchMap]);

  const moveTokenOptimistic = async (tokenId: string, x: number, y: number) => {
    if (!mapData || !token) return;

    setMapData(prev => {
      if (!prev) return null;
      return {
        ...prev,
        tokens: prev.tokens.map(t =>
          t.id === tokenId ? { ...t, grid_x: x, grid_y: y } : t
        )
      };
    });

    try {
      await updateToken(mapData.id, tokenId, { grid_x: x, grid_y: y }, token);
    } catch (err) {
      toast({
        title: "Move failed",
        description: "Could not sync move with server",
        variant: "destructive"
      });
    }
  };

  const updateTokenProperty = async (tokenId: string, updates: Partial<UpdateTokenRequest>) => {
    if (!mapData || !token) return;

    setMapData(prev => {
      if (!prev) return null;
      return {
        ...prev,
        tokens: prev.tokens.map(t =>
          t.id === tokenId ? { ...t, ...updates } as MapToken : t // Cast to MapToken
        )
      };
    });

    try {
      await updateToken(mapData.id, tokenId, updates, token);
    } catch (err) {
      toast({
        title: "Update failed",
        description: "Could not sync token update with server",
        variant: "destructive"
      });
      // Revert will happen on next poll if server update truly failed
    }
  };


  return { mapData, loading, error, moveTokenOptimistic, updateTokenProperty, refresh: fetchMap };
}
