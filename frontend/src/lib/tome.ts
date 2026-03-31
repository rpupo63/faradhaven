/**
 * Tome design tokens – hand-drawn, ink-and-quill aesthetic.
 * Use for inline styles or as reference when building components.
 */

/** Hand-drawn border CSS object – irregular quill feel */
export const handDrawnBorder = {
  borderStyle: "solid",
  borderWidth: "3px 4px 3px 5px",
  borderRadius: "4% 95% 6% 95% / 95% 4% 92% 5%",
  borderColor: "#2C2421",
} as const;

/** Stained Vellum palette – hex codes for non-Tailwind usage */
export const tomeColors = {
  parchment: "#F4EBD0",
  ink: "#2C2421",
  driedBlood: "#7A201C",
  fadedGold: "#B8860B",
  charcoal: "#4A4A4A",
} as const;

/** CSS filter for ink-wash / sketch-style images */
export const inkSketchFilter = "sepia(0.5) contrast(1.2) grayscale(1)";
