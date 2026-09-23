package cli

import (
	"fmt"
	"runtime"

	"github.com/AliFnieer/morph/internal/registry"
)

func RunDoctor() {
	fmt.Print("morph doctor\n\n")
	fmt.Printf("Version: v0.0.1\n")
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Go: %s\n", runtime.Version())
	fmt.Printf("Converters loaded: %d\n", registry.Count())
	fmt.Printf("Mode: offline-only\n")
}
