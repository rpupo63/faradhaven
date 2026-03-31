import React from 'react';
import { Layer } from 'react-konva';
import { Token } from './Token';
import { MapToken } from '@/types/map';
import Konva from 'konva';

interface TokenLayerProps {
  tokens: MapToken[];
  tileSize: number;
  currentUserId?: string;
  isDM: boolean;
  onMoveToken: (id: string, x: number, y: number) => void;
  onTokenContextMenu: (token: MapToken, event: Konva.KonvaEventObject<MouseEvent>) => void; // Added onTokenContextMenu prop
}

export const TokenLayer: React.FC<TokenLayerProps> = ({ tokens, tileSize, currentUserId, isDM, onMoveToken, onTokenContextMenu }) => {
  return (
    <Layer>
      {tokens.map((token) => {
        if (!token.visible && !isDM) return null;

        // Determine if draggable
        // DM can move anything. User can move their own assigned tokens.
        const isOwner = token.assigned_user_id === currentUserId;
        const canMove = isDM || isOwner;

        return (
          <Token
            key={token.id}
            token={token}
            tileSize={tileSize}
            isDraggable={canMove}
            onDragEnd={onMoveToken}
            onContextMenu={onTokenContextMenu} // Pass onTokenContextMenu down
          />
        );
      })}
    </Layer>
  );
};
