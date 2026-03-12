"""
CLI for the Okta Terraform Provider Generator.
"""

import pathlib
import subprocess
import click

from jinja2 import Template

from . import setup
from . import openapi

# Path of this file inside the package: .generator/src/generator/cli.py
# Repo root is 4 levels up.
_REPO_ROOT = pathlib.Path(__file__).resolve().parent.parent.parent.parent


@click.command()
@click.argument(
    "spec_path",
    type=click.Path(exists=True, file_okay=True, dir_okay=False, path_type=pathlib.Path),
)
@click.argument(
    "config_path",
    type=click.Path(exists=True, file_okay=True, dir_okay=False, path_type=pathlib.Path),
)
@click.option("--go-fmt/--no-go-fmt", default=True, help="Run go fmt on generated files")
@click.option(
    "--output", "-o",
    type=click.Path(path_type=pathlib.Path),
    help="Output directory for .go files (default: <repo_root>/okta/fwprovider)",
)
@click.option(
    "--examples-output", "-e",
    type=click.Path(path_type=pathlib.Path),
    help="Output directory for example .tf files (default: <repo_root>/examples)",
)
def cli(
    spec_path: pathlib.Path,
    config_path: pathlib.Path,
    go_fmt: bool,
    output: pathlib.Path,
    examples_output: pathlib.Path,
):
    """
    Generate Terraform Provider code from OpenAPI specification.

    SPEC_PATH: Path to the OpenAPI specification YAML file.
    CONFIG_PATH: Path to the generator configuration YAML file.

    This generator supports multiple APIs from a single spec file.
    Each resource/datasource can specify an api_tag in the config
    to indicate which API client to use.
    """
    click.echo(f"Loading OpenAPI spec from {spec_path}")
    click.echo(f"Loading config from {config_path}")

    # Load environment and templates
    env = setup.load_environment()
    templates = setup.load_templates(env)

    # Load spec and config
    spec = setup.load(str(spec_path))
    config = setup.load(str(config_path))

    # Resolve output directories
    resolved_output = (output or _REPO_ROOT / "okta" / "fwprovider").resolve()
    resolved_examples = (examples_output or _REPO_ROOT / "examples").resolve()

    click.echo(f"Go output directory:      {resolved_output}")
    click.echo(f"Examples output directory: {resolved_examples}")

    # Generate data sources
    data_sources = openapi.get_data_sources(spec, config)
    click.echo(f"Found {len(data_sources)} data source(s) to generate")

    for name, data_source in data_sources.items():
        click.echo(f"  Generating data source: okta_{name}")
        generate_data_source(
            name=name,
            data_source=data_source,
            templates=templates,
            output=resolved_output,
            go_fmt=go_fmt,
        )

    # Generate resources
    resources = openapi.get_resources(spec, config)
    click.echo(f"Found {len(resources)} resource(s) to generate")

    for name, resource in resources.items():
        click.echo(f"  Generating resource: okta_{name}")
        generate_resource(
            name=name,
            resource=resource,
            templates=templates,
            output=resolved_output,
            examples_output=resolved_examples,
            go_fmt=go_fmt,
        )

    click.echo("Generation complete!")


def generate_data_source(
    name: str,
    data_source: dict,
    templates: dict[str, Template],
    output: pathlib.Path,
    go_fmt: bool,
) -> None:
    """
    Generate a data source file.

    Args:
        name: The data source name.
        data_source: The data source operations and config.
        templates: The loaded templates.
        output: Output directory path.
        go_fmt: Whether to run go fmt.
    """
    output.mkdir(parents=True, exist_ok=True)
    filename = output / f"data_source_okta_{name}_generated.go"

    api_tag = data_source.get("api_tag", "Default")

    with filename.open("w") as fp:
        fp.write(templates["datasource"].render(
            name=name,
            operations=data_source,
            api_tag=api_tag,
        ))

    if go_fmt:
        subprocess.call(["go", "fmt", str(filename)])


def generate_resource(
    name: str,
    resource: dict,
    templates: dict[str, Template],
    output: pathlib.Path,
    examples_output: pathlib.Path,
    go_fmt: bool,
) -> None:
    """
    Generate resource files: main .go, test .go, example .tf, import .sh.

    Args:
        name: The resource name.
        resource: The resource operations and config.
        templates: The loaded templates.
        output: Output directory for Go files.
        examples_output: Root examples directory (resources go in <examples_output>/resources/okta_<name>/).
        go_fmt: Whether to run go fmt.
    """
    output.mkdir(parents=True, exist_ok=True)

    api_tag = resource.get("api_tag", "Default")
    ctx = dict(name=name, operations=resource, api_tag=api_tag)

    # Main resource file
    resource_filename = output / f"resource_okta_{name}_generated.go"
    with resource_filename.open("w") as fp:
        fp.write(templates["base"].render(**ctx))
    if go_fmt:
        subprocess.call(["go", "fmt", str(resource_filename)])

    # Test file
    test_filename = output / f"resource_okta_{name}_generated_test.go"
    with test_filename.open("w") as fp:
        fp.write(templates["test"].render(**ctx))
    if go_fmt:
        subprocess.call(["go", "fmt", str(test_filename)])

    # Examples
    examples_dir = examples_output / "resources" / f"okta_{name}"
    examples_dir.mkdir(parents=True, exist_ok=True)

    example_filename = examples_dir / "resource.tf"
    with example_filename.open("w") as fp:
        fp.write(templates["example"].render(**ctx))

    import_filename = examples_dir / "import.sh"
    with import_filename.open("w") as fp:
        fp.write(templates["import"].render(**ctx))


if __name__ == "__main__":
    cli()
