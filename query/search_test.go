package query

import (
	"path/filepath"
	"testing"

	"drivingmaterials/domain"
	"drivingmaterials/persistence"
)

func TestSearchRanksAndGroups(t *testing.T) {
	store, err := persistence.Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	for _, m := range []domain.Material{
		domain.NewMaterial("m1", "Signs", domain.CourseMaterial, "one", "north", "2026-01-01", "u1", "c"),
		domain.NewMaterial("m2", "Lights", domain.LightSimulation, "one", "north", "2026-01-02", "u2", "c"),
	} {
		m.Status = domain.StatusPublished
		if err := store.PutMaterial(m); err != nil {
			t.Fatal(err)
		}
	}
	results, err := NewSearcher(store).Search(domain.MaterialFilter{Subject: "one", Campus: "north"})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results", len(results))
	}
	if len(GroupByKind(results)) != 2 {
		t.Fatal("expected two groups")
	}
	if results[0].Score < results[1].Score {
		t.Fatal("results are not ranked")
	}
}
