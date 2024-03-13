package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"opencv/src/recognizer"
	"time"
)

type App struct {
	router http.Handler
	config Config
}

func New(config Config) *App {
	app := &App{
		config: config,
	}
	app.loadRoutes()

	return app
}

func (app *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.ServerPort),
		Handler: app.router,
	}

	ch := make(chan error, 1)

	go func() {
		log.Println("Server start")
		log.Printf("Threshold set to %f", recognizer.ClassificationThreshold)
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	case err := <-ch:
		return err
	}
}
