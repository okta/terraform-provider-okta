// Package main terraform initial entrypoint & redirect to the okta package
package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
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

	// Upgrade the SDKv2 provider (protocol v5) to protocol v6 so it can be
	// muxed alongside the framework provider (which uses protocol v6 natively).
	// This allows framework resources to use SingleNestedAttribute and other
	// features not available in protocol v5.
	upgradedPrimary, err := tf5to6server.UpgradeServer(context.Background(), primary.GRPCProvider)
	if err != nil {
		log.Fatal(err.Error())
	}

	providers := []func() tfprotov6.ProviderServer{
		// SDKv2 provider upgraded to protocol v6
		func() tfprotov6.ProviderServer {
			return upgradedPrimary
		},
		// Framework provider at protocol v6 (supports SingleNestedAttribute)
		providerserver.NewProtocol6(fwprovider.NewFrameworkProvider(version.OktaTerraformProviderVersion, primary)),
	}

	muxServer, err := tf6muxserver.NewMuxServer(context.Background(), providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"okta/okta",
		muxServer.ProviderServer,
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
}
