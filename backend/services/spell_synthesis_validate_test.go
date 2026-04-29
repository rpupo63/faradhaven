package services

import (
	"testing"

	"github.com/rpupo63/faradhaven/backend/models"
)

func testComp(name string, category models.ComponentCategory) models.Component {
	return models.Component{Name: name, Category: category}
}

func hasErrContaining(errs []string, want string) bool {
	for _, err := range errs {
		if err == want {
			return true
		}
	}
	return false
}

func TestSpellSynthesisValidate_PhaseAwareFormaScopusRules(t *testing.T) {
	svc := &SpellSynthesisService{}

	tests := []struct {
		name         string
		components   []models.Component
		wantNoErrors bool
		wantErr      string
	}{
		{
			name: "single phase requires scopus",
			components: []models.Component{
				testComp("Projectile", models.CategoryForma),
				testComp("Ignis", models.CategoryEssentia),
			},
			wantNoErrors: false,
			wantErr:      "Phase 1 requires exactly 1 Scopus (targeting) component",
		},
		{
			name: "single phase passes with one forma and one scopus",
			components: []models.Component{
				testComp("Projectile", models.CategoryForma),
				testComp("Target", models.CategoryScopus),
				testComp("Ignis", models.CategoryEssentia),
			},
			wantNoErrors: true,
		},
		{
			name: "logica fails when phase 2 has no forma/scopus",
			components: []models.Component{
				testComp("Projectile", models.CategoryForma),
				testComp("Target", models.CategoryScopus),
				testComp("Ignis", models.CategoryEssentia),
				testComp("Then", models.CategoryLogica),
				testComp("Push", models.CategoryActio),
				testComp("Aqua", models.CategoryEssentia),
			},
			wantNoErrors: false,
			wantErr:      "Phase 2 requires exactly 1 Forma (shape) component",
		},
		{
			name: "logica fails when forma appears after logica without second-phase scopus",
			components: []models.Component{
				testComp("Projectile", models.CategoryForma),
				testComp("Target", models.CategoryScopus),
				testComp("Ignis", models.CategoryEssentia),
				testComp("Then", models.CategoryLogica),
				testComp("Nova", models.CategoryForma),
				testComp("Aqua", models.CategoryEssentia),
			},
			wantNoErrors: false,
			wantErr:      "Phase 2 requires exactly 1 Scopus (targeting) component",
		},
		{
			name: "logica per-phase passes with one forma and one scopus in each phase",
			components: []models.Component{
				testComp("Projectile", models.CategoryForma),
				testComp("Target", models.CategoryScopus),
				testComp("Ignis", models.CategoryEssentia),
				testComp("Then", models.CategoryLogica),
				testComp("Nova", models.CategoryForma),
				testComp("Ground", models.CategoryScopus),
				testComp("Aqua", models.CategoryEssentia),
			},
			wantNoErrors: true,
		},
		{
			name: "logica fails when second phase lacks scopus",
			components: []models.Component{
				testComp("Projectile", models.CategoryForma),
				testComp("Target", models.CategoryScopus),
				testComp("Ignis", models.CategoryEssentia),
				testComp("Then", models.CategoryLogica),
				testComp("Nova", models.CategoryForma),
				testComp("Aqua", models.CategoryEssentia),
			},
			wantNoErrors: false,
			wantErr:      "Phase 2 requires exactly 1 Scopus (targeting) component",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := svc.Validate(tt.components)
			if tt.wantNoErrors && len(errs) > 0 {
				t.Fatalf("expected no validation errors, got: %v", errs)
			}
			if !tt.wantNoErrors {
				if len(errs) == 0 {
					t.Fatalf("expected validation errors, got none")
				}
				if tt.wantErr != "" && !hasErrContaining(errs, tt.wantErr) {
					t.Fatalf("expected error %q, got: %v", tt.wantErr, errs)
				}
			}
		})
	}
}
