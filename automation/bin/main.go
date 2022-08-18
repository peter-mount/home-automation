package main

import (
	"github.com/peter-mount/go-kernel"
	"github.com/peter-mount/go-kernel/rest"
	"github.com/peter-mount/home-automation"
	"github.com/peter-mount/home-automation/automation/api"
	"github.com/peter-mount/home-automation/util/graphite"
	"log"
)

func main() {
	err := kernel.Launch(
		&graphite.Graphite{},
		&automation.Zigbee{},
		&api.Service{},
		&rest.Server{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
