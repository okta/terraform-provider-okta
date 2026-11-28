package governance

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type dateNotMoreThanFiveYearsInFuture struct{}

func (dateNotMoreThanFiveYearsInFuture) Description(ctx context.Context) string {
	return "must be an ISO 8601 date not more than 5 years in the future"
}

func (dateNotMoreThanFiveYearsInFuture) MarkdownDescription(ctx context.Context) string {
	return "must be an ISO 8601 date not more than 5 years in the future"
}

func (v dateNotMoreThanFiveYearsInFuture) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid date format",
			"Value must be in ISO 8601 (RFC3339) format, e.g. 2025-10-04T13:43:40.000Z",
		)
		return
	}
	now := time.Now().UTC()
	maxDate := now.AddDate(5, 0, 0)
	parsedUTC := parsed.UTC()
	if parsedUTC.After(maxDate) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Date out of range",
			"Value must not be more than 5 years in the future",
		)
	}
}
