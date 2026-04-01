import { useState, useEffect } from 'react';
import { Dices, RotateCcw } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { ApiRace } from '@/types/game';
import { cn } from '@/lib/utils';
import { roll4d6DropLowest } from '@/lib/dice';

interface AbilityScoreSelectionProps {
  method: 'standard' | 'pointbuy' | 'roll';
  setMethod: (m: 'standard' | 'pointbuy' | 'roll') => void;
  scores: Record<string, number>;
  setScores: (s: Record<string, number> | ((prev: Record<string, number>) => Record<string, number>)) => void;
  race: ApiRace | null;
  lineageId: string | null;
}

const ABILITIES = ['strength', 'dexterity', 'constitution', 'intelligence', 'wisdom', 'charisma'];
const STANDARD_ARRAY = [15, 14, 13, 12, 10, 8];

export function AbilityScoreSelection({ method, setMethod, scores, setScores, race, lineageId }: AbilityScoreSelectionProps) {
  // Pool state for Standard Array and Roll methods
  // Stores the actual values available to be assigned
  const [pool, setPool] = useState<number[]>([]);
  
  // Maps ability name -> index in the pool array
  const [assignedIndices, setAssignedIndices] = useState<Record<string, number>>({});

  // Store detailed roll info for the "roll" method
  const [rollDetails, setRollDetails] = useState<Record<number, { rolls: number[], dropped: number }>>({});

  // Initialize Standard Array on first load if method is standard
  useEffect(() => {
    if (method === 'standard' && pool.length === 0) {
      setTimeout(() => {
        setPool([...STANDARD_ARRAY]);
        // Default assignment: 0->Str, 1->Dex, etc.
        const initialIndices: Record<string, number> = {};
        ABILITIES.forEach((ability, i) => {
          initialIndices[ability] = i;
        });
        setAssignedIndices(initialIndices);
      }, 0);
    }
  }, [method, pool.length]);

  // Sync scores when pool or assignments change (for Standard/Roll)
  useEffect(() => {
    if (method === 'pointbuy') return; // Point buy manages scores directly

    if (pool.length > 0 && Object.keys(assignedIndices).length > 0) {
      const newScores: Record<string, number> = {};
      ABILITIES.forEach(ability => {
        const index = assignedIndices[ability];
        if (index !== undefined && pool[index] !== undefined) {
          newScores[ability] = pool[index];
        } else {
          newScores[ability] = 10; // Fallback
        }
      });
      // Only update if different to prevent loops
      if (JSON.stringify(newScores) !== JSON.stringify(scores)) {
        setScores(newScores);
      }
    }
  }, [pool, assignedIndices, method, setScores, scores]);

  // Calculate total racial bonuses
  const getBonus = (ability: string) => {
    let bonus = 0;
    if (race?.ability_score_bonuses?.[ability]) bonus += race.ability_score_bonuses[ability];
    
    // Find lineage bonus
    if (lineageId && race?.traits) {
      for (const t of race.traits) {
        const opt = t.options?.find(o => o.id === lineageId);
        if (opt?.ability_score_bonuses?.[ability]) {
          bonus += opt.ability_score_bonuses[ability];
        }
      }
    }
    return bonus;
  };

  const updatePointBuyScore = (ability: string, val: number) => {
    setScores((prev: Record<string, number>) => ({ ...prev, [ability]: val }));
  };

  const resetPointBuy = () => {
    const newScores: Record<string, number> = {};
    ABILITIES.forEach(a => newScores[a] = 8);
    setScores(newScores);
  };

  const getPointCost = (score: number) => {
    if (score <= 8) return 0;
    if (score === 9) return 1;
    if (score === 10) return 2;
    if (score === 11) return 3;
    if (score === 12) return 4;
    if (score === 13) return 5;
    if (score === 14) return 7;
    if (score === 15) return 9;
    return 0;
  };

  const usedPoints = Object.values(scores).reduce((acc: number, val: number) => acc + getPointCost(val), 0);
  const remainingPoints = 27 - usedPoints;

  // Auto Roll Feature
  const handleAutoRoll = async () => {
    // Roll silently (no animation) to avoid 6 sequential 3D dice animations.
    const rawResults = await Promise.all(Array(6).fill(0).map(() => roll4d6DropLowest(true)));
    const results = rawResults.filter((r): r is NonNullable<typeof r> => r !== null);
    if (results.length < 6) return;
    setPool(results.map(r => r.total));

    const details: Record<number, { rolls: number[], dropped: number }> = {};
    results.forEach((r, i) => {
      details[i] = { rolls: r.rolls, dropped: r.dropped };
    });
    setRollDetails(details);
    
    // Auto-assign indices if not set
    if (Object.keys(assignedIndices).length === 0) {
      const newIndices: Record<string, number> = {};
      ABILITIES.forEach((ability, i) => {
        newIndices[ability] = i;
      });
      setAssignedIndices(newIndices);
    }
  };

  const handleAssignmentChange = (ability: string, poolIndexStr: string) => {
    const poolIndex = parseInt(poolIndexStr);
    
    // Swap logic: Find who currently has this index and give them my old index
    const myOldIndex = assignedIndices[ability];
    const otherAbility = Object.keys(assignedIndices).find(key => assignedIndices[key] === poolIndex && key !== ability);

    const newIndices = { ...assignedIndices, [ability]: poolIndex };
    
    if (otherAbility && myOldIndex !== undefined) {
      newIndices[otherAbility] = myOldIndex;
    }

    setAssignedIndices(newIndices);
  };

  return (
    <div className="space-y-6 h-full">
      <Tabs value={method} onValueChange={(v: string) => {
        const m = v as 'standard' | 'pointbuy' | 'roll';
        setMethod(m);
        // Reset/Init logic on tab change
        if (m === 'standard') {
            setPool([...STANDARD_ARRAY]);
            const idx: Record<string, number> = {}; ABILITIES.forEach((a, i) => idx[a] = i);
            setAssignedIndices(idx);
        } else if (m === 'pointbuy') {
            // Optional: reset to 8s? Or keep current if valid?
            // Let's not force reset to allow switching back and forth if compatible, 
            // but usually point buy starts low.
        } else if (m === 'roll') {
            // Clear pool to force roll? Or keep if existing?
            if (pool.length === 0) handleAutoRoll();
        }
      }}>
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="standard">Standard Array</TabsTrigger>
          <TabsTrigger value="pointbuy">Point Buy</TabsTrigger>
          <TabsTrigger value="roll">Manual / Roll</TabsTrigger>
        </TabsList>
      </Tabs>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mt-6">
        <div className="space-y-4">
          {ABILITIES.map(ability => {
            const base = scores[ability as keyof typeof scores];
            const bonus = getBonus(ability);
            const total = base + bonus;
            const mod = Math.floor((total - 10) / 2);

            return (
              <div key={ability} className="flex items-center justify-between p-2 rounded border bg-card">
                <div className="w-32">
                  <Label className="uppercase text-xs font-bold text-muted-foreground block mb-1">
                    {ability}
                  </Label>
                  <div className="flex items-center gap-2">
                    <span className="font-display text-2xl text-primary">{total}</span>
                    <span className="text-xs text-muted-foreground">(Mod: {mod >= 0 ? '+' : ''}{mod})</span>
                  </div>
                </div>

                <div className="flex items-center gap-3">
                  {method === 'pointbuy' ? (
                    <div className="flex items-center gap-2">
                      <Button 
                        variant="outline" size="icon" className="h-8 w-8"
                        onClick={() => updatePointBuyScore(ability, Math.max(8, base - 1))}
                        disabled={base <= 8}
                      >
                        -
                      </Button>
                      <span className="w-8 text-center font-mono">{base}</span>
                      <Button 
                        variant="outline" size="icon" className="h-8 w-8"
                        onClick={() => updatePointBuyScore(ability, Math.min(15, base + 1))}
                        disabled={base >= 15 || (getPointCost(base + 1) - getPointCost(base) > remainingPoints)}
                      >
                        +
                      </Button>
                    </div>
                  ) : (method === 'standard' || method === 'roll') ? (
                    <div className="w-24">
                       <Select 
                         value={assignedIndices[ability]?.toString()} 
                         onValueChange={(val) => handleAssignmentChange(ability, val)}
                         disabled={pool.length === 0}
                       >
                        <SelectTrigger className="h-8">
                          <SelectValue placeholder="-" />
                        </SelectTrigger>
                        <SelectContent>
                          {pool.map((val, idx) => (
                            <SelectItem key={idx} value={idx.toString()}>
                              {val}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                  ) : (
                    <Input 
                        type="number" 
                        min="1" max="20" 
                        value={base} 
                        onChange={(e) => updatePointBuyScore(ability, parseInt(e.target.value) || 8)}
                        className="w-20 text-center"
                    />
                  )}
                  
                  {bonus > 0 && (
                    <span className="text-xs text-green-500 font-bold ml-2">
                      +{bonus} Race
                    </span>
                  )}
                </div>
              </div>
            );
          })}
        </div>

        <div className="space-y-6">
          {method === 'pointbuy' && (
            <div className="p-4 border rounded-lg bg-muted/20 text-center space-y-4">
               <div>
                <div className="text-sm text-muted-foreground uppercase mb-1">Remaining Points</div>
                <div className={cn("font-display text-4xl", remainingPoints < 0 ? "text-red-500" : "text-primary")}>
                    {remainingPoints} <span className="text-lg text-muted-foreground">/ 27</span>
                </div>
              </div>
              
              <div className="grid grid-cols-2 gap-2 text-xs text-muted-foreground border-t pt-4">
                 <span>8: 0 pts</span>
                 <span>9: 1 pt</span>
                 <span>10: 2 pts</span>
                 <span>11: 3 pts</span>
                 <span>12: 4 pts</span>
                 <span>13: 5 pts</span>
                 <span>14: 7 pts</span>
                 <span>15: 9 pts</span>
              </div>

              <Button variant="outline" size="sm" onClick={resetPointBuy} className="w-full">
                <RotateCcw className="h-3 w-3 mr-2" /> Reset All to 8
              </Button>
            </div>
          )}
          
          {method === 'standard' && (
            <div className="p-4 border rounded-lg bg-muted/20">
              <h4 className="font-bold mb-2 text-sm">Standard Array</h4>
              <p className="text-xs text-muted-foreground mb-4">
                Assign the standard values to your abilities. Select a value for each ability using the dropdowns.
              </p>
              <div className="flex flex-wrap gap-2 justify-center">
                {STANDARD_ARRAY.map((v, i) => (
                  <span key={i} className={cn(
                    "px-3 py-1 rounded font-mono text-sm shadow-sm border",
                    Object.values(assignedIndices).includes(i) ? "bg-muted text-muted-foreground opacity-50" : "bg-card text-foreground"
                  )}>
                    {v}
                  </span>
                ))}
              </div>
            </div>
          )}

          {method === 'roll' && (
            <div className="p-4 border rounded-lg bg-muted/20 space-y-4">
              <div className="text-center">
                <Dices className="mx-auto h-8 w-8 text-primary mb-2" />
                <p className="text-sm text-muted-foreground mb-4">
                  Roll 4d6, drop the lowest die for each ability score.
                </p>
                <Button onClick={handleAutoRoll} className="gap-2 w-full">
                  <Dices className="h-4 w-4" />
                  Auto Roll Stats
                </Button>
              </div>
              
              {pool.length > 0 && (
                <div className="border-t pt-4">
                    <h4 className="font-bold mb-2 text-sm text-center">Your Pool</h4>
                    <div className="flex flex-wrap gap-4 justify-center">
                        {pool.map((val, idx) => (
                        <div key={idx} className="flex flex-col items-center">
                            <span className={cn(
                                "px-3 py-1 rounded font-mono text-sm shadow-sm border",
                                "bg-card text-foreground"
                            )}>
                                {val}
                            </span>
                            {rollDetails[idx] && (
                                <span className="text-micro text-muted-foreground mt-1 flex flex-col items-center">
                                    <span>[{rollDetails[idx].rolls.join(', ')}]</span>
                                    <span className="text-destructive/70 line-through" title="Dropped Die">{rollDetails[idx].dropped}</span>
                                </span>
                            )}
                        </div>
                        ))}
                    </div>
                    <p className="text-xs text-center text-muted-foreground mt-4">
                        Values are auto-assigned. Use the dropdowns on the left to swap them.
                    </p>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
