.PHONY: \
	generate \
	install \
	clean

clean:
	rm -rf ./mocks/*

generate: clean
	go generate
	oapi-codegen -package api ./clockodo/api/apispec.yaml > ./clockodo/api/clockodo.gen.go
	mockery

test: generate
	go test -v ./...

install: generate
	go install ./cmd/track
