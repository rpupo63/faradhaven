import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Label } from '@/components/ui/label';
import { type LevelUpPreview } from '@/lib/api';

interface ModularChoiceSelectionProps {
  preview: LevelUpPreview;
  onSelect: (choice: { mp: number; br: number }) => void;
}

interface Choice {
  label: string;
  mp: number;
  br: number;
}

function parseFeatureDescription(description: string): Choice[] {
    // Example: "Choose one: (A) Doctor's Intuition (+1 MP) or (B) Primal Senses (+1 BR)."
    const choices: Choice[] = [];
    const parts = description.split(' or ');
    for (const part of parts) {
        const labelMatch = part.match(/\) (.*?)\s*\(/);
        const mpMatch = part.match(/\(\+([0-9]+)\s*MP\)/);
        const brMatch = part.match(/\(\+([0-9]+)\s*BR\)/);

        if (labelMatch) {
            const label = labelMatch[1].trim();
            const mp = mpMatch ? parseInt(mpMatch[1], 10) : 0;
            const br = brMatch ? parseInt(brMatch[1], 10) : 0;
            choices.push({ label, mp, br });
        }
    }
    return choices;
}


export function ModularChoiceSelection({ preview, onSelect }: ModularChoiceSelectionProps) {
  const modularChoiceFeature = preview.class_level.level_features?.find(f => f.name.startsWith('Modular Choice'));

  if (!modularChoiceFeature) {
    return <div>No modular choices available for this level.</div>;
  }
  
  const choices = parseFeatureDescription(modularChoiceFeature.description);

  return (
    <div className="space-y-4">
      <h3 className="text-lg font-medium text-primary">{modularChoiceFeature.name}</h3>
      <RadioGroup
        onValueChange={(value) => {
          const choice = choices[parseInt(value, 10)];
          onSelect({ mp: choice.mp, br: choice.br });
        }}
      >
        {choices.map((choice, index) => (
          <div key={index} className="flex items-center space-x-2">
            <RadioGroupItem value={index.toString()} id={`choice-${index}`} />
            <Label htmlFor={`choice-${index}`}>{choice.label} (+{choice.mp > 0 ? `${choice.mp} MP` : `${choice.br} BR`})</Label>
          </div>
        ))}
      </RadioGroup>
    </div>
  );
}
