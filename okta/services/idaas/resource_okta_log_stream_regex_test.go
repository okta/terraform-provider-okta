package idaas

import "testing"

func TestAwsEventBridgeEventSourceNameRegex(t *testing.T) {
	valid := []string{"my-event-source", "my_event_source", "my.event.source", "MyEventSource1", "a-b_c.d"}
	for _, name := range valid {
		if !awsEventBridgeEventSourceNameRegex.MatchString(name) {
			t.Errorf("expected %q to be accepted", name)
		}
	}
	invalid := []string{"", "has space", "back\\slash", "bracket]", "caret^", "slash/name"}
	for _, name := range invalid {
		if awsEventBridgeEventSourceNameRegex.MatchString(name) {
			t.Errorf("expected %q to be rejected", name)
		}
	}
}
