package models

import (
	"testing"
)

func TestArchitectureSupport(t *testing.T) {
	arch := ArchitectureSupport{
		OperatorName:   "test-operator",
		BundleName:     "test-operator.v1.0.0",
		Package:        "test-operator",
		AMD64:          "supported",
		ARM64:          "supported",
		SupportedArchs: []string{"amd64", "arm64"},
	}

	if arch.OperatorName != "test-operator" {
		t.Errorf("Expected OperatorName to be 'test-operator', got '%s'", arch.OperatorName)
	}

	if len(arch.SupportedArchs) != 2 {
		t.Errorf("Expected 2 supported architectures, got %d", len(arch.SupportedArchs))
	}
}

func TestOperatorIndexEntry(t *testing.T) {
	entry := OperatorIndexEntry{
		Schema:  "olm.bundle",
		Name:    "test-bundle",
		Package: "test-package",
	}

	if entry.Schema != "olm.bundle" {
		t.Errorf("Expected Schema to be 'olm.bundle', got '%s'", entry.Schema)
	}

	if entry.Name != "test-bundle" {
		t.Errorf("Expected Name to be 'test-bundle', got '%s'", entry.Name)
	}
}

