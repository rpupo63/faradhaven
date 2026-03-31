export interface ApiClassResourceDefinition {
  key: string;
  display_name: string;
  category: 'pool' | 'die_size' | 'limit' | 'slot_count' | 'modifier' | 'state';
  description?: string;
  is_trackable: boolean;
  restore_on_short_rest: boolean;
  restore_on_long_rest: boolean;
  display_order: number;
}
