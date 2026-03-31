import React, { useState, useEffect, useMemo } from 'react';
import { Check, ChevronsUpDown } from 'lucide-react';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { getAllCharacters } from '@/lib/api/character';
import { getMonstersByUser } from '@/lib/api/monster';
import { getPartiesByUserId } from '@/lib/api/party';
import type { ApiCharacter } from '@/types/game';
import type { Monster } from '@/types/monster';

export interface TokenEntity {
  id: string;
  name: string;
  image_url?: string;
  entity_type: 'character' | 'monster' | 'party';
  owner_user_id?: string;
  subtitle?: string;
  /** Populated when entity_type === 'party' */
  party_members?: TokenEntity[];
}

interface EntityComboboxProps {
  value: TokenEntity | null;
  onSelect: (entity: TokenEntity) => void;
  userId: string;
  authToken: string;
}

export const EntityCombobox: React.FC<EntityComboboxProps> = ({
  value,
  onSelect,
  userId,
  authToken,
}) => {
  const [open, setOpen] = useState(false);
  const [characters, setCharacters] = useState<ApiCharacter[]>([]);
  const [monsters, setMonsters] = useState<Monster[]>([]);
  const [parties, setParties] = useState<import('@/types/game/api').ApiParty[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    let cancelled = false;
    setTimeout(() => setLoading(true), 0);

    Promise.all([
      getAllCharacters(authToken),
      getMonstersByUser(userId, authToken),
      getPartiesByUserId(userId, authToken),
    ])
      .then(([charsRes, mons, partiesRes]) => {
        if (!cancelled) {
          setCharacters(charsRes.characters || []);
          setMonsters(mons);
          setParties(partiesRes);
        }
      })
      .catch(console.error)
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    return () => { cancelled = true; };
  }, [userId, authToken]);

  const { yourChars, otherChars, monsterEntities, partyEntities } = useMemo(() => {
    const yours: TokenEntity[] = [];
    const others: TokenEntity[] = [];

    for (const c of characters) {
      const subtitle = [
        c.level ? `Lvl ${c.level}` : null,
        c.race?.name,
        c.class?.name,
      ].filter(Boolean).join(' ');

      const entity: TokenEntity = {
        id: c.id,
        name: c.name,
        image_url: c.image_url,
        entity_type: 'character',
        owner_user_id: c.user_id,
        subtitle,
      };

      if (c.user_id === userId) {
        yours.push(entity);
      } else {
        others.push(entity);
      }
    }

    const mons: TokenEntity[] = monsters.map((m) => ({
      id: m.id,
      name: m.name,
      image_url: m.image_url,
      entity_type: 'monster',
      subtitle: `CR ${m.challenge_rating} ${m.type}`,
    }));

    const partyEnts: TokenEntity[] = parties.map((p) => {
      const memberCount = p.members?.length ?? 0;
      const members: TokenEntity[] = (p.members ?? []).map((c) => {
        const subtitle = [
          c.level ? `Lvl ${c.level}` : null,
          c.race?.name,
          c.class?.name,
        ].filter(Boolean).join(' ');
        return {
          id: c.id,
          name: c.name,
          image_url: c.image_url,
          entity_type: 'character',
          owner_user_id: c.user_id,
          subtitle,
        };
      });
      return {
        id: p.id,
        name: p.name,
        entity_type: 'party',
        subtitle: `${memberCount} member${memberCount !== 1 ? 's' : ''}`,
        party_members: members,
      };
    });

    return { yourChars: yours, otherChars: others, monsterEntities: mons, partyEntities: partyEnts };
  }, [characters, monsters, parties, userId]);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="w-full justify-between bg-input border-input text-left font-normal"
        >
          {value ? (
            <span className="truncate">{value.name}</span>
          ) : (
            <span className="text-muted-foreground">
              {loading ? 'Loading...' : 'Select entity...'}
            </span>
          )}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[--radix-popover-trigger-width] p-0" align="start">
        <Command>
          <CommandInput placeholder="Search characters & monsters..." />
          <CommandList>
            <CommandEmpty>No results found.</CommandEmpty>

            {partyEntities.length > 0 && (
              <CommandGroup heading="Parties">
                {partyEntities.map((entity) => (
                  <CommandItem
                    key={`party-${entity.id}`}
                    value={`party ${entity.name} ${entity.subtitle ?? ''}`}
                    onSelect={() => {
                      onSelect(entity);
                      setOpen(false);
                    }}
                  >
                    <Check
                      className={cn(
                        'mr-2 h-4 w-4',
                        value?.id === entity.id ? 'opacity-100' : 'opacity-0'
                      )}
                    />
                    <div className="flex flex-col">
                      <span>{entity.name}</span>
                      {entity.subtitle && (
                        <span className="text-xs text-muted-foreground">{entity.subtitle}</span>
                      )}
                    </div>
                  </CommandItem>
                ))}
              </CommandGroup>
            )}

            {yourChars.length > 0 && (
              <CommandGroup heading="Your Characters">
                {yourChars.map((entity) => (
                  <CommandItem
                    key={`char-${entity.id}`}
                    value={`${entity.name} ${entity.subtitle ?? ''}`}
                    onSelect={() => {
                      onSelect(entity);
                      setOpen(false);
                    }}
                  >
                    <Check
                      className={cn(
                        'mr-2 h-4 w-4',
                        value?.id === entity.id ? 'opacity-100' : 'opacity-0'
                      )}
                    />
                    <div className="flex flex-col">
                      <span>{entity.name}</span>
                      {entity.subtitle && (
                        <span className="text-xs text-muted-foreground">{entity.subtitle}</span>
                      )}
                    </div>
                  </CommandItem>
                ))}
              </CommandGroup>
            )}

            {otherChars.length > 0 && (
              <CommandGroup heading="Other Characters">
                {otherChars.map((entity) => (
                  <CommandItem
                    key={`char-${entity.id}`}
                    value={`${entity.name} ${entity.subtitle ?? ''}`}
                    onSelect={() => {
                      onSelect(entity);
                      setOpen(false);
                    }}
                  >
                    <Check
                      className={cn(
                        'mr-2 h-4 w-4',
                        value?.id === entity.id ? 'opacity-100' : 'opacity-0'
                      )}
                    />
                    <div className="flex flex-col">
                      <span>{entity.name}</span>
                      {entity.subtitle && (
                        <span className="text-xs text-muted-foreground">{entity.subtitle}</span>
                      )}
                    </div>
                  </CommandItem>
                ))}
              </CommandGroup>
            )}

            {monsterEntities.length > 0 && (
              <CommandGroup heading="Beasts and NPCs">
                {monsterEntities.map((entity) => (
                  <CommandItem
                    key={`mon-${entity.id}`}
                    value={`${entity.name} ${entity.subtitle ?? ''}`}
                    onSelect={() => {
                      onSelect(entity);
                      setOpen(false);
                    }}
                  >
                    <Check
                      className={cn(
                        'mr-2 h-4 w-4',
                        value?.id === entity.id ? 'opacity-100' : 'opacity-0'
                      )}
                    />
                    <div className="flex flex-col">
                      <span>{entity.name}</span>
                      {entity.subtitle && (
                        <span className="text-xs text-muted-foreground">{entity.subtitle}</span>
                      )}
                    </div>
                  </CommandItem>
                ))}
              </CommandGroup>
            )}
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
};
