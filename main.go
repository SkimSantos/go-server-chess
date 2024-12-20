package main

import (
	"fmt"
	"go-http-server/routes"
	"net/http"
)

func main() {
	// Serve static files
	var fs = http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", routes.HomeHandler) // Define the route and handler
	http.HandleFunc("/about", routes.AboutHandler)
	http.HandleFunc("/chess", routes.ChessHandler)

	var port string = "8080"
	fmt.Printf("Starting server on port %s...\n", port)
	var err error = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
