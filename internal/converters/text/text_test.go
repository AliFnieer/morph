package textconv

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestGetExt(t *testing.T) {
	cases := map[string]string{
		"a/b/c.txt":   "txt",
		"file":        "",
		"a.b/c":       "",
		"archive.tar": "tar",
	}
	for path, want := range cases {
		if got := getExt(path); got != want {
			t.Errorf("getExt(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestToString(t *testing.T) {
	if got := toString("x"); got != "x" {
		t.Errorf("toString(string) = %q, want x", got)
	}
	if got := toString(map[string]int{"a": 1}); got != `{"a":1}` {
		t.Errorf("toString(map) = %q, want {\"a\":1}", got)
	}
}

func TestCSVToJSON(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.json")
	if err := csvToJSON([]byte("name,age\nalice,30\nbob,25\n"), out); err != nil {
		t.Fatal(err)
	}

	var rows []map[string]string
	if err := json.Unmarshal(readFile(t, out), &rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("len(rows) = %d, want 2", len(rows))
	}
	if rows[0]["name"] != "alice" || rows[0]["age"] != "30" {
		t.Fatalf("unexpected first row: %v", rows[0])
	}
}

func TestCSVToJSONEmpty(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.json")
	if err := csvToJSON(nil, out); err != nil {
		t.Fatal(err)
	}
	if got := string(readFile(t, out)); got != "[]" {
		t.Fatalf("output = %q, want []", got)
	}
}

func TestJSONToCSV(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.csv")
	input := []byte(`[{"name":"alice","age":"30"},{"name":"bob","age":"25"}]`)
	if err := jsonToCSV(input, out); err != nil {
		t.Fatal(err)
	}

	records, err := csv.NewReader(bytes.NewReader(readFile(t, out))).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 3 {
		t.Fatalf("len(records) = %d, want 3", len(records))
	}
	idx := map[string]int{}
	for i, h := range records[0] {
		idx[h] = i
	}
	if records[1][idx["name"]] != "alice" || records[1][idx["age"]] != "30" {
		t.Fatalf("unexpected first row: %v", records[1])
	}
}

func TestJSONToCSVEmpty(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.csv")
	if err := jsonToCSV([]byte("[]"), out); err != nil {
		t.Fatal(err)
	}
	if got := len(readFile(t, out)); got != 0 {
		t.Fatalf("output length = %d, want 0", got)
	}
}

func TestConvertCopy(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "a.txt")
	out := filepath.Join(dir, "b.txt")
	if err := os.WriteFile(in, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := (&conv{}).Convert(in, out); err != nil {
		t.Fatal(err)
	}
	if got := string(readFile(t, out)); got != "hello" {
		t.Fatalf("output = %q, want hello", got)
	}
}

func TestConvertCSVToJSON(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.csv")
	out := filepath.Join(dir, "out.json")
	if err := os.WriteFile(in, []byte("a,b\n1,2\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := (&conv{}).Convert(in, out); err != nil {
		t.Fatal(err)
	}
	var rows []map[string]string
	if err := json.Unmarshal(readFile(t, out), &rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0]["a"] != "1" || rows[0]["b"] != "2" {
		t.Fatalf("unexpected rows: %v", rows)
	}
}

func TestConvertMissingInput(t *testing.T) {
	dir := t.TempDir()
	if err := (&conv{}).Convert(filepath.Join(dir, "missing.txt"), filepath.Join(dir, "out.txt")); err == nil {
		t.Fatal("expected error for missing input")
	}
}

func TestMetadata(t *testing.T) {
	c := &conv{}
	if c.Name() != "stdlib-text" {
		t.Errorf("Name() = %q, want stdlib-text", c.Name())
	}
	if !slices.Contains(c.From(), "csv") {
		t.Errorf("From() missing csv: %v", c.From())
	}
	if !slices.Contains(c.To(), "json") {
		t.Errorf("To() missing json: %v", c.To())
	}
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
