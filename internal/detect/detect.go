package detect

import (
	"errors"
	"os"
	"strings"
)

func Detect(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Read magic bytes
	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	buf = buf[:n]

	ext := strings.ToLower(strings.TrimPrefix(getExt(path), "."))

	// Try magic first for common types
	if len(buf) >= 4 {
		// PNG
		if buf[0] == 0x89 && buf[1] == 0x50 && buf[2] == 0x4E && buf[3] == 0x47 {
			return "png", nil
		}
		// JPG
		if buf[0] == 0xFF && buf[1] == 0xD8 {
			return "jpg", nil
		}
		// GIF
		if string(buf[0:3]) == "GIF" {
			return "gif", nil
		}
		// PDF
		if string(buf[0:4]) == "%PDF" {
			return "pdf", nil
		}
		// ZIP based (docx/xlsx/pptx/odt etc often zip)
		if buf[0] == 0x50 && buf[1] == 0x4B {
			// Could be zip or office openxml; fallback to ext
			if ext != "" {
				return ext, nil
			}
			return "zip", nil
		}
	}

	if ext == "" {
		return "", errors.New("could not detect format (no extension/magic match)")
	}
	return ext, nil
}

func getExt(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i:]
		}
		if path[i] == '/' {
			break
		}
	}
	return ""
}
