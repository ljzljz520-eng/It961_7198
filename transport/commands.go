package transport

import (
	"fmt"
	"strings"

	"drivingmaterials/domain"
	"drivingmaterials/workflow"
)

type Command struct {
	Name string
	Args map[string]string
}

func ParseCommand(line string) (Command, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return Command{}, fmt.Errorf("empty command")
	}
	c := Command{Name: strings.ToLower(fields[0]), Args: map[string]string{}}
	for _, field := range fields[1:] {
		parts := strings.SplitN(field, "=", 2)
		if len(parts) != 2 {
			return Command{}, fmt.Errorf("argument %q must use key=value", field)
		}
		c.Args[strings.ToLower(parts[0])] = parts[1]
	}
	return c, nil
}

func BuildIngest(c Command) (workflow.IngestRequest, error) {
	required := []string{"id", "title", "kind", "subject", "campus", "version", "uri", "creator"}
	for _, key := range required {
		if c.Args[key] == "" {
			return workflow.IngestRequest{}, fmt.Errorf("missing %s", key)
		}
	}
	kind := domain.MaterialKind(c.Args["kind"])
	return workflow.IngestRequest{ID: c.Args["id"], Title: c.Args["title"], Kind: kind, Subject: c.Args["subject"], Campus: c.Args["campus"], VersionDate: c.Args["version"], URI: c.Args["uri"], Description: c.Args["description"], Creator: c.Args["creator"]}, nil
}
