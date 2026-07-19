package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <path-to-openapi-spec.yaml>")
	}

	specPath := os.Args[1]
	outputDir := "./generated"

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Read the OpenAPI spec
	specBytes, err := os.ReadFile(specPath)
	if err != nil {
		log.Fatalf("Failed to read spec file: %v", err)
	}

	// Parse YAML directly
	var spec map[string]interface{}
	if err := yaml.Unmarshal(specBytes, &spec); err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	fmt.Printf("Generating Go models from OpenAPI spec: %s\n", specPath)
	fmt.Printf("Output directory: %s\n\n", outputDir)

	// Extract schemas from components
	components, ok := spec["components"].(map[string]interface{})
	if !ok {
		log.Fatal("No components found in spec")
	}

	schemas, ok := components["schemas"].(map[string]interface{})
	if !ok {
		log.Fatal("No schemas found in components")
	}

	fmt.Printf("Found %d schemas\n", len(schemas))

	// Generate models
	schemaCount := 0
	for schemaName, schemaData := range schemas {
		schemaMap, ok := schemaData.(map[string]interface{})
		if !ok {
			continue
		}

		goCode := generateGoStruct(schemaName, schemaMap)

		fileName := toSnakeCase(schemaName) + ".go"
		filePath := filepath.Join(outputDir, fileName)

		if err := os.WriteFile(filePath, []byte(goCode), 0644); err != nil {
			log.Printf("Warning: Failed to write file %s: %v", fileName, err)
			continue
		}

		schemaCount++
		if schemaCount%100 == 0 {
			fmt.Printf("Generated %d models...\n", schemaCount)
		}
	}

	fmt.Printf("\nGeneration complete! Generated %d Go model files in: %s\n", schemaCount, outputDir)
}

func generateGoStruct(name string, schema map[string]interface{}) string {
	var sb strings.Builder

	// Package and imports
	sb.WriteString("// Code generated from OpenAPI spec. DO NOT EDIT.\n")
	sb.WriteString("package models\n\n")

	// Check if we need time import
	needsTime := false
	if props, ok := schema["properties"].(map[string]interface{}); ok {
		for _, propData := range props {
			if propMap, ok := propData.(map[string]interface{}); ok {
				if format, _ := propMap["format"].(string); format == "date-time" {
					needsTime = true
					break
				}
			}
		}
	}

	if needsTime {
		sb.WriteString("import \"time\"\n\n")
	}

	// Struct definition
	sb.WriteString(fmt.Sprintf("// %s represents the %s schema\n", name, name))
	if desc, ok := schema["description"].(string); ok && desc != "" {
		sb.WriteString(fmt.Sprintf("// %s\n", cleanDescription(desc)))
	}
	sb.WriteString(fmt.Sprintf("type %s struct {\n", name))

	// Generate fields from properties
	if props, ok := schema["properties"].(map[string]interface{}); ok {
		required := make(map[string]bool)
		if reqList, ok := schema["required"].([]interface{}); ok {
			for _, r := range reqList {
				if reqStr, ok := r.(string); ok {
					required[reqStr] = true
				}
			}
		}

		for propName, propData := range props {
			propMap, ok := propData.(map[string]interface{})
			if !ok {
				continue
			}

			fieldName := toCamelCase(propName)
			fieldType := getGoType(propMap)
			jsonTag := propName

			if !required[propName] {
				jsonTag += ",omitempty"
			}

			if desc, ok := propMap["description"].(string); ok && desc != "" {
				sb.WriteString(fmt.Sprintf("\t// %s\n", cleanDescription(desc)))
			}
			sb.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldType, jsonTag))
		}
	}

	sb.WriteString("}\n")

	return sb.String()
}

func getGoType(schema map[string]interface{}) string {
	// Handle array types
	if typ, ok := schema["type"].(string); ok && typ == "array" {
		if items, ok := schema["items"].(map[string]interface{}); ok {
			return "[]" + getGoType(items)
		}
		return "[]interface{}"
	}

	// Handle basic types
	if typ, ok := schema["type"].(string); ok {
		format, _ := schema["format"].(string)
		switch typ {
		case "string":
			if format == "date-time" {
				return "*time.Time"
			}
			return "string"
		case "integer":
			if format == "int64" {
				return "int64"
			}
			return "int"
		case "number":
			return "float64"
		case "boolean":
			return "bool"
		case "object":
			return "map[string]interface{}"
		}
	}

	return "interface{}"
}

func toCamelCase(s string) string {
	if s == "" {
		return s
	}
	parts := strings.Split(s, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	result := strings.Join(parts, "")
	// Handle special cases
	if result == "Id" {
		return "ID"
	}
	return result
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func cleanDescription(desc string) string {
	desc = strings.ReplaceAll(desc, "\n", " ")
	desc = strings.ReplaceAll(desc, "\r", "")
	desc = strings.TrimSpace(desc)
	if len(desc) > 200 {
		desc = desc[:197] + "..."
	}
	return desc
}
