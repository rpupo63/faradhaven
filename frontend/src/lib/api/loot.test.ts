import { beforeEach, describe, expect, it, vi } from 'vitest';
import { generateLootPreview, confirmLootPickup } from './loot';
import type { LootPreviewResponse, LootResult } from '@/types/game/api';

const { apiFetchMock, handleResponseMock } = vi.hoisted(() => ({
  apiFetchMock: vi.fn(),
  handleResponseMock: vi.fn(),
}));

vi.mock('./base', () => ({
  getBaseUrl: () => '',
  apiFetch: apiFetchMock,
  handleResponse: handleResponseMock,
}));

describe('loot api', () => {
  beforeEach(() => {
    apiFetchMock.mockReset();
    handleResponseMock.mockReset();
  });

  it('sends source/theme/location and loot level payload', async () => {
    const response = {} as Response;
    const expected: LootPreviewResponse = {
      session_id: 'session-1',
      loot: {
        items: [],
        weapons: [],
        gold_earned: 0,
        total_money: 0,
        message: 'ok',
        items_rolled: 0,
        weapons_rolled: 0,
        drops: [],
      },
      party_members: [{ id: 'char-1', name: 'Alice' }],
    };
    apiFetchMock.mockResolvedValueOnce(response);
    handleResponseMock.mockResolvedValueOnce(expected);

    await generateLootPreview(
      'char-1',
      'room',
      'dungeon',
      'underground',
      8,
      'token-abc'
    );

    expect(apiFetchMock).toHaveBeenCalledTimes(1);
    const [url, init] = apiFetchMock.mock.calls[0];
    expect(url).toBe('/api/character/char-1/loot');
    const body = JSON.parse((init as RequestInit).body as string);
    expect(body).toMatchObject({
      character_id: 'char-1',
      source: 'room',
      room_theme: 'dungeon',
      location: 'underground',
      loot_level: 8,
    });
  });

  it('sends loot confirmation assignments payload', async () => {
    const response = {} as Response;
    const expected: LootResult = {
      items: [],
      weapons: [],
      gold_earned: 0,
      total_money: 100,
      message: 'confirmed',
      items_rolled: 0,
      weapons_rolled: 0,
      drops: [],
    };
    apiFetchMock.mockResolvedValueOnce(response);
    handleResponseMock.mockResolvedValueOnce(expected);

    await confirmLootPickup(
      'char-1',
      'session-1',
      [{ drop_index: 0, character_id: 'char-2' }],
      'token-abc'
    );

    expect(apiFetchMock).toHaveBeenCalledTimes(1);
    const [url, init] = apiFetchMock.mock.calls[0];
    expect(url).toBe('/api/character/char-1/loot/confirm');
    const body = JSON.parse((init as RequestInit).body as string);
    expect(body).toMatchObject({
      session_id: 'session-1',
      assignments: [{ drop_index: 0, character_id: 'char-2' }],
    });
  });
});
