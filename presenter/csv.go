package presenter

import (
	"encoding/csv"
	"io"
	"strconv"

	"drivingmaterials/query"
)

func WriteCSV(w io.Writer, results []query.SearchResult) error {
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"id", "title", "kind", "subject", "campus", "version", "score", "labels"}); err != nil {
		return err
	}
	for _, result := range results {
		row := []string{result.Material.ID, result.Material.Title, string(result.Material.Kind), result.Material.Subject, result.Material.Campus, result.Material.VersionDate, strconv.Itoa(result.Score), joinLabels(result.Labels)}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func joinLabels(labels []string) string {
	value := ""
	for i, label := range labels {
		if i > 0 {
			value += ","
		}
		value += label
	}
	return value
}

func WriteIDs(w io.Writer, results []query.SearchResult) error {
	for _, result := range results {
		if _, err := io.WriteString(w, result.Material.ID+"\n"); err != nil {
			return err
		}
	}
	return nil
}
