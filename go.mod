module github.com/terraform-providers/terraform-provider-okta

go 1.12

require (
	github.com/articulate/oktasdk-go v1.0.3-0.20200311150058-f2661b7a273f
	github.com/beevik/etree v1.1.0 // indirect
	github.com/bflad/tfproviderlint v0.4.0
	github.com/client9/misspell v0.3.4
	github.com/crewjam/saml v0.0.0-20180831135026-ebc5f787b786
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/terraform-plugin-sdk v1.15.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.3
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/okta/okta-sdk-golang v0.1.0
	github.com/peterhellberg/link v1.0.0
	github.com/russellhaering/goxmldsig v0.0.0-20180430223755-7acd5e4a6ef7 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
)

replace github.com/okta/okta-sdk-golang => github.com/articulate/okta-sdk-golang v1.1.1
