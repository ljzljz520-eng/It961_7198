package workflow

import (
	"path/filepath"
	"testing"

	"drivingmaterials/catalog"
	"drivingmaterials/domain"
	"drivingmaterials/persistence"
	"drivingmaterials/query"
)

func TestRetrieveBuildsSnapshot(t *testing.T) {
	store, err := persistence.Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	c := catalog.NewService(store, catalog.StaticClock{Date: "2026-04-01"})
	s := NewIngestService(c, store)
	_, err = s.Ingest(IngestRequest{ID: "m1", Title: "Route", Kind: domain.RouteMap, Subject: "two", Campus: "central", VersionDate: "2026-04-01", URI: "route", Creator: "coach"})
	if err != nil {
		t.Fatal(err)
	}
	r := NewRetrieveService(c, query.NewSearcher(store), store)
	report, err := r.Retrieve(domain.MaterialFilter{Campus: "central"})
	if err != nil {
		t.Fatal(err)
	}
	if report.Total != 1 || report.Snapshot.Total() != 1 {
		t.Fatalf("unexpected report: %#v", report)
	}
}
