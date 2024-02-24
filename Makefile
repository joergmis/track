.PHONY: \
	help \
	generate \
	build

help: ## help: you are looking at it
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | sed 's/^[^:]*://g' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

generate: ## generate: generate client code based on the api spec
	@echo "+ $@"
	@go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
	@oapi-codegen -package clockodo ./clockodo/apispec.yaml > ./clockodo/clockodo.gen.go
	@go generate

build: generate
	@echo "+ $@"
	@go build -o track cmd/clockodo/main.go
