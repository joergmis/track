package main

import (
	"context"
)

type Track struct{}

func (m *Track) Base() *Container {
	return dag.Container().
		From("alpine:latest").
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
			"gem",
			"install",
			"asciidoctor-pdf",
		})
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
