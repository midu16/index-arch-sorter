package models

// OperatorIndexEntry represents a single entry in the operator index
type OperatorIndexEntry struct {
	Schema        string         `json:"schema"`
	Name          string         `json:"name"`
	Package       string         `json:"package"`
	Image         string         `json:"image,omitempty"`
	Properties    []Property     `json:"properties,omitempty"`
	RelatedImages []RelatedImage `json:"relatedImages,omitempty"`
}

// Property represents a property of a bundle
type Property struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// RelatedImage represents a related image in a bundle
type RelatedImage struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// ArchitectureSupport represents the architecture support for an operator
type ArchitectureSupport struct {
	OperatorName   string   `json:"operator_name"`
	BundleName     string   `json:"bundle_name"`
	Package        string   `json:"package"`
	Image          string   `json:"image"`
	SupportedArchs []string `json:"supported_architectures"`
	AMD64          string   `json:"amd64,omitempty"`
	ARM64          string   `json:"arm64,omitempty"`
	PPC64LE        string   `json:"ppc64le,omitempty"`
	S390X          string   `json:"s390x,omitempty"`
}
