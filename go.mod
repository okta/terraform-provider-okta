module github.com/articulate/terraform-provider-okta

go 1.12

require (
	github.com/articulate/id-mapper v0.0.0-20190620150652-1eb71788adf7 // indirect
	github.com/articulate/oktasdk-go v0.0.0-20190417182045-e41ed7befc56
	github.com/beevik/etree v1.1.0 // indirect
	github.com/crewjam/saml v0.0.0-20180831135026-ebc5f787b786
	github.com/hashicorp/go-cleanhttp v0.5.0
	github.com/hashicorp/go-hclog v0.0.0-20190109152822-4783caec6f2e // indirect
	github.com/hashicorp/terraform v0.12.0
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/okta/okta-sdk-golang v0.1.0
	github.com/peterhellberg/link v1.0.0
	github.com/russellhaering/goxmldsig v0.0.0-20180430223755-7acd5e4a6ef7 // indirect
)

replace github.com/okta/okta-sdk-golang => github.com/articulate/okta-sdk-golang v1.0.5
