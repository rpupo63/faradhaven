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
  { path: '/dm-tools', label: 'DM Tools', renderIcon: () => <RaIcon name="anvil" className="text-xl shrink-0" /> },
];

const gmNavItem: NavItem = { path: '/gm/spells', label: 'GM Review', renderIcon: () => <RaIcon name="shield" className="text-xl shrink-0" /> };

const loreNavItems: NavItem[] = [
  { path: '/lore', label: 'Lore', renderIcon: () => <ScrollText className="w-5 h-5 shrink-0" /> },
];

export function Layout({ children }: LayoutProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const { user, logout, isAuthenticated } = useAuth();
  const [diceCustomizerOpen, setDiceCustomizerOpen] = useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const mainRef = useRef<HTMLElement>(null);

  const mainNavItems = user?.email === GM_EMAIL ? [...baseNavItems, gmNavItem] : baseNavItems;
  const navSections: NavSection[] = [
    { items: mainNavItems },
    { title: 'Lore', items: loreNavItems },
  ];

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <MainContentBoundsProvider mainRef={mainRef}>
    <div className="h-dvh min-h-0 w-full flex flex-col overflow-hidden">
      {/* Tome header – bar like a handbook title strip */}
      <header className="shrink-0 border-b-2 border-faded-gold/50 bg-card/80 backdrop-blur-sm z-50 hand-drawn-border border-t-0 border-l-0 border-r-0 rounded-none">
        <div className="w-full py-4 px-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-6">
              {/* Sidebar Toggle / Mobile Hamburger Menu */}
              <div className="flex items-center">
                {/* Mobile Menu */}
                <div className="md:hidden">
                  <Sheet>
                    <SheetTrigger asChild>
                      <Button variant="ghost" size="icon" className="text-muted-foreground hover:text-primary">
                        <Menu className="w-6 h-6" />
                      </Button>
                    </SheetTrigger>
                    <SheetContent side="left" className="w-[85vw] max-w-[320px] bg-card/95 backdrop-blur-sm border-r-2 border-faded-gold/50 p-0">
                      <SheetHeader className="p-6 border-b border-faded-gold/20">
                        <SheetTitle className="font-tome-heading text-xl text-primary text-left flex items-center gap-3">
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
                                  'px-4 text-xs font-tome-marginalia uppercase tracking-widest text-muted-foreground',
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
                                    'flex items-center gap-3 px-4 py-3 rounded-lg transition-colors',
                                    isActive
                                      ? 'bg-primary/10 text-primary border border-primary/20'
                                      : 'text-muted-foreground hover:text-foreground hover:bg-primary/5'
                                  )}
                                >
                                  {item.renderIcon()}
                                  <span className="font-tome-subheading text-sm uppercase tracking-wide">{item.label}</span>
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
                    className="text-muted-foreground hover:text-primary"
                    onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
                  >
                    <Menu className="w-6 h-6" />
                  </Button>
                </div>
              </div>

              <Link to="/" className="flex items-center gap-3 group shrink-0">
                <div>
                  <h1 className="font-tome-heading text-xl text-primary tracking-wide leading-tight">Faradhaven</h1>
                  <p className="text-xs text-muted-foreground font-tome-marginalia">Steampunk RPG System</p>
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
                      className="flex items-center gap-2 px-2 py-1 h-auto"
                    >
                      <div className="p-1.5 rounded-full bg-primary/20 border border-faded-gold/40">
                        <User className="w-4 h-4 text-primary" />
                      </div>
                      <span className="hidden sm:inline text-sm font-tome-marginalia text-foreground/80 max-w-[120px] truncate">
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
                  className="border-faded-gold/40 hover:bg-primary/10 font-tome-marginalia uppercase tracking-wider text-xs"
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
            "hidden md:flex md:flex-col shrink-0 min-h-0 border-r-2 border-faded-gold/50 bg-card/40 backdrop-blur-sm hand-drawn-border border-t-0 border-l-0 border-b-0 rounded-none transition-all duration-300 ease-in-out",
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
                        'px-3 text-[10px] font-tome-marginalia uppercase tracking-widest text-muted-foreground/80',
                        sIdx > 0 && 'pt-2'
                      )}
                    >
                      {section.title}
                    </p>
                  )}
                  {section.title && sidebarCollapsed && sIdx > 0 && (
                    <div className="mx-1 my-1 h-px shrink-0 bg-faded-gold/25" aria-hidden />
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
                          'flex items-center gap-3 px-3 py-3 rounded-lg transition-all text-sm font-tome-subheading uppercase tracking-wide border border-transparent',
                          isActive
                            ? 'text-primary bg-primary/15 border-primary/20 shadow-seal'
                            : 'text-muted-foreground hover:text-primary hover:bg-primary/10',
                          sidebarCollapsed ? 'justify-center px-0' : 'px-4'
                        )}
                        title={sidebarCollapsed ? item.label : undefined}
                      >
                        {item.renderIcon()}
                        {!sidebarCollapsed && (
                          <span className="truncate opacity-100 transition-opacity duration-300">{item.label}</span>
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
          className="min-h-0 min-w-0 flex-1 w-full touch-pan-y overflow-y-auto overscroll-y-contain px-6 pt-4 pb-10 sm:px-10 sm:pt-6 sm:pb-16 md:px-12 md:pt-8 md:pb-12 lg:px-16 lg:pt-10 lg:pb-16"
        >
          <div className="relative w-full max-w-6xl mx-auto min-w-0">
            {children}
          </div>
        </main>
      </div>

      {/* Footer – ledger-style line */}
      <footer className="shrink-0 border-t-2 border-faded-gold/40 py-4">
        <div className="container mx-auto px-4 text-center">
          <p className="text-sm text-muted-foreground font-tome-marginalia">
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