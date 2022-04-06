package okta

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
