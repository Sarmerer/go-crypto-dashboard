package main

import (
	"fmt"

	"github.com/sarmerer/go-crypto-dashboard/passivbotmgr/droplet"
)

func main() {
	pm := droplet.NewProcessManager()
	pid, err := pm.Start([]string{"python", "test.py", ">", "test.log"})
	if err != nil {
		panic(err)
	}

	fmt.Printf("PID: %d\n", pid)
}
