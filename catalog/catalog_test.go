package catalog

import (
	"path/filepath"
	"testing"

	"drivingmaterials/domain"
	"drivingmaterials/persistence"
)

func TestCatalogStatusTransitions(t *testing.T) {
	store, err := persistence.Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	service := NewService(store, StaticClock{Date: "2026-03-01"})
	m := domain.NewMaterial("m1", "Lights", domain.LightSimulation, "three", "west", "2026-03-01", "lights", "coach")
	if err := service.RegisterMaterial(m); err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetStatus("m1", domain.StatusPublished, "coach"); err == nil {
		t.Fatal("draft should not skip review")
	}
	if _, err := service.SetStatus("m1", domain.StatusReview, "coach"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetStatus("m1", domain.StatusPublished, "reviewer"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetStatus("m1", domain.StatusArchived, "coach"); err != nil {
		t.Fatal(err)
	}
}
