/**
 * One lore entry per markdown file in repo-root lore-notes/.
 * Add new .md files there; they appear automatically. Use `slugOrder` to control ordering.
 */

export type LoreEntry = {
  id: string;
  label: string;
  filename: string;
  markdown: string;
};

/** Slugs listed first (recommended reading order). Any other file sorts after, A–Z by slug. */
const slugOrder = ['world-overview', 'copper-veins'];

/** Sidebar / tab labels (optional; otherwise slug is title-cased). */
const labelsBySlug: Record<string, string> = {
  'world-overview': 'World overview',
  'copper-veins': 'The Copper Veins',
};

const rawModules = import.meta.glob<string>('../../../../lore-notes/*.md', {
  eager: true,
  query: '?raw',
  import: 'default',
});

function filenameFromPath(path: string): string {
  const seg = path.split('/').pop();
  return seg ?? path;
}

function slugFromFilename(filename: string): string {
  return filename.replace(/\.md$/i, '');
}

function defaultLabelFromSlug(slug: string): string {
  return slug
    .split('-')
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(' ');
}

function sortEntries(a: LoreEntry, b: LoreEntry): number {
  const ia = slugOrder.indexOf(a.id);
  const ib = slugOrder.indexOf(b.id);
  if (ia !== -1 && ib !== -1) return ia - ib;
  if (ia !== -1) return -1;
  if (ib !== -1) return 1;
  return a.id.localeCompare(b.id);
}

export function getLoreEntries(): LoreEntry[] {
  const entries: LoreEntry[] = Object.entries(rawModules).map(([path, markdown]) => {
    const filename = filenameFromPath(path);
    const id = slugFromFilename(filename);
    return {
      id,
      filename,
      label: labelsBySlug[id] ?? defaultLabelFromSlug(id),
      markdown,
    };
  });
  return entries.sort(sortEntries);
}
