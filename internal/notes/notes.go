package notes

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// LoadFromFile executes: readelf -n <path> and parses the output.
func LoadFromFile(path string) ([]Probe, error) {
	cmd := exec.Command("readelf", "-n", path)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running readelf: %w", err)
	}

	return ParseSpecification(strings.NewReader(string(output)))
}

// ParseSpecification parses readelf -n output.
func ParseSpecification(r io.Reader) ([]Probe, error) {
	scanner := bufio.NewScanner(r)

	inStapsdt := false
	var current *Probe
	var probes []Probe

	flushCurrent := func() {
		if current != nil {
			p := *current
			probes = append(probes, p)
			current = nil
		}
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Detect .note.stapsdt entering/exiting
		if strings.HasPrefix(line, "Displaying notes found in:") {
			if inStapsdt && !strings.Contains(line, ".note.stapsdt") {
				flushCurrent()
				inStapsdt = false
			}
			if strings.Contains(line, ".note.stapsdt") {
				inStapsdt = true
			}
			continue
		}

		if !inStapsdt {
			continue
		}

		// New STAPSDT record
		if strings.HasPrefix(line, "stapsdt") && strings.Contains(line, "NT_STAPSDT") {
			flushCurrent()
			current = &Probe{}
			continue
		}

		if current == nil {
			continue
		}

		switch {
		case strings.HasPrefix(line, "Provider:"):
			current.Provider = strings.TrimSpace(strings.TrimPrefix(line, "Provider:"))

		case strings.HasPrefix(line, "Name:"):
			current.Name = strings.TrimSpace(strings.TrimPrefix(line, "Name:"))

		case strings.HasPrefix(line, "Arguments:"):
			argsStr := strings.TrimSpace(strings.TrimPrefix(line, "Arguments:"))

			if argsStr != "" {
				current.Arguments = append(current.Arguments, strings.Fields(argsStr)...)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan readelf output: %w", err)
	}

	flushCurrent()

	return probes, nil
}
