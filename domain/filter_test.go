package domain

import "testing"

func TestMaterialFilterAcceptsVisibleRecords(t *testing.T) {
	m := NewMaterial("m1", "Signs", CourseMaterial, "one", "north", "2026-01-01", "uri", "coach")
	m.Status = StatusPublished
	if !(MaterialFilter{Subject: "one", Campus: "north", Kind: CourseMaterial}.Accepts(m)) {
		t.Fatal("expected record to match")
	}
	m.Status = StatusArchived
	if (MaterialFilter{}).Accepts(m) {
		t.Fatal("archived record should be hidden")
	}
	if !(MaterialFilter{IncludeArchived: true}.Accepts(m)) {
		t.Fatal("explicit archive filter should match")
	}
}
