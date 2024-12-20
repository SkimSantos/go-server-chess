package routes

import (
	"html/template"
	"net/http"
)

// AboutHandler renders the about.html page
func ChessHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/chess.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"Title": "Chess Page",
	}
	tmpl.Execute(w, data)
}
