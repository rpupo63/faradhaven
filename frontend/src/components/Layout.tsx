import { ReactNode, useRef, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { cn } from '@/lib/utils';
import { LogOut, User, Menu, Newspaper, Store, ScrollText } from 'lucide-react';
import { RaIcon } from '@/components/ui/RaIcon';
import { useAuth } from '@/context/AuthContext';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Sheet, SheetContent, SheetTrigger, SheetHeader, SheetTitle } from '@/components/ui/sheet';
import { DiceCustomizer } from '@/components/DiceCustomizer';
import { DiceAnimation } from '@/components/DiceAnimation';
import { MainContentBoundsProvider } from '@/context/MainContentBoundsContext';

interface LayoutProps {
  children: ReactNode;
}

const GM_EMAIL = 'rpupo63@gmail.com';

type NavItem = { path: string; label: string; renderIcon: () => React.ReactNode };

type NavSection = { title?: string; items: NavItem[] };

const baseNavItems: NavItem[] = [
  { path: '/characters', label: 'Characters', renderIcon: () => <RaIcon name="player" className="text-xl shrink-0" /> },
  { path: '/bulletin', label: 'Bulletin', renderIcon: () => <Newspaper className="w-5 h-5 shrink-0" /> },
  { path: '/arcana', label: 'Arcana', renderIcon: () => <RaIcon name="aura" className="text-xl shrink-0" /> },
  { path: '/shop', label: 'Shop', renderIcon: () => <Store className="w-5 h-5 shrink-0" /> },
  { path: '/game-rules', label: 'Game Rules', renderIcon: () => <RaIcon name="book" className="text-xl shrink-0" /> },
  { path: '/lore', label: 'Lore', renderIcon: () => <ScrollText className="w-5 h-5 shrink-0" /> },
  { path: '/dm-tools', label: 'DM Tools', renderIcon: () => <RaIcon name="anvil" className="text-xl shrink-0" /> },
];

const gmNavItem: NavItem = { path: '/gm/spells', label: 'GM Review', renderIcon: () => <RaIcon name="shield" className="text-xl shrink-0" /> };

export function Layout({ children }: LayoutProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const { user, logout, isAuthenticated } = useAuth();
  const [diceCustomizerOpen, setDiceCustomizerOpen] = useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const mainRef = useRef<HTMLElement>(null);

  const mainNavItems = user?.email === GM_EMAIL ? [...baseNavItems, gmNavItem] : baseNavItems;
  const navSections: NavSection[] = [{ items: mainNavItems }];

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <MainContentBoundsProvider mainRef={mainRef}>
    <div className="h-dvh min-h-0 w-full flex flex-col overflow-hidden">
      {/* Tome header – leather-bound cover title strip */}
      <header className="shrink-0 z-50 tome-leather-chrome tome-leather-header rounded-none">
        <div className="w-full py-4 px-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-6">
              {/* Sidebar Toggle / Mobile Hamburger Menu */}
              <div className="flex items-center">
                {/* Mobile Menu */}
                <div className="md:hidden">
                  <Sheet>
                    <SheetTrigger asChild>
                      <Button variant="ghost" size="icon" className="text-amber-200/85 hover:text-amber-50 hover:bg-black/20">
                        <Menu className="w-6 h-6" />
                      </Button>
                    </SheetTrigger>
                    <SheetContent side="left" className="w-[85vw] max-w-[320px] tome-leather-chrome border-r-2 border-faded-gold/45 p-0 shadow-2xl">
                      <SheetHeader className="p-6 border-b border-faded-gold/25">
                        <SheetTitle className="font-tome-heading text-xl tome-leather-foil-title text-left flex items-center gap-3">
                          <RaIcon name="fire" className="text-xl" />
                          Faradhaven
                        </SheetTitle>
                      </SheetHeader>
                      <nav className="flex flex-col p-4 gap-2">
                        {navSections.map((section, sIdx) => (
                          <div key={sIdx} className="flex flex-col gap-2">
                            {section.title && (
                              <p
                                className={cn(
                                  'px-4 text-xs font-tome-marginalia uppercase tracking-widest text-amber-200/50',
                                  sIdx > 0 && 'pt-4'
                                )}
                              >
                                {section.title}
                              </p>
                            )}
                            {section.items.map((item) => {
                              const isActive =
                                location.pathname === item.path ||
                                location.pathname.startsWith(`${item.path}/`);
                              return (
                                <Link
                                  key={item.path}
                                  to={item.path}
                                  className={cn(
                                    'tome-leather-nav-link flex items-center gap-3 px-4 py-3 rounded-lg transition-colors border border-transparent font-tome-subheading text-sm uppercase tracking-wide',
                                    isActive && 'tome-leather-nav-link-active'
                                  )}
                                >
                                  {item.renderIcon()}
                                  <span>{item.label}</span>
                                </Link>
                              );
                            })}
                          </div>
                        ))}
                      </nav>
                    </SheetContent>
                  </Sheet>
                </div>

                {/* Desktop Toggle Button */}
                <div className="hidden md:block">
                  <Button 
                    variant="ghost" 
                    size="icon" 
                    className="text-amber-200/85 hover:text-amber-50 hover:bg-black/20"
                    onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
                  >
                    <Menu className="w-6 h-6" />
                  </Button>
                </div>
              </div>

              <Link to="/" className="flex items-center gap-3 group shrink-0">
                <div>
                  <h1 className="font-tome-heading text-xl tome-leather-foil-title tracking-wide leading-tight">Faradhaven</h1>
                  <p className="text-xs text-amber-200/65 font-tome-marginalia">Steampunk RPG System</p>
                </div>
              </Link>
            </div>

            <div className="flex items-center gap-2">
              {/* User menu / Sign in */}
              {isAuthenticated ? (
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="ghost"
                      className="flex items-center gap-2 px-2 py-1 h-auto text-amber-100/90 hover:text-amber-50 hover:bg-black/20"
                    >
                      <div className="p-1.5 rounded-full bg-black/25 border border-faded-gold/45 shadow-[inset_0_1px_0_rgba(255,255,255,0.08)]">
                        <User className="w-4 h-4 text-faded-gold" />
                      </div>
                      <span className="hidden sm:inline text-sm font-tome-marginalia text-amber-100/88 max-w-[120px] truncate">
                        {user?.name || 'User'}
                      </span>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="w-48">
                    <div className="px-2 py-1.5">
                      <p className="text-sm font-medium">{user?.name}</p>
                      <p className="text-xs text-muted-foreground truncate">{user?.email}</p>
                    </div>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem onClick={() => setDiceCustomizerOpen(true)} className="cursor-pointer">
                      <RaIcon name="perspective-dice-five" className="text-sm mr-2" />
                      Dice Appearance
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem onClick={handleLogout} className="text-destructive cursor-pointer">
                      <LogOut className="w-4 h-4 mr-2" />
                      Sign out
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              ) : (
                <Button 
                  onClick={() => navigate('/login')} 
                  variant="outline" 
                  size="sm" 
                  className="border-faded-gold/55 bg-black/15 text-faded-gold hover:bg-black/28 hover:text-amber-50 font-tome-marginalia uppercase tracking-wider text-xs shadow-[inset_0_1px_0_rgba(255,255,255,0.06)]"
                >
                  Sign In
                </Button>
              )}
            </div>
          </div>
        </div>
      </header>

      <div className="flex min-h-0 flex-1 overflow-hidden">
        {/* Desktop Sidebar Navigation */}
        <aside 
          className={cn(
            "hidden md:flex md:flex-col shrink-0 min-h-0 tome-leather-chrome tome-leather-sidebar rounded-none transition-all duration-300 ease-in-out",
            sidebarCollapsed ? "w-16" : "w-60"
          )}
        >
          <div className="flex min-h-0 flex-1 flex-col overflow-y-auto overflow-x-hidden p-4">
            <nav className="flex flex-col gap-2">
              {navSections.map((section, sIdx) => (
                <div key={sIdx} className="flex flex-col gap-2">
                  {section.title && !sidebarCollapsed && (
                    <p
                      className={cn(
                        'px-3 text-[10px] font-tome-marginalia uppercase tracking-widest text-amber-200/45',
                        sIdx > 0 && 'pt-2'
                      )}
                    >
                      {section.title}
                    </p>
                  )}
                  {section.title && sidebarCollapsed && sIdx > 0 && (
                    <div className="mx-1 my-1 h-px shrink-0 bg-faded-gold/20" aria-hidden />
                  )}
                  {section.items.map((item) => {
                    const isActive =
                      location.pathname === item.path ||
                      location.pathname.startsWith(`${item.path}/`);
                    return (
                      <Link
                        key={item.path}
                        to={item.path}
                        className={cn(
                          'tome-leather-nav-link flex items-center gap-3 py-3 rounded-lg transition-all text-sm font-tome-subheading uppercase tracking-wide border border-transparent',
                          isActive && 'tome-leather-nav-link-active',
                          sidebarCollapsed ? 'justify-center px-0' : 'px-4'
                        )}
                        title={sidebarCollapsed ? item.label : undefined}
                      >
                        {item.renderIcon()}
                        {!sidebarCollapsed && (
                          <span className="truncate transition-opacity duration-300">{item.label}</span>
                        )}
                      </Link>
                    );
                  })}
                </div>
              ))}
            </nav>
          </div>
        </aside>

        {/* Main content – only this column scrolls; header/sidebar/footer stay fixed */}
        <main
          ref={mainRef}
          className="tome-inner-folio min-h-0 min-w-0 flex-1 w-full touch-pan-y overflow-y-auto overscroll-y-contain px-6 pt-4 pb-10 sm:px-10 sm:pt-6 sm:pb-16 md:px-12 md:pt-8 md:pb-12 lg:px-16 lg:pt-10 lg:pb-16"
        >
          <div className="relative w-full max-w-6xl mx-auto min-w-0">
            {children}
          </div>
        </main>
      </div>

      {/* Footer – leather tooling strip */}
      <footer className="shrink-0 tome-leather-chrome tome-leather-footer py-4">
        <div className="container mx-auto px-4 text-center">
          <p className="text-sm text-amber-200/75 font-tome-marginalia tracking-wide">
            A Pupo Production
          </p>
        </div>
      </footer>

      {isAuthenticated && (
        <DiceCustomizer open={diceCustomizerOpen} onOpenChange={setDiceCustomizerOpen} />
      )}

      <DiceAnimation />
    </div>
    </MainContentBoundsProvider>
  );
}