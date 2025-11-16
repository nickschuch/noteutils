package validate

import "github.com/nickschuch/noteutils/internal/notes"

// Specification used to define a set of notes and compare during validation.
type Specification struct {
	// List of probes keyed by architecture.
	Probes map[string][]notes.Probe `yaml:"probes,omitempty"`
}
