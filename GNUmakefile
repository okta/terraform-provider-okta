VERSION=v2.2.2

SWEEP?=global
TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: deps build

deps:
	curl -s https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	${GOPATH}/bin/dep ensure

build: fmtcheck build-macos build-linux

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o terraform-provider-okta_${VERSION}
	mkdir -p ~/.terraform.d/plugins/darwin_amd64/
	mv terraform-provider-okta_${VERSION} ~/.terraform.d/plugins/darwin_amd64/

build-linux:
	GOOS=linux GOARCH=amd64 go build -o terraform-provider-okta_${VERSION}
	mkdir -p ~/.terraform.d/plugins/linux_amd64/
	mv terraform-provider-okta_${VERSION} ~/.terraform.d/plugins/linux_amd64/

ship: build
	exists=$$(aws s3api list-objects --bucket articulate-terraform-providers --profile prod --prefix terraform-provider-okta --query Contents[].Key | jq 'contains(["${VERSION}"])' ) \
	&& if [ $$exists == "true" ]; then \
	  echo "[ERROR] terraform-provider-okta_${VERSION} already exists in s3://${TERRAFORM_PLUGINS_BUCKET} - don't forget to bump the version."; else \
		echo "copying terraform-provider-okta_${VERSION} to s3://${TERRAFORM_PLUGINS_BUCKET}"; \
	  aws s3 cp ~/.terraform.d/plugins/linux_amd64/terraform-provider-okta_${VERSION}  s3://${TERRAFORM_PLUGINS_BUCKET}/ --profile ${TERRAFORM_PLUGINS_PROFILE}; \
	fi

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck sweep
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

# Sweeps up leaked dangling resources
sweep:
	@echo "WARNING: This will destroy resources. Use only in development accounts."
	go test $(TEST) -v -sweep=$(SWEEP) $(SWEEPARGS)

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./aws"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile
