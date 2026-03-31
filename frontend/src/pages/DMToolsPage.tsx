import React, { useState, useEffect } from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Layout } from '@/components/Layout';
import { Hammer, Bone, Swords } from 'lucide-react';
import { useLocation } from 'react-router-dom';

// Import the content components for each tab
import BattleMapPage from './BattleMapPage';
import MonsterBuilderDashboardPage from './MonsterBuilderDashboardPage';

type DMToolTab = 'battle-map' | 'monster-builder';

export default function DMToolsPage() {
  const location = useLocation();
  const queryParams = React.useMemo(() => new URLSearchParams(location.search), [location.search]);

  const [activeTab, setActiveTab] = useState<DMToolTab>('battle-map');
  const [mapIdFromQuery, setMapIdFromQuery] = useState<string | undefined>(undefined);

  useEffect(() => {
    const tabParam = queryParams.get('tab');
    if (tabParam && (tabParam === 'battle-map' || tabParam === 'monster-builder')) {
      setTimeout(() => setActiveTab(tabParam), 0);
    }

    const mapIdParam = queryParams.get('mapId');
    setTimeout(() => setMapIdFromQuery(mapIdParam || undefined), 0);
  }, [location.search, queryParams]);

  return (
      <div className="space-y-6">
        <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as DMToolTab)} className="w-full">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border pb-4">
            <h1 className="text-3xl font-bold text-primary glow-text">DM Tools</h1>
            <div className="order-last md:order-none w-full md:w-auto overflow-x-auto">
              <TabsList className="grid grid-cols-2 w-full md:w-auto min-w-[300px]">
                <TabsTrigger value="battle-map" className="gap-2 px-3">
                  <Swords className="h-4 w-4" />
                  <span className="hidden sm:inline">Battle Map</span>
                </TabsTrigger>
                <TabsTrigger value="monster-builder" className="gap-2 px-3">
                  <Hammer className="h-4 w-4" />
                  <span className="hidden sm:inline">Monster Builder</span>
                </TabsTrigger>
              </TabsList>
            </div>
          </div>

          <TabsContent value="battle-map" className="mt-6">
            <BattleMapPage initialMapId={mapIdFromQuery} />
          </TabsContent>
          <TabsContent value="monster-builder" className="mt-6">
            <MonsterBuilderDashboardPage />
          </TabsContent>
        </Tabs>
      </div>
  );
}
