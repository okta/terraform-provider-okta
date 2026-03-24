package openapi

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Spec is a simplified representation of an OpenAPI 3.x spec
type Spec struct {
	Components Components          `yaml:"components"`
	Paths      map[string]PathItem `yaml:"paths"`
}

type Components struct {
	Schemas map[string]Schema `yaml:"schemas"`
}

// PathItem holds HTTP method operations keyed by method
type PathItem struct {
	Get    *Operation `yaml:"get"`
	Post   *Operation `yaml:"post"`
	Put    *Operation `yaml:"put"`
	Patch  *Operation `yaml:"patch"`
	Delete *Operation `yaml:"delete"`
}

// Operation holds a single HTTP operation
type Operation struct {
	OperationID     string              `yaml:"operationId"`
	RequestBodyName string              `yaml:"x-codegen-request-body-name"`
	Tags            []string            `yaml:"tags"`
	Summary         string              `yaml:"summary"`
	Description     string              `yaml:"description"`
	Parameters      []Parameter         `yaml:"parameters"`
	RequestBody     *RequestBody        `yaml:"requestBody"`
	Responses       map[string]Response `yaml:"responses"`
}

type Parameter struct {
	Name        string `yaml:"name"`
	In          string `yaml:"in"`
	Required    bool   `yaml:"required"`
	Description string `yaml:"description"`
	Schema      Schema `yaml:"schema"`
	Ref         string `yaml:"$ref"`
}

type RequestBody struct {
	Content map[string]MediaType `yaml:"content"`
}

type MediaType struct {
	Schema Schema `yaml:"schema"`
}

type Response struct {
	Content map[string]MediaType `yaml:"content"`
	Ref     string               `yaml:"$ref"`
}

// Discriminator holds the OAS3 discriminator object
type Discriminator struct {
	PropertyName string            `yaml:"propertyName"`
	Mapping      map[string]string `yaml:"mapping"`
}

// Schema is a simplified OpenAPI schema
type Schema struct {
	Ref           string            `yaml:"$ref"`
	Type          string            `yaml:"type"`
	Format        string            `yaml:"format"`
	Description   string            `yaml:"description"`
	Properties    map[string]Schema `yaml:"properties"`
	Items         *Schema           `yaml:"items"`
	AllOf         []Schema          `yaml:"allOf"`
	OneOf         []Schema          `yaml:"oneOf"`
	Discriminator *Discriminator    `yaml:"discriminator"`
	ReadOnly      bool              `yaml:"readOnly"`
	Enum          []interface{}     `yaml:"enum"`
	Required      []string          `yaml:"required"`
}

// Property is a resolved schema property with its name
type Property struct {
	Name        string
	Schema      Schema
	Description string
	Required    bool
	ReadOnly    bool
}

// Load reads and parses an OpenAPI YAML spec file
func Load(path string) (*Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading spec file: %w", err)
	}
	var spec Spec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("parsing spec YAML: %w", err)
	}
	return &spec, nil
}

// GetOperation returns the operation for a given HTTP method and path
func (s *Spec) GetOperation(method, path string) *Operation {
	item, ok := s.Paths[path]
	if !ok {
		return nil
	}
	switch strings.ToLower(method) {
	case "get":
		return item.Get
	case "post":
		return item.Post
	case "put":
		return item.Put
	case "patch":
		return item.Patch
	case "delete":
		return item.Delete
	}
	return nil
}

// ResolveSchema follows $ref chains to return the final schema
func (s *Spec) ResolveSchema(schema Schema) Schema {
	if schema.Ref == "" {
		return schema
	}
	name := refName(schema.Ref)
	resolved, ok := s.Components.Schemas[name]
	if !ok {
		return schema
	}
	return s.ResolveSchema(resolved)
}

// GetProperties returns all properties from a schema, resolving $refs and allOf
func (s *Spec) GetProperties(schema Schema) []Property {
	resolved := s.ResolveSchema(schema)

	var props []Property

	// Handle allOf
	for _, sub := range resolved.AllOf {
		subResolved := s.ResolveSchema(sub)
		props = append(props, s.GetProperties(subResolved)...)
	}

	// Handle direct properties
	requiredSet := make(map[string]bool)
	for _, r := range resolved.Required {
		requiredSet[r] = true
	}

	for name, propSchema := range resolved.Properties {
		resolvedProp := s.ResolveSchema(propSchema)
		desc := resolvedProp.Description
		if desc == "" {
			desc = propSchema.Description
		}
		props = append(props, Property{
			Name:        name,
			Schema:      resolvedProp,
			Description: desc,
			Required:    requiredSet[name],
			ReadOnly:    resolvedProp.ReadOnly || propSchema.ReadOnly,
		})
	}
	return props
}

// GetResponseSchema returns the schema from a 200/201 response body
func (s *Spec) GetResponseSchema(op *Operation) *Schema {
	if op == nil {
		return nil
	}
	for _, code := range []string{"200", "201"} {
		resp, ok := op.Responses[code]
		if !ok {
			continue
		}
		// resolve $ref on the response itself
		if resp.Ref != "" {
			// not handling response $refs for now
			continue
		}
		if ct, ok := resp.Content["application/json"]; ok {
			sc := ct.Schema
			if sc.Ref != "" {
				resolved := s.ResolveSchema(sc)
				return &resolved
			}
			return &sc
		}
	}
	return nil
}

// GetResponseSchemaUnwrapArray is like GetResponseSchema but also handles array responses.
// If the response schema is type:array, it returns the items schema (the element type)
// and sets isArray=true. This is used for list-only operations that return []SomeType.
func (s *Spec) GetResponseSchemaUnwrapArray(op *Operation) (schema *Schema, isArray bool) {
	raw := s.GetResponseSchema(op)
	if raw == nil {
		return nil, false
	}
	if raw.Type == "array" && raw.Items != nil {
		if raw.Items.Ref != "" {
			resolved := s.ResolveSchema(*raw.Items)
			return &resolved, true
		}
		return raw.Items, true
	}
	return raw, false
}

// GetRequestBodySchema returns the schema from the request body
func (s *Spec) GetRequestBodySchema(op *Operation) *Schema {
	if op == nil || op.RequestBody == nil {
		return nil
	}
	if ct, ok := op.RequestBody.Content["application/json"]; ok {
		sc := ct.Schema
		if sc.Ref != "" {
			resolved := s.ResolveSchema(sc)
			return &resolved
		}
		return &sc
	}
	return nil
}

// refName extracts the schema name from a $ref string like "#/components/schemas/Group"
func refName(ref string) string {
	parts := strings.Split(ref, "/")
	return parts[len(parts)-1]
}

// GetSchemaByRef looks up a component schema by its short name (e.g. "SamlApplication").
// Returns nil if not found.
func (s *Spec) GetSchemaByRef(name string) *Schema {
	sc, ok := s.Components.Schemas[name]
	if !ok {
		return nil
	}
	resolved := s.ResolveSchema(sc)
	return &resolved
}

// GetDiscriminatorFromSchema returns the discriminator if the schema root is a oneOf+discriminator,
// otherwise returns nil.
func GetDiscriminatorFromSchema(schema Schema) *Discriminator {
	if len(schema.OneOf) > 0 && schema.Discriminator != nil {
		return schema.Discriminator
	}
	return nil
}

// GetRequestBodySchemaRef returns the bare schema name ($ref tail) from the request body,
// e.g. "CAPTCHAInstance" for $ref: '#/components/schemas/CAPTCHAInstance'.
// Returns "" if the request body has no $ref (inline schema).
func (s *Spec) GetRequestBodySchemaRef(op *Operation) string {
	if op == nil || op.RequestBody == nil {
		return ""
	}
	if ct, ok := op.RequestBody.Content["application/json"]; ok {
		if ct.Schema.Ref != "" {
			return refName(ct.Schema.Ref)
		}
	}
	return ""
}

// GetUnionTypeName derives the SDK union wrapper type name for a polymorphic (oneOf) response.
// For a singular read path like /api/v1/logStreams/{logStreamId}, it finds the GET on the
// collection path /api/v1/logStreams, reads its operationId (e.g. "listLogStreams"), and
// returns "List<PascalCase(opId)>200ResponseInner" — the name the SDK generator uses.
// Returns "" if the collection path or its GET operation cannot be found.
func (s *Spec) GetUnionTypeName(readPath string) string {
	// Strip the trailing path segment that contains a path parameter, e.g.
	// /api/v1/logStreams/{logStreamId} → /api/v1/logStreams
	// /api/v1/behaviors/{behaviorId}  → /api/v1/behaviors
	idx := strings.LastIndex(readPath, "/")
	if idx < 0 {
		return ""
	}
	// The last segment might be a path param like "{logStreamId}" or a literal.
	// Only strip if it looks like a path param.
	lastSeg := readPath[idx+1:]
	var collectionPath string
	if strings.HasPrefix(lastSeg, "{") && strings.HasSuffix(lastSeg, "}") {
		collectionPath = readPath[:idx]
	} else {
		// Singular path doesn't end in a param — cannot derive collection path reliably.
		return ""
	}

	item, ok := s.Paths[collectionPath]
	if !ok {
		return ""
	}
	listOp := item.Get
	if listOp == nil || listOp.OperationID == "" {
		return ""
	}
	// listOpId is lowerCamelCase e.g. "listLogStreams".
	// SDK union type: "List<UpperFirst(listOpId)>200ResponseInner"
	// = "List" + upper-first the operationId (already CamelCase except first char).
	opId := listOp.OperationID

	// opId is lowerCamelCase already starting with "list", e.g. "listLogStreams".
	// SDK type name: UpperFirst(opId) + "200ResponseInner" = "ListLogStreams200ResponseInner".
	unionType := strings.ToUpper(opId[:1]) + opId[1:] + "200ResponseInner"
	return unionType
}
