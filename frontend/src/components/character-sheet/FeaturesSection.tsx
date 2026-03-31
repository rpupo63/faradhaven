import { useState, useEffect } from 'react';
import { ScrollText, Zap, Languages, ChevronDown, ChevronUp, ChevronsDown, ChevronsUp, FileText, Save, Loader2, Brain } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { NormalizedCharacterSheet } from '@/types/game';
import { cn } from '@/lib/utils';

interface FeaturesSectionProps {
  sheet: NormalizedCharacterSheet;
  onNotesChange?: (notes: string) => void | Promise<void>;
}

export function FeaturesSection({ sheet, onNotesChange }: FeaturesSectionProps) {
  const { character, race_traits } = sheet;
  const cls = sheet.class;

  // State for expanded traits
  const [expandedTraits, setExpandedTraits] = useState<Set<string>>(new Set());
  const [proficienciesExpanded, setProficienciesExpanded] = useState(false);
  const [madnessExpanded, setMadnessExpanded] = useState(false);

  // State for notes editing
  const [isEditingNotes, setIsEditingNotes] = useState(false);
  const [notesDraft, setNotesDraft] = useState(character.notes || '');
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    setNotesDraft(character.notes || '');
  }, [character.notes]);

  const handleSaveNotes = async () => {
    if (!onNotesChange) return;
    setIsSaving(true);
    try {
      await onNotesChange(notesDraft);
      setIsEditingNotes(false);
    } catch (err) {
      console.error('Failed to save notes:', err);
    } finally {
      setIsSaving(false);
    }
  };

  const toggleTrait = (id: string) => {
    const newExpanded = new Set(expandedTraits);
    if (newExpanded.has(id)) {
      newExpanded.delete(id);
    } else {
      newExpanded.add(id);
    }
    setExpandedTraits(newExpanded);
  };

  const expandAll = () => {
    if (race_traits) {
      setExpandedTraits(new Set(race_traits.map(t => t.id)));
    }
  };

  const collapseAll = () => {
    setExpandedTraits(new Set());
  };

  // State for expanded features
  const [expandedFeatures, setExpandedFeatures] = useState<Set<string>>(new Set());

  const toggleFeature = (id: string) => {
    const newExpanded = new Set(expandedFeatures);
    if (newExpanded.has(id)) {
      newExpanded.delete(id);
    } else {
      newExpanded.add(id);
    }
    setExpandedFeatures(newExpanded);
  };

  const expandAllFeatures = () => {
    if (sheet.class_level.level_features) {
      setExpandedFeatures(new Set(sheet.class_level.level_features.map(f => f.id)));
    }
  };

  const collapseAllFeatures = () => {
    setExpandedFeatures(new Set());
  };

  const filteredRaceTraits = (race_traits || []).filter(
    (trait) => trait.name.toLowerCase() !== 'ability score improvement'
  );

  const filteredClassFeatures = (sheet.class_level.level_features || []).filter(
    (feature) => feature.name.toLowerCase() !== 'ability score improvement'
  );

  return (
    <div className="space-y-4">
      {/* Madness / Feral Table */}
      {sheet.madness_table && (
        <Card className="arcane-border bg-card">
          <CardHeader
            className="pb-2 cursor-pointer hover:bg-primary/5 transition-colors"
            onClick={() => setMadnessExpanded(!madnessExpanded)}
          >
            <CardTitle className="flex items-center justify-between text-base font-tome-subheading text-primary">
              <div className="flex items-center gap-2">
                <Brain className="h-4 w-4" />
                {sheet.class.name === 'The Mutagen' ? 'Feral Table' : 'Madness Table'}
              </div>
              {madnessExpanded ? (
                <ChevronUp className="h-4 w-4 text-muted-foreground" />
              ) : (
                <ChevronDown className="h-4 w-4 text-muted-foreground" />
              )}
            </CardTitle>
          </CardHeader>
          {madnessExpanded && (
            <CardContent className="animate-in fade-in slide-in-from-top-2 duration-200 p-0">
              <div className="overflow-y-auto max-h-[300px]">
                <table className="w-full text-sm">
                  <thead className="bg-muted/50 text-xs uppercase text-muted-foreground font-tome-marginalia sticky top-0">
                    <tr>
                      <th className="px-3 py-2 text-center w-12">
                        {sheet.class.name === 'The Mutagen' ? 'd20' : 'd50'}
                      </th>
                      <th className="px-3 py-2 text-left">Effect</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border/50">
                    {Object.entries(sheet.madness_table)
                      .sort((a, b) => Number(a[0]) - Number(b[0]))
                      .map(([roll, effect]) => (
                        <tr key={roll} className="hover:bg-muted/30">
                          <td className="px-3 py-2 text-center font-tome-subheading text-primary">{roll}</td>
                          <td className="px-3 py-2 font-tome-marginalia text-foreground/90">{effect}</td>
                        </tr>
                      ))}
                  </tbody>
                </table>
              </div>
            </CardContent>
          )}
        </Card>
      )}

      {/* Proficiencies */}
      {(cls.proficiencies || cls.tools?.length || cls.skill_focus?.length) && (
        <Card className="arcane-border bg-card">
          <CardHeader 
            className="pb-2 cursor-pointer hover:bg-primary/5 transition-colors"
            onClick={() => setProficienciesExpanded(!proficienciesExpanded)}
          >
            <CardTitle className="flex items-center justify-between text-base font-tome-subheading text-primary">
              <div className="flex items-center gap-2">
                <ScrollText className="h-4 w-4" />
                Proficiencies
              </div>
              {proficienciesExpanded ? (
                <ChevronUp className="h-4 w-4 text-muted-foreground" />
              ) : (
                <ChevronDown className="h-4 w-4 text-muted-foreground" />
              )}
            </CardTitle>
          </CardHeader>
          {proficienciesExpanded && (
            <CardContent className="space-y-3 animate-in fade-in slide-in-from-top-2 duration-200">
              {cls.proficiencies && (
                <div>
                  <p className="text-xs text-muted-foreground font-tome-marginalia mb-1">Weapons & Armor</p>
                  <p className="text-sm font-tome-marginalia">{cls.proficiencies}</p>
                </div>
              )}
              {cls.tools && cls.tools.length > 0 && (
                <div>
                  <p className="text-xs text-muted-foreground font-tome-marginalia mb-1">Tools</p>
                  <p className="text-sm font-tome-marginalia">{cls.tools.join(', ')}</p>
                </div>
              )}
              {cls.skill_focus && cls.skill_focus.length > 0 && (
                <div>
                  <p className="text-xs text-muted-foreground font-tome-marginalia mb-1">Skill Focus</p>
                  <p className="text-sm font-tome-marginalia">{cls.skill_focus.join(', ')}</p>
                </div>
              )}
            </CardContent>
          )}
        </Card>
      )}

      {/* Racial Traits */}
      {filteredRaceTraits.length > 0 && (
        <Card className="arcane-border bg-card">
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center gap-2 text-base font-tome-subheading text-primary">
                <Zap className="h-4 w-4" />
                Racial Traits
              </CardTitle>
              <div className="flex gap-1">
                <Button 
                  variant="ghost" 
                  size="icon" 
                  className="h-6 w-6" 
                  onClick={expandAll} 
                  title="Expand All"
                >
                  <ChevronsDown className="h-3 w-3" />
                </Button>
                <Button 
                  variant="ghost" 
                  size="icon" 
                  className="h-6 w-6" 
                  onClick={collapseAll} 
                  title="Collapse All"
                >
                  <ChevronsUp className="h-3 w-3" />
                </Button>
              </div>
            </div>
            <p className="text-xs text-muted-foreground font-tome-marginalia">
              {character.lineageName || character.raceName} heritage
            </p>
          </CardHeader>
          <CardContent>
            <ul className="space-y-3">
              {filteredRaceTraits.map((trait) => {
                const isExpanded = expandedTraits.has(trait.id);
                return (
                  <li key={trait.id} className="border-l-2 border-primary/30 pl-3">
                    <div 
                      className="flex items-center gap-2 cursor-pointer hover:text-primary/80 transition-colors"
                      onClick={() => toggleTrait(trait.id)}
                    >
                      {isExpanded ? (
                        <ChevronUp className="h-3 w-3 flex-shrink-0 text-muted-foreground" />
                      ) : (
                        <ChevronDown className="h-3 w-3 flex-shrink-0 text-muted-foreground" />
                      )}
                      <span className="font-tome-subheading text-primary text-sm select-none">{trait.name}</span>
                      {trait.action_type && (
                        <span className="text-xs text-muted-foreground ml-1">
                          ({trait.action_type})
                        </span>
                      )}
                    </div>
                    {isExpanded && trait.description && (
                      <p className="text-xs font-tome-marginalia text-foreground/90 mt-1 pl-5 animate-in fade-in zoom-in-95 duration-200">
                        {trait.description}
                      </p>
                    )}
                  </li>
                );
              })}
            </ul>
          </CardContent>
        </Card>
      )}

      {/* Class Features */}
      <Card className="arcane-border bg-card">
        <CardHeader className="pb-2">
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2 text-base font-tome-subheading text-primary">
              <Zap className="h-4 w-4" />
              Class Features
            </CardTitle>
            <div className="flex gap-1">
              <Button 
                variant="ghost" 
                size="icon" 
                className="h-6 w-6" 
                onClick={expandAllFeatures} 
                title="Expand All"
              >
                <ChevronsDown className="h-3 w-3" />
              </Button>
              <Button 
                variant="ghost" 
                size="icon" 
                className="h-6 w-6" 
                onClick={collapseAllFeatures} 
                title="Collapse All"
              >
                <ChevronsUp className="h-3 w-3" />
              </Button>
            </div>
          </div>
          <p className="text-xs text-muted-foreground font-tome-marginalia">
            {character.className} abilities
          </p>
        </CardHeader>
        <CardContent>
          {filteredClassFeatures.length > 0 ? (
            <ul className="space-y-3">
              {filteredClassFeatures.map((feature) => {
                const isExpanded = expandedFeatures.has(feature.id);
                return (
                  <li key={feature.id} className="border-l-2 border-primary/30 pl-3">
                    <div 
                      className="flex items-center gap-2 cursor-pointer hover:text-primary/80 transition-colors"
                      onClick={() => toggleFeature(feature.id)}
                    >
                      {isExpanded ? (
                        <ChevronUp className="h-3 w-3 flex-shrink-0 text-muted-foreground" />
                      ) : (
                        <ChevronDown className="h-3 w-3 flex-shrink-0 text-muted-foreground" />
                      )}
                      <span className="font-tome-subheading text-primary text-sm select-none">{feature.name}</span>
                    </div>
                    {isExpanded && feature.description && (
                      <p className="text-xs font-tome-marginalia text-foreground/90 mt-1 pl-5 animate-in fade-in zoom-in-95 duration-200">
                        {feature.description}
                      </p>
                    )}
                  </li>
                );
              })}
            </ul>
          ) : (
            <p className="text-sm font-tome-marginalia text-muted-foreground italic">
              No features gained at this level
            </p>
          )}
        </CardContent>
      </Card>

      {/* Languages */}
      <Card className="arcane-border bg-card">
        <CardHeader className="pb-2">
          <CardTitle className="flex items-center gap-2 text-base font-tome-subheading text-primary">
            <Languages className="h-4 w-4" />
            Languages
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm font-tome-marginalia text-muted-foreground">
            Common
          </p>
        </CardContent>
      </Card>

      {/* Character Notes */}
      <Card 
        className={cn(
          "arcane-border bg-card transition-colors",
          !isEditingNotes && onNotesChange && "hover:bg-primary/5 cursor-pointer"
        )}
        onClick={() => {
          if (!isEditingNotes && onNotesChange) {
            setIsEditingNotes(true);
          }
        }}
      >
        <CardHeader className="pb-2">
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2 text-base font-tome-subheading text-primary">
              <FileText className="h-4 w-4" />
              Notes
            </CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          {isEditingNotes ? (
            <div className="space-y-3" onClick={(e) => e.stopPropagation()}>
              <Textarea
                autoFocus
                value={notesDraft}
                onChange={(e) => setNotesDraft(e.target.value)}
                placeholder="Add miscellaneous notes about your character..."
                className="min-h-[150px] text-sm font-tome-marginalia bg-background/50"
              />
              <div className="flex justify-end gap-2">
                <Button 
                  variant="outline" 
                  size="sm" 
                  onClick={() => {
                    setIsEditingNotes(false);
                    setNotesDraft(character.notes || '');
                  }}
                  disabled={isSaving}
                >
                  Cancel
                </Button>
                <Button 
                  size="sm" 
                  onClick={handleSaveNotes}
                  disabled={isSaving || notesDraft === character.notes}
                  className="gap-2"
                >
                  {isSaving ? (
                    <Loader2 className="h-3 w-3 animate-spin" />
                  ) : (
                    <Save className="h-3 w-3" />
                  )}
                  Save
                </Button>
              </div>
            </div>
          ) : (
            <div className="text-sm font-tome-marginalia text-muted-foreground whitespace-pre-wrap">
              {character.notes ? (
                character.notes
              ) : (
                <span className="italic">No notes yet. Click here to add some.</span>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
