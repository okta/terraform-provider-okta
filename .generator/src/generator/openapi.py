"""
OpenAPI specification parsing and processing.
"""

from .utils import (
    GET_OPERATION,
    CREATE_OPERATION,
    UPDATE_OPERATION,
    DELETE_OPERATION,
    get_api_tag_from_operation,
    get_primary_id_from_path,
)


def get_name(schema: dict) -> str:
    """Get the name from a schema reference."""
    if hasattr(schema, "__reference__"):
        ref = schema.__reference__.get("$ref", "")
        if ref:
            return ref.split("/")[-1]
    return None


def _get_response_schema(operation_schema: dict) -> dict:
    """
    Extract the primary JSON response schema from an operation's responses.
    Returns the resolved schema dict or an empty dict if not found.
    """
    for status_code in ("200", "201", "202"):
        resp = operation_schema.get("responses", {}).get(status_code, {})
        content = resp.get("content", {})
        json_content = content.get("application/json", {})
        schema = json_content.get("schema", {})
        if schema:
            # Unwrap array wrappers so callers always get the object schema
            if schema.get("type") == "array":
                return schema.get("items", {})
            return schema
    return {}


def get_resources(spec: dict, config: dict) -> dict:
    """
    Generate a dictionary of resources and their CRUD operations.
    Supports multiple APIs from a single spec file.

    Each resource entry in config may contain:
        api_tag: (optional) Override the API tag used to select the SDK client.
        read / create / update / delete:
            method: HTTP method (get/post/put/patch/delete)
            path:   Path in the OpenAPI spec (e.g. /api/v1/groups/{groupId})

    Args:
        spec: The OpenAPI specification.
        config: The configuration defining resources to generate.

    Returns:
        A dictionary of resources with their operation schemas.
    """
    resources_to_generate = {}

    for resource_name in config.get("resources", []):
        resource_config = config["resources"][resource_name]

        # Get optional API tag override
        api_tag = resource_config.get("api_tag")

        for crud_operation, value in resource_config.items():
            # Skip non-CRUD config keys
            if crud_operation in ("api_tag", "description", "options"):
                continue

            # value must be a dict with method + path
            if not isinstance(value, dict):
                continue

            method = value.get("method", "").lower()
            path = value.get("path", "")

            if not path or path not in spec.get("paths", {}):
                print(f"Warning: Path {path!r} not found in spec for {resource_name}.{crud_operation}")
                continue

            path_item = spec["paths"][path]
            if method not in path_item:
                print(f"Warning: Method {method!r} not found at {path!r} for {resource_name}.{crud_operation}")
                continue

            operation_schema = path_item[method]

            # Map config key → canonical operation name
            operation_map = {
                "read": GET_OPERATION,
                "create": CREATE_OPERATION,
                "update": UPDATE_OPERATION,
                "delete": DELETE_OPERATION,
            }
            operation = operation_map.get(crud_operation)
            if not operation:
                print(f"Warning: Unknown operation {crud_operation!r} for {resource_name}")
                continue

            # Detect API tag from spec if not supplied in config
            if not api_tag:
                api_tag = get_api_tag_from_operation(operation_schema)

            # Resolve response schema for template access
            response_schema = _get_response_schema(operation_schema)

            resources_to_generate.setdefault(resource_name, {
                "api_tag": api_tag,
                "config": resource_config,
            })[operation] = {
                "schema": operation_schema,
                "path": path,
                "method": method,
                "operation_id": operation_schema.get("operationId", ""),
                "response_schema": response_schema,
            }

    return resources_to_generate


def get_data_sources(spec: dict, config: dict) -> dict:
    """
    Generate a dictionary of data sources and their endpoints.
    Supports multiple APIs from a single spec file.

    Each data source entry in config may contain:
        api_tag: (optional) Override the API tag.
        singular:
            method: HTTP method
            path:   Path to the singular (get-by-id) endpoint
        plural:
            method: HTTP method
            path:   Path to the list endpoint (optional)

    Args:
        spec: The OpenAPI specification.
        config: The configuration defining data sources to generate.

    Returns:
        A dictionary of data sources with their operation schemas.
    """
    data_sources_to_generate = {}

    for ds_name in config.get("datasources", []):
        ds_config = config["datasources"][ds_name]

        # Get optional API tag override
        api_tag = ds_config.get("api_tag")

        # ── Singular endpoint ──────────────────────────────────────────────
        singular_entry = ds_config.get("singular")
        if isinstance(singular_entry, dict):
            singular_path = singular_entry.get("path", "")
            singular_method = singular_entry.get("method", "get").lower()
        else:
            # Legacy: string value treated as path with GET
            singular_path = singular_entry or ""
            singular_method = "get"

        if singular_path and singular_path in spec.get("paths", {}):
            singular_op = spec["paths"][singular_path].get(singular_method, {})
            if not api_tag:
                api_tag = get_api_tag_from_operation(singular_op)

            response_schema = _get_response_schema(singular_op)
            data_sources_to_generate.setdefault(ds_name, {
                "api_tag": api_tag,
                "config": ds_config,
            })["singular"] = {
                "schema": singular_op,
                "path": singular_path,
                "method": singular_method,
                "operation_id": singular_op.get("operationId", ""),
                "response_schema": response_schema,
            }
        elif singular_path:
            print(f"Warning: Singular path {singular_path!r} not found in spec for datasource {ds_name!r}")

        # ── Plural / list endpoint (optional) ─────────────────────────────
        plural_entry = ds_config.get("plural")
        if isinstance(plural_entry, dict):
            plural_path = plural_entry.get("path", "")
            plural_method = plural_entry.get("method", "get").lower()
        else:
            plural_path = plural_entry or ""
            plural_method = "get"

        if plural_path and plural_path in spec.get("paths", {}):
            plural_op = spec["paths"][plural_path].get(plural_method, {})

            response_schema = _get_response_schema(plural_op)
            data_sources_to_generate.setdefault(ds_name, {
                "api_tag": api_tag,
                "config": ds_config,
            })["plural"] = {
                "schema": plural_op,
                "path": plural_path,
                "method": plural_method,
                "operation_id": plural_op.get("operationId", ""),
                "response_schema": response_schema,
            }
        elif plural_path:
            print(f"Warning: Plural path {plural_path!r} not found in spec for datasource {ds_name!r}")

    return data_sources_to_generate


def get_terraform_primary_id(operations: dict, operation_key: str = UPDATE_OPERATION) -> dict:
    """
    Get the primary ID parameter from operations.

    Looks in the given operation_key first (default: 'update'); falls back to
    'read' if not found.  Returns a dict with 'name' and 'schema'.

    Args:
        operations: The operations dictionary for a resource or data source.
        operation_key: Which operation to extract the ID from.

    Returns:
        A dictionary with 'name' and 'schema' for the primary ID,
        or an empty dict if no path parameters exist.
    """
    # Try requested key, then fall back to read/get
    for key in (operation_key, GET_OPERATION, "singular"):
        operation = operations.get(key, {})
        path = operation.get("path", "")
        if path:
            break
    else:
        return {}

    primary_id_name = get_primary_id_from_path(path)
    if not primary_id_name:
        return {}

    # Try to locate the schema for this parameter in the operation
    operation_schema = operation.get("schema", {})
    path_params_list = parameters(operation_schema).get("path", [])
    primary_id_schema = {"type": "string"}
    for p in path_params_list:
        if p.get("name") == primary_id_name:
            primary_id_schema = parameter_schema(p)
            break

    return {
        "name": primary_id_name,
        "schema": primary_id_schema,
    }


def parameters(operation: dict) -> dict:
    """
    Extract parameters from an operation schema, grouped by location.

    Returns a dict with keys:
        "path"  → list of path parameter dicts
        "query" → list of query parameter dicts
        "body"  → list of body property dicts (expanded from requestBody)
    """
    result: dict = {"path": [], "query": [], "body": []}

    # Path and query parameters
    for param in operation.get("parameters", []):
        loc = param.get("in", "")
        if loc == "path":
            result["path"].append(param)
        elif loc == "query":
            result["query"].append(param)

    # Request body – expand JSON properties as individual body params
    if "requestBody" in operation:
        request_body = operation["requestBody"]
        content = request_body.get("content", {})

        if "multipart/form-data" in content:
            form_schema = content["multipart/form-data"].get("schema", {})
            required_fields = form_schema.get("required", [])
            for prop_name, prop_schema in form_schema.get("properties", {}).items():
                result["body"].append({
                    "in": "form",
                    "name": prop_name,
                    "schema": prop_schema,
                    "description": prop_schema.get("description", ""),
                    "required": prop_name in required_fields,
                })
        elif "application/json" in content:
            json_schema = content["application/json"].get("schema", {})
            required_fields = json_schema.get("required", [])
            for prop_name, prop_schema in json_schema.get("properties", {}).items():
                result["body"].append({
                    "in": "body",
                    "name": prop_name,
                    "schema": prop_schema,
                    "description": prop_schema.get("description", ""),
                    "required": prop_name in required_fields,
                })

    return result


def parameter_schema(param: dict) -> dict:
    """
    Extract the schema from a parameter definition.

    Args:
        param: The parameter definition.

    Returns:
        The schema dictionary.
    """
    if "schema" in param:
        return param["schema"]
    if "content" in param:
        for content_type in param["content"].values():
            if "schema" in content_type:
                return content_type["schema"]
    return {}


def is_json_api(schema: dict) -> bool:
    """
    Check if a schema follows JSON:API conventions.

    Args:
        schema: The schema to check.

    Returns:
        True if the schema appears to be JSON:API format.
    """
    if not schema:
        return False

    properties = schema.get("properties", {})
    if "data" in properties:
        data_props = properties["data"].get("properties", {})
        if "type" in data_props and "attributes" in data_props:
            return True

    return False


def json_api_attributes_schema(schema: dict) -> dict:
    """
    Extract the attributes schema from a JSON:API schema.

    Args:
        schema: The JSON:API schema.

    Returns:
        The attributes schema.
    """
    return (
        schema.get("properties", {})
        .get("data", {})
        .get("properties", {})
        .get("attributes", {})
    )
