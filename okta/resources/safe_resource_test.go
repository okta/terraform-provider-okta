package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// mockPanicResource is a test resource that panics on various operations
type mockPanicResource struct {
	panicOnCreate bool
	panicOnRead   bool
	panicOnUpdate bool
	panicOnDelete bool
}

func (m *mockPanicResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mock_panic"
}
