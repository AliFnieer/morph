package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func runMorph(t *testing.T, args ...string) {
	t.Helper()
	cmd := exec.Command("./morph", args...)
	cmd.Dir = filepath.Join("..")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("morph %v failed: %v\n%s", args, err, string(out))
	}
}

func TestHelp(t *testing.T) {
	runMorph(t, "--help")
}

func TestList(t *testing.T) {
	runMorph(t, "--list")
}

func TestListJSON(t *testing.T) {
	runMorph(t, "--list-json")
}

func TestVersion(t *testing.T) {
	runMorph(t, "--version")
}

func TestTextConversions(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(in, []byte("hello world\nfoo bar"), 0644); err != nil {
		t.Fatal(err)
	}
	runMorph(t, "--in", in, "--to", "txt", "--out", filepath.Join(dir, "out.txt"))
	runMorph(t, "--in", in, "--to", "md", "--out", filepath.Join(dir, "out.md"))
	runMorph(t, "--in", in, "--to", "json", "--out", filepath.Join(dir, "out.json"))
	runMorph(t, "--in", in, "--to", "html", "--out", filepath.Join(dir, "out.html"))
	runMorph(t, "--in", in, "--to", "pdf", "--out", filepath.Join(dir, "out.pdf"))
}

func TestCSVJSON(t *testing.T) {
	dir := t.TempDir()
	csv := filepath.Join(dir, "data.csv")
	if err := os.WriteFile(csv, []byte("name,age\nAlice,30\nBob,25\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runMorph(t, "--in", csv, "--to", "json", "--out", filepath.Join(dir, "data.json"))
	runMorph(t, "--in", filepath.Join(dir, "data.json"), "--to", "csv", "--out", filepath.Join(dir, "data2.csv"))
	runMorph(t, "--in", csv, "--to", "txt", "--out", filepath.Join(dir, "data.txt"))
}

func TestYAML(t *testing.T) {
	dir := t.TempDir()
	yaml := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(yaml, []byte("foo:\n  bar: 123\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runMorph(t, "--in", yaml, "--to", "json", "--out", filepath.Join(dir, "config.json"))
	runMorph(t, "--in", filepath.Join(dir, "config.json"), "--to", "yaml", "--out", filepath.Join(dir, "config2.yaml"))
	runMorph(t, "--in", yaml, "--to", "toml", "--out", filepath.Join(dir, "config.toml"))
}

func TestJSONL(t *testing.T) {
	dir := t.TempDir()
	ndjson := filepath.Join(dir, "data.ndjson")
	if err := os.WriteFile(ndjson, []byte("{\"a\":1}\n{\"a\":2}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runMorph(t, "--in", ndjson, "--to", "json", "--out", filepath.Join(dir, "data.json"))
}

func TestHTML(t *testing.T) {
	dir := t.TempDir()
	html := filepath.Join(dir, "page.html")
	if err := os.WriteFile(html, []byte("<html><body><h1>Hi</h1></body></html>"), 0644); err != nil {
		t.Fatal(err)
	}
	runMorph(t, "--in", html, "--to", "txt", "--out", filepath.Join(dir, "page.txt"))
	runMorph(t, "--in", filepath.Join(dir, "page.txt"), "--to", "html", "--out", filepath.Join(dir, "page2.html"))
}

func TestRTF(t *testing.T) {
	dir := t.TempDir()
	rtf := filepath.Join(dir, "test.rtf")
	if err := os.WriteFile(rtf, []byte("{\\rtf1\\ansi Hello World}"), 0644); err != nil {
		t.Fatal(err)
	}
	runMorph(t, "--in", rtf, "--to", "txt", "--out", filepath.Join(dir, "test.txt"))
	runMorph(t, "--in", rtf, "--to", "pdf", "--out", filepath.Join(dir, "test.pdf"))
}

func TestArchives(t *testing.T) {
	dir := t.TempDir()
	// create zip
	zipf := filepath.Join(dir, "test.zip")
	// use zip command
	cmd := exec.Command("zip", "-q", zipf, filepath.Join("..", "README.md"))
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Skip("zip not available")
	}
	runMorph(t, "--in", zipf, "--to", "tar.gz", "--out", filepath.Join(dir, "test.tgz"))
	runMorph(t, "--in", filepath.Join(dir, "test.tgz"), "--to", "tar.xz", "--out", filepath.Join(dir, "test.tar.xz"))
	runMorph(t, "--in", filepath.Join(dir, "test.tgz"), "--to", "tar.bz2", "--out", filepath.Join(dir, "test.tar.bz2"))
}

func TestImages(t *testing.T) {
	dir := t.TempDir()
	// try to create bmp with python if available
	bmp := filepath.Join(dir, "test.bmp")
	cmd := exec.Command("python3", "-c", "from PIL import Image; Image.new('RGB',(10,10),(255,0,0)).save('"+bmp+"')")
	if err := cmd.Run(); err == nil {
		runMorph(t, "--in", bmp, "--to", "png", "--out", filepath.Join(dir, "test.png"))
	}
}

func TestSVG(t *testing.T) {
	dir := t.TempDir()
	svg := filepath.Join(dir, "test.svg")
	if err := os.WriteFile(svg, []byte("<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"10\" height=\"10\"><rect width=\"10\" height=\"10\" fill=\"blue\"/></svg>"), 0644); err != nil {
		t.Fatal(err)
	}
	runMorph(t, "--in", svg, "--to", "png", "--out", filepath.Join(dir, "test.png"))
}

func TestINI(t *testing.T) {
	dir := t.TempDir()
	ini := filepath.Join(dir, "config.ini")
	if err := os.WriteFile(ini, []byte("[section]\nkey=value\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runMorph(t, "--in", ini, "--to", "json", "--out", filepath.Join(dir, "config.json"))
}
