package vectorconv

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/AliFnieer/morph/internal/registry"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

type conv struct{}

func init() {
	registry.Register(&conv{})
}

func (c *conv) Name() string   { return "oksvg" }
func (c *conv) From() []string { return []string{"svg"} }
func (c *conv) To() []string   { return []string{"png", "jpg", "jpeg"} }

func (c *conv) Convert(in, out string) error {
	icon, err := oksvg.ReadIcon(in)
	if err != nil {
		return err
	}
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	if w == 0 {
		w = 1024
	}
	if h == 0 {
		h = 1024
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	scanner := rasterx.NewScannerGV(w, h, img, img.Bounds())
	raster := rasterx.NewDasher(w, h, scanner)
	icon.Draw(raster, 1.0)
	outf, err := os.Create(out)
	if err != nil {
		return err
	}
	defer outf.Close()
	ext := strings.ToLower(getExt(out))
	if ext == "jpg" || ext == "jpeg" {
		return jpeg.Encode(outf, img, &jpeg.Options{Quality: 90})
	}
	return png.Encode(outf, img)
}

func getExt(path string) string {
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

type Conv = conv
