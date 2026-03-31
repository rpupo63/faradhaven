import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getCharacterSpells, getSpells, castSpell, getCharacterSheet } from '@/lib/api'; // NEW: Added getSpells
import { hasAllComponents } from '@/lib/spellUtils';
import { CharacterSpellForge } from '@/components/CharacterSpellForge';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { LoadingButton } from '@/components/ui/loading-button';
import { LoadingQuill } from '@/components/LoadingQuill';
import { Sparkles, Wand2, BookOpen, Shield, Zap, Clock, Target, Brain, Heart, Swords, Edit, AlertTriangle, Timer } from 'lucide-react';
import type { ApiSpell, ApiComponent, ApiCharacterComponent, ApiCharacterSheet, ApiClassResource, CastSpellResponse } from '@/types/game';
import { toast } from 'sonner';
import { Switch } from '@/components/ui/switch'; // NEW: Import Switch component
import { Label } from '@/components/ui/label';   // NEW: Import Label for Switch
import { SpellListPagination } from './SpellListPagination'; // NEW: Import SpellListPagination
import { buildCastToast, RESOURCE_KEY_LABELS } from '@/lib/toastUtils';
import { formatSpellRangeFeet, formatSpellDamageDice } from '@/lib/spellMechanics';

function getClassResourceValue(resources: ApiClassResource[] | undefined, key: string): number | undefined {
  return resources?.find(r => r.key === key)?.current_value;
}

function getClassResourceMaxOrValue(resources: ApiClassResource[] | undefined, key: string): number | undefined {
  const r = resources?.find(r => r.key === key);
  return r?.max_value ?? r?.value;
}

function getPoolResource(sheet?: ApiCharacterSheet): ApiClassResource | undefined {
  return sheet?.class_resources?.find(r => r.category === 'pool' && r.is_trackable);
}


function isPoolComponent(compId: string, classComponents: ApiComponent[] | undefined): boolean {
  return classComponents?.some(c => c.id === compId) ?? false;
}

interface SpellCostInfo {
  cost: number;
  label: string;
  hasEnough: boolean;
  warningText?: string;
  /** Mutagen: whether this cast would trigger a madness save */
  madnessSaveRequired?: boolean;
}

/**
 * Computes cost, label, availability, and warning for a spell based on the character's class.
 */
function getSpellCostInfo(sheet: ApiCharacterSheet | undefined, spell: ApiSpell): SpellCostInfo {
  if (!sheet) return { cost: spell.slot_level, label: 'Slot', hasEnough: true };

  const className = sheet.class?.name;
  const poolResource = getPoolResource(sheet);

  switch (className) {
    case 'The Rift Weaver': {
      const cost = (spell.components?.length || 0) * 2;
      const current = poolResource?.current_value ?? 0;
      return { cost, label: 'SP', hasEnough: current >= cost };
    }

    case 'The Piston Brawler': {
      const cost = spell.slot_level;
      const current = poolResource?.current_value ?? 0;
      return { cost, label: 'Stab', hasEnough: current >= cost };
    }

    case 'The Sanguinist': {
      const cost = spell.slot_level;
      const current = poolResource?.current_value ?? 0;
      return { cost, label: 'Ichor', hasEnough: current >= cost };
    }

    case 'The Vapor Blade': {
      const cost = spell.slot_level;
      const current = poolResource?.current_value ?? 0;
      return { cost, label: 'SP', hasEnough: current >= cost };
    }

    case 'The Ironwright':
    case 'The Lorewright': {
      return { cost: spell.components?.length || 0, label: 'Comp', hasEnough: true };
    }

    case 'The Mutagen': {
      const profBonus = sheet.class_level?.proficiency_bonus || 2;
      const conMod = Math.floor(((sheet.character?.constitution ?? 10) - 10) / 2);
      const castCount = sheet.character?.madness_cast_count ?? 0;
      const threshold = profBonus + conMod;
      const componentCount = spell.components?.length || 0;

      const willExceedCastLimit = castCount + 1 > threshold;
      const exceedsComponentLimit = componentCount > profBonus;
      const madnessSaveRequired = spell.slot_level >= 1 && (willExceedCastLimit || exceedsComponentLimit);

      let warningText: string | undefined;
      if (madnessSaveRequired) {
        const dc = sheet.character?.current_madness_dc ?? 10;
        warningText = `CON save DC ${dc} or go Feral`;
      }

      return {
        cost: 0,
        label: 'Madness',
        hasEnough: true,
        warningText,
        madnessSaveRequired,
      };
    }

    case 'The Powder Mage': {
      return { cost: 0, label: 'Speed Dial', hasEnough: true };
    }

    default: {
      if (poolResource) {
        const cost = spell.slot_level;
        const current = poolResource.current_value ?? 0;
        const label = RESOURCE_KEY_LABELS[poolResource.key] || poolResource.display_name;
        return { cost, label, hasEnough: current >= cost };
      }
      return { cost: spell.slot_level, label: 'Slot', hasEnough: true };
    }
  }
}

interface PreparedSpellsProps {
  characterId: string;
  token?: string;
  characterName?: string;
}

const PAGINATION_LIMIT = 10; // Number of spells to display per page

export function PreparedSpells({ characterId, token, characterName }: PreparedSpellsProps) {
  const queryClient = useQueryClient();
  const [editingSpell, setEditingSpell] = useState<ApiSpell | null>(null);
  const [showAllSpells, setShowAllSpells] = useState(false); // NEW: State for toggle
  const [pointsFilter, setPointsFilter] = useState<number | 'all'>('all'); // NEW: pointsFilter state
  const [currentPage, setCurrentPage] = useState(1); // NEW: currentPage state

  // Conditional fetching of spells
  const { data: paginatedCharacterSpells, isLoading: isCharacterSpellsLoading, error: characterSpellsError } = useQuery({
    queryKey: ['character-spells', characterId, { page: currentPage, limit: PAGINATION_LIMIT, pointsFilter }],
    queryFn: () => getCharacterSpells(characterId, token, { page: currentPage, limit: PAGINATION_LIMIT, slot_level: pointsFilter }),
    enabled: !!characterId && !!token && !showAllSpells,
    staleTime: 30_000,
  });

  const { data: paginatedAllSpells, isLoading: isAllSpellsLoading, error: allSpellsError } = useQuery({
    queryKey: ['all-spells', { page: currentPage, limit: PAGINATION_LIMIT, pointsFilter }],
    queryFn: () => getSpells(token, { page: currentPage, limit: PAGINATION_LIMIT, slot_level: pointsFilter }),
    enabled: !!token && showAllSpells,
    staleTime: 30_000,
  });

  const spells = showAllSpells ? paginatedAllSpells?.spells : paginatedCharacterSpells?.spells;
  const totalCount = showAllSpells ? (paginatedAllSpells?.total_count || 0) : (paginatedCharacterSpells?.total_count || 0);
  const totalPages = Math.ceil(totalCount / PAGINATION_LIMIT);

  const isSpellsLoading = showAllSpells ? isAllSpellsLoading : isCharacterSpellsLoading;
  const spellsError = showAllSpells ? allSpellsError : characterSpellsError;

  const { data: sheet, isLoading: isSheetLoading } = useQuery({
    queryKey: ['character-sheet', characterId],
    queryFn: () => getCharacterSheet(characterId, token),
    enabled: !!characterId && !!token,
    staleTime: 30_000,
  });

  const castMutation = useMutation({
    mutationFn: (spell: ApiSpell) => castSpell(characterId, spell.slot_level, token, spell.id),
    onSuccess: (data, spell) => {
      const message = buildCastToast(sheet, spell.name, data);

      const isMadnessWarning =
        sheet?.class?.name === 'The Mutagen' &&
        data.madness_cast_count !== undefined &&
        data.madness_cast_count > (sheet?.class_level?.proficiency_bonus || 2) +
          Math.floor(((sheet?.character?.constitution ?? 10) - 10) / 2);

      if (isMadnessWarning) {
        toast.warning(message);
      } else {
        toast.success(message);
      }

      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      queryClient.invalidateQueries({ queryKey: ['character', characterId] });
      queryClient.invalidateQueries({ queryKey: ['character-spells', characterId] });
      queryClient.invalidateQueries({ queryKey: ['all-spells'] }); // Invalidate all-spells too
    },
    onError: (err: Error) => {
      const errorMsg = err.message || 'Failed to cast spell';
      toast.error(errorMsg);
      console.error(err);
    },
  });

  const isLoading = isSpellsLoading || isSheetLoading;
  const error = spellsError;

  if (!token) {
    return (
      <Card className="arcane-border bg-card">
        <CardContent className="py-12 text-center">
          <Sparkles className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
          <h3 className="text-lg font-tome-heading text-primary mb-2">Login Required</h3>
          <p className="text-muted-foreground font-tome-marginalia">
            Please log in to view prepared spells.
          </p>
        </CardContent>
      </Card>
    );
  }

  if (isLoading) {
    return (
      <Card className="arcane-border bg-card">
        <CardContent className="py-6">
          <LoadingQuill label="Loading prepared spells..." />
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className="arcane-border bg-card">
        <CardContent className="py-12 text-center">
          <Wand2 className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
          <h3 className="text-lg font-tome-heading text-primary mb-2">Error</h3>
          <p className="text-muted-foreground font-tome-marginalia">
            Failed to load prepared spells.
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <>
      <Card className="arcane-border bg-card">
        <CardHeader className="pb-3">
          <CardTitle className="flex items-center gap-2 text-lg font-tome-heading text-primary">
            <Sparkles className="h-5 w-5" />
            {characterName ? `${characterName}'s Spells` : 'Spells'}
            <Badge variant="secondary" className="ml-auto">
              {spells?.length || 0} {spells?.length === 1 ? 'spell' : 'spells'}
            </Badge>
          </CardTitle>
          {/* NEW: Toggle for showing all spells */}
          <div className="flex items-center space-x-2 mt-2">
            <Switch
              id="show-all-spells"
              checked={showAllSpells}
              onCheckedChange={setShowAllSpells}
            />
            <Label htmlFor="show-all-spells">Show All Spells</Label>
          </div>
        </CardHeader>
        <CardContent>
          {/* Using SpellListPagination */}
          <SpellListPagination
            initialSpells={spells}
            isLoading={isLoading}
            error={error}
            emptyMessage={
              showAllSpells
                ? "No spells have been crafted in the realm."
                : `${characterName || "This character"} hasn't prepared any spells yet. Use the Spell Forge to create spells.`
            }
            headerTitle={showAllSpells ? "All Spells" : `${characterName || "Prepared"} Spells`}
            showFilters={showAllSpells} // Only show filters for "All Spells" view
            pointsFilter={pointsFilter} // NEW: Pass pointsFilter
            setPointsFilter={setPointsFilter} // NEW: Pass setPointsFilter
            currentPage={currentPage} // NEW: Pass currentPage
            setCurrentPage={setCurrentPage} // NEW: Pass setCurrentPage
            totalPages={totalPages} // Corrected: Pass the calculated totalPages
            renderSpellItem={(spell: ApiSpell) => ( // Custom render prop
              <PreparedSpellCard 
                key={spell.id} 
                spell={spell} 
                sheet={sheet}
                onCast={() => castMutation.mutate(spell)}
                onEdit={() => setEditingSpell(spell)}
                isCasting={castMutation.isPending && castMutation.variables?.id === spell.id}
                isDisabled={!hasAllComponents(spell, sheet?.class?.components, sheet?.components)}
              />
            )}
          />
        </CardContent>
      </Card>

      <Dialog open={!!editingSpell} onOpenChange={(isOpen) => { if (!isOpen) setEditingSpell(null); }}>
        <DialogContent className="max-w-7xl w-full h-[90vh] flex flex-col">
          <DialogHeader>
            <DialogTitle>Edit Spell: {editingSpell?.name}</DialogTitle>
          </DialogHeader>
          <ScrollArea className="flex-grow pr-6 -mr-6">
            {editingSpell && sheet && (
              <CharacterSpellForge
                spellToEdit={editingSpell}
                onSaveComplete={() => setEditingSpell(null)}
                userId={sheet.character.user_id}
                characterId={characterId}
                token={token}
                availableComponents={sheet.class?.components}
                components={sheet.components}
                currentStability={getClassResourceValue(sheet.class_resources, 'max_stability')}
                maxStability={getClassResourceMaxOrValue(sheet.class_resources, 'max_stability')}
                maxBlueprintSlots={getClassResourceMaxOrValue(sheet.class_resources, 'speed_dial_slots')}
              />
            )}
          </ScrollArea>
        </DialogContent>
      </Dialog>
    </>
  );
}

export interface PreparedSpellCardProps {
  spell: ApiSpell;
  sheet?: ApiCharacterSheet;
  onCast?: () => void;
  onEdit?: () => void;
  isCasting?: boolean;
  isDisabled?: boolean;
}

export function PreparedSpellCard({ spell, sheet, onCast, onEdit, isCasting, isDisabled: isComponentsMissing }: PreparedSpellCardProps) {
  const characterComponents = sheet?.components;
  const classComponents = sheet?.class?.components;
  const costInfo = getSpellCostInfo(sheet, spell);

  const isDisabled = isComponentsMissing || !costInfo.hasEnough;

  const getMod = (ability: string) => {
    if (!sheet?.character) return 0;
    let score: number;
    switch (ability.toLowerCase()) {
      case 'strength': score = sheet.character.strength; break;
      case 'dexterity': score = sheet.character.dexterity; break;
      case 'constitution': score = sheet.character.constitution; break;
      case 'intelligence': score = sheet.character.intelligence; break;
      case 'wisdom': score = sheet.character.wisdom; break;
      case 'charisma': score = sheet.character.charisma; break;
      default: score = 10;
    }
    return Math.floor((score - 10) / 2);
  };

  const primaryAbility = sheet?.class?.primary_ability || 'intelligence';
  const spellMod = getMod(primaryAbility);
  const profBonus = sheet?.class_level?.proficiency_bonus || 2;
  const saveDC = 8 + profBonus + spellMod;

  const renderCostBadge = () => {
    const { cost, label, hasEnough } = costInfo;
    const insufficientClass = !hasEnough ? 'text-destructive border-destructive/50' : 'bg-background/50';

    if (label === 'Madness') {
      return (
        <Badge
          variant={costInfo.madnessSaveRequired ? 'warning' : 'outline'}
          theme={costInfo.madnessSaveRequired ? undefined : "background-subtle"}
          className="shrink-0"
        >
          Lvl {spell.slot_level}
        </Badge>
      );
    }

    if (label === 'Speed Dial') {
      return (
        <Badge variant="outline" theme="background-subtle" className="shrink-0">
          <Timer className="h-3 w-3 mr-1" />
          Speed Dial
        </Badge>
      );
    }

    if (label === 'Comp') {
      return (
        <Badge variant="outline" theme="background-subtle" className="shrink-0">
          {cost} Comp
        </Badge>
      );
    }

    return (
      <Badge variant="outline" className={`shrink-0 ${insufficientClass}`}>
        {cost} {label}
      </Badge>
    );
  };

  const renderMechanics = () => {
    const parts = [];
    
    if (spell.type && spell.type !== 'Utility') {
      let icon = <Zap className="h-3 w-3 mr-1" />;
      let themeVariant: "spell-type-default" | "spell-type-attack" | "spell-type-healing" | "spell-type-save" = "spell-type-default";

      if (spell.type === 'Attack') { icon = <Swords className="h-3 w-3 mr-1" />; themeVariant="spell-type-attack"; }
      if (spell.type === 'Healing') { icon = <Heart className="h-3 w-3 mr-1" />; themeVariant="spell-type-healing"; }
      if (spell.type === 'Save') { icon = <Shield className="h-3 w-3 mr-1" />; themeVariant="spell-type-save"; }

      parts.push(
        <Badge key="type" variant="outline" theme={themeVariant} className="mr-2">
          {icon} {spell.type}
        </Badge>
      );
    }

    if (typeof spell.range === 'number') {
      parts.push(
        <span key="range" className="flex items-center text-xs text-muted-foreground mr-3" title="Range">
          <Target className="h-3 w-3 mr-1" /> {formatSpellRangeFeet(spell.range)}
        </span>
      );
    }

    if (spell.duration) {
      parts.push(
        <span key="duration" className="flex items-center text-xs text-muted-foreground mr-3" title="Duration">
          <Clock className="h-3 w-3 mr-1" /> {spell.duration}
        </span>
      );
    }

    if (spell.concentration) {
      parts.push(
        <Badge key="conc" variant="secondary" size="h5-sm" className="mr-2" title="Requires Concentration">
          <Brain className="h-3 w-3 mr-1" /> Conc.
        </Badge>
      );
    }

    return <div className="flex flex-wrap items-center mt-2 mb-2">{parts}</div>;
  };

  const renderImpact = () => {
    const hasDice =
      typeof spell.damage_dice_count === 'number' && typeof spell.damage_die_size === 'number';
    if (!hasDice && !spell.save_attr) return null;

    return (
      <div className="bg-muted/30 p-2 rounded-md mb-3 flex items-center justify-between text-sm">
        {hasDice && (
          <div className="flex items-center gap-2">
             <span className="font-bold text-primary">
               {formatSpellDamageDice(spell.damage_dice_count!, spell.damage_die_size!)}
               {spell.add_modifier && spellMod > 0 && ` + ${spellMod}`}
             </span>
             <span className="text-muted-foreground text-xs uppercase tracking-wider">
               {spell.damage_type} {spell.type === 'Healing' ? 'Healing' : 'Damage'}
             </span>
          </div>
        )}
        
        {spell.save_attr && (
          <div className="flex items-center gap-2">
             <span className="font-bold text-amber-500">
               DC {saveDC}
             </span>
             <span className="text-muted-foreground text-xs uppercase tracking-wider">
               {spell.save_attr} Save
             </span>
          </div>
        )}
      </div>
    );
  };

  const renderMadnessWarning = () => {
    if (!costInfo.madnessSaveRequired || !costInfo.warningText) return null;
    return (
      <div className="flex items-center gap-2 px-2 py-1.5 rounded-md bg-amber-500/10 border border-amber-500/30 mb-3 text-xs text-amber-400">
        <AlertTriangle className="h-3.5 w-3.5 shrink-0" />
        <span>{costInfo.warningText}</span>
      </div>
    );
  };

  const castButtonLabel = isComponentsMissing
    ? 'Missing Components'
    : !costInfo.hasEnough
      ? `Insufficient ${costInfo.label}`
      : 'Cast';

  return (
    <div className={`p-4 rounded-lg border border-border bg-card/50 hover:bg-card/80 transition-all group ${isDisabled ? 'opacity-50 grayscale blur-[1px]' : ''} ${costInfo.madnessSaveRequired ? 'border-amber-500/30' : ''}`}>
      <div className="flex items-start justify-between gap-4 mb-2">
        <div>
          <h4 className="font-tome-heading text-primary text-lg leading-none">{spell.name}</h4>
          <p className="text-xs text-muted-foreground font-tome-marginalia mt-1">
            Level {spell.slot_level} Evocation
          </p>
        </div>
        <div className="flex items-center gap-2">
          {renderCostBadge()}
          {onCast && (
            <LoadingButton
              size="sm"
              variant="secondary"
              className="h-7 px-3 text-xs"
              onClick={onCast}
              isLoading={isCasting}
              disabled={isDisabled}
              loadingText="Casting..."
            >
              <Wand2 className="h-3 w-3 mr-1" />
              {castButtonLabel}
            </LoadingButton>
          )}
          <Button
            size="icon"
            variant="ghost"
            className="h-7 w-7 opacity-0 group-hover:opacity-100 transition-opacity"
            onClick={onEdit}
            title="Edit Spell"
          >
            <Edit className="h-4 w-4" />
          </Button>
        </div>
      </div>
      
      {renderMadnessWarning()}
      {renderMechanics()}
      {renderImpact()}
      
      {spell.description && (
        <p className="text-sm text-muted-foreground mb-3 line-clamp-3">
          {spell.description}
        </p>
      )}
      
      {spell.components && spell.components.length > 0 && (
        <div className="flex flex-wrap gap-1 mt-auto pt-2 border-t border-border/30">
          {spell.components.map((component: ApiComponent) => {
            const inPool = isPoolComponent(component.id, classComponents);
            if (inPool) {
              return (
                <ComponentBadge 
                  key={component.id} 
                  component={component} 
                />
              );
            }
            const charComp = characterComponents?.find(cc => cc.component_id === component.id);
            const count = charComp?.count || 0;
            return (
              <ComponentBadge 
                key={component.id} 
                component={component} 
                count={count}
                isMissing={count === 0}
              />
            );
          })}
        </div>
      )}
    </div>
  );
}

interface ComponentBadgeProps {
  component: ApiComponent;
  count?: number;
  isMissing?: boolean;
}

const categoryColors: Record<string, string> = {
  Forma: 'bg-yellow-400/20 text-yellow-800 border-yellow-400/30',
  Scopus: 'bg-sky-400/20 text-sky-800 border-sky-400/30',
  Essentia: 'bg-amber-400/20 text-amber-800 border-amber-400/30',
  Actio: 'bg-red-400/20 text-red-800 border-red-400/30',
  Magnitudo: 'bg-violet-400/20 text-violet-800 border-violet-400/30',
};

function ComponentBadge({ component, count, isMissing }: ComponentBadgeProps) {
  const colorClass = categoryColors[component.category] || 'bg-primary/20 text-foreground border-primary/30';
  
  return (
    <span
      className={`inline-flex items-center gap-1 px-2 py-0.5 text-xs rounded-full border ${colorClass} ${isMissing ? 'border-destructive/50 text-destructive' : ''}`}
      title={`${component.name} (${component.category})`}
    >
      <span className="font-mono font-bold">{component.symbol}</span>
      <span className="hidden sm:inline">{component.name}</span>
      {count !== undefined && (
        <Badge variant="secondary" className="h-4 px-1 text-[10px] min-w-[1.2rem] flex justify-center">
          {count}
        </Badge>
      )}
    </span>
  );
}