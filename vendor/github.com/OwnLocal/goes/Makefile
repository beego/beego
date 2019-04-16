help:
	@echo "Available targets:"
	@echo "- test: run tests"
	@echo "- deps: installs dependencies with glide"
	@echo "- watch: watch for changes and re-run tests"

deps:
	glide install	&& mkdir -p vendor/bin && go build -o vendor/bin/ginkgo ./vendor/github.com/onsi/ginkgo/ginkgo


test: deps
	vendor/bin/ginkgo -race -randomizeAllSpecs -r -skipPackage vendor -progress .

watch: deps
	vendor/bin/ginkgo watch -race -randomizeAllSpecs -r -skipPackage vendor -progress -notify .
