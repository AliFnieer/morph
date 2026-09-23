package tests

import (
	"testing"

	"github.com/AliFnieer/morph/internal/registry"
)

func TestRegistryHasConversions(t *testing.T) {
	// just ensure registry is accessible
	_ = registry.Get("png", "jpg")
}
