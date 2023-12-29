package main

import "dagger.io/dagger"

const (
	// TestsImage is the docker image used to execute tests.
	TestsImage = "golang:1.21.4"
)

// Tests returns a dagger container to run tests.
func Tests(client *dagger.Client) *dagger.Container {
	return client.Container().
		From(TestsImage).
		With(sourceAsWorkdir(client)).
		WithExec([]string{"go", "test", "./..."})
}
