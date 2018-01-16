package okta

import (
	"log"

	"github.com/articulate/oktasdk-go/okta"
)

type Config struct {
	orgName  string
	domain   string
	apiToken string

	oktaClient  *okta.Client
}

func (c *Config) loadAndValidate() error {
	client, err := okta.NewClientWithDomain(nil, c.orgName, c.domain, c.apiToken)
	if err != nil {
		log.Println("[ERROR] Error Initializing Okta client:\n \t%v", err)
		return err
	}
	c.oktaClient = client
	return nil
}
