package okta

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/okta/okta-sdk-golang/okta"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
)

// camel cased strings from okta responses become underscore separated to match
// the terraform configs for state file setting (ie. firstName from okta response becomes first_name)
func camelCaseToUnderscore(s string) string {
	a := []rune(s)

	for i, r := range a {
		if !unicode.IsLower(r) {
			a = append(a, 0)
			a[i] = unicode.ToLower(r)
			copy(a[i+1:], a[i:])
			a[i] = []rune("_")[0]
		}
	}

	s = string(a)

	return s
}

func conditionalRequire(d *schema.ResourceData, propList []string, reason string) error {
	var missing []string

	for _, prop := range propList {
		if _, ok := d.GetOkExists(prop); !ok {
			missing = append(missing, prop)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing conditionally required fields, reason: %s, missing fields: %s", reason, strings.Join(missing, ", "))
	}

	return nil
}

// Conditionally validates a slice of strings for required and valid values.
func conditionalValidator(field string, typeValue string, require []string, valid []string, actual []string) error {
	explanation := fmt.Sprintf("failed conditional validation for field \"%s\" of type \"%s\", it can contain %s", field, typeValue, strings.Join(valid, ", "))

	if len(require) > 0 {
		explanation = fmt.Sprintf("%s and must contain %s", explanation, strings.Join(require, ", "))
	}

	for _, val := range require {
		if !contains(actual, val) {
			return fmt.Errorf("%s, received %s", explanation, strings.Join(actual, ", "))
		}
	}

	for _, val := range actual {
		if !contains(valid, val) {
			return fmt.Errorf("%s, received %s", explanation, strings.Join(actual, ", "))
		}
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

// Ensures all elements are contained in slice.
func containsAll(s []string, elements ...string) bool {
	for _, a := range elements {
		if !contains(s, a) {
			return false
		}
	}

	return true
}

// Ensures at least one element is contained in provided slice. More performant version of contains(..) || contains(..)
func containsOne(s []string, elements ...string) bool {
	for _, a := range s {
		if contains(elements, a) {
			return true
		}
	}

	return false
}

func convertBoolToInt(b bool) int {
	if b == true {
		return 1
	}
	return 0
}

func convertIntToBool(i int) bool {
	if i > 0 {
		return true
	}

	return false
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

// Converts interface to string array, if there are no elements it returns nil to conform with optional properties.
func convertInterfaceToStringArrNullable(purportedList interface{}) []string {
	arr := convertInterfaceToStringArr(purportedList)

	if len(arr) < 1 {
		return nil
	}

	return arr
}

func convertStringArrToInterface(stringList []string) []interface{} {
	arr := make([]interface{}, len(stringList))
	for i, str := range stringList {
		arr[i] = str
	}
	return arr
}

// Allows you to chain multiple validation functions
func createValidationChain(validationChain ...schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	return func(val interface{}, key string) ([]string, []error) {
		var warningList []string
		var errorList []error

		for _, cb := range validationChain {
			warnings, errors := cb(val, key)
			errorList = append(errorList, errors...)
			warningList = append(warningList, warnings...)
		}

		return warningList, errorList
	}
}

func createValueDiffSuppression(newValueToIgnore string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return new == newValueToIgnore
	}
}

func ensureNotDefault(d *schema.ResourceData, t string) error {
	thing := fmt.Sprintf("Default %s", t)

	if d.Get("name").(string) == thing {
		return fmt.Errorf("%s is immutable", thing)
	}

	return nil
}

func getClientFromMetadata(meta interface{}) *articulateOkta.Client {
	return meta.(*Config).articulateOktaClient
}

func getOktaClientFromMetadata(meta interface{}) *okta.Client {
	return meta.(*Config).oktaClient
}

func is404(client *articulateOkta.Client) bool {
	return client.OktaErrorCode == "E0000007"
}

// regex lovingly lifted from: http://www.golangprograms.com/regular-expression-to-validate-email-address.html
func matchEmailRegexp(val interface{}, key string) (warnings []string, errors []error) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(val.(string)) == false {
		errors = append(errors, fmt.Errorf("%s field not a valid email address", key))
	}
	return warnings, errors
}

func requireOneOf(d *schema.ResourceData, propList ...string) error {
	for _, prop := range propList {
		if _, ok := d.GetOkExists(prop); !ok {
			return nil
		}
	}

	return fmt.Errorf("One of the following fields must be set: %s", strings.Join(propList, ", "))
}

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

func suppressDefaultedArrayDiff(k, old, new string, d *schema.ResourceData) bool {
	return new == "0"
}

func suppressDefaultedDiff(k, old, new string, d *schema.ResourceData) bool {
	return new == ""
}

// Matching level of validation done by Okta API
func validateIsURL(val interface{}, b string) ([]string, []error) {
	doesMatch, err := regexp.Match(`^(http|https):\/\/.*`, []byte(val.(string)))

	if err != nil {
		return nil, []error{err}
	} else if !doesMatch {
		return nil, []error{fmt.Errorf("failed to validate url, \"%s\"", val)}
	}

	return nil, nil
}

func validatePriority(in int, out int) error {
	if in > 0 && in != out {
		return fmt.Errorf("provided priority was not valid, got: %d, API responded with: %d. See schema for attribute details", in, out)
	}

	return nil
}
