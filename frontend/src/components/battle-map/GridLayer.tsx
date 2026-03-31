import React from 'react';
import { Layer, Line } from 'react-konva';

interface GridLayerProps {
  cols: number;
  rows: number;
  tileSize: number;
  width: number;
  height: number;
}

export const GridLayer: React.FC<GridLayerProps> = ({ cols, rows, tileSize, width, height }) => {
  const lines = [];

  // Vertical lines
  for (let i = 0; i <= cols; i++) {
    lines.push(
      <Line
        key={`v-${i}`}
        points={[i * tileSize, 0, i * tileSize, height]}
        stroke="hsl(43, 89%, 70%)" // Brighter gold for visibility on dark background
        strokeWidth={1}
        listening={false}
      />
    );
  }

  // Horizontal lines
  for (let i = 0; i <= rows; i++) {
    lines.push(
      <Line
        key={`h-${i}`}
        points={[0, i * tileSize, width, i * tileSize]}
        stroke="hsl(43, 89%, 70%)" // Brighter gold for visibility on dark background
        strokeWidth={1}
        listening={false}
      />
    );
  }

  return <Layer>{lines}</Layer>;
};
