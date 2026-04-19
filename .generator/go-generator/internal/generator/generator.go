package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/okta/terraform-provider-okta/generator/internal/config"
	"github.com/okta/terraform-provider-okta/generator/internal/formatter"
	"github.com/okta/terraform-provider-okta/generator/internal/openapi"
	oktypes "github.com/okta/terraform-provider-okta/generator/internal/types"
)

// TemplateData is the data passed to all Go templates
type TemplateData struct {
	Name      string // lowerCamel, e.g. "group"
	TitleName string // TitleCase, e.g. "Group"
	TFName    string // snake_case terraform name, e.g. "group"
	APITag    string // e.g. "Group"
	// AllNestedStructDefs contains all nested Go struct type definitions for the resource,
	// computed in one pass with a shared dedup map to prevent redeclarations.
	AllNestedStructDefs string
	Properties          []PropData
	// RequestBodyFields is the subset of Properties that are settable (non-computed, non-id).
	// Used to generate the request body construction in Create/Update.
	RequestBodyFields []PropData
	ParentParams      []ParentParamData
	HasParent         bool
	HasCreate         bool
	HasUpdate         bool
	HasDelete         bool
	ReadMethod        string
	CreateMethod      string
	UpdateMethod      string
	DeleteMethod      string
	ListMethod        string
	// IDField is the Go struct field name of the ID property in the API response model.
	// Defaults to "Id" when the response schema has an "id" property.
	IDField string
	// BodySetterMethod is the CamelCased name of the SDK request builder method
	// that sets the Create request body, e.g. "Rule", "FederatedClaimRequestBody".
	// Derived from x-codegen-request-body-name → request body $ref name → "Body".
	BodySetterMethod string
	// UpdateBodySetterMethod is the SDK builder method name for the Update request body.
	// Differs from BodySetterMethod when POST and PUT use different body schemas
	// (e.g. POST uses FederatedClaimRequestBody, PUT uses FederatedClaim).
	// Falls back to BodySetterMethod when Update has no distinct request body schema.
	UpdateBodySetterMethod string
	// SDKTypeName is the Go SDK struct name used for Create request body construction.
	// For normal resources it equals TitleName. For variants it is the DiscriminatorSchemaRef
	// (e.g. "BehaviorRuleASN") which matches the actual okta.New<SDKTypeName>WithDefaults() constructor.
	SDKTypeName string
	// UpdateSDKTypeName is the Go SDK struct name for Update request body construction.
	// Differs from SDKTypeName when the PUT request body uses a different schema than POST
	// (e.g. TrustedOrigin for PUT vs TrustedOriginWrite for POST).
	// Falls back to SDKTypeName when the Update op has no distinct $ref.
	UpdateSDKTypeName string
	// UpdateRequestBodyFields are the scalar writable fields for the Update request body.
	// Derived from the Update op's request body schema — may differ from RequestBodyFields
	// when Create and Update use different schemas (e.g. AppServiceAccount vs AppServiceAccountForUpdate).
	// Falls back to RequestBodyFields when Update has no distinct schema.
	UpdateRequestBodyFields []PropData
	// Discriminator fields — set when this is one variant of a polymorphic (oneOf) resource
	DiscriminatorField     string // e.g. "signOnMode", "type" (raw spec name)
	DiscriminatorTFAttr    string // e.g. "sign_on_mode", "type" (snake_case for tfsdk tags)
	DiscriminatorValue     string // e.g. "SAML_2_0", "ACCESS_POLICY"
	DiscriminatorSchemaRef string // e.g. "SamlApplication"
	IsVariant              bool   // true when this TemplateData came from a variants expansion
	// UnionTypeName is the SDK union wrapper type for polymorphic (oneOf) responses.
	// e.g. "ListLogStreams200ResponseInner" — returned by GetLogStream / CreateLogStream Execute().
	// Used to unwrap the concrete variant: result.<DiscriminatorSchemaRef> and to
	// build the body wrapper: okta.<SDKTypeName>As<UnionTypeName>(body).
	UnionTypeName string
	// CreateIDField: when set, Create does not call result.GetId() (no response body).
	// Instead, the resource ID is taken from this plan field (Go struct field name, e.g. "MemberExternalId").
	CreateIDField string
	// HasComputedProperties is true when at least one property is Computed (read-only from the API).
	// Used by the template to decide whether to emit a follow-up Read after a 204 Create.
	HasComputedProperties bool
	// ReadListIDField: when set, Read calls the list endpoint and searches this response field
	// (Go struct field name, e.g. "MemberExternalIds") for the resource ID to verify existence.
	// The list field must be a []string (or similar) in the SDK response struct.
	ReadListIDField string
	// ReadNoop: when true, the Read function does not call any API — it preserves existing state unchanged.
	// Used for POST-only resources where there is no GET endpoint.
	ReadNoop bool
	// DeleteNoop: when true, the Delete function does not call any API.
	// Used for resources that cannot be deleted via the API.
	DeleteNoop bool
}

// ParentParamData holds a rendered parent param for templates
type ParentParamData struct {
	GoField     string // e.g. "AppID"
	TFAttr      string // e.g. "app_id"
	Description string // e.g. "The ID of the parent application"
	PathParam   string // e.g. "{appId}"
}

// PropData holds one schema property for template rendering
type PropData struct {
	GoField           string
	GoType            string
	TFAttr            string
	TFSchemaType      string
	ElementType       string     // non-empty for flat array types: the attr.Type expression for ElementType
	IsListNested      bool       // true when the array items are objects → use schema.ListNestedAttribute
	IsListScalar      bool       // true when the array items are scalars ([]string, []int64) — uses schema.ListAttribute + ElementType
	NestedProps       []PropData // non-empty when property is an inline object/list-of-objects — use SingleNestedAttribute or ListNestedAttribute
	NestedModelName   string     // e.g. "DevicePostureCheckRemediationSettingsModel" — the generated nested struct name
	SDKNestedTypeName string     // e.g. "IdentitySourceGroupProfileForUpsert" — the SDK struct to construct for nested body fields; derived from $ref
	TFPlanNestedProps []PropData // when this is a nested request-body field, holds the *response-schema* nested props so the template knows the plan access path (e.g. plan.Profile.Profile.DisplayName vs plan.Profile.DisplayName)
	SchemaAttrBlock   string     // pre-rendered full schema attribute Go literal (handles arbitrary depth)
	NestedStructDefs  string     // pre-rendered nested struct type definitions for all sub-models
	Description       string
	Required          bool
	Computed          bool
	IsDateTime        bool // true when OAS format=date-time; SDK getter returns time.Time, needs .Format(time.RFC3339)
	WriteOnly         bool // true when field exists only in the request body schema, not in the response — never call Get<Field>() on result
}

// Generator holds templates and spec
type Generator struct {
	spec      *openapi.Spec
	templates *template.Template
	outputDir string
	goFmt     bool
	log       *log.Logger // nil = silent
}

// New creates a new Generator
func New(spec *openapi.Spec, templatesDir, outputDir string, goFmt bool, logger *log.Logger) (*Generator, error) {
	tmpl := template.New("").Funcs(template.FuncMap{
		"lower": strings.ToLower,
		"title": strings.Title, //nolint:staticcheck
		"add":   func(a, b int) int { return a + b },
	})

	pattern := filepath.Join(templatesDir, "*.tmpl")
	tmpl, err := tmpl.ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("parsing templates from %s: %w", pattern, err)
	}

	if logger != nil {
		for _, t := range tmpl.Templates() {
			logger.Printf("Loaded template: %s", t.Name())
		}
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("creating output dir: %w", err)
	}

	return &Generator{
		spec:      spec,
		templates: tmpl,
		outputDir: outputDir,
		goFmt:     goFmt,
		log:       logger,
	}, nil
}

// GenerateResource generates resource .go and _test.go files for a named resource.
// For polymorphic (oneOf) resources with variants defined, one file set is emitted per variant.
func (g *Generator) GenerateResource(name string, cfg config.ResourceConfig) error {
	dataList := g.buildResourceDataList(name, cfg)

	for _, data := range dataList {
		if g.log != nil {
			g.log.Printf("─── resource: %s ───────────────────────────────", data.TFName)
			g.log.Printf("  APITag      : %s", data.APITag)
			g.log.Printf("  HasParent   : %v", data.HasParent)
			for i, pp := range data.ParentParams {
				g.log.Printf("  ParentParam[%d]: GoField=%s  TFAttr=%s  PathParam=%s", i, pp.GoField, pp.TFAttr, pp.PathParam)
			}
			g.log.Printf("  HasUpdate   : %v  HasDelete : %v", data.HasUpdate, data.HasDelete)
			g.log.Printf("  ReadMethod  : %s", data.ReadMethod)
			g.log.Printf("  CreateMethod: %s", data.CreateMethod)
			if data.DiscriminatorField != "" {
				g.log.Printf("  Discriminator: field=%s value=%s schemaRef=%s",
					data.DiscriminatorField, data.DiscriminatorValue, data.DiscriminatorSchemaRef)
			}
			g.log.Printf("  Properties  : %d fields", len(data.Properties))
			for _, p := range data.Properties {
				g.log.Printf("    %-30s %-15s required=%-5v computed=%v", p.TFAttr, p.GoType, p.Required, p.Computed)
			}
		}

		if err := g.renderToFile("resource.go.tmpl",
			filepath.Join(g.outputDir, fmt.Sprintf("resource_okta_%s_generated.go", data.TFName)),
			data,
		); err != nil {
			return fmt.Errorf("generating resource %s: %w", data.TFName, err)
		}

		if err := g.renderToFile("resource_test.go.tmpl",
			filepath.Join(g.outputDir, fmt.Sprintf("resource_okta_%s_generated_test.go", data.TFName)),
			data,
		); err != nil {
			return fmt.Errorf("generating resource test %s: %w", data.TFName, err)
		}
	}
	return nil
}

// GenerateDataSource generates a data source .go file for a named data source.
// For polymorphic (oneOf) resources with variants defined, one file is emitted per variant.
func (g *Generator) GenerateDataSource(name string, cfg config.DataSourceConfig) error {
	dataList := g.buildDataSourceDataList(name, cfg)

	for _, data := range dataList {
		if g.log != nil {
			g.log.Printf("─── datasource: %s ──────────────────────────────", data.TFName)
			g.log.Printf("  APITag    : %s", data.APITag)
			g.log.Printf("  HasParent : %v", data.HasParent)
			for i, pp := range data.ParentParams {
				g.log.Printf("  ParentParam[%d]: GoField=%s  TFAttr=%s  PathParam=%s", i, pp.GoField, pp.TFAttr, pp.PathParam)
			}
			if data.DiscriminatorField != "" {
				g.log.Printf("  Discriminator: field=%s value=%s schemaRef=%s",
					data.DiscriminatorField, data.DiscriminatorValue, data.DiscriminatorSchemaRef)
			}
			g.log.Printf("  ReadMethod: %s  ListMethod: %s", data.ReadMethod, data.ListMethod)
			g.log.Printf("  Properties: %d fields", len(data.Properties))
			for _, p := range data.Properties {
				g.log.Printf("    %-30s %-15s computed=%v", p.TFAttr, p.GoType, p.Computed)
			}
		}

		if err := g.renderToFile("data_source.go.tmpl",
			filepath.Join(g.outputDir, fmt.Sprintf("data_source_okta_%s_generated.go", data.TFName)),
			data,
		); err != nil {
			return fmt.Errorf("generating datasource %s: %w", data.TFName, err)
		}
	}
	return nil
}

// buildResourceDataList returns one TemplateData per resource to generate.
// For polymorphic resources with cfg.Variants, this returns one entry per variant.
// If no variants are configured but the response schema has a discriminator.mapping,
// variants are auto-expanded from the mapping (Issue 2 fix).
// For normal resources it returns a single-element slice.
func (g *Generator) buildResourceDataList(name string, cfg config.ResourceConfig) []TemplateData {
	if len(cfg.Variants) == 0 {
		// Issue 2: auto-detect variants from discriminator.mapping when no variants configured
		if autoVariants := g.autoDetectVariants(name, cfg); len(autoVariants) > 0 {
			if g.log != nil {
				g.log.Printf("  [buildResourceDataList] auto-detected %d variants from discriminator.mapping for %s", len(autoVariants), name)
			}
			cfgCopy := cfg
			cfgCopy.Variants = autoVariants
			return g.buildResourceDataList(name, cfgCopy)
		}
		return []TemplateData{g.buildResourceData(name, cfg)}
	}

	var out []TemplateData
	for _, v := range cfg.Variants {
		variantName := name + "_" + v.Suffix
		title := formatter.CamelCase(variantName)

		if g.log != nil {
			g.log.Printf("  [buildResourceDataList] variant %s: schemaRef=%s discriminatorValue=%s",
				variantName, v.SchemaRef, v.DiscriminatorValue)
		}

		var props []PropData
		if v.SchemaRef != "" {
			sc := g.spec.GetSchemaByRef(v.SchemaRef)
			if sc == nil {
				if g.log != nil {
					g.log.Printf("  [buildResourceDataList] WARNING: schemaRef %q not found in spec", v.SchemaRef)
				}
			} else {
				if g.log != nil {
					g.log.Printf("  [buildResourceDataList] resolved schemaRef %q: type=%q allOf=%d props=%d",
						v.SchemaRef, sc.Type, len(sc.AllOf), len(sc.Properties))
				}
				props = g.schemaToProps(sc, false, title+"Model")
				// Remove the discriminator field — it is emitted separately in the template.
				if v.DiscriminatorField != "" {
					props = filterOutProp(props, formatter.TerraformAttrName(v.DiscriminatorField))
				}
			}
		}

		parentParams := buildParentParams(cfg.ParentParams, g.log)

		// Prefer operationId from the spec for accurate SDK method names.
		readOpID, createOpID, updateOpID, deleteOpID := "", "", "", ""
		if cfg.Read != nil {
			if op := g.spec.GetOperation(cfg.Read.Method, cfg.Read.Path); op != nil {
				readOpID = op.OperationID
			}
		}
		if cfg.Create != nil {
			if op := g.spec.GetOperation(cfg.Create.Method, cfg.Create.Path); op != nil {
				createOpID = op.OperationID
			}
		}
		if cfg.Update != nil {
			if op := g.spec.GetOperation(cfg.Update.Method, cfg.Update.Path); op != nil {
				updateOpID = op.OperationID
			}
		}
		if cfg.Delete != nil {
			if op := g.spec.GetOperation(cfg.Delete.Method, cfg.Delete.Path); op != nil {
				deleteOpID = op.OperationID
			}
		}

		readMethod := formatter.OperationIDToMethodName(readOpID, "get", name)
		createMethod := formatter.OperationIDToMethodName(createOpID, "create", name)
		updateMethod := formatter.OperationIDToMethodName(updateOpID, "update", name)
		deleteMethod := formatter.OperationIDToMethodName(deleteOpID, "delete", name)
		listMethod := formatter.ListAPIMethodName(name)

		// RequestBodyFields for variant: the request body is itself the oneOf union, not a concrete
		// schema $ref. So we derive body fields directly from the variant's own schema (v.SchemaRef)
		// — the same schema that was used to build props. Filter to non-computed scalar fields.
		var requestBodyFields []PropData
		for _, p := range props {
			if !p.Computed && isScalarGoType(p.GoType) {
				requestBodyFields = append(requestBodyFields, p)
			}
		}

		// UnionTypeName: the SDK wrapper type returned by Execute() for polymorphic responses.
		// e.g. "ListLogStreams200ResponseInner" for /api/v1/logStreams/{logStreamId}.
		unionTypeName := ""
		if cfg.Read != nil {
			unionTypeName = g.spec.GetUnionTypeName(cfg.Read.Path)
		}

		sdkTypeName := title
		if v.SchemaRef != "" {
			sdkTypeName = v.SchemaRef
		}

		// Derive body setter method from x-codegen-request-body-name
		variantBodySetter := "Body"
		if cfg.Create != nil {
			if op := g.spec.GetOperation(cfg.Create.Method, cfg.Create.Path); op != nil && op.RequestBodyName != "" {
				variantBodySetter = formatter.GoFieldName(op.RequestBodyName)
			}
		} else if cfg.Update != nil {
			if op := g.spec.GetOperation(cfg.Update.Method, cfg.Update.Path); op != nil && op.RequestBodyName != "" {
				variantBodySetter = formatter.GoFieldName(op.RequestBodyName)
			}
		}

		out = append(out, TemplateData{
			Name:                   formatter.LowerFirst(title),
			TitleName:              title,
			TFName:                 variantName,
			APITag:                 cfg.APITag,
			AllNestedStructDefs:    renderNestedStructs(props, map[string]bool{}),
			Properties:             props,
			RequestBodyFields:      requestBodyFields,
			ParentParams:           parentParams,
			HasParent:              len(parentParams) > 0,
			HasCreate:              cfg.Create != nil,
			HasUpdate:              cfg.Update != nil,
			HasDelete:              cfg.Delete != nil,
			ReadMethod:             readMethod,
			CreateMethod:           createMethod,
			UpdateMethod:           updateMethod,
			DeleteMethod:           deleteMethod,
			ListMethod:             listMethod,
			IDField:                "Id",
			SDKTypeName:            sdkTypeName,
			UpdateSDKTypeName:      sdkTypeName,
			BodySetterMethod:       variantBodySetter,
			DiscriminatorField:     v.DiscriminatorField,
			DiscriminatorTFAttr:    formatter.TerraformAttrName(v.DiscriminatorField),
			DiscriminatorValue:     v.DiscriminatorValue,
			DiscriminatorSchemaRef: v.SchemaRef,
			IsVariant:              true,
			UnionTypeName:          unionTypeName,
		})
	}
	return out
}

// autoDetectVariants inspects the response schema of the read/create operation.
// If the schema has a oneOf+discriminator.mapping at the root, it auto-builds
// VariantConfig entries from the mapping — one per discriminator value.
// Returns nil if no discriminator is found or schema cannot be resolved.
func (g *Generator) autoDetectVariants(name string, cfg config.ResourceConfig) []config.VariantConfig {
	// Try read op first, then create op
	var schema *openapi.Schema
	for _, opCfg := range []*config.OperationConfig{cfg.Read, cfg.Create} {
		if opCfg == nil {
			continue
		}
		op := g.spec.GetOperation(opCfg.Method, opCfg.Path)
		if op == nil {
			continue
		}

		//The reason Read is preferred over Create is that GET responses are the most complete —
		//they include all fields including readOnly ones (id, created, lastUpdated, etc.)
		//that the API echoes back, but you'd never send in a POST body.
		s := g.spec.GetResponseSchema(op)
		if s != nil {
			schema = s
			break
		}
	}
	if schema == nil {
		return nil
	}

	disc := openapi.GetDiscriminatorFromSchema(*schema)
	if disc == nil || len(disc.Mapping) == 0 {
		return nil
	}

	if g.log != nil {
		g.log.Printf("  [autoDetectVariants] found discriminator field=%q mapping=%d entries for %s",
			disc.PropertyName, len(disc.Mapping), name)
	}

	var variants []config.VariantConfig
	for discValue, schemaRefFull := range disc.Mapping {
		schemaRef := schemaRefFull
		// strip "#/components/schemas/" prefix if present
		if idx := len("#/components/schemas/"); len(schemaRefFull) > idx {
			schemaRef = schemaRefFull[idx:]
		}
		// derive suffix: replace any non-alphanumeric character with '_', then lowercase.
		// e.g. "token:hardware" → "token_hardware", "SAML_2_0" → "saml_2_0"
		sanitized := strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				return r
			}
			return '_'
		}, discValue)
		suffix := strings.ToLower(sanitized)
		variants = append(variants, config.VariantConfig{
			Suffix:             suffix,
			SchemaRef:          schemaRef,
			DiscriminatorValue: discValue,
			DiscriminatorField: disc.PropertyName,
		})
	}
	return variants
}

// buildDataSourceDataList mirrors buildResourceDataList for data sources.
func (g *Generator) buildDataSourceDataList(name string, cfg config.DataSourceConfig) []TemplateData {
	if len(cfg.Variants) == 0 {
		return []TemplateData{
			g.buildDataSourceData(name, cfg)}
	}

	var out []TemplateData
	for _, v := range cfg.Variants {
		variantName := name + "_" + v.Suffix
		title := formatter.CamelCase(variantName)

		if g.log != nil {
			g.log.Printf("  [buildDataSourceDataList] variant %s: schemaRef=%s", variantName, v.SchemaRef)
		}

		var props []PropData
		if v.SchemaRef != "" {
			sc := g.spec.GetSchemaByRef(v.SchemaRef)
			if sc != nil {
				props = g.schemaToProps(sc, true, title+"DataSourceModel")
				// Remove the discriminator field — it is emitted separately in the template.
				if v.DiscriminatorField != "" {
					props = filterOutProp(props, v.DiscriminatorField)
				}
			}
		}

		parentParams := buildParentParams(cfg.ParentParams, g.log)
		readMethod := formatter.APIMethodName("get", name)
		listMethod := formatter.ListAPIMethodName(name)

		out = append(out, TemplateData{
			Name:                   formatter.LowerFirst(title),
			TitleName:              title,
			TFName:                 variantName,
			APITag:                 cfg.APITag,
			AllNestedStructDefs:    renderNestedStructs(props, map[string]bool{}),
			Properties:             props,
			ParentParams:           parentParams,
			HasParent:              len(parentParams) > 0,
			HasCreate:              false,
			HasUpdate:              false,
			HasDelete:              false,
			ReadMethod:             readMethod,
			ListMethod:             listMethod,
			DiscriminatorField:     v.DiscriminatorField,
			DiscriminatorTFAttr:    formatter.TerraformAttrName(v.DiscriminatorField),
			DiscriminatorValue:     v.DiscriminatorValue,
			DiscriminatorSchemaRef: v.SchemaRef,
			IsVariant:              true,
		})
	}
	return out
}

func (g *Generator) buildResourceData(name string, cfg config.ResourceConfig) TemplateData {
	title := formatter.CamelCase(name)
	apiTag := cfg.APITag

	if g.log != nil {
		g.log.Printf("  [buildResourceData] name=%s  title=%s  apiTag=%s", name, title, apiTag)
	}

	var props []PropData
	// Schema priority: GET response → POST response → PUT response → POST request body → PUT request body
	type schemaCandidate struct {
		label  string
		opCfg  *config.OperationConfig
		source string // "response" or "requestbody"
	}
	candidates := []schemaCandidate{
		{"read (GET response)", cfg.Read, "response"},
		{"create (POST response)", cfg.Create, "response"},
		{"update (PUT response)", cfg.Update, "response"},
		{"create (POST request body)", cfg.Create, "requestbody"},
		{"update (PUT request body)", cfg.Update, "requestbody"},
	}
	for _, c := range candidates {
		if c.opCfg == nil {
			continue
		}
		if g.log != nil {
			g.log.Printf("  [buildResourceData] trying %s: method=%s path=%s", c.label, c.opCfg.Method, c.opCfg.Path)
		}
		op := g.spec.GetOperation(c.opCfg.Method, c.opCfg.Path)
		if op == nil {
			if g.log != nil {
				g.log.Printf("  [buildResourceData] WARNING: operation not found in spec for %s %s", c.opCfg.Method, c.opCfg.Path)
			}
			continue
		}
		var schema *openapi.Schema
		if c.source == "response" {
			schema = g.spec.GetResponseSchema(op)
		} else {
			schema = g.spec.GetRequestBodySchema(op)
		}
		if schema == nil {
			if g.log != nil {
				g.log.Printf("  [buildResourceData] no schema in %s, trying next candidate", c.label)
			}
			continue
		}
		if g.log != nil {
			g.log.Printf("  [buildResourceData] using schema from %s: ref=%q type=%q", c.label, schema.Ref, schema.Type)
		}
		props = g.schemaToProps(schema, false, title+"Model")
		break
	}
	if props == nil && g.log != nil {
		g.log.Printf("  [buildResourceData] WARNING: no schema found from any candidate — resource will have no properties")
	}

	// Merge request-body fields into props.
	// Two cases are handled:
	// 1. Request-body-only fields (not in response) are added as WriteOnly so they appear in
	//    the TF schema/model but are never read back from the API response.
	// 2. Writable nested fields that exist in both request body and response schemas but with
	//    different nesting structures: the request-body version is preferred because it directly
	//    mirrors the API body shape that users configure, avoiding spurious extra wrapper levels
	//    introduced by the response schema (e.g. GroupsResponseSchema wraps
	//    IdentitySourceGroupProfileForUpsert in an extra "profile" object that does not exist
	//    in GroupsRequestSchema, which would force users to write profile { profile { ... } }
	//    instead of the natural profile { display_name = ... }).
	{
		existingFields := make(map[string]bool, len(props))
		for _, p := range props {
			existingFields[p.TFAttr] = true
		}
		for _, opCfg := range []*config.OperationConfig{cfg.Create, cfg.Update} {
			if opCfg == nil {
				continue
			}
			op := g.spec.GetOperation(opCfg.Method, opCfg.Path)
			if op == nil {
				continue
			}
			rbSchema := g.spec.GetRequestBodySchema(op)
			if rbSchema == nil {
				continue
			}
			for _, rbProp := range g.schemaToProps(rbSchema, false, title+"Model") {
				if !existingFields[rbProp.TFAttr] {
					// Field only in request body — add as WriteOnly
					if g.log != nil {
						g.log.Printf("  [buildResourceData] merging request-body-only field %q (%s) into props (WriteOnly)", rbProp.TFAttr, rbProp.GoType)
					}
					rbProp.WriteOnly = true
					props = append(props, rbProp)
					existingFields[rbProp.TFAttr] = true
				} else if !rbProp.Computed {
					// Field exists in both schemas and request body version is writable.
					for i, existing := range props {
						if existing.TFAttr != rbProp.TFAttr {
							continue
						}
						if len(rbProp.NestedProps) > 0 {
							// Request body has a writable nested structure — replace the response-schema
							// version so the TF schema matches the API body shape (avoids extra wrapper
							// nesting levels introduced by the response schema).
							if g.log != nil {
								g.log.Printf("  [buildResourceData] replacing response-schema field %q with request-body version (shallower nesting)", rbProp.TFAttr)
							}
							props[i] = rbProp
						} else if existing.Computed {
							// Scalar field: response schema marks it readOnly/Computed, but the request
							// body accepts it as user input — flip to Required or Optional so users can
							// set it. (e.g. UserResponseSchema.externalId has readOnly:true, but
							// UserRequestSchema.externalId is a plain writable field.)
							if g.log != nil {
								g.log.Printf("  [buildResourceData] un-computing response-schema field %q — it is writable in the request body (required=%v)", rbProp.TFAttr, rbProp.Required)
							}
							props[i].Computed = false
							props[i].Required = rbProp.Required
						}
						break
					}
				}
			}
			break // only process Create body once; Update body is handled in UpdateRequestBodyFields
		}
	}

	parentParams := buildParentParams(cfg.ParentParams, g.log)

	// Remove any property that is already emitted as a parent param to avoid duplicate schema attributes.
	if len(parentParams) > 0 {
		parentAttrs := make(map[string]bool, len(parentParams))
		for _, pp := range parentParams {
			parentAttrs[pp.TFAttr] = true
		}
		props = filterByTFAttr(props, parentAttrs)
	}

	// Prefer operationId from the spec for accurate SDK method names.
	// Fall back to convention-derived names if operationId is absent.
	readOpID, createOpID, updateOpID, deleteOpID := "", "", "", ""
	if cfg.Read != nil {
		if op := g.spec.GetOperation(cfg.Read.Method, cfg.Read.Path); op != nil {
			readOpID = op.OperationID
		}
	}
	if cfg.Create != nil {
		if op := g.spec.GetOperation(cfg.Create.Method, cfg.Create.Path); op != nil {
			createOpID = op.OperationID
		}
	}
	if cfg.Update != nil {
		if op := g.spec.GetOperation(cfg.Update.Method, cfg.Update.Path); op != nil {
			updateOpID = op.OperationID
		}
	}
	if cfg.Delete != nil {
		if op := g.spec.GetOperation(cfg.Delete.Method, cfg.Delete.Path); op != nil {
			deleteOpID = op.OperationID
		}
	}

	readMethod := formatter.OperationIDToMethodName(readOpID, "get", name)
	createMethod := formatter.OperationIDToMethodName(createOpID, "create", name)
	updateMethod := formatter.OperationIDToMethodName(updateOpID, "update", name)
	deleteMethod := formatter.OperationIDToMethodName(deleteOpID, "delete", name)
	listMethod := formatter.ListAPIMethodName(name)

	// deriveBodySetter returns the SDK builder method name for a given operation's request body.
	// Priority: x-codegen-request-body-name → request body $ref name → "Body".
	deriveBodySetter := func(opCfg *config.OperationConfig) string {
		if opCfg == nil {
			return ""
		}
		op := g.spec.GetOperation(opCfg.Method, opCfg.Path)
		if op == nil {
			return ""
		}
		if op.RequestBodyName != "" {
			return formatter.GoFieldName(op.RequestBodyName)
		}
		if ref := g.spec.GetRequestBodySchemaRef(op); ref != "" {
			return ref // already CamelCase (e.g. "FederatedClaimRequestBody")
		}
		return ""
	}

	// Derive Create body setter (fall back to Update if no Create op)
	bodySetterMethod := "Body"
	if s := deriveBodySetter(cfg.Create); s != "" {
		bodySetterMethod = s
	} else if s := deriveBodySetter(cfg.Update); s != "" {
		bodySetterMethod = s
	}

	// Derive Update body setter independently (may differ from Create)
	updateBodySetterMethod := bodySetterMethod // default: same as create
	if s := deriveBodySetter(cfg.Update); s != "" {
		updateBodySetterMethod = s
	}

	// Derive SDKTypeName (Create) and UpdateSDKTypeName (Update) separately,
	// since POST and PUT may use different request body schemas
	// (e.g. TrustedOriginWrite for POST, TrustedOrigin for PUT).
	sdkTypeName := title
	if cfg.Create != nil {
		if op := g.spec.GetOperation(cfg.Create.Method, cfg.Create.Path); op != nil {
			if ref := g.spec.GetRequestBodySchemaRef(op); ref != "" {
				sdkTypeName = ref
			}
		}
	} else if cfg.Update != nil {
		if op := g.spec.GetOperation(cfg.Update.Method, cfg.Update.Path); op != nil {
			if ref := g.spec.GetRequestBodySchemaRef(op); ref != "" {
				sdkTypeName = ref
			}
		}
	}
	updateSDKTypeName := sdkTypeName // default: same as Create
	if cfg.Update != nil {
		if op := g.spec.GetOperation(cfg.Update.Method, cfg.Update.Path); op != nil {
			if ref := g.spec.GetRequestBodySchemaRef(op); ref != "" {
				updateSDKTypeName = ref
			}
		}
	}

	// RequestBodyFields: fields that exist in the Create request body schema.
	// We derive these from the Create op's request body schema directly — NOT from the response
	// schema — so that response-only readOnly fields are never included in the body.
	// We restrict to scalar Go types and exclude computed (readOnly) fields.
	var requestBodyFields []PropData
	{
		var rbSchema *openapi.Schema
		rbParentModelName := title + "Model"
		if cfg.Create != nil {
			if op := g.spec.GetOperation(cfg.Create.Method, cfg.Create.Path); op != nil {
				rbSchema = g.spec.GetRequestBodySchema(op)
				if ref := g.spec.GetRequestBodySchemaRef(op); ref != "" {
					rbParentModelName = ref + "Model"
				}
			}
		}
		if rbSchema == nil && cfg.Update != nil {
			// No Create op — fall back to Update schema
			if op := g.spec.GetOperation(cfg.Update.Method, cfg.Update.Path); op != nil {
				rbSchema = g.spec.GetRequestBodySchema(op)
				if ref := g.spec.GetRequestBodySchemaRef(op); ref != "" {
					rbParentModelName = ref + "Model"
				}
			}
		}
		if rbSchema != nil {
			rbProps := g.schemaToProps(rbSchema, false, rbParentModelName)
			// Build a lookup of response-schema props by TFAttr for enriching nested fields.
			respPropsByAttr := make(map[string]PropData, len(props))
			for _, rp := range props {
				respPropsByAttr[rp.TFAttr] = rp
			}
			for _, p := range rbProps {
				if p.Computed {
					continue
				}
				if isScalarGoType(p.GoType) {
					requestBodyFields = append(requestBodyFields, p)
				} else if p.IsListScalar {
					// Scalar list ([]string, []int64, etc.) — include as-is.
					requestBodyFields = append(requestBodyFields, p)
				} else if p.IsListNested && len(p.NestedProps) > 0 && p.SDKNestedTypeName != "" {
					// List-of-nested-objects: enrich with response schema props if available.
					if respProp, ok := respPropsByAttr[p.TFAttr]; ok && len(respProp.NestedProps) > 0 {
						p.TFPlanNestedProps = respProp.NestedProps
					} else {
						p.TFPlanNestedProps = p.NestedProps
					}
					requestBodyFields = append(requestBodyFields, p)
				} else if len(p.NestedProps) > 0 && p.SDKNestedTypeName != "" {
					// Enrich with TFPlanNestedProps from the response schema so the template
					// knows the correct plan field path (response schema may have deeper nesting).
					if respProp, ok := respPropsByAttr[p.TFAttr]; ok && len(respProp.NestedProps) > 0 {
						p.TFPlanNestedProps = respProp.NestedProps
					} else {
						p.TFPlanNestedProps = p.NestedProps
					}
					requestBodyFields = append(requestBodyFields, p)
				}
			}
		}
	}

	// UpdateRequestBodyFields: fields for the Update request body schema.
	// Derived separately because the Update op may use a different (and narrower) schema
	// than the Create op (e.g. AppServiceAccountForUpdate vs AppServiceAccount for Create).
	// Falls back to requestBodyFields when Update has no request body schema.
	var updateRequestBodyFields []PropData
	{
		var updateRbSchema *openapi.Schema
		updateRbParentModelName := title + "Model"
		if cfg.Update != nil {
			if op := g.spec.GetOperation(cfg.Update.Method, cfg.Update.Path); op != nil {
				updateRbSchema = g.spec.GetRequestBodySchema(op)
				if ref := g.spec.GetRequestBodySchemaRef(op); ref != "" {
					updateRbParentModelName = ref + "Model"
				}
			}
		}
		if updateRbSchema != nil {
			rbProps := g.schemaToProps(updateRbSchema, false, updateRbParentModelName)
			// Build a lookup of response-schema props by TFAttr for enriching nested fields.
			respPropsByAttr := make(map[string]PropData, len(props))
			for _, rp := range props {
				respPropsByAttr[rp.TFAttr] = rp
			}
			for _, p := range rbProps {
				if p.Computed {
					continue
				}
				if isScalarGoType(p.GoType) {
					updateRequestBodyFields = append(updateRequestBodyFields, p)
				} else if p.IsListScalar {
					// Scalar list ([]string, []int64, etc.) — include as-is.
					updateRequestBodyFields = append(updateRequestBodyFields, p)
				} else if p.IsListNested && len(p.NestedProps) > 0 && p.SDKNestedTypeName != "" {
					// List-of-nested-objects: enrich with response schema props if available.
					if respProp, ok := respPropsByAttr[p.TFAttr]; ok && len(respProp.NestedProps) > 0 {
						p.TFPlanNestedProps = respProp.NestedProps
					} else {
						p.TFPlanNestedProps = p.NestedProps
					}
					updateRequestBodyFields = append(updateRequestBodyFields, p)
				} else if len(p.NestedProps) > 0 && p.SDKNestedTypeName != "" {
					if respProp, ok := respPropsByAttr[p.TFAttr]; ok && len(respProp.NestedProps) > 0 {
						p.TFPlanNestedProps = respProp.NestedProps
					} else {
						p.TFPlanNestedProps = p.NestedProps
					}
					updateRequestBodyFields = append(updateRequestBodyFields, p)
				}
			}
		}
		if updateRequestBodyFields == nil {
			// No distinct update schema — reuse create fields
			updateRequestBodyFields = requestBodyFields
		}
	}

	// IDField: name of the Go SDK struct setter/getter for the ID field.
	// Most Okta resources use "Id"; we default to that.
	idField := "Id"

	// Auto-detect CreateIDField when:
	//   1. No create_id_field was set in config, AND
	//   2. The Create operation exists, AND
	//   3. The Create operation returns no response schema (i.e. the POST response has no body).
	// In that case the resource ID cannot be read from the response; we pick it from the request
	// body instead.  Selection priority:
	//   a) a field whose snake_case name is "external_id"
	//   b) a field whose snake_case name ends in "_id"
	//   c) the first non-nested, non-computed string field in the request body
	// If nothing matches, we leave CreateIDField empty and let the template fall back to
	// result.GetId() — which will likely cause a compile error that the developer must resolve
	// manually (by adding create_id_field to the config).
	resolvedCreateIDField := cfg.CreateIDField
	if resolvedCreateIDField == "" && cfg.Create != nil {
		createOp := g.spec.GetOperation(cfg.Create.Method, cfg.Create.Path)
		if createOp != nil && g.spec.GetResponseSchema(createOp) == nil {
			// POST returns no body — try to find a suitable ID field in the request body.
			var candidate string
			if rbSchema := g.spec.GetRequestBodySchema(createOp); rbSchema != nil {
				rbProps := g.schemaToProps(rbSchema, false, title+"Model")
				// Pass 1: prefer a field literally named "external_id"
				for _, p := range rbProps {
					if p.TFAttr == "external_id" && len(p.NestedProps) == 0 {
						candidate = p.TFAttr
						break
					}
				}
				// Pass 2: any field whose name ends in "_id"
				if candidate == "" {
					for _, p := range rbProps {
						if strings.HasSuffix(p.TFAttr, "_id") && len(p.NestedProps) == 0 {
							candidate = p.TFAttr
							break
						}
					}
				}
				// Pass 3: first non-nested, non-computed string field
				if candidate == "" {
					for _, p := range rbProps {
						if p.GoType == "types.String" && !p.Computed && len(p.NestedProps) == 0 {
							candidate = p.TFAttr
							break
						}
					}
				}
			}
			if candidate != "" {
				resolvedCreateIDField = candidate
				if g.log != nil {
					g.log.Printf("  [buildResourceData] auto-detected create_id_field=%q (POST returns no body)", candidate)
				}
			} else if g.log != nil {
				g.log.Printf("  [buildResourceData] WARNING: POST returns no body and no suitable create_id_field found; generated Create will likely not compile")
			}
		}
	}

	hasComputedProperties := false
	for _, p := range props {
		if p.Computed {
			hasComputedProperties = true
			break
		}
	}

	if g.log != nil {
		g.log.Printf("  [buildResourceData] operationIds: read=%q create=%q update=%q delete=%q",
			readOpID, createOpID, updateOpID, deleteOpID)
		g.log.Printf("  [buildResourceData] resolved methods: read=%s create=%s update=%s delete=%s list=%s",
			readMethod, createMethod, updateMethod, deleteMethod, listMethod)
		g.log.Printf("  [buildResourceData] hasUpdate=%v hasDelete=%v props=%d parentParams=%d requestBodyFields=%d updateRequestBodyFields=%d",
			cfg.Update != nil, cfg.Delete != nil, len(props), len(parentParams), len(requestBodyFields), len(updateRequestBodyFields))
	}

	return TemplateData{
		Name:                    formatter.LowerFirst(title),
		TitleName:               title,
		TFName:                  name,
		APITag:                  apiTag,
		AllNestedStructDefs:     renderNestedStructs(props, map[string]bool{}),
		Properties:              props,
		RequestBodyFields:       requestBodyFields,
		UpdateRequestBodyFields: updateRequestBodyFields,
		ParentParams:            parentParams,
		HasParent:               len(parentParams) > 0,
		HasCreate:               cfg.Create != nil,
		HasUpdate:               cfg.Update != nil,
		HasDelete:               cfg.Delete != nil,
		ReadMethod:              readMethod,
		CreateMethod:            createMethod,
		UpdateMethod:            updateMethod,
		DeleteMethod:            deleteMethod,
		ListMethod:              listMethod,
		IDField:                 idField,
		SDKTypeName:             sdkTypeName,
		UpdateSDKTypeName:       updateSDKTypeName,
		BodySetterMethod:        bodySetterMethod,
		UpdateBodySetterMethod:  updateBodySetterMethod,
		CreateIDField:           formatter.GoFieldName(resolvedCreateIDField),
		ReadListIDField:         formatter.GoFieldName(cfg.ReadListIDField),
		HasComputedProperties:   hasComputedProperties,
		ReadNoop:                cfg.ReadNoop,
		DeleteNoop:              cfg.DeleteNoop,
	}
}

func (g *Generator) buildDataSourceData(name string, cfg config.DataSourceConfig) TemplateData {
	title := formatter.CamelCase(name)
	apiTag := cfg.APITag

	if g.log != nil {
		g.log.Printf("  [buildDataSourceData] name=%s  title=%s  apiTag=%s", name, title, apiTag)
	}

	var props []PropData
	if cfg.Singular != nil {
		if g.log != nil {
			g.log.Printf("  [buildDataSourceData] looking up singular op: method=%s path=%s", cfg.Singular.Method, cfg.Singular.Path)
		}
		op := g.spec.GetOperation(cfg.Singular.Method, cfg.Singular.Path)
		if op == nil {
			if g.log != nil {
				g.log.Printf("  [buildDataSourceData] WARNING: singular operation not found in spec for %s %s", cfg.Singular.Method, cfg.Singular.Path)
			}
		} else {
			schema := g.spec.GetResponseSchema(op)
			if schema == nil {
				if g.log != nil {
					g.log.Printf("  [buildDataSourceData] WARNING: no response schema for %s %s", cfg.Singular.Method, cfg.Singular.Path)
				}
			} else {
				if g.log != nil {
					g.log.Printf("  [buildDataSourceData] response schema ref=%q type=%q", schema.Ref, schema.Type)
				}
				props = g.schemaToProps(schema, true, title+"DataSourceModel")
			}
		}
	} else if g.log != nil {
		g.log.Printf("  [buildDataSourceData] no singular operation configured")
	}

	// If no singular op, derive properties from the plural (list) op by unwrapping the array items schema.
	var listOpID string
	if cfg.Plural != nil {
		if g.log != nil {
			g.log.Printf("  [buildDataSourceData] plural op: method=%s path=%s", cfg.Plural.Method, cfg.Plural.Path)
		}
		if pluralOp := g.spec.GetOperation(cfg.Plural.Method, cfg.Plural.Path); pluralOp != nil {
			listOpID = pluralOp.OperationID
			if props == nil {
				// list-only datasource: unwrap the array to get the item schema for properties
				itemSchema, isArray := g.spec.GetResponseSchemaUnwrapArray(pluralOp)
				if g.log != nil {
					g.log.Printf("  [buildDataSourceData] plural response isArray=%v schema=%v", isArray, itemSchema != nil)
				}
				if itemSchema != nil {
					props = g.schemaToProps(itemSchema, true, title+"DataSourceModel")
				}
			}
		}
	}

	parentParams := buildParentParams(cfg.ParentParams, g.log)

	// Remove any property that is already emitted as a parent param to avoid duplicate schema attributes.
	if len(parentParams) > 0 {
		parentAttrs := make(map[string]bool, len(parentParams))
		for _, pp := range parentParams {
			parentAttrs[pp.TFAttr] = true
		}
		props = filterByTFAttr(props, parentAttrs)
	}

	// Prefer operationId from the spec for accurate SDK method names.
	readOpID := ""
	if cfg.Singular != nil {
		if op := g.spec.GetOperation(cfg.Singular.Method, cfg.Singular.Path); op != nil {
			readOpID = op.OperationID
		}
	}

	readMethod := formatter.OperationIDToMethodName(readOpID, "get", name)
	// Use the actual operationId from the list op if available, else fall back to convention
	listMethod := formatter.OperationIDToMethodName(listOpID, "list", name)

	if g.log != nil {
		g.log.Printf("  [buildDataSourceData] operationIds: read=%q list=%q", readOpID, listOpID)
		g.log.Printf("  [buildDataSourceData] resolved methods: read=%s list=%s", readMethod, listMethod)
		g.log.Printf("  [buildDataSourceData] props=%d parentParams=%d", len(props), len(parentParams))
	}

	return TemplateData{
		Name:                formatter.LowerFirst(title),
		TitleName:           title,
		TFName:              name,
		APITag:              apiTag,
		AllNestedStructDefs: renderNestedStructs(props, map[string]bool{}),
		Properties:          props,
		ParentParams:        parentParams,
		HasParent:           len(parentParams) > 0,
		HasCreate:           false,
		HasUpdate:           false,
		HasDelete:           false,
		ReadMethod:          readMethod,
		CreateMethod:        "",
		UpdateMethod:        "",
		DeleteMethod:        "",
		ListMethod:          listMethod,
	}
}

// isScalarGoType returns true for TF framework types that map directly to a Go primitive.
// These are the only types the template can safely serialize into an SDK setter call.
// Nested objects (SingleNestedAttribute), lists, and unknown types must be skipped.
func isScalarGoType(goType string) bool {
	switch goType {
	case "types.String", "types.Int64", "types.Bool":
		return true
	}
	return false
}

// filterOutProp returns a new slice with any property whose TFAttr matches fieldName removed.
// Used to prevent the discriminator field from appearing twice in variant resources
// (once from the explicit discriminator block in the template, once from the properties loop).
func filterOutProp(props []PropData, fieldName string) []PropData {
	out := props[:0:0] // same backing-array type, zero len
	for _, p := range props {
		if p.TFAttr != fieldName {
			out = append(out, p)
		}
	}
	return out
}

// filterByTFAttr removes any prop whose TFAttr key exists in the excluded map.
// Used to prevent parent-param fields from being emitted a second time in the properties loop.
func filterByTFAttr(props []PropData, excluded map[string]bool) []PropData {
	out := props[:0:0]
	for _, p := range props {
		if !excluded[p.TFAttr] {
			out = append(out, p)
		}
	}
	return out
}

// excludedTFAttrs lists TF attribute names that should never be surfaced in the
// generated schema/model. These are internal Okta API fields that are meaningless
// or harmful in a Terraform provider:
//   - "links" / "_links": HAL hypermedia links — purely navigational, not config
var excludedTFAttrs = map[string]bool{
	"links":  true,
	"_links": true,
}

// filterExcludedProps removes globally-excluded fields from a props slice.
func filterExcludedProps(props []PropData) []PropData {
	out := props[:0:0]
	for _, p := range props {
		if !excludedTFAttrs[p.TFAttr] {
			out = append(out, p)
		}
	}
	return out
}

func buildParentParams(params []config.ParentParam, logger *log.Logger) []ParentParamData {
	if logger != nil {
		logger.Printf("  [buildParentParams] %d param(s) to process", len(params))
	}
	var out []ParentParamData
	for _, p := range params {
		desc := p.Description
		if desc == "" {
			desc = "The ID of the parent " + p.Name
		}
		goField := formatter.GoFieldName(p.Name)
		if logger != nil {
			logger.Printf("  [buildParentParams]   name=%s → GoField=%s  TFAttr=%s  PathParam=%s",
				p.Name, goField, p.Name, p.PathParam)
		}

		/*ParentParamData{
		    GoField:     "AuthServerID",      // used in model struct:  AuthServerID types.String
		    TFAttr:      "auth_server_id",    // used in schema:        "auth_server_id": schema.StringAttribute{...}
		    Description: "The ID of the authorization server",
		    PathParam:   "{authServerId}",    // used in API call comment to show which path var it fills
		}*/
		out = append(out, ParentParamData{
			GoField:     goField,
			TFAttr:      p.Name,
			Description: desc,
			PathParam:   p.PathParam,
		})
	}
	return out
}

func (g *Generator) schemaToProps(schema *openapi.Schema, allComputed bool, parentModelName string) []PropData {
	return g.schemaToPropsDepth(schema, allComputed, 0, parentModelName)
}

// schemaToPropsDepth is the recursive inner implementation.
// depth guards against truly circular inline schemas (limit: 10).
// In practice the Okta management spec has no circular inline objects; max real depth is 5.
func (g *Generator) schemaToPropsDepth(schema *openapi.Schema, allComputed bool, depth int, parentModelName string) []PropData {
	rawProps := g.spec.GetProperties(*schema)

	if g.log != nil {
		g.log.Printf("  [schemaToProps] raw properties from spec: %d  allComputed=%v", len(rawProps), allComputed)
	}

	// Sort for deterministic output
	sort.Slice(rawProps, func(i, j int) bool {
		return rawProps[i].Name < rawProps[j].Name
	})

	// Deduplicate by tfAttr — allOf composition can yield the same property from multiple levels.
	seenTFAttrs := make(map[string]bool)

	var props []PropData
	for _, p := range rawProps {
		name := p.Name
		tfAttr := formatter.TerraformAttrName(name)
		goField := formatter.GoFieldName(name)

		// Skip if would collide with ID
		if tfAttr == "id" {
			if g.log != nil {
				g.log.Printf("  [schemaToProps]   SKIP %q (collides with id)", name)
			}
			continue
		}

		// Skip globally-excluded fields (e.g. HAL _links / links)
		if excludedTFAttrs[tfAttr] {
			if g.log != nil {
				g.log.Printf("  [schemaToProps]   SKIP %q (globally excluded)", name)
			}
			continue
		}

		// Skip duplicate tfAttrs (can arise from allOf merging or camelCase/snake_case variants)
		if seenTFAttrs[tfAttr] {
			if g.log != nil {
				g.log.Printf("  [schemaToProps]   SKIP %q (duplicate tfAttr %q)", name, tfAttr)
			}
			continue
		}
		seenTFAttrs[tfAttr] = true

		desc := formatter.SanitizeDescription(p.Description)
		if desc == "" {
			desc = goField
		}

		goType := oktypes.GoType(p.Schema)
		tfSchemaType := oktypes.TFSchemaType(p.Schema)
		elementType := ""
		sdkNestedTypeName := ""
		var nestedProps []PropData

		// Only flatten nested objects into typed structs when the field is writable.
		// Read-only (computed) nested objects are represented as types.Object — users never
		// configure them, so there's no benefit in expanding the full struct hierarchy.
		isWritable := !p.ReadOnly && !allComputed

		switch {
		case p.Schema.Type == "array":
			// Resolve items schema to determine whether items are scalars or objects.
			if p.Schema.Items != nil && isWritable && depth < 10 {
				resolved := g.spec.ResolveSchema(*p.Schema.Items)
				// Object items (inline or via $ref) → ListNestedAttribute
				if resolved.Type == "object" || resolved.Ref != "" || len(resolved.Properties) > 0 {
					nestedModelName := parentModelName + formatter.GoFieldName(name) + "Model"
					candidateNestedProps := g.schemaToPropsDepth(&resolved, allComputed, depth+1, nestedModelName)
					if len(candidateNestedProps) > 0 {
						tfSchemaType = "schema.ListNestedAttribute"
						goType = "[]" + nestedModelName
						nestedProps = candidateNestedProps
						// Capture SDK item type name from $ref if available.
						// Prefer items $ref (direct item type); fall back to OriginalRef.
						// When OriginalRef is a named array schema (e.g. IdentitySourceGroupMembershipsDeleteProfile
						// which is type:array), the SDK codegen appends "Inner" to name the item struct.
						// When both are absent (inline items), the SDK codegen convention is
						// <ParentSchemaName><FieldNamePascal>Inner.
						const refPrefix = "#/components/schemas/"
						if p.Schema.Items.Ref != "" {
							sdkNestedTypeName = p.Schema.Items.Ref
						} else if p.OriginalRef != "" {
							candidateName := p.OriginalRef
							if strings.HasPrefix(candidateName, refPrefix) {
								candidateName = candidateName[len(refPrefix):]
							}
							if named, ok := g.spec.Components.Schemas[candidateName]; ok && named.Type == "array" {
								// Named array schema → SDK item struct is <Name>Inner
								sdkNestedTypeName = candidateName + "Inner"
							} else {
								sdkNestedTypeName = candidateName
							}
						} else {
							// Inline items with no $ref — derive from parent schema name + field + "Inner".
							// parentModelName is e.g. "BulkUpsertRequestBodyModel"; strip the "Model" suffix
							// to get "BulkUpsertRequestBody", then append PascalCase(name) + "Inner".
							parentBase := strings.TrimSuffix(parentModelName, "Model")
							if parentBase != "" {
								sdkNestedTypeName = parentBase + formatter.GoFieldName(name) + "Inner"
							}
						}
						if sdkNestedTypeName != "" && strings.HasPrefix(sdkNestedTypeName, refPrefix) {
							sdkNestedTypeName = sdkNestedTypeName[len(refPrefix):]
						}
					} else {
						// Object items but no properties resolved — fall back to scalar list
						elementType = oktypes.ElementTypeStr(resolved)
					}
				} else {
					// Scalar items (string/bool/int)
					elementType = oktypes.ElementTypeStr(resolved)
				}
			} else if p.Schema.Items != nil {
				resolved := g.spec.ResolveSchema(*p.Schema.Items)
				elementType = oktypes.ElementTypeStr(resolved)
			} else {
				elementType = "types.StringType" // safe fallback
			}
		case isWritable && p.Schema.Type == "object" && len(p.Schema.Properties) > 0 && p.Schema.Ref == "" && depth < 10:
			// Inline writable object with known sub-fields → recurse to build SingleNestedAttribute.
			nestedModelName := parentModelName + formatter.GoFieldName(name) + "Model"
			candidateNestedProps := g.schemaToPropsDepth(&p.Schema, allComputed, depth+1, nestedModelName)
			if len(candidateNestedProps) > 0 {
				tfSchemaType = "schema.SingleNestedAttribute"
				goType = "*" + nestedModelName
				nestedProps = candidateNestedProps
				// If this inline object originated from a $ref, capture the SDK type name
				// so the template can call okta.New<SDKNestedTypeName>WithDefaults().
				if p.OriginalRef != "" {
					sdkNestedTypeName = p.OriginalRef
					const refPrefix = "#/components/schemas/"
					if strings.HasPrefix(sdkNestedTypeName, refPrefix) {
						sdkNestedTypeName = sdkNestedTypeName[len(refPrefix):]
					}
				}
			}

		case isWritable && p.Schema.Ref != "" && depth < 10:
			// $ref to a named schema — resolve and expand as SingleNestedAttribute if it has properties.
			resolvedRef := g.spec.ResolveSchema(p.Schema)
			refProps := g.spec.GetProperties(resolvedRef)
			if len(refProps) > 0 {
				nestedModelName := parentModelName + formatter.GoFieldName(name) + "Model"
				candidateNestedProps := g.schemaToPropsDepth(&resolvedRef, allComputed, depth+1, nestedModelName)
				if len(candidateNestedProps) > 0 {
					tfSchemaType = "schema.SingleNestedAttribute"
					goType = "*" + nestedModelName
					nestedProps = candidateNestedProps
					// Capture the SDK type name from the $ref so the template can call okta.New<SDKNestedTypeName>WithDefaults()
					// Strip the "#/components/schemas/" prefix to get just the type name.
					sdkNestedTypeName = p.Schema.Ref
					const refPrefix = "#/components/schemas/"
					if strings.HasPrefix(sdkNestedTypeName, refPrefix) {
						sdkNestedTypeName = sdkNestedTypeName[len(refPrefix):]
					}
				}
			}

		case isWritable && p.Schema.Type == "" && p.Schema.Ref == "" && (len(p.Schema.AllOf) > 0 || len(p.Schema.OneOf) > 0) && depth < 10:
			// allOf/oneOf composition — treat as object, collect all properties.
			composedProps := g.spec.GetProperties(p.Schema)
			if len(composedProps) > 0 {
				nestedModelName := parentModelName + formatter.GoFieldName(name) + "Model"
				candidateNestedProps := g.schemaToPropsDepth(&p.Schema, allComputed, depth+1, nestedModelName)
				if len(candidateNestedProps) > 0 {
					tfSchemaType = "schema.SingleNestedAttribute"
					goType = "*" + nestedModelName
					nestedProps = candidateNestedProps
					// If this composed object originated from a $ref, capture the SDK type name
					// so the template can call okta.New<SDKNestedTypeName>WithDefaults().
					if p.OriginalRef != "" {
						sdkNestedTypeName = p.OriginalRef
						const refPrefix = "#/components/schemas/"
						if strings.HasPrefix(sdkNestedTypeName, refPrefix) {
							sdkNestedTypeName = sdkNestedTypeName[len(refPrefix):]
						}
					}
				}
			}
		}

		required := p.Required && !allComputed
		computed := allComputed || p.ReadOnly

		if g.log != nil {
			g.log.Printf("  [schemaToProps]   ACCEPT %-25s → GoField=%-20s GoType=%-15s TFSchema=%-25s required=%-5v computed=%v readOnly=%v",
				name, goField, goType, tfSchemaType, required, computed, p.ReadOnly)
		}

		nestedModelName := ""
		schemaAttrBlock := ""
		isListNested := tfSchemaType == "schema.ListNestedAttribute"
		isListScalar := elementType != "" && !isListNested
		if len(nestedProps) > 0 {
			nestedModelName = parentModelName + formatter.GoFieldName(name) + "Model"
			// Pre-render the schema attribute block (schema only, no struct defs here).
			pd := PropData{
				TFAttr:       tfAttr,
				TFSchemaType: tfSchemaType,
				IsListNested: isListNested,
				IsListScalar: isListScalar,
				ElementType:  elementType,
				NestedProps:  nestedProps,
				Description:  desc,
				Required:     required,
				Computed:     computed,
			}
			schemaAttrBlock = renderSchemaAttr(pd, 3)
		}
		props = append(props, PropData{
			GoField:           goField,
			GoType:            goType,
			TFAttr:            tfAttr,
			TFSchemaType:      tfSchemaType,
			IsListNested:      isListNested,
			IsListScalar:      isListScalar,
			ElementType:       elementType,
			NestedProps:       nestedProps,
			NestedModelName:   nestedModelName,
			SDKNestedTypeName: sdkNestedTypeName,
			SchemaAttrBlock:   schemaAttrBlock,
			Description:       desc,
			Required:          required,
			Computed:          computed,
			IsDateTime:        goType == "types.String" && p.Schema.Format == "date-time",
		})
	}

	if g.log != nil {
		g.log.Printf("  [schemaToProps] accepted %d / %d properties", len(props), len(rawProps))
	}
	return props
}

// renderSchemaAttr returns the Go source for a single schema.Attribute entry at the given
// indentation level. It recurses into NestedProps so there is no depth cap in the template.
func renderSchemaAttr(p PropData, indent int) string {
	tabs := strings.Repeat("\t", indent)
	var b strings.Builder

	if len(p.NestedProps) > 0 && p.IsListNested {
		// ListNestedAttribute — wrap nested attrs in a NestedAttributeObject
		b.WriteString("schema.ListNestedAttribute{\n")
		b.WriteString(fmt.Sprintf("%s\tDescription: %q,\n", tabs, p.Description))
		if p.Required {
			b.WriteString(fmt.Sprintf("%s\tRequired: true,\n", tabs))
		} else if !p.Computed {
			b.WriteString(fmt.Sprintf("%s\tOptional: true,\n", tabs))
		}
		if p.Computed {
			b.WriteString(fmt.Sprintf("%s\tComputed: true,\n", tabs))
		}
		b.WriteString(fmt.Sprintf("%s\tNestedObject: schema.NestedAttributeObject{\n", tabs))
		b.WriteString(fmt.Sprintf("%s\t\tAttributes: map[string]schema.Attribute{\n", tabs))
		for _, sub := range p.NestedProps {
			b.WriteString(fmt.Sprintf("%s\t\t\t%q: %s,\n", tabs, sub.TFAttr, renderSchemaAttr(sub, indent+3)))
		}
		b.WriteString(fmt.Sprintf("%s\t\t},\n", tabs))
		b.WriteString(fmt.Sprintf("%s\t},\n", tabs))
		b.WriteString(fmt.Sprintf("%s}", tabs))
	} else if len(p.NestedProps) > 0 {
		// SingleNestedAttribute — emit Attributes map recursively
		b.WriteString("schema.SingleNestedAttribute{\n")
		b.WriteString(fmt.Sprintf("%s\tDescription: %q,\n", tabs, p.Description))
		if p.Required {
			b.WriteString(fmt.Sprintf("%s\tRequired: true,\n", tabs))
		} else if !p.Computed {
			b.WriteString(fmt.Sprintf("%s\tOptional: true,\n", tabs))
		}
		if p.Computed {
			b.WriteString(fmt.Sprintf("%s\tComputed: true,\n", tabs))
		}
		b.WriteString(fmt.Sprintf("%s\tAttributes: map[string]schema.Attribute{\n", tabs))
		for _, sub := range p.NestedProps {
			b.WriteString(fmt.Sprintf("%s\t\t%q: %s,\n", tabs, sub.TFAttr, renderSchemaAttr(sub, indent+2)))
		}
		b.WriteString(fmt.Sprintf("%s\t},\n", tabs))
		b.WriteString(fmt.Sprintf("%s}", tabs))
	} else {
		b.WriteString(fmt.Sprintf("%s{\n", p.TFSchemaType))
		b.WriteString(fmt.Sprintf("%s\tDescription: %q,\n", tabs, p.Description))
		if p.ElementType != "" {
			b.WriteString(fmt.Sprintf("%s\tElementType: %s,\n", tabs, p.ElementType))
		}
		if p.Required {
			b.WriteString(fmt.Sprintf("%s\tRequired: true,\n", tabs))
		} else if !p.Computed {
			b.WriteString(fmt.Sprintf("%s\tOptional: true,\n", tabs))
		}
		if p.Computed {
			b.WriteString(fmt.Sprintf("%s\tComputed: true,\n", tabs))
		}
		b.WriteString(fmt.Sprintf("%s}", tabs))
	}
	return b.String()
}

// renderNestedStructs recursively emits Go struct type definitions for all nested models
// reachable from props. Deduplicates by struct name to avoid redeclarations.
func renderNestedStructs(props []PropData, seen map[string]bool) string {
	var b strings.Builder
	for _, p := range props {
		if p.NestedModelName == "" {
			continue
		}
		if seen[p.NestedModelName] {
			continue
		}
		seen[p.NestedModelName] = true
		b.WriteString(fmt.Sprintf("// %s is the nested model for %s.\n", p.NestedModelName, p.TFAttr))
		b.WriteString(fmt.Sprintf("type %s struct {\n", p.NestedModelName))
		for _, sub := range p.NestedProps {
			b.WriteString(fmt.Sprintf("\t%s %s `tfsdk:%q`\n", sub.GoField, sub.GoType, sub.TFAttr))
		}
		b.WriteString("}\n\n")
		// Recurse into children using the shared seen map to deduplicate.
		b.WriteString(renderNestedStructs(p.NestedProps, seen))
	}
	return b.String()
}

func (g *Generator) renderToFile(tmplName, outPath string, data TemplateData) error {
	if g.log != nil {
		g.log.Printf("  Rendering template=%s → %s", tmplName, filepath.Base(outPath))
	}

	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, tmplName, data); err != nil {
		return fmt.Errorf("executing template %s: %w", tmplName, err)
	}

	content := buf.Bytes()
	if g.log != nil {
		g.log.Printf("  Template rendered: %d bytes", len(content))
	}

	if g.goFmt {
		formatted, err := format.Source(content)
		if err != nil {
			// Write unformatted on error so user can debug
			fmt.Fprintf(os.Stderr, "Warning: gofmt failed for %s: %v\n", outPath, err)
			if g.log != nil {
				g.log.Printf("  gofmt FAILED: %v", err)
			}
		} else {
			content = formatted
			if g.log != nil {
				g.log.Printf("  gofmt OK: %d bytes", len(content))
			}
		}
	}

	if err := os.WriteFile(outPath, content, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", outPath, err)
	}
	if g.log != nil {
		g.log.Printf("  Written: %s", outPath)
	}
	return nil
}
