package okta

import (
	"github.com/articulate/oktasdk-go/okta"
)

// Config is a struct containing our provider schema values
// plus the okta client object
type Config struct {
	orgName  string
	domain   string
	apiToken string

	oktaClient  *okta.Client
}

func (c *Config) loadAndValidate() error {
	client, err := okta.NewClientWithDomain(nil, c.orgName, c.domain, c.apiToken)
	if err != nil {
		return err
	}
	c.oktaClient = client
	return nil
}
