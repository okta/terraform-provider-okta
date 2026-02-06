package utils

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	v6okta "github.com/okta/okta-sdk-golang/v6/okta"

	"github.com/cenkalti/backoff"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v4/okta"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/sdk"
)

const DefaultPaginationLimit int64 = 200

var ErrInvalidElemFormat = errors.New("element type does not match the value provided in 'array_type' or 'type'")

func BuildSchema(schemas ...map[string]*schema.Schema) map[string]*schema.Schema {
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
func CamelCaseToUnderscore(s string) string {
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

// UnderscoreToCamelCase converts underscore separated strings to camel case
// (ie. first_name becomes firstName)
func UnderscoreToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			runes := []rune(parts[i])
			runes[0] = unicode.ToUpper(runes[0])
			parts[i] = string(runes)
		}
	}
	return strings.Join(parts, "")
}

func ConditionalRequire(d *schema.ResourceData, propList []string, reason string) error {
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
func ConditionalValidator(field, typeValue string, require, valid, actual []string) error {
	explanation := fmt.Sprintf("failed conditional validation for field \"%s\" of type \"%s\", it can contain %s", field, typeValue, strings.Join(valid, ", "))

	if len(require) > 0 {
		explanation = fmt.Sprintf("%s and must contain %s", explanation, strings.Join(require, ", "))
	}

	for _, val := range require {
		if !Contains(actual, val) {
			return fmt.Errorf("%s, received %s", explanation, strings.Join(actual, ", "))
		}
	}

	for _, val := range actual {
		if !Contains(valid, val) {
			return fmt.Errorf("%s, received %s", explanation, strings.Join(actual, ", "))
		}
	}

	return nil
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsInt(codes []int, code int) bool {
	for _, a := range codes {
		if a == code {
			return true
		}
	}
	return false
}

// Ensures at least one element is contained in provided slice. More performant version of contains(..) || contains(..)
func ContainsOne(s []string, elements ...string) bool {
	for _, a := range s {
		if Contains(elements, a) {
			return true
		}
	}
	return false
}

func ConvertInterfaceToStringSet(purportedSet interface{}) []string {
	return ConvertInterfaceToStringArr(purportedSet.(*schema.Set).List())
}

func ConvertInterfaceToStringSetNullable(purportedSet interface{}) []string {
	set, ok := purportedSet.(*schema.Set)
	if ok {
		return ConvertInterfaceToStringArrNullable(set.List())
	}
	return nil
}

func ConvertInterfaceToStringArr(purportedList interface{}) []string {
	var arr []string
	rawArr, ok := purportedList.([]interface{})
	if ok {
		arr = ConvertInterfaceArrToStringArr(rawArr)
	}
	return arr
}

func ConvertInterfaceArrToStringArr(rawArr []interface{}) []string {
	arr := make([]string, len(rawArr))
	for i, thing := range rawArr {
		if a, ok := thing.(string); ok {
			arr[i] = a
		}
	}
	return arr
}

// Converts interface to string array, if there are no elements it returns nil to conform with optional properties.
func ConvertInterfaceToStringArrNullable(purportedList interface{}) []string {
	arr := ConvertInterfaceToStringArr(purportedList)
	if len(arr) < 1 {
		return nil
	}
	return arr
}

func CreateNestedResourceImporter(fields []string) *schema.ResourceImporter {
	return CreateCustomNestedResourceImporter(fields, fmt.Sprintf("Expecting the following format %s", strings.Join(fields, "/")))
}

// CreateCustomNestedResourceImporter Fields making up the ID should be in
// order, for instance, []string{"auth_server_id", "policy_id", "id"} However,
// extra fields can be specified after as well, []string{"auth_server_id",
// "policy_id", "id", "extra"}
func CreateCustomNestedResourceImporter(fields []string, errMessage string) *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
			parts := strings.Split(d.Id(), "/")
			if len(parts) != len(fields) {
				return nil, fmt.Errorf("expected %d import fields %q, got %d fields %q", len(fields), strings.Join(fields, "/"), len(parts), d.Id())
			}
			for i, field := range fields {
				if field == "id" {
					d.SetId(parts[i])
					continue
				}
				var value interface{}
				if i < len(parts) {
					// deal with the import parameter being a boolean "true" / "false"
					if bValue, err := strconv.ParseBool(parts[i]); err == nil {
						value = bValue
					} else {
						value = parts[i]
					}
				}
				// lintignore:R001
				_ = d.Set(field, value)
			}

			return []*schema.ResourceData{d}, nil
		},
	}
}

func ConvertStringSliceToInterfaceSlice(stringList []string) []interface{} {
	if len(stringList) == 0 {
		return nil
	}
	arr := make([]interface{}, len(stringList))
	for i, str := range stringList {
		arr[i] = str
	}
	return arr
}

func ConvertStringSliceToSet(stringList []string) *schema.Set {
	arr := make([]interface{}, len(stringList))
	for i, str := range stringList {
		arr[i] = str
	}
	return schema.NewSet(schema.HashString, arr)
}

func ConvertStringSliceToSetNullable(stringList []string) *schema.Set {
	if len(stringList) == 0 {
		return nil
	}
	arr := make([]interface{}, len(stringList))
	for i, str := range stringList {
		arr[i] = str
	}
	return schema.NewSet(schema.HashString, arr)
}

func CreateValueDiffSuppression(newValueToIgnore string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return new == newValueToIgnore
	}
}

func EnsureNotDefault(d *schema.ResourceData, t string) error {
	thing := fmt.Sprintf("Default %s", t)

	if d.Get("name").(string) == thing {
		return fmt.Errorf("%s is immutable", thing)
	}

	return nil
}

func GetMapString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if res, ok := v.(string); ok {
			return res
		}
	}
	return ""
}

// BoolPtr return bool pointer to b's value
func BoolPtr(b bool) (ptr *bool) {
	ptr = &b
	return
}

// BoolFromBoolPtr if b is nil returns false, otherwise return boolean value of b
func BoolFromBoolPtr(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func StringPtr(s string) (ptr *string) {
	ptr = &s
	return
}

func DoesResourceExist(response *sdk.Response, err error) (bool, error) {
	if response == nil {
		return false, err
	}
	// We don't want to consider a 404 an error in some cases and thus the delineation
	if response.StatusCode == 404 {
		return false, nil
	}
	if err != nil {
		return false, ResponseErr(response, err)
	}

	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return false, ResponseErr(response, err)
	}
	// some of the API response can be 200 and return an empty object or list meaning nothing was found
	body := string(b)
	if body == "{}" || body == "[]" {
		return false, nil
	}

	return true, nil
}

func DoesResourceExistV3(response *okta.APIResponse, err error) (bool, error) {
	if response == nil {
		return false, err
	}
	// We don't want to consider a 404 an error in some cases and thus the delineation
	if response.StatusCode == 404 {
		return false, nil
	}
	if err != nil {
		return false, ResponseErr_V3(response, err)
	}

	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return false, ResponseErr_V3(response, err)
	}
	// some of the API response can be 200 and return an empty object or list meaning nothing was found
	body := string(b)
	if body == "{}" || body == "[]" {
		return false, nil
	}

	return true, nil
}

func DoesResourceExistV5(response *v5okta.APIResponse, err error) (bool, error) {
	if response == nil {
		return false, err
	}
	// We don't want to consider a 404 an error in some cases and thus the delineation
	if response.StatusCode == 404 {
		return false, nil
	}
	if err != nil {
		return false, ResponseErr_V5(response, err)
	}

	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return false, ResponseErr_V5(response, err)
	}
	// some of the API response can be 200 and return an empty object or list meaning nothing was found
	body := string(b)
	if body == "{}" || body == "[]" {
		return false, nil
	}

	return true, nil
}

func DoesResourceExistV6(response *v6okta.APIResponse, err error) (bool, error) {
	if response == nil {
		return false, err
	}
	// We don't want to consider a 404 an error in some cases and thus the delineation
	if response.StatusCode == 404 {
		return false, nil
	}
	if err != nil {
		return false, ResponseErr_V6(response, err)
	}

	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return false, ResponseErr_V6(response, err)
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
func SuppressErrorOn404(resp *sdk.Response, err error) error {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return nil
	}
	return ResponseErr(resp, err)
}

// TODO switch to suppressErrorOn404 when migration complete
func SuppressErrorOn404_V3(resp *okta.APIResponse, err error) error {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return nil
	}
	return ResponseErr_V3(resp, err)
}

// TODO switch to suppressErrorOn404 when migration complete
func SuppressErrorOn404_V5(resp *v5okta.APIResponse, err error) error {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return nil
	}
	return ResponseErr_V5(resp, err)
}

func SuppressErrorOn404_V6(resp *v6okta.APIResponse, err error) error {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return nil
	}
	return ResponseErr_V6(resp, err)
}

// Useful shortcut for suppressing errors from Okta's SDK when a Org does not
// have permission to access a feature.
func SuppressErrorOn401(what string, meta interface{}, resp *sdk.Response, err error) error {
	if resp != nil && resp.StatusCode == http.StatusUnauthorized {
		return nil
	}
	return ResponseErr(resp, err)
}

// Useful shortcut for suppressing errors from Okta's SDK when a Org does not
// have permission to access a feature.
func SuppressErrorOn403(what string, meta interface{}, resp *sdk.Response, err error) error {
	if resp != nil && resp.StatusCode == http.StatusForbidden {
		return nil
	}
	return ResponseErr(resp, err)
}

func Is404(resp *sdk.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusNotFound
}

func NormalizeDataJSON(val interface{}) string {
	dataMap := map[string]interface{}{}

	// Ignoring errors since we know it is valid
	_ = json.Unmarshal([]byte(val.(string)), &dataMap)
	ret, _ := json.Marshal(dataMap)

	return string(ret)
}

// Removes nulls from group profile map and returns, since Okta does not render nulls in profile
func NormalizeGroupProfile(profile sdk.GroupProfileMap) sdk.GroupProfileMap {
	trimedProfile := make(sdk.GroupProfileMap)
	for k, v := range profile {
		if v != nil {
			trimedProfile[k] = v
		}
	}
	return trimedProfile
}

// Opposite of append
func Remove(arr []string, el string) []string {
	var newArr []string

	for _, item := range arr {
		if item != el {
			newArr = append(newArr, item)
		}
	}
	return newArr
}

// AppendUnique appends el to arr if el isn't already present in arr
func AppendUnique(arr []string, el string) []string {
	found := false
	for _, item := range arr {
		if item == el {
			found = true
			break
		}
	}
	if found {
		return arr
	}
	return append(arr, el)
}

// The best practices states that aggregate types should have error handling (think non-primitive). This will not attempt to set nil values.
func SetNonPrimitives(d *schema.ResourceData, valueMap map[string]interface{}) error {
	for k, v := range valueMap {
		if v != nil {
			// lintignore:R001
			if err := d.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return nil
}

// Okta SDK will (not often) return just `Okta API has returned an error: ""“ when the error is not valid JSON.
// The status should help with debugability. Potentially also could check for an empty error and omit
// it when it occurs and build some more context.
func ResponseErr(resp *sdk.Response, err error) error {
	if err != nil {
		msg := err.Error()
		if resp != nil {
			msg += fmt.Sprintf(", Status: %s", resp.Status)
		}
		return errors.New(msg)
	}
	return nil
}

func ResponseErr_V3(resp *okta.APIResponse, err error) error {
	if err != nil {
		msg := err.Error()
		if resp != nil {
			msg += fmt.Sprintf(", Status: %s", resp.Status)
		}
		return errors.New(msg)
	}
	return nil
}

// TODO switch to responseErr when migration complete
func ResponseErr_V5(resp *v5okta.APIResponse, err error) error {
	if err != nil {
		msg := err.Error()
		if resp != nil {
			msg += fmt.Sprintf(", Status: %s", resp.Status)
		}
		return errors.New(msg)
	}
	return nil
}

func ResponseErr_V6(resp *v6okta.APIResponse, err error) error {
	if err != nil {
		msg := err.Error()
		if resp != nil {
			msg += fmt.Sprintf(", Status: %s", resp.Status)
		}
		return errors.New(msg)
	}
	return nil
}

func ValidatePriority(in, out int64) error {
	if in > 0 && in != out {
		return fmt.Errorf("provided priority was not valid, got: %d, API responded with: %d. See schema for attribute details", in, out)
	}
	return nil
}

func BuildEnum(ae []interface{}, elemType string) ([]interface{}, error) {
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
				return nil, ErrInvalidElemFormat
			}
			enum[i] = f
		case "integer":
			f, err := strconv.Atoi(aeStr)
			if err != nil {
				return nil, ErrInvalidElemFormat
			}
			enum[i] = f
		default:
			enum[i] = aeStr
		}
	}
	return enum, nil
}

// LocalFileStateFunc - helper for schema.Schema StateFunc checking if a the
// blob of a local file has changed - is not file path dependant.
func LocalFileStateFunc(val interface{}) string {
	filePath := val.(string)
	if filePath == "" {
		return ""
	}
	return ComputeFileHash(filePath)
}

// ComputeFileHash - equivalent to  `shasum -a 256 filepath`
func ComputeFileHash(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer file.Close()
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// SuppressDuringCreateFunc - attribute has changed assume this is a create and
// treat the properties as readers not caring about what would otherwise apear
// to be drift.
func SuppressDuringCreateFunc(attribute string) func(k, old, new string, d *schema.ResourceData) bool {
	return func(k, old, new string, d *schema.ResourceData) bool {
		if d.HasChange(attribute) {
			return true
		}
		return old == new
	}
}

// Normalizes to certificate object when it's passed as a raw b64 block instead of a full pem file
func RawCertNormalize(certContents string) (*x509.Certificate, error) {
	certContents = strings.ReplaceAll(strings.TrimSpace(certContents), " ", "")
	certDecoded, err := base64.StdEncoding.DecodeString(certContents)
	if err != nil {
		return nil, fmt.Errorf("failed to decode b64: %s", err)
	}
	cert, err := x509.ParseCertificate(certDecoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode pem certificate: %s", err)
	}

	return cert, nil
}

// Normalizes to certificate object when passed as PEM formatted certificate file
func PemCertNormalize(certContents string) (*x509.Certificate, error) {
	certContents = strings.TrimSpace(certContents)
	cert, rest := pem.Decode([]byte(certContents))
	if cert == nil {
		return nil, fmt.Errorf("failed to decode PEM file, rest: %s", rest)
	}

	parsedCert, err := x509.ParseCertificate(cert.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %s", err)
	}

	return parsedCert, nil
}

func CertNormalize(certContents string) (*x509.Certificate, error) {
	certDecoded, err := PemCertNormalize(certContents)
	if err == nil {
		return certDecoded, nil
	}
	certDecoded, err = RawCertNormalize(certContents)
	if err != nil {
		return nil, err
	}
	return certDecoded, nil
}

// NoChangeInObjectFromUnmarshaledJSON Intended for use by a DiffSuppressFunc,
// returns true if old and new JSONs are equivalent object representations ...
// It is true, there is no change!  Edge chase if newJSON is blank, will also
// return true which cover the new resource case.
func NoChangeInObjectFromUnmarshaledJSON(k, oldJSON, newJSON string, d *schema.ResourceData) bool {
	if newJSON == "" {
		return true
	}
	var oldObj map[string]any
	var newObj map[string]any
	if err := json.Unmarshal([]byte(oldJSON), &oldObj); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(newJSON), &newObj); err != nil {
		return false
	}

	return reflect.DeepEqual(oldObj, newObj)
}

// NoChangeInObjectWithSortedSlicesFromUnmarshaledJSON Intended for use by a DiffSuppressFunc,
// returns true if old and new JSONs are equivalent object representations no matter the order of any slices...
// It is true, there is no change!  Edge chase if newJSON is blank, will also
// return true which cover the new resource case.
func NoChangeInObjectWithSortedSlicesFromUnmarshaledJSON(k, oldJSON, newJSON string, d *schema.ResourceData) bool {
	if newJSON == "" {
		return true
	}

	var oldObj any
	var newObj any

	if err := json.Unmarshal([]byte(oldJSON), &oldObj); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(newJSON), &newObj); err != nil {
		return false
	}

	oldObj = sortSlices(oldObj)
	newObj = sortSlices(newObj)

	return reflect.DeepEqual(oldObj, newObj)
}

func sortSlices(obj any) any {
	switch v := obj.(type) {
	case []any:
		for i := range v {
			v[i] = sortSlices(v[i])
		}
		sort.SliceStable(v, func(i, j int) bool {
			return less(v[i], v[j])
		})
		return v
	case map[string]any:
		for key, val := range v {
			v[key] = sortSlices(val)
		}
		return v
	default:
		return v
	}
}

func less(a, b any) bool {
	// Unmarshaled JSON into any can only have string, float64, bool, nil, []any, map[string]any
	switch aTyped := a.(type) {
	case string:
		if bTyped, ok := b.(string); ok {
			return aTyped < bTyped
		}
	case float64:
		if bTyped, ok := b.(float64); ok {
			return aTyped < bTyped
		}
	case bool:
		if bTyped, ok := b.(bool); ok {
			return !aTyped && bTyped
		}
	}

	if a == nil {
		return true
	}
	if b == nil {
		return false
	}

	// Fallback: use type name as last resort to ensure consistency
	return reflect.TypeOf(a).String() < reflect.TypeOf(b).String()
}

func Intersection(old, new []string) (intersection, exclusiveOld, exclusiveNew []string) {
	intersection = make([]string, 0)
	exclusiveOld = make([]string, 0)
	exclusiveNew = make([]string, 0)
	oldElementMap := make(map[string]bool)
	newElementMap := make(map[string]bool)
	for _, o := range old {
		oldElementMap[o] = true
	}
	for _, n := range new {
		newElementMap[n] = true
	}
	for _, n := range new {
		if oldElementMap[n] {
			intersection = append(intersection, n)
		} else {
			exclusiveNew = append(exclusiveNew, n)
		}
	}
	for _, o := range old {
		if !newElementMap[o] {
			exclusiveOld = append(exclusiveOld, o)
		}
	}
	return
}

func LogoFileIsValid() schema.SchemaValidateDiagFunc {
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

// NewExponentialBackOffWithContext helper to dry up creating a backoff object that is exponential and has context
func NewExponentialBackOffWithContext(ctx context.Context, maxElapsedTime time.Duration) backoff.BackOffContext {
	bOff := backoff.NewExponentialBackOff()
	bOff.MaxElapsedTime = maxElapsedTime

	// NOTE: backoff.BackOffContext is an interface that embeds backoff.Backoff
	// so the greater context is considered on backoff.Retry
	return backoff.WithContext(bOff, ctx)
}

func Int64Ptr(what int) *int64 {
	result := int64(what)
	return &result
}

func ResourceFuncNoOp(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}

func LinksValue(links interface{}, keys ...string) string {
	if links == nil {
		return ""
	}
	sl, ok := links.([]interface{})
	if ok {
		if len(sl) == 0 {
			links = map[string]interface{}{}
		} else {
			links = sl[0]
		}
	}
	if len(keys) == 0 {
		v, ok := links.(string)
		if !ok {
			return ""
		}
		return v
	}
	l, ok := links.(map[string]interface{})
	if !ok {
		return ""
	}
	if len(keys) == 1 {
		return LinksValue(l[keys[0]])
	}
	return LinksValue(l[keys[0]], keys[1:]...)
}

// StrMaxLength validates that the string is not longer than the specified maximum length.
func StrMaxLength(max int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		v, ok := i.(string)
		if !ok {
			return diag.Errorf("expected type of %s to be string", k)
		}

		// https://github.com/okta/terraform-provider-okta/issues/2396
		runes := []rune(v)
		if len(runes) > max {
			return diag.Errorf("%s cannot be longer than %d runes", k, max)
		}
		return nil
	}
}
