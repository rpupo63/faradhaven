import { useState } from 'react';
import { Coins, Plus, Minus, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { NormalizedCharacterSheet } from '@/types/game';
import { toast } from 'sonner';

interface MoneyPanelProps {
  sheet: NormalizedCharacterSheet;
  onMoneyChange?: (newTotal: number) => void | Promise<void>;
}

export function CurrencyDisplay({ cp, className = "" }: { cp: number; className?: string }) {
  const gp = Math.floor(cp / 100);
  const sp = Math.floor((cp % 100) / 10);
  const remainingCp = cp % 10;

  return (
    <div className={`flex items-center gap-3 font-tome-marginalia text-sm ${className}`}>
      {gp > 0 && (
        <div className="flex items-center gap-1">
          <span className="text-faded-gold font-bold">{gp}</span>
          <span className="text-muted-foreground text-micro uppercase">gp</span>
        </div>
      )}
      {sp > 0 && (
        <div className="flex items-center gap-1">
          <span className="text-slate-400 font-bold">{sp}</span>
          <span className="text-muted-foreground text-micro uppercase">sp</span>
        </div>
      )}
      <div className="flex items-center gap-1">
        <span className="text-orange-700 font-bold">{remainingCp}</span>
        <span className="text-muted-foreground text-micro uppercase">cp</span>
      </div>
    </div>
  );
}

export function MoneyPanel({ sheet, onMoneyChange }: MoneyPanelProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [amount, setAmount] = useState('0');
  const [unit, setUnit] = useState<'gp' | 'sp' | 'cp'>('gp');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleUpdateMoney = async (isAdding: boolean) => {
    if (!onMoneyChange) return;
    
    const value = parseInt(amount, 10);
    if (isNaN(value) || value <= 0) return;

    let cpDelta = value;
    if (unit === 'gp') cpDelta *= 100;
    if (unit === 'sp') cpDelta *= 10;

    const newMoney = isAdding 
      ? sheet.money + cpDelta 
      : Math.max(0, sheet.money - cpDelta);

    setIsSubmitting(true);
    try {
      await onMoneyChange(newMoney);
      setIsEditing(false);
      setAmount('0');
      toast.success(`Updated money: ${isAdding ? '+' : '-'}${value} ${unit}`);
    } catch (error) {
      toast.error('Failed to update money');
      console.error(error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between bg-muted/20 p-2 rounded border border-border/50">
        <div className="flex items-center gap-2">
          <Coins className="h-4 w-4 text-faded-gold" />
          <span className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Wealth</span>
        </div>
        
        <div className="flex items-center gap-3">
          <CurrencyDisplay cp={sheet.money} />
          {onMoneyChange && (
            <Button 
              variant="ghost" 
              size="sm" 
              className="h-6 w-6 p-0" 
              onClick={() => setIsEditing(!isEditing)}
            >
              <Plus className={`h-3 w-3 transition-transform ${isEditing ? 'rotate-45' : ''}`} />
            </Button>
          )}
        </div>
      </div>

      {isEditing && (
        <div className="p-3 bg-card border border-primary/20 rounded-lg animate-in fade-in slide-in-from-top-1 duration-200">
          <div className="flex gap-2">
            <div className="flex-1 flex gap-1">
              <Input
                type="number"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                className="h-8 text-sm"
                min="0"
              />
              <select 
                value={unit} 
                onChange={(e) => setUnit(e.target.value as 'gp' | 'sp' | 'cp')}
                className="h-8 bg-background border border-input rounded px-1 text-xs outline-none focus:ring-1 focus:ring-ring"
              >
                <option value="gp">GP</option>
                <option value="sp">SP</option>
                <option value="cp">CP</option>
              </select>
            </div>
            <div className="flex gap-1">
              <Button 
                size="sm" 
                variant="outline" 
                className="h-8 w-8 p-0" 
                onClick={() => handleUpdateMoney(false)}
                disabled={isSubmitting}
              >
                {isSubmitting ? <Loader2 className="h-3 w-3 animate-spin" /> : <Minus className="h-3 w-3" />}
              </Button>
              <Button 
                size="sm" 
                variant="default" 
                className="h-8 w-8 p-0" 
                onClick={() => handleUpdateMoney(true)}
                disabled={isSubmitting}
              >
                {isSubmitting ? <Loader2 className="h-3 w-3 animate-spin" /> : <Plus className="h-3 w-3" />}
              </Button>
            </div>
          </div>
          <p className="text-micro text-muted-foreground mt-2 text-center italic">
            Add or subtract from your total wealth.
          </p>
        </div>
      )}
    </div>
  );
}
