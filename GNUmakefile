VERSION=v1.1.5

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

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

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
