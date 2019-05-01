SWEEP?=global

# If test pkgs is not set find/use all non-vendor packages
ifndef TEST_PKGS
	TEST_PKGS := $$(go list ./... |grep -v 'vendor')
endif

# Expression to match against tests
# go test -run <filter>
# e.g. Iden will run all TestAccIdentity tests
ifdef TEST_FILTER
	TEST_FILTER := -run $(TEST_FILTER)
endif

# Pass additional go test <args> via TEST_ARGS
# TEST_ARGS

GOFMT_FILES := $(filter-out vendor, $(wildcard **/*.go))
export GO111MODULE=on

# Last tagged version
VERSION = $$(git tag --sort=v:refname | tail -1)


default: build-plugins

# Builds a binary for current OS and Arch
build: fmtcheck
	@mkdir -p ~/.terraform.d/plugins/
	@go build -o terraform-provider-okta_${VERSION}

# Builds a binary for Linux, Windows, and OSX and installs it in the default terraform plugins directory
build-plugins: fmtcheck
	@mkdir -p ~/.terraform.d/plugins/
	gox -osarch="linux/amd64 darwin/amd64 windows/amd64" \
	  -output="${HOME}/.terraform.d/plugins/{{.OS}}_{{.Arch}}/terraform-provider-okta_${VERSION}_x4" .

ship: build-plugins
	exists=$$(aws s3api list-objects --bucket articulate-terraform-providers --profile prod --prefix terraform-provider-okta --query Contents[].Key | jq 'contains(["${VERSION}"])' ) \
	&& if [ $$exists == "true" ]; then \
	  echo "[ERROR] terraform-provider-okta_${VERSION} already exists in s3://${TERRAFORM_PLUGINS_BUCKET} - don't forget to bump the version."; else \
		echo "copying terraform-provider-okta_${VERSION} to s3://${TERRAFORM_PLUGINS_BUCKET}"; \
	  aws s3 cp ~/.terraform.d/plugins/linux_amd64/terraform-provider-okta_${VERSION}_x4  s3://${TERRAFORM_PLUGINS_BUCKET}/ --profile ${TERRAFORM_PLUGINS_PROFILE}; \
	fi

test: fmtcheck
	go test -i $(TEST_PKGS) || exit 1
	echo $(TEST_PKGS) | \
		xargs -t -n4 go test $(TEST_ARGS) $(TEST_FILTER) \
			-timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test -v $(filter-out -v,$(TEST_ARGS)) $(TEST_FILTER) \
		$(TEST_PKGS) -timeout 120m

# Sweeps up leaked dangling resources
sweep:
	@echo "WARNING: This will destroy resources. Use only in development accounts."
	go test $(TEST_PKGS) -v -sweep=$(SWEEP) $(SWEEPARGS)

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
	@if [ "$(TEST_PKGS)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./aws"; \
		exit 1; \
	fi
	go test -c $(TEST_PKGS) $(TEST_ARGS)

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile

# Print the value of any variable as make print-VAR
print-%  : ; @echo $* = $($*)
