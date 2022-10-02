package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.html")
	})

	fmt.Println("Starting front end service on port 82")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Panic(err)
	}
}

//go:embed templates
var templateFS embed.FS

func render(w http.ResponseWriter, t string) {

	partials := []string{
		"templates/base.layout.html",
		"templates/header.partial.html",
		"templates/footer.partial.html",
		// "./cmd/web/templates/base.layout.html",
		// "./cmd/web/templates/header.partial.html",
		// "./cmd/web/templates/footer.partial.html",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("templates/%s", t))
	// templateSlice = append(templateSlice, fmt.Sprintf("./cmd/web/templates/%s", t))

	templateSlice = append(templateSlice, partials...)

	// tmpl, err := template.ParseFiles(templateSlice...)
	tmpl, err := template.ParseFS(templateFS, templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := &struct {
		BrokerURL string
	}{ 
		BrokerURL: os.Getenv("BROKER_URL"),
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
