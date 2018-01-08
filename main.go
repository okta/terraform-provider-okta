package main

import (
    "github.com/articulate/terraform-provider-okta/okta?ref=okta-provider"
    "github.com/hashicorp/terraform/plugin"
    "github.com/hashicorp/terraform/terraform"
)

var Version string

func main() {
    plugin.Serve(&plugin.ServeOpts{
        ProviderFunc: okta.Provider,
    })
}
