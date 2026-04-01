import { ReactNode, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { cn } from '@/lib/utils';
import { Users, Skull, BookMarked, Sparkles, LogOut, User, Menu, Atom, Flame, Zap, Newspaper, Hammer, Bone, Store, Shield, Dice5 } from 'lucide-react';
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

interface LayoutProps {
  children: ReactNode;
}

const GM_EMAIL = 'rpupo63@gmail.com';

const baseNavItems = [
  { path: '/characters', label: 'Characters', icon: Users },
  { path: '/bulletin', label: 'Bulletin', icon: Newspaper },
  { path: '/arcana', label: 'Arcana', icon: Atom },
  { path: '/shop', label: 'Shop', icon: Store },
  { path: '/game-rules', label: 'Game Rules', icon: BookMarked },
  { path: '/dm-tools', label: 'DM Tools', icon: Hammer },
];

const gmNavItem = { path: '/gm/spells', label: 'GM Review', icon: Shield };

export function Layout({ children }: LayoutProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const { user, logout, isAuthenticated } = useAuth();
  const [diceCustomizerOpen, setDiceCustomizerOpen] = useState(false);

  const navItems = user?.email === GM_EMAIL ? [...baseNavItems, gmNavItem] : baseNavItems;

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="min-h-screen w-full flex flex-col overflow-x-hidden">
      {/* Tome header – bar like a handbook title strip */}
      <header className="border-b-2 border-faded-gold/50 bg-card/80 backdrop-blur-sm sticky top-0 z-50 hand-drawn-border border-t-0 border-l-0 border-r-0 rounded-none">
        <div className="container mx-auto w-full px-6 py-4 sm:px-10">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-6">
              {/* Mobile Hamburger Menu - now on the far left */}
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
                        <Flame className="w-5 h-5" />
                        Faradhaven
                      </SheetTitle>
                    </SheetHeader>
                    <nav className="flex flex-col p-4 gap-2">
                      {navItems.map((item) => {
                        const Icon = item.icon;
                        const isActive = location.pathname.startsWith(item.path);
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
                            <Icon className="w-5 h-5" />
                            <span className="font-tome-subheading text-sm uppercase tracking-wide">{item.label}</span>
                          </Link>
                        );
                      })}
                    </nav>
                  </SheetContent>
                </Sheet>
              </div>

              <Link to="/" className="flex items-center gap-3 group shrink-0">
                <div>
                  <h1 className="font-tome-heading text-xl text-primary tracking-wide leading-tight">Faradhaven</h1>
                  <p className="text-xs text-muted-foreground font-tome-marginalia">Steampunk RPG System</p>
                </div>
              </Link>

              {/* Desktop Horizontal Navigation (Top Menu) */}
              <nav className="hidden md:flex items-center gap-1 ml-4">
                {navItems.map((item) => {
                  const Icon = item.icon;
                  const isActive = location.pathname === item.path ||
                    location.pathname.startsWith(`${item.path}/`)
                  return (
                    <Link
                      key={item.path}
                      to={item.path}
                      className={cn(
                        'flex items-center gap-2 px-3 py-2 rounded-md transition-all text-sm font-tome-subheading uppercase tracking-wide',
                        isActive
                          ? 'text-primary bg-primary/15 border border-primary/20 shadow-seal'
                          : 'text-muted-foreground hover:text-primary hover:bg-primary/10'
                      )}
                      title={item.label}
                    >
                      <Icon className="w-4 h-4 shrink-0" />
                      <span className="hidden lg:inline">{item.label}</span>
                    </Link>
                  );
                })}
              </nav>
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
                      <Dice5 className="w-4 h-4 mr-2" />
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

      <div className="flex flex-1">
        {/* Main content – tighter padding on mobile, generous on desktop */}
        <main className="flex-1 w-full px-6 py-10 sm:px-10 sm:py-16 md:px-16 md:py-20 lg:px-24 lg:py-24">
          <div className="w-full max-w-6xl mx-auto overflow-x-hidden">
            {children}
          </div>
        </main>
      </div>

      {/* Footer – ledger-style line */}
      <footer className="border-t-2 border-faded-gold/40 py-4 mt-auto">
        <div className="container mx-auto px-4 text-center">
          <p className="text-sm text-muted-foreground font-tome-marginalia">
            A Pupo Production
          </p>
        </div>
      </footer>

      {isAuthenticated && (
        <DiceCustomizer open={diceCustomizerOpen} onOpenChange={setDiceCustomizerOpen} />
      )}
    </div>
  );
}