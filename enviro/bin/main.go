package main

import (
	"github.com/peter-mount/go-kernel"
	"github.com/peter-mount/go-kernel/rest"
	"github.com/peter-mount/home-automation/enviro"
	"github.com/peter-mount/home-automation/graphite"
	"log"
)

// main standalone app for developing the graphite population
func main() {
	err := kernel.Launch(
		&graphite.Graphite{},
		&enviro.Enviro{},
		&rest.Server{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
