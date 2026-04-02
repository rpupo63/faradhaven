package models

import (
	"testing"

	"github.com/google/uuid"
)

func mkLink(order int, id uuid.UUID, cat ComponentCategory) SpellComponent {
	return SpellComponent{
		SortOrder:   order,
		ComponentID: id,
		Component:   Component{ID: id, Category: cat},
	}
}

func TestComponentFingerprint_noLogica_orderIndependentWithinSinglePhase(t *testing.T) {
	a := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	b := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	c := uuid.MustParse("33333333-3333-3333-3333-333333333333")

	// A, B, C
	seq1 := []SpellComponent{
		mkLink(0, a, CategoryForma),
		mkLink(1, b, CategoryEssentia),
		mkLink(2, c, CategoryActio),
	}
	// C, A, B — same multiset
	seq2 := []SpellComponent{
		mkLink(0, c, CategoryActio),
		mkLink(1, a, CategoryForma),
		mkLink(2, b, CategoryEssentia),
	}

	fp1 := ComponentFingerprint(seq1)
	fp2 := ComponentFingerprint(seq2)
	if fp1 != fp2 {
		t.Fatalf("same multiset should match: %s vs %s", fp1, fp2)
	}
}

func TestComponentFingerprint_noLogica_duplicatesPreserved(t *testing.T) {
	a := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	b := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")

	// A, A, B
	seq1 := []SpellComponent{
		mkLink(0, a, CategoryForma),
		mkLink(1, a, CategoryForma),
		mkLink(2, b, CategoryEssentia),
	}
	// A, B, A
	seq2 := []SpellComponent{
		mkLink(0, a, CategoryForma),
		mkLink(1, b, CategoryEssentia),
		mkLink(2, a, CategoryForma),
	}

	if ComponentFingerprint(seq1) != ComponentFingerprint(seq2) {
		t.Fatal("duplicate counts should match after canonicalization")
	}
}

func TestComponentFingerprint_logica_separatesPhasesMultisetPerPhase(t *testing.T) {
	a := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	b := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	c := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	logica := uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd")

	// Phase1: A,B — If — Phase2: C  vs  Phase1: B,A — If — Phase2: C
	seq1 := []SpellComponent{
		mkLink(0, a, CategoryForma),
		mkLink(1, b, CategoryEssentia),
		mkLink(2, logica, CategoryLogica),
		mkLink(3, c, CategoryActio),
	}
	seq2 := []SpellComponent{
		mkLink(0, b, CategoryEssentia),
		mkLink(1, a, CategoryForma),
		mkLink(2, logica, CategoryLogica),
		mkLink(3, c, CategoryActio),
	}

	if ComponentFingerprint(seq1) != ComponentFingerprint(seq2) {
		t.Fatal("order within phase 1 should not change fingerprint")
	}
}

func TestComponentFingerprint_logica_orderOfLinksMatters(t *testing.T) {
	a := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	if1 := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	if2 := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	// A — If1 — A — If2
	seq1 := []SpellComponent{
		mkLink(0, a, CategoryForma),
		mkLink(1, if1, CategoryLogica),
		mkLink(2, a, CategoryForma),
		mkLink(3, if2, CategoryLogica),
	}
	// A — If2 — A — If1 (swap connectors)
	seq2 := []SpellComponent{
		mkLink(0, a, CategoryForma),
		mkLink(1, if2, CategoryLogica),
		mkLink(2, a, CategoryForma),
		mkLink(3, if1, CategoryLogica),
	}

	if ComponentFingerprint(seq1) == ComponentFingerprint(seq2) {
		t.Fatal("Logica order should change fingerprint")
	}
}

func TestComponentFingerprint_logicaLeadingPhase(t *testing.T) {
	a := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	logica := uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd")

	seq := []SpellComponent{
		mkLink(0, logica, CategoryLogica),
		mkLink(1, a, CategoryForma),
	}
	fp := ComponentFingerprint(seq)
	if fp == "" {
		t.Fatal("expected non-empty fingerprint")
	}
}
