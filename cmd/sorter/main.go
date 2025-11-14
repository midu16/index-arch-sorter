package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/midu/index-arch-sorter/pkg/models"
	"github.com/midu/index-arch-sorter/pkg/parser"
)

// Valid architectures
var validArchitectures = []string{"amd64", "arm64", "ppc64le", "s390x"}

func main() {
	// Parse command line flags
	filePath := flag.String("file", "redhat-operator-index.json", "Path to the operator index JSON file")
	outputFormat := flag.String("format", "table", "Output format: table, json, or csv")
	filterArch := flag.String("arch", "", "Filter by architecture: amd64, arm64, ppc64le, s390x")
	filterOperator := flag.String("operator", "", "Filter by operator/package name (partial match)")
	flag.StringVar(filterOperator, "operator-name", "", "Filter by operator/package name (alias for --operator)")
	listArchs := flag.Bool("list-archs", false, "List all available architectures and exit")
	flag.Parse()

	// Handle list architectures flag
	if *listArchs {
		fmt.Println("Available architectures:")
		for _, arch := range validArchitectures {
			fmt.Printf("  - %s\n", arch)
		}
		os.Exit(0)
	}

	// Validate architecture filter if provided
	if *filterArch != "" {
		if !isValidArchitecture(*filterArch) {
			fmt.Fprintf(os.Stderr, "Error: Invalid architecture '%s'\n", *filterArch)
			fmt.Fprintf(os.Stderr, "Valid architectures are: %s\n", strings.Join(validArchitectures, ", "))
			os.Exit(1)
		}
	}

	// Parse the operator index
	fmt.Printf("Parsing operator index from: %s\n", *filePath)
	results, err := parser.ParseOperatorIndex(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing operator index: %v\n", err)
		os.Exit(1)
	}

	// Filter by operator name if specified
	if *filterOperator != "" {
		fmt.Printf("Filtering by operator name: %s\n", *filterOperator)
		results = filterByOperatorName(results, *filterOperator)
		fmt.Printf("Found %d bundles matching '%s'\n", len(results), *filterOperator)
	}

	// Filter by architecture if specified
	if *filterArch != "" {
		fmt.Printf("Filtering by architecture: %s\n", strings.ToUpper(*filterArch))
		beforeArchFilter := len(results)
		results = filterByArchitecture(results, *filterArch)
		fmt.Printf("Found %d operator bundles supporting %s (out of %d bundles)\n\n",
			len(results), strings.ToUpper(*filterArch), beforeArchFilter)
	} else if *filterOperator == "" {
		fmt.Printf("Found %d operator bundles with architecture information\n\n", len(results))
	}

	// Sort results by package name for consistent output
	sort.Slice(results, func(i, j int) bool {
		if results[i].Package == results[j].Package {
			return results[i].BundleName < results[j].BundleName
		}
		return results[i].Package < results[j].Package
	})

	// Output results in the specified format
	switch *outputFormat {
	case "json":
		outputJSON(results)
	case "csv":
		outputCSV(results)
	case "table":
		fallthrough
	default:
		outputTable(results)
	}

	// Print summary statistics
	printSummary(results, *filterArch)
}

// isValidArchitecture checks if the provided architecture is valid
func isValidArchitecture(arch string) bool {
	for _, valid := range validArchitectures {
		if strings.ToLower(arch) == valid {
			return true
		}
	}
	return false
}

// filterByOperatorName filters results to only include bundles matching the operator name
// The filter does a case-insensitive partial match on both package name and bundle name
func filterByOperatorName(results []models.ArchitectureSupport, operatorName string) []models.ArchitectureSupport {
	var filtered []models.ArchitectureSupport
	operatorName = strings.ToLower(operatorName)

	for _, result := range results {
		packageName := strings.ToLower(result.Package)
		bundleName := strings.ToLower(result.BundleName)

		// Match if the operator name is found in either package or bundle name
		if strings.Contains(packageName, operatorName) || strings.Contains(bundleName, operatorName) {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

// filterByArchitecture filters results to only include bundles supporting the specified architecture
func filterByArchitecture(results []models.ArchitectureSupport, arch string) []models.ArchitectureSupport {
	var filtered []models.ArchitectureSupport
	arch = strings.ToLower(arch)

	for _, result := range results {
		supported := false
		switch arch {
		case "amd64":
			supported = result.AMD64 == "supported"
		case "arm64":
			supported = result.ARM64 == "supported"
		case "ppc64le":
			supported = result.PPC64LE == "supported"
		case "s390x":
			supported = result.S390X == "supported"
		}

		if supported {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

func outputTable(results []models.ArchitectureSupport) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PACKAGE\tBUNDLE NAME\tAMD64\tARM64\tPPC64LE\tS390X")
	fmt.Fprintln(w, "-------\t-----------\t-----\t-----\t-------\t-----")

	for _, result := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			result.Package,
			result.BundleName,
			formatSupport(result.AMD64),
			formatSupport(result.ARM64),
			formatSupport(result.PPC64LE),
			formatSupport(result.S390X),
		)
	}
	w.Flush()
}

func outputJSON(results []models.ArchitectureSupport) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func outputCSV(results []models.ArchitectureSupport) {
	fmt.Println("package,bundle_name,amd64,arm64,ppc64le,s390x")
	for _, result := range results {
		fmt.Printf("%s,%s,%s,%s,%s,%s\n",
			result.Package,
			result.BundleName,
			result.AMD64,
			result.ARM64,
			result.PPC64LE,
			result.S390X,
		)
	}
}

func formatSupport(support string) string {
	if support == "supported" {
		return "✓"
	} else if support == "" {
		return "-"
	}
	return support
}

func printSummary(results []models.ArchitectureSupport, filterArch string) {
	amd64Count := 0
	arm64Count := 0
	ppc64leCount := 0
	s390xCount := 0

	packageMap := make(map[string]bool)

	for _, result := range results {
		packageMap[result.Package] = true

		if result.AMD64 == "supported" {
			amd64Count++
		}
		if result.ARM64 == "supported" {
			arm64Count++
		}
		if result.PPC64LE == "supported" {
			ppc64leCount++
		}
		if result.S390X == "supported" {
			s390xCount++
		}
	}

	fmt.Printf("\n--- Summary ---\n")
	if filterArch != "" {
		fmt.Printf("Filtered by architecture: %s\n", strings.ToUpper(filterArch))
	}
	fmt.Printf("Total unique packages: %d\n", len(packageMap))
	fmt.Printf("Total bundles: %d\n", len(results))

	if filterArch == "" {
		fmt.Printf("\nArchitecture Support:\n")
		fmt.Printf("  AMD64:   %d bundles\n", amd64Count)
		fmt.Printf("  ARM64:   %d bundles\n", arm64Count)
		fmt.Printf("  PPC64LE: %d bundles\n", ppc64leCount)
		fmt.Printf("  S390X:   %d bundles\n", s390xCount)
	} else {
		fmt.Printf("\nNote: All bundles shown support %s architecture\n", strings.ToUpper(filterArch))
		if filterArch != "amd64" && amd64Count > 0 {
			fmt.Printf("  Also support AMD64:   %d bundles\n", amd64Count)
		}
		if filterArch != "arm64" && arm64Count > 0 {
			fmt.Printf("  Also support ARM64:   %d bundles\n", arm64Count)
		}
		if filterArch != "ppc64le" && ppc64leCount > 0 {
			fmt.Printf("  Also support PPC64LE: %d bundles\n", ppc64leCount)
		}
		if filterArch != "s390x" && s390xCount > 0 {
			fmt.Printf("  Also support S390X:   %d bundles\n", s390xCount)
		}
	}
}
