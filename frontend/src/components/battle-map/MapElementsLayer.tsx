import React from 'react';
import { Layer, Rect, Text } from 'react-konva';
import { GameMap, MapElement, MapElementType } from '@/types/map';

interface MapElementsLayerProps {
  elements: MapElement[];
  tileSize: number;
  width: number;
  height: number;
  isDM: boolean;
}

export const MapElementsLayer: React.FC<MapElementsLayerProps> = ({ elements, tileSize, width, height, isDM }) => {
  if (!elements || !Array.isArray(elements)) {
    return null; // Don't render anything if elements is not an array
  }

  return (
    <Layer>
      {elements.map((element) => {
        const x = element.grid_x * tileSize;
        const y = element.grid_y * tileSize;

        switch (element.type) {
          case MapElementType.Wall:
            return (
              <Rect
                key={element.id}
                x={x}
                y={y}
                width={tileSize}
                height={tileSize}
                fill="#333333" // Dark grey for walls
              />
            );
          case MapElementType.DifficultTerrain:
            return (
              <Rect
                key={element.id}
                x={x}
                y={y}
                width={tileSize}
                height={tileSize}
                fill={element.difficult_terrain_properties?.color || "#6B8E2380"} // Olive green with transparency
              />
            );
          case MapElementType.Trap:
            // Only visible to DM or if explicitly visible to players
            if (isDM || element.trap_properties?.visible_to_players) {
              return (
                <Rect
                  key={element.id}
                  x={x}
                  y={y}
                  width={tileSize}
                  height={tileSize}
                  fill="rgba(255, 0, 0, 0.3)" // Red translucent overlay
                />
              );
            }
            return null;
          case MapElementType.Elevation: {
            // Represent with a subtle border and a text label for now
            const level = element.elevation_properties?.level || 0;
            const borderColor = level > 0 ? "rgba(0, 100, 0, 0.5)" : "rgba(100, 0, 0, 0.5)"; // Green for up, Red for down
            const elevationText = level > 0 ? `+${level}` : `${level}`;
            return (
              <React.Fragment key={element.id}>
                <Rect
                  x={x}
                  y={y}
                  width={tileSize}
                  height={tileSize}
                  stroke={borderColor}
                  strokeWidth={2}
                />
                <Text
                  x={x + tileSize / 4}
                  y={y + tileSize / 4}
                  text={elevationText}
                  fontSize={tileSize / 2}
                  fill={level > 0 ? "green" : "red"}
                  fontFamily="Arial"
                  align="center"
                  verticalAlign="middle"
                  width={tileSize / 2}
                  height={tileSize / 2}
                />
              </React.Fragment>
            );
          }
          default:
            return null;
        }
      })}
    </Layer>
  );
};
