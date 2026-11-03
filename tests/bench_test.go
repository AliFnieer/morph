package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func BenchmarkTextToPDF(b *testing.B) {
	dir := b.TempDir()
	in := filepath.Join(dir, "b.txt")
	os.WriteFile(in, []byte("hello world\nfoo bar\n"), 0644)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("./morph", "--in", in, "--to", "pdf", "--out", filepath.Join(dir, "out.pdf"), "--force")
		cmd.Dir = filepath.Join("..")
		if err := cmd.Run(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCSVToJSON(b *testing.B) {
	dir := b.TempDir()
	in := filepath.Join(dir, "b.csv")
	os.WriteFile(in, []byte("a,b\n1,2\n3,4\n"), 0644)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("./morph", "--in", in, "--to", "json", "--out", filepath.Join(dir, "out.json"), "--force")
		cmd.Dir = filepath.Join("..")
		if err := cmd.Run(); err != nil {
			b.Fatal(err)
		}
	}
}
