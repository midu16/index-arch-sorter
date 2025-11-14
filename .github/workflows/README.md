# GitHub Workflows

## build-and-analyze.yml

This workflow automates the complete operator architecture analysis pipeline.

### What it does

1. **Test** - Runs Go tests to ensure code quality
2. **Download OPM** - Downloads the OpenShift `opm` CLI tool (with caching)
3. **Generate Index** - Pulls and renders the Red Hat operator index
4. **Build** - Compiles the Go application
5. **Analyze** - Generates comprehensive reports in multiple formats
6. **Archive** - Uploads all outputs as GitHub artifacts

### Triggers

- **Push/PR** to main/master branches
- **Manual dispatch** with custom parameters (index version and image)
- **Scheduled** weekly runs (Mondays at 2 AM UTC)

### Outputs

The workflow produces two artifact packages:

#### 1. operator-arch-analysis-{version}-{run_number}

Contains:
- `operator-arch-analysis.json` - Complete analysis in JSON
- `operator-arch-analysis.csv` - Complete analysis in CSV
- `operator-summary.txt` - Human-readable summary
- `architectures/` directory with per-architecture reports:
  - `{arch}-operators.json` - JSON format
  - `{arch}-operators.csv` - CSV format
  
Supported architectures: amd64, arm64, ppc64le, s390x

#### 2. redhat-operator-index-{version}-{run_number}

Contains:
- `redhat-operator-index.json` - Raw operator index from Red Hat registry

### Manual Trigger

To manually trigger the workflow with custom parameters:

1. Go to Actions tab in GitHub
2. Select "Build and Analyze Operator Index"
3. Click "Run workflow"
4. Enter custom values (optional):
   - **index_version**: e.g., `v4.20`, `v4.19`
   - **index_image**: e.g., `registry.redhat.io/redhat/redhat-operator-index:v4.20`

### Environment Variables

- `GO_VERSION`: Go version to use (default: 1.21)
- `INDEX_VERSION`: Operator index version (default: v4.20)
- `INDEX_IMAGE`: Full image path (default: registry.redhat.io/redhat/redhat-operator-index:v4.20)

### Caching

The workflow uses caching to speed up runs:
- **Go modules** - Cached based on go.sum hash
- **OPM binary** - Cached based on INDEX_VERSION

### Viewing Results

After a workflow run completes:

1. Go to the workflow run page
2. Scroll to "Artifacts" section at the bottom
3. Download the artifact packages
4. Extract and explore the reports

The workflow also adds a summary to the GitHub Actions run page showing key statistics.

