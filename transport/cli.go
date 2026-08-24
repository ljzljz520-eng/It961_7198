package transport

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"drivingmaterials/domain"
	"drivingmaterials/presenter"
	"drivingmaterials/workflow"
)

type CLI struct {
	In       io.Reader
	Out      io.Writer
	Ingest   *workflow.IngestService
	Retrieve *workflow.RetrieveService
}

func (c *CLI) Run() error {
	scanner := bufio.NewScanner(c.In)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "quit" || line == "exit" {
			break
		}
		if err := c.execute(line); err != nil {
			fmt.Fprintln(c.Out, "error:", err)
		}
	}
	return scanner.Err()
}

func (c *CLI) execute(line string) error {
	cmd, err := ParseCommand(line)
	if err != nil {
		return err
	}
	switch cmd.Name {
	case "add":
		req, err := BuildIngest(cmd)
		if err != nil {
			return err
		}
		m, err := c.Ingest.Ingest(req)
		if err == nil {
			err = presenter.WriteMaterial(c.Out, m)
		}
		return err
	case "publish":
		m, err := c.Ingest.Publish(cmd.Args["id"], cmd.Args["actor"])
		if err == nil {
			err = presenter.WriteMaterial(c.Out, m)
		}
		return err
	case "search":
		filter := domain.MaterialFilter{Subject: cmd.Args["subject"], Campus: cmd.Args["campus"], VersionDate: cmd.Args["version"], Kind: domain.MaterialKind(cmd.Args["kind"]), Query: cmd.Args["q"]}
		report, err := c.Retrieve.Retrieve(filter)
		if err != nil {
			return err
		}
		return presenter.WriteTable(c.Out, report.Results)
	default:
		return fmt.Errorf("unknown command %q", cmd.Name)
	}
}
