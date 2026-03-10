package resources

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestWrapSDKDataSource_ReadContext_PanicRecovery tests that data source ReadContext panics are recovered
func TestWrapSDKDataSource_ReadContext_PanicRecovery(t *testing.T) {
	d := &schema.Resource{
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			var nilPtr *string
			_ = *nilPtr
			return nil
		},
	}

	wrapped := WrapSDKDataSource(d)
	diags := wrapped.ReadContext(context.Background(), nil, nil)

	if len(diags) == 0 {
		t.Fatal("expected diagnostics from panic recovery, got none")
	}

	if !strings.Contains(diags[0].Summary, "Provider Crash in Read") {
		t.Errorf("expected summary to contain 'Provider Crash in Read', got %s", diags[0].Summary)
	}
}

// TestWrapSDKDataSource_LegacyRead_PanicRecovery tests that data source legacy Read panics are recovered
func TestWrapSDKDataSource_LegacyRead_PanicRecovery(t *testing.T) {
	d := &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			var nilPtr *string
			_ = *nilPtr
			return nil
		},
	}

	wrapped := WrapSDKDataSource(d)
	err := wrapped.Read(nil, nil)

	if err == nil {
		t.Fatal("expected error from panic recovery, got nil")
	}

	if !strings.Contains(err.Error(), "provider crashed during the Read operation") {
		t.Errorf("expected error to contain 'provider crashed during the Read operation', got %s", err.Error())
	}
}

// TestWrapSDKDataSources_Map tests that a map of data sources is wrapped correctly
func TestWrapSDKDataSources_Map(t *testing.T) {
	dataSources := map[string]*schema.Resource{
		"okta_datasource_a": {
			ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
				var nilPtr *string
				_ = *nilPtr
				return nil
			},
		},
		"okta_datasource_b": {
			ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
				var nilPtr *string
				_ = *nilPtr
				return nil
			},
		},
	}

	wrapped := WrapSDKDataSources(dataSources)

	// Verify both data sources are wrapped
	if len(wrapped) != 2 {
		t.Errorf("expected 2 wrapped data sources, got %d", len(wrapped))
	}

	// Call both and verify panic recovery
	for name, d := range wrapped {
		diags := d.ReadContext(context.Background(), nil, nil)
		if len(diags) == 0 {
			t.Errorf("expected diagnostics from panic recovery for %s", name)
		}
	}
}
