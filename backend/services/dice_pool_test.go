package services

import (
	"testing"

	"github.com/rpupo63/faradhaven/backend/models"
)

func TestCalculateSpellEffect(t *testing.T) {
	tests := []struct {
		name       string
		level      int
		damageType models.DamageType
		formaName  string
		scopusName string
		magnitudes []string
		wantCount  int
		wantFaces  int
	}{
		{
			name:       "Standard Ranged Single Target (Rare)",
			level:      3,
			damageType: models.DamageThunder,
			formaName:  "Projectile",
			scopusName: "Target",
			magnitudes: []string{},
			wantCount:  2,
			wantFaces:  8,
		},
		{
			name:       "Touch Single Target (Rare)",
			level:      3,
			damageType: models.DamageThunder,
			formaName:  "Projectile", // Projectile + Self = effectively touch/melee delivery
			scopusName: "Self",
			magnitudes: []string{},
			wantCount:  2,
			wantFaces:  10, // 8 + 2 (Self)
		},
		{
			name:       "Touch Common Damage",
			level:      3,
			damageType: models.DamageFire,
			formaName:  "Projectile",
			scopusName: "Self",
			magnitudes: []string{},
			wantCount:  2,
			wantFaces:  12, // 8 + 2 (Self) + 4 (Fire) = 14, clamped to 12
		},
		{
			name:       "Nova AoE (Rare)",
			level:      3,
			damageType: models.DamageThunder,
			formaName:  "Nova",
			scopusName: "Target",
			magnitudes: []string{},
			wantCount:  2,
			wantFaces:  6, // 8 - 2 (Nova)
		},
		{
			name:       "Zone Persistent AoE (Rare)",
			level:      3,
			damageType: models.DamageThunder,
			formaName:  "Zone",
			scopusName: "Target",
			magnitudes: []string{},
			wantCount:  1, // 2 - 1 (Zone penalty)
			wantFaces:  4, // 8 - 4 (Zone penalty)
		},
		{
			name:       "High Level Zone (Common)",
			level:      7,
			damageType: models.DamageFire,
			formaName:  "Zone",
			scopusName: "Target",
			magnitudes: []string{},
			wantCount:  7, // 8 - 1
			wantFaces:  8, // 8 - 4 (Zone) + 4 (Fire)
		},
		{
			name:       "Lance + Marked boosts precision damage",
			level:      3,
			damageType: models.DamageThunder,
			formaName:  "Lance",
			scopusName: "Marked",
			magnitudes: []string{},
			wantCount:  2,
			wantFaces:  12, // 8 + 2 (Lance) + 2 (Marked)
		},
		{
			name:       "Orbit + Area-First lowers sustained area damage",
			level:      4,
			damageType: models.DamageCold,
			formaName:  "Orbit",
			scopusName: "Area-First",
			magnitudes: []string{},
			wantCount:  2, // 3 - 1
			wantFaces:  4, // 8 - 4 (Orbit) - 2 (Area-First) + 2 (Cold)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateSpellEffect(tt.level, tt.damageType, tt.formaName, tt.scopusName, tt.magnitudes)
			if got.Count != tt.wantCount {
				t.Errorf("CalculateSpellEffect() gotCount = %v, want %v", got.Count, tt.wantCount)
			}
			if got.Faces != tt.wantFaces {
				t.Errorf("CalculateSpellEffect() gotFaces = %v, want %v", got.Faces, tt.wantFaces)
			}
		})
	}
}
