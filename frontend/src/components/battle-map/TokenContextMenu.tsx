import React, { useState, useEffect, useRef } from 'react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { MapToken } from '@/types/map';
import { EntityCombobox, TokenEntity } from './EntityCombobox';
import { useNavigate } from 'react-router-dom';
import { useToast } from '@/hooks/use-toast';

interface TokenContextMenuProps {
  token: MapToken;
  x: number;
  y: number;
  isOpen: boolean;
  onClose: () => void;    // Hides the dropdown only (keeps component mounted)
  onDismiss: () => void;  // Fully unmounts the component
  onChangeAffiliation: (tokenId: string, updates: Partial<MapToken>) => void;
  onDeleteToken: (tokenId: string) => void;
  userId: string;
  authToken: string;
}

export const TokenContextMenu: React.FC<TokenContextMenuProps> = ({
  token,
  x,
  y,
  isOpen,
  onClose,
  onDismiss,
  onChangeAffiliation,
  onDeleteToken,
  userId,
  authToken,
}) => {
  const navigate = useNavigate();
  const { toast } = useToast();
  const [affiliationDialogOpen, setAffiliationDialogOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  // Close menu on outside click
  useEffect(() => {
    if (!isOpen) return;
    const handleClick = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        onDismiss();
      }
    };
    document.addEventListener('mousedown', handleClick);
    return () => document.removeEventListener('mousedown', handleClick);
  }, [isOpen, onDismiss]);

  const handleGoToSheet = () => {
    onDismiss();
    if (token.token_type === 'pc' && token.character_id) {
      navigate(`/character/${token.character_id}/sheet`);
    } else if (token.token_type === 'npc' && token.monster_id) {
      navigate(`/monster/${token.monster_id}`);
    } else {
      toast({
        title: 'No Sheet Found',
        description: 'This token is not linked to a character or monster sheet.',
        variant: 'destructive',
      });
    }
  };

  const handleOpenAffiliationDialog = () => {
    onClose(); // Hide the dropdown, but keep component mounted
    setAffiliationDialogOpen(true);
  };

  const handleDialogClose = (open: boolean) => {
    setAffiliationDialogOpen(open);
    if (!open) {
      onDismiss(); // Fully unmount when dialog closes
    }
  };

  const handleSelectEntity = (entity: TokenEntity) => {
    const updates: Partial<MapToken> = {
      name: entity.name,
      image_url: entity.image_url ?? '',
      token_type: entity.entity_type === 'character' ? 'pc' : 'npc',
    };

    if (entity.entity_type === 'character') {
      (updates as Record<string, unknown>).character_id = entity.id;
      (updates as Record<string, unknown>).monster_id = null;
      if (entity.owner_user_id) {
        (updates as Record<string, unknown>).assigned_user_id = entity.owner_user_id;
      }
    } else {
      (updates as Record<string, unknown>).monster_id = entity.id;
      (updates as Record<string, unknown>).character_id = null;
    }

    onChangeAffiliation(token.id, updates);
    setAffiliationDialogOpen(false);
    onDismiss();
  };

  const handleDelete = () => {
    if (window.confirm(`Are you sure you want to remove ${token.name} from the map?`)) {
      onDeleteToken(token.id);
      onDismiss();
    }
  };

  const hasSheet =
    (token.token_type === 'pc' && token.character_id) ||
    (token.token_type === 'npc' && token.monster_id);

  if (!isOpen && !affiliationDialogOpen) return null;

  return (
    <>
      {isOpen && (
        <div
          ref={menuRef}
          className="w-52 rounded-md border bg-popover p-1 text-popover-foreground shadow-md"
          style={{
            position: 'fixed',
            top: y,
            left: x,
            zIndex: 1000,
          }}
        >
          <div className="px-2 py-1.5 text-sm font-semibold">{token.name}</div>
          <div className="-mx-1 my-1 h-px bg-muted" />

          {hasSheet && (
            <button
              className="relative flex w-full cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground"
              onClick={handleGoToSheet}
            >
              Go to {token.token_type === 'pc' ? 'Character' : 'Monster'} Sheet
            </button>
          )}

          <button
            className="relative flex w-full cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground"
            onClick={handleOpenAffiliationDialog}
          >
            Change Affiliation
          </button>

          <div className="-mx-1 my-1 h-px bg-muted" />
          <button
            className="relative flex w-full cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-destructive hover:text-destructive-foreground"
            onClick={handleDelete}
          >
            Remove from Map
          </button>
        </div>
      )}

      <Dialog open={affiliationDialogOpen} onOpenChange={handleDialogClose}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Change Affiliation — {token.name}</DialogTitle>
          </DialogHeader>
          <div className="py-4">
            <EntityCombobox
              value={null}
              onSelect={handleSelectEntity}
              userId={userId}
              authToken={authToken}
            />
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
};
