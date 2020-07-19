module github.com/terraform-providers/terraform-provider-okta

go 1.12

require (
	github.com/articulate/oktasdk-go v1.0.3-0.20200311150058-f2661b7a273f
	github.com/beevik/etree v1.1.0 // indirect
	github.com/bflad/tfproviderlint v0.4.0
	github.com/client9/misspell v0.3.4
	github.com/crewjam/saml v0.0.0-20180831135026-ebc5f787b786
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.9.0
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/okta/okta-sdk-golang/v2 v2.0.0
	github.com/peterhellberg/link v1.0.0
	github.com/russellhaering/goxmldsig v0.0.0-20180430223755-7acd5e4a6ef7 // indirect
	github.com/ulikunitz/xz v0.5.6 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	go.opencensus.io v0.22.1 // indirect
	google.golang.org/api v0.10.0 // indirect
	google.golang.org/appengine v1.6.2 // indirect
)

replace github.com/okta/okta-sdk-golang => github.com/articulate/okta-sdk-golang v1.1.1
