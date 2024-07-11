package main

import (
	"context"
)

type Track struct{}

func (m *Track) Base() *Container {
	return dag.Container().
		From("golang:alpine").
		WithExec([]string{
			"apk",
			"update",
		}).
		WithExec([]string{
			"apk",
			"add",
			"ruby",
		}).
		WithExec([]string{
			"apk",
			"add",
			"make",
		}).
		WithExec([]string{
			"gem",
			"install",
			"asciidoctor-pdf",
		}).
		WithExec([]string{
			"go",
			"install",
			"github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.0",
		}).
		WithExec([]string{
			"go",
			"install",
			"github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest",
		}).
		WithExec([]string{
			"go",
			"install",
			"github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest",
		})
}

func (m *Track) Test(ctx context.Context, source *Directory) (string, error) {
	return m.Base().
		WithMountedDirectory("/mnt", source).
		WithWorkdir("/mnt").
		WithExec([]string{
			"printf", "%s", "$(git describe --tags --dirty)", ">", "VERSION.txt",
		}).
		WithExec([]string{
			"make",
			"tests",
		}).
		WithExec([]string{
			"make",
			"mutation-tests",
		}).
		Stdout(ctx)
}

func (m *Track) GenerateArchitecture(ctx context.Context, source *Directory) *File {
	return m.Base().
		WithMountedDirectory("/mnt", source).
		WithWorkdir("/mnt/doc").
		WithExec([]string{
			"asciidoctor-pdf",
			"architecture.adoc",
		}).
		File("/mnt/doc/architecture.pdf")
}
