package okta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/terraform-provider-okta/sdk"
)

const defaultPaginationLimit int64 = 200

func buildSchema(schemas ...map[string]*schema.Schema) map[string]*schema.Schema {
	r := map[string]*schema.Schema{}
	for _, s := range schemas {
		for key, val := range s {
			r[key] = val
		}
	}
	return r
}

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
		if _, ok := d.GetOk(prop); !ok {
			missing = append(missing, prop)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing conditionally required fields, reason: '%s', missing fields: %s", reason, strings.Join(missing, ", "))
	}

	return nil
}

// Conditionally validates a slice of strings for required and valid values.
func conditionalValidator(field, typeValue string, require, valid, actual []string) error {
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

func containsInt(codes []int, code int) bool {
	for _, a := range codes {
		if a == code {
			return true
		}
	}
	return false
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

func convertInterfaceToStringSet(purportedSet interface{}) []string {
	return convertInterfaceToStringArr(purportedSet.(*schema.Set).List())
}

func convertInterfaceToStringSetNullable(purportedSet interface{}) []string {
	set, ok := purportedSet.(*schema.Set)
	if ok {
		return convertInterfaceToStringArrNullable(set.List())
	}
	return nil
}

func convertInterfaceToStringArr(purportedList interface{}) []string {
	var arr []string
	rawArr, ok := purportedList.([]interface{})
	if ok {
		arr = convertInterfaceArrToStringArr(rawArr)
	}
	return arr
}

func convertInterfaceArrToStringArr(rawArr []interface{}) []string {
	arr := make([]string, len(rawArr))
	for i, thing := range rawArr {
		arr[i] = thing.(string)
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

func createNestedResourceImporter(fields []string) *schema.ResourceImporter {
	return createCustomNestedResourceImporter(fields, fmt.Sprintf("Expecting the following format %s", strings.Join(fields, "/")))
}

// Fields should be in order, for instance, []string{"auth_server_id", "policy_id", "id"}
func createCustomNestedResourceImporter(fields []string, errMessage string) *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
			parts := strings.Split(d.Id(), "/")
			if len(parts) != len(fields) {
				return nil, fmt.Errorf("invalid resource import specifier. %s", errMessage)
			}

			for i, field := range fields {
				if field == "id" {
					d.SetId(parts[i])
					continue
				}
				_ = d.Set(field, parts[i])
			}

			return []*schema.ResourceData{d}, nil
		},
	}
}

func convertStringSliceToInterfaceSlice(stringList []string) []interface{} {
	if len(stringList) == 0 {
		return nil
	}
	arr := make([]interface{}, len(stringList))
	for i, str := range stringList {
		arr[i] = str
	}
	return arr
}

func convertStringSliceToSet(stringList []string) *schema.Set {
	arr := make([]interface{}, len(stringList))
	for i, str := range stringList {
		arr[i] = str
	}
	return schema.NewSet(schema.HashString, arr)
}

func convertStringSliceToSetNullable(stringList []string) *schema.Set {
	if len(stringList) == 0 {
		return nil
	}
	arr := make([]interface{}, len(stringList))
	for i, str := range stringList {
		arr[i] = str
	}
	return schema.NewSet(schema.HashString, arr)
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

func getMapString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		return v.(string)
	}
	return ""
}

func boolPtr(b bool) (ptr *bool) {
	ptr = &b
	return
}

func stringPtr(s string) (ptr *string) {
	ptr = &s
	return
}

func doesResourceExist(response *okta.Response, err error) (bool, error) {
	if response == nil {
		return false, err
	}
	// We don't want to consider a 404 an error in some cases and thus the delineation
	if response.StatusCode == 404 {
		return false, nil
	}
	if err != nil {
		return false, responseErr(response, err)
	}

	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return false, responseErr(response, err)
	}
	// some of the API response can be 200 and return an empty object or list meaning nothing was found
	body := string(b)
	if body == "{}" || body == "[]" {
		return false, nil
	}

	return true, nil
}

// Useful shortcut for suppressing errors from Okta's SDK when a resource does not exist. Usually used during deletion
// of nested resources.
func suppressErrorOn404(resp *okta.Response, err error) error {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return nil
	}
	return responseErr(resp, err)
}

func getParallelismFromMetadata(meta interface{}) int {
	return meta.(*Config).parallelism
}

func getOktaClientFromMetadata(meta interface{}) *okta.Client {
	return meta.(*Config).oktaClient
}

func getSupplementFromMetadata(meta interface{}) *sdk.APISupplement {
	return meta.(*Config).supplementClient
}

func getRequestExecutor(m interface{}) *okta.RequestExecutor {
	return getOktaClientFromMetadata(m).GetRequestExecutor()
}

func is404(resp *okta.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusNotFound
}

func logger(meta interface{}) hclog.Logger {
	return meta.(*Config).logger
}

func normalizeDataJSON(val interface{}) string {
	dataMap := map[string]interface{}{}

	// Ignoring errors since we know it is valid
	_ = json.Unmarshal([]byte(val.(string)), &dataMap)
	ret, _ := json.Marshal(dataMap)

	return string(ret)
}

// Opposite of append
func remove(arr []string, el string) []string {
	var newArr []string

	for _, item := range arr {
		if item != el {
			newArr = append(newArr, item)
		}
	}
	return newArr
}

// The best practices states that aggregate types should have error handling (think non-primitive). This will not attempt to set nil values.
func setNonPrimitives(d *schema.ResourceData, valueMap map[string]interface{}) error {
	for k, v := range valueMap {
		if v != nil {
			if err := d.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return nil
}

// Okta SDK will (not often) return just `Okta API has returned an error: ""`` when the error is not valid JSON.
// The status should help with debugability. Potentially also could check for an empty error and omit
// it when it occurs and build some more context.
func responseErr(resp *okta.Response, err error) error {
	if err != nil {
		msg := err.Error()
		if resp != nil {
			msg += fmt.Sprintf(", Status: %s", resp.Status)
		}
		return errors.New(msg)
	}
	return nil
}

func validatePriority(in, out int64) error {
	if in > 0 && in != out {
		return fmt.Errorf("provided priority was not valid, got: %d, API responded with: %d. See schema for attribute details", in, out)
	}
	return nil
}

func buildEnum(ae []interface{}, elemType string) ([]interface{}, error) {
	enum := make([]interface{}, len(ae))
	for i, aeVal := range ae {
		if aeVal == nil {
			switch elemType {
			case "number":
				enum[i] = float64(0)
			case "integer":
				enum[i] = 0
			default:
				enum[i] = ""
			}
			continue
		}

		aeStr, ok := aeVal.(string)
		if !ok {
			return nil, fmt.Errorf("expected %+v value to cast to string", aeVal)
		}
		switch elemType {
		case "number":
			f, err := strconv.ParseFloat(aeStr, 64)
			if err != nil {
				return nil, errInvalidElemFormat
			}
			enum[i] = f
		case "integer":
			f, err := strconv.Atoi(aeStr)
			if err != nil {
				return nil, errInvalidElemFormat
			}
			enum[i] = f
		default:
			enum[i] = aeStr
		}
	}
	return enum, nil
}
