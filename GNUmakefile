SWEEP?=global
TEST?=$$(go list ./... |grep -v 'vendor')
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=okta
GOFMT:=gofumpt
TFPROVIDERLINT=tfproviderlint
STATICCHECK=staticcheck

# Expression to match against tests
# go test -run <filter>
# e.g. Iden will run all TestAccIdentity tests
ifdef TEST_FILTER
	TEST_FILTER := -run $(TEST_FILTER)
endif

DEFAULT_SMOKE_TESTS?=\
  TestAccDataSourceOktaUser_read \
  TestAccResourceOktaUserSchema_crud

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
	go clean -cache -testcache ./...

clean-all:
	go clean -cache -testcache -modcache ./...

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -v -sweep=$(SWEEP) $(SWEEPARGS)

test:
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) $(TEST_FILTER) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) $(TEST_FILTER) -timeout 120m

test-play-vcr-acc:
	OKTA_VCR_TF_ACC=play TF_ACC=1 go test -tags unit -mod=readonly -test.v -timeout 120m ./okta

smoke-test-play-vcr-acc:
	OKTA_VCR_TF_ACC=play TF_ACC=1 go test -tags unit -mod=readonly -test.v -timeout 120m -run ^$(smoke_tests)$$ ./okta

test-record-vcr-acc:
	OKTA_VCR_TF_ACC=record TF_ACC=1 go test -tags unit -mod=readonly -test.v -timeout 120m ./okta

vet:
	@echo "==> Checking source code against go vet and staticcheck"
	@go vet ./...
	@staticcheck ./...

fmt: tools # Format the code
	@echo "formatting the code with $(GOFMT)..."
	@$(GOFMT) -l -w .

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
	@echo "==> Checking source code against linters..."
	@$(TFPROVIDERLINT) \
		-c 1 \
		-AT001 \
    -R004 \
		-S001 \
		-S002 \
		-S003 \
		-S004 \
		-S005 \
		-S007 \
		-S008 \
		-S009 \
		-S010 \
		-S011 \
		-S012 \
		-S013 \
		-S014 \
		-S015 \
		-S016 \
		-S017 \
		-S019 \
		./$(PKG_NAME)

tools:
	@which $(GOFMT) || go install mvdan.cc/gofumpt@v0.4.0
	@which $(TFPROVIDERLINT) || go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@v0.28.1
	@which $(STATICCHECK) || go install honnef.co/go/tools/cmd/staticcheck@v0.4.2

tools-update:
	@go install mvdan.cc/gofumpt@v0.4.0
	@go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@v0.28.1
	@go install honnef.co/go/tools/cmd/staticcheck@v0.4.2

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
