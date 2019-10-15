package okta

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/peterhellberg/link"
	sdk "github.com/terraform-providers/terraform-provider-okta/sdk"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
)

func buildSchema(s, t map[string]*schema.Schema) map[string]*schema.Schema {
	for key, val := range s {
		t[key] = val
	}

	return t
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

func condenseError(errorList []error) error {
	if len(errorList) < 1 {
		return nil
	}
	msgList := make([]string, len(errorList))
	for i, err := range errorList {
		msgList[i] = err.Error()
	}

	return fmt.Errorf("Series of errors occurred: %s", strings.Join(msgList, ", "))
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

func convertInterfaceToStringSet(purportedSet interface{}) []string {
	return convertInterfaceToStringArr(purportedSet.(*schema.Set).List())
}

func convertInterfaceToStringSetNullable(purportedSet interface{}) []string {
	return convertInterfaceToStringArrNullable(purportedSet.(*schema.Set).List())
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

func createNestedResourceImporter(fields []string) *schema.ResourceImporter {
	return createCustomNestedResourceImporter(fields, fmt.Sprintf("Expecting the following format %s", strings.Join(fields, "/")))
}

// Fields should be in order, for instance, []string{"auth_server_id", "policy_id", "id"}
func createCustomNestedResourceImporter(fields []string, errMessage string) *schema.ResourceImporter {
	return &schema.ResourceImporter{
		State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			parts := strings.Split(d.Id(), "/")
			if len(parts) != len(fields) {
				return nil, fmt.Errorf("Invalid resource import specifier. %s", errMessage)
			}

			for i, field := range fields {
				if field == "id" {
					d.SetId(parts[i])
					continue
				}
				d.Set(field, parts[i])
			}

			return []*schema.ResourceData{d}, nil
		},
	}
}

func convertStringArrToInterface(stringList []string) []interface{} {
	arr := make([]interface{}, len(stringList))
	for i, str := range stringList {
		arr[i] = str
	}
	return arr
}

func convertStringSetToInterface(stringList []string) *schema.Set {
	arr := make([]interface{}, len(stringList))
	for i, str := range stringList {
		arr[i] = str
	}
	return schema.NewSet(schema.HashString, arr)
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

// Grabs after link from link headers if it exists
func getAfterParam(res *okta.Response) string {
	if res == nil {
		return ""
	}

	linkList := link.ParseHeader(res.Header)
	for _, l := range linkList {
		if l.Rel == "next" {
			parsedURL, err := url.Parse(l.URI)
			if err != nil {
				continue
			}
			q := parsedURL.Query()
			return q.Get("after")
		}
	}

	return ""
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

func doesResourceExist(response *okta.Response, err error) (bool, error) {
	// We don't want to consider a 404 an error in some cases and thus the delineation
	if response.StatusCode == 404 {
		return false, nil
	}

	if err != nil {
		return false, responseErr(response, err)
	}

	return true, err
}

func intPtr(b int) (ptr *int) {
	ptr = &b
	return
}

func getNullableInt(d *schema.ResourceData, key string) *int {
	if v, ok := d.GetOk(key); ok {
		i := v.(int)

		return &i
	}

	return nil
}

// Useful shortcut for suppressing errors from Okta's SDK when a resource does not exist. Usually used during deletion
// of nested resources.
func suppressErrorOn404(resp *okta.Response, err error) error {
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	return responseErr(resp, err)
}

func getApiToken(m interface{}) string {
	return m.(*Config).apiToken
}

func getBaseUrl(m interface{}) string {
	c := m.(*Config)
	return fmt.Sprintf("https://%v.%v", c.orgName, c.domain)
}

// Safely get string value
func getStringValue(d *schema.ResourceData, key string) string {
	if v, ok := d.GetOk(key); ok {
		return v.(string)
	}
	return ""
}

func getParallelismFromMetadata(meta interface{}) int {
	return meta.(*Config).parallelism
}

func getClientFromMetadata(meta interface{}) *articulateOkta.Client {
	return meta.(*Config).articulateOktaClient
}

func getOktaClientFromMetadata(meta interface{}) *okta.Client {
	return meta.(*Config).oktaClient
}

func getSupplementFromMetadata(meta interface{}) *sdk.ApiSupplement {
	return meta.(*Config).supplementClient
}

func getRequestExecutor(m interface{}) *okta.RequestExecutor {
	return getOktaClientFromMetadata(m).GetRequestExecutor()
}

func is404(status int) bool {
	return status == http.StatusNotFound
}

// regex lovingly lifted from: http://www.golangprograms.com/regular-expression-to-validate-email-address.html
func matchEmailRegexp(val interface{}, key string) (warnings []string, errors []error) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(val.(string)) == false {
		errors = append(errors, fmt.Errorf("%s field not a valid email address", key))
	}
	return warnings, errors
}

func mergeMaps(target, source map[string]interface{}) map[string]interface{} {
	for key, value := range source {
		target[key] = value
	}

	return target
}

func normalizeDataJSON(val interface{}) string {
	dataMap := map[string]interface{}{}

	// Ignoring errors since we know it is valid
	json.Unmarshal([]byte(val.(string)), &dataMap)
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

func requireOneOf(d *schema.ResourceData, propList ...string) error {
	for _, prop := range propList {
		if _, ok := d.GetOkExists(prop); !ok {
			return nil
		}
	}

	return fmt.Errorf("One of the following fields must be set: %s", strings.Join(propList, ", "))
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

func validateDataJSON(val interface{}, k string) ([]string, []error) {
	err := json.Unmarshal([]byte(val.(string)), &map[string]interface{}{})
	if err != nil {
		return nil, []error{err}
	}
	return nil, nil
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
