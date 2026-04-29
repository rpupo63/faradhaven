import type { SpellForgeVesselComponent } from '@/components/spell-forge-vessels/types';
import { DefaultForgeVessel } from '@/components/spell-forge-vessels/DefaultForgeVessel';
import { MutagenForgeVessel } from '@/components/spell-forge-vessels/classes/MutagenForgeVessel';
import { IronwrightForgeVessel } from '@/components/spell-forge-vessels/classes/IronwrightForgeVessel';
import { SanguinistForgeVessel } from '@/components/spell-forge-vessels/classes/SanguinistForgeVessel';
import { ElixiristForgeVessel } from '@/components/spell-forge-vessels/classes/ElixiristForgeVessel';
import { PistonBrawlerForgeVessel } from '@/components/spell-forge-vessels/classes/PistonBrawlerForgeVessel';
import { PowderMageForgeVessel } from '@/components/spell-forge-vessels/classes/PowderMageForgeVessel';
import { LorewrightForgeVessel } from '@/components/spell-forge-vessels/classes/LorewrightForgeVessel';
import { SyllogistForgeVessel } from '@/components/spell-forge-vessels/classes/SyllogistForgeVessel';

/** Matches `backend/seed/faradhaven_classes` AllClasses `Name` fields. */
const REGISTRY: Record<string, SpellForgeVesselComponent> = {
  'The Mutagen': MutagenForgeVessel,
  'The Ironwright': IronwrightForgeVessel,
  'The Sanguinist': SanguinistForgeVessel,
  'The Elixirist': ElixiristForgeVessel,
  'The Piston Brawler': PistonBrawlerForgeVessel,
  'The Powder Mage': PowderMageForgeVessel,
  'The Lorewright': LorewrightForgeVessel,
  'The Syllogist': SyllogistForgeVessel,
};

export function resolveForgeVessel(className: string | undefined): SpellForgeVesselComponent {
  if (!className) return DefaultForgeVessel;
  return REGISTRY[className] ?? DefaultForgeVessel;
}
