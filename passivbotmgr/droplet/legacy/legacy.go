package legacy

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sarmerer/go-crypto-dashboard/passivbotmgr/droplet"
)

type legacyDroplet struct {
	port int
}

func New(port int) droplet.Droplet {
	return &legacyDroplet{port: port}
}

func (d *legacyDroplet) ListenAndServe() error {
	addr := fmt.Sprintf("localhost:%d", d.port)

	http.HandleFunc("/", d.index)
	http.HandleFunc("/command/", d.HandleCommand)

	log.Printf("Listening on http://%s", addr)

	return http.ListenAndServe(addr, nil)
}

func (d *legacyDroplet) index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}

func (d *legacyDroplet) HandleCommand(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s! (from legacy)", r.URL.Path[1:])
}

func (d *legacyDroplet) Start(instances []droplet.PassivbotInstance) error {
	return nil
}

func (d *legacyDroplet) Stop(instances []droplet.PassivbotInstance) error {
	return nil
}

func (d *legacyDroplet) Restart(instances []droplet.PassivbotInstance) error {
	return nil
}
