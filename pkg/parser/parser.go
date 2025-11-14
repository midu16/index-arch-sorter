package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/midu/index-arch-sorter/pkg/models"
)

// ParseOperatorIndex parses the operator index file and extracts architecture support information
func ParseOperatorIndex(filePath string) ([]models.ArchitectureSupport, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var results []models.ArchitectureSupport
	decoder := json.NewDecoder(file)

	entryCount := 0
	for {
		var entry models.OperatorIndexEntry
		err := decoder.Decode(&entry)
		if err != nil {
			// Check if we've reached the end of the file
			if err == io.EOF {
				break
			}
			// Check for end of input in error message
			if strings.Contains(err.Error(), "EOF") {
				break
			}
			fmt.Printf("Warning: Failed to parse entry %d: %v\n", entryCount+1, err)
			// Try to continue parsing
			continue
		}

		entryCount++

		// We only care about bundles
		if entry.Schema != "olm.bundle" {
			continue
		}

		archSupport := extractArchitectureSupport(entry)
		if archSupport != nil {
			results = append(results, *archSupport)
		}
	}

	fmt.Printf("Processed %d entries\n", entryCount)
	return results, nil
}

// extractArchitectureSupport extracts architecture support from a bundle entry
func extractArchitectureSupport(entry models.OperatorIndexEntry) *models.ArchitectureSupport {
	archSupport := &models.ArchitectureSupport{
		OperatorName:   entry.Name,
		BundleName:     entry.Name,
		Package:        entry.Package,
		Image:          entry.Image,
		SupportedArchs: []string{},
	}

	// Look for architecture labels in properties
	for _, prop := range entry.Properties {
		if prop.Type == "olm.csv.metadata" {
			// Extract CSV metadata
			if propMap, ok := prop.Value.(map[string]interface{}); ok {
				if labels, hasLabels := propMap["labels"].(map[string]interface{}); hasLabels {
					extractArchLabels(labels, archSupport)
				}
			}
		}
	}

	// Only return if we found at least one architecture
	if len(archSupport.SupportedArchs) > 0 {
		return archSupport
	}

	return nil
}

// extractArchLabels extracts architecture labels and populates the ArchitectureSupport struct
func extractArchLabels(labels map[string]interface{}, archSupport *models.ArchitectureSupport) {
	archMap := map[string]*string{
		"operatorframework.io/arch.amd64":   &archSupport.AMD64,
		"operatorframework.io/arch.arm64":   &archSupport.ARM64,
		"operatorframework.io/arch.ppc64le": &archSupport.PPC64LE,
		"operatorframework.io/arch.s390x":   &archSupport.S390X,
	}

	archNames := map[string]string{
		"operatorframework.io/arch.amd64":   "amd64",
		"operatorframework.io/arch.arm64":   "arm64",
		"operatorframework.io/arch.ppc64le": "ppc64le",
		"operatorframework.io/arch.s390x":   "s390x",
	}

	for label, field := range archMap {
		if value, exists := labels[label]; exists {
			if strValue, ok := value.(string); ok {
				*field = strValue
				if strValue == "supported" {
					archSupport.SupportedArchs = append(archSupport.SupportedArchs, archNames[label])
				}
			}
		}
	}
}
