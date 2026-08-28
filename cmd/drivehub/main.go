package main

import (
	"flag"
	"fmt"
	"os"

	"drivingmaterials/catalog"
	"drivingmaterials/persistence"
	"drivingmaterials/query"
	"drivingmaterials/transport"
	"drivingmaterials/workflow"
)

func main() {
	path := flag.String("db", "drivehub.db", "bbolt database path")
	flag.Parse()
	store, err := persistence.Open(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer store.Close()
	serviceCatalog := catalog.NewService(store, catalog.StaticClock{Date: "2026-01-01"})
	ingest := workflow.NewIngestService(serviceCatalog, store)
	retrieve := workflow.NewRetrieveService(serviceCatalog, query.NewSearcher(store), store)
	cli := transport.CLI{In: os.Stdin, Out: os.Stdout, Ingest: ingest, Retrieve: retrieve}
	if err := cli.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
