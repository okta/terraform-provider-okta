// Package main terraform initial entrypoint & redirect to the okta package
package main

import (
	"github.com/articulate/terraform-provider-okta/okta"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: okta.Provider,
	})
}
