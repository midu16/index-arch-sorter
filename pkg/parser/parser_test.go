package parser

import (
	"os"
	"testing"

	"github.com/midu/index-arch-sorter/pkg/models"
)

func TestExtractArchLabels(t *testing.T) {
	labels := map[string]interface{}{
		"operatorframework.io/arch.amd64":   "supported",
		"operatorframework.io/arch.arm64":   "supported",
		"operatorframework.io/arch.ppc64le": "supported",
	}

	archSupport := &models.ArchitectureSupport{}
	extractArchLabels(labels, archSupport)

	if archSupport.AMD64 != "supported" {
		t.Errorf("Expected AMD64 to be 'supported', got '%s'", archSupport.AMD64)
	}

	if archSupport.ARM64 != "supported" {
		t.Errorf("Expected ARM64 to be 'supported', got '%s'", archSupport.ARM64)
	}

	if archSupport.PPC64LE != "supported" {
		t.Errorf("Expected PPC64LE to be 'supported', got '%s'", archSupport.PPC64LE)
	}

	if len(archSupport.SupportedArchs) != 3 {
		t.Errorf("Expected 3 supported architectures, got %d", len(archSupport.SupportedArchs))
	}
}

func TestParseOperatorIndexNonExistentFile(t *testing.T) {
	_, err := ParseOperatorIndex("nonexistent-file.json")
	if err == nil {
		t.Error("Expected error when parsing nonexistent file, got nil")
	}
}

func TestParseOperatorIndexEmptyFile(t *testing.T) {
	// Create a temporary empty file
	tmpFile, err := os.CreateTemp("", "test-index-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	results, err := ParseOperatorIndex(tmpFile.Name())
	if err != nil {
		t.Errorf("Expected no error for empty file, got: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty file, got %d", len(results))
	}
}

func TestExtractArchitectureSupport(t *testing.T) {
	// Test with a bundle that has no architecture labels
	entry := models.OperatorIndexEntry{
		Schema:  "olm.bundle",
		Name:    "test-bundle",
		Package: "test-package",
		Properties: []models.Property{
			{
				Type: "olm.package",
				Value: map[string]interface{}{
					"packageName": "test-package",
					"version":     "1.0.0",
				},
			},
		},
	}

	result := extractArchitectureSupport(entry)
	if result != nil {
		t.Error("Expected nil for bundle without architecture labels")
	}

	// Test with a bundle that has architecture labels
	entryWithArch := models.OperatorIndexEntry{
		Schema:  "olm.bundle",
		Name:    "test-bundle",
		Package: "test-package",
		Properties: []models.Property{
			{
				Type: "olm.csv.metadata",
				Value: map[string]interface{}{
					"labels": map[string]interface{}{
						"operatorframework.io/arch.amd64": "supported",
						"operatorframework.io/arch.arm64": "supported",
					},
				},
			},
		},
	}

	resultWithArch := extractArchitectureSupport(entryWithArch)
	if resultWithArch == nil {
		t.Error("Expected non-nil result for bundle with architecture labels")
	} else {
		if resultWithArch.AMD64 != "supported" {
			t.Errorf("Expected AMD64 to be 'supported', got '%s'", resultWithArch.AMD64)
		}
		if len(resultWithArch.SupportedArchs) != 2 {
			t.Errorf("Expected 2 supported architectures, got %d", len(resultWithArch.SupportedArchs))
		}
	}
}

