---
name: run-generator
description: Use this skill when the user wants to build the tf-generator binary, run the generator, and clean up generated files with goimports. Trigger on phrases like "run the generator", "regenerate", "rebuild the generator", "run goimports on generated files", "build and run the generator", or any time you need to regenerate provider files from config_full.yaml without adding new config entries.
user_invocable: true
---

# Build Generator, Run Generator, Fix Imports

Three-step workflow: build binary → run generator (output to this repo) → run goimports on all generated files.

---

## STEP 1 — Build the generator binary

The generator source lives in `dcp-terraform-codegen`. Always rebuild before running so the binary reflects the latest template changes.

```bash
cd /Users/aditya.anand/okta-sdks/dcp-terraform-codegen/.generator/go-generator
GOTOOLCHAIN=auto go build -o tf-generator ./cmd/...
```

If the build fails, the most common causes are:
- A new flag added to `config_full.yaml` that has no corresponding struct field in `internal/config/config.go` — add the yaml-tagged field.
- A template syntax error in `templates/resource.go.tmpl` or another `.tmpl` file — read the error, fix the template.

---

## STEP 2 — Run the generator

Generated files are written to this repo's staging area (`terraform-provider-okta/.generator/go-generator/okta/fwprovider/`), **not** back into `dcp-terraform-codegen`.

```bash
cd /Users/aditya.anand/okta-sdks/dcp-terraform-codegen/.generator/go-generator

./tf-generator \
  -output /Users/aditya.anand/okta-sdks/terraform-provider-okta/.generator/go-generator/okta/fwprovider \
  -templates ./templates \
  /Users/aditya.anand/okta-sdks/okta-sdk-golang/.generator/okta-management-APIs-oasv3-noEnums-inheritance.yaml \
  /Users/aditya.anand/okta-sdks/dcp-terraform-codegen/.generator/config_full.yaml
```

Verify the expected files were (re)created:

```bash
ls -la /Users/aditya.anand/okta-sdks/terraform-provider-okta/.generator/go-generator/okta/fwprovider/*_generated.go
```

The generator overwrites existing `*_generated.go` files in place. Any file in that directory whose name ends with `_generated.go` was produced by this generator.

---

## STEP 3 — Run goimports on all generated files

`goimports` removes unused imports, adds missing ones, and sorts import blocks.

```bash
~/go/bin/goimports -w \
  /Users/aditya.anand/okta-sdks/terraform-provider-okta/.generator/go-generator/okta/fwprovider/*_generated.go
```

If you need to run it on a single file:

```bash
~/go/bin/goimports -w \
  /Users/aditya.anand/okta-sdks/terraform-provider-okta/.generator/go-generator/okta/fwprovider/resource_okta_<name>_generated.go
```

---

## Key paths

| Item | Path |
|------|------|
| Generator source | `/Users/aditya.anand/okta-sdks/dcp-terraform-codegen/.generator/go-generator/` |
| Generator binary | `/Users/aditya.anand/okta-sdks/dcp-terraform-codegen/.generator/go-generator/tf-generator` |
| Generator templates | `/Users/aditya.anand/okta-sdks/dcp-terraform-codegen/.generator/go-generator/templates/` |
| Generator config | `/Users/aditya.anand/okta-sdks/dcp-terraform-codegen/.generator/config_full.yaml` |
| OAS spec | `/Users/aditya.anand/okta-sdks/okta-sdk-golang/.generator/okta-management-APIs-oasv3-noEnums-inheritance.yaml` |
| Generated output (this repo) | `/Users/aditya.anand/okta-sdks/terraform-provider-okta/.generator/go-generator/okta/fwprovider/` |
| goimports binary | `~/go/bin/goimports` |
| Build env | `GOTOOLCHAIN=auto` (required — go.mod specifies go 1.25) |

---

## Checklist

- [ ] Generator built without errors
- [ ] Generator ran and produced/updated `*_generated.go` files in `.generator/go-generator/okta/fwprovider/`
- [ ] `goimports` run on all generated files
