package registry

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

type fakeConverter struct {
	name string
	from []string
	to   []string
}

func (f fakeConverter) Convert(in, out string) error { return nil }
func (f fakeConverter) Name() string                 { return f.name }
func (f fakeConverter) From() []string               { return f.from }
func (f fakeConverter) To() []string                 { return f.to }

func withRegistry(t *testing.T, cs ...Converter) {
	t.Helper()
	old := converters
	converters = append([]Converter(nil), cs...)
	t.Cleanup(func() { converters = old })
}

func TestRegisterAndGet(t *testing.T) {
	withRegistry(t, fakeConverter{name: "fake", from: []string{"a"}, to: []string{"b"}})

	if got := Get("a", "b"); got == nil || got.Name() != "fake" {
		t.Fatalf("Get(a, b) = %v, want fake", got)
	}
	if got := Get("a", "c"); got != nil {
		t.Fatalf("Get(a, c) = %v, want nil", got)
	}
	if got := Get("x", "b"); got != nil {
		t.Fatalf("Get(x, b) = %v, want nil", got)
	}
}

func TestRegisterAppends(t *testing.T) {
	withRegistry(t)
	Register(fakeConverter{name: "one", from: []string{"a"}, to: []string{"b"}})
	Register(fakeConverter{name: "two", from: []string{"b"}, to: []string{"c"}})
	if len(converters) != 2 {
		t.Fatalf("len(converters) = %d, want 2", len(converters))
	}
}

func TestAllPairsAndCount(t *testing.T) {
	withRegistry(t,
		fakeConverter{name: "one", from: []string{"a"}, to: []string{"b", "c"}},
		fakeConverter{name: "two", from: []string{"a"}, to: []string{"b"}},
	)

	if got := len(AllPairs()); got != 2 {
		t.Fatalf("len(AllPairs()) = %d, want 2", got)
	}
	if got := Count(); got != 2 {
		t.Fatalf("Count() = %d, want 2", got)
	}
}

func TestCategory(t *testing.T) {
	cases := map[string]string{
		"stdlib-image":     "Images",
		"x/image":          "Images",
		"oksvg":            "Vector",
		"morph-data":       "Data/Config",
		"stdlib-text":      "Text/Markup",
		"morph-pdf":        "Documents",
		"morph-ooxml":      "Office (OOXML)",
		"morph-archive":    "Archives",
		"morph-audio":      "Audio",
		"morph-video-stub": "Video",
		"unknown":          "Other",
	}
	for name, want := range cases {
		if got := category(name); got != want {
			t.Errorf("category(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestListJSON(t *testing.T) {
	withRegistry(t, fakeConverter{name: "morph-data", from: []string{"yaml"}, to: []string{"json"}})

	var out []map[string]string
	if err := json.Unmarshal([]byte(captureStdout(t, ListJSON)), &out); err != nil {
		t.Fatalf("ListJSON output is not valid JSON: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("len(output) = %d, want 1", len(out))
	}
	entry := out[0]
	if entry["from"] != "yaml" || entry["to"] != "json" || entry["via"] != "morph-data" || entry["category"] != "Data/Config" {
		t.Fatalf("unexpected entry: %v", entry)
	}
}

func TestList(t *testing.T) {
	withRegistry(t, fakeConverter{name: "morph-data", from: []string{"yaml"}, to: []string{"json"}})

	got := captureStdout(t, List)
	if !strings.Contains(got, "[Data/Config]") || !strings.Contains(got, "yaml -> json") {
		t.Fatalf("List output missing expected content:\n%s", got)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = old

	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Close(); err != nil {
		t.Fatal(err)
	}
	return string(data)
}
