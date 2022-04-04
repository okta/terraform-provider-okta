// Package main terraform initial entrypoint & redirect to the okta package
package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/okta/terraform-provider-okta/okta"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: okta.Provider,
	})
}
