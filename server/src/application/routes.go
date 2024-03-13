package application

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"opencv/src/handler"
)

func (app *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Route("/recognizer", app.loadRecognizerRoutes)

	app.router = router
}

func (app *App) loadRecognizerRoutes(router chi.Router) {
	recognizerHandler := &handler.Recognizer{}
	router.Post("/", recognizerHandler.RecognizeTwoPhoto)
	router.Put("/{id}", recognizerHandler.AddToRecognize)
	router.Post("/{id}", recognizerHandler.Recognize)
}
