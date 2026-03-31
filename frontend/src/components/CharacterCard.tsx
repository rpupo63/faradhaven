import { cn } from '@/lib/utils';
import { Character } from '@/types/game';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { useNavigate } from 'react-router-dom';
import {
  Trash2,
  BookOpen,
  User,
  GraduationCap,
  ChevronRight,
  Star,
} from 'lucide-react';

interface CharacterCardProps {
  character: Character;
  isActive: boolean;
  onSelect: () => void;
  onDelete: () => void;
  onSetActive?: () => void;
}

export function CharacterCard({ character, isActive, onSelect, onDelete, onSetActive }: CharacterCardProps) {
  const navigate = useNavigate();
  const ClassIcon = GraduationCap;

  const handleCardClick = () => {
    onSelect();
    navigate(`/character/${character.id}/sheet`);
  };

  return (
    <div
      className={cn(
        'arcane-border rounded-lg cursor-pointer transition-all duration-300 group overflow-hidden flex flex-col',
        isActive
          ? 'ring-2 ring-primary shadow-glow'
          : 'hover:ring-1 hover:ring-primary/50 hover:shadow-lg'
      )}
      onClick={handleCardClick}
    >
      {/* Character Image Header */}
      {character.image_url && (
        <div className="relative w-full h-40 overflow-hidden border-b border-border/50">
          <img 
            src={character.image_url} 
            alt={character.name} 
            className="w-full h-full object-cover scale-125 object-top transition-transform group-hover:scale-[1.35]"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-background/90 via-background/20 to-transparent" />
        </div>
      )}

      <div className="p-4 flex-1 flex flex-col">
        <div className="flex items-start justify-between mb-3">
          <div className="flex items-center gap-3">
            {!character.image_url && (
              <div className="text-3xl flex items-center justify-center w-10 h-10 text-primary shrink-0 overflow-hidden rounded-full border border-primary/20 bg-primary/5">
                <User className="w-6 h-6" />
              </div>
            )}
            <div>
              <div className="flex items-center gap-2">
                <h3 className="font-display text-lg text-primary glow-text leading-tight">{character.name}</h3>
                {isActive && (
                  <Badge theme="primary-light-outline" size="sm">
                    <Star className="w-2.5 h-2.5 mr-0.5 fill-current" />
                    Active
                  </Badge>
                )}
              </div>
              <p className="text-xs text-muted-foreground mt-0.5">
                Level {character.level} {character.race} {character.class || 'Unknown'}
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            {onSetActive && (
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  onSetActive();
                }}
                className={cn(
                  "p-1 rounded-full transition-colors hover:bg-primary/10",
                  isActive ? "text-primary" : "text-muted-foreground hover:text-primary"
                )}
                title={isActive ? "Active Character" : "Set as Active Character"}
              >
                <Star className={cn("w-4 h-4", isActive && "fill-current")} />
              </button>
            )}
            <span className="flex items-center justify-center w-6 h-6 text-primary">
              <ClassIcon className="w-5 h-5" />
            </span>
            <ChevronRight className="w-4 h-4 text-muted-foreground group-hover:text-primary transition-colors" />
          </div>
        </div>

        <div className="flex items-center justify-between mt-auto pt-3 border-t border-border/50">
          <div className="flex items-center gap-3">
            <span className="flex items-center gap-1 text-muted-foreground text-xs">
              <BookOpen className="w-3.5 h-3.5" />
              <span>{character.spellbook.length} spells</span>
            </span>
          </div>

          <div className="flex gap-2" onClick={e => e.stopPropagation()}>
            <Button
              variant="ghost"
              size="icon"
              onClick={onDelete}
              className="h-8 w-8 text-destructive hover:text-destructive hover:bg-destructive/10"
              title="Delete character"
            >
              <Trash2 className="w-3.5 h-3.5" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}