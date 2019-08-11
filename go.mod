module github.com/articulate/terraform-provider-okta

go 1.12

require (
	github.com/articulate/oktasdk-go v0.0.0-20190417182045-e41ed7befc56
	github.com/articulate/terraform v0.6.12 // indirect
	github.com/aws/aws-sdk-go v1.20.3 // indirect
	github.com/beevik/etree v1.1.0 // indirect
	github.com/crewjam/saml v0.0.0-20180831135026-ebc5f787b786
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-cleanhttp v0.5.0
	github.com/hashicorp/go-hclog v0.0.0-20190109152822-4783caec6f2e // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform v0.12.0
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/okta/okta-sdk-golang v0.1.0
	github.com/peterhellberg/link v1.0.0
	github.com/russellhaering/goxmldsig v0.0.0-20180430223755-7acd5e4a6ef7 // indirect
	golang.org/x/crypto v0.0.0-20190513172903-22d7a77e9e5f // indirect
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092 // indirect
	golang.org/x/sys v0.0.0-20190528012530-adf421d2caf4 // indirect
)

replace github.com/okta/okta-sdk-golang => github.com/articulate/okta-sdk-golang v0.0.0-20190810203837-596e830a6cb1
