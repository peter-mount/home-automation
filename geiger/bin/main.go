package main

import (
	"github.com/peter-mount/go-kernel"
	"github.com/peter-mount/go-kernel/rest"
	"github.com/peter-mount/home-automation/geiger"
	"github.com/peter-mount/home-automation/graphite"
	"log"
)

// main standalone app for developing the graphite population
func main() {
	err := kernel.Launch(
		&graphite.Graphite{},
		&geiger.Geiger{},
		&rest.Server{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
