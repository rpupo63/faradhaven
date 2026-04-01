import { useEffect, useMemo, useState } from 'react';
import { cn } from '@/lib/utils';
import { categoryColors } from './ElementTile'; // We'll reuse the colors
import { categoryOrder } from './ElementTable';
import type { ApiComponent } from '@/types/game';

// Standard ANSI Keyboard Layout (approximate)
const KEY_ROWS = [
  ['1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-', '='],
  ['q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p', '[', ']'],
  ['a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l', ';', "'"],
  ['z', 'x', 'c', 'v', 'b', 'n', 'm', ',', '.', '/'],
];

const SHIFT_MAP: Record<string, string> = {
  '1': '!', '2': '@', '3': '#', '4': '$', '5': '%', '6': '^', '7': '&', '8': '*', '9': '(', '0': ')', '-': '_', '=': '+'
};

const SHIFT_KEYS = ['!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '+'];

// Flatten for easy lookup - append shift keys for additional components
const ALL_KEYS = [...KEY_ROWS.flat(), ...SHIFT_KEYS];

interface ComponentKeyboardProps {
  availableComponents: ApiComponent[];
  selectedComponents: ApiComponent[];
  onComponentSelect: (component: ApiComponent) => void;
  className?: string;
  isTimerActive?: boolean; // New prop
}

export function ComponentKeyboard({
  availableComponents,
  selectedComponents,
  onComponentSelect,
  className,
  isTimerActive = false, // Default to false
}: ComponentKeyboardProps) {
  const [isShiftPressed, setIsShiftPressed] = useState(false);

  // 1. Map components to keys
  // We want a consistent mapping if possible, but dynamic based on availability is okay for now.
  // To keep it somewhat consistent, we'll sort available components by category then name.
  
  const mappedComponents = useMemo(() => {
    // Sort components: Categories in order, then alphabetically
    const sorted = [...availableComponents].sort((a, b) => {
      const catIndexA = categoryOrder.indexOf(a.category);
      const catIndexB = categoryOrder.indexOf(b.category);
      if (catIndexA !== catIndexB) return catIndexA - catIndexB;
      return a.name.localeCompare(b.name);
    });

    // Create map: Key -> Component
    const map = new Map<string, ApiComponent>();
    
    sorted.forEach((comp, index) => {
      if (index < ALL_KEYS.length) {
        map.set(ALL_KEYS[index], comp);
      }
    });

    return map;
  }, [availableComponents]);

  // Reverse map for highlighting: Component ID -> Key
  const componentKeyMap = useMemo(() => {
    const map = new Map<string, string>();
    mappedComponents.forEach((comp, key) => {
      map.set(comp.id, key);
    });
    return map;
  }, [mappedComponents]);

  // 2. Handle Keyboard Events
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!isTimerActive) return; // Disable keyboard input if timer is not active

      if (e.key === 'Shift') {
        setIsShiftPressed(true);
      }

      // Ignore if user is typing in an input/textarea
      if (
        document.activeElement?.tagName === 'INPUT' ||
        document.activeElement?.tagName === 'TEXTAREA' ||
        (document.activeElement instanceof HTMLElement && document.activeElement.isContentEditable)
      ) {
        return;
      }

      const key = e.key.toLowerCase();
      // Handle symbols which don't have a lowercase (or are already symbols)
      const component = mappedComponents.get(e.key) || mappedComponents.get(key);

      if (component) {
        e.preventDefault(); // Prevent default browser actions for mapped keys
        onComponentSelect(component);
      }
    };

    const handleKeyUp = (e: KeyboardEvent) => {
      if (!isTimerActive) return; // Disable keyboard input if timer is not active

      if (e.key === 'Shift') {
        setIsShiftPressed(false);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    window.addEventListener('keyup', handleKeyUp);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
      window.removeEventListener('keyup', handleKeyUp);
    };
  }, [mappedComponents, onComponentSelect, isTimerActive]); // Add isTimerActive to dependencies

  // Helper to render a key
  const renderKey = (keyChar: string) => {
    const shiftKey = SHIFT_MAP[keyChar];
    const activeKey = (isShiftPressed && shiftKey) ? shiftKey : keyChar;
    
    const component = mappedComponents.get(activeKey);
    const isSelected = component ? selectedComponents.some(c => c.id === component.id) : false;
    
    if (!component) {
      // Empty key
      return (
        <div 
          key={keyChar} 
          className="w-12 h-12 m-0.5 rounded-md bg-muted/20 border border-border/30 flex items-center justify-center opacity-30 select-none relative"
        >
          <span className="text-micro absolute top-0.5 left-1 font-mono uppercase text-muted-foreground">{keyChar}</span>
        </div>
      );
    }

    const colors = categoryColors[component.category] || categoryColors.Essentia;

    return (
      <button
        key={keyChar}
        onClick={() => onComponentSelect(component)}
        disabled={!isTimerActive} // Disable button if timer is not active
        className={cn(
          "w-12 h-12 m-0.5 relative group rounded-md border-2 transition-all duration-100 flex flex-col items-center justify-between p-1",
          // Colors
          colors.bg,
          colors.border,
          // Hover/Active states
          isTimerActive && "hover:scale-105 hover:z-10 hover:shadow-md", // Only apply hover if active
          isTimerActive && "active:scale-95", // Only apply active if active
          // Selected state
          isSelected ? "ring-2 ring-primary ring-offset-1 z-10 scale-105" : "",
          "focus:outline-none focus:ring-2 focus:ring-primary/50",
          !isTimerActive && "opacity-50 cursor-not-allowed" // Grey out and change cursor if disabled
        )}
        title={`${component.name} (${component.symbol}) - Press '${activeKey.toUpperCase()}'`}
      >
        {/* Key Label (Top Left) */}
        <span className="absolute top-0.5 left-1 text-micro font-mono font-bold opacity-50 uppercase">
          {activeKey}
        </span>

        {/* Shift Hint (Top Right) */}
        {shiftKey && !isShiftPressed && mappedComponents.has(shiftKey) && (
          <span className="absolute top-0.5 right-1 text-nano font-mono font-bold opacity-30 uppercase">
            ⇧{shiftKey}
          </span>
        )}

        {/* Component Symbol (Center) */}
        <span className={cn("text-lg font-bold font-mono mt-1", colors.text)}>
          {component.symbol}
        </span>

        {/* Component Name (Bottom, truncated) */}
        <span className="text-micro text-center opacity-80 font-medium">
          {component.name}
        </span>
      </button>
    );
  };

  return (
    <div className={cn("w-full flex flex-col items-center gap-1 select-none p-4 rounded-xl border border-border/50 bg-card/30 backdrop-blur-sm", className)}>
      <div className="text-center mb-4">
        <h3 className="text-lg font-tome-heading text-primary">Arcane Keyboard</h3>
        <p className="text-xs text-muted-foreground font-tome-marginalia">
          Press keys to forge your spell. Be careful not to lose focus!
        </p>
        {!isTimerActive && (
          <p className="text-xs text-amber-500 font-medium mt-2">
            Keyboard input is disabled until the timer starts.
          </p>
        )}
      </div>

      <div className="flex flex-col gap-1 items-center">
        {/* Row 1: Numbers */}
        <div className="flex">
          {KEY_ROWS[0].map(renderKey)}
        </div>
        
        {/* Row 2: QWERTY... (Indent slightly) */}
        <div className="flex ml-4">
          {KEY_ROWS[1].map(renderKey)}
        </div>

        {/* Row 3: ASDF... (Indent more) */}
        <div className="flex ml-8">
          {KEY_ROWS[2].map(renderKey)}
        </div>

        {/* Row 4: ZXCV... (Indent even more) */}
        <div className="flex ml-12">
          {KEY_ROWS[3].map(renderKey)}
        </div>
      </div>

      {/* Unmapped Components Warning */}
      {availableComponents.length > ALL_KEYS.length && (
        <div className="mt-4 text-xs text-amber-500 font-medium">
          Note: {availableComponents.length - ALL_KEYS.length} components could not be mapped to the keyboard.
        </div>
      )}
    </div>
  );
}
