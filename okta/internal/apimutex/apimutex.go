package apimutex

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	APPS_KEY  = "apps"
	USERS_KEY = "users"
	OTHER_KEY = "other"
)

// APIMutex synchronizes keeping account of current known rate limit values
// from Okta management endpoints. Specifically apps, users, and other, see:
// https://developer.okta.com/docs/reference/rl-global-mgmt/ The Okta Terraform
// Provider can not account for all other kinds of clients utilization of API
// limits but it can account for its own usage and attempt to preemptively
// react appropriately.
type APIMutex struct {
	lock     sync.Mutex
	status   map[string]*APIStatus
	capacity int
}

// APIStatus is used to hold rate limit information from Okta's API, see:
// https://developer.okta.com/docs/reference/rl-best-practices/
type APIStatus struct {
	limit     int
	remaining int
	reset     int64 // UTC epoch time in seconds
}

// NewAPIMutex returns a new api mutex object that represents untilized
// capacity under the specified capacity percentage.
func NewAPIMutex(capacity int) (*APIMutex, error) {
	if capacity < 1 || capacity > 100 {
		return nil, fmt.Errorf("expecting capacity as whole number > 0 and <= 100, was %d", capacity)
	}
	status := map[string]*APIStatus{
		APPS_KEY:  {},
		USERS_KEY: {},
		OTHER_KEY: {},
	}
	return &APIMutex{
		capacity: capacity,
		status:   status,
	}, nil
}

// HasCapacity approximates if there is capacity below the api mutex's maximum
// capacity threshold.
func (m *APIMutex) HasCapacity(endPoint string) bool {
	status := m.get(endPoint)

	// if the status hasn't been updated recently assume there is capacity
	if status.reset+60 < time.Now().Unix() {
		return true
	}

	// calculate utilization
	utilization := 100.0 * (float32(status.limit-status.remaining) / float32(status.limit))

	return utilization <= float32(m.capacity)
}

// Update updates the known status for the given API endpoint. It is synchronous
// and intelligently accounts for new values regardless of parallelism.
func (m *APIMutex) Update(endPoint string, limit, remaining int, reset int64) {
	m.lock.Lock()
	defer m.lock.Unlock()

	key := m.normalizeKey(endPoint)
	status := m.status[key]

	if reset > status.reset {
		// reset value greater than current reset implies we are in a new Okta API
		// one minute window. set/reset values.
		status.reset = reset
		status.remaining = remaining
		status.limit = limit
		return
	}

	if reset <= (status.reset - 60) {
		// these values are from the previous one minute window, ignore
		return
	}

	if remaining < status.remaining {
		status.remaining = remaining
	}
}

// Status return the APIStatus for the given class of endpoint.
func (m *APIMutex) Status(endPoint string) *APIStatus {
	return m.get(endPoint)
}

func (m *APIMutex) normalizeKey(endPoint string) string {
	var result string
	switch {
	case strings.HasPrefix(endPoint, "/api/v1/apps"):
		result = "apps"
	case strings.HasPrefix(endPoint, "/api/v1/users"):
		result = "users"
	default:
		result = "other"
	}
	return result
}

// Reset returns the current reset value of the api status object.
func (s *APIStatus) Reset() int64 {
	return s.reset
}

// Limit returns the current limit value of the api status object.
func (s *APIStatus) Limit() int {
	return s.limit
}

// Remaining returns the current remaining value of the api status object.
func (s *APIStatus) Remaining() int {
	return s.remaining
}

func (m *APIMutex) get(endPoint string) *APIStatus {
	key := m.normalizeKey(endPoint)
	return m.status[key]
}
