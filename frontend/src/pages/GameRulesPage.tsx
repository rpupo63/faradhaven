import { useState } from 'react';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { RaIcon } from '@/components/ui/RaIcon';

import { BestiaryTabContent } from './BestiaryPage';
import { ClassesTabContent } from './ClassesPage';
import { RacesTabContent } from './RacesPage';
import { EffectsTabContent } from './EffectsPage';


export default function GameRulesPage() {
  const { '*': splat } = useParams<{ '*': string }>(); // Captures /classes/:id or /races/:id etc.
  const location = useLocation();
  const navigate = useNavigate();

  const [activeTab, setActiveTab] = useState(() => {
    if (location.pathname.includes('/classes')) return 'classes';
    if (location.pathname.includes('/races')) return 'races';
    if (location.pathname.includes('/effects')) return 'effects';
    if (location.pathname.includes('/bestiary')) return 'bestiary';
    return 'classes'; // Default tab
  });

  // Extract ID if present for detail pages
  const pathParts = splat?.split('/') || [];
  const resourceType = pathParts[0]; // e.g., 'classes', 'races'
  const resourceId = pathParts[1]; // e.g., '123'

  const handleTabChange = (value: string) => {
    setActiveTab(value);
    navigate(`/game-rules/${value}`);
  };

  return (
    <div className="w-full space-y-12">
      <Tabs value={activeTab} onValueChange={handleTabChange} className="w-full">
        <TabsList className="mb-6 grid w-full grid-cols-4 arcane-border bg-card">
          <TabsTrigger value="classes">
            <RaIcon name="book" className="text-sm mr-2" /> Classes
          </TabsTrigger>
          <TabsTrigger value="races">
            <RaIcon name="aura" className="text-sm mr-2" /> Races
          </TabsTrigger>
          <TabsTrigger value="effects">
            <RaIcon name="lightning-bolt" className="text-sm mr-2" /> Effects
          </TabsTrigger>
          <TabsTrigger value="bestiary">
            <RaIcon name="skull" className="text-sm mr-2" /> Bestiary
          </TabsTrigger>
        </TabsList>

        <TabsContent value="classes" className="mt-0">
          <ClassesTabContent classId={activeTab === 'classes' ? resourceId : undefined} />
        </TabsContent>
        <TabsContent value="races" className="mt-0">
          <RacesTabContent raceId={activeTab === 'races' ? resourceId : undefined} />
        </TabsContent>
        <TabsContent value="effects" className="mt-0">
          <EffectsTabContent />
        </TabsContent>
        <TabsContent value="bestiary" className="mt-0">
          <BestiaryTabContent />
        </TabsContent>
      </Tabs>
    </div>
  );
}