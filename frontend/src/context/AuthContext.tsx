import React, { createContext, useContext, useEffect, useState, useCallback } from 'react';
import { login as apiLogin, register as apiRegister, logoutApi, getUserById, setActiveCharacter as apiSetActiveCharacter, updateUserDicePrefs, type ApiUser, type DicePrefs } from '@/lib/api';

export const DICE_THEME_DEFAULT = 'default';
export const DICE_THEME_COLOR_DEFAULT = '#7A201C';
export const DICE_FONT_COLOR_DEFAULT = '#B8860B';

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
  /** Resolved dice prefs: character override if set, otherwise user defaults */
  activeDicePrefs: Required<DicePrefs>;
  /**
   * The character's dice prefs as saved in the DB (null fields = no override).
   * Set by CharacterSheetPage on load, cleared on unmount.
   */
  characterDicePrefs: DicePrefs | null;
  /** Set a per-character dice override (called by CharacterSheetPage on load, cleared on unmount) */
  setDicePrefsOverride: (prefs: DicePrefs | null) => void;
  /** Persist user-level dice preferences to the backend */
  updateUserDicePreferences: (prefs: DicePrefs) => Promise<void>;
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
  const [dicePrefsOverride, setDicePrefsOverrideState] = useState<DicePrefs | null>(null);

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
  }, []);

  const setActiveCharacter = useCallback(async (characterId: string | null) => {
    if (!userId || !token) {
      throw new Error('Not authenticated');
    }

    await apiSetActiveCharacter(userId, characterId, token);

    // Update local user state
    setUser(prev => prev ? { ...prev, active_character_id: characterId } : null);
  }, [userId, token]);

  const setDicePrefsOverride = useCallback((prefs: DicePrefs | null) => {
    setDicePrefsOverrideState(prefs);
  }, []);

  const updateUserDicePreferences = useCallback(async (prefs: DicePrefs) => {
    if (!userId || !token) throw new Error('Not authenticated');
    const updated = await updateUserDicePrefs(userId, prefs, token);
    setUser(updated);
  }, [userId, token]);

  // Derive activeCharacterId from user state
  const activeCharacterId = user?.active_character_id ?? null;

  // Resolved dice prefs: character override wins over user defaults, user defaults win over hardcoded
  const activeDicePrefs: Required<DicePrefs> = {
    dice_theme:
      dicePrefsOverride?.dice_theme ??
      user?.dice_theme ??
      DICE_THEME_DEFAULT,
    dice_theme_color:
      dicePrefsOverride?.dice_theme_color ??
      user?.dice_theme_color ??
      DICE_THEME_COLOR_DEFAULT,
    dice_font_color:
      dicePrefsOverride?.dice_font_color ??
      user?.dice_font_color ??
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
        characterDicePrefs: dicePrefsOverride,
        setDicePrefsOverride,
        updateUserDicePreferences,
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
