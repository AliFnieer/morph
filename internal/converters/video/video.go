package videoconv

import (
	"os"
	"strings"

	"github.com/AliFnieer/morph/internal/registry"
)

type conv struct{}

func init() {
	registry.Register(&conv{})
}

func (c *conv) Name() string   { return "morph-video-stub" }
func (c *conv) From() []string { return []string{"mp4", "mkv", "webm", "avi", "mov", "flv"} }
func (c *conv) To() []string   { return []string{"mp4", "webm", "mkv"} }

func (c *conv) Convert(in, out string) error {
	if getExt(in) == getExt(out) {
		data, err := os.ReadFile(in)
		if err != nil {
			return err
		}
		return os.WriteFile(out, data, 0644)
	}
	return os.ErrNotExist
}

func getExt(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return strings.ToLower(path[i+1:])
		}
		if path[i] == '/' || path[i] == '\\' {
			break
		}
	}
	return ""
}

type Conv = conv
