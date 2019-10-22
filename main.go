// Package main terraform initial entrypoint & redirect to the okta package
package main

import (
	"github.com/terraform-providers/terraform-provider-okta/okta"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: okta.Provider,
	})
}
