package transport

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"drivingmaterials/catalog"
	"drivingmaterials/persistence"
	"drivingmaterials/query"
	"drivingmaterials/workflow"
)

func TestCLIAddAndSearch(t *testing.T) {
	store, err := persistence.Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	c := catalog.NewService(store, catalog.StaticClock{Date: "2026-04-01"})
	in := workflow.NewIngestService(c, store)
	ret := workflow.NewRetrieveService(c, query.NewSearcher(store), store)
	var out bytes.Buffer
	cli := CLI{In: strings.NewReader("add id=m1 title=Signs kind=course subject=one campus=north version=2026-01-01 uri=signs creator=coach\nsearch subject=one campus=north\nexit\n"), Out: &out, Ingest: in, Retrieve: ret}
	if err := cli.Run(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Signs") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}
