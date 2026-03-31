import { useState } from 'react';
import { BestiaryEntry, Attack, CreatureSize, CreatureType, DamageType, CREATURE_SIZES, CREATURE_TYPES, DAMAGE_TYPES } from '@/types/game';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Plus, Trash2, Skull, X } from 'lucide-react';

interface BestiaryFormProps {
  entry?: BestiaryEntry;
  onSubmit: (data: Omit<BestiaryEntry, 'id' | 'createdAt'>) => void;
  onCancel: () => void;
}

export function BestiaryForm({ entry, onSubmit, onCancel }: BestiaryFormProps) {
  const [name, setName] = useState(entry?.name || '');
  const [imageUrl, setImageUrl] = useState(entry?.imageUrl || '');
  const [size, setSize] = useState<CreatureSize>(entry?.size || 'Medium');
  const [type, setType] = useState<CreatureType>(entry?.type || 'Humanoid');
  const [alignment, setAlignment] = useState(entry?.alignment || 'Neutral');
  const [armorClass, setArmorClass] = useState(entry?.armorClass || 10);
  const [hitPoints, setHitPoints] = useState(entry?.hitPoints || 10);
  const [hitDice, setHitDice] = useState(entry?.hitDice || '2d8');
  const [speed, setSpeed] = useState(entry?.speed || '30 ft.');
  const [strength, setStrength] = useState(entry?.strength || 10);
  const [dexterity, setDexterity] = useState(entry?.dexterity || 10);
  const [constitution, setConstitution] = useState(entry?.constitution || 10);
  const [intelligence, setIntelligence] = useState(entry?.intelligence || 10);
  const [wisdom, setWisdom] = useState(entry?.wisdom || 10);
  const [charisma, setCharisma] = useState(entry?.charisma || 10);
  const [challengeRating, setChallengeRating] = useState(entry?.challengeRating || '1');
  const [attacks, setAttacks] = useState<Attack[]>(entry?.attacks || []);
  const [abilities, setAbilities] = useState(entry?.abilities || '');
  const [description, setDescription] = useState(entry?.description || '');

  const addAttack = () => {
    const newAttack: Attack = {
      id: crypto.randomUUID(),
      name: 'New Attack',
      attackBonus: 4,
      damageType: 'Slashing',
      damageDice: '1d6+2',
      reach: '5 ft.',
    };
    setAttacks([...attacks, newAttack]);
  };

  const updateAttack = (id: string, updates: Partial<Attack>) => {
    setAttacks(attacks.map(a => a.id === id ? { ...a, ...updates } : a));
  };

  const removeAttack = (id: string) => {
    setAttacks(attacks.filter(a => a.id !== id));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    onSubmit({
      name: name.trim(),
      imageUrl: imageUrl.trim() || undefined,
      size,
      type,
      alignment: alignment.trim(),
      armorClass,
      hitPoints,
      hitDice: hitDice.trim(),
      speed: speed.trim(),
      strength,
      dexterity,
      constitution,
      intelligence,
      wisdom,
      charisma,
      challengeRating: challengeRating.trim(),
      attacks,
      abilities: abilities.trim() || undefined,
      description: description.trim(),
    });
  };

  return (
    <Card className="arcane-border bg-card/90 backdrop-blur-sm">
      <CardHeader className="pb-4">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2 font-display text-primary">
            <Skull className="w-5 h-5" />
            {entry ? 'Edit Creature' : 'Add New Creature'}
          </CardTitle>
          <Button variant="ghost" size="icon" onClick={onCancel}>
            <X className="w-4 h-4" />
          </Button>
        </div>
      </CardHeader>

      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Basic Info */}
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="name">Creature Name *</Label>
              <Input
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Goblin, Dragon, etc."
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="imageUrl">Image URL</Label>
              <Input
                id="imageUrl"
                value={imageUrl}
                onChange={(e) => setImageUrl(e.target.value)}
                placeholder="https://example.com/monster.png"
              />
            </div>
          </div>

          {/* Type & Size */}
          <div className="grid gap-4 md:grid-cols-3">
            <div className="space-y-2">
              <Label>Size</Label>
              <Select value={size} onValueChange={(v) => setSize(v as CreatureSize)}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {CREATURE_SIZES.map((s) => (
                    <SelectItem key={s} value={s}>{s}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>Type</Label>
              <Select value={type} onValueChange={(v) => setType(v as CreatureType)}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {CREATURE_TYPES.map((t) => (
                    <SelectItem key={t} value={t}>{t}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="alignment">Alignment</Label>
              <Input
                id="alignment"
                value={alignment}
                onChange={(e) => setAlignment(e.target.value)}
                placeholder="Chaotic Evil"
              />
            </div>
          </div>

          {/* Combat Stats */}
          <div className="grid gap-4 md:grid-cols-4">
            <div className="space-y-2">
              <Label htmlFor="ac">Armor Class</Label>
              <Input
                id="ac"
                type="number"
                value={armorClass}
                onChange={(e) => setArmorClass(parseInt(e.target.value) || 10)}
                min={1}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="hp">Hit Points</Label>
              <Input
                id="hp"
                type="number"
                value={hitPoints}
                onChange={(e) => setHitPoints(parseInt(e.target.value) || 1)}
                min={1}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="hitDice">Hit Dice</Label>
              <Input
                id="hitDice"
                value={hitDice}
                onChange={(e) => setHitDice(e.target.value)}
                placeholder="2d8+4"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cr">Challenge Rating</Label>
              <Input
                id="cr"
                value={challengeRating}
                onChange={(e) => setChallengeRating(e.target.value)}
                placeholder="1/4, 1, 5, etc."
              />
            </div>
          </div>

          {/* Speed */}
          <div className="space-y-2">
            <Label htmlFor="speed">Speed</Label>
            <Input
              id="speed"
              value={speed}
              onChange={(e) => setSpeed(e.target.value)}
              placeholder="30 ft., fly 60 ft."
            />
          </div>

          {/* Ability Scores */}
          <div>
            <Label className="mb-3 block">Ability Scores</Label>
            <div className="grid grid-cols-6 gap-2">
              {[
                { label: 'STR', value: strength, setter: setStrength },
                { label: 'DEX', value: dexterity, setter: setDexterity },
                { label: 'CON', value: constitution, setter: setConstitution },
                { label: 'INT', value: intelligence, setter: setIntelligence },
                { label: 'WIS', value: wisdom, setter: setWisdom },
                { label: 'CHA', value: charisma, setter: setCharisma },
              ].map(({ label, value, setter }) => (
                <div key={label} className="space-y-1">
                  <Label className="text-xs text-center block text-muted-foreground">{label}</Label>
                  <Input
                    type="number"
                    value={value}
                    onChange={(e) => setter(parseInt(e.target.value) || 1)}
                    min={1}
                    max={30}
                    className="text-center"
                  />
                </div>
              ))}
            </div>
          </div>

          {/* Attacks */}
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <Label>Attacks</Label>
              <Button type="button" variant="outline" size="sm" onClick={addAttack}>
                <Plus className="w-4 h-4" />
                Add Attack
              </Button>
            </div>
            {attacks.map((attack) => (
              <div key={attack.id} className="p-3 rounded-lg border border-border bg-muted/30 space-y-3">
                <div className="flex items-center justify-between">
                  <Input
                    value={attack.name}
                    onChange={(e) => updateAttack(attack.id, { name: e.target.value })}
                    placeholder="Attack name"
                    className="max-w-[200px]"
                  />
                  <Button type="button" variant="ghost" size="icon" onClick={() => removeAttack(attack.id)}>
                    <Trash2 className="w-4 h-4 text-destructive" />
                  </Button>
                </div>
                <div className="grid gap-2 md:grid-cols-4">
                  <div className="space-y-1">
                    <Label className="text-xs">Attack Bonus</Label>
                    <Input
                      type="number"
                      value={attack.attackBonus}
                      onChange={(e) => updateAttack(attack.id, { attackBonus: parseInt(e.target.value) || 0 })}
                    />
                  </div>
                  <div className="space-y-1">
                    <Label className="text-xs">Damage</Label>
                    <Input
                      value={attack.damageDice}
                      onChange={(e) => updateAttack(attack.id, { damageDice: e.target.value })}
                      placeholder="1d6+2"
                    />
                  </div>
                  <div className="space-y-1">
                    <Label className="text-xs">Damage Type</Label>
                    <Select
                      value={attack.damageType}
                      onValueChange={(v) => updateAttack(attack.id, { damageType: v as DamageType })}
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        {DAMAGE_TYPES.map((dt) => (
                          <SelectItem key={dt} value={dt}>{dt}</SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-1">
                    <Label className="text-xs">Reach</Label>
                    <Input
                      value={attack.reach || ''}
                      onChange={(e) => updateAttack(attack.id, { reach: e.target.value })}
                      placeholder="5 ft."
                    />
                  </div>
                </div>
              </div>
            ))}
          </div>

          {/* Abilities & Description */}
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="abilities">Special Abilities</Label>
              <Textarea
                id="abilities"
                value={abilities}
                onChange={(e) => setAbilities(e.target.value)}
                placeholder="Darkvision 60 ft., Pack Tactics, etc."
                rows={3}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Description / Lore</Label>
              <Textarea
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="A fearsome creature that lurks in the shadows..."
                rows={3}
              />
            </div>
          </div>

          {/* Actions */}
          <div className="flex gap-3 justify-end pt-4 border-t border-border">
            <Button type="button" variant="outline" onClick={onCancel}>
              Cancel
            </Button>
            <Button type="submit" disabled={!name.trim()}>
              {entry ? 'Update Creature' : 'Add to Bestiary'}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
