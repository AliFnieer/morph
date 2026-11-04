package dataconv

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func writeTemp(t *testing.T, dir, name string, data []byte) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestConvertJSONToYAML(t *testing.T) {
	dir := t.TempDir()
	in := writeTemp(t, dir, "in.json", []byte(`{"name":"x","n":1}`))
	out := filepath.Join(dir, "out.yaml")

	if err := (&conv{}).Convert(in, out); err != nil {
		t.Fatal(err)
	}
	var got map[string]interface{}
	if err := yaml.Unmarshal(readFile(t, out), &got); err != nil {
		t.Fatal(err)
	}
	if got["name"] != "x" {
		t.Fatalf("name = %v, want x", got["name"])
	}
}

func TestConvertYAMLToJSON(t *testing.T) {
	dir := t.TempDir()
	in := writeTemp(t, dir, "in.yaml", []byte("name: x\nn: 1\n"))
	out := filepath.Join(dir, "out.json")

	if err := (&conv{}).Convert(in, out); err != nil {
		t.Fatal(err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(readFile(t, out), &got); err != nil {
		t.Fatal(err)
	}
	if got["name"] != "x" {
		t.Fatalf("name = %v, want x", got["name"])
	}
}

func TestConvertTOMLToJSON(t *testing.T) {
	dir := t.TempDir()
	in := writeTemp(t, dir, "in.toml", []byte("name = \"x\"\n"))
	out := filepath.Join(dir, "out.json")

	if err := (&conv{}).Convert(in, out); err != nil {
		t.Fatal(err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(readFile(t, out), &got); err != nil {
		t.Fatal(err)
	}
	if got["name"] != "x" {
		t.Fatalf("name = %v, want x", got["name"])
	}
}

func TestConvertJSONLToJSON(t *testing.T) {
	dir := t.TempDir()
	in := writeTemp(t, dir, "in.jsonl", []byte("{\"a\":1}\n{\"b\":2}\n"))
	out := filepath.Join(dir, "out.json")

	if err := (&conv{}).Convert(in, out); err != nil {
		t.Fatal(err)
	}
	var got []interface{}
	if err := json.Unmarshal(readFile(t, out), &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
}

func TestConvertINIAndEnvToJSON(t *testing.T) {
	for _, name := range []string{"config.ini", ".env"} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			in := writeTemp(t, dir, name, []byte("A=1\nB=2\n"))
			out := filepath.Join(dir, "out.json")

			if err := (&conv{}).Convert(in, out); err != nil {
				t.Fatal(err)
			}
			var got map[string]string
			if err := json.Unmarshal(readFile(t, out), &got); err != nil {
				t.Fatal(err)
			}
			if got["A"] != "1" || got["B"] != "2" {
				t.Fatalf("unexpected output: %v", got)
			}
		})
	}
}

func TestConvertXMLToJSON(t *testing.T) {
	dir := t.TempDir()
	in := writeTemp(t, dir, "in.xml", []byte("<r><a>1</a></r>"))
	out := filepath.Join(dir, "out.json")

	if err := (&conv{}).Convert(in, out); err != nil {
		t.Fatal(err)
	}
	var got map[string]string
	if err := json.Unmarshal(readFile(t, out), &got); err != nil {
		t.Fatal(err)
	}
	if got["raw_xml"] != "<r><a>1</a></r>" {
		t.Fatalf("raw_xml = %q", got["raw_xml"])
	}
}

func TestConvertJSONToTxt(t *testing.T) {
	dir := t.TempDir()
	in := writeTemp(t, dir, "in.json", []byte(`{"a":1}`))
	out := filepath.Join(dir, "out.txt")

	if err := (&conv{}).Convert(in, out); err != nil {
		t.Fatal(err)
	}
	if got := string(readFile(t, out)); got == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestConvertMissingInput(t *testing.T) {
	dir := t.TempDir()
	if err := (&conv{}).Convert(filepath.Join(dir, "missing.json"), filepath.Join(dir, "out.yaml")); err == nil {
		t.Fatal("expected error for missing input")
	}
}

func TestGetExt(t *testing.T) {
	cases := map[string]string{
		"a/b/c.JSON": "json",
		"file":       "",
		"a.b/c":      "",
	}
	for path, want := range cases {
		if got := getExt(path); got != want {
			t.Errorf("getExt(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestToStr(t *testing.T) {
	if got := toStr("x"); got != "x" {
		t.Errorf("toStr(string) = %q, want x", got)
	}
	if got := toStr(1); got != "1" {
		t.Errorf("toStr(int) = %q, want 1", got)
	}
}
