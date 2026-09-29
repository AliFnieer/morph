package imageconv

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/AliFnieer/morph/internal/registry"
)

type conv struct{}

func init() {
	registry.Register(&conv{})
}

func (c *conv) Name() string { return "stdlib-image" }
func (c *conv) From() []string {
	return []string{"png", "jpg", "jpeg", "gif", "bmp"}
}
func (c *conv) To() []string {
	return []string{"png", "jpg", "jpeg", "gif"}
}

func (c *conv) Convert(in, out string) error {
	f, err := os.Open(in)
	if err != nil {
		return err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	outf, err := os.Create(out)
	if err != nil {
		return err
	}
	defer outf.Close()

	ext := strings.ToLower(getExt(out))
	switch ext {
	case "jpg", "jpeg":
		return jpeg.Encode(outf, img, &jpeg.Options{Quality: 90})
	case "png":
		return png.Encode(outf, img)
	case "gif":
		return gif.Encode(outf, img, &gif.Options{})
	default:
		// try png
		return png.Encode(outf, img)
	}
}

func getExt(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
		if path[i] == '/' {
			break
		}
	}
	return ""
}
