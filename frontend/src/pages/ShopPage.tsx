import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  getItems,
  getWeapons,
  getCharacterSheet,
  purchaseItem,
  getStoreOwners,
} from '@/lib/api';
import { ApiItem, ApiWeapon, ApiStoreOwner } from '@/types/game';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { LoadingQuill } from '@/components/LoadingQuill';
import { RefreshIndicator } from '@/components/ui/refresh-indicator';
import { LoadingButton } from '@/components/ui/loading-button';
import {
  Package,
  Scale,
  Coins,
  Zap,
  Search,
  SortAsc,
  ShoppingCart,
  AlertCircle,
  Sword,
  Target,
  Store,
  MapPin,
  Users,
} from 'lucide-react';
import { Link } from 'react-router-dom';
import { cn } from '@/lib/utils';
import { useState, useMemo, useEffect } from 'react';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { useAuth } from '@/context/AuthContext';
import { toast } from 'sonner';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { ApiWeaponDamage } from '@/types/game';

const PAGINATION_LIMIT = 12; // Number of items to display per page

function RarityBadge({ rarity }: { rarity: string }) {
  const normalized = rarity.toLowerCase();
  let colorClass =
    'bg-muted/50 text-muted-foreground border-muted-foreground/30';

  if (normalized === 'common')
    colorClass = 'bg-slate-100/50 text-slate-700 border-slate-200';
  else if (normalized === 'uncommon')
    colorClass =
      'bg-element-nature/10 text-element-nature border-element-nature/30';
  else if (normalized === 'rare')
    colorClass = 'bg-element-ice/10 text-element-ice border-element-ice/30';
  else if (normalized === 'very rare')
    colorClass =
      'bg-element-dark/10 text-element-dark border-element-dark/30';
  else if (normalized === 'legendary')
    colorClass =
      'bg-element-fire/10 text-element-fire border-element-fire/30 shadow-glow';

  return (
    <Badge
      variant="outline"
      className={cn(
        'font-tome-marginalia text-[10px] uppercase tracking-wider',
        colorClass
      )}
    >
      {rarity}
    </Badge>
  );
}

function CurrencyDisplay({ cp }: { cp: number }) {
  const gp = Math.floor(cp / 100);
  const sp = Math.floor((cp % 100) / 10);
  const remainingCp = cp % 10;

  return (
    <div className="flex items-center gap-3 font-tome-marginalia text-sm">
      {gp > 0 && (
        <div className="flex items-center gap-1">
          <span className="text-faded-gold font-bold">{gp}</span>
          <span className="text-muted-foreground text-[10px] uppercase">
            gp
          </span>
        </div>
      )}
      {sp > 0 && (
        <div className="flex items-center gap-1">
          <span className="text-slate-400 font-bold">{sp}</span>
          <span className="text-muted-foreground text-[10px] uppercase">
            sp
          </span>
        </div>
      )}
      <div className="flex items-center gap-1">
        <span className="text-orange-700 font-bold">{remainingCp}</span>
        <span className="text-muted-foreground text-[10px] uppercase">cp</span>
      </div>
    </div>
  );
}

function DamageDisplay({ damage }: { damage: ApiWeaponDamage }) {
  return (
    <div className="flex items-center gap-1.5 text-xs">
      <Badge variant="secondary" font="mono" size="sm">
        {damage.damage_dice}
      </Badge>
      <span className="text-muted-foreground">{damage.damage_type}</span>
      {damage.damage_category !== 'Base' && (
        <span className="text-[10px] text-muted-foreground/70 italic">
          ({damage.damage_category})
        </span>
      )}
    </div>
  );
}

type ShopItem = (ApiItem & { type: 'item' }) | (ApiWeapon & { type: 'weapon' });

function storeOwnerSellsItem(owner: ApiStoreOwner, shopItem: ShopItem): boolean {
  const rules = owner.catalog_rules ?? [];
  const isWeapon = shopItem.type === 'weapon';
  const id = shopItem.id;
  const category = shopItem.category;
  const rarity = shopItem.rarity;
  for (const rule of rules) {
    if (rule.item_id && rule.item_id === id && !isWeapon) return true;
    if (rule.weapon_id && rule.weapon_id === id && isWeapon) return true;
    if (
      rule.category != null &&
      rule.category !== '' &&
      category === rule.category
    ) {
      const allowed = rule.allowed_rarities ?? [];
      if (allowed.length === 0) return true;
      if (allowed.includes(rarity)) return true;
    }
  }
  return false;
}

function ShopItemCard({
  item,
  onBuy,
  isBuying,
  canAfford,
  vendorPriceCp,
  listPriceCp,
}: {
  item: ShopItem;
  onBuy?: () => void;
  isBuying?: boolean;
  canAfford?: boolean;
  /** When set (vendor selected), show this cp price and optional list comparison */
  vendorPriceCp?: number | null;
  listPriceCp?: number | null;
}) {
  const isWeapon = item.type === 'weapon';

  return (
    <Card className="arcane-border h-full hover:bg-primary/5 transition-all hover:shadow-lg hover:shadow-primary/10 group overflow-hidden flex flex-col">
      <CardHeader className="pb-3 pt-4 space-y-1">
        <div className="flex items-start justify-between gap-2">
          <h3 className="font-tome-heading text-lg text-primary group-hover:glow-text transition-colors leading-tight">
            {item.name}
          </h3>
          <RarityBadge rarity={item.rarity} />
        </div>
        <div className="flex items-center gap-2 text-xs text-muted-foreground font-tome-marginalia">
          <span>{item.category}</span>
          {isWeapon && (
            <>
              <span>•</span>
              <span>{(item as ApiWeapon).range_type}</span>
            </>
          )}
          {!isWeapon && (item as ApiItem).is_consumable && (
            <>
              <span>•</span>
              <span className="text-element-nature">Consumable</span>
            </>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-4 flex-grow flex flex-col">
        {/* Stats Grid */}
        <div className="grid grid-cols-2 gap-2 text-xs font-tome-marginalia bg-muted/20 p-2 rounded-lg border border-border/50">
          <div className="flex flex-col gap-1 text-muted-foreground min-w-0">
            <div className="flex items-start gap-1.5">
              <Coins className="w-3 h-3 shrink-0 mt-0.5" />
              {vendorPriceCp != null ? (
                <div className="flex flex-col gap-0.5 min-w-0">
                  <CurrencyDisplay cp={vendorPriceCp} />
                  {listPriceCp != null && listPriceCp !== vendorPriceCp && (
                    <div className="text-[10px] text-muted-foreground/80 flex flex-wrap items-center gap-x-2 gap-y-0.5">
                      <span>List:</span>
                      <CurrencyDisplay cp={listPriceCp} />
                    </div>
                  )}
                </div>
              ) : (
                <span>{item.cost}</span>
              )}
            </div>
          </div>
          <div className="flex items-center gap-1.5 text-muted-foreground">
            <Scale className="w-3 h-3" />
            <span>{item.weight}</span>
          </div>
          {isWeapon && (
            <div className="flex items-center gap-1.5 text-muted-foreground col-span-2">
              <Target className="w-3 h-3" />
              <span>
                Range {(item as ApiWeapon).range_normal} ft.
                {(item as ApiWeapon).range_long > 0 && (
                  <span className="opacity-70">
                    {' '}
                    / {(item as ApiWeapon).range_long} ft.
                  </span>
                )}
              </span>
            </div>
          )}
        </div>

        {/* Damage */}
        {isWeapon && (item as ApiWeapon).damages && (
          <div className="space-y-1.5">
            <div className="text-xs font-tome-subheading text-muted-foreground uppercase tracking-wider flex items-center gap-1.5">
              <Sword className="w-3 h-3" />
              Damage
            </div>
            <div className="space-y-1">
              {(item as ApiWeapon).damages?.map(d => (
                <DamageDisplay key={d.id} damage={d} />
              ))}
              {(item as ApiWeapon).versatile_damage_dice && (
                <div className="text-xs text-muted-foreground italic pl-1">
                  Versatile: {(item as ApiWeapon).versatile_damage_dice}
                </div>
              )}
            </div>
          </div>
        )}

        {/* Properties */}
        {isWeapon &&
          (item as ApiWeapon).properties &&
          (item as ApiWeapon).properties.length > 0 && (
            <div className="flex flex-wrap gap-1">
              {(item as ApiWeapon).properties.map((prop, idx) => (
                <Badge
                  key={idx}
                  variant="outline"
                  className="text-[10px] font-tome-marginalia border-muted-foreground/20 text-muted-foreground"
                >
                  {prop}
                </Badge>
              ))}
            </div>
          )}

        {/* Description & Effects */}
        <div className="text-sm space-y-2 flex-grow">
          <p className="text-muted-foreground leading-relaxed italic">
            {item.description}
          </p>
          {(isWeapon
            ? (item as ApiWeapon).secondary_effect
            : (item as ApiItem).effects) && (
            <div className="bg-primary/5 p-2 rounded border border-primary/10 text-xs mt-2">
              <span className="font-semibold text-primary flex items-center gap-1 mb-1">
                <Zap className="w-3 h-3" />{' '}
                {isWeapon ? 'Special Effect' : 'Effects'}
              </span>
              <span className="text-muted-foreground">
                {isWeapon
                  ? (item as ApiWeapon).secondary_effect
                  : (item as ApiItem).effects}
              </span>
            </div>
          )}
        </div>

        {/* Buy Button */}
        {onBuy && (
          <LoadingButton
            onClick={e => {
              e.stopPropagation();
              onBuy();
            }}
            isLoading={isBuying}
            disabled={!canAfford}
            variant="quill"
            size="sm"
            className="w-full mt-auto"
            loadingText="Buying..."
          >
            <ShoppingCart className="w-4 h-4 mr-2" />
            {canAfford ? 'Buy' : 'Insufficient Funds'}
          </LoadingButton>
        )}
      </CardContent>
    </Card>
  );
}

export default function ShopPage() {
  const queryClient = useQueryClient();
  const { token, activeCharacterId } = useAuth();

  const [search, setSearch] = useState('');
  const [itemType, setItemType] = useState('all');
  const [category, setCategory] = useState('all');
  const [rarity, setRarity] = useState('all');
  const [sortBy, setSortBy] = useState('name-asc');
  const [itemToBuy, setItemToBuy] = useState<ShopItem | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedVendorId, setSelectedVendorId] = useState<string | null>(null);

  const {
    data: items,
    isLoading: isLoadingItems,
    isFetching: isFetchingItems,
  } = useQuery({
    queryKey: ['items'],
    queryFn: () => getItems(token ?? undefined),
    staleTime: 60_000 * 5,
  });

  const {
    data: weapons,
    isLoading: isLoadingWeapons,
    isFetching: isFetchingWeapons,
  } = useQuery({
    queryKey: ['weapons'],
    queryFn: () => getWeapons(token ?? undefined),
    staleTime: 60_000 * 5,
  });

  const { data: sheet } = useQuery({
    queryKey: ['character-sheet', activeCharacterId],
    queryFn: () => getCharacterSheet(activeCharacterId!, token ?? undefined),
    enabled: !!activeCharacterId && !!token,
  });

  const {
    data: storeOwners,
    isLoading: isLoadingStoreOwners,
    isFetching: isFetchingStoreOwners,
  } = useQuery({
    queryKey: ['store-owners'],
    queryFn: () => getStoreOwners(token ?? undefined),
    staleTime: 60_000 * 5,
  });

  const selectedVendor = useMemo((): ApiStoreOwner | null => {
    if (!selectedVendorId || !storeOwners?.length) return null;
    return storeOwners.find(o => o.id === selectedVendorId) ?? null;
  }, [selectedVendorId, storeOwners]);

  const purchaseMutation = useMutation({
    mutationFn: ({
      item,
      storeOwnerId,
    }: {
      item: ShopItem;
      storeOwnerId: string | null;
    }) =>
      purchaseItem(
        activeCharacterId!,
        item.id,
        item.type,
        token ?? undefined,
        storeOwnerId
      ),
    onSuccess: data => {
      toast.success(data.message);
      queryClient.invalidateQueries({
        queryKey: ['character-sheet', activeCharacterId],
      });
      setItemToBuy(null);
    },
    onError: (error: unknown) => {
      toast.error((error as Error).message || 'Failed to purchase item');
      setItemToBuy(null);
    },
  });

  const combinedItems = useMemo((): ShopItem[] => {
    const allItems: ShopItem[] = [];
    if (items) {
      allItems.push(...items.map(i => ({ ...i, type: 'item' as const })));
    }
    if (weapons) {
      allItems.push(...weapons.map(w => ({ ...w, type: 'weapon' as const })));
    }
    return allItems;
  }, [items, weapons]);

  const categories = useMemo(() => {
    if (!combinedItems) return [];
    let source = combinedItems;

    const armorCategories = ["Light Armor", "Medium Armor", "Heavy Armor"];

    if (itemType === 'weapon') {
        source = combinedItems.filter(i => i.type === 'weapon');
    } else if (itemType === 'armor') {
        source = combinedItems.filter(i => i.type === 'item' && armorCategories.includes(i.category));
    } else if (itemType === 'item') {
        source = combinedItems.filter(i => i.type === 'item' && !armorCategories.includes(i.category));
    }
    
    const cats = new Set(source.map(i => i.category));
    return Array.from(cats).sort();
  }, [combinedItems, itemType]);

  const rarities = useMemo(() => {
    if (!combinedItems) return [];
    const rs = new Set(combinedItems.map(i => i.rarity));
    return Array.from(rs).sort();
  }, [combinedItems]);

  const parseCost = (costStr: string): number => {
    const match = costStr.match(/(\d+)\s*(gp|sp|cp|pp)/i);
    if (!match) return 0;
    const value = parseInt(match[1], 10);
    const unit = match[2].toLowerCase();
    switch (unit) {
      case 'pp':
        return value * 1000;
      case 'gp':
        return value * 100;
      case 'sp':
        return value * 10;
      case 'cp':
        return value * 1;
      default:
        return value;
    }
  };

  const filteredAndSortedItems = useMemo(() => {
    if (!combinedItems) return [];
    const result = combinedItems.filter(item => {
      const isArmor = item.type === 'item' && item.category.includes('Armor');

      const matchesSearch =
        item.name.toLowerCase().includes(search.toLowerCase()) ||
        item.description.toLowerCase().includes(search.toLowerCase());
      
      let matchesItemType = false;
      if (itemType === 'all') {
          matchesItemType = true;
      } else if (itemType === 'weapon') {
          matchesItemType = item.type === 'weapon';
      } else if (itemType === 'armor') {
          matchesItemType = isArmor;
      } else if (itemType === 'item') {
          matchesItemType = item.type === 'item' && !isArmor;
      }

      const matchesCategory = category === 'all' || item.category === category;
      const matchesRarity = rarity === 'all' || item.rarity === rarity;
      const matchesVendor =
        !selectedVendor || storeOwnerSellsItem(selectedVendor, item);
      return (
        matchesSearch &&
        matchesItemType &&
        matchesCategory &&
        matchesRarity &&
        matchesVendor
      );
    });

    const priceCp = (item: ShopItem) => {
      const base = parseCost(item.cost);
      if (selectedVendor) {
        return Math.round(base * selectedVendor.exchange_rate);
      }
      return base;
    };

    result.sort((a, b) => {
      if (sortBy === 'name-asc') return a.name.localeCompare(b.name);
      if (sortBy === 'name-desc') return b.name.localeCompare(a.name);
      if (sortBy === 'price-asc') return priceCp(a) - priceCp(b);
      if (sortBy === 'price-desc') return priceCp(b) - priceCp(a);
      return 0;
    });

    return result;
  }, [
    combinedItems,
    search,
    itemType,
    category,
    rarity,
    sortBy,
    selectedVendor,
  ]);
  
  const totalPages = useMemo(() => {
    return Math.ceil(filteredAndSortedItems.length / PAGINATION_LIMIT);
  }, [filteredAndSortedItems]);

  const paginatedShopItems = useMemo(() => {
    const startIndex = (currentPage - 1) * PAGINATION_LIMIT;
    const endIndex = startIndex + PAGINATION_LIMIT;
    return filteredAndSortedItems.slice(startIndex, endIndex);
  }, [filteredAndSortedItems, currentPage]);

  useEffect(() => {
    setTimeout(() => setCurrentPage(1), 0);
  }, [search, itemType, category, rarity, sortBy, selectedVendorId]);

  const isLoading =
    isLoadingItems || isLoadingWeapons || isLoadingStoreOwners;
  const isFetching =
    isFetchingItems || isFetchingWeapons || isFetchingStoreOwners;

  const handleBuyClick = (item: ShopItem) => {
    const storeOwnerId = selectedVendorId;
    if (item.type === 'item') {
      purchaseMutation.mutate({ item, storeOwnerId });
    } else {
      setItemToBuy(item);
    }
  };

  return (
    <div className="w-full space-y-8 p-6 relative">
      <RefreshIndicator isFetching={isFetching && !isLoading} />

      {/* Header */}
      <div className="flex flex-col xl:flex-row xl:items-center justify-between gap-4">
        <div className="flex items-center gap-4">
          <div className="p-3 rounded-full border-2 border-faded-gold/50 bg-primary/10">
            <Store className="w-6 h-6 text-primary" />
          </div>
          <div>
            <h1 className="font-tome-heading text-3xl text-primary glow-text">
              Faradhaven Market
            </h1>
            <div className="flex items-center gap-4 mt-1">
              <p className="text-muted-foreground text-sm font-tome-marginalia">
                Browse goods, gear, and weaponry from across the realm.
              </p>
              {sheet && (
                <>
                  <span className="text-muted-foreground/30">•</span>
                  <div className="flex items-center gap-2">
                    <span className="text-xs text-muted-foreground uppercase tracking-widest font-tome-marginalia">
                      Wallet:
                    </span>
                    <CurrencyDisplay cp={sheet.money} />
                  </div>
                </>
              )}
            </div>
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-3">
          <div className="relative w-full md:w-64">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
            <Input
              placeholder="Search items..."
              value={search}
              onChange={e => setSearch(e.target.value)}
              className="pl-9"
            />
          </div>
           <Select value={itemType} onValueChange={(v) => { setItemType(v); setCategory('all'); }}>
            <SelectTrigger className="w-[140px]">
              <SelectValue placeholder="Item Type" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Types</SelectItem>
              <SelectItem value="weapon">Weapons</SelectItem>
              <SelectItem value="armor">Armor</SelectItem>
              <SelectItem value="item">Items</SelectItem>
            </SelectContent>
          </Select>
          <Select value={category} onValueChange={setCategory}>
            <SelectTrigger className="w-[140px]">
              <SelectValue placeholder="Category" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Categories</SelectItem>
              {categories.map(cat => (
                <SelectItem key={cat} value={cat}>
                  {cat}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Select value={rarity} onValueChange={setRarity}>
            <SelectTrigger className="w-[140px]">
              <SelectValue placeholder="Rarity" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Rarities</SelectItem>
              {rarities.map(r => (
                <SelectItem key={r} value={r}>
                  {r}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Select value={sortBy} onValueChange={setSortBy}>
            <SelectTrigger className="w-[160px]">
              <div className="flex items-center gap-2 text-muted-foreground">
                <SortAsc className="w-4 h-4" />
                <SelectValue placeholder="Sort by" />
              </div>
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="name-asc">Name (A-Z)</SelectItem>
              <SelectItem value="name-desc">Name (Z-A)</SelectItem>
              <SelectItem value="price-asc">Price (Low-High)</SelectItem>
              <SelectItem value="price-desc">Price (High-Low)</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Vendors */}
      {storeOwners && storeOwners.length > 0 && (
        <div className="space-y-4">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <h2 className="font-tome-heading text-lg text-primary flex items-center gap-2">
              <Users className="w-5 h-5" />
              Shopkeepers
            </h2>
            <button
              type="button"
              onClick={() => setSelectedVendorId(null)}
              className={cn(
                'text-sm font-tome-marginalia px-3 py-1.5 rounded-md border transition-colors',
                selectedVendorId === null
                  ? 'border-primary bg-primary/10 text-primary'
                  : 'border-border text-muted-foreground hover:border-primary/40'
              )}
            >
              All vendors
            </button>
          </div>
          <div className="flex gap-3 overflow-x-auto pb-2 scrollbar-thin">
            {storeOwners.map(owner => {
              const selected = selectedVendorId === owner.id;
              return (
                <button
                  key={owner.id}
                  type="button"
                  onClick={() => setSelectedVendorId(owner.id)}
                  className={cn(
                    'flex-shrink-0 w-[min(100%,280px)] text-left arcane-border rounded-xl p-4 transition-all',
                    selected
                      ? 'bg-primary/10 border-primary shadow-md shadow-primary/10'
                      : 'bg-background/50 hover:bg-primary/5'
                  )}
                >
                  <p className="font-tome-heading text-base text-primary leading-tight">
                    {owner.name}
                  </p>
                  <p className="text-xs text-muted-foreground font-tome-marginalia mt-1 flex items-start gap-1">
                    <MapPin className="w-3 h-3 shrink-0 mt-0.5" />
                    <span className="line-clamp-2">{owner.location}</span>
                  </p>
                  <div className="flex flex-wrap gap-1.5 mt-2">
                    <Badge
                      variant="outline"
                      className="text-[10px] font-tome-marginalia"
                    >
                      Buy ×{owner.exchange_rate.toFixed(2)}
                    </Badge>
                    <Badge
                      variant="outline"
                      className="text-[10px] font-tome-marginalia border-muted-foreground/30"
                    >
                      Sell ×{owner.willingness_to_purchase.toFixed(2)}
                    </Badge>
                  </div>
                </button>
              );
            })}
          </div>
          {selectedVendor && (
            <div className="arcane-border rounded-xl p-5 bg-muted/10 space-y-3">
              <div>
                <p className="text-xs uppercase tracking-widest text-muted-foreground font-tome-marginalia mb-1">
                  At the counter
                </p>
                <p className="text-sm text-foreground leading-relaxed font-tome-marginalia">
                  {selectedVendor.personality}
                </p>
              </div>
              <p className="text-xs text-muted-foreground font-tome-marginalia flex items-start gap-2">
                <MapPin className="w-3.5 h-3.5 shrink-0 mt-0.5" />
                {selectedVendor.location}
              </p>
              {(selectedVendor.categories_obtained?.length ?? 0) > 0 && (
                <div className="flex flex-wrap gap-1.5">
                  <span className="text-[10px] uppercase text-muted-foreground font-tome-marginalia w-full">
                    Seeks / sources
                  </span>
                  {selectedVendor.categories_obtained!.map(cat => (
                    <Badge
                      key={cat}
                      variant="secondary"
                      className="text-[10px] font-tome-marginalia"
                    >
                      {cat}
                    </Badge>
                  ))}
                </div>
              )}
              <p className="text-xs text-muted-foreground font-tome-marginalia border-t border-border/50 pt-3">
                Listed prices are adjusted by this vendor&apos;s exchange rate (
                {selectedVendor.exchange_rate.toFixed(2)}×). Purchases use the
                vendor price shown on each card.
              </p>
            </div>
          )}
        </div>
      )}

      {/* No Active Character Banner */}
      {!activeCharacterId && (
        <div className="arcane-border rounded-lg p-4 bg-muted/20 flex items-center gap-3">
          <AlertCircle className="w-5 h-5 text-muted-foreground flex-shrink-0" />
          <p className="text-sm text-muted-foreground">
            Select an active character on the{' '}
            <Link to="/" className="text-primary hover:underline font-medium">
              Your Characters
            </Link>{' '}
            page to make purchases.
          </p>
        </div>
      )}

      {/* Items Grid */}
      {isLoading ? (
        <LoadingQuill label="Loading wares from the market..." />
      ) : filteredAndSortedItems.length > 0 ? (
        <>
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {paginatedShopItems.map(item => {
              const listCp = parseCost(item.cost);
              const vendorCp = selectedVendor
                ? Math.round(listCp * selectedVendor.exchange_rate)
                : null;
              const affordCp = vendorCp ?? listCp;
              return (
                <ShopItemCard
                  key={item.id}
                  item={item}
                  onBuy={
                    activeCharacterId
                      ? () => handleBuyClick(item)
                      : undefined
                  }
                  isBuying={
                    purchaseMutation.variables?.item?.id === item.id &&
                    purchaseMutation.isPending
                  }
                  canAfford={!sheet || sheet.money >= affordCp}
                  vendorPriceCp={vendorCp}
                  listPriceCp={selectedVendor ? listCp : null}
                />
              );
            })}
          </div>

          {totalPages > 1 && (
            <div className="flex items-center justify-between mt-8">
              <button
                onClick={() => setCurrentPage((prev) => Math.max(1, prev - 1))}
                disabled={currentPage === 1}
                className="px-4 py-2 rounded-md bg-muted-foreground/20 text-muted-foreground hover:bg-muted-foreground/30 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Previous
              </button>
              <span className="font-tome-marginalia text-sm text-muted-foreground">
                Page {currentPage} of {totalPages}
              </span>
              <button
                onClick={() => setCurrentPage((prev) => Math.min(totalPages, prev + 1))}
                disabled={currentPage === totalPages}
                className="px-4 py-2 rounded-md bg-muted-foreground/20 text-muted-foreground hover:bg-muted-foreground/30 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Next
              </button>
            </div>
          )}
        </>
      ) : (
        <div className="arcane-border rounded-xl p-12 text-center">
          <Store className="w-16 h-16 mx-auto mb-6 text-muted-foreground" />
          <h2 className="font-tome-heading text-2xl text-primary mb-2">
            The Market is Empty
          </h2>
          <p className="text-muted-foreground font-tome-marginalia">
            No items or weapons match your search. Try a different filter.
          </p>
        </div>
      )}

       {/* Confirmation Dialog for Weapons */}
       <AlertDialog open={!!itemToBuy && itemToBuy.type === 'weapon'} onOpenChange={(open) => !open && setItemToBuy(null)}>
        <AlertDialogContent className="arcane-border">
          <AlertDialogHeader>
            <AlertDialogTitle className="font-tome-heading text-xl text-primary">Confirm Purchase</AlertDialogTitle>
            <AlertDialogDescription className="font-tome-marginalia">
              Are you sure you want to purchase the <span className="text-primary font-bold">{itemToBuy?.name}</span> for <span className="text-faded-gold font-bold">{itemToBuy?.cost}</span>? 
              This will add the weapon to your character's inventory and deduct the cost from your wallet.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel className="font-tome-marginalia">Cancel</AlertDialogCancel>
            <AlertDialogAction 
              className="bg-primary text-primary-foreground hover:bg-primary/90 font-tome-marginalia"
              onClick={() =>
                itemToBuy &&
                purchaseMutation.mutate({
                  item: itemToBuy,
                  storeOwnerId: selectedVendorId,
                })
              }
            >
              Confirm Purchase
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
