package models

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

// ComponentFingerprint computes a canonical fingerprint for a spell's component chain.
//
// Rules (aligned with the forge / prepared-spell matching):
//   - CategoryLogica components (If, Then, Therefore) are phase separators; their order
//     and identity are preserved in the serialized form.
//   - Between separators, non-Logica components form buckets: order within a phase does
//     not matter. Canonical form is lexicographically sorted UUID strings joined by
//     commas, with duplicate IDs preserved (multiset).
//
// Requires each SpellComponent.Component.Category to be set (e.g. spell_repo
// computeAndStoreFingerprint loads categories; DB preloads must include Component).
func ComponentFingerprint(links []SpellComponent) string {
	ordered := make([]SpellComponent, len(links))
	copy(ordered, links)
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].SortOrder < ordered[j].SortOrder
	})

	var sb strings.Builder
	for i := 0; i < len(ordered); {
		var phase []string
		for i < len(ordered) && ordered[i].Component.Category != CategoryLogica {
			phase = append(phase, ordered[i].ComponentID.String())
			i++
		}
		sb.WriteString(canonicalPhaseMultiset(phase))
		if i < len(ordered) && ordered[i].Component.Category == CategoryLogica {
			sb.WriteByte('|')
			sb.WriteString(ordered[i].ComponentID.String())
			sb.WriteByte('|')
			i++
		}
	}

	h := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(h[:])
}

// canonicalPhaseMultiset encodes one phase’s multiset of component UUIDs in an
// order-independent way (sorted IDs, duplicates kept).
func canonicalPhaseMultiset(phase []string) string {
	if len(phase) == 0 {
		return ""
	}
	cp := append([]string(nil), phase...)
	sort.Strings(cp)
	return strings.Join(cp, ",")
}
