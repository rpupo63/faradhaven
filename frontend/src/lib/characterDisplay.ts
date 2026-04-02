/** Display class name without leading "The " (e.g. "The Lorewright" → "Lorewright"). */
export function displayClassName(name: string | undefined): string {
  if (!name) return '';
  return name.replace(/^The\s+/i, '').trim();
}
