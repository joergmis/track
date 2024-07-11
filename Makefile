.PHONY: \
	generate \
	install \
	test \
	mutation-testing \
	lint \
	completion \
	clean

generate: clean
	go generate
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	oapi-codegen -package api ./clockodo/api/apispec.yaml > ./clockodo/api/clockodo.gen.go

test: lint
	go test -coverpkg=./... ./...

mutation-tests:
	go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest
	go-mutesting github.com/joergmis/track || true
	go-mutesting github.com/joergmis/track/local || true

install: test generate
	go install ./cmd/track

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.0
	golangci-lint run ./...

swagger:
	docker pull swaggerapi/swagger-editor
	docker run -p 80:8080 -v ./:/tmp -e SWAGGER_FILE=/tmp/clockodo/api/apispec.yaml swaggerapi/swagger-editor

completion:
	track completion bash > /etc/bash_completion.d/track
