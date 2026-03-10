// Package resources contains safe wrappers for Terraform resources and data sources.
// This test file intentionally uses panic() to test the panic recovery mechanism.
//
//nolint:all
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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// ============================================
// Mock Resource for Testing
// ============================================

type mockResource struct {
	panicOnCreate bool
	panicOnRead   bool
	panicOnUpdate bool
	panicOnDelete bool
	panicMessage  string
}

func (m *mockResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mock"
}

func (m *mockResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		Attributes: map[string]resourceschema.Attribute{
			"id": resourceschema.StringAttribute{Computed: true},
		},
	}
}

func (m *mockResource) Create(_ context.Context, _ resource.CreateRequest, _ *resource.CreateResponse) {
	if m.panicOnCreate {
		panic(m.panicMessage) //nolint:R009 // intentional panic for testing SafeResource recovery
	}
}

func (m *mockResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	if m.panicOnRead {
		panic(m.panicMessage) //nolint:R009 // intentional panic for testing SafeResource recovery
	}
}

func (m *mockResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	if m.panicOnUpdate {
		panic(m.panicMessage) //nolint:R009 // intentional panic for testing SafeResource recovery
	}
}

func (m *mockResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	if m.panicOnDelete {
		panic(m.panicMessage) //nolint:R009 // intentional panic for testing SafeResource recovery
	}
}

// ============================================
// Mock DataSource for Testing
// ============================================

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
		panic(m.panicMessage) //nolint:R009 // intentional panic for testing SafeDataSource recovery
	}
}

// ============================================
// SafeResource Tests
// ============================================

func TestSafeResource_Create_PanicRecovery(t *testing.T) {
	mock := &mockResource{
		panicOnCreate: true,
		panicMessage:  "test panic in Create",
	}
	safe := NewSafeResource(mock)

	resp := &resource.CreateResponse{
		Diagnostics: diag.Diagnostics{},
	}

	// This should NOT panic - SafeResource should catch it
	safe.Create(context.Background(), resource.CreateRequest{}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Expected diagnostics to have error after panic")
	}

	// Check error message contains panic info
	errSummary := resp.Diagnostics.Errors()[0].Summary()
	if !strings.Contains(errSummary, "Provider Crash in Create") {
		t.Errorf("Expected error summary to contain 'Provider Crash in Create', got: %s", errSummary)
	}

	errDetail := resp.Diagnostics.Errors()[0].Detail()
	if !strings.Contains(errDetail, "test panic in Create") {
		t.Errorf("Expected error detail to contain panic message, got: %s", errDetail)
	}

	if !strings.Contains(errDetail, "Stack trace") {
		t.Errorf("Expected error detail to contain stack trace, got: %s", errDetail)
	}
}

func TestSafeResource_Read_PanicRecovery(t *testing.T) {
	mock := &mockResource{
		panicOnRead:  true,
		panicMessage: "test panic in Read",
	}
	safe := NewSafeResource(mock)

	resp := &resource.ReadResponse{
		Diagnostics: diag.Diagnostics{},
	}

	safe.Read(context.Background(), resource.ReadRequest{}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Expected diagnostics to have error after panic")
	}

	errSummary := resp.Diagnostics.Errors()[0].Summary()
	if !strings.Contains(errSummary, "Provider Crash in Read") {
		t.Errorf("Expected error summary to contain 'Provider Crash in Read', got: %s", errSummary)
	}
}

func TestSafeResource_Update_PanicRecovery(t *testing.T) {
	mock := &mockResource{
		panicOnUpdate: true,
		panicMessage:  "test panic in Update",
	}
	safe := NewSafeResource(mock)

	resp := &resource.UpdateResponse{
		Diagnostics: diag.Diagnostics{},
	}

	safe.Update(context.Background(), resource.UpdateRequest{}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Expected diagnostics to have error after panic")
	}

	errSummary := resp.Diagnostics.Errors()[0].Summary()
	if !strings.Contains(errSummary, "Provider Crash in Update") {
		t.Errorf("Expected error summary to contain 'Provider Crash in Update', got: %s", errSummary)
	}
}

func TestSafeResource_Delete_PanicRecovery(t *testing.T) {
	mock := &mockResource{
		panicOnDelete: true,
		panicMessage:  "test panic in Delete",
	}
	safe := NewSafeResource(mock)

	resp := &resource.DeleteResponse{
		Diagnostics: diag.Diagnostics{},
	}

	safe.Delete(context.Background(), resource.DeleteRequest{}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Expected diagnostics to have error after panic")
	}

	errSummary := resp.Diagnostics.Errors()[0].Summary()
	if !strings.Contains(errSummary, "Provider Crash in Delete") {
		t.Errorf("Expected error summary to contain 'Provider Crash in Delete', got: %s", errSummary)
	}
}

func TestSafeResource_NoPanic_PassesThrough(t *testing.T) {
	mock := &mockResource{
		panicOnCreate: false,
		panicMessage:  "",
	}
	safe := NewSafeResource(mock)

	resp := &resource.CreateResponse{
		Diagnostics: diag.Diagnostics{},
	}

	safe.Create(context.Background(), resource.CreateRequest{}, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("Expected no errors when resource doesn't panic, got: %v", resp.Diagnostics.Errors())
	}
}

func TestSafeResource_ConcurrentPanics(t *testing.T) {
	// Test that SafeResource handles concurrent panics correctly
	var wg sync.WaitGroup
	numGoroutines := 10

	results := make([]*resource.CreateResponse, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			mock := &mockResource{
				panicOnCreate: true,
				panicMessage:  "concurrent panic " + string(rune('A'+index)),
			}
			safe := NewSafeResource(mock)

			resp := &resource.CreateResponse{
				Diagnostics: diag.Diagnostics{},
			}

			safe.Create(context.Background(), resource.CreateRequest{}, resp)
			results[index] = resp
		}(i)
	}

	wg.Wait()

	// All goroutines should have completed with errors, not crashed
	for i, resp := range results {
		if !resp.Diagnostics.HasError() {
			t.Errorf("Goroutine %d: Expected error but got none", i)
		}
	}
}

func TestWrapResources(t *testing.T) {
	constructors := []func() resource.Resource{
		func() resource.Resource { return &mockResource{} },
		func() resource.Resource { return &mockResource{} },
	}

	wrapped := WrapResources(constructors)

	if len(wrapped) != len(constructors) {
		t.Errorf("Expected %d wrapped constructors, got %d", len(constructors), len(wrapped))
	}

	// Verify each wrapped constructor returns a SafeResource
	for i, constructor := range wrapped {
		r := constructor()
		if _, ok := r.(*SafeResource); !ok {
			t.Errorf("Constructor %d did not return a SafeResource", i)
		}
	}
}

// ============================================
// SafeDataSource Tests
// ============================================

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

// TestSafeResource_Create_GoroutineWithChannel tests resource panic recovery in goroutine
func TestSafeResource_Create_GoroutineWithChannel(t *testing.T) {
	mock := &mockResource{
		panicOnCreate: true,
		panicMessage:  "test panic in Create from goroutine",
	}

	safe := NewSafeResource(mock)
	resp := &resource.CreateResponse{
		Diagnostics: diag.Diagnostics{},
	}

	done := make(chan bool)

	go func() {
		safe.Create(context.Background(), resource.CreateRequest{}, resp)
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
