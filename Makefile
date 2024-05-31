.PHONY: \
	generate \
	install \
	test \
	mutation-testing \
	lint \
	clean

generate: clean
	go generate
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
	oapi-codegen -package api ./clockodo/api/apispec.yaml > ./clockodo/api/clockodo.gen.go

test: lint
	go test -coverpkg=./... ./...

mutation-tests:
	go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest
	go-mutesting github.com/joergmis/track
	go-mutesting github.com/joergmis/track/local

install: test generate
	go install ./cmd/track

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.0
	golangci-lint run ./...
