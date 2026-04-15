package config

// Config holds the generator configuration loaded from config.yaml
type Config struct {
	Resources   map[string]ResourceConfig   `yaml:"resources"`
	DataSources map[string]DataSourceConfig `yaml:"datasources"`
}

// ResourceConfig holds the CRUD operation config for a single resource
type ResourceConfig struct {
	APITag       string           `yaml:"api_tag"`
	ParentParams []ParentParam    `yaml:"parent_params"`
	Variants     []VariantConfig  `yaml:"variants"`
	Read         *OperationConfig `yaml:"read"`
	Create       *OperationConfig `yaml:"create"`
	Update       *OperationConfig `yaml:"update"`
	Delete       *OperationConfig `yaml:"delete"`
	// CreateIDField: when the Create operation returns no body (e.g. 204 No Content),
	// set the resource ID from this plan field name (snake_case TF attr) instead of result.GetId().
	// Example: "member_external_id"
	CreateIDField string `yaml:"create_id_field"`
	// ReadListIDField: when the Read operation returns a list (not a single item),
	// search this array field name in the response for the resource ID to verify existence.
	// Example: "member_external_ids"
	ReadListIDField string `yaml:"read_list_id_field"`
}

// DataSourceConfig holds the singular/plural fetch config for a data source
type DataSourceConfig struct {
	APITag       string           `yaml:"api_tag"`
	ParentParams []ParentParam    `yaml:"parent_params"`
	Variants     []VariantConfig  `yaml:"variants"`
	Singular     *OperationConfig `yaml:"singular"`
	Plural       *OperationConfig `yaml:"plural"`
}

// VariantConfig describes one concrete variant of a polymorphic (oneOf) resource.
// When variants are defined, the generator emits one resource/datasource per variant
// instead of a single merged resource.
type VariantConfig struct {
	// Suffix is appended to the base resource name with underscore, e.g. "saml" → "okta_application_saml".
	Suffix string `yaml:"suffix"`
	// SchemaRef is the component schema name for this variant, e.g. "SamlApplication".
	SchemaRef string `yaml:"schema_ref"`
	// DiscriminatorValue is the value of the discriminator field for this variant, e.g. "SAML_2_0".
	DiscriminatorValue string `yaml:"discriminator_value"`
	// DiscriminatorField is the property name used to distinguish variants (e.g. "signOnMode", "type").
	DiscriminatorField string `yaml:"discriminator_field"`
}

// OperationConfig holds the HTTP method and path for a single operation
type OperationConfig struct {
	Method string `yaml:"method"`
	Path   string `yaml:"path"`
}

// ParentParam describes a required parent resource ID parameter (for nested resources)
type ParentParam struct {
	Name        string `yaml:"name"`        // e.g. "app_id"
	Description string `yaml:"description"` // e.g. "The ID of the parent application"
	PathParam   string `yaml:"path_param"`  // e.g. "{appId}" as it appears in the path
}
