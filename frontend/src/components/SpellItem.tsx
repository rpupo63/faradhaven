import React from 'react';
import { RaIcon } from '@/components/ui/RaIcon';
import { cn } from '@/lib/utils';
import { formatSpellRangeFeet, formatSpellDamageDice } from '@/lib/spellMechanics';
import type { SpellbookListSpell } from '@/components/spellbook/spellbookTypes';
import { getSpellComponentCount } from '@/lib/spellUtils';

interface SpellItemProps {
  spell: SpellbookListSpell;
}

export const SpellItem: React.FC<SpellItemProps> = ({ spell }) => {
  const compCount = getSpellComponentCount(spell);
  return (
    <div className="arcane-scroll-box bg-card/50 text-card-foreground rounded-lg p-2 mb-2">
      <div className="flex justify-between items-center">
        <div className="flex items-center space-x-2">
          <RaIcon name="aura" className="text-base text-purple-400" />
          <span className="font-tome-item-heading text-lg">{spell.name}</span>
          <span className="text-sm text-muted-foreground">
            ({compCount > 0 ? `${compCount} comp.` : '—'})
          </span>
        </div>
      </div>

      <div className="overflow-hidden pt-2">
        <p className="text-sm text-muted-foreground mb-2">{spell.description || "No description provided."}</p>
        
        <div className="grid grid-cols-2 gap-x-4 gap-y-1 text-sm font-tome-marginalia">
          {spell.type && (
            <div className="flex items-center gap-1">
              <RaIcon name="aura" className="text-sm text-gray-500" /> Type: {spell.type}
            </div>
          )}
          {typeof spell.range === 'number' && (
            <div className="flex items-center gap-1">
              <RaIcon name="archery-target" className="text-sm text-gray-500" /> Range: {formatSpellRangeFeet(spell.range)}
            </div>
          )}
          {spell.duration && (
            <div className="flex items-center gap-1">
              <RaIcon name="stopwatch" className="text-sm text-gray-500" /> Duration: {spell.duration}
            </div>
          )}
          {spell.concentration && (
            <div className="flex items-center gap-1">
              <RaIcon name="bolt-shield" className="text-sm text-gray-500" /> Concentration
            </div>
          )}
          {spell.save_attr && (
            <div className="flex items-center gap-1">
              <RaIcon name="sword" className="text-sm text-gray-500" /> Save: {spell.save_attr}
            </div>
          )}
          {typeof spell.damage_dice_count === 'number' &&
            typeof spell.damage_die_size === 'number' && (
            <div className="flex items-center gap-1">
              <RaIcon name="sword" className="text-sm text-gray-500" /> Damage:{' '}
              {formatSpellDamageDice(spell.damage_dice_count, spell.damage_die_size)} {spell.damage_type}
            </div>
          )}
          {spell.character_name && (
            <div className="flex items-center gap-1 col-span-2">
              <RaIcon name="player" className="text-sm text-green-400" /> From: {spell.character_name} ({spell.character_class})
            </div>
          )}
        </div>
        
        {spell.components && spell.components.length > 0 && (
          <div className="mt-2 text-sm font-tome-marginalia">
            <p className="font-semibold">Components:</p>
            <ul className="list-disc list-inside pl-2">
              {spell.components.map(comp => (
                <li key={comp.id}>{comp.name}</li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
};