.PHONY: \
	help \
	generate \
	build

help: ## help: you are looking at it
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | sed 's/^[^:]*://g' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

generate: ## generate: generate mocks, version and api client code based on the api spec
	@echo "+ $@"
	@go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
	@go install github.com/vektra/mockery/v2@v2.42.0
	@oapi-codegen -package clockodo ./clockodo/apispec.yaml > ./clockodo/clockodo.gen.go
	@go generate
	@mockery

build: generate
	@echo "+ $@"
	@go build -o track cmd/clockodo/main.go
