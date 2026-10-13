package idaas

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestBuildAppGroupAssignmentPriority(t *testing.T) {
	for _, test := range []struct {
		name         string
		configured   cty.Value
		planned      int
		wantPriority bool
	}{
		{"zero", cty.NumberIntVal(0), 0, true},
		{"nonzero", cty.NumberIntVal(1), 1, true},
		{"omitted", cty.NullVal(cty.Number), 0, false},
		{"unknown", cty.UnknownVal(cty.Number), 0, false},
		{"planned_value", cty.NumberIntVal(0), 2, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			resource := resourceAppGroupAssignment()
			data, err := schema.InternalMap(resource.Schema).Data(nil, &terraform.InstanceDiff{
				RawConfig: cty.ObjectVal(map[string]cty.Value{
					"priority": test.configured,
				}),
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"priority": {New: strconv.Itoa(test.planned)},
					"profile":  {New: `{"customAttribute":"value"}`},
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			assignment := buildAppGroupAssignment(data)
			payload, err := json.Marshal(&assignment)
			if err != nil {
				t.Fatal(err)
			}
			var body map[string]json.RawMessage
			if err := json.Unmarshal(payload, &body); err != nil {
				t.Fatal(err)
			}
			priority, exists := body["priority"]
			if exists != test.wantPriority {
				t.Fatalf("priority present = %t, want %t; payload = %s", exists, test.wantPriority, payload)
			}
			if exists && string(priority) != strconv.Itoa(test.planned) {
				t.Fatalf("priority = %s, want %d", priority, test.planned)
			}
			if string(body["profile"]) != `{"customAttribute":"value"}` {
				t.Fatalf("unexpected profile: %s", body["profile"])
			}
		})
	}
}
