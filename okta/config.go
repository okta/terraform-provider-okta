package okta

import (
	"fmt"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/okta/okta-sdk-golang/okta"
)

// Config is a struct containing our provider schema values
// plus the okta client object
type Config struct {
	orgName  string
	domain   string
	apiToken string

	articulateOktaClient *articulateOkta.Client
	oktaClient           *okta.Client
}

func (c *Config) loadAndValidate() error {
	articulateClient, err := articulateOkta.NewClientWithDomain(nil, c.orgName, c.domain, c.apiToken)

	// add the Articulaet Okta client object to Config
	c.articulateOktaClient = articulateClient

	if err != nil {
		return fmt.Errorf("[ERROR] Error creating Articulate Okta client: %v", err)
	}

	orgUrl := fmt.Sprintf("https://%v.%v", c.orgName, c.domain)

	config := okta.NewConfig().WithOrgUrl(orgUrl).WithToken(c.apiToken)
	client := okta.NewClient(config, nil, nil)

	// add the Okta SDK client object to Config
	c.oktaClient = client
	return nil
}
