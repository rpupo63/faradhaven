# Caster resource balance (vs D&D 5e slot tables)

Reference spells use component counts **3 / 6 / 9 / 12** as low–high complexity stand-ins. D&D **slot weight** = sum over prepared slots of `(spell level × slots)` per long rest (cantrips excluded).

## D&D slot weight (PHB progression)

Sum of `spell level × number of slots` per long rest (prepared slots only).

| Level | Wizard (full) | Paladin (half) |
|------:|--------------:|-----------------|
| 1 | 2 | 0 |
| 5 | 16 | 8 |
| 11 | 47 | 16 |
| 17 | 71 | 36 |
| 20 | 89 | 41 |

*Wizard 20:* 4×1 + 3×(2+3+4+5) + 2×(6+7) + 8 + 9 = 89.  
*Paladin 20:* 4×1 + 3×(2+3+4) + 2×5 = 41.

## Faradhaven pools (after tuning)

**Max spell points** (most classes’ `ClassLevel.max_spell_points`): `48 + 2 × level` (levels 1–20 → 50–88).  
**Rift Weaver** spends **2 SP per component**; `CharacterResource.spell_points` matches that max and restores on short/long rest with the character column.

| Level | Max SP | RW cost @ n comps (2n) | Casts/day @ n=6 | Casts/day @ n=12 |
|------:|-------:|----------------------:|----------------:|-----------------:|
| 1 | 50 | 12 | 4 | 2 |
| 5 | 58 | 12 | 4 | 2 |
| 11 | 70 | 12 | 5 | 2 |
| 17 | 82 | 12 | 6 | 3 |
| 20 | 88 | 12 | 7 | 3 |

**Piston Brawler:** stability cost = **sum of component tiers** (not raw count). Pool `max_stability` tracks paladin-ish half-caster cadence in class seed.

**Sanguinist:** `max_blood_ichor` scales with level; high levels 17–20 plateau slightly vs pure +1/level to avoid outrunning half-caster feel alongside martial features.

**Vapor Blade:** `shadow_points` max = **Dex modifier + proficiency bonus** (min 1), per class rules; casting spends 1 shadow per component.

## Implementation notes

- Rift Weaver: `class_level_resources.spell_points` seeded per level; `castSpell` syncs `current_spell_points` after deducting `CharacterResource.spell_points`.
- Mutagen: no spell-point deduction; madness-only economy.
- Piston Brawler: server computes tier-sum cost when casting by `spell_id`.

Reseed classes after changing seeds (`GENERATE_MODELS` / your seed workflow) so `class_level_resources` picks up `spell_points` for Rift Weaver.
