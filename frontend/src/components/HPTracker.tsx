import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Plus, Minus } from 'lucide-react';
import { RaIcon } from '@/components/ui/RaIcon';
import { cn } from '@/lib/utils';

interface HPTrackerProps {
  currentHP: number;
  maxHP: number;
  tempHP: number;
  onHPChange: (delta: number) => void | Promise<void>;
  onTempHPChange: (value: number) => void | Promise<void>;
  disabled?: boolean;
  className?: string;
}

export function HPTracker({
  currentHP,
  maxHP,
  tempHP,
  onHPChange,
  onTempHPChange,
  disabled = false,
  className,
}: HPTrackerProps) {
  const [hpInput, setHpInput] = useState('');
  const [tempHpInput, setTempHpInput] = useState('');

  // Calculate health percentage for color coding
  const healthPercent = maxHP > 0 ? (currentHP / maxHP) * 100 : 100;
  const getHealthColor = () => {
    if (healthPercent > 50) return 'bg-green-500';
    if (healthPercent > 25) return 'bg-yellow-500';
    return 'bg-red-500';
  };

  const handleDamage = () => {
    const value = parseInt(hpInput);
    if (!isNaN(value) && value > 0) {
      onHPChange(-value);
      setHpInput('');
    }
  };

  const handleHeal = () => {
    const value = parseInt(hpInput);
    if (!isNaN(value) && value > 0) {
      onHPChange(value);
      setHpInput('');
    }
  };

  const handleSetTempHP = () => {
    const value = parseInt(tempHpInput);
    if (!isNaN(value) && value >= 0) {
      onTempHPChange(value);
      setTempHpInput('');
    }
  };

  const handleQuickChange = (delta: number) => {
    onHPChange(delta);
  };

  return (
    <Card className={cn('arcane-border bg-card', className)}>
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center gap-2 text-base font-tome-subheading text-primary">
          <RaIcon name="health" className="text-sm" />
          Hit Points
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* HP Display */}
        <div className="text-center">
          <div className="font-display text-4xl text-primary">
            {currentHP}
            <span className="text-2xl text-muted-foreground">/{maxHP}</span>
          </div>
          {tempHP > 0 && (
            <div className="flex items-center justify-center gap-1 text-cyan-400 mt-1">
              <RaIcon name="shield" className="text-sm" />
              <span className="font-medium">+{tempHP} temp</span>
            </div>
          )}
        </div>

        {/* Health Bar */}
        <div className="h-3 rounded-full bg-muted overflow-hidden">
          <div
            className={cn('h-full transition-all duration-300', getHealthColor())}
            style={{ width: `${Math.min(100, healthPercent)}%` }}
          />
        </div>

        {/* Quick Damage/Heal Buttons */}
        <div className="flex items-center justify-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => handleQuickChange(-5)}
            disabled={disabled}
            className="text-red-400 hover:text-red-300 hover:border-red-500/50"
          >
            -5
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => handleQuickChange(-1)}
            disabled={disabled}
            className="text-red-400 hover:text-red-300 hover:border-red-500/50"
          >
            -1
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => handleQuickChange(1)}
            disabled={disabled}
            className="text-green-400 hover:text-green-300 hover:border-green-500/50"
          >
            +1
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => handleQuickChange(5)}
            disabled={disabled}
            className="text-green-400 hover:text-green-300 hover:border-green-500/50"
          >
            +5
          </Button>
        </div>

        {/* Custom HP Input */}
        <div className="flex items-center gap-2">
          <Input
            type="number"
            min="1"
            placeholder="Amount"
            value={hpInput}
            onChange={(e) => setHpInput(e.target.value)}
            className="flex-1"
            disabled={disabled}
          />
          <Button
            variant="outline"
            size="icon"
            onClick={handleDamage}
            disabled={disabled || !hpInput}
            className="text-red-400 hover:text-red-300 hover:border-red-500/50"
            title="Take Damage"
          >
            <Minus className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            size="icon"
            onClick={handleHeal}
            disabled={disabled || !hpInput}
            className="text-green-400 hover:text-green-300 hover:border-green-500/50"
            title="Heal"
          >
            <Plus className="h-4 w-4" />
          </Button>
        </div>

        {/* Temp HP Input */}
        <div className="pt-2 border-t border-border/50">
          <p className="text-xs text-muted-foreground mb-2">Temporary HP</p>
          <div className="flex items-center gap-2">
            <Input
              type="number"
              min="0"
              placeholder="Temp HP"
              value={tempHpInput}
              onChange={(e) => setTempHpInput(e.target.value)}
              className="flex-1"
              disabled={disabled}
            />
            <Button
              variant="outline"
              size="sm"
              onClick={handleSetTempHP}
              disabled={disabled || tempHpInput === ''}
              className="text-cyan-400 hover:text-cyan-300 hover:border-cyan-500/50"
            >
              <RaIcon name="shield" className="text-sm mr-1" />
              Set
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
