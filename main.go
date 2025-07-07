// Package main terraform initial entrypoint & redirect to the okta package
package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6to5server"
	"github.com/okta/terraform-provider-okta/okta/fwprovider"
	"github.com/okta/terraform-provider-okta/okta/provider"
	"github.com/okta/terraform-provider-okta/okta/version"
)

// Ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Generate documentation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	primary := provider.Provider()
	// TODO: Uses v5 protocol for now, however lets swap to v6 when a drop of support for TF versions prior to 1.0 can be made
	framework, err := tf6to5server.DowngradeServer(context.Background(), providerserver.NewProtocol6(fwprovider.NewFrameworkProvider(version.OktaTerraformProviderVersion, primary)))
	if err != nil {
		log.Fatal(err.Error())
	}

	providers := []func() tfprotov5.ProviderServer{
		// v2 plugin
		primary.GRPCProvider,
		// v3 plugin
		func() tfprotov5.ProviderServer {
			return framework
		},
	}
	// At the moment, there is no way to use tf6muxserver to mux the new and old provider because the okta.Provider().GRPCProvider only return protocolv5
	// Most likely we will have to convert all of the old provider to new provider to use v6
	// There are constraint of using v5, such as not being able to use NestedAttributes, we have to use NestedBlock instead (see okta_org_metada)
	// use the muxer
	muxServer, err := tf5muxserver.NewMuxServer(context.Background(), providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	err = tf5server.Serve(
		"okta/okta",
		muxServer.ProviderServer,
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
}
