package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, name string, data []byte) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestDetectMagic(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		want string
	}{
		{"image.txt", []byte{0x89, 'P', 'N', 'G', 0x0D}, "png"},
		{"photo.bin", []byte{0xFF, 0xD8, 0xFF, 0xE0}, "jpg"},
		{"anim.dat", []byte("GIF89a"), "gif"},
		{"document.weird", []byte("%PDF-1.4"), "pdf"},
	}

	for _, tc := range cases {
		got, err := Detect(writeFile(t, tc.name, tc.data))
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if got != tc.want {
			t.Errorf("%s: Detect() = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestDetectZip(t *testing.T) {
	got, err := Detect(writeFile(t, "book.docx", []byte{'P', 'K', 3, 4}))
	if err != nil {
		t.Fatal(err)
	}
	if got != "docx" {
		t.Errorf("Detect() = %q, want docx", got)
	}

	got, err = Detect(writeFile(t, "archive", []byte{'P', 'K', 3, 4}))
	if err != nil {
		t.Fatal(err)
	}
	if got != "zip" {
		t.Errorf("Detect() = %q, want zip", got)
	}
}

func TestDetectByExtension(t *testing.T) {
	got, err := Detect(writeFile(t, "data.json", []byte(`{"a":1}`)))
	if err != nil {
		t.Fatal(err)
	}
	if got != "json" {
		t.Errorf("Detect() = %q, want json", got)
	}
}

func TestDetectNoExtension(t *testing.T) {
	if _, err := Detect(writeFile(t, "plain", []byte("hello"))); err == nil {
		t.Fatal("expected error for file without extension or magic")
	}
}

func TestDetectMissingFile(t *testing.T) {
	if _, err := Detect(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestGetExt(t *testing.T) {
	cases := map[string]string{
		"a/b/c.txt": ".txt",
		"file":      "",
		"a.b/c":     "",
		".env":      ".env",
	}
	for path, want := range cases {
		if got := getExt(path); got != want {
			t.Errorf("getExt(%q) = %q, want %q", path, got, want)
		}
	}
}
