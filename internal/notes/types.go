package notes

// Probe represents a single probe configuration.
type Probe struct {
	Provider  string   `yaml:"provider,omitempty"`
	Name      string   `yaml:"name,omitempty"`
	Arguments []string `yaml:"arguments,omitempty"`
}
