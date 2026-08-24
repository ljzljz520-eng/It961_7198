package integration

import (
	"path/filepath"
	"testing"

	"drivingmaterials/catalog"
	"drivingmaterials/domain"
	"drivingmaterials/persistence"
	"drivingmaterials/query"
	"drivingmaterials/workflow"
)

func TestCourseSearchResultsStayIndependent(t *testing.T) {
	store, _, in, _, _ := setup(t)
	defer store.Close()
	_, err := in.Ingest(workflow.IngestRequest{ID: "course-1", Title: "Signs", Kind: domain.CourseMaterial, Subject: "one", Campus: "north", VersionDate: "2026-05-01", URI: "course", Creator: "coach"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = in.Ingest(workflow.IngestRequest{ID: "light-1", Title: "Lights", Kind: domain.LightSimulation, Subject: "three", Campus: "south", VersionDate: "2026-06-01", URI: "lights", Creator: "coach"})
	if err != nil {
		t.Fatal(err)
	}
	searcher := query.NewSearcher(store)
	first, err := searcher.Search(domain.MaterialFilter{Kind: domain.CourseMaterial})
	if err != nil {
		t.Fatal(err)
	}
	second, err := searcher.Search(domain.MaterialFilter{Kind: domain.LightSimulation})
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 1 || len(second) != 1 {
		t.Fatalf("unexpected result sizes: %d %d", len(first), len(second))
	}
	if first[0].Labels[0] != string(domain.CourseMaterial) || first[0].Labels[2] != "north" {
		t.Fatalf("first result labels changed: %#v", first[0].Labels)
	}
	if second[0].Labels[0] != string(domain.LightSimulation) || second[0].Labels[2] != "south" {
		t.Fatalf("second result labels wrong: %#v", second[0].Labels)
	}
}

func TestWorkflowReopenSearchesPersistedMaterials(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	store, err := persistence.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	c := catalog.NewService(store, catalog.StaticClock{Date: "2026-05-01"})
	in := workflow.NewIngestService(c, store)
	if _, err := in.Ingest(workflow.IngestRequest{ID: "m1", Title: "Route", Kind: domain.RouteMap, Subject: "two", Campus: "north", VersionDate: "2026-05-01", URI: "route", Creator: "coach"}); err != nil {
		t.Fatal(err)
	}
	store.Close()
	reopened, err := persistence.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	results, err := query.NewSearcher(reopened).Search(domain.MaterialFilter{Kind: domain.RouteMap})
	if err != nil || len(results) != 1 {
		t.Fatalf("reopen search failed: %v %d", err, len(results))
	}
}
