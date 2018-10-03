package okta

import (
	"strings"
	"testing"
)

func TestContainsAll(t *testing.T) {
	testArr := []string{"1", "2", "3"}
	tests := []struct {
		elements []string
		expected bool
	}{
		{[]string{"2", "3"}, true},
		{[]string{"1", "2"}, true},
		{[]string{"1", "4"}, false},
		{[]string{"1"}, true},
		{[]string{"10"}, false},
	}

	for _, test := range tests {
		actual := containsAll(testArr, test.elements...)

		if actual != test.expected {
			t.Errorf("containsAll test failed, test array: %s, elements %s, expected %t, actual %t", strings.Join(testArr, ", "), strings.Join(test.elements, ", "), test.expected, actual)
		}
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
