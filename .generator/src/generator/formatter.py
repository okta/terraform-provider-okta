"""
String formatting utilities for the generator.
"""

import re


def snake_case(s: str) -> str:
    """Convert a string to snake_case."""
    if not s:
        return s
    # Handle acronyms and camelCase
    s = re.sub(r'([A-Z]+)([A-Z][a-z])', r'\1_\2', s)
    s = re.sub(r'([a-z\d])([A-Z])', r'\1_\2', s)
    s = re.sub(r'[-\s]', '_', s)
    return s.lower()


def camel_case(s: str) -> str:
    """Convert a string to CamelCase (PascalCase)."""
    if not s:
        return s
    # Split on underscores, hyphens, and spaces
    parts = re.split(r'[_\-\s]+', s)
    return ''.join(word.capitalize() for word in parts if word)


def untitle_case(s: str) -> str:
    """Convert first character to lowercase (camelCase from PascalCase)."""
    if not s:
        return s
    return s[0].lower() + s[1:]


def variable_name(s: str) -> str:
    """Convert a string to a valid Go variable name (camelCase)."""
    return untitle_case(camel_case(s))


def attribute_name(s: str) -> str:
    """Convert a string to a Terraform attribute name (snake_case)."""
    return snake_case(s)


def sanitize_description(s: str) -> str:
    """Sanitize a description string for use in Go code."""
    if not s:
        return ""
    # Escape backticks and newlines for Go string literals
    s = s.replace('`', "'")
    s = s.replace('\n', ' ')
    s = s.replace('\r', '')
    s = s.replace('"', '\\"')
    return s.strip()


def go_to_terraform_type_formatter(go_type: str) -> str:
    """Convert a Go type to a Terraform Framework type."""
    type_map = {
        "string": "types.String",
        "int": "types.Int64",
        "int32": "types.Int64",
        "int64": "types.Int64",
        "float32": "types.Float64",
        "float64": "types.Float64",
        "bool": "types.Bool",
        "[]string": "types.List",
        "[]int": "types.List",
        "map[string]string": "types.Map",
    }
    return type_map.get(go_type, "types.String")


def get_terraform_schema_type(schema: dict) -> str:
    """Get the Terraform schema type from an OpenAPI schema."""
    if not schema:
        return "String"
    
    schema_type = schema.get("type", "string")
    schema_format = schema.get("format", "")
    
    if schema_type == "string":
        return "String"
    elif schema_type == "integer":
        return "Int64"
    elif schema_type == "number":
        return "Float64"
    elif schema_type == "boolean":
        return "Bool"
    elif schema_type == "array":
        return "List"
    elif schema_type == "object":
        return "Object"
    
    return "String"


def simple_type(schema: dict) -> str:
    """Get a simple Go type from an OpenAPI schema."""
    if not schema:
        return "string"
    
    schema_type = schema.get("type", "string")
    schema_format = schema.get("format", "")
    
    if schema_type == "string":
        if schema_format == "date-time":
            return "time.Time"
        return "string"
    elif schema_type == "integer":
        if schema_format == "int32":
            return "int32"
        return "int64"
    elif schema_type == "number":
        if schema_format == "float":
            return "float32"
        return "float64"
    elif schema_type == "boolean":
        return "bool"
    elif schema_type == "array":
        items = schema.get("items", {})
        item_type = simple_type(items)
        return f"[]{item_type}"
    elif schema_type == "object":
        return "map[string]interface{}"
    
    return "string"
