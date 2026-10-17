package rawconv

import (
	"os"
	"strings"

	"github.com/AliFnieer/morph/internal/registry"
)

type conv struct{}

func init() {
	registry.Register(&conv{})
}

func (c *conv) Name() string { return "morph-raw-copy" }
func (c *conv) From() []string {
	return []string{"bin", "dat", "raw", "hex", "bmp", "ico", "tiff", "tif", "webp", "avif", "heic", "heif"}
}
func (c *conv) To() []string {
	return []string{"bin", "dat", "raw", "hex", "bmp", "ico", "tiff", "tif"}
}

func (c *conv) Convert(in, out string) error {
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	// For unsupported encoded formats, try byte copy if user explicitly requests
	return os.WriteFile(out, data, 0644)
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
