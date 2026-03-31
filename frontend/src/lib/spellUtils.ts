import type { ApiSpell, ApiComponent, ApiCharacterComponent } from '@/types/game';

/**
 * Per-component availability check: pool components (from class) are always available;
 * non-pool components must exist in the character's inventory with count >= 1.
 */
export function hasAllComponents(
  spell: ApiSpell,
  classComponents: ApiComponent[] | undefined,
  inventory: ApiCharacterComponent[] | undefined,
): boolean {
  if (!spell.components?.length) return true;
  const poolIds = new Set(classComponents?.map(c => c.id) ?? []);
  const needById = new Map<string, number>();
  for (const comp of spell.components) {
    if (poolIds.has(comp.id)) continue;
    needById.set(comp.id, (needById.get(comp.id) ?? 0) + 1);
  }
  for (const [compId, need] of needById) {
    const inv = inventory?.find(cc => cc.component_id === compId);
    if ((inv?.count ?? 0) < need) return false;
  }
  return true;
}
