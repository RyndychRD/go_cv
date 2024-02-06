package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb *redis.Client
}

func New() *App {
	app := &App{
		router: loadRoutes(),
		rdb: redis.NewClient(&redis.Options{}),
	}

	return app
}

func (this *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: this.router,
	}
	err:=this.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	
	defer func(){
		if err:= this.rdb.Close(); err!= nil {
			fmt.Println("failed to close redis",err)
		}
	}()

	ch :=make(chan error,1)

	go func(){
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select{
		case <-ctx.Done():
			timeout, cancel:=context.WithTimeout(context.Background(),time.Second*10)
			defer cancel()
			return server.Shutdown(timeout)
		case err=<-ch:
			return err
	}
}