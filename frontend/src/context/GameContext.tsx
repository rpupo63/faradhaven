import React, { createContext, useContext, useEffect, useState } from 'react';
import { Character, CraftedSpell, GameState, ElementType, BestiaryEntry, SPELL_ELEMENTS } from '@/types/game';

interface GameContextType {
  state: GameState;
  activeCharacter: Character | null;
  createCharacter: (character: Omit<Character, 'id' | 'spellbook' | 'createdAt'>) => void;
  updateCharacter: (id: string, updates: Partial<Character>) => void;
  deleteCharacter: (id: string) => void;
  setActiveCharacter: (id: string | null) => void;
  craftSpell: (name: string, description: string, elements: ElementType[]) => CraftedSpell;
  deleteSpell: (id: string) => void;
  assignSpellToCharacter: (characterId: string, spellId: string) => void;
  removeSpellFromCharacter: (characterId: string, spellId: string) => void;
  getSpellById: (id: string) => CraftedSpell | undefined;
  addBestiaryEntry: (entry: Omit<BestiaryEntry, 'id' | 'createdAt'>) => void;
  updateBestiaryEntry: (id: string, updates: Partial<BestiaryEntry>) => void;
  deleteBestiaryEntry: (id: string) => void;
}

const GameContext = createContext<GameContextType | undefined>(undefined);

const STORAGE_KEY = 'faradhaven-game-state';

const getInitialState = (): GameState => {
  if (typeof window === 'undefined') {
    return { characters: [], activeCharacterId: null, craftedSpells: [], bestiary: [] };
  }
  
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored) {
    try {
      const parsed = JSON.parse(stored);
      return { ...parsed, bestiary: parsed.bestiary || [] };
    } catch {
      return { characters: [], activeCharacterId: null, craftedSpells: [], bestiary: [] };
    }
  }
  return { characters: [], activeCharacterId: null, craftedSpells: [], bestiary: [] };
};

export function GameProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<GameState>(getInitialState);

  // Persist to localStorage
  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
  }, [state]);

  const activeCharacter = state.characters.find(c => c.id === state.activeCharacterId) || null;

  const createCharacter = (character: Omit<Character, 'id' | 'spellbook' | 'createdAt'>) => {
    const newCharacter: Character = {
      ...character,
      id: crypto.randomUUID(),
      spellbook: [],
      createdAt: Date.now(),
    };
    setState(prev => ({
      ...prev,
      characters: [...prev.characters, newCharacter],
      activeCharacterId: newCharacter.id,
    }));
  };

  const updateCharacter = (id: string, updates: Partial<Character>) => {
    setState(prev => ({
      ...prev,
      characters: prev.characters.map(c => 
        c.id === id ? { ...c, ...updates } : c
      ),
    }));
  };

  const deleteCharacter = (id: string) => {
    setState(prev => ({
      ...prev,
      characters: prev.characters.filter(c => c.id !== id),
      activeCharacterId: prev.activeCharacterId === id 
        ? (prev.characters.find(c => c.id !== id)?.id || null)
        : prev.activeCharacterId,
    }));
  };

  const setActiveCharacter = (id: string | null) => {
    setState(prev => ({ ...prev, activeCharacterId: id }));
  };

  const craftSpell = (name: string, description: string, elements: ElementType[]): CraftedSpell => {
    const newSpell: CraftedSpell = {
      id: crypto.randomUUID(),
      name,
      description,
      elements,
      level: elements.length,
      createdAt: Date.now(),
    };

    setState(prev => ({
      ...prev,
      craftedSpells: [...prev.craftedSpells, newSpell],
    }));

    return newSpell;
  };

  const deleteSpell = (id: string) => {
    setState(prev => ({
      ...prev,
      craftedSpells: prev.craftedSpells.filter(s => s.id !== id),
      characters: prev.characters.map(c => ({
        ...c,
        spellbook: c.spellbook.filter(sId => sId !== id),
      })),
    }));
  };

  const assignSpellToCharacter = (characterId: string, spellId: string) => {
    setState(prev => ({
      ...prev,
      characters: prev.characters.map(c => 
        c.id === characterId && !c.spellbook.includes(spellId)
          ? { ...c, spellbook: [...c.spellbook, spellId] }
          : c
      ),
    }));
  };

  const removeSpellFromCharacter = (characterId: string, spellId: string) => {
    setState(prev => ({
      ...prev,
      characters: prev.characters.map(c => 
        c.id === characterId
          ? { ...c, spellbook: c.spellbook.filter(id => id !== spellId) }
          : c
      ),
    }));
  };

  const getSpellById = (id: string) => state.craftedSpells.find(s => s.id === id);

  // Bestiary functions
  const addBestiaryEntry = (entry: Omit<BestiaryEntry, 'id' | 'createdAt'>) => {
    const newEntry: BestiaryEntry = {
      ...entry,
      id: crypto.randomUUID(),
      createdAt: Date.now(),
    };
    setState(prev => ({
      ...prev,
      bestiary: [...prev.bestiary, newEntry],
    }));
  };

  const updateBestiaryEntry = (id: string, updates: Partial<BestiaryEntry>) => {
    setState(prev => ({
      ...prev,
      bestiary: prev.bestiary.map(e => 
        e.id === id ? { ...e, ...updates } : e
      ),
    }));
  };

  const deleteBestiaryEntry = (id: string) => {
    setState(prev => ({
      ...prev,
      bestiary: prev.bestiary.filter(e => e.id !== id),
    }));
  };

  return (
    <GameContext.Provider value={{
      state,
      activeCharacter,
      createCharacter,
      updateCharacter,
      deleteCharacter,
      setActiveCharacter,
      craftSpell,
      deleteSpell,
      assignSpellToCharacter,
      removeSpellFromCharacter,
      getSpellById,
      addBestiaryEntry,
      updateBestiaryEntry,
      deleteBestiaryEntry,
    }}>
      {children}
    </GameContext.Provider>
  );
}

export function useGame() {
  const context = useContext(GameContext);
  if (!context) {
    throw new Error('useGame must be used within a GameProvider');
  }
  return context;
}
