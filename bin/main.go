package main

import (
	"github.com/peter-mount/go-kernel"
	"github.com/peter-mount/go-kernel/rest"
	"github.com/peter-mount/home-automation"
	"github.com/peter-mount/home-automation/api"
	"log"
)

func main() {
	err := kernel.Launch(
		&automation.Service{},
		&api.Service{},
		&rest.Server{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
