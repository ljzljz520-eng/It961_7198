package persistence

import (
	"path/filepath"
	"testing"

	"drivingmaterials/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "materials.db")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	m := domain.NewMaterial("persist-1", "Route A", domain.RouteMap, "two", "central", "2026-02-01", "route-a", "coach")
	m.Status = domain.StatusPublished
	if err := store.PutMaterial(m); err != nil {
		t.Fatal(err)
	}
	if err := store.PutCourse(domain.Course{Code: "C1", Name: "Road", Subject: "two", Campus: "central", Instructor: "Lee", Active: true}); err != nil {
		t.Fatal(err)
	}
	if err := store.PutReview(domain.Review{MaterialID: m.ID, Reviewer: "Lee", Decision: "approve", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if err := store.AppendAudit(AuditEvent{ID: "a1", Action: "seed", MaterialID: m.ID, Actor: "Lee", Detail: "seed", CreatedAt: "2026-02-01"}); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	got, err := reopened.GetMaterial(m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != m.Title || got.Status != domain.StatusPublished {
		t.Fatalf("reopened material mismatch: %#v", got)
	}
	counts, err := reopened.Counts()
	if err != nil {
		t.Fatal(err)
	}
	if counts.Materials != 1 || counts.Courses != 1 || counts.Reviews != 1 || counts.Audits != 1 {
		t.Fatalf("unexpected counts: %#v", counts)
	}
}

func TestStoreListsAndExports(t *testing.T) {
	path := filepath.Join(t.TempDir(), "materials.db")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	for i := 0; i < 2; i++ {
		m := domain.NewMaterial(string(rune('a'+i)), "Title", domain.CourseMaterial, "one", "north", "2026-01-01", "uri", "coach")
		if err := store.PutMaterial(m); err != nil {
			t.Fatal(err)
		}
	}
	values, err := store.ListMaterials()
	if err != nil || len(values) != 2 {
		t.Fatalf("list failed: %v %d", err, len(values))
	}
	if err := store.Export(filepath.Join(t.TempDir(), "copy.db")); err != nil {
		t.Fatal(err)
	}
}
