package audioconv

import (
	"os"
	"strings"

	"github.com/AliFnieer/morph/internal/registry"
)

type conv struct{}

func init() {
	registry.Register(&conv{})
}

func (c *conv) Name() string { return "morph-audio" }
func (c *conv) From() []string {
	return []string{"mp3", "ogg", "flac", "aac", "m4a", "opus", "wma", "wav"}
}
func (c *conv) To() []string { return []string{"wav", "raw", "pcm"} }

func (c *conv) Convert(in, out string) error {
	inExt := getExt(in)
	outExt := getExt(out)
	if inExt == "wav" && outExt == "raw" {
		data, err := os.ReadFile(in)
		if err != nil {
			return err
		}
		if len(data) > 44 {
			return os.WriteFile(out, data[44:], 0644)
		}
		return os.WriteFile(out, data, 0644)
	}
	if inExt == "raw" && outExt == "wav" {
		data, err := os.ReadFile(in)
		if err != nil {
			return err
		}
		hdr := make([]byte, 44)
		// minimal header placeholder
		copy(hdr[0:4], "RIFF")
		copy(hdr[8:12], "WAVE")
		copy(hdr[12:16], "fmt ")
		return os.WriteFile(out, append(hdr, data...), 0644)
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
