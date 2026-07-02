package resources

import (
	"context"
	"fmt"
	"reflect"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Ensure SafeDataSource implements all required interfaces
var (
	_ datasource.DataSource              = &SafeDataSource{}
	_ datasource.DataSourceWithConfigure = &SafeDataSource{}
)

// SafeDataSource wraps a data source with panic recovery to prevent provider crashes
type SafeDataSource struct {
	underlying     datasource.DataSource
	nameOnce       sync.Once
	dataSourceName atomic.Value // string
}

// NewSafeDataSource creates a new SafeDataSource wrapper around the given data source
func NewSafeDataSource(d datasource.DataSource) datasource.DataSource {
	return &SafeDataSource{underlying: d}
}

// WrapDataSources wraps multiple data source constructors with SafeDataSource
func WrapDataSources(constructors []func() datasource.DataSource) []func() datasource.DataSource {
	wrapped := make([]func() datasource.DataSource, len(constructors))
	for i, constructor := range constructors {
		c := constructor // capture loop variable
		wrapped[i] = func() datasource.DataSource {
			return NewSafeDataSource(c())
		}
	}
	return wrapped
}

// recoverPanic handles panic recovery and adds appropriate diagnostics
func (s *SafeDataSource) recoverPanic(diags *diag.Diagnostics, operation string) {
	if r := recover(); r != nil {
		stackTrace := string(debug.Stack())
		dataSourceName, _ := s.dataSourceName.Load().(string)
		if dataSourceName == "" && s.underlying != nil {
			dataSourceName = typeBaseName(reflect.TypeOf(s.underlying))
		}
		if dataSourceName == "" {
			dataSourceName = "unknown"
		}

		diags.AddError(
			fmt.Sprintf("Provider Crash in %s operation of data source %s", operation, dataSourceName),
			fmt.Sprintf(
				"The Terraform Provider Okta crashed during the %s operation of data source %s.\n\n"+
					"Please check if this issue has already been reported on\n"+
					"https://github.com/okta/terraform-provider-okta/issues\n"+
					"or create a new issue with this stack trace.\n"+
					"Error: %v\n\nStack trace:\n%s\n\n",
				operation, dataSourceName, r, stackTrace,
			),
		)
	}
}

// Metadata delegates to the underlying data source
func (s *SafeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	s.underlying.Metadata(ctx, req, resp)
	if resp.TypeName != "" {
		s.nameOnce.Do(func() {
			s.dataSourceName.Store(resp.TypeName)
		})
	}
}

// Schema delegates to the underlying data source
func (s *SafeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	s.underlying.Schema(ctx, req, resp)
}

// Read wraps the underlying Read with panic recovery
func (s *SafeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer s.recoverPanic(&resp.Diagnostics, "Read")
	s.underlying.Read(ctx, req, resp)
}

// Configure delegates to the underlying data source if it implements DataSourceWithConfigure
func (s *SafeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	defer s.recoverPanic(&resp.Diagnostics, "Configure")
	if dc, ok := s.underlying.(datasource.DataSourceWithConfigure); ok {
		dc.Configure(ctx, req, resp)
	}
}
