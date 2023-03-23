package sdk

// APISupplement not all APIs are supported by okta-sdk-golang, this will act as a supplement to the Okta SDK
type APISupplement struct {
	RequestExecutor *RequestExecutor
}

// CloneRequestExecutor create a clone of the underlying request executor
func (m *APISupplement) cloneRequestExecutor() *RequestExecutor {
	a := *m.RequestExecutor
	return &a
}
