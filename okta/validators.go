package okta

import (
	"os"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
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

func strMaxLength(max int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		v, ok := i.(string)
		if !ok {
			return diag.Errorf("expected type of %s to be string", k)
		}
		if len(v) > max {
			return diag.Errorf("%s cannot be longer than %d characters", k, max)
		}
		return nil
	}
}

func logoFileIsValid() schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		v, ok := i.(string)
		if !ok {
			return diag.Errorf("expected type of %v to be string", k)
		}
		stat, err := os.Stat(v)
		if err != nil {
			return diag.Errorf("invalid '%s' file: %v", v, err)
		}
		if stat.Size() > 1<<20 { // should be less than 1 MB in size.
			return diag.Errorf("file '%s' should be less than 1 MB in size", v)
		}
		return nil
	}
}

func stringIsJSON(i interface{}, k cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.Errorf("expected type of %s to be string", k)
	}
	if v == "" {
		return diag.Errorf("expected %q JSON to not be empty, got %v", k, i)
	}
	if _, err := structure.NormalizeJsonString(v); err != nil {
		return diag.Errorf("%q contains an invalid JSON: %s", k, err)
	}
	return nil
}
