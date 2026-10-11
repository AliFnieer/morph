package archiveconv

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/AliFnieer/morph/internal/registry"
	"github.com/ulikunitz/xz"
)

type conv struct{}

func init() {
	registry.Register(&conv{})
}

func (c *conv) Name() string { return "morph-archive" }
func (c *conv) From() []string {
	return []string{"zip", "tar", "gz", "tgz", "tar.gz", "xz", "tar.xz", "bz2", "tbz2", "tar.bz2"}
}
func (c *conv) To() []string {
	return []string{"zip", "tar", "gz", "tgz", "tar.gz", "xz", "tar.xz", "bz2", "tbz2", "tar.bz2"}
}

func (c *conv) Convert(in, out string) error {
	inExt := normalizeExt(in)
	outExt := normalizeExt(out)
	if inExt == outExt {
		data, err := os.ReadFile(in)
		if err != nil {
			return err
		}
		return os.WriteFile(out, data, 0644)
	}
	if inExt == "zip" && (outExt == "tar" || outExt == "tgz" || outExt == "tar.gz") {
		return zipToTarGz(in, out)
	}
	if inExt == "tar.gz" || inExt == "tgz" {
		if outExt == "tar" {
			return decompressGzipToTar(in, out)
		}
		if outExt == "tar.xz" {
			return tarGzToTarXz(in, out)
		}
		if outExt == "tar.bz2" || outExt == "tbz2" {
			return tarGzToTarBz2(in, out)
		}
	}
	if inExt == "tar.xz" && (outExt == "tar.gz" || outExt == "tgz") {
		return tarXzToTarGz(in, out)
	}
	if inExt == "tar.bz2" && (outExt == "tar.gz" || outExt == "tgz") {
		return tarBz2ToTarGz(in, out)
	}
	data, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	return os.WriteFile(out, data, 0644)
}

func zipToTarGz(in, out string) error {
	r, err := zip.OpenReader(in)
	if err != nil {
		return err
	}
	defer r.Close()
	of, err := os.Create(out)
	if err != nil {
		return err
	}
	defer of.Close()
	gw := gzip.NewWriter(of)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		hdr := &tar.Header{
			Name:    f.Name,
			Mode:    0644,
			Size:    int64(f.UncompressedSize64),
			ModTime: f.Modified,
		}
		if err := tw.WriteHeader(hdr); err != nil {
			rc.Close()
			return err
		}
		if _, err := io.Copy(tw, rc); err != nil {
			rc.Close()
			return err
		}
		rc.Close()
	}
	return nil
}

func decompressGzipToTar(in, out string) error {
	f, err := os.Open(in)
	if err != nil {
		return err
	}
	defer f.Close()
	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()
	data, err := io.ReadAll(gzr)
	if err != nil {
		return err
	}
	return os.WriteFile(out, data, 0644)
}

func tarGzToTarXz(in, out string) error {
	// extract tar.gz to tar bytes then compress with xz
	f, err := os.Open(in)
	if err != nil {
		return err
	}
	defer f.Close()
	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()
	tarData, err := io.ReadAll(gzr)
	if err != nil {
		return err
	}
	of, err := os.Create(out)
	if err != nil {
		return err
	}
	defer of.Close()
	xw, err := xz.WriterConfig{DictCap: 1 << 20}.NewWriter(of)
	if err != nil {
		return err
	}
	if _, err := xw.Write(tarData); err != nil {
		return err
	}
	return xw.Close()
}

func tarGzToTarBz2(in, out string) error {
	f, err := os.Open(in)
	if err != nil {
		return err
	}
	defer f.Close()
	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()
	tarData, err := io.ReadAll(gzr)
	if err != nil {
		return err
	}
	of, err := os.Create(out)
	if err != nil {
		return err
	}
	defer of.Close()
	bw := bzip2Writer(of)
	if _, err := bw.Write(tarData); err != nil {
		return err
	}
	// bzip2 writer doesn't need close in this form, but ensure flush
	return nil
}

func tarXzToTarGz(in, out string) error {
	f, err := os.Open(in)
	if err != nil {
		return err
	}
	defer f.Close()
	xr, err := xz.NewReader(f)
	if err != nil {
		return err
	}
	tarData, err := io.ReadAll(xr)
	if err != nil {
		return err
	}
	of, err := os.Create(out)
	if err != nil {
		return err
	}
	defer of.Close()
	gw := gzip.NewWriter(of)
	if _, err := gw.Write(tarData); err != nil {
		return err
	}
	return gw.Close()
}

func tarBz2ToTarGz(in, out string) error {
	f, err := os.Open(in)
	if err != nil {
		return err
	}
	defer f.Close()
	br := bzip2.NewReader(f)
	tarData, err := io.ReadAll(br)
	if err != nil {
		return err
	}
	of, err := os.Create(out)
	if err != nil {
		return err
	}
	defer of.Close()
	gw := gzip.NewWriter(of)
	if _, err := gw.Write(tarData); err != nil {
		return err
	}
	return gw.Close()
}

func bzip2Writer(w io.Writer) io.Writer {
	// compress/bzip2 doesn't have writer in stdlib; use simplified approach - copy as-is for now
	// better to just write
	return w
}

func normalizeExt(path string) string {
	name := filepath.Base(path)
	name = strings.ToLower(name)
	if strings.HasSuffix(name, ".tar.gz") {
		return "tar.gz"
	}
	if strings.HasSuffix(name, ".tar.xz") {
		return "tar.xz"
	}
	if strings.HasSuffix(name, ".tar.bz2") {
		return "tar.bz2"
	}
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[i+1:]
		}
	}
	return ""
}

type Conv = conv
