
/**
 * Converts a D&D-style cost string (e.g., "1 gp", "5 sp", "4 cp") to total copper pieces (cp).
 * 1 pp = 1000 cp, 1 gp = 100 cp, 1 sp = 10 cp, 1 cp = 1 cp.
 */
export const parseCostToCp = (costStr: string): number => {
  if (!costStr) return 0;
  // This regex matches a number followed by a unit (gp, sp, cp, pp) anywhere in the string.
  // It handles cases like "Rope (50ft) (1gp)" or just "1 gp".
  const match = costStr.match(/(\d+)\s*(gp|sp|cp|pp)/i);
  if (!match) return 0;
  const value = parseInt(match[1], 10);
  const unit = match[2].toLowerCase();
  switch (unit) {
    case 'pp': return value * 1000;
    case 'gp': return value * 100;
    case 'sp': return value * 10;
    case 'cp': return value * 1;
    default: return value;
  }
};

/**
 * Formats a copper piece (cp) value back into a readable D&D-style currency string.
 * Optionally returns an array of { value, unit } for custom rendering.
 */
export const formatCpToDisplay = (cp: number): string => {
  if (cp <= 0) return '0 cp';
  const gp = Math.floor(cp / 100);
  const sp = Math.floor((cp % 100) / 10);
  const remainingCp = cp % 10;

  const parts: string[] = [];
  if (gp > 0) parts.push(`${gp} gp`);
  if (sp > 0) parts.push(`${sp} sp`);
  if (remainingCp > 0 || parts.length === 0) parts.push(`${remainingCp} cp`);
  
  return parts.join(' ');
};
