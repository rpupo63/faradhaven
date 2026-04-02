import type { ApiSpell, ApiComponent, ApiCharacterComponent, ApiCharacterSheet } from '@/types/game';

/**
 * Spell tier for grouping, filters, and cast cost. Prefers live `components` length, then API `level`
 * (matches backend synthesis: level = component count), then optional `slot_level`.
 */
export function getSpellComponentCount(
  spell: Pick<ApiSpell, 'components' | 'level' | 'slot_level'>,
): number {
  const fromComponents = spell.components?.length;
  if (typeof fromComponents === 'number' && fromComponents > 0) return fromComponents;
  if (typeof spell.level === 'number' && Number.isFinite(spell.level) && spell.level >= 1) {
    return spell.level;
  }
  if (typeof spell.slot_level === 'number' && Number.isFinite(spell.slot_level)) {
    return spell.slot_level;
  }
  return 0;
}

/**
 * Union of class + race spell pools (matches backend MergeSpellPoolComponents). Dedupes by component id.
 */
export function mergeSpellPoolComponents(
  classComponents: ApiComponent[] | undefined,
  raceComponents: ApiComponent[] | undefined,
): ApiComponent[] {
  const byId = new Map<string, ApiComponent>();
  for (const c of classComponents ?? []) {
    byId.set(c.id, c);
  }
  for (const c of raceComponents ?? []) {
    if (!byId.has(c.id)) byId.set(c.id, c);
  }
  return Array.from(byId.values()).sort((a, b) => {
    const ca = a.category ?? '';
    const cb = b.category ?? '';
    if (ca !== cb) return String(ca).localeCompare(String(cb));
    return a.name.localeCompare(b.name);
  });
}

/**
 * Prefer `available_components` from the character sheet (server-merged class + race).
 * If missing or empty, merge embedded `class.components` and `character.race.components`.
 */
export function resolveSpellPoolComponents(sheet: ApiCharacterSheet | undefined): ApiComponent[] | undefined {
  if (!sheet) return undefined;
  const direct = sheet.available_components;
  if (direct && direct.length > 0) return direct;
  const merged = mergeSpellPoolComponents(sheet.class?.components, sheet.character?.race?.components);
  return merged.length > 0 ? merged : undefined;
}

/**
 * Per-component availability: unlimited pool components (class + race) are always available;
 * non-pool components must exist in inventory with sufficient count.
 * If `spellPoolComponents` is undefined (pool not loaded), returns true so the UI does not false-negative; the server still validates on cast.
 */
export function hasAllComponents(
  spell: ApiSpell,
  spellPoolComponents: ApiComponent[] | undefined,
  inventory: ApiCharacterComponent[] | undefined,
): boolean {
  if (!spell.components?.length) return true;
  if (spellPoolComponents === undefined) return true;
  const poolIds = new Set(spellPoolComponents.map((c) => c.id));
  const needById = new Map<string, number>();
  for (const comp of spell.components) {
    if (poolIds.has(comp.id)) continue;
    needById.set(comp.id, (needById.get(comp.id) ?? 0) + 1);
  }
  for (const [compId, need] of needById) {
    const inv = inventory?.find((cc) => cc.component_id === compId);
    if ((inv?.count ?? 0) < need) return false;
  }
  return true;
}
