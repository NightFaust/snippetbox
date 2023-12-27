package main

import "net/http"

func (app *application) routes(staticDir string) *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir(staticDir)})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}
