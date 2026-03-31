/** Battle map token role (matches backend models.MapTokenType). */
export type MapTokenType = 'pc' | 'npc';

export interface MapToken {
  id: string;
  map_id: string;
  character_id?: string;
  monster_id?: string;
  assigned_user_id?: string;
  name: string;
  image_url: string;
  token_type: MapTokenType;
  grid_x: number;
  grid_y: number;
  size: number;
  color: string;
  visible: boolean;
  initiative_order?: number;
}

export enum MapElementType {
  Wall = "wall",
  DifficultTerrain = "difficult_terrain",
  Trap = "trap",
  Elevation = "elevation",
}

export interface TrapProperties {
  id: string;
  map_element_id: string;
  name: string;
  description: string;
  visible_to_players: boolean;
}

export interface DifficultTerrainProperties {
  id: string;
  map_element_id: string;
  color: string;
}

export interface ElevationProperties {
  id: string;
  map_element_id: string;
  level: number;
}

export interface WallProperties {
  id: string;
  map_element_id: string;
  details: string;
}

export interface MapElement {
  id: string;
  map_id: string;
  type: MapElementType;
  grid_x: number;
  grid_y: number;
  trap_properties?: TrapProperties;
  difficult_terrain_properties?: DifficultTerrainProperties;
  elevation_properties?: ElevationProperties;
  wall_properties?: WallProperties;
}

export interface GameMap {
  id: string;
  owner_id: string;
  room_code: string;
  name: string;
  background_url: string;
  grid_rows: number;
  grid_cols: number;
  tile_size: number;
  tokens: MapToken[];
  elements: MapElement[]; // New field
}

export interface CreateMapRequest {
  name: string;
  room_code: string;
  background_url: string;
  grid_rows: number;
  grid_cols: number;
  tile_size: number;
}

export interface UpdateMapRequest {
  name?: string;
  background_url?: string;
  grid_rows?: number;
  grid_cols?: number;
  tile_size?: number;
}

export interface CreateTokenRequest {
  character_id?: string;
  monster_id?: string;
  assigned_user_id?: string;
  name: string;
  image_url: string;
  token_type: MapTokenType;
  grid_x: number;
  grid_y: number;
  size: number;
  color: string;
  visible: boolean;
}

export interface CreateMapElementRequest {
  type: MapElementType;
  grid_x: number;
  grid_y: number;
  trap_name?: string;
  trap_description?: string;
  trap_visible_to_players?: boolean;
  difficult_terrain_color?: string;
  elevation_level?: number;
  wall_details?: string;
}

export interface UpdateMapElementRequest {
  type?: MapElementType;
  grid_x?: number;
  grid_y?: number;
  trap_name?: string;
  trap_description?: string;
  trap_visible_to_players?: boolean;
  difficult_terrain_color?: string;
  elevation_level?: number;
  wall_details?: string;
}

export interface UpdateTokenRequest {
  grid_x?: number;
  grid_y?: number;
  visible?: boolean;
  size?: number;
  color?: string;
  assigned_user_id?: string;
  character_id?: string | null;
  monster_id?: string | null;
  name?: string;
  image_url?: string;
  token_type?: MapTokenType;
}
