package main

import (
	"net/http"

	"github.com/nightfaust/snippetbox/ui"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/ping", ping)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.about))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignUp))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignUpPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/account/view", protected.ThenFunc(app.userAccount))
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}

// If you want to restrict access to static files
//
//type neuteredFileSystem struct {
//	fs http.FileSystem
//}
//
//func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
//	f, err := nfs.fs.Open(path)
//	if err != nil {
//		return nil, err
//	}
//
//	s, _ := f.Stat()
//	if s.IsDir() {
//		index := filepath.Join(path, "index.html")
//		if _, err := nfs.fs.Open(index); err != nil {
//			closeErr := f.Close()
//			if closeErr != nil {
//				return nil, closeErr
//			}
//
//			return nil, err
//		}
//	}
//
//	return f, nil
//}
