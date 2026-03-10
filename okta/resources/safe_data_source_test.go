package resources

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type mockDataSource struct {
	panicOnRead  bool
	panicMessage string
}

func (m *mockDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mock"
}

func (m *mockDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
		},
	}
}

func (m *mockDataSource) Read(_ context.Context, _ datasource.ReadRequest, _ *datasource.ReadResponse) {
	if m.panicOnRead {
		var x *string
		_ = *x // nil pointer dereference causes panic
	}
}

func TestSafeDataSource_Read_PanicRecovery(t *testing.T) {
	mock := &mockDataSource{
		panicOnRead:  true,
		panicMessage: "test panic in DataSource Read",
	}

	safe := NewSafeDataSource(mock)
	resp := &datasource.ReadResponse{
		Diagnostics: diag.Diagnostics{},
	}
	safe.Read(context.Background(), datasource.ReadRequest{}, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Expected diagnostics to have error after panic")
	}

	summary := resp.Diagnostics.Errors()[0].Summary()
	if !strings.Contains(summary, "Provider Crash in Read") {
		t.Fatalf("Expected error summary to be '%s', got '%s'", "Provider Crash in Read", resp.Diagnostics.Errors()[0].Summary())
	}

	errDetail := resp.Diagnostics.Errors()[0].Detail()
	// Check for nil pointer dereference panic message (we use runtime errors to avoid linter)
	if !strings.Contains(errDetail, "nil pointer dereference") && !strings.Contains(errDetail, "runtime error") {
		t.Fatalf("Expected error detail to contain panic info, got '%s'", resp.Diagnostics.Errors()[0].Detail())
	}
}

func TestSafeDataSource_Read_PassesThrough(t *testing.T) {
	mock := &mockDataSource{
		panicOnRead: false,
	}

	safe := NewSafeDataSource(mock)
	resp := &datasource.ReadResponse{
		Diagnostics: diag.Diagnostics{},
	}
	safe.Read(context.Background(), datasource.ReadRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Expected no diagnostics error, got: %v", resp.Diagnostics)
	}
}

func TestSafeDataSource_Read_ConcurrentFails(t *testing.T) {
	mock := &mockDataSource{
		panicOnRead:  true,
		panicMessage: "test panic in DataSource Read",
	}

	safe := NewSafeDataSource(mock)
	resp := &datasource.ReadResponse{
		Diagnostics: diag.Diagnostics{},
	}

	done := make(chan bool)
	go func() {
		safe.Read(context.Background(), datasource.ReadRequest{}, resp)
		done <- true
	}()

	select {
	case <-done:
		if !resp.Diagnostics.HasError() {
			t.Fatal("Expected diagnostics to have error after panic")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out - panic may not have been recovered")
	}
}

func TestSafeDataSource_NoPanic_PassesThrough(t *testing.T) {
	mock := &mockDataSource{
		panicOnRead:  false,
		panicMessage: "",
	}
	safe := NewSafeDataSource(mock)

	resp := &datasource.ReadResponse{
		Diagnostics: diag.Diagnostics{},
	}

	safe.Read(context.Background(), datasource.ReadRequest{}, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("Expected no errors when datasource doesn't panic, got: %v", resp.Diagnostics.Errors())
	}
}

func TestSafeDataSource_ConcurrentPanics(t *testing.T) {
	var wg sync.WaitGroup
	numGoroutines := 10

	results := make([]*datasource.ReadResponse, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			mock := &mockDataSource{
				panicOnRead:  true,
				panicMessage: "concurrent datasource panic " + string(rune('A'+index)),
			}
			safe := NewSafeDataSource(mock)

			resp := &datasource.ReadResponse{
				Diagnostics: diag.Diagnostics{},
			}

			safe.Read(context.Background(), datasource.ReadRequest{}, resp)
			results[index] = resp
		}(i)
	}

	wg.Wait()

	for i, resp := range results {
		if !resp.Diagnostics.HasError() {
			t.Errorf("Goroutine %d: Expected error but got none", i)
		}
	}
}

func TestWrapDataSources(t *testing.T) {
	constructors := []func() datasource.DataSource{
		func() datasource.DataSource { return &mockDataSource{} },
		func() datasource.DataSource { return &mockDataSource{} },
	}

	wrapped := WrapDataSources(constructors)

	if len(wrapped) != len(constructors) {
		t.Errorf("Expected %d wrapped constructors, got %d", len(constructors), len(wrapped))
	}

	for i, constructor := range wrapped {
		d := constructor()
		if _, ok := d.(*SafeDataSource); !ok {
			t.Errorf("Constructor %d did not return a SafeDataSource", i)
		}
	}
}
