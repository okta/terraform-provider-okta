package fwprovider

// TODO placeholder to add VCR testing to plugin framework provider based
// resources and datasources
func NewFrameworkProviderTest(testName string) *frameworkProviderTest {
	return &frameworkProviderTest{
		FrameworkProvider: FrameworkProvider{
			Version: "test",
		},
		TestName: testName,
	}
}

type frameworkProviderTest struct {
	FrameworkProvider
	TestName string
}
