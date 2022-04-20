package main

import (
	"log"
	"os"

	"github.com/sarmerer/go-crypto-dashboard/passivbotmgr/droplet"
	"github.com/sarmerer/go-crypto-dashboard/passivbotmgr/droplet/legacy"
)

func main() {
	var d droplet.Droplet

	args := os.Args[1:]
	if args[0] == "-l" || args[0] == "--legacy" {
		d = legacy.New(8080)
	} else {
		d = droplet.New(8080)
	}

	log.Fatal(d.ListenAndServe())
}
