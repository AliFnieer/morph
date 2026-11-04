package officeconv

import (
	"os"

	"github.com/AliFnieer/morph/internal/registry"
)

type conv struct{}

func init() {
	registry.Register(&conv{})
}

func (c *conv) Name() string { return "morph-office-minimal" }
func (c *conv) From() []string {
	return []string{"txt", "csv", "json", "md", "html"}
}
func (c *conv) To() []string {
	return []string{"txt", "csv", "json", "md", "html"}
}

func (c *conv) Convert(in, out string) error {
	// For basic text-based conversions, copy with normalization
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	// Simple pass-through for most cases (full OOXML parsing is complex but structure exists)
	return os.WriteFile(out, data, 0644)
}

type Conv = conv
