# Okta Terraform Provider Generator

Generate Terraform Provider code for Okta resources from OpenAPI specifications.

## Overview

This generator creates Terraform Plugin Framework resources and data sources from OpenAPI specifications. It supports generating code from a single spec file that contains multiple APIs.

## Requirements

- Python 3.10+
- Poetry
- Go (for `go fmt`)
- An OpenAPI 3.0.x specification

## Installation

```bash
cd .generator
poetry install
```

## Usage

### Configuration File

Create a YAML configuration file that defines which resources and data sources to generate:

```yaml
# config.yaml
resources:
  group:
    # API tag (optional, auto-detected from spec)
    api_tag: "Group"
    read:
      method: get
      path: /api/v1/groups/{groupId}
    create:
      method: post
      path: /api/v1/groups
    update:
      method: put
      path: /api/v1/groups/{groupId}
    delete:
      method: delete
      path: /api/v1/groups/{groupId}
  
  user:
    api_tag: "User"
    read:
      method: get
      path: /api/v1/users/{userId}
    create:
      method: post
      path: /api/v1/users
    update:
      method: put
      path: /api/v1/users/{userId}
    delete:
      method: delete
      path: /api/v1/users/{userId}

datasources:
  group:
    singular: /api/v1/groups/{groupId}
    plural: /api/v1/groups
  
  user:
    singular: /api/v1/users/{userId}
    plural: /api/v1/users
```

### Running the Generator

```bash
# From the .generator directory
poetry run python -m generator <openapi_spec_path> <config_path>

# Example
poetry run python -m generator ./specs/okta-management.yaml ./config.yaml
```

### Output

Generated files are placed in:
- Resources: `okta/fwprovider/resource_okta_{name}_generated.go`
- Data sources: `okta/fwprovider/data_source_okta_{name}_generated.go`
- Tests: `okta/fwprovider/resource_okta_{name}_generated_test.go`
- Examples: `examples/resources/okta_{name}/resource.tf`

## Multiple APIs Support

The generator handles OpenAPI specs with multiple APIs (different tags/servers). Each resource in the config can specify an `api_tag` to indicate which API client to use:

```yaml
resources:
  entitlement:
    api_tag: "GovernanceEntitlement"
    # ... CRUD operations
  
  group:
    api_tag: "Group"
    # ... CRUD operations
```

## Customization

### Templates

Templates are located in `src/generator/templates/`:

- `base_resource.j2` - Main resource template
- `data_source/base.j2` - Main data source template
- `schema.j2` - Schema generation
- `types.j2` - Go type definitions
- `utils/` - Helper macros

### Adding Custom Logic

Extend the generator by:
1. Adding new Jinja2 filters in `src/generator/formatter.py`
2. Adding new utility functions in `src/generator/utils.py`
3. Creating new templates in `src/generator/templates/`

## Notes

> **Warning**: This generator creates scaffolding code. Generated code should be reviewed and may require manual adjustments for complex resources.

The generated code follows these patterns:
- Uses Terraform Plugin Framework (not SDK v2)
- Implements `resource.ResourceWithConfigure` and `resource.ResourceWithImportState`
- Uses the Okta SDK v6 client by default
- Includes panic recovery wrappers if configured
