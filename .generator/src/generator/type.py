"""
Type handling utilities for OpenAPI to Terraform type conversion.
"""

from . import formatter
from . import utils


def type_to_go(schema: dict) -> str:
    """Convert an OpenAPI schema to a Go type string."""
    if not schema:
        return "string"

    schema_type = schema.get("type", "")
    schema_format = schema.get("format", "")

    # Check for $ref
    if hasattr(schema, "__reference__"):
        ref = schema.__reference__.get("$ref", "")
        if ref:
            name = ref.split("/")[-1]
            return name

    if schema_type == "string":
        if schema_format == "date-time":
            return "*time.Time"
        if schema_format == "binary":
            return "[]byte"
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
        item_type = type_to_go(items)
        return f"[]{item_type}"
    elif schema_type == "object":
        additional = schema.get("additionalProperties")
        if additional:
            value_type = type_to_go(additional)
            return f"map[string]{value_type}"
        return "map[string]interface{}"

    return "interface{}"


def get_type_for_parameter(param: dict) -> dict:
    """Get the schema from a parameter definition."""
    if "schema" in param:
        return param["schema"]
    if "content" in param:
        for content in param.get("content", {}).values():
            if "schema" in content:
                return content["schema"]
    return {}


def get_type_for_response(responses: dict) -> dict:
    """Get the schema from response definitions."""
    for code in ["200", "201", "202"]:
        if code in responses:
            content = responses[code].get("content", {})
            for media_type in ["application/json", "*/*"]:
                if media_type in content:
                    return content[media_type].get("schema", {})
    return {}


def get_schema_from_response(responses: dict) -> dict:
    """Extract the schema from OpenAPI response definitions."""
    return get_type_for_response(responses)


def return_type(schema: dict) -> str:
    """Get the return type for a schema."""
    return type_to_go(schema)


def sort_schemas_by_type(schemas: dict) -> tuple:
    """
    Sort a flat {name: schema} dict into four buckets:
        (primitive, primitive_array, non_primitive_list, non_primitive_obj)
    """
    primitive = {}
    primitive_array = {}
    non_primitive_list = {}
    non_primitive_obj = {}

    for name, schema in schemas.items():
        if not schema:
            continue

        # Unwrap parameter wrapper {"schema": {...}}
        actual_schema = schema.get("schema", schema) if isinstance(schema, dict) else schema

        schema_type = actual_schema.get("type", "string")

        if schema_type == "array":
            items = actual_schema.get("items", {})
            if utils.is_primitive(items):
                primitive_array[name] = actual_schema
            else:
                non_primitive_list[name] = actual_schema
        elif schema_type == "object":
            non_primitive_obj[name] = actual_schema
        elif utils.is_primitive(actual_schema):
            primitive[name] = actual_schema
        else:
            primitive[name] = actual_schema

    return (primitive, primitive_array, non_primitive_list, non_primitive_obj)


def _params_list_to_dict(params_list: list) -> dict:
    """
    Convert a list of parameter dicts (each with 'name' + 'schema') into a
    flat {name: schema} dict suitable for sort_schemas_by_type.
    """
    result = {}
    for param in params_list:
        name = param.get("name", "")
        if not name:
            continue
        schema = param.get("schema", {})
        # Carry description / required onto the schema so templates can access them
        annotated = dict(schema) if schema else {}
        if "description" in param:
            annotated.setdefault("description", param["description"])
        if "required" in param:
            annotated.setdefault("required", param["required"])
        result[name] = annotated
    return result


def tf_sort_params_by_type(params) -> tuple:
    """
    Sort parameters by type for Terraform schema generation.

    Accepts either:
      - The grouped dict returned by openapi.parameters():
            {"path": [...], "query": [...], "body": [...]}
      - A legacy flat {name: schema} dict.
    """
    if isinstance(params, dict) and any(k in params for k in ("path", "query", "body")):
        # New grouped format — merge body params (the ones we put in the schema)
        flat: dict = {}
        flat.update(_params_list_to_dict(params.get("body", [])))
        # Include query params too (they may be Optional attributes)
        flat.update(_params_list_to_dict(params.get("query", [])))
        return sort_schemas_by_type(flat)
    # Legacy flat dict
    return sort_schemas_by_type(params)


def tf_sort_properties_by_type(schema: dict) -> tuple:
    """Sort schema properties by type for Terraform schema generation."""
    properties = schema.get("properties", {}) if isinstance(schema, dict) else {}
    return sort_schemas_by_type(properties)
