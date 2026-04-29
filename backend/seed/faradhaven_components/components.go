package faradhaven_components

import (
	"github.com/rpupo63/faradhaven/backend/models"
)

// ComponentSeed defines component data for seeding.
type ComponentSeed struct {
	Name           string
	Symbol         string
	RpgAwesomeIcon string
	Category       models.ComponentCategory
	Description    string
	Element        string
	Tier           int
}

// AllComponents returns all component seeds.
func AllComponents() []ComponentSeed {
	return []ComponentSeed{
		// =====================================================================
		// FORMA (SHAPE) - REQUIRED PILLAR
		// The physical manifestation and geometric delivery of the magic.
		// =====================================================================
		{Name: "Projectile", Symbol: "Pj", RpgAwesomeIcon: "supersonic-arrow", Category: models.CategoryForma, Description: "Launches a discrete manifestation at a point or target. Ranged.", Tier: 1},
		{Name: "Beam", Symbol: "Bm", RpgAwesomeIcon: "laser-blast", Category: models.CategoryForma, Description: "Projects a continuous linear manifestation. Pierces linearly.", Tier: 1},
		{Name: "Nova", Symbol: "Nv", RpgAwesomeIcon: "explosion", Category: models.CategoryForma, Description: "Explodes outward spherically from an origin point.", Tier: 1},
		{Name: "Wall", Symbol: "Wl", RpgAwesomeIcon: "barrier", Category: models.CategoryForma, Description: "Constructs a linear or shaped barrier/partition.", Tier: 1},
		{Name: "Zone", Symbol: "Zn", RpgAwesomeIcon: "circle-of-circles", Category: models.CategoryForma, Description: "Manifests a persistent volumetric area of effect.", Tier: 1},
		{Name: "Cone", Symbol: "Cn", RpgAwesomeIcon: "burst-blob", Category: models.CategoryForma, Description: "Radiates outward in a widening directional spread.", Tier: 1},
		{Name: "Aura", Symbol: "Au", RpgAwesomeIcon: "aura", Category: models.CategoryForma, Description: "Clings to a target, radiating a localized field.", Tier: 1},
		{Name: "Touch", Symbol: "Tc", RpgAwesomeIcon: "hand", Category: models.CategoryForma, Description: "Delivers the magic through direct physical contact. Range: self/melee.", Tier: 1},
		{Name: "Arc", Symbol: "Ac", RpgAwesomeIcon: "plain-dagger", Category: models.CategoryForma, Description: "Follows a curved trajectory that can bend around partial cover.", Tier: 2},
		{Name: "Ring", Symbol: "Rg", RpgAwesomeIcon: "ring", Category: models.CategoryForma, Description: "Creates a hollow perimeter effect, leaving the center unaffected.", Tier: 2},
		{Name: "Pillar", Symbol: "Pr", RpgAwesomeIcon: "stone-pillar", Category: models.CategoryForma, Description: "Erupts as a vertical cylinder from a chosen point.", Tier: 2},
		{Name: "Orbit", Symbol: "Ob", RpgAwesomeIcon: "orbit", Category: models.CategoryForma, Description: "Circles around an anchor while persistently affecting nearby space.", Tier: 2},
		{Name: "Lance", Symbol: "La", RpgAwesomeIcon: "spear-head", Category: models.CategoryForma, Description: "Compresses the spell into an ultra-narrow, high-precision line.", Tier: 2},

		// =====================================================================
		// SCOPUS (TARGETING)
		// The anchor point or entity the magic is allowed to interact with.
		// =====================================================================
		{Name: "Target", Symbol: "Tg", RpgAwesomeIcon: "on-target", Category: models.CategoryScopus, Description: "Locks onto a specific, distinct external entity or object.", Tier: 1},
		{Name: "Self", Symbol: "Sf", RpgAwesomeIcon: "player", Category: models.CategoryScopus, Description: "Anchors the magic exclusively to the caster.", Tier: 1},
		{Name: "Ground", Symbol: "Gd", RpgAwesomeIcon: "groundbreaker", Category: models.CategoryScopus, Description: "Anchors the magic to a spatial coordinate or physical surface.", Tier: 1},
		{Name: "Chain", Symbol: "Ch", RpgAwesomeIcon: "chain", Category: models.CategoryScopus, Description: "Allows the magic to jump between proximate valid targets.", Tier: 2},
		{Name: "Ally", Symbol: "Al", RpgAwesomeIcon: "team-upgrade", Category: models.CategoryScopus, Description: "Restricts valid targets to friendly entities.", Tier: 1},
		{Name: "Enemy", Symbol: "En", RpgAwesomeIcon: "crossed-swords", Category: models.CategoryScopus, Description: "Restricts valid targets to hostile entities.", Tier: 1},
		{Name: "Object", Symbol: "Oj", RpgAwesomeIcon: "wooden-sign", Category: models.CategoryScopus, Description: "Targets unattended objects and structures rather than creatures.", Tier: 1},
		{Name: "Marked", Symbol: "Mk", RpgAwesomeIcon: "targeting", Category: models.CategoryScopus, Description: "Can only affect entities previously marked by the caster.", Tier: 2},
		{Name: "Area-First", Symbol: "Af", RpgAwesomeIcon: "reticle", Category: models.CategoryScopus, Description: "Selects a location first, then resolves to nearest valid entities in the area.", Tier: 2},
		{Name: "LOS-Only", Symbol: "Lo", RpgAwesomeIcon: "focused-lightning", Category: models.CategoryScopus, Description: "Requires uninterrupted line of sight between caster and anchor.", Tier: 1},
		{Name: "Through-Walls", Symbol: "Tw", RpgAwesomeIcon: "stone-wall", Category: models.CategoryScopus, Description: "Can anchor through opaque barriers at increased instability.", Tier: 2},

		// =====================================================================
		// ESSENTIA (DOMAINS & MATTER)
		// The fundamental matter, energy, or abstract concept being manipulated.
		// =====================================================================
		{Name: "Ignis", Symbol: "Ig", RpgAwesomeIcon: "fire", Category: models.CategoryEssentia, Description: "The essence of fire, heat, and combustion.", Element: "fire", Tier: 1},
		{Name: "Aqua", Symbol: "Aq", RpgAwesomeIcon: "water-drop", Category: models.CategoryEssentia, Description: "The essence of water, ice, and fluidity.", Element: "cold", Tier: 1},
		{Name: "Terra", Symbol: "Te", RpgAwesomeIcon: "mountains", Category: models.CategoryEssentia, Description: "The essence of earth, stone, and particulates.", Element: "bludgeoning", Tier: 1},
		{Name: "Aer", Symbol: "Ae", RpgAwesomeIcon: "feather-wing", Category: models.CategoryEssentia, Description: "The essence of air, wind, and gas.", Element: "thunder", Tier: 1},
		{Name: "Ferrum", Symbol: "Fe", RpgAwesomeIcon: "anvil", Category: models.CategoryEssentia, Description: "The essence of metal and magnetism.", Element: "slashing", Tier: 1},
		{Name: "Fulgur", Symbol: "Fu", RpgAwesomeIcon: "lightning-bolt", Category: models.CategoryEssentia, Description: "The essence of lightning and electrical charge.", Element: "lightning", Tier: 1},
		{Name: "Lux", Symbol: "Lu", RpgAwesomeIcon: "sun", Category: models.CategoryEssentia, Description: "The essence of light and radiance.", Element: "radiant", Tier: 1},
		{Name: "Umbra", Symbol: "Um", RpgAwesomeIcon: "hood", Category: models.CategoryEssentia, Description: "The essence of shadow and darkness.", Element: "necrotic", Tier: 1},
		{Name: "Sonus", Symbol: "So", RpgAwesomeIcon: "microphone", Category: models.CategoryEssentia, Description: "The essence of sound and acoustic vibration.", Element: "thunder", Tier: 1},
		{Name: "Odor", Symbol: "Od", RpgAwesomeIcon: "bottle-vapors", Category: models.CategoryEssentia, Description: "The essence of smell, scent, and pheromones.", Tier: 1},
		{Name: "Vita", Symbol: "Vi", RpgAwesomeIcon: "health", Category: models.CategoryEssentia, Description: "The abstract essence of life force, healing, and biological growth.", Element: "radiant", Tier: 1},
		{Name: "Mortis", Symbol: "Mo", RpgAwesomeIcon: "death-skull", Category: models.CategoryEssentia, Description: "The abstract essence of death, decay, and rot.", Element: "necrotic", Tier: 1},
		{Name: "Arcanum", Symbol: "Ar", RpgAwesomeIcon: "crystal-ball", Category: models.CategoryEssentia, Description: "The raw, unshaped essence of pure magical force.", Element: "force", Tier: 1},
		{Name: "Spatium", Symbol: "Sp", RpgAwesomeIcon: "perspective-dice-random", Category: models.CategoryEssentia, Description: "The abstract fabric of physical space and dimensions.", Tier: 2},
		{Name: "Chronos", Symbol: "Ti", RpgAwesomeIcon: "hourglass", Category: models.CategoryEssentia, Description: "The abstract flow of time and temporal causality.", Tier: 2},
		{Name: "Venenum", Symbol: "Vn", RpgAwesomeIcon: "poison-cloud", Category: models.CategoryEssentia, Description: "The essence of poison, toxin, and biological corruption.", Element: "poison", Tier: 1},
		{Name: "Acidum", Symbol: "Ac", RpgAwesomeIcon: "acid", Category: models.CategoryEssentia, Description: "The essence of acid, corrosion, and chemical dissolution.", Element: "acid", Tier: 1},
		{Name: "Mens", Symbol: "Mn", RpgAwesomeIcon: "overmind", Category: models.CategoryEssentia, Description: "The abstract essence of thought, memory, and conscious will.", Element: "psychic", Tier: 2},

		// Pathos (Emotional Spectrum)
		{Name: "Anger", Symbol: "An", RpgAwesomeIcon: "crossed-swords", Category: models.CategoryEssentia, Description: "The emotional resonance of rage, hostility, and aggression.", Element: "psychic", Tier: 1},
		{Name: "Sadness", Symbol: "Sd", RpgAwesomeIcon: "cold-heart", Category: models.CategoryEssentia, Description: "The emotional resonance of grief, despair, and lethargy.", Element: "psychic", Tier: 1},
		{Name: "Fear", Symbol: "Fr", RpgAwesomeIcon: "player-dodge", Category: models.CategoryEssentia, Description: "The emotional resonance of terror, dread, and flight.", Element: "psychic", Tier: 1},
		{Name: "Happiness", Symbol: "Hp", RpgAwesomeIcon: "hearts", Category: models.CategoryEssentia, Description: "The emotional resonance of joy, euphoria, and laughter.", Element: "psychic", Tier: 1},
		{Name: "Disgust", Symbol: "Di", RpgAwesomeIcon: "biohazard", Category: models.CategoryEssentia, Description: "The emotional resonance of revulsion, nausea, and rejection.", Element: "psychic", Tier: 1},
		{Name: "Embarrassment", Symbol: "Em", RpgAwesomeIcon: "glass-heart", Category: models.CategoryEssentia, Description: "The emotional resonance of shame, awkwardness, and self-doubt.", Element: "psychic", Tier: 1},

		// =====================================================================
		// ACTIO (KINETIC VERBS)
		// What the magic physically does to the Essentia or Scopus.
		// =====================================================================
		{Name: "Push", Symbol: "Pu", RpgAwesomeIcon: "forward", Category: models.CategoryActio, Description: "Applies outward or repelling kinetic force.", Tier: 1},
		{Name: "Pull", Symbol: "Pl", RpgAwesomeIcon: "magnet", Category: models.CategoryActio, Description: "Applies inward or attracting kinetic force.", Tier: 1},
		{Name: "Grab", Symbol: "Gr", RpgAwesomeIcon: "grappling-hook", Category: models.CategoryActio, Description: "Seizes, holds, or grapples the target in place.", Tier: 1},
		{Name: "Spin", Symbol: "Sn", RpgAwesomeIcon: "spiral-shell", Category: models.CategoryActio, Description: "Applies rotational or spiraling force.", Tier: 1},
		{Name: "Crush", Symbol: "Cr", RpgAwesomeIcon: "crush", Category: models.CategoryActio, Description: "Applies compressing, inward-crushing pressure.", Tier: 1},
		{Name: "Pierce", Symbol: "Pi", RpgAwesomeIcon: "spear-head", Category: models.CategoryActio, Description: "Applies highly concentrated penetrative force.", Tier: 1},
		{Name: "Create", Symbol: "Ct", RpgAwesomeIcon: "sprout", Category: models.CategoryActio, Description: "Manifests the Essentia from nothing into reality.", Tier: 2},
		{Name: "Destroy", Symbol: "Ds", RpgAwesomeIcon: "demolish", Category: models.CategoryActio, Description: "Obliterates or nullifies the target or Essentia.", Tier: 2},
		{Name: "Mutate", Symbol: "Mu", RpgAwesomeIcon: "recycle", Category: models.CategoryActio, Description: "Alters the fundamental state, shape, or nature of the target.", Tier: 2},
		{Name: "Bind", Symbol: "Bn", RpgAwesomeIcon: "three-keys", Category: models.CategoryActio, Description: "Tethers two entities, Essentias, or forces together.", Tier: 2},
		{Name: "Sense", Symbol: "Se", RpgAwesomeIcon: "telescope", Category: models.CategoryActio, Description: "Extends perception to detect, identify, or observe a target or phenomenon.", Tier: 1},
		{Name: "Ward", Symbol: "Wd", RpgAwesomeIcon: "shield", Category: models.CategoryActio, Description: "Erects a protective barrier or nullifying field against a force or entity.", Tier: 1},
		{Name: "Conceal", Symbol: "Co", RpgAwesomeIcon: "cloak-and-dagger", Category: models.CategoryActio, Description: "Hides, obscures, or projects a false appearance over a target.", Tier: 1},
		{Name: "Summon", Symbol: "Sm", RpgAwesomeIcon: "spawn-node", Category: models.CategoryActio, Description: "Calls an entity or object from another location or plane into presence.", Tier: 2},
		{Name: "Move", Symbol: "Mv", RpgAwesomeIcon: "player-teleport", Category: models.CategoryActio, Description: "Instantly translates a target across space, bypassing the intervening distance.", Tier: 2},
		{Name: "Compel", Symbol: "Cm", RpgAwesomeIcon: "crown", Category: models.CategoryActio, Description: "Overrides or bends the conscious will of a target, forcing behavior.", Tier: 2},

		// =====================================================================
		// MAGNITUDO (SCALE & MODIFIERS)
		// The mathematical dials that adjust the spell's parameters.
		// =====================================================================
		{Name: "Increase", Symbol: "In", RpgAwesomeIcon: "health-increase", Category: models.CategoryMagnitudo, Description: "Raises a specific property (e.g., speeding up time, raising temperature).", Tier: 1},
		{Name: "Decrease", Symbol: "Dc", RpgAwesomeIcon: "health-decrease", Category: models.CategoryMagnitudo, Description: "Lowers a specific property (e.g., slowing time, dropping temperature).", Tier: 1},
		{Name: "Strong", Symbol: "St", RpgAwesomeIcon: "muscle-up", Category: models.CategoryMagnitudo, Description: "Amplifies the overall power, damage, or resistance to nullification.", Tier: 1},
		{Name: "Weak", Symbol: "Wk", RpgAwesomeIcon: "small-fire", Category: models.CategoryMagnitudo, Description: "Diminishes the overall power, making the spell subtle or highly efficient to cast.", Tier: 1},
		{Name: "Extreme", Symbol: "Xt", RpgAwesomeIcon: "nuclear", Category: models.CategoryMagnitudo, Description: "Pushes the magic to its absolute dimensional limit. Highly volatile.", Tier: 2},
		{Name: "Inverse", Symbol: "Iv", RpgAwesomeIcon: "reverse", Category: models.CategoryMagnitudo, Description: "Inverts the fundamental nature of the paired Essentia or Actio (e.g. fire→cold, push→pull, create→destroy).", Tier: 1},

		// =====================================================================
		// LOGICA (SEQUENTIAL LINKS)
		// Establishes order and narrative causality between component groups (reuse allowed).
		// =====================================================================
		{Name: "If", Symbol: "If", RpgAwesomeIcon: "uncertainty", Category: models.CategoryLogica, Description: "Opens a conditional branch: the following segment applies when the stated condition (table agreement) is met before resolving later phases.", Tier: 1},
		{Name: "Then", Symbol: "Th", RpgAwesomeIcon: "forward", Category: models.CategoryLogica, Description: "Declares a sequenced phase: earlier magic completes first; this segment builds on that outcome (e.g. water pooled, then cold freezes it).", Tier: 1},
		{Name: "Therefore", Symbol: "Tf", RpgAwesomeIcon: "gavel", Category: models.CategoryLogica, Description: "Marks a concluding causal beat tying prior phases into a final effect or summary outcome.", Tier: 1},
	}
}
