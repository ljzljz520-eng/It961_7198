package presenter

import (
	"encoding/json"
	"io"

	"drivingmaterials/workflow"
)

func WriteReport(w io.Writer, report workflow.RetrievalReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func MarshalReport(report workflow.RetrievalReport) ([]byte, error) { return json.Marshal(report) }
