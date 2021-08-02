package apimutex

import (
	"strings"
	"sync"
)

// ApiMutex synchronizes keeping account of current known rate limit values from
// Okta management endpoints. Specifically apps, users, and other, see:
// https://developer.okta.com/docs/reference/rl-global-mgmt/ The Okta Terraform
// Provider can not account for all other kinds of clients utilization of API
// limits but it account for its own usage and attempt to preemptively react
// appropriately.
type ApiMutex struct {
	lock   sync.Mutex
	status map[string]*ApiStatus
}

// ApiStatus is used to hold rate limit information from Okta's API, see:
// https://developer.okta.com/docs/reference/rl-best-practices/
type ApiStatus struct {
	limit     int
	remaining int
	reset     int64 // UTC epoch time in seconds
}

// NewApiMutex returns a new api mutex object.
func NewApiMutex() *ApiMutex {
	return &ApiMutex{
		status: make(map[string]*ApiStatus),
	}
}

// Update updates the known status for the given API endpoint. It is synchronous
// and intellegently accounts for new values regardless of parallelism.
func (m *ApiMutex) Update(endPoint string, limit, remaining int, reset int64) {
	_ = m.get(endPoint) // get will initialize api status structs

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

	if remaining > status.remaining {
		status.remaining = remaining
	}
}

// Status return the ApiStatus for the given class of endpoint.
func (m *ApiMutex) Status(endPoint string) *ApiStatus {
	return m.get(endPoint)
}

func (m *ApiMutex) normalizeKey(endPoint string) string {
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

func (m *ApiMutex) get(endPoint string) *ApiStatus {
	m.lock.Lock()
	defer m.lock.Unlock()

	key := m.normalizeKey(endPoint)
	status, found := m.status[key]
	if !found {
		status = &ApiStatus{}
		m.status[key] = status
	}

	return status
}
