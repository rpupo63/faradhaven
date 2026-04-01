import { useState } from 'react';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { BookOpen, Sparkles, Zap, Skull } from 'lucide-react';

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
        <TabsList className="grid w-full grid-cols-4 arcane-border bg-card">
          <TabsTrigger value="classes">
            <BookOpen className="w-4 h-4 mr-2" /> Classes
          </TabsTrigger>
          <TabsTrigger value="races">
            <Sparkles className="w-4 h-4 mr-2" /> Races
          </TabsTrigger>
          <TabsTrigger value="effects">
            <Zap className="w-4 h-4 mr-2" /> Effects
          </TabsTrigger>
          <TabsTrigger value="bestiary">
            <Skull className="w-4 h-4 mr-2" /> Bestiary
          </TabsTrigger>
        </TabsList>

        <TabsContent value="classes">
          <ClassesTabContent classId={activeTab === 'classes' ? resourceId : undefined} />
        </TabsContent>
        <TabsContent value="races">
          <RacesTabContent raceId={activeTab === 'races' ? resourceId : undefined} />
        </TabsContent>
        <TabsContent value="effects">
          <EffectsTabContent />
        </TabsContent>
        <TabsContent value="bestiary">
          <BestiaryTabContent />
        </TabsContent>
      </Tabs>
    </div>
  );
}