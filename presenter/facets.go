package presenter

import (
	"fmt"
	"io"

	"drivingmaterials/query"
)

func WriteFacets(w io.Writer, facets query.Facets) error {
	if _, err := fmt.Fprintln(w, "subjects:"); err != nil {
		return err
	}
	if err := writeFacetGroup(w, facets.Subjects); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "campuses:"); err != nil {
		return err
	}
	if err := writeFacetGroup(w, facets.Campuses); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "kinds:"); err != nil {
		return err
	}
	return writeFacetGroup(w, facets.Kinds)
}

func writeFacetGroup(w io.Writer, values []query.Facet) error {
	for _, facet := range values {
		if _, err := fmt.Fprintf(w, "- %s (%d)\n", facet.Value, facet.Count); err != nil {
			return err
		}
	}
	return nil
}

func FormatFacetCount(facet query.Facet) string {
	return fmt.Sprintf("%s: %d", facet.Value, facet.Count)
}
