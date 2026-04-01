package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rpupo63/unified-personal-site-backend/internal/bootstrap"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/services"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Info: No .env file found: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set. Set it in .env or environment.")
	}

	db, err := bootstrap.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// All spells with checked = false (GM not yet approved), regardless of existing AI columns.
	var spells []*models.Spell
	err = db.Preload("ComponentLinks", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC") }).
		Preload("ComponentLinks.Component").
		Where("checked = false").
		Order("created_at ASC").
		Find(&spells).Error
	if err != nil {
		log.Fatalf("Failed to query spells: %v", err)
	}
	for _, sp := range spells {
		sp.HydrateComponentsFromLinks()
	}

	total := len(spells)
	fmt.Printf("Found %d spell(s) needing AI review\n", total)
	if total == 0 {
		fmt.Println("Nothing to do.")
		return
	}

	var llmClient services.LLMClient
	if apiKey := os.Getenv("GOOGLE_API_KEY"); apiKey != "" {
		llmClient = services.NewGeminiLLMClient(apiKey)
		fmt.Println("Using Gemini LLM client")
	} else {
		llmClient = services.NewMockLLMClient()
		fmt.Println("GOOGLE_API_KEY not set — using mock LLM client")
	}
	spellAIService := services.NewSpellAIService(llmClient)
	// Spell review text and JSON schema live in services (BuildSpellAIOpinionPrompt, SpellAIOpinionResponseSchema).
	// Always go through GetSpellOpinion so batch review matches the HTTP API.

	successCount := 0
	errorCount := 0

	for i, spell := range spells {
		fmt.Printf("[%d/%d] %s (id=%s, level=%d) ... ", i+1, total, spell.Name, spell.ID, spell.Level)

		opinion, raw, err := spellAIService.GetSpellOpinion(context.Background(), spell, spell.Components)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			errorCount++
			continue
		}

		// Only write AI columns — leave all other spell fields untouched.
		// Mechanics recommendations are normalized (SpellType, feet, validated duration, enums, dice pair).
		aiT, aiR, aiD := models.NormalizeSpellAIRecommendations(opinion.RecommendedType, opinion.RecommendedRange, opinion.RecommendedDuration)
		sa, dt, aiDC, aiDS := models.NormalizeSpellAIRecommendationsExtras(opinion.RecommendedSaveAttr, opinion.RecommendedDamageType, opinion.RecommendedDamageDiceCount, opinion.RecommendedDamageDieSize)
		updates := map[string]interface{}{
			"ai_description_opinion":       &opinion.DescriptionOpinion,
			"ai_damage_opinion":            &opinion.DamageOpinion,
			"ai_effect_opinion":            &opinion.EffectOpinion,
			"ai_overall_verdict":           &opinion.OverallVerdict,
			"ai_raw_output":                &raw,
			"ai_recommended_name":          opinion.RecommendedName,
			"ai_recommended_description":   opinion.RecommendedDescription,
			"ai_recommended_save_attr":     sa,
			"ai_recommended_damage_type":   dt,
			"ai_recommended_damage_dice_count": aiDC,
			"ai_recommended_damage_die_size":   aiDS,
			"updated_at":                   time.Now(),
		}
		if aiT != nil {
			updates["ai_recommended_type"] = string(*aiT)
		} else {
			updates["ai_recommended_type"] = nil
		}
		if aiR != nil {
			updates["ai_recommended_range"] = *aiR
		} else {
			updates["ai_recommended_range"] = nil
		}
		if aiD != nil {
			updates["ai_recommended_duration"] = *aiD
		} else {
			updates["ai_recommended_duration"] = nil
		}

		if err := db.Model(&models.Spell{}).Where("id = ?", spell.ID).Updates(updates).Error; err != nil {
			fmt.Printf("SAVE ERROR: %v\n", err)
			errorCount++
			continue
		}

		fmt.Println("OK")
		successCount++

		if i < total-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	fmt.Printf("\nDone. %d succeeded, %d failed (out of %d total)\n", successCount, errorCount, total)
}
