package transport

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"drivingmaterials/workflow"
)

type BatchRunner struct {
	Ingest *workflow.IngestService
	Out    io.Writer
}

func (b *BatchRunner) Run(lines []string) (int, int) {
	accepted := 0
	failed := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		command, err := ParseCommand(line)
		if err != nil {
			failed++
			fmt.Fprintln(b.Out, "error:", err)
			continue
		}
		request, err := BuildIngest(command)
		if err != nil {
			failed++
			fmt.Fprintln(b.Out, "error:", err)
			continue
		}
		if _, err := b.Ingest.Ingest(request); err != nil {
			failed++
			fmt.Fprintln(b.Out, "error:", err)
			continue
		}
		accepted++
	}
	return accepted, failed
}

func (b *BatchRunner) RunReader(reader io.Reader) (int, int, error) {
	scanner := bufio.NewScanner(reader)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}
	accepted, failed := b.Run(lines)
	return accepted, failed, nil
}

func (b *BatchRunner) Summary(accepted, failed int) string {
	return fmt.Sprintf("accepted=%d failed=%d", accepted, failed)
}
