.PHONY: \
	generate \
	install \
	clean

generate: clean
	go generate
	oapi-codegen -package api ./clockodo/api/apispec.yaml > ./clockodo/api/clockodo.gen.go

test:
	go test -coverpkg=./... ./...

install: generate
	go install ./cmd/track
