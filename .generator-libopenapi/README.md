# libopenapi Model Generator

This tool uses [libopenapi](https://github.com/pb33f/libopenapi) to generate Go models from an OpenAPI specification.

## Usage

1. Install dependencies:
```bash
go mod download
```

2. Run the generator:
```bash
go run main.go /path/to/okta-management-APIs-oasv3-noEnums-inheritance.yaml
```

3. Generated models will be in the `./generated` directory.

## What it generates

- Go structs for each schema in the OpenAPI spec
- Proper JSON tags with omitempty for optional fields
- Type mappings (string, int, bool, arrays, objects)
- Time handling for date-time fields
- Documentation from schema descriptions
