package okta

import (
	"net/url"
	"regexp"
	"strings"

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

func intAtLeast(min int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		v, ok := i.(int)
		if !ok {
			return diag.Errorf("expected type of %s to be integer", k)
		}
		if v < min {
			return diag.Errorf("expected %s to be at least (%d), got %d", k, min, v)
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

func stringInSlice(valid []string) schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		v, ok := i.(string)
		if !ok {
			return diag.Errorf("expected type of %v to be string", k)
		}
		for _, str := range valid {
			if v == str {
				return nil
			}
		}
		return diag.Errorf("expected %v to be one of %v, got %s", k, strings.Join(valid, ","), v)
	}
}

var validURLSchemes = []string{"http", "https"}

func stringIsURL(schemes ...string) schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		val := k[0].(cty.GetAttrStep).Name // 'k' should not be map, slice or array
		v, ok := i.(string)
		if !ok {
			return diag.Errorf("invalid URL: expected type of '%s' to be string", val)
		}
		if v == "" {
			return diag.Errorf("invalid URL: expected '%s' to be not empty", val)
		}
		u, err := url.Parse(v)
		if err != nil {
			return diag.Errorf("invalid URL: expected '%s' to be a valid url: %v", val, err)
		}
		if u.Host == "" {
			return diag.Errorf("invalid URL: expected '%s' to have a host", val)
		}
		for _, s := range schemes {
			if u.Scheme == s {
				return nil
			}
		}
		return diag.Errorf("invalid URL: expected %s to have a url with schema of: %q, got %v", val, strings.Join(schemes, ","), v)
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

func stringLenBetween(min, max int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		v, ok := i.(string)
		if !ok {
			return diag.Errorf("expected type of %s to be string", k)
		}
		if len(v) < min || len(v) > max {
			return diag.Errorf("expected length of %s to be in the range (%d - %d), got %d", k, min, max, len(v))
		}
		return nil
	}
}

func stringAtLeast(min int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		v, ok := i.(string)
		if !ok {
			return diag.Errorf("expected type of %s to be string", k)
		}
		if len(strings.TrimSpace(v)) < min {
			return diag.Errorf("expected minimum length of %s to be %d, got %d", k, min, len(v))
		}
		return nil
	}
}

// regex lovingly lifted from: http://www.golangprograms.com/regular-expression-to-validate-email-address.html
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func stringIsEmail(i interface{}, k cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.Errorf("expected type of %s to be string", k)
	}
	if v == "" {
		return diag.Errorf("expected %s email to not be empty", k)
	}
	if !emailRegex.MatchString(v) {
		return diag.Errorf("%s field is not a valid email address", k)
	}
	return nil
}

var versionRegex = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)

func stringIsVersion(i interface{}, k cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.Errorf("expected type of %s to be string", k)
	}
	if v == "" {
		return diag.Errorf("expected %s version to not be empty", k)
	}
	if !versionRegex.MatchString(v) {
		return diag.Errorf("%s field is not a valid version", k)
	}
	return nil
}
