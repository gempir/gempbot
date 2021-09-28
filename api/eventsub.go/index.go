package eventsubgo

import (
	"fmt"
	"net/http"

	"github.com/gempir/gempbot/pkg/log"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Info(r.URL.Query())
	log.Info(r.Body)

	fmt.Fprintf(w, "Hello")
}
