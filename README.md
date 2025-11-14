# Index Architecture Sorter

A Go-based tool to analyze Red Hat Operator Index files and determine architecture support for each operator bundle.

## Overview

This tool parses the Red Hat operator index JSON file and extracts architecture support information from operator bundles. It identifies which architectures (amd64, arm64, ppc64le, s390x) each operator supports based on the labels in the bundle metadata.

## Project Structure

```
index-arch-sorter/
├── cmd/
│   └── sorter/
│       └── main.go              # Main application entry point
├── pkg/
│   ├── models/
│   │   └── operator.go          # Data structures
│   └── parser/
│       └── parser.go            # Parser logic
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## Building

To build the application:

```bash
go build -o index-arch-sorter ./cmd/sorter
```

## Usage

### Basic Usage

```bash
./bin/index-arch-sorter -file redhat-operator-index.json
```

### Command Line Options

- `-file <path>`: Path to the operator index JSON file (default: "redhat-operator-index.json")
- `-format <format>`: Output format - `table`, `json`, or `csv` (default: "table")
- `-arch <architecture>`: Filter results to show only operators supporting the specified architecture - `amd64`, `arm64`, `ppc64le`, or `s390x`
- `-operator <name>` or `-operator-name <name>`: Filter by operator/package name (partial, case-insensitive match)
- `-list-archs`: List all available architectures and exit

### Examples

1. **Display results in table format** (default):
   ```bash
   ./bin/index-arch-sorter -file redhat-operator-index.json
   ```

2. **List available architectures**:
   ```bash
   ./bin/index-arch-sorter -list-archs
   ```

3. **Filter by specific architecture** (shows only operators supporting that architecture):
   ```bash
   ./bin/index-arch-sorter -file redhat-operator-index.json -arch arm64
   ```

4. **Filter ARM64 operators and output as JSON**:
   ```bash
   ./bin/index-arch-sorter -file redhat-operator-index.json -arch arm64 -format json
   ```

5. **Output as CSV**:
   ```bash
   ./bin/index-arch-sorter -file redhat-operator-index.json -format csv > operators.csv
   ```

6. **Find all operators supporting s390x**:
   ```bash
   ./bin/index-arch-sorter -arch s390x -format csv > s390x-operators.csv
   ```

7. **Filter by operator name** (partial match):
   ```bash
   ./bin/index-arch-sorter -operator amq
   ```

8. **Combine operator name and architecture filters**:
   ```bash
   ./bin/index-arch-sorter -operator cluster -arch arm64
   ```

9. **Find specific operator supporting a specific architecture**:
   ```bash
   ./bin/index-arch-sorter -operator-name "advanced-cluster-management" -arch ppc64le -format json
   ```

10. **Search for broker operators**:
    ```bash
    ./bin/index-arch-sorter -operator broker
    ```

11. **Using the Makefile**:
    ```bash
    # Show all available targets
    make help
    
    # Build and run
    make run
    
    # Filter by architecture
    make run-arch ARCH=arm64
    
    # Filter by operator name
    make run-operator OPERATOR=amq
    
    # Filter by both operator and architecture
    make run-filter OPERATOR=cluster ARCH=arm64
    
    # Generate the index and run
    make full
    ```

## Output Format

### Table Format

The table format displays operator bundles with checkmarks (✓) for supported architectures:

```
PACKAGE               BUNDLE NAME                    AMD64  ARM64  PPC64LE  S390X
-------               -----------                    -----  -----  -------  -----
3scale-operator       3scale-operator.v0.10.0-mas    ✓      -      ✓        ✓
```

### JSON Format

The JSON format provides detailed information in a structured format suitable for programmatic processing:

```json
[
  {
    "operator_name": "3scale-operator.v0.10.0-mas",
    "bundle_name": "3scale-operator.v0.10.0-mas",
    "package": "3scale-operator",
    "image": "registry.redhat.io/3scale-mas/3scale-operator-bundle@sha256:...",
    "supported_architectures": ["amd64", "ppc64le", "s390x"],
    "amd64": "supported",
    "ppc64le": "supported",
    "s390x": "supported"
  }
]
```

### CSV Format

The CSV format is suitable for importing into spreadsheet applications:

```csv
package,bundle_name,amd64,arm64,ppc64le,s390x
3scale-operator,3scale-operator.v0.10.0-mas,supported,,supported,supported
```

## Summary Statistics

The tool also provides summary statistics including:
- Total number of unique packages
- Total number of bundles analyzed
- Number of bundles supporting each architecture

## Architecture Labels

The tool looks for the following labels in operator bundle metadata:

- `operatorframework.io/arch.amd64`
- `operatorframework.io/arch.arm64`
- `operatorframework.io/arch.ppc64le`
- `operatorframework.io/arch.s390x`

Each label can have the value `"supported"` to indicate that the operator bundle supports that architecture.

## Makefile Targets

The project includes a comprehensive Makefile for common tasks:

- `make build` - Build the application
- `make run` - Run the application with the default index file
- `make run-arch ARCH=<arch>` - Filter by architecture
- `make run-operator OPERATOR=<name>` - Filter by operator name
- `make run-filter OPERATOR=<name> ARCH=<arch>` - Combine both filters
- `make run-json` - Output in JSON format
- `make run-csv` - Output in CSV format
- `make download-opm` - Download the opm CLI tool
- `make generate-index` - Generate the operator index JSON using opm
- `make test` - Run tests
- `make clean` - Clean build artifacts
- `make help` - Show all available targets

## Requirements

- Go 1.21 or later
- The operator index JSON file (newline-delimited JSON format)
- Optional: Docker/Podman for pulling operator index images (when using `make generate-index`)

## CI/CD

The project includes a GitHub Actions workflow (`.github/workflows/build-and-analyze.yml`) that:

1. Runs tests
2. Downloads the opm CLI tool
3. Generates the operator index from the Red Hat registry
4. Builds the application
5. Generates comprehensive reports in multiple formats (JSON, CSV, plain text)
6. Creates architecture-specific reports for amd64, arm64, ppc64le, and s390x
7. Uploads all outputs as GitHub artifacts

### Workflow Triggers

- **Push** to main/master branches
- **Pull requests** to main/master branches
- **Manual dispatch** with configurable index version and image
- **Scheduled** runs weekly on Mondays at 2 AM UTC

### Artifacts

The workflow produces the following artifacts:

- `operator-arch-analysis-<version>-<run_number>` - Complete analysis in JSON, CSV, and text formats, plus architecture-specific reports
- `redhat-operator-index-<version>-<run_number>` - Raw operator index JSON file

## Testing

Run tests with:

```bash
make test
```

Tests cover the data models and parser functionality.

## License

MIT License

