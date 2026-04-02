import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { GameProvider } from "@/context/GameContext";
import { AuthProvider, useAuth } from "@/context/AuthContext";
import { Layout } from "@/components/Layout";
import { LoadingQuill } from "@/components/LoadingQuill";
import Index from "./pages/Index";
import CharactersPage from "./pages/CharactersPage";
import ArcanumSpellbookPage from "./pages/ArcanumSpellbookPage";
import ShopPage from "./pages/ShopPage";
import BulletinPage from "./pages/BulletinPage";
import LoginPage from "./pages/LoginPage";
import NotFound from "./pages/NotFound";
import DMToolsPage from "./pages/DMToolsPage";
import MonsterSheetPage from "./pages/MonsterSheetPage";
import BattleMapPage from "./pages/BattleMapPage";
import GameRulesPage from "./pages/GameRulesPage";
import CharacterSheetPage from "./pages/CharacterSheetPage";
import GMSpellReviewPage from "./pages/GMSpellReviewPage";

const queryClient = new QueryClient();

// Protected route wrapper - redirects to login if not authenticated
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingQuill label="Verifying credentials..." />
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
}

// Public route wrapper - redirects to characters if already authenticated
function PublicOnlyRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingQuill label="Loading..." />
      </div>
    );
  }

  if (isAuthenticated) {
    return <Navigate to="/characters" replace />;
  }

  return <>{children}</>;
}

const AppRoutes = () => {
  return (
    <Routes>
      {/* Public-only route (login page) */}
      <Route
        path="/login"
        element={
          <PublicOnlyRoute>
            <LoginPage />
          </PublicOnlyRoute>
        }
      />

      {/* Public routes - accessible without authentication */}
      <Route
        path="/"
        element={
          <Layout>
            <Index />
          </Layout>
        }
      />
      <Route
        path="/characters"
        element={
          <ProtectedRoute>
            <Layout>
              <CharactersPage />
            </Layout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/spellbook"
        element={
          <Layout>
            <ArcanumSpellbookPage />
          </Layout>
        }
      />
      <Route
        path="/shop"
        element={
          <Layout>
            <ShopPage />
          </Layout>
        }
      />
      <Route
        path="/game-rules/*"
        element={
          <Layout>
            <GameRulesPage />
          </Layout>
        }
      />
      <Route
        path="/arcana"
        element={
          <Layout>
            <ArcanumSpellbookPage />
          </Layout>
        }
      />
      <Route
        path="/bulletin"
        element={
          <Layout>
            <BulletinPage />
          </Layout>
        }
      />

      <Route
        path="/dm-tools"
        element={
          <Layout>
            <DMToolsPage />
          </Layout>
        }
      />
      {/* Public routes for BattleMapPage - BattleMapPage handles auth internally */}
      <Route
        path="/battle-map"
        element={
          <Layout>
            <BattleMapPage />
          </Layout>
        }
      />
      <Route
        path="/battle-map/:mapId"
        element={
          <Layout>
            <BattleMapPage />
          </Layout>
        }
      />
      <Route
        path="/battle-map/room/:roomCode"
        element={
          <Layout>
            <BattleMapPage />
          </Layout>
        }
      />

      {/* Protected routes - require authentication */}
      <Route
        path="/character/:id/sheet"
        element={
          <ProtectedRoute>
            <Layout>
              <CharacterSheetPage />
            </Layout>
          </ProtectedRoute>
        }
      />
      {/* NEW: Protected route for Monster Sheet */}
      <Route
        path="/monster/:id"
        element={
          <ProtectedRoute>
            <Layout>
              <MonsterSheetPage />
            </Layout>
          </ProtectedRoute>
        }
      />
      {/* GM-only route */}
      <Route
        path="/gm/spells"
        element={
          <ProtectedRoute>
            <Layout>
              <GMSpellReviewPage />
            </Layout>
          </ProtectedRoute>
        }
      />
      <Route path="*" element={<NotFound />} />
    </Routes>
  );
};

const App = () => (
  <QueryClientProvider client={queryClient}>
    <AuthProvider>
      <GameProvider>
        <TooltipProvider>
          <Toaster />
          <Sonner />
          <BrowserRouter
            future={{
              v7_startTransition: true,
              v7_relativeSplatPath: true,
            }}
          >
            <AppRoutes />
          </BrowserRouter>
        </TooltipProvider>
      </GameProvider>
    </AuthProvider>
  </QueryClientProvider>
);

export default App;