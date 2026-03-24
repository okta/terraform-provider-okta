package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"

	gconfig "github.com/okta/terraform-provider-okta/generator/internal/config"
	"github.com/okta/terraform-provider-okta/generator/internal/generator"
	"github.com/okta/terraform-provider-okta/generator/internal/openapi"
)

func main() {
	var (
		outputDir    string
		noGoFmt      bool
		templatesDir string
		debug        bool
	)

	flag.StringVar(&outputDir, "output", "", "Output directory for generated .go files (default: <repo_root>/okta/fwprovider)")
	flag.StringVar(&templatesDir, "templates", "", "Directory containing .go.tmpl template files (default: <binary_dir>/../templates)")
	flag.BoolVar(&noGoFmt, "no-go-fmt", false, "Skip running gofmt on generated files")
	flag.BoolVar(&debug, "debug", false, "Print detailed debug output for every resource/datasource")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tf-generator [options] <spec_path> <config_path>\n\n")
		fmt.Fprintf(os.Stderr, "  spec_path    Path to OpenAPI YAML spec file\n")
		fmt.Fprintf(os.Stderr, "  config_path  Path to generator config YAML file\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	specPath := args[0]
	configPath := args[1]

	// Resolve templates dir relative to binary or fallback to nearby path
	if templatesDir == "" {
		exe, err := os.Executable()
		if err != nil {
			exe = "."
		}
		templatesDir = filepath.Join(filepath.Dir(exe), "..", "templates")
		if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
			// Fallback: relative to current working directory
			cwd, _ := os.Getwd()
			templatesDir = filepath.Join(cwd, "templates")
		}
	}

	// Resolve output dir
	if outputDir == "" {
		cwd, _ := os.Getwd()
		outputDir = filepath.Join(cwd, "okta", "fwprovider")
	}

	fmt.Printf("Loading OpenAPI spec from %s\n", specPath)
	spec, err := openapi.Load(specPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading spec: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loading config from %s\n", configPath)
	cfg, err := loadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Go output directory: %s\n", outputDir)
	fmt.Printf("Templates directory: %s\n", templatesDir)

	var dbgLogger *log.Logger
	if debug {
		dbgLogger = log.New(os.Stderr, "[DEBUG] ", 0)
		dbgLogger.Println("Debug mode enabled")
		dbgLogger.Printf("Spec paths loaded: %d", len(spec.Paths))
		dbgLogger.Printf("Resources in config: %d", len(cfg.Resources))
		dbgLogger.Printf("DataSources in config: %d", len(cfg.DataSources))
	}

	gen, err := generator.New(spec, templatesDir, outputDir, !noGoFmt, dbgLogger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing generator: %v\n", err)
		os.Exit(1)
	}

	// Generate data sources
	dsNames := sortedKeys(cfg.DataSources)
	fmt.Printf("Found %d data source(s) to generate\n", len(dsNames))
	for _, name := range dsNames {
		fmt.Printf("  Generating data source: okta_%s\n", name)
		if err := gen.GenerateDataSource(name, cfg.DataSources[name]); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating datasource %s: %v\n", name, err)
			os.Exit(1)
		}
	}

	// Generate resources
	resNames := sortedKeys(cfg.Resources)
	fmt.Printf("Found %d resource(s) to generate\n", len(resNames))
	for _, name := range resNames {
		fmt.Printf("  Generating resource: okta_%s\n", name)
		if err := gen.GenerateResource(name, cfg.Resources[name]); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating resource %s: %v\n", name, err)
			os.Exit(1)
		}
	}
	_ = debug // consumed via dbgLogger

	fmt.Println("Generation complete!")
}

func loadConfig(path string) (*gconfig.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}
	var cfg gconfig.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config YAML: %w", err)
	}
	if cfg.Resources == nil {
		cfg.Resources = make(map[string]gconfig.ResourceConfig)
	}
	if cfg.DataSources == nil {
		cfg.DataSources = make(map[string]gconfig.DataSourceConfig)
	}
	return &cfg, nil
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
