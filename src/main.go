package main

import (
	"log"
	"net/http"
)

func main() {
	// Create a file server handler pointing to your local directory
	fileServer := http.FileServer(http.Dir("./aws"))

	// Strip the URL prefix so the file server looks for files correctly
	http.Handle("/aws/", http.StripPrefix("/aws/", fileServer))

	http.HandleFunc("/aws", func(w http.ResponseWriter, r *http.Request) {
		// Target URL can be relative to domain ("/new-path") or absolute ("https://example.com")
		// Redirects the client temporarily to a new URL using HTTP 302 via http.StatusFound
		http.Redirect(w, r, "/aws/", http.StatusFound)
	})

	// Default path redirects to web page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/aws/", http.StatusFound)
	})

	log.Println("Server starting on :8888...")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}
}
