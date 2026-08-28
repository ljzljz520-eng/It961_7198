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

func setup(t *testing.T) (*persistence.Store, *catalog.Service, *workflow.IngestService, *workflow.RetrieveService, *workflow.ReviewService) {
	t.Helper()
	store, err := persistence.Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	c := catalog.NewService(store, catalog.StaticClock{Date: "2026-05-01"})
	return store, c, workflow.NewIngestService(c, store), workflow.NewRetrieveService(c, query.NewSearcher(store), store), workflow.NewReviewService(c, store)
}

func TestWorkflowMaterialIngestion(t *testing.T) {
	store, _, in, ret, _ := setup(t)
	defer store.Close()
	if _, err := in.Ingest(IngestRequest{ID: "m1", Title: "Course one", Kind: domain.CourseMaterial, Subject: "one", Campus: "north", VersionDate: "2026-05-01", URI: "course", Creator: "coach"}); err != nil {
		t.Fatal(err)
	}
	report, err := ret.Retrieve(domain.MaterialFilter{Subject: "one", Campus: "north"})
	if err != nil {
		t.Fatal(err)
	}
	if report.Total != 1 {
		t.Fatalf("expected one result, got %d", report.Total)
	}
}

func TestWorkflowCrossCampusSearch(t *testing.T) {
	store, _, in, ret, _ := setup(t)
	defer store.Close()
	for _, args := range []IngestRequest{{ID: "m1", Title: "North lights", Kind: domain.LightSimulation, Subject: "three", Campus: "north", VersionDate: "2026-05-01", URI: "n", Creator: "coach"}, {ID: "m2", Title: "South route", Kind: domain.RouteMap, Subject: "three", Campus: "south", VersionDate: "2026-05-01", URI: "s", Creator: "coach"}} {
		if _, err := in.Ingest(args); err != nil {
			t.Fatal(err)
		}
	}
	report, err := ret.Retrieve(domain.MaterialFilter{Campus: "south"})
	if err != nil {
		t.Fatal(err)
	}
	if report.Total != 1 || report.Results[0].Material.Campus != "south" {
		t.Fatalf("wrong campus results: %#v", report.Results)
	}
}

func TestWorkflowStatusReview(t *testing.T) {
	store, cat, in, _, review := setup(t)
	defer store.Close()
	m, err := in.Ingest(IngestRequest{ID: "m1", Title: "Safety", Kind: domain.SafetyVideo, Subject: "four", Campus: "east", VersionDate: "2026-05-01", URI: "safety", Creator: "coach"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cat.SetStatus(m.ID, domain.StatusReview, "coach"); err != nil {
		t.Fatal(err)
	}
	published, err := review.Submit(domain.Review{MaterialID: m.ID, Reviewer: "reviewer", Decision: "approve", Version: 1})
	if err != nil {
		t.Fatal(err)
	}
	if published.Status != domain.StatusPublished {
		t.Fatalf("status=%s", published.Status)
	}
}

type IngestRequest = workflow.IngestRequest
