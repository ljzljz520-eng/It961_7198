package workflow

import (
	"path/filepath"
	"testing"

	"drivingmaterials/catalog"
	"drivingmaterials/domain"
	"drivingmaterials/persistence"
)

func TestIngestNormalizesAndPublishes(t *testing.T) {
	store, err := persistence.Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	c := catalog.NewService(store, catalog.StaticClock{Date: "2026-04-01"})
	s := NewIngestService(c, store)
	m, err := s.Ingest(IngestRequest{ID: "m1", Title: "Safety", Kind: domain.SafetyVideo, Subject: "one", Campus: "north", VersionDate: "2026-04-01", URI: "video", Tags: []string{" Safety ", "intro"}, Creator: "coach"})
	if err != nil {
		t.Fatal(err)
	}
	if !m.HasTag("safety") {
		t.Fatal("safety policy missing")
	}
	if _, err := s.Publish(m.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
}
