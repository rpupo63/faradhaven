import { useState } from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Spellbook } from '@/components/Spellbook';
import { ElementTable } from '@/components/arcanum/ElementTable';
import { BookOpen, Atom } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { getComponents } from '@/lib/api';
import { LoadingQuill } from '@/components/LoadingQuill';

export default function ArcanumSpellbookPage() {
  const [activeTab, setActiveTab] = useState('spellbook');

  const { data: components, isLoading: isLoadingComponents, error: componentsError } = useQuery({
    queryKey: ['components'],
    queryFn: getComponents,
  });

  return (
    <div className="container mx-auto py-8">
      <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
        <TabsList className="grid w-full grid-cols-2 md:w-auto min-w-[300px] mb-8">
          <TabsTrigger value="spellbook" className="gap-2 px-3">
            <BookOpen className="w-4 h-4" />
            <span>Spellbook</span>
          </TabsTrigger>
          <TabsTrigger value="arcanum" className="gap-2 px-3">
            <Atom className="w-4 h-4" />
            <span>Arcanum</span>
          </TabsTrigger>
        </TabsList>

        <TabsContent value="spellbook" className="mt-6">
          <Spellbook />
        </TabsContent>

        <TabsContent value="arcanum" className="mt-6">
          {isLoadingComponents ? (
            <div className="flex items-center justify-center py-12">
              <LoadingQuill label="Loading arcanum components..." />
            </div>
          ) : componentsError ? (
            <div className="arcane-border rounded-lg p-8 text-center">
              <Atom className="w-12 h-12 mx-auto mb-4 text-destructive" />
              <h3 className="font-tome-heading text-lg text-destructive">Failed to Load Arcanum</h3>
              <p className="text-sm text-muted-foreground font-tome-marginalia">
                {componentsError instanceof Error ? componentsError.message : 'An error occurred'}
              </p>
            </div>
          ) : (
            <ElementTable components={components || []} />
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
}
