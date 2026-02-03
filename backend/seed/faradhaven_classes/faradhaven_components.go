package faradhaven_classes

import (
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// ComponentClassMapping pairs a component with the class that has access to it.
// Supports many-to-many: a component could appear in multiple ClassNames in the future.
type ComponentClassMapping struct {
	Component ComponentSeed
	ClassName string
}

// AllComponentClassMappings returns all spell components and their class associations.
// Components are the building blocks classes use (e.g., Aether Jack uses Ether, Scry, Phase, Weave).
// ClassComponent links are created from these mappings.
func AllComponentClassMappings() []ComponentClassMapping {
	return []ComponentClassMapping{
		// The Rift Weaver (formerly Aether Jack / Ether-Mage)
		{ComponentSeed{Name: "Ether", Type: models.ComponentTypeBase, Description: "Force damage / ethereal touch", Element: "force"}, "The Rift Weaver"},
		{ComponentSeed{Name: "Scry", Type: models.ComponentTypeBase, Description: "Divination / detection", Element: "divination"}, "The Rift Weaver"},
		{ComponentSeed{Name: "Phase", Type: models.ComponentTypeModifier, Description: "Teleportation / pass through barriers", Element: "transmutation"}, "The Rift Weaver"},
		{ComponentSeed{Name: "Weave", Type: models.ComponentTypeModifier, Description: "Enchantment / buff or debuff", Element: "enchantment"}, "The Rift Weaver"},

		// The Mutagen (formerly Chimera / Were-beast)
		{ComponentSeed{Name: "Feral Muscle", Type: models.ComponentTypeBase, Description: "Temp STR buff", Element: "nature"}, "The Mutagen"},
		{ComponentSeed{Name: "Chitin/Hide", Type: models.ComponentTypeBase, Description: "AC buff", Element: "shield"}, "The Mutagen"},
		{ComponentSeed{Name: "Adrenaline", Type: models.ComponentTypeModifier, Description: "Speed", Element: "nature"}, "The Mutagen"},
		{ComponentSeed{Name: "Regenerate", Type: models.ComponentTypeModifier, Description: "Heal over time", Element: "heal"}, "The Mutagen"},

		// The Ironwright (formerly Cog-Weaver / Tinkerer)
		{ComponentSeed{Name: "Spark", Type: models.ComponentTypeBase, Description: "Lightning damage", Element: "lightning"}, "The Ironwright"},
		{ComponentSeed{Name: "Steam", Type: models.ComponentTypeBase, Description: "Obscurement/Fog", Element: "nature"}, "The Ironwright"},
		{ComponentSeed{Name: "Clockwork Trigger", Type: models.ComponentTypeModifier, Description: "Delayed effect", Element: ""}, "The Ironwright"},
		{ComponentSeed{Name: "Construct", Type: models.ComponentTypeModifier, Description: "Manifests a small walking turret", Element: ""}, "The Ironwright"},

		// The Powder Mage (formerly Fusilier)
		{ComponentSeed{Name: "Ignite", Type: models.ComponentTypeBase, Description: "Fire damage", Element: "fire"}, "The Powder Mage"},
		{ComponentSeed{Name: "Frost", Type: models.ComponentTypeBase, Description: "Slows movement", Element: "ice"}, "The Powder Mage"},
		{ComponentSeed{Name: "Pierce", Type: models.ComponentTypeModifier, Description: "Lines/goes through cover", Element: ""}, "The Powder Mage"},
		{ComponentSeed{Name: "Seeker", Type: models.ComponentTypeModifier, Description: "Ignores cover/homing", Element: ""}, "The Powder Mage"},

		// The Piston Brawler (formerly Vi / Hextech Brawler)
		{ComponentSeed{Name: "Kinetic", Type: models.ComponentTypeBase, Description: "Force damage/Knockback", Element: "push"}, "The Piston Brawler"},
		{ComponentSeed{Name: "Weight", Type: models.ComponentTypeBase, Description: "Heavy impact/Stun", Element: "push"}, "The Piston Brawler"},
		{ComponentSeed{Name: "Explosive", Type: models.ComponentTypeModifier, Description: "Splash damage on hit", Element: "fire"}, "The Piston Brawler"},
		{ComponentSeed{Name: "Sunder", Type: models.ComponentTypeModifier, Description: "Reduces enemy AC", Element: ""}, "The Piston Brawler"},

		// The Sanguinist (formerly Vitalist / Bloodborne)
		{ComponentSeed{Name: "Sanguine", Type: models.ComponentTypeBase, Description: "Necrotic damage/Life steal", Element: "dark"}, "The Sanguinist"},
		{ComponentSeed{Name: "Mend", Type: models.ComponentTypeBase, Description: "Healing", Element: "heal"}, "The Sanguinist"},
		{ComponentSeed{Name: "Contagion", Type: models.ComponentTypeModifier, Description: "Spreads to nearby targets", Element: "dark"}, "The Sanguinist"},
		{ComponentSeed{Name: "Boil", Type: models.ComponentTypeModifier, Description: "Adds Fire damage", Element: "fire"}, "The Sanguinist"},

		// The Vapor Blade (formerly Shadow-Stalker / Assassin)
		{ComponentSeed{Name: "Venom", Type: models.ComponentTypeBase, Description: "Poison damage over time", Element: "poison"}, "The Vapor Blade"},
		{ComponentSeed{Name: "Shadow", Type: models.ComponentTypeBase, Description: "Stealth/obscurement", Element: "dark"}, "The Vapor Blade"},
		{ComponentSeed{Name: "Silence", Type: models.ComponentTypeModifier, Description: "Muffles sound/suppresses verbal", Element: ""}, "The Vapor Blade"},
		{ComponentSeed{Name: "Lethal", Type: models.ComponentTypeModifier, Description: "Extra damage on crit/surprise", Element: ""}, "The Vapor Blade"},

		// The Lorewright (formerly Field Scholar / Inquisitive)
		{ComponentSeed{Name: "Decipher", Type: models.ComponentTypeBase, Description: "Decode ancient languages and extract information", Element: "divination"}, "The Lorewright"},
		{ComponentSeed{Name: "Weak Point", Type: models.ComponentTypeBase, Description: "Identify vulnerabilities in creatures and machinery", Element: "divination"}, "The Lorewright"},
		{ComponentSeed{Name: "Record", Type: models.ComponentTypeModifier, Description: "Preserve findings in field journal for the party", Element: ""}, "The Lorewright"},
		{ComponentSeed{Name: "Observe", Type: models.ComponentTypeModifier, Description: "Enhanced perception and environmental awareness", Element: "divination"}, "The Lorewright"},
	}
}
