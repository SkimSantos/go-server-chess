package routes

import (
	"html/template"
	"net/http"
)

// HomeHandler renders the home.html page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/home.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"Title": "Home Page",
	}
	tmpl.Execute(w, data)
}
