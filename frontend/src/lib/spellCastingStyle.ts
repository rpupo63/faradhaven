/** How this class satisfies Faradhaven casting at the table (matches backend seed values). */
export type SpellCastingComponentType = 'verbal' | 'somatic' | 'material';

export function normalizeSpellCastingComponent(
  raw: string | undefined
): SpellCastingComponentType | undefined {
  if (!raw) return undefined;
  const s = raw.toLowerCase().trim();
  if (s === 'verbal' || s === 'somatic' || s === 'material') return s;
  return undefined;
}

export const SPELL_CASTING_COMPONENT_LABELS: Record<SpellCastingComponentType, string> = {
  verbal: 'Verbal',
  somatic: 'Somatic',
  material: 'Material',
};

export const SPELL_CASTING_ABBREV: Record<SpellCastingComponentType, string> = {
  verbal: 'V',
  somatic: 'S',
  material: 'M',
};

export type ComponentPlayGuidance = {
  advantages: string[];
  risks: string[];
};

/** Table-facing pros/cons for each component; class flavor comes from the API. */
export const SPELL_CASTING_GUIDANCE: Record<SpellCastingComponentType, ComponentPlayGuidance> = {
  verbal: {
    advantages: [
      'If your hands are bound or full, you can still complete a cast as long as you can speak clearly.',
      'You do not need a free hand or visible kit to start the effect—only your voice and breath.',
    ],
    risks: [
      'Area silence, gags, or being underwater can shut you down as surely as losing a focus.',
      'Whispering is quieter but may not satisfy a full verbal cast; loud work ruins stealth and pinpoints you.',
    ],
  },
  somatic: {
    advantages: [
      'You can often cast in dead silence if your body is free to move—no incantation to hear.',
      'Subtle micro-gestures may fly under notice in a crowd if the rules at your table allow reduced motion.',
    ],
    risks: [
      'Obvious motion, stance shifts, and contortions are hard to hide—stealth and stakeouts can be blown.',
      'Restraints, grappling, or packed-in armor can block the motion your form needs.',
    ],
  },
  material: {
    advantages: [
      'Tactile work with real reagents and tools can feel more reliable when the table tracks hands and props.',
      'If your voice is gone, you may still progress a cast by kit and contact (when the situation allows).',
    ],
    risks: [
      'You need space, time, and often a free hand—search, disarm, or strip your kit and the cast stalls.',
      'The table can see you preparing: telltale vials, scrap, and lab motion are not subtle by default.',
    ],
  },
};
