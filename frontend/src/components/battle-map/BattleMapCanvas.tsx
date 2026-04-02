import React, { useRef, useState, useEffect, useCallback } from 'react';
import { Stage, Layer, Image as KonvaImage, Rect } from 'react-konva';
import useImage from 'use-image';
import Konva from 'konva'; // Import Konva for event types
import { ZoomIn, ZoomOut, Maximize } from 'lucide-react';
import { GridLayer } from './GridLayer';
import { TokenLayer } from './TokenLayer';
import { TokenContextMenu } from './TokenContextMenu'; // Import TokenContextMenu
import { MapElementsLayer } from './MapElementsLayer'; // Import MapElementsLayer
import { GameMap, MapToken } from '@/types/map';

interface BattleMapCanvasProps {
  mapData: GameMap;
  currentUserId?: string;
  userId: string;
  authToken: string;
  onMoveToken: (id: string, x: number, y: number) => void;
  onUpdateToken: (tokenId: string, updates: Partial<MapToken>) => void;
  onDeleteToken: (tokenId: string) => void;
  isDM: boolean;
  onCanvasClickForPlacement: (gridCoords: { gridX: number, gridY: number }[]) => void; // Updated prop signature
  placementModeActive: boolean; // New prop
}

export const BattleMapCanvas: React.FC<BattleMapCanvasProps> = ({
  mapData,
  currentUserId,
  userId,
  authToken,
  onMoveToken,
  onUpdateToken,
  onDeleteToken,
  isDM,
  onCanvasClickForPlacement, // Destructure new prop
  placementModeActive, // Destructure new prop
}) => {
  const [bgImage] = useImage(mapData.background_url || '');
  const containerRef = useRef<HTMLDivElement>(null);
  const stageRef = useRef<Konva.Stage>(null);
  const isInitialized = useRef(false);
  const [containerWidth, setContainerWidth] = useState(0);
  const [containerHeight, setContainerHeight] = useState(0);

  // Zoom/pan state
  const [mapScale, setMapScale] = useState(1);
  const [mapOffset, setMapOffset] = useState({ x: 0, y: 0 });
  const [isPanning, setIsPanning] = useState(false);
  const [panStart, setPanStart] = useState<{ x: number; y: number; offsetX: number; offsetY: number } | null>(null);

  // State for drag-and-select placement
  const [isDragging, setIsDragging] = useState(false);
  const [dragStart, setDragStart] = useState<{ x: number; y: number } | null>(null);
  const [dragEnd, setDragEnd] = useState<{ x: number; y: number } | null>(null);
  const [selectedGridCells, setSelectedGridCells] = useState<Set<string>>(new Set());

  // Context menu state
  const [contextMenuState, setContextMenuState] = useState<{
    isOpen: boolean;
    x: number;
    y: number;
    token: MapToken | null;
  }>({
    isOpen: false,
    x: 0,
    y: 0,
    token: null,
  });

  useEffect(() => {
    const checkSize = () => {
      if (containerRef.current) {
        setContainerWidth(containerRef.current.offsetWidth);
        setContainerHeight(containerRef.current.offsetHeight);
      }
    };

    checkSize();
    window.addEventListener('resize', checkSize);

    const observer = new ResizeObserver(entries => {
      for (const entry of entries) {
        setContainerWidth(entry.contentRect.width);
        setContainerHeight(entry.contentRect.height);
      }
    });

    const currentContainerRef = containerRef.current; // Capture current value

    if (currentContainerRef) {
      observer.observe(currentContainerRef);
    }

    return () => {
      window.removeEventListener('resize', checkSize);
      if (currentContainerRef) { // Use captured value
        observer.unobserve(currentContainerRef);
      }
    };
  }, []);

  // Natural map dimensions — stage transform handles all visual scaling
  const tileSize = mapData.tile_size;
  const mapNaturalWidth = mapData.grid_cols * tileSize;
  const mapNaturalHeight = mapData.grid_rows * tileSize;

  // Initialize fit-scale once we have valid container dimensions
  useEffect(() => {
    if (containerWidth === 0 || containerHeight === 0) return;
    if (isInitialized.current) return;
    const fitScale = Math.min(containerWidth / mapNaturalWidth, containerHeight / mapNaturalHeight);
    setTimeout(() => {
      setMapScale(fitScale);
      setMapOffset({
        x: (containerWidth - mapNaturalWidth * fitScale) / 2,
        y: (containerHeight - mapNaturalHeight * fitScale) / 2,
      });
    }, 0);
    isInitialized.current = true;
  }, [containerWidth, containerHeight, mapNaturalWidth, mapNaturalHeight]);

  const handleTokenContextMenu = (token: MapToken, event: Konva.KonvaEventObject<MouseEvent>) => {
    // Get mouse position relative to the document
    const mouseX = event.evt.clientX;
    const mouseY = event.evt.clientY;

    setContextMenuState({
      isOpen: true,
      x: mouseX,
      y: mouseY,
      token: token,
    });
  };

  const handleCloseContextMenu = () => {
    setContextMenuState(prev => ({ ...prev, isOpen: false }));
  };

  const handleDismissContextMenu = () => {
    setContextMenuState(prev => ({ ...prev, isOpen: false, token: null }));
  };

  const handleChangeAffiliation = (tokenId: string, updates: Partial<MapToken>) => {
    onUpdateToken(tokenId, updates);
  };

  // Reverse stage transform to convert pointer position to map coords
  const getGridCoordinates = useCallback((pointerPosition: { x: number; y: number }) => {
    const mapX = (pointerPosition.x - mapOffset.x) / mapScale;
    const mapY = (pointerPosition.y - mapOffset.y) / mapScale;
    return {
      x: Math.max(0, Math.min(Math.floor(mapX / tileSize), mapData.grid_cols - 1)),
      y: Math.max(0, Math.min(Math.floor(mapY / tileSize), mapData.grid_rows - 1)),
    };
  }, [mapOffset, mapScale, tileSize, mapData.grid_cols, mapData.grid_rows]);

  // --- Zoom/pan handlers ---

  const handleWheel = useCallback((e: Konva.KonvaEventObject<WheelEvent>) => {
    e.evt.preventDefault();
    const stage = stageRef.current;
    if (!stage) return;
    const pointer = stage.getPointerPosition();
    if (!pointer) return;
    const FACTOR = 1.1;
    const direction = e.evt.deltaY > 0 ? -1 : 1;
    const newScale = Math.max(0.1, Math.min(8, mapScale * (direction > 0 ? FACTOR : 1 / FACTOR)));
    setMapScale(newScale);
    setMapOffset({
      x: pointer.x - (pointer.x - mapOffset.x) * (newScale / mapScale),
      y: pointer.y - (pointer.y - mapOffset.y) * (newScale / mapScale),
    });
  }, [mapScale, mapOffset]);

  const handleFitToScreen = useCallback(() => {
    if (!containerWidth || !containerHeight) return;
    const fitScale = Math.min(containerWidth / mapNaturalWidth, containerHeight / mapNaturalHeight);
    setMapScale(fitScale);
    setMapOffset({
      x: (containerWidth - mapNaturalWidth * fitScale) / 2,
      y: (containerHeight - mapNaturalHeight * fitScale) / 2,
    });
  }, [containerWidth, containerHeight, mapNaturalWidth, mapNaturalHeight]);

  const handleZoomIn = useCallback(() => {
    const newScale = Math.min(8, mapScale * 1.2);
    const cx = containerWidth / 2, cy = containerHeight / 2;
    setMapOffset({ x: cx - (cx - mapOffset.x) * (newScale / mapScale), y: cy - (cy - mapOffset.y) * (newScale / mapScale) });
    setMapScale(newScale);
  }, [mapScale, mapOffset, containerWidth, containerHeight]);

  const handleZoomOut = useCallback(() => {
    const newScale = Math.max(0.1, mapScale / 1.2);
    const cx = containerWidth / 2, cy = containerHeight / 2;
    setMapOffset({ x: cx - (cx - mapOffset.x) * (newScale / mapScale), y: cy - (cy - mapOffset.y) * (newScale / mapScale) });
    setMapScale(newScale);
  }, [mapScale, mapOffset, containerWidth, containerHeight]);

  // --- Mouse handlers ---

  const handleStageMouseDown = (e: Konva.KonvaEventObject<MouseEvent>) => {
    if (placementModeActive) {
      e.evt.preventDefault();
      const stage = e.target.getStage();
      if (stage) {
        const pointerPosition = stage.getPointerPosition();
        if (pointerPosition) {
          const gridCoords = getGridCoordinates(pointerPosition);
          setIsDragging(true);
          setDragStart(gridCoords);
          setDragEnd(gridCoords); // Initialize dragEnd to dragStart
          setSelectedGridCells(new Set([`${gridCoords.x},${gridCoords.y}`]));
        }
      }
    } else if (e.evt.button === 0 && e.target === e.target.getStage()) {
      setIsPanning(true);
      setPanStart({ x: e.evt.clientX, y: e.evt.clientY, offsetX: mapOffset.x, offsetY: mapOffset.y });
      e.evt.preventDefault();
    }
    handleCloseContextMenu();
  };

  const handleStageMouseMove = (e: Konva.KonvaEventObject<MouseEvent>) => {
    if (isDragging && placementModeActive) {
      e.evt.preventDefault();
      const stage = e.target.getStage();
      if (stage) {
        const pointerPosition = stage.getPointerPosition();
        if (pointerPosition && dragStart) {
          const currentGridCoords = getGridCoordinates(pointerPosition);
          setDragEnd(currentGridCoords);

          const newSelectedCells = new Set<string>();
          const startX = Math.min(dragStart.x, currentGridCoords.x);
          const endX = Math.max(dragStart.x, currentGridCoords.x);
          const startY = Math.min(dragStart.y, currentGridCoords.y);
          const endY = Math.max(dragStart.y, currentGridCoords.y);

          for (let x = startX; x <= endX; x++) {
            for (let y = startY; y <= endY; y++) {
              newSelectedCells.add(`${x},${y}`);
            }
          }
          setSelectedGridCells(newSelectedCells);
        }
      }
    } else if (isPanning && panStart) {
      setMapOffset({ x: panStart.offsetX + (e.evt.clientX - panStart.x), y: panStart.offsetY + (e.evt.clientY - panStart.y) });
    }
  };

  const handleStageMouseUp = async (e: Konva.KonvaEventObject<MouseEvent>) => {
    if (isDragging && placementModeActive) {
      e.evt.preventDefault();
      setIsDragging(false);
      setDragStart(null);
      setDragEnd(null);

      if (selectedGridCells.size > 0) {
        const coordsToPlace = Array.from(selectedGridCells).map(cell => {
          const [gridXStr, gridYStr] = cell.split(',');
          return { gridX: parseInt(gridXStr, 10), gridY: parseInt(gridYStr, 10) };
        });
        await onCanvasClickForPlacement(coordsToPlace);
      }
      setSelectedGridCells(new Set()); // Clear selection after placement
    } else if (placementModeActive) {
        // Single click placement if not dragging
        const stage = e.target.getStage();
        if (stage) {
            const pointerPosition = stage.getPointerPosition();
            if (pointerPosition) {
                const gridCoords = getGridCoordinates(pointerPosition);
                await onCanvasClickForPlacement([{ gridX: gridCoords.x, gridY: gridCoords.y }]); // Pass as array
            }
        }
    }
    if (isPanning) {
      setIsPanning(false);
      setPanStart(null);
    }
  };

  return (
    <div
        ref={containerRef}
        className="border border-border rounded bg-card shadow-xl w-full h-full relative"
        style={{ cursor: placementModeActive ? (isDragging ? 'grabbing' : 'crosshair') : (isPanning ? 'grabbing' : 'default') }}
    >
      {containerWidth > 0 && containerHeight > 0 && (
        <Stage
            ref={stageRef}
            width={containerWidth}
            height={containerHeight}
            x={mapOffset.x}
            y={mapOffset.y}
            scaleX={mapScale}
            scaleY={mapScale}
            onWheel={handleWheel}
            onContextMenu={(e) => { e.evt.preventDefault(); handleCloseContextMenu(); }}
            onMouseDown={handleStageMouseDown}
            onMouseMove={handleStageMouseMove}
            onMouseUp={handleStageMouseUp}
        >
          <Layer>
              {/* Background Image */}
              {bgImage && (
                  <KonvaImage
                      image={bgImage}
                      width={mapNaturalWidth}
                      height={mapNaturalHeight}
                  />
              )}
              {/* Fallback background color if no image */}
              {!bgImage && (
                  <Rect
                      width={mapNaturalWidth}
                      height={mapNaturalHeight}
                      fill="hsl(var(--background))"
                  />
              )}
          </Layer>

          <GridLayer
            cols={mapData.grid_cols}
            rows={mapData.grid_rows}
            tileSize={tileSize}
            width={mapNaturalWidth}
            height={mapNaturalHeight}
          />

          <MapElementsLayer
            elements={mapData.elements}
            tileSize={tileSize}
            width={mapNaturalWidth}
            height={mapNaturalHeight}
            isDM={isDM}
          />

          <TokenLayer
            tokens={mapData.tokens}
            tileSize={tileSize}
            currentUserId={currentUserId}
            isDM={isDM}
            onMoveToken={onMoveToken}
            onTokenContextMenu={isDM ? handleTokenContextMenu : () => {}} // Only allow DM to open context menu
          />
          <Layer>
            {isDragging && dragStart && dragEnd && (
              <Rect
                x={Math.min(dragStart.x, dragEnd.x) * tileSize}
                y={Math.min(dragStart.y, dragEnd.y) * tileSize}
                width={(Math.abs(dragStart.x - dragEnd.x) + 1) * tileSize}
                height={(Math.abs(dragStart.y - dragEnd.y) + 1) * tileSize}
                fill="rgba(0, 161, 255, 0.3)" // Translucent blue
                stroke="rgb(0, 161, 255)"
                strokeWidth={2}
              />
            )}
          </Layer>
        </Stage>
      )}

      {/* Zoom toolbar */}
      <div className="absolute bottom-4 right-4 flex items-center gap-1 bg-card/90 border border-border rounded-md shadow-lg p-1 z-10">
        <button
          onClick={handleZoomOut}
          className="h-8 w-8 flex items-center justify-center rounded hover:bg-accent text-foreground"
          title="Zoom Out"
        >
          <ZoomOut className="h-4 w-4" />
        </button>
        <span className="text-xs font-mono text-muted-foreground w-12 text-center select-none">
          {Math.round(mapScale * 100)}%
        </span>
        <button
          onClick={handleZoomIn}
          className="h-8 w-8 flex items-center justify-center rounded hover:bg-accent text-foreground"
          title="Zoom In"
        >
          <ZoomIn className="h-4 w-4" />
        </button>
        <div className="w-px h-5 bg-border mx-1" />
        <button
          onClick={handleFitToScreen}
          className="h-8 w-8 flex items-center justify-center rounded hover:bg-accent text-foreground"
          title="Fit to Screen"
        >
          <Maximize className="h-4 w-4" />
        </button>
      </div>

      {contextMenuState.token && isDM && (
        <TokenContextMenu
          token={contextMenuState.token}
          x={contextMenuState.x}
          y={contextMenuState.y}
          isOpen={contextMenuState.isOpen}
          onClose={handleCloseContextMenu}
          onDismiss={handleDismissContextMenu}
          onChangeAffiliation={handleChangeAffiliation}
          onDeleteToken={onDeleteToken}
          userId={userId}
          authToken={authToken}
        />
      )}

    </div>
  );
};
