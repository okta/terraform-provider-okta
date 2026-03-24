package formatter

import (
	"regexp"
	"strings"
	"unicode"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
var nonAlphaNum = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// SnakeCase converts CamelCase or mixed strings to snake_case
func SnakeCase(s string) string {
	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = nonAlphaNum.ReplaceAllString(snake, "_")
	snake = strings.Trim(snake, "_")
	return strings.ToLower(snake)
}

// CamelCase converts snake_case or kebab-case to CamelCase (PascalCase).
// Splits on both '_' and '-' so that e.g. "grant-creatable" → "GrantCreatable".
func CamelCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool { return r == '_' || r == '-' })
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

// LowerFirst lowercases the first character of a string
func LowerFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// UpperFirst uppercases the first character of a string
func UpperFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// TerraformAttrName converts a property name to a Terraform attribute name (snake_case)
func TerraformAttrName(s string) string {
	// Handle names that start with _ (like _embedded, _links)
	s = strings.TrimPrefix(s, "_")
	return SnakeCase(s)
}

// GoFieldName converts a property name to a Go struct field name (CamelCase)
func GoFieldName(s string) string {
	// Handle names that start with _ (like _embedded, _links)
	s = strings.TrimPrefix(s, "_")
	return CamelCase(SnakeCase(s))
}

// APIMethodName derives a Go SDK method name from operationId or path+method
// e.g. "get" + "group" -> "GetGroup", "create" + "group" -> "CreateGroup"
func APIMethodName(operation, resourceName string) string {
	op := strings.ToLower(operation)
	switch op {
	case "get", "read":
		return "Get" + CamelCase(resourceName)
	case "post", "create":
		return "Create" + CamelCase(resourceName)
	case "put", "patch", "update":
		return "Update" + CamelCase(resourceName)
	case "delete":
		return "Delete" + CamelCase(resourceName)
	case "list":
		return "List" + CamelCase(resourceName) + "s"
	default:
		return CamelCase(operation) + CamelCase(resourceName)
	}
}

// ListAPIMethodName derives the list/plural SDK method name
func ListAPIMethodName(resourceName string) string {
	return "List" + CamelCase(resourceName) + "s"
}

// OperationIDToMethodName converts an OAS operationId (lowerCamelCase) to a Go SDK method name (UpperCamelCase).
// e.g. "getGroup" → "GetGroup", "replaceGroup" → "ReplaceGroup"
// Falls back to APIMethodName convention if operationId is empty.
func OperationIDToMethodName(operationID, fallbackOp, resourceName string) string {
	if operationID != "" {
		return UpperFirst(operationID)
	}
	return APIMethodName(fallbackOp, resourceName)
}

// SanitizeDescription removes newlines and trims descriptions for use in Go comments/struct tags
func SanitizeDescription(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, `"`, `'`)
	// Trim very long descriptions
	if len(s) > 200 {
		s = s[:197] + "..."
	}
	return strings.TrimSpace(s)
}

// ExtractIDParam extracts the last path parameter name from a path like /api/v1/groups/{groupId}
// Returns "groupId" -> "groupId"
func ExtractIDParam(path string) string {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(path, -1)
	if len(matches) == 0 {
		return "id"
	}
	return matches[len(matches)-1][1]
}
