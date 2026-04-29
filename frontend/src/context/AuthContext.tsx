import React, { createContext, useContext, useEffect, useState, useCallback } from 'react';
import { login as apiLogin, register as apiRegister, logoutApi, getUserById, setActiveCharacter as apiSetActiveCharacter, type ApiUser, type DicePrefs } from '@/lib/api';

export const DICE_THEME_DEFAULT = 'default';
export const DICE_THEME_COLOR_DEFAULT = '#7A201C';
export const DICE_FONT_COLOR_DEFAULT = '#B8860B';

/** Themes with files under `public/assets/dice-box/themes/` (see @3d-dice/dice-box). */
export const SHIPPED_DICE_THEMES = ['default'] as const;
const SHIPPED_DICE_THEME_SET = new Set<string>(SHIPPED_DICE_THEMES);

function coerceShippedDiceTheme(theme: string | undefined): string {
  const t = theme ?? DICE_THEME_DEFAULT;
  return SHIPPED_DICE_THEME_SET.has(t) ? t : DICE_THEME_DEFAULT;
}

interface AuthContextType {
  token: string | null;
  user: ApiUser | null;
  userId: string | null;
  activeCharacterId: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  logout: () => void;
  setActiveCharacter: (characterId: string | null) => Promise<void>;
  /** Resolved: customizer preview → character colors → standard dice (user account does not affect rolls). */
  activeDicePrefs: Required<DicePrefs>;
  /**
   * Character dice prefs from the API while the sheet is mounted (null = no per-character override).
   * Set only by CharacterSheetPage.
   */
  characterDicePrefs: DicePrefs | null;
  /** Set per-character dice from the sheet (not used by the dice customizer preview). */
  setCharacterDicePrefs: (prefs: DicePrefs | null) => void;
  /** Live preview while Dice Appearance dialog is open; cleared when it closes. */
  setDicePreview: (prefs: DicePrefs | null) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const TOKEN_KEY = 'faradhaven-auth-token';
const USER_ID_KEY = 'faradhaven-user-id';

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = useState<string | null>(() => {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(TOKEN_KEY);
  });
  const [userId, setUserId] = useState<string | null>(() => {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(USER_ID_KEY);
  });
  const [user, setUser] = useState<ApiUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [characterDicePrefsState, setCharacterDicePrefsState] = useState<DicePrefs | null>(null);
  const [dicePreview, setDicePreviewState] = useState<DicePrefs | null>(null);

  // Validate token on mount by fetching user data
  useEffect(() => {
    const validateToken = async () => {
      if (!token || !userId) {
        setIsLoading(false);
        return;
      }

      try {
        const userData = await getUserById(userId, token);
        setUser(userData);
      } catch {
        // Token is invalid, clear auth state
        localStorage.removeItem(TOKEN_KEY);
        localStorage.removeItem(USER_ID_KEY);
        setToken(null);
        setUserId(null);
        setUser(null);
      } finally {
        setIsLoading(false);
      }
    };

    validateToken();
  }, [token, userId]);

  const login = useCallback(async (email: string, password: string) => {
    const response = await apiLogin({ email, password });

    localStorage.setItem(TOKEN_KEY, response.token);
    localStorage.setItem(USER_ID_KEY, response.user_id);
    setToken(response.token);
    setUserId(response.user_id);

    // Fetch user data
    try {
      const userData = await getUserById(response.user_id, response.token);
      setUser(userData);
    } catch {
      // If fetching user fails, we still have the token
      setUser(null);
    }
  }, []);

  const register = useCallback(async (name: string, email: string, password: string) => {
    const response = await apiRegister({ name, email, password });

    localStorage.setItem(TOKEN_KEY, response.token);
    localStorage.setItem(USER_ID_KEY, response.user_id);
    setToken(response.token);
    setUserId(response.user_id);

    // Fetch user data
    try {
      const userData = await getUserById(response.user_id, response.token);
      setUser(userData);
    } catch {
      // If fetching user fails, we still have the token
      setUser(null);
    }
  }, []);

  const logout = useCallback(() => {
    logoutApi(); // fire and forget
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_ID_KEY);
    setToken(null);
    setUserId(null);
    setUser(null);
    setCharacterDicePrefsState(null);
    setDicePreviewState(null);
  }, []);

  const setActiveCharacter = useCallback(async (characterId: string | null) => {
    if (!userId || !token) {
      throw new Error('Not authenticated');
    }

    await apiSetActiveCharacter(userId, characterId, token);

    // Update local user state
    setUser(prev => prev ? { ...prev, active_character_id: characterId } : null);
  }, [userId, token]);

  const setCharacterDicePrefs = useCallback((prefs: DicePrefs | null) => {
    setCharacterDicePrefsState(prefs);
  }, []);

  const setDicePreview = useCallback((prefs: DicePrefs | null) => {
    setDicePreviewState(prefs);
  }, []);

  // Derive activeCharacterId from user state
  const activeCharacterId = user?.active_character_id ?? null;

  const activeDicePrefs: Required<DicePrefs> = {
    dice_theme: coerceShippedDiceTheme(
      dicePreview?.dice_theme ?? characterDicePrefsState?.dice_theme ?? DICE_THEME_DEFAULT
    ),
    dice_theme_color:
      dicePreview?.dice_theme_color ??
      characterDicePrefsState?.dice_theme_color ??
      DICE_THEME_COLOR_DEFAULT,
    dice_font_color:
      dicePreview?.dice_font_color ??
      characterDicePrefsState?.dice_font_color ??
      DICE_FONT_COLOR_DEFAULT,
  };

  return (
    <AuthContext.Provider
      value={{
        token,
        user,
        userId,
        activeCharacterId,
        isAuthenticated: !!token,
        isLoading,
        login,
        register,
        logout,
        setActiveCharacter,
        activeDicePrefs,
        characterDicePrefs: characterDicePrefsState,
        setCharacterDicePrefs,
        setDicePreview,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
