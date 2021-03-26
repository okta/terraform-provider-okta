SWEEP?=global
TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=okta
GOLINT=./bin/golangci-lint
GOFMT:=gofumpt
TFPROVIDERLINT=tfproviderlint

# Expression to match against tests
# go test -run <filter>
# e.g. Iden will run all TestAccIdentity tests
ifdef TEST_FILTER
	TEST_FILTER := -run $(TEST_FILTER)
endif

default: build

dep: # Download required dependencies
	go mod tidy
	go mod vendor

build: fmtcheck
	go install

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -v -sweep=$(SWEEP) $(SWEEPARGS)

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) $(TEST_FILTER) -timeout=30s -parallel=4
    test: fmtcheck

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) $(TEST_FILTER) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: fmt
fmt: check-fmt # Format the code
	@echo "formatting the code..."
	@$(GOFMT) -l -w $$(find . -name '*.go' |grep -v vendor)

check-fmt:
	@which $(GOFMT) > /dev/null || (echo "downloading formatter..." && GO111MODULE=on go get mvdan.cc/gofumpt)

fmtcheck: dep
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

lint: tools
	@echo "==> Checking source code against linters..."
	@GOGC=30 $(GOLINT) run ./$(PKG_NAME)
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
	@which $(GOLINT) || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.33.0
	@which $(TFPROVIDERLINT) || go install github.com/bflad/tfproviderlint/cmd/tfproviderlint

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
