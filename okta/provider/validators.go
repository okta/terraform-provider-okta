package provider

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func intBetween(min, max int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		v, ok := i.(int)
		if !ok {
			return diag.Errorf("expected type of %s to be integer", k)
		}
		if v < min || v > max {
			return diag.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v)
		}
		return nil
	}
}

func intAtMost(max int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		v, ok := i.(int)
		if !ok {
			return diag.Errorf("expected type of %s to be integer", k)
		}
		if v > max {
			return diag.Errorf("expected %s to be at most (%d), got %d", k, max, v)
		}
		return nil
	}
}
