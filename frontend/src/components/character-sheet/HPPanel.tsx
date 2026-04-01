import { Minus, Plus, Shield } from 'lucide-react';
import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Separator } from '@/components/ui/separator';

interface HPPanelProps {
  currentHP: number;
  maxHP: number;
  tempHP: number;
  onHPChange: (delta: number) => void | Promise<void>;
  onTempHPChange?: (value: number) => void | Promise<void>;
}

export function HPPanel({ 
  currentHP, 
  maxHP, 
  tempHP, 
  onHPChange, 
  onTempHPChange 
}: HPPanelProps) {
  const [damageInput, setDamageInput] = useState('');
  const [tempInput, setTempInput] = useState('');
  const [isProcessing, setIsProcessing] = useState(false);

  const handleDamage = async (amount: number) => {
    if (amount > 0) {
      setIsProcessing(true);
      try {
        await onHPChange(-amount);
        setDamageInput('');
      } finally {
        setIsProcessing(false);
      }
    }
  };

  const handleHeal = async (amount: number) => {
    if (amount > 0) {
      setIsProcessing(true);
      try {
        await onHPChange(amount);
        setDamageInput('');
      } finally {
        setIsProcessing(false);
      }
    }
  };

  const handleSetTemp = async (amount: number) => {
    if (onTempHPChange && amount >= 0) {
      setIsProcessing(true);
      try {
        await onTempHPChange(amount);
        setTempInput('');
      } finally {
        setIsProcessing(false);
      }
    }
  };

  return (
    <div className="space-y-4 py-2">
      {/* Current HP Display */}
      <div className="text-center space-y-1">
        <div className="text-3xl font-display text-primary">
          {currentHP} <span className="text-lg text-muted-foreground">/ {maxHP}</span>
        </div>
        {tempHP > 0 && (
          <div className="text-sm text-blue-400 font-tome-marginalia flex items-center justify-center gap-1">
            <Shield className="h-3 w-3" />
            <span>+{tempHP} Temp HP</span>
          </div>
        )}
      </div>

      <Separator className="opacity-50" />

      {/* Custom Amount Input */}
      <div className="space-y-2">
        <p className="text-micro uppercase tracking-wider text-muted-foreground font-semibold text-center">Damage / Healing</p>
        <div className="flex items-center gap-2">
          <Input
            type="number"
            min="1"
            placeholder="Amount"
            value={damageInput}
            onChange={(e) => setDamageInput(e.target.value)}
            disabled={isProcessing}
            className="flex-1"
          />
          <Button
            variant="outline"
            size="icon"
            onClick={() => {
              const val = parseInt(damageInput);
              if (!isNaN(val)) handleDamage(val);
            }}
            disabled={!damageInput || isProcessing}
            className="text-red-400 hover:text-red-300 hover:border-red-500/50"
            title="Take Damage"
          >
            <Minus className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            size="icon"
            onClick={() => {
              const val = parseInt(damageInput);
              if (!isNaN(val)) handleHeal(val);
            }}
            disabled={!damageInput || isProcessing}
            className="text-green-400 hover:text-green-300 hover:border-green-500/50"
            title="Heal"
          >
            <Plus className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {onTempHPChange && (
        <>
          <Separator className="opacity-50" />
          <div className="space-y-2">
            <p className="text-micro uppercase tracking-wider text-muted-foreground font-semibold text-center">Temporary HP</p>
            <div className="flex items-center gap-2">
              <Input
                type="number"
                min="0"
                placeholder="New Temp HP"
                value={tempInput}
                onChange={(e) => setTempInput(e.target.value)}
                disabled={isProcessing}
                className="flex-1"
              />
              <Button
                variant="outline"
                size="sm"
                onClick={() => {
                  const val = parseInt(tempInput);
                  if (!isNaN(val)) handleSetTemp(val);
                }}
                disabled={!tempInput || isProcessing}
                className="text-blue-400 hover:text-blue-300 hover:border-blue-500/50"
              >
                Set
              </Button>
            </div>
          </div>
        </>
      )}
    </div>
  );
}