SWEEP?=global
UNIT_TESTS?=$$(go list ./... |grep -v 'vendor'|grep -v 'okta.services.idaas')
ACC_TESTS=./okta/services/idaas
TEST?=$$(go list ./... |grep -v 'vendor')
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=./okta/services/idaas
GOFMT:=gofumpt
TFPROVIDERLINT=tfproviderlint
TFPROVIDERLINTX=tfproviderlintx
STATICCHECK=staticcheck

# Expression to match against tests
# go test -run <filter>
# e.g. Iden will run all TestAccIdentity tests
ifdef TEST_FILTER
	TEST_FILTER := -run $(TEST_FILTER)
endif

TESTARGS?=-test.v

DEFAULT_SMOKE_TESTS?=\
  TestAccDataSourceOktaAppSaml_read \
  TestAccDataSourceOktaApp_read \
  TestAccDataSourceOktaGroup_read \
  TestAccDataSourceOktaGroups_read \
  TestAccDataSourceOktaPolicy_read \
  TestAccDataSourceOktaUser_read \
  TestAccResourceOktaAppAutoLoginApplication_crud \
  TestAccResourceOktaAppBasicAuthApplication_crud \
  TestAccResourceOktaAppBookmarkApplication_crud \
  TestAccResourceOktaAppOauth_basic \
  TestAccResourceOktaAppOauth_serviceWithJWKS \
  TestAccResourceOktaAppSaml_crud \
  TestAccResourceOktaAppSwaApplication_crud \
  TestAccResourceOktaAppThreeFieldApplication_crud \
  TestAccResourceOktaAppUser_crud \
  TestAccResourceOktaDefaultMFAPolicy \
  TestAccResourceOktaGroup_crud \
  TestAccResourceOktaMfaPolicyRule_crud \
  TestAccResourceOktaMfaPolicy_crud \
  TestAccResourceOktaOrgConfiguration_crud \
  TestAccResourceOktaPolicyRulePassword_crud \
  TestAccResourceOktaUser_updateAllAttributes

ifeq ($(strip $(SMOKE_TESTS)),)
	SMOKE_TESTS = $(DEFAULT_SMOKE_TESTS)
endif

space := $(subst ,, )
smoke_tests := $(subst $(space),\|,$(SMOKE_TESTS))

default: build

dep: # Download required dependencies
	go mod tidy

docs:
	go generate

build: fmtcheck
	go install

clean:
	go clean -cache -testcache

clean-all:
	go clean -cache -testcache -modcache

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -sweep=$(SWEEP) $(SWEEPARGS)

test:
	echo $(UNIT_TESTS) | \
		xargs -t -n4 go test $(TESTARGS) $(TEST_FILTER) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(ACC_TESTS) $(TESTARGS) $(TEST_FILTER) -timeout 120m

test-play-vcr-acc:
	OKTA_VCR_TF_ACC=play TF_ACC=1 go test -tags unit -mod=readonly -test.v -timeout 120m $(PKG_NAME)

test-play-vcr-acc-governance:
	OKTA_VCR_TF_ACC=play TF_ACC=1 go test -tags unit -mod=readonly -test.v -timeout 120m ./okta/services/governance

smoke-test-play-vcr-acc:
	OKTA_VCR_TF_ACC=play TF_ACC=1 go test -tags unit -mod=readonly -test.v -timeout 120m -run ^$(smoke_tests)$$ $(ACC_TESTS)

test-record-vcr-acc:
	OKTA_VCR_TF_ACC=record TF_ACC=1 go test -tags unit -mod=readonly -test.v -timeout 120m $(ACC_TESTS)

qc: fmtcheck vet staticcheck lint

vet:
	@echo "==> Checking source code against go vet"
	@go vet ./...

staticcheck:
	@echo "==> Checking source code against staticcheck"
	@staticcheck ./...

fmt: tools # Format the code
	@echo "formatting the code with $(GOFMT)..."
	@$(GOFMT) -l -w .
	@terraform fmt -recursive ./examples/

fmtcheck:
	@gofumpt -d -l .

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

lint:
	@echo "==> Checking source code against linter tfproviderlint ..."
	@$(TFPROVIDERLINT) -c 1 ./...

lintx:
	# NOTE tfproviderlintx is very opinionated, don't add it to qc target
	@echo "==> Checking source code against linter tfproviderlintx ..."
	@$(TFPROVIDERLINTX) -c 1 ./...

tools:
	@which $(GOFMT) || go install mvdan.cc/gofumpt@v0.7.0
	@which $(TFPROVIDERLINT) || go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@v0.31.0
	@which $(TFPROVIDERLINTX) || go install github.com/bflad/tfproviderlint/cmd/tfproviderlintx@v0.31.0
	@which $(STATICCHECK) || go install honnef.co/go/tools/cmd/staticcheck@v0.6.1

tools-update:
	@go install mvdan.cc/gofumpt@v0.7.0
	@go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@v0.31.0
	@go install github.com/bflad/tfproviderlint/cmd/tfproviderlintx@v0.31.0
	@go install honnef.co/go/tools/cmd/staticcheck@v0.6.1

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile website website-test
