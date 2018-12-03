TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
TARGETS=darwin linux windows

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test -v -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test -v $(TESTARGS) -timeout=30s -parallel=4

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

targets: $(TARGETS)

$(TARGETS):
	GOOS=$@ GOARCH=amd64 CGO_ENABLED=0 go build -o "dist/$@/terraform-provider-kubernetes_${TRAVIS_TAG}_x4" -a -ldflags '-extldflags "-static"'
	zip -j dist/terraform-provider-kubernetes_${TRAVIS_TAG}_$@_amd64.zip dist/$@/terraform-provider-kubernetes_${TRAVIS_TAG}_x4

changelog:
	github_changelog_generator --user sl1pm4t --project terraform-provider-kubernetes --release-branch custom

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile targets $(TARGETS)

