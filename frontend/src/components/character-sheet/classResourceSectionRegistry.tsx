import type { FC } from 'react';
import type { ClassResourceExtraSectionProps } from './classResourceSectionTypes';
import {
  ElixiristSection,
  IronwrightSection,
  LorewrightSection,
  MutagenActions,
  PistonBrawlerActions,
  PowderMageActions,
  SanguinistSection,
} from './ClassResourceClassSections';

export type { ClassResourceExtraSectionProps } from './classResourceSectionTypes';

type RegistryEntry = { className: string; Component: FC<ClassResourceExtraSectionProps> };

/** Ordered list: first match by exact class display name. */
const CLASS_RESOURCE_EXTRA_SECTIONS: RegistryEntry[] = [
  { className: 'The Elixirist', Component: ElixiristSection },
  { className: 'The Piston Brawler', Component: PistonBrawlerActions },
  { className: 'The Mutagen', Component: MutagenActions },
  { className: 'The Powder Mage', Component: PowderMageActions },
  { className: 'The Ironwright', Component: IronwrightSection },
  { className: 'The Lorewright', Component: LorewrightSection },
  { className: 'The Sanguinist', Component: SanguinistSection },
];

/** Renders the extra interactive section for the character's class, if any. */
export function ClassResourceExtraSections(props: ClassResourceExtraSectionProps) {
  const name = props.sheet.class.name;
  const entry = CLASS_RESOURCE_EXTRA_SECTIONS.find((e) => e.className === name);
  if (!entry) return null;
  const C = entry.Component;
  return <C {...props} />;
}
