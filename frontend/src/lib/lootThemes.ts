import type { LootRoomTheme, LootRewardAmount, LootSource, LootTier, LootLevelBand } from '@/types/game/api';

export interface LootThemeOption {
  value: LootRoomTheme;
  label: string;
  description: string;
  compatibleSources: LootSource[];
  flavor: string;
}

export interface LootRewardOption {
  value: LootRewardAmount;
  label: string;
  description: string;
}

export const LOOT_THEME_OPTIONS: LootThemeOption[] = [
  {
    value: 'dungeon',
    label: 'Dungeon',
    description: 'Dusty ruins, trapped halls, and forgotten keeps.',
    compatibleSources: ['common_enemy', 'boss_enemy', 'room'],
    flavor: 'Expect old steel, relic scraps, and hidden stashes.',
  },
  {
    value: 'office',
    label: 'Office',
    description: 'Administrative compounds and paper-heavy compounds.',
    compatibleSources: ['common_enemy', 'room'],
    flavor: 'More utility finds, fewer martial prizes.',
  },
  {
    value: 'rich',
    label: 'Rich Quarter',
    description: 'Luxury estates, guild vaults, and noble holdings.',
    compatibleSources: ['boss_enemy', 'room'],
    flavor: 'Higher odds of premium valuables and polished gear.',
  },
  {
    value: 'poor',
    label: 'Poor Quarter',
    description: 'Shanties, market debris, and survival caches.',
    compatibleSources: ['common_enemy', 'room'],
    flavor: 'Lean hauls with occasional practical treasures.',
  },
  {
    value: 'gangster',
    label: 'Gangster Turf',
    description: 'Smuggler hideouts and black-market corners.',
    compatibleSources: ['common_enemy', 'boss_enemy', 'room'],
    flavor: 'Weapon-heavy loot and dirty coin.',
  },
  {
    value: 'arcane',
    label: 'Arcane Lab',
    description: 'Magic facilities, ritual chambers, and mana debris.',
    compatibleSources: ['boss_enemy', 'room'],
    flavor: 'Component-rich caches and unstable curios.',
  },
  {
    value: 'wilderness',
    label: 'Wilderness',
    description: 'Campsites, ruins in the wild, and hunting grounds.',
    compatibleSources: ['common_enemy', 'boss_enemy', 'room'],
    flavor: 'Balanced drops from scavenged and natural sources.',
  },
];

export const LOOT_REWARD_OPTIONS: LootRewardOption[] = [
  { value: 'scarce', label: 'Scarce', description: 'Low-expectation payout with sparse drops.' },
  { value: 'standard', label: 'Standard', description: 'Default reward curve and risk profile.' },
  { value: 'bountiful', label: 'Bountiful', description: 'Generous haul with elevated drop quality.' },
  { value: 'jackpot', label: 'Jackpot', description: 'High-roll fantasy mode with premium outcomes.' },
];

export function deriveLootLevelBand(level: number): LootLevelBand {
  if (level <= 4) return 'novice';
  if (level <= 9) return 'adventurer';
  if (level <= 14) return 'veteran';
  return 'legend';
}

export function getThemesForSource(source: LootSource): LootThemeOption[] {
  return LOOT_THEME_OPTIONS.filter((theme) => theme.compatibleSources.includes(source));
}

export function recommendedTierForReward(rewardAmount: LootRewardAmount): LootTier {
  switch (rewardAmount) {
    case 'scarce':
      return 'low';
    case 'jackpot':
      return 'high';
    default:
      return 'medium';
  }
}
