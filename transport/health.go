package transport

import (
	"fmt"
	"io"

	"drivingmaterials/persistence"
)

func WriteHealth(w io.Writer, store *persistence.Store) error {
	counts, err := store.Counts()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "materials=%d courses=%d reviews=%d audits=%d\n", counts.Materials, counts.Courses, counts.Reviews, counts.Audits)
	return err
}
