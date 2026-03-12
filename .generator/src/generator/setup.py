"""
Setup and configuration for the template environment.
"""

import pathlib
import yaml
from jinja2 import Environment, FileSystemLoader, Template
from jsonref import JsonRef

from . import openapi
from . import formatter
from . import utils
from . import type as type_module


def load_environment() -> Environment:
    """
    Load and configure the Jinja2 template environment.

    Returns:
        Configured Jinja2 Environment.
    """
    template_path = pathlib.Path(__file__).parent / "templates"
    env = Environment(
        loader=FileSystemLoader(str(template_path)),
        trim_blocks=True,
        lstrip_blocks=True,
    )

    # Register filters
    env.filters["attribute_name"] = formatter.attribute_name
    env.filters["camel_case"] = formatter.camel_case
    env.filters["sanitize_description"] = formatter.sanitize_description
    env.filters["snake_case"] = formatter.snake_case
    env.filters["untitle_case"] = formatter.untitle_case
    env.filters["variable_name"] = formatter.variable_name
    env.filters["go_to_terraform_type_formatter"] = formatter.go_to_terraform_type_formatter
    env.filters["get_terraform_schema_type"] = formatter.get_terraform_schema_type
    
    env.filters["parameter_schema"] = openapi.parameter_schema
    env.filters["parameters"] = openapi.parameters
    env.filters["is_json_api"] = openapi.is_json_api
    
    env.filters["capitalize"] = utils.capitalize
    env.filters["is_primitive"] = utils.is_primitive
    env.filters["debug"] = utils.debug_filter
    env.filters["only_keep_filters"] = utils.only_keep_filters
    env.filters["clean_response_for_datasource"] = utils.clean_response_for_datasource
    
    env.filters["response_type"] = type_module.get_type_for_response
    env.filters["get_schema_from_response"] = type_module.get_schema_from_response
    env.filters["return_type"] = type_module.return_type
    env.filters["sort_schemas_by_type"] = type_module.sort_schemas_by_type
    env.filters["tf_sort_params_by_type"] = type_module.tf_sort_params_by_type
    env.filters["tf_sort_properties_by_type"] = type_module.tf_sort_properties_by_type

    # Register globals
    env.globals["enumerate"] = enumerate
    env.globals["get_name"] = openapi.get_name
    env.globals["get_terraform_primary_id"] = openapi.get_terraform_primary_id
    env.globals["json_api_attributes_schema"] = openapi.json_api_attributes_schema
    env.globals["get_terraform_schema_type"] = formatter.get_terraform_schema_type
    env.globals["get_type_for_parameter"] = type_module.get_type_for_parameter
    env.globals["get_type"] = type_module.type_to_go
    env.globals["is_required"] = utils.is_required
    env.globals["is_computed"] = utils.is_computed
    env.globals["is_enum"] = utils.is_enum
    env.globals["is_nullable"] = utils.is_nullable
    env.globals["simple_type"] = formatter.simple_type

    env.globals["GET_OPERATION"] = utils.GET_OPERATION
    env.globals["CREATE_OPERATION"] = utils.CREATE_OPERATION
    env.globals["UPDATE_OPERATION"] = utils.UPDATE_OPERATION
    env.globals["DELETE_OPERATION"] = utils.DELETE_OPERATION

    return env


def load_templates(env: Environment) -> dict[str, Template]:
    """
    Load all template files.

    Args:
        env: The Jinja2 environment.

    Returns:
        Dictionary of template name to Template object.
    """
    templates = {
        "base": env.get_template("base_resource.j2"),
        "test": env.get_template("resource_test.j2"),
        "example": env.get_template("resource_example.j2"),
        "import": env.get_template("resource_import_example.j2"),
        "datasource": env.get_template("data_source/base.j2"),
    }
    return templates


def load(filename: str) -> dict:
    """
    Load a YAML file with JSON reference resolution.

    Args:
        filename: Path to the YAML file.

    Returns:
        Parsed and dereferenced dictionary.
    """
    path = pathlib.Path(filename)
    with path.open() as fp:
        return JsonRef.replace_refs(yaml.safe_load(fp))
