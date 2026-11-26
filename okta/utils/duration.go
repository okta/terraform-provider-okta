package utils

import (
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// iso8601Duration matches PnYnMnWnDTnHnMnS (integer components, seconds may be fractional).
var iso8601Duration = regexp.MustCompile(`^P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)W)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+(?:\.\d+)?)S)?)?$`)

type isoDuration struct {
	years, months int64
	fixed         time.Duration // weeks, days, hours, minutes, seconds
}

func parseISO8601Duration(s string) (isoDuration, bool) {
	m := iso8601Duration.FindStringSubmatch(s)
	if m == nil || s == "P" || s == "PT" {
		return isoDuration{}, false
	}
	n := func(i int) int64 {
		if m[i] == "" {
			return 0
		}
		v, _ := strconv.ParseInt(m[i], 10, 64)
		return v
	}
	var d isoDuration
	d.years, d.months = n(1), n(2)
	d.fixed = time.Duration(n(3))*7*24*time.Hour +
		time.Duration(n(4))*24*time.Hour +
		time.Duration(n(5))*time.Hour +
		time.Duration(n(6))*time.Minute
	if m[7] != "" {
		secs, _ := strconv.ParseFloat(m[7], 64)
		d.fixed += time.Duration(secs * float64(time.Second))
	}
	return d, true
}

// SameISO8601Duration reports whether a and b are the same ISO 8601 duration written
// differently, e.g. "P7D" and "PT168H" (Okta normalizes durations to hours). Years and
// months have no fixed length, so they only match the same number of years and months.
// Strings that are not durations are compared as-is.
func SameISO8601Duration(a, b string) bool {
	if a == b {
		return true
	}
	da, okA := parseISO8601Duration(a)
	db, okB := parseISO8601Duration(b)
	return okA && okB && da == db
}

// PreferPriorDuration returns prior when the API value api is the same duration, so a
// configured "P7D" is not replaced by Okta's "PT168H" (no spurious diff, no "inconsistent
// result after apply"); otherwise it returns the API value.
func PreferPriorDuration(prior types.String, api string) types.String {
	if !prior.IsNull() && !prior.IsUnknown() && SameISO8601Duration(prior.ValueString(), api) {
		return prior
	}
	return types.StringValue(api)
}
