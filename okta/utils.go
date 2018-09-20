package okta

import (
	"fmt"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
)

// The best practices states that aggregate types should have error handling (think non-primitive). This will not attempt to set nil values.
func setNonPrimitives(data *schema.ResourceData, valueMap map[string]interface{}) error {
	for k, v := range valueMap {
		if v != nil {
			if err := data.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, data.Id(), err)
			}
		}
	}

	return nil
}

func convertStringArrToInterface(stringList []string) []interface{} {
	arr := make([]interface{}, len(stringList))
	for i, str := range stringList {
		arr[i] = str
	}
	return arr
}

func convertInterfaceToStringArr(purportedList interface{}) []string {
	var arr []string
	rawArr, ok := purportedList.([]interface{})

	if ok {
		arr = make([]string, len(rawArr))
		for i, thing := range rawArr {
			arr[i] = thing.(string)
		}
	}

	return arr
}

func convertIntToBool(i int) bool {
	if i > 0 {
		return true
	}

	return false
}

func convertBoolToInt(b bool) int {
	if b == true {
		return 1
	}
	return 0
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func suppressDefaultedArrayDiff(k, old, new string, d *schema.ResourceData) bool {
	return new == "0"
}

func createValueDiffSuppression(newValueToIgnore string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return new == newValueToIgnore
	}
}

func getClientFromMetadata(meta interface{}) *articulateOkta.Client {
	return meta.(*Config).articulateOktaClient
}

func validatePriority(in int, out int) error {
	if in > 0 && in != out {
		return fmt.Errorf("provided priority was not valid, got: %d, API responded with: %d. See schema for attribute details", in, out)
	}

	return nil
}

func is404(client *articulateOkta.Client) bool {
	return client.OktaErrorCode == "E0000007"
}
