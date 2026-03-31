import React, { useRef } from 'react';
import { Group, Circle, Image as KonvaImage } from 'react-konva';
import useImage from 'use-image';
import Konva from 'konva';
import { MapToken } from '@/types/map';

interface TokenProps {
  token: MapToken;
  tileSize: number;
  isDraggable: boolean;
  onDragEnd: (id: string, x: number, y: number) => void;
  onContextMenu: (token: MapToken, event: Konva.KonvaEventObject<MouseEvent>) => void; // Added onContextMenu prop
}

export const Token: React.FC<TokenProps> = ({ token, tileSize, isDraggable, onDragEnd, onContextMenu }) => {
  const [image] = useImage(token.image_url || '/placeholder.svg'); // Fallback image
  const groupRef = useRef<Konva.Group>(null);

  // Calculate pixel position
  const x = token.grid_x * tileSize;
  const y = token.grid_y * tileSize;
  const size = token.size * tileSize;
  const radius = size / 2;

  const handleDragEnd = (e: Konva.KonvaEventObject<DragEvent>) => {
    const node = groupRef.current;
    if (!node) return;

    // Snap to grid logic
    const rawX = e.target.x();
    const rawY = e.target.y();

    const snappedX = Math.round(rawX / tileSize) * tileSize;
    const snappedY = Math.round(rawY / tileSize) * tileSize;

    // Calculate grid coordinates
    const gridX = Math.round(snappedX / tileSize);
    const gridY = Math.round(snappedY / tileSize);

    node.position({ x: snappedX, y: snappedY });

    onDragEnd(token.id, gridX, gridY);
  };

  const handleContextMenu = (e: Konva.KonvaEventObject<MouseEvent>) => {
    e.evt.preventDefault(); // Prevent default browser context menu
    e.cancelBubble = true;  // Stop Konva event from bubbling to Stage
    onContextMenu(token, e);
  };

  return (
    <Group
      ref={groupRef}
      x={x}
      y={y}
      draggable={isDraggable}
      onDragEnd={handleDragEnd}
      onContextmenu={handleContextMenu} // Attach context menu handler
      width={size}
      height={size}
    >
      {/* Token Border/Background */}
      <Circle
        x={radius}
        y={radius}
        radius={radius - 2}
        fill="hsl(var(--card))"
        stroke={token.color || 'hsl(var(--foreground))'}
        strokeWidth={3}
      />

      {/* Token Image */}
      {image && (
        <KonvaImage
          image={image}
          x={2} // Padding
          y={2}
          width={size - 4}
          height={size - 4}
          cornerRadius={radius} // Make it circular (approx) - actually cornerRadius for Rect, we might need to clip
        />
      )}
    </Group>
  );
};
