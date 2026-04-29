import { describe, expect, it } from 'vitest';
import { deriveLootLevelBand, getThemesForSource, recommendedTierForReward } from './lootThemes';

describe('lootThemes', () => {
  it('derives level bands from character level', () => {
    expect(deriveLootLevelBand(1)).toBe('novice');
    expect(deriveLootLevelBand(7)).toBe('adventurer');
    expect(deriveLootLevelBand(12)).toBe('veteran');
    expect(deriveLootLevelBand(18)).toBe('legend');
  });

  it('filters themes by source compatibility', () => {
    const commonEnemyThemes = getThemesForSource('common_enemy').map((theme) => theme.value);
    const bossThemes = getThemesForSource('boss_enemy').map((theme) => theme.value);

    expect(commonEnemyThemes).not.toContain('rich');
    expect(bossThemes).toContain('rich');
  });

  it('maps reward amount to recommended tier', () => {
    expect(recommendedTierForReward('scarce')).toBe('low');
    expect(recommendedTierForReward('standard')).toBe('medium');
    expect(recommendedTierForReward('bountiful')).toBe('medium');
    expect(recommendedTierForReward('jackpot')).toBe('high');
  });
});
