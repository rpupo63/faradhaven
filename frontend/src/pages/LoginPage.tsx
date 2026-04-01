import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Flame, LogIn, UserPlus, Loader2 } from 'lucide-react';

type AuthMode = 'login' | 'register';

export default function LoginPage() {
  const navigate = useNavigate();
  const { login, register } = useAuth();
  
  const [mode, setMode] = useState<AuthMode>('login');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Form fields
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    
    if (mode === 'register' && password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }
    
    if (mode === 'register' && password.length < 6) {
      setError('Password must be at least 6 characters');
      return;
    }
    
    setIsLoading(true);
    
    try {
      if (mode === 'login') {
        await login(email, password);
      } else {
        await register(name, email, password);
      }
      navigate('/');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Authentication failed');
    } finally {
      setIsLoading(false);
    }
  };

  const toggleMode = () => {
    setMode(mode === 'login' ? 'register' : 'login');
    setError(null);
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center px-8 sm:px-12 md:px-20 py-12 md:py-24">
      <div className="w-full max-w-md">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center p-4 rounded-full border-2 border-faded-gold/50 bg-primary/10 mb-4">
            <Flame className="w-10 h-10 text-primary" />
          </div>
          <h1 className="font-tome-heading text-3xl text-primary tracking-wide">Faradhaven</h1>
          <p className="text-muted-foreground font-tome-marginalia mt-1">
            {mode === 'login' ? 'Welcome back, adventurer' : 'Begin your legend'}
          </p>
        </div>

        {/* Auth Form Card */}
        <div className="arcane-border rounded-xl p-8 bg-card/80 backdrop-blur-sm">
          <h2 className="font-tome-heading text-xl text-center mb-6">
            {mode === 'login' ? 'Sign In' : 'Create Account'}
          </h2>

          <form onSubmit={handleSubmit} className="space-y-4">
            {mode === 'register' && (
              <div className="space-y-2">
                <Label htmlFor="name" className="font-tome-marginalia text-sm">
                  Name
                </Label>
                <Input
                  id="name"
                  type="text"
                  placeholder="Your character's keeper"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  required
                  disabled={isLoading}
                  className="bg-background/50"
                />
              </div>
            )}

            <div className="space-y-2">
              <Label htmlFor="email" className="font-tome-marginalia text-sm">
                Email
              </Label>
              <Input
                id="email"
                type="email"
                placeholder="your@email.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                disabled={isLoading}
                className="bg-background/50"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="password" className="font-tome-marginalia text-sm">
                Password
              </Label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                disabled={isLoading}
                className="bg-background/50"
              />
            </div>

            {mode === 'register' && (
              <div className="space-y-2">
                <Label htmlFor="confirmPassword" className="font-tome-marginalia text-sm">
                  Confirm Password
                </Label>
                <Input
                  id="confirmPassword"
                  type="password"
                  placeholder="••••••••"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  required
                  disabled={isLoading}
                  className="bg-background/50"
                />
              </div>
            )}

            {error && (
              <div className="text-destructive text-sm text-center p-2 rounded bg-destructive/10 border border-destructive/20">
                {error}
              </div>
            )}

            <Button
              type="submit"
              variant="seal"
              size="lg"
              className="w-full"
              disabled={isLoading}
            >
              {isLoading ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : mode === 'login' ? (
                <>
                  <LogIn className="w-4 h-4" />
                  Sign In
                </>
              ) : (
                <>
                  <UserPlus className="w-4 h-4" />
                  Create Account
                </>
              )}
            </Button>
          </form>


          <div className="text-center">
            <p className="text-muted-foreground text-sm mb-2">
              {mode === 'login' ? "Don't have an account?" : 'Already have an account?'}
            </p>
            <Button
              variant="ghost"
              onClick={toggleMode}
              disabled={isLoading}
              className="text-primary hover:text-primary/80"
            >
              {mode === 'login' ? 'Create one here' : 'Sign in instead'}
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
