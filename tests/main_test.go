package tests

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	build := exec.Command("go", "build", "-o", "morph", "./cmd/morph")
	build.Dir = ".."
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		os.Exit(1)
	}

	code := m.Run()

	_ = os.Remove("../morph")
	os.Exit(code)
}
