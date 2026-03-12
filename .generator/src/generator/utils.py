"""
Utility functions and constants for the generator.
"""

# Operation constants
GET_OPERATION = "read"
CREATE_OPERATION = "create"
UPDATE_OPERATION = "update"
DELETE_OPERATION = "delete"

# Terraform type mappings
OPENAPI_TO_TF_TYPES = {
    "string": "String",
    "integer": "Int64",
    "number": "Float64",
    "boolean": "Bool",
    "array": "List",
    "object": "Object",
}


def capitalize(s: str) -> str:
    """Capitalize the first letter of a string."""
    if not s:
        return s
    return s[0].upper() + s[1:]


def is_primitive(schema: dict) -> bool:
    """Check if a schema represents a primitive type."""
    if not schema:
        return False
    schema_type = schema.get("type", "")
    return schema_type in ("string", "integer", "number", "boolean")


def is_required(schema: dict) -> bool:
    """Check if a schema property is required."""
    return schema.get("required", False)


def is_computed(schema: dict) -> bool:
    """Check if a schema property is computed (read-only)."""
    return schema.get("readOnly", False)


def is_enum(schema: dict) -> bool:
    """Check if a schema is an enum."""
    return "enum" in schema


def is_nullable(schema: dict) -> bool:
    """Check if a schema is nullable."""
    return schema.get("nullable", False)


def debug_filter(value):
    """Debug filter to print values during template rendering."""
    print(f"DEBUG: {value}")
    return value


def only_keep_filters(params: dict) -> dict:
    """Filter parameters to only keep query parameters (filters)."""
    return {
        name: param for name, param in params.items() 
        if param.get("in") == "query"
    }


def clean_response_for_datasource(schema: dict) -> dict:
    """Clean up response schema for datasource generation."""
    if not schema:
        return {}
    
    # If it's a JSON:API response, extract the data attributes
    if "properties" in schema:
        data = schema.get("properties", {}).get("data", {})
        if "properties" in data:
            attrs = data.get("properties", {}).get("attributes", {})
            if attrs and "properties" in attrs:
                return attrs
    
    return schema.get("properties", schema)


def get_api_tag_from_operation(operation: dict) -> str:
    """Extract the API tag from an operation schema."""
    tags = operation.get("tags", [])
    if tags:
        # Clean up the tag name for use as Go identifier
        tag = tags[0].replace(" ", "").replace("-", "")
        return tag
    return "Default"


def extract_path_parameters(path: str) -> list[str]:
    """Extract path parameters from a URL path."""
    import re
    return re.findall(r'\{(\w+)\}', path)


def get_primary_id_from_path(path: str) -> str:
    """Extract the primary ID parameter from a path (usually the last path param)."""
    params = extract_path_parameters(path)
    if params:
        return params[-1]
    return "id"
