package routes

import (
	"fmt"
	"net/http"
)

func APIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Testing API Page Of Server")
}
