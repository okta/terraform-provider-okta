package types

import (
	"github.com/okta/terraform-provider-okta/generator/internal/openapi"
)

// TFSchemaType returns the Terraform schema attribute type for an OpenAPI schema
func TFSchemaType(schema openapi.Schema) string {
	switch schema.Type {
	case "string":
		return "schema.StringAttribute"
	case "boolean":
		return "schema.BoolAttribute"
	case "integer", "number":
		return "schema.Int64Attribute"
	case "array":
		return "schema.ListAttribute"
	case "object":
		return "schema.ObjectAttribute"
	default:
		if schema.Ref != "" {
			return "schema.ObjectAttribute"
		}
		return "schema.StringAttribute"
	}
}

// GoType returns the Go types.X type for an OpenAPI schema
func GoType(schema openapi.Schema) string {
	switch schema.Type {
	case "string":
		return "types.String"
	case "boolean":
		return "types.Bool"
	case "integer", "number":
		return "types.Int64"
	case "array":
		return "types.List"
	case "object":
		return "types.Object"
	default:
		if schema.Ref != "" {
			return "types.Object"
		}
		return "types.String"
	}
}

// IsPrimitive returns true if the schema type maps to a primitive Terraform attribute
func IsPrimitive(schema openapi.Schema) bool {
	switch schema.Type {
	case "string", "boolean", "integer", "number":
		return true
	}
	return false
}

// ElementTypeStr returns the attr.Type expression used in schema.ListAttribute{ElementType: ...}
// for the given items schema. Falls back to types.StringType for unknown/object/array items.
func ElementTypeStr(schema openapi.Schema) string {
	switch schema.Type {
	case "string":
		return "types.StringType"
	case "boolean":
		return "types.BoolType"
	case "integer", "number":
		return "types.Int64Type"
	default:
		// object, nested array, or unresolved $ref — use StringType as safe fallback
		return "types.StringType"
	}
}
