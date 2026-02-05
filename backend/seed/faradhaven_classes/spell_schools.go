package faradhaven_classes

import (
	"github.com/rpupo63/unified-personal-site-backend/models"
)

func schoolComponents() []ComponentSeed {
	return []ComponentSeed{
		// VII. ABJURATION
		{
			Name:        "Reflect",
			Symbol:      "Rf",
			Category:    models.CategoryAbjuration,
			Description: "Bounces incoming projectiles back at the source (mirror)",
			Element:     "",
		},
		{
			Name:        "Dispel",
			Symbol:      "Dp",
			Category:    models.CategoryAbjuration,
			Description: "Removes active status effects or destroys enemy traps/zones (nullify)",
			Element:     "",
		},
		{
			Name:        "Absorb",
			Symbol:      "Ab",
			Category:    models.CategoryAbjuration,
			Description: "Converts incoming damage into Mana for the player",
			Element:     "",
		},
		{
			Name:        "Seal",
			Symbol:      "Sl",
			Category:    models.CategoryAbjuration,
			Description: "Permanently closes a door, chest, or portal (lock)",
			Element:     "",
		},

		// VIII. CONJURATION
		{
			Name:        "Summon",
			Symbol:      "Su",
			Category:    models.CategoryConjuration,
			Description: "Spawns a friendly entity (AI companion, minion)",
			Element:     "",
		},
		{
			Name:        "Portal",
			Symbol:      "Po",
			Category:    models.CategoryConjuration,
			Description: "Creates two linked points; walking in one exits the other (gate)",
			Element:     "",
		},
		{
			Name:        "Grease",
			Symbol:      "Gr",
			Category:    models.CategoryConjuration,
			Description: "Covers the floor in a zero-friction substance (slick, flammable based on other components)",
			Element:     "",
		},

		// IX. DIVINATION
		{
			Name:        "Sight",
			Symbol:      "Si",
			Category:    models.CategoryDivination,
			Description: "Allows seeing enemies/loot through walls (x-ray)",
			Element:     "",
		},
		{
			Name:        "Identify",
			Symbol:      "Id",
			Category:    models.CategoryDivination,
			Description: "Reveals the HP, weaknesses, and stats of a target (lore)",
			Element:     "",
		},
		{
			Name:        "Predict",
			Symbol:      "Pd",
			Category:    models.CategoryDivination,
			Description: "Shows the path of enemy projectiles before they fire (trajectory)",
			Element:     "",
		},
		{
			Name:        "Link",
			Symbol:      "Lk",
			Category:    models.CategoryDivination,
			Description: "Connects two targets; damage dealt to one is felt by the other (share)",
			Element:     "",
		},

		// X. ENCHANTMENT
		{
			Name:        "Fear",
			Symbol:      "Fr",
			Category:    models.CategoryEnchantment,
			Description: "Forces enemies to run directly away from the caster (panic)",
			Element:     "",
		},
		{
			Name:        "Taunt",
			Symbol:      "Ta",
			Category:    models.CategoryEnchantment,
			Description: "Forces enemies to attack a specific target (aggro, good for decoys)",
			Element:     "",
		},
		{
			Name:        "Frenzy",
			Symbol:      "Fz",
			Category:    models.CategoryEnchantment,
			Description: "Increases enemy damage but causes them to attack their own allies (berserk)",
			Element:     "",
		},
		{
			Name:        "Command",
			Symbol:      "Cm",
			Category:    models.CategoryEnchantment,
			Description: "A generic trigger component that forces actions (obey, e.g., Command + Jump)",
			Element:     "",
		},

		// XI. EVOCATION
		{
			Name:        "Channel",
			Symbol:      "Ch",
			Category:    models.CategoryEvocation,
			Description: "A continuous flow of energy like a flamethrower that drains mana over time (stream)",
			Element:     "",
		},
		{
			Name:        "Chain",
			Symbol:      "Ci",
			Category:    models.CategoryEvocation,
			Description: "On hit, the spell jumps to a nearby enemy (arc)",
			Element:     "",
		},
		{
			Name:        "Imbue",
			Symbol:      "Im",
			Category:    models.CategoryEvocation,
			Description: "Instead of firing a spell, adds the element to your sword/melee attack (enchant weapon)",
			Element:     "",
		},
		{
			Name:        "Rupture",
			Symbol:      "Rp",
			Category:    models.CategoryEvocation,
			Description: "A delayed explosion that triggers seconds after the initial hit (aftershock)",
			Element:     "",
		},

		// XII. ILLUSION
		{
			Name:        "Decoy",
			Symbol:      "Dc",
			Category:    models.CategoryIllusion,
			Description: "Creates a fake copy of the player that attracts attacks but deals no damage (clone)",
			Element:     "",
		},
		{
			Name:        "Invisible",
			Symbol:      "Iv",
			Category:    models.CategoryIllusion,
			Description: "Makes the target unseeable by AI (stealth)",
			Element:     "",
		},
		{
			Name:        "Silence",
			Symbol:      "Si",
			Category:    models.CategoryIllusion,
			Description: "Prevents the target from casting spells or calling for help (mute)",
			Element:     "",
		},
		{
			Name:        "Disguise",
			Symbol:      "Dg",
			Category:    models.CategoryIllusion,
			Description: "Makes the player look like an enemy faction (mask, prevents aggro)",
			Element:     "",
		},
		{
			Name:        "Phantom",
			Symbol:      "Ph",
			Category:    models.CategoryIllusion,
			Description: "Creates a visual spell effect that does no damage but triggers Fear or Dodge AI (fake)",
			Element:     "",
		},

		// XIII. NECROMANCY
		{
			Name:        "Bone",
			Symbol:      "Bn",
			Category:    models.CategoryNecromancy,
			Description: "Uses physical remains to build walls or armor (structure)",
			Element:     "",
		},
		{
			Name:        "Soul",
			Symbol:      "So",
			Category:    models.CategoryNecromancy,
			Description: "A component required for targeting ghosts or interacting with the dead (spirit)",
			Element:     "",
		},
		{
			Name:        "Curse",
			Symbol:      "Cu",
			Category:    models.CategoryNecromancy,
			Description: "A debuff that lowers enemy stats (Attack/Defense) without doing damage (hex)",
			Element:     "",
		},
		{
			Name:        "Sacrifice",
			Symbol:      "Sf",
			Category:    models.CategoryNecromancy,
			Description: "Casts a spell using HP instead of Mana (blood)",
			Element:     "",
		},

		// XIV. TRANSMUTATION
		{
			Name:        "Transmute",
			Symbol:      "Tm",
			Category:    models.CategoryTransmutation,
			Description: "The 'Make X' command. Turn Flesh to Stone, Wood to Gold (change)",
			Element:     "",
		},
		{
			Name:        "Burrow",
			Symbol:      "Bu",
			Category:    models.CategoryTransmutation,
			Description: "Allows movement through the ground/terrain (phase)",
			Element:     "",
		},
		{
			Name:        "Fuse",
			Symbol:      "Fs",
			Category:    models.CategoryTransmutation,
			Description: "Permanently attaches two objects together physically (combine)",
			Element:     "",
		},
	}
}
