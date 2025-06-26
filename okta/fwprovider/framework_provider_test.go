package fwprovider

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
