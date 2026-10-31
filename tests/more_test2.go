package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func rm(t *testing.T, args ...string) {
	t.Helper()
	cmd := exec.Command("./morph", args...)
	cmd.Dir = filepath.Join("..")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("morph %v failed: %v\n%s", args, err, string(out))
	}
}

func TestDoctorCmd2(t *testing.T) {
	cmd := exec.Command("./morph", "doctor")
	cmd.Dir = filepath.Join("..")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "Mode") {
		t.Fatal("bad")
	}
}

func TestCompletion2(t *testing.T) {
	cmd := exec.Command("./morph", "completion", "bash")
	cmd.Dir = filepath.Join("..")
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	cmd = exec.Command("./morph", "completion", "zsh")
	cmd.Dir = filepath.Join("..")
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	cmd = exec.Command("./morph", "completion", "fish")
	cmd.Dir = filepath.Join("..")
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestStdinStdout2(t *testing.T) {
	cmd := exec.Command("bash", "-c", "echo '{\"a\":1}' | ./morph --in - --to yaml --out -")
	cmd.Dir = filepath.Join("..")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err, string(out))
	}
	if !strings.Contains(string(out), "a") {
		t.Fatal(string(out))
	}
}

func TestForce2(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "t.txt")
	out := filepath.Join(dir, "o.txt")
	os.WriteFile(in, []byte("hi"), 0644)
	rm(t, "--in", in, "--to", "txt", "--out", out)
	cmd := exec.Command("./morph", "--in", in, "--to", "txt", "--out", out)
	cmd.Dir = filepath.Join("..")
	if err := cmd.Run(); err == nil {
		t.Fatal("bad")
	}
	rm(t, "--in", in, "--to", "txt", "--out", out, "--force")
}

func TestTOML2(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "c.toml"), []byte("name=\"x\"\n"), 0644)
	rm(t, "--in", filepath.Join(dir, "c.toml"), "--to", "json", "--out", filepath.Join(dir, "c.json"))
	rm(t, "--in", filepath.Join(dir, "c.json"), "--to", "toml", "--out", filepath.Join(dir, "c2.toml"))
}

func TestPDFMulti2(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "m.txt"), []byte("a\nb\nc\n"), 0644)
	rm(t, "--in", filepath.Join(dir, "m.txt"), "--to", "pdf", "--out", filepath.Join(dir, "m.pdf"))
}

func TestPDFToTXT2(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "m.txt"), []byte("hello test"), 0644)
	rm(t, "--in", filepath.Join(dir, "m.txt"), "--to", "pdf", "--out", filepath.Join(dir, "m.pdf"))
	rm(t, "--in", filepath.Join(dir, "m.pdf"), "--to", "txt", "--out", filepath.Join(dir, "m2.txt"))
}

func TestMD2(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "r.md"), []byte("# H\nhi"), 0644)
	rm(t, "--in", filepath.Join(dir, "r.md"), "--to", "html", "--out", filepath.Join(dir, "r.html"))
	rm(t, "--in", filepath.Join(dir, "r.md"), "--to", "pdf", "--out", filepath.Join(dir, "r.pdf"))
	rm(t, "--in", filepath.Join(dir, "r.md"), "--to", "txt", "--out", filepath.Join(dir, "r.txt"))
}

func TestXML2(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "d.xml"), []byte("<r><a>1</a></r>"), 0644)
	rm(t, "--in", filepath.Join(dir, "d.xml"), "--to", "json", "--out", filepath.Join(dir, "d.json"))
}

func TestENV2(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".env"), []byte("A=1\nB=2\n"), 0644)
	rm(t, "--in", filepath.Join(dir, ".env"), "--to", "json", "--out", filepath.Join(dir, "env.json"))
}

func TestJSONL2YAML(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "d.ndjson"), []byte("{\"a\":1}\n{\"b\":2}\n"), 0644)
	rm(t, "--in", filepath.Join(dir, "d.ndjson"), "--to", "yaml", "--out", filepath.Join(dir, "d.yaml"))
}

func TestArchiveMulti2(t *testing.T) {
	dir := t.TempDir()
	z := filepath.Join(dir, "t.zip")
	cmd := exec.Command("zip", "-q", z, "README.md")
	cmd.Dir = filepath.Join("..")
	if err := cmd.Run(); err != nil {
		t.Skip("zip")
	}
	rm(t, "--in", z, "--to", "tar.gz", "--out", filepath.Join(dir, "t.tgz"))
	rm(t, "--in", filepath.Join(dir, "t.tgz"), "--to", "tar.xz", "--out", filepath.Join(dir, "t.tar.xz"))
	rm(t, "--in", filepath.Join(dir, "t.tgz"), "--to", "tar.bz2", "--out", filepath.Join(dir, "t.tar.bz2"))
}

func TestQuiet2(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "q.txt")
	os.WriteFile(in, []byte("q"), 0644)
	cmd := exec.Command("./morph", "--in", in, "--to", "txt", "--out", filepath.Join(dir, "q2.txt"), "--quiet")
	cmd.Dir = filepath.Join("..")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "Converting") {
		t.Fatal("bad")
	}
}
