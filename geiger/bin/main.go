package main

import (
	"github.com/peter-mount/go-kernel"
	"github.com/peter-mount/go-kernel/rest"
	"github.com/peter-mount/home-automation/geiger"
	"log"
)

// main standalone app for developing the graphite population
func main() {
	err := kernel.Launch(
		&geiger.Geiger{},
		&rest.Server{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
