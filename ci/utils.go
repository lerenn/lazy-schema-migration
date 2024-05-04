package main

import (
	"context"
	"sync"

	"dagger.io/dagger"
)

func sourceAsWorkdir(client *dagger.Client) func(r *dagger.Container) *dagger.Container {
	// Set path where the source code is mounted.
	containerDir := "/go/src/github.com/lerenn/lazy-schema-migration"

	return func(r *dagger.Container) *dagger.Container {
		return r.
			// Add Go caches
			WithMountedCache("/root/.cache/go-build", client.CacheVolume("gobuild")).
			WithMountedCache("/go/pkg/mod", client.CacheVolume("gocache")).

			// Add source code
			WithMountedDirectory(containerDir, client.Host().Directory(".")).

			// Add workdir
			WithWorkdir(containerDir)
	}
}

func executeContainers(ctx context.Context, containers ...[]*dagger.Container) {
	// Regroup arg
	rContainers := make([]*dagger.Container, 0)
	for _, c := range containers {
		rContainers = append(rContainers, c...)
	}

	// Excute containers
	var wg sync.WaitGroup
	for _, ec := range rContainers {
		go func(e *dagger.Container) {
			_, err := e.Stderr(ctx)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(ec)

		wg.Add(1)
	}

	wg.Wait()
}
