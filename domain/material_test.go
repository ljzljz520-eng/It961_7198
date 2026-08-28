package domain

import "testing"

func TestMaterialValidationAndSearchText(t *testing.T) {
	m := NewMaterial("m1", "Signs", CourseMaterial, "subject-one", "north", "2026-01-01", "file://signs", "coach")
	m.Tags = []string{"priority"}
	if err := m.Validate(); err != nil {
		t.Fatal(err)
	}
	if !m.MatchesSubject("SUBJECT-ONE") || !m.MatchesCampus("NORTH") {
		t.Fatal("case-insensitive matching failed")
	}
	if !m.HasTag("PRIORITY") {
		t.Fatal("tag matching failed")
	}
	if m.SearchText() == "" {
		t.Fatal("search text is empty")
	}
}
