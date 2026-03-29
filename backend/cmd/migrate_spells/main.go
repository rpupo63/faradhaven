// Unified spell table migration: mechanics (type, range feet, duration, AI recommendations)
// plus damage dice split and save/damage-type normalization, then drops legacy text columns.
//
// Replaces separate migrate_spell_mechanics + migrate_spell_damage_fields for one compliant pass.
//
// Usage (from backend/):
//
//	go run ./cmd/migrate_spells              # dry-run: list planned row updates
//	go run ./cmd/migrate_spells -apply       # apply updates + drop legacy columns (when present)
//	go run ./cmd/migrate_spells -mechanics-only       # only mechanics columns (no dice schema/drop)
//	go run ./cmd/migrate_spells -damage-only -apply   # only dice + normalize + drop legacy
//	go run ./cmd/migrate_spells -report-parse-failures  # list rows where legacy text fails parsers (no writes)
//
// Requires DATABASE_URL. Optional safety before -apply: pg_dump spells or export CSV.
//
// Suggested order with other data migrations: migrate_spells first (legacy spell columns), then
// migrate_weapon_damages, migrate_creatures, migrate_map_tokens, migrate_class_resources (independent).
//
// If legacy text columns somehow remain after migrate_spells, use:
//	go run ./cmd/drop_spell_legacy_text_columns -apply
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/joho/godotenv"
	"github.com/rpupo63/unified-personal-site-backend/internal/bootstrap"
	"github.com/rpupo63/unified-personal-site-backend/internal/migrationspell"
)

func main() {
	apply := flag.Bool("apply", false, "write updates and drop legacy damage_dice* text columns (when damage phase runs)")
	mechanicsOnly := flag.Bool("mechanics-only", false, "only migrate type/range/duration and related AI columns")
	damageOnly := flag.Bool("damage-only", false, "only migrate dice ints + save/damage type + drop legacy text columns")
	reportParseFailures := flag.Bool("report-parse-failures", false, "print spell rows where legacy dice/save/damage text fails parsers (no writes)")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Printf("info: no .env in cwd: %v", err)
	}
	bootstrap.LoadEnv()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := bootstrap.InitDB(dsn)
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	if *reportParseFailures {
		issues, err := migrationspell.ScanLegacyDamageParseIssues(db)
		if err != nil {
			log.Fatalf("report: %v", err)
		}
		fmt.Printf("Legacy parse issues (spell rows): %d\n", len(issues))
		for _, it := range issues {
			fmt.Printf("[%s] %q field=%s raw=%q (%s)\n", it.SpellID, it.Name, it.Field, it.Raw, it.Detail)
		}
		return
	}

	opts := migrationspell.Options{
		Apply:         *apply,
		MechanicsOnly: *mechanicsOnly,
		DamageOnly:    *damageOnly,
	}

	changes, err := migrationspell.ComputeChanges(db, opts)
	if err != nil {
		log.Fatalf("plan: %v", err)
	}

	total, _ := migrationspell.CountSpells(db)
	fmt.Printf("Spells in database: %d\n", total)
	fmt.Printf("Spells with planned updates: %d\n", len(changes))
	if len(changes) > 0 && !*apply {
		fmt.Println("Dry-run (no writes). Pass -apply to persist.")
	}

	for _, c := range changes {
		fmt.Printf("[%s] %q\n", c.ID, c.Name)
		for _, k := range sortedKeys(c.Updates) {
			fmt.Printf("  %s: %v\n", k, c.Updates[k])
		}
	}

	if *apply {
		if err := migrationspell.ApplyPlanned(db, opts, changes); err != nil {
			log.Fatalf("apply: %v", err)
		}
		fmt.Println()
		if !*mechanicsOnly {
			fmt.Println("Applied updates. Legacy text columns dropped if they existed: damage_dice, suggested_damage_dice, ai_recommended_damage_dice")
		} else {
			fmt.Println("Applied mechanics-only updates.")
		}
	}

	fmt.Printf("\nDone. %d spell(s) with planned changes.\n", len(changes))
}

func sortedKeys(m map[string]interface{}) []string {
	pref := []string{
		"type", "range", "duration",
		"ai_recommended_type", "ai_recommended_range", "ai_recommended_duration",
		"damage_dice_count", "damage_die_size",
		"suggested_damage_dice_count", "suggested_damage_die_size",
		"ai_recommended_damage_dice_count", "ai_recommended_damage_die_size",
		"save_attr", "ai_recommended_save_attr", "ai_recommended_damage_type",
	}
	seen := map[string]bool{}
	var out []string
	for _, k := range pref {
		if _, ok := m[k]; ok {
			out = append(out, k)
			seen[k] = true
		}
	}
	var rest []string
	for k := range m {
		if !seen[k] {
			rest = append(rest, k)
		}
	}
	sort.Strings(rest)
	return append(out, rest...)
}
