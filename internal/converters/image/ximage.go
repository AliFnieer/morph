package imageconv

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/AliFnieer/morph/internal/registry"
	_ "golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

type xconv struct{}

func init() {
	registry.Register(&xconv{})
}

func (c *xconv) Name() string   { return "x/image" }
func (c *xconv) From() []string { return []string{"tiff", "tif", "bmp", "webp"} }
func (c *xconv) To() []string   { return []string{"png", "jpg", "jpeg"} }

func (c *xconv) Convert(in, out string) error {
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

	ext := strings.ToLower(getExtPath(out))
	if ext == "jpg" || ext == "jpeg" {
		return jpeg.Encode(outf, img, &jpeg.Options{Quality: 90})
	}
	return png.Encode(outf, img)
}

func getExtPath(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
		if path[i] == '/' || path[i] == '\\' {
			break
		}
	}
	return ""
}

type xconv2 struct{}

func init() {
	registry.Register(&xconv2{})
}

func (c *xconv2) Name() string   { return "x/image-tiff" }
func (c *xconv2) From() []string { return []string{"tiff", "tif"} }
func (c *xconv2) To() []string   { return []string{"png", "jpg", "jpeg", "tiff", "tif"} }

func (c *xconv2) Convert(in, out string) error {
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
	ext := strings.ToLower(getExtPath(out))
	if ext == "tif" || ext == "tiff" {
		return tiff.Encode(outf, img, &tiff.Options{})
	}
	if ext == "jpg" || ext == "jpeg" {
		return jpeg.Encode(outf, img, &jpeg.Options{Quality: 90})
	}
	return png.Encode(outf, img)
}
