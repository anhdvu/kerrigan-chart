package handler

import (
	"html/template"
	"log"
	"net/http"
)

func AulabIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/aulab" {
		http.NotFound(w, r)
		return
	}

	// Include the footer partial in the template files.
	files := []string{
		"./ui/html/home.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
