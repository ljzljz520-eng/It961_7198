package presenter

import (
	"fmt"
	"io"
	"strings"

	"drivingmaterials/query"
)

func WriteTable(w io.Writer, results []query.SearchResult) error {
	if _, err := fmt.Fprintln(w, "ID\tTITLE\tKIND\tCAMPUS\tVERSION\tLABELS"); err != nil {
		return err
	}
	for _, r := range results {
		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", r.Material.ID, r.Material.Title, r.Material.Kind, r.Material.Campus, r.Material.VersionDate, strings.Join(r.Labels, ",")); err != nil {
			return err
		}
	}
	return nil
}

func FormatSummary(results []query.SearchResult) string {
	if len(results) == 0 {
		return "no materials"
	}
	parts := make([]string, 0, len(results))
	for _, r := range results {
		parts = append(parts, fmt.Sprintf("%s (%s)", r.Material.Title, r.Material.Campus))
	}
	return strings.Join(parts, "; ")
}
