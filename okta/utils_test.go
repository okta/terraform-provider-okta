package okta

import (
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestRemove(t *testing.T) {
	tests := []struct {
		elements []string
		toRemove string
		expected []string
	}{
		{[]string{"one", "two", "three"}, "dne", []string{"one", "two", "three"}},
		{[]string{"one", "two", "three"}, "", []string{"one", "two", "three"}},
		{[]string{"one", "two", "three"}, "one", []string{"two", "three"}},
		{[]string{"one", "two", "three"}, "two", []string{"one", "three"}},
		{[]string{"one", "two", "three"}, "three", []string{"one", "two"}},
	}

	for _, test := range tests {
		result := remove(test.elements, test.toRemove)
		require.Equal(t, test.expected, result)
	}
}

func TestAppendUnique(t *testing.T) {
	tests := []struct {
		elements []string
		toAdd    string
		expected []string
	}{
		{[]string{"one", "two"}, "one", []string{"one", "two"}},
		{[]string{"one", "two"}, "three", []string{"one", "two", "three"}},
	}

	for _, test := range tests {
		result := appendUnique(test.elements, test.toAdd)
		require.Equal(t, test.expected, result)
	}
}

func TestContainsOne(t *testing.T) {
	testArr := []string{"1", "2", "3"}
	tests := []struct {
		elements []string
		expected bool
	}{
		{[]string{"2", "3"}, true},
		{[]string{"1", "2"}, true},
		{[]string{"1", "4"}, true},
		{[]string{"1"}, true},
		{[]string{"10"}, false},
	}

	for _, test := range tests {
		actual := containsOne(testArr, test.elements...)

		if actual != test.expected {
			t.Errorf("containsOne test failed, test array: %s, elements %s, expected %t, actual %t", strings.Join(testArr, ", "), strings.Join(test.elements, ", "), test.expected, actual)
		}
	}
}

func TestContains(t *testing.T) {
	testArr := []string{"1", "2", "3"}
	tests := []struct {
		element  string
		expected bool
	}{
		{"3", true},
		{"1", true},
		{"4", false},
		{"10", false},
		{"", false},
	}

	for _, test := range tests {
		actual := containsOne(testArr, test.element)

		if actual != test.expected {
			t.Errorf("contains test failed, test array: %s, elements %s, expected %t, actual %t", strings.Join(testArr, ", "), test.element, test.expected, actual)
		}
	}
}

func TestBuildSchema(t *testing.T) {
	sampleSchema1, sampleSchema2 := new(schema.Schema), new(schema.Schema)
	testMaps := []map[string]*schema.Schema{
		{
			"A": nil,
			"B": nil,
			"C": sampleSchema2,
		},
		{
			"A": sampleSchema1,
			"D": sampleSchema1,
			"E": sampleSchema1,
		},
		{
			"C": nil,
			"D": sampleSchema2,
		},
	}

	tests := []struct {
		element  string
		expected *schema.Schema
	}{
		{"A", sampleSchema1},
		{"B", nil},
		{"C", nil},
		{"D", sampleSchema2},
		{"E", sampleSchema1},
	}

	output := buildSchema(testMaps...)
	for _, test := range tests {
		actual := output[test.element]
		if actual != test.expected {
			t.Errorf("buildSchema test failed, test maps: %v, elements %s, expected %v, actual %v",
				testMaps, test.element, test.expected, actual)
		}
	}
}

func TestBuildEnum(t *testing.T) {
	tests := []struct {
		name        string
		ae          []interface{}
		elemType    string
		expected    []interface{}
		shouldError bool
	}{
		{
			"number slice including empty value",
			[]interface{}{"1.1", nil},
			"number",
			[]interface{}{1.1, 0.0},
			false,
		},
		{
			"integer slice including empty value",
			[]interface{}{"1", nil},
			"integer",
			[]interface{}{1, 0},
			false,
		},
		{
			"string slice including empty value",
			[]interface{}{"one", nil},
			"string",
			[]interface{}{"one", ""},
			false,
		},
		{
			"number slice with invalid value and empty value",
			[]interface{}{"one", nil},
			"number",
			nil,
			true,
		},
		{
			"integer slice with invalid value and empty value",
			[]interface{}{"one", nil},
			"integer",
			nil,
			true,
		},
		{
			"ae slice is not string slice",
			[]interface{}{1, 2, 3},
			"integer",
			nil,
			true,
		},
	}

	for _, test := range tests {
		actual, err := buildEnum(test.ae, test.elemType)
		if test.shouldError && err == nil {
			t.Errorf("%q - buildEnum should have errored on %+v, %s, got %+v", test.name, test.ae, test.elemType, actual)
		}
		if !reflect.DeepEqual(test.expected, actual) {
			t.Errorf("%q - buildEnum expected %+v, got %+v", test.name, test.expected, actual)
		}
	}
}

func TestNormalizeGroupProfile(t *testing.T) {
	profileWithNils := sdk.GroupProfileMap{
		"test1": "test",
		"test2": nil,
		"test3": true,
		"test4": nil,
		"test5": []string{"a", "b", "c"},
		"test6": 1234,
	}
	normalizedProfile := normalizeGroupProfile(profileWithNils)
	expectedProfile := sdk.GroupProfileMap{
		"test1": "test",
		"test3": true,
		"test5": []string{"a", "b", "c"},
		"test6": 1234,
	}
	assert.Equal(t, normalizedProfile, expectedProfile)
}

func TestCertNormalize(t *testing.T) {
	testCert := `-----BEGIN CERTIFICATE-----
MIIDpDCCAoygAwIBAgIGAXL+Po5gMA0GCSqGSIb3DQEBCwUAMIGSMQswCQYDVQQGEwJVUzETMBEG
A1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzENMAsGA1UECgwET2t0YTEU
MBIGA1UECwwLU1NPUHJvdmlkZXIxEzARBgNVBAMMCmRldi0zODU2NjcxHDAaBgkqhkiG9w0BCQEW
DWluZm9Ab2t0YS5jb20wHhcNMjAwNjI5MDQwMjMyWhcNMzAwNjI5MDQwMzMyWjCBkjELMAkGA1UE
BhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xDTALBgNV
BAoMBE9rdGExFDASBgNVBAsMC1NTT1Byb3ZpZGVyMRMwEQYDVQQDDApkZXYtMzg1NjY3MRwwGgYJ
KoZIhvcNAQkBFg1pbmZvQG9rdGEuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA
iF3iIkzvBZ43ObOfvcWB71EHBlJL/LXJOpomnGdpQQ+ZkxuIVxqnvhHTfIlub3ob5mgGjofI/B12
xQ0CKVpWyTtl6mFVRJsNnu0IJ+64RrA9wXiJhObF5aqHEweLDiZMR/QVFc0MtisjpCoewNSxmWLB
JYaJ84SvvETUM8dvwe7YQ5fU+/psI1w6ydkrcehAWnJ2MC4eFRqNOTM+x/4c4QyL084U1J5azLjY
UtOfbp5bKSoWcSc6mUyNryJfSjKhLba1hrdjBz8hvpmRxUb2rPP1d9IKhZ4s8h+p9dN/IIW6yQZ3
/CKA92ibK3ErHO5x7ivZs11H09UKsKdiRPG8pQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQANhYa/
qyuBoqw6QFJrr1fQxwXfa+zazDcTW1sCXtofgZ77CQoWKqc84C8fCZneDRVExYIYcxfSPY5l75Fv
yag6gpSCa5GsqNKf6AefjXE1gi5mfEqIHCaFcNQX9mxe6ML3zfsqV0rmOLfAiExS28V2rdjIWrKO
pEkANWvDbqL4TOKq5Kr9nD9ItLM2WOBI8SWfNDtGfHiNa1ytVrFNeSBPanTxV1pi50BovU4/JWff
3/ptuMQhKYs9kIP4CFtsQ5ezIFJRq5l9/XiwNYOfP++R4QNKSfCJt6D6ZKN9iLq9YIMJBgb/fd5B
xqneNjZf70DMNAFNXG1VltldQ3hOnRML
-----END CERTIFICATE-----`

	cert, err := certNormalize(testCert)
	if err != nil {
		t.Fatal("failed to normalize PEM cert")
	}

	testCert2 := `MIIDpDCCAoygAwIBAgIGAXL+Po5gMA0GCSqGSIb3DQEBCwUAMIGSMQswCQYDVQQGEwJVUzETMBEG A1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzENMAsGA1UECgwET2t0YTEU MBIGA1UECwwLU1NPUHJvdmlkZXIxEzARBgNVBAMMCmRldi0zODU2NjcxHDAaBgkqhkiG9w0BCQEW DWluZm9Ab2t0YS5jb20wHhcNMjAwNjI5MDQwMjMyWhcNMzAwNjI5MDQwMzMyWjCBkjELMAkGA1UE BhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xDTALBgNV BAoMBE9rdGExFDASBgNVBAsMC1NTT1Byb3ZpZGVyMRMwEQYDVQQDDApkZXYtMzg1NjY3MRwwGgYJ KoZIhvcNAQkBFg1pbmZvQG9rdGEuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA iF3iIkzvBZ43ObOfvcWB71EHBlJL/LXJOpomnGdpQQ+ZkxuIVxqnvhHTfIlub3ob5mgGjofI/B12 xQ0CKVpWyTtl6mFVRJsNnu0IJ+64RrA9wXiJhObF5aqHEweLDiZMR/QVFc0MtisjpCoewNSxmWLB JYaJ84SvvETUM8dvwe7YQ5fU+/psI1w6ydkrcehAWnJ2MC4eFRqNOTM+x/4c4QyL084U1J5azLjY UtOfbp5bKSoWcSc6mUyNryJfSjKhLba1hrdjBz8hvpmRxUb2rPP1d9IKhZ4s8h+p9dN/IIW6yQZ3 /CKA92ibK3ErHO5x7ivZs11H09UKsKdiRPG8pQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQANhYa/ qyuBoqw6QFJrr1fQxwXfa+zazDcTW1sCXtofgZ77CQoWKqc84C8fCZneDRVExYIYcxfSPY5l75Fv yag6gpSCa5GsqNKf6AefjXE1gi5mfEqIHCaFcNQX9mxe6ML3zfsqV0rmOLfAiExS28V2rdjIWrKO pEkANWvDbqL4TOKq5Kr9nD9ItLM2WOBI8SWfNDtGfHiNa1ytVrFNeSBPanTxV1pi50BovU4/JWff 3/ptuMQhKYs9kIP4CFtsQ5ezIFJRq5l9/XiwNYOfP++R4QNKSfCJt6D6ZKN9iLq9YIMJBgb/fd5B xqneNjZf70DMNAFNXG1VltldQ3hOnRML`

	cert2, err := certNormalize(testCert2)
	if err != nil {
		log.Fatal("failed to normalize raw cert")
	}

	if !cert.Equal(cert2) {
		t.Fatalf("certs do not match: A: %s, B: %s", cert.Issuer.CommonName, cert2.Issuer.CommonName)
	}
}

func TestNoChangeInObjectUnmarshaledFromJSON(t *testing.T) {
	testCases := []struct {
		name     string
		oldJSON  string
		newJSON  string
		expected bool
	}{
		{
			name: "there is no change - same same",
			oldJSON: `{
				"one": 1,
				"some": [
					1,
					"a"
				]
			}`,
			newJSON: `{
				"one": 1,
				"some": [
					1,
					"a"
				]
			}`,
			expected: true,
		},
		{
			name: "there is no change - same objects, different string formatting",
			oldJSON: `{
				"one": 1,
				"some": [
					1,
					"a"
				]
			}`,
			newJSON:  `{ "one": 1, "some": [ 1, "a" ] }`,
			expected: true,
		},
		{
			name: "there is no change - attributes in different order ",
			oldJSON: `{
				"one": 1,
				"some": [
					1,
					"a"
				]
			}`,
			newJSON: `{
				"some": [
					1,
					"a"
				],
				"one": 1
			}`,
			expected: true,
		},
		{
			name: "there is change - different values",
			oldJSON: `{
				"one": 1,
				"some": [
					1,
					"a"
				]
			}`,
			newJSON: `{
				"one": 2,
				"some": [
					"a",
					1
				]
			}`,
			expected: false,
		},
		{
			name: "there is change - slice out of order",
			oldJSON: `{
				"one": 1,
				"some": [
					1,
					"a"
				]
			}`,
			newJSON: `{
				"one": 1,
				"some": [
					"a",
					1
				]
			}`,
			expected: false,
		},
		{
			name: "there is no change - new resource value will be blank",
			oldJSON: `{
				"one": 1,
				"some": [
					1,
					"a"
				]
			}`,
			newJSON:  "",
			expected: true,
		},
	}
	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := noChangeInObjectFromUnmarshaledJSON("", tc.oldJSON, tc.newJSON, nil)
			if tc.expected != result {
				t.Errorf("expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	old := []string{"a", "b", "c", "d", "e"}
	new := []string{"c", "d", "e", "f", "g"}
	intersection, exclusiveOld, exclusiveNew := Intersection(old, new)
	assert.Equal(t, []string{"c", "d", "e"}, intersection)
	assert.Equal(t, []string{"a", "b"}, exclusiveOld)
	assert.Equal(t, []string{"f", "g"}, exclusiveNew)
}
