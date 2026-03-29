package faradhaven_components

import (
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// ComponentSeed defines component data for seeding.
type ComponentSeed struct {
	Name        string
	Symbol      string
	Category    models.ComponentCategory
	Description string
	Element     string
	Tier        int
}

// AllComponents returns all component seeds.
func AllComponents() []ComponentSeed {
	return []ComponentSeed{
		// =====================================================================
		// FORMA (SHAPE) - REQUIRED PILLAR
		// The physical manifestation and geometric delivery of the magic.
		// =====================================================================
		{Name: "Projectile", Symbol: "Pj", Category: models.CategoryForma, Description: "Launches a discrete manifestation at a point or target. Ranged.", Tier: 1},
		{Name: "Beam", Symbol: "Bm", Category: models.CategoryForma, Description: "Projects a continuous linear manifestation. Pierces linearly.", Tier: 1},
		{Name: "Nova", Symbol: "Nv", Category: models.CategoryForma, Description: "Explodes outward spherically from an origin point.", Tier: 1},
		{Name: "Wall", Symbol: "Wl", Category: models.CategoryForma, Description: "Constructs a linear or shaped barrier/partition.", Tier: 1},
		{Name: "Zone", Symbol: "Zn", Category: models.CategoryForma, Description: "Manifests a persistent volumetric area of effect.", Tier: 1},
		{Name: "Cone", Symbol: "Cn", Category: models.CategoryForma, Description: "Radiates outward in a widening directional spread.", Tier: 1},
		{Name: "Aura", Symbol: "Au", Category: models.CategoryForma, Description: "Clings to a target, radiating a localized field.", Tier: 1},
		{Name: "Touch", Symbol: "Tc", Category: models.CategoryForma, Description: "Delivers the magic through direct physical contact. Range: self/melee.", Tier: 1},

		// =====================================================================
		// SCOPUS (TARGETING)
		// The anchor point or entity the magic is allowed to interact with.
		// =====================================================================
		{Name: "Target", Symbol: "Tg", Category: models.CategoryScopus, Description: "Locks onto a specific, distinct external entity or object.", Tier: 1},
		{Name: "Self", Symbol: "Sf", Category: models.CategoryScopus, Description: "Anchors the magic exclusively to the caster.", Tier: 1},
		{Name: "Ground", Symbol: "Gd", Category: models.CategoryScopus, Description: "Anchors the magic to a spatial coordinate or physical surface.", Tier: 1},
		{Name: "Chain", Symbol: "Ch", Category: models.CategoryScopus, Description: "Allows the magic to jump between proximate valid targets.", Tier: 2},

		// =====================================================================
		// ESSENTIA (DOMAINS & MATTER)
		// The fundamental matter, energy, or abstract concept being manipulated.
		// =====================================================================
		{Name: "Ignis", Symbol: "Ig", Category: models.CategoryEssentia, Description: "The essence of fire, heat, and combustion.", Element: "fire", Tier: 1},
		{Name: "Aqua", Symbol: "Aq", Category: models.CategoryEssentia, Description: "The essence of water, ice, and fluidity.", Element: "cold", Tier: 1},
		{Name: "Terra", Symbol: "Te", Category: models.CategoryEssentia, Description: "The essence of earth, stone, and particulates.", Element: "bludgeoning", Tier: 1},
		{Name: "Aer", Symbol: "Ae", Category: models.CategoryEssentia, Description: "The essence of air, wind, and gas.", Element: "thunder", Tier: 1},
		{Name: "Ferrum", Symbol: "Fe", Category: models.CategoryEssentia, Description: "The essence of metal and magnetism.", Element: "slashing", Tier: 1},
		{Name: "Fulgur", Symbol: "Fu", Category: models.CategoryEssentia, Description: "The essence of lightning and electrical charge.", Element: "lightning", Tier: 1},
		{Name: "Lux", Symbol: "Lu", Category: models.CategoryEssentia, Description: "The essence of light and radiance.", Element: "radiant", Tier: 1},
		{Name: "Umbra", Symbol: "Um", Category: models.CategoryEssentia, Description: "The essence of shadow and darkness.", Element: "necrotic", Tier: 1},
		{Name: "Sonus", Symbol: "So", Category: models.CategoryEssentia, Description: "The essence of sound and acoustic vibration.", Element: "thunder", Tier: 1},
		{Name: "Odor", Symbol: "Od", Category: models.CategoryEssentia, Description: "The essence of smell, scent, and pheromones.", Tier: 1},
		{Name: "Vita", Symbol: "Vi", Category: models.CategoryEssentia, Description: "The abstract essence of life force, healing, and biological growth.", Element: "radiant", Tier: 1},
		{Name: "Mortis", Symbol: "Mo", Category: models.CategoryEssentia, Description: "The abstract essence of death, decay, and rot.", Element: "necrotic", Tier: 1},
		{Name: "Arcanum", Symbol: "Ar", Category: models.CategoryEssentia, Description: "The raw, unshaped essence of pure magical force.", Element: "force", Tier: 1},
		{Name: "Spatium", Symbol: "Sp", Category: models.CategoryEssentia, Description: "The abstract fabric of physical space and dimensions.", Tier: 2},
		{Name: "Chronos", Symbol: "Ti", Category: models.CategoryEssentia, Description: "The abstract flow of time and temporal causality.", Tier: 2},
		{Name: "Venenum", Symbol: "Vn", Category: models.CategoryEssentia, Description: "The essence of poison, toxin, and biological corruption.", Element: "poison", Tier: 1},
		{Name: "Acidum", Symbol: "Ac", Category: models.CategoryEssentia, Description: "The essence of acid, corrosion, and chemical dissolution.", Element: "acid", Tier: 1},
		{Name: "Mens", Symbol: "Mn", Category: models.CategoryEssentia, Description: "The abstract essence of thought, memory, and conscious will.", Element: "psychic", Tier: 2},

		// Pathos (Emotional Spectrum)
		{Name: "Anger", Symbol: "An", Category: models.CategoryEssentia, Description: "The emotional resonance of rage, hostility, and aggression.", Element: "psychic", Tier: 1},
		{Name: "Sadness", Symbol: "Sd", Category: models.CategoryEssentia, Description: "The emotional resonance of grief, despair, and lethargy.", Element: "psychic", Tier: 1},
		{Name: "Fear", Symbol: "Fr", Category: models.CategoryEssentia, Description: "The emotional resonance of terror, dread, and flight.", Element: "psychic", Tier: 1},
		{Name: "Happiness", Symbol: "Hp", Category: models.CategoryEssentia, Description: "The emotional resonance of joy, euphoria, and laughter.", Element: "psychic", Tier: 1},
		{Name: "Disgust", Symbol: "Di", Category: models.CategoryEssentia, Description: "The emotional resonance of revulsion, nausea, and rejection.", Element: "psychic", Tier: 1},
		{Name: "Embarrassment", Symbol: "Em", Category: models.CategoryEssentia, Description: "The emotional resonance of shame, awkwardness, and self-doubt.", Element: "psychic", Tier: 1},

		// =====================================================================
		// ACTIO (KINETIC VERBS)
		// What the magic physically does to the Essentia or Scopus.
		// =====================================================================
		{Name: "Push", Symbol: "Pu", Category: models.CategoryActio, Description: "Applies outward or repelling kinetic force.", Tier: 1},
		{Name: "Pull", Symbol: "Pl", Category: models.CategoryActio, Description: "Applies inward or attracting kinetic force.", Tier: 1},
		{Name: "Grab", Symbol: "Gr", Category: models.CategoryActio, Description: "Seizes, holds, or grapples the target in place.", Tier: 1},
		{Name: "Spin", Symbol: "Sn", Category: models.CategoryActio, Description: "Applies rotational or spiraling force.", Tier: 1},
		{Name: "Crush", Symbol: "Cr", Category: models.CategoryActio, Description: "Applies compressing, inward-crushing pressure.", Tier: 1},
		{Name: "Pierce", Symbol: "Pi", Category: models.CategoryActio, Description: "Applies highly concentrated penetrative force.", Tier: 1},
		{Name: "Create", Symbol: "Ct", Category: models.CategoryActio, Description: "Manifests the Essentia from nothing into reality.", Tier: 2},
		{Name: "Destroy", Symbol: "Ds", Category: models.CategoryActio, Description: "Obliterates or nullifies the target or Essentia.", Tier: 2},
		{Name: "Mutate", Symbol: "Mu", Category: models.CategoryActio, Description: "Alters the fundamental state, shape, or nature of the target.", Tier: 2},
		{Name: "Bind", Symbol: "Bn", Category: models.CategoryActio, Description: "Tethers two entities, Essentias, or forces together.", Tier: 2},
		{Name: "Sense", Symbol: "Se", Category: models.CategoryActio, Description: "Extends perception to detect, identify, or observe a target or phenomenon.", Tier: 1},
		{Name: "Ward", Symbol: "Wd", Category: models.CategoryActio, Description: "Erects a protective barrier or nullifying field against a force or entity.", Tier: 1},
		{Name: "Conceal", Symbol: "Co", Category: models.CategoryActio, Description: "Hides, obscures, or projects a false appearance over a target.", Tier: 1},
		{Name: "Summon", Symbol: "Sm", Category: models.CategoryActio, Description: "Calls an entity or object from another location or plane into presence.", Tier: 2},
		{Name: "Move", Symbol: "Mv", Category: models.CategoryActio, Description: "Instantly translates a target across space, bypassing the intervening distance.", Tier: 2},
		{Name: "Compel", Symbol: "Cm", Category: models.CategoryActio, Description: "Overrides or bends the conscious will of a target, forcing behavior.", Tier: 2},

		// =====================================================================
		// MAGNITUDO (SCALE & MODIFIERS)
		// The mathematical dials that adjust the spell's parameters.
		// =====================================================================
		{Name: "Increase", Symbol: "In", Category: models.CategoryMagnitudo, Description: "Raises a specific property (e.g., speeding up time, raising temperature).", Tier: 1},
		{Name: "Decrease", Symbol: "Dc", Category: models.CategoryMagnitudo, Description: "Lowers a specific property (e.g., slowing time, dropping temperature).", Tier: 1},
		{Name: "Strong", Symbol: "St", Category: models.CategoryMagnitudo, Description: "Amplifies the overall power, damage, or resistance to nullification.", Tier: 1},
		{Name: "Weak", Symbol: "Wk", Category: models.CategoryMagnitudo, Description: "Diminishes the overall power, making the spell subtle or highly efficient to cast.", Tier: 1},
		{Name: "Extreme", Symbol: "Xt", Category: models.CategoryMagnitudo, Description: "Pushes the magic to its absolute dimensional limit. Highly volatile.", Tier: 2},
		{Name: "Inverse", Symbol: "Iv", Category: models.CategoryMagnitudo, Description: "Inverts the fundamental nature of the paired Essentia or Actio (e.g. fire→cold, push→pull, create→destroy).", Tier: 1},
	}
}
