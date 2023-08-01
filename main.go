// Package main terraform initial entrypoint & redirect to the okta package
package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/okta/terraform-provider-okta/okta"
)

// Ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Generate documentation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	var debug bool
	providers := []func() tfprotov5.ProviderServer{
		// v2 plugin
		okta.Provider().GRPCProvider,
		// v3 plugin
		providerserver.NewProtocol5(okta.NewFWProvider(okta.OktaTerraformProviderVersion)),
	}

	// use the muxer
	muxServer, err := tf5muxserver.NewMuxServer(context.Background(), providers...)
	if err != nil {
		log.Fatalf(err.Error())
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
