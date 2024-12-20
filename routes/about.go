package routes

import (
	"html/template"
	"net/http"
)

// AboutHandler renders the about.html page
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/about.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"Title": "About Page",
	}
	tmpl.Execute(w, data)
}
