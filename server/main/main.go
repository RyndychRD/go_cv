package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"opencv/main/application"
	"opencv/main/recognizer"
	"os"
	"os/signal"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found, using default values")
	}
	recognizer.InitEnv()
}

func main() {
	app := application.New(application.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app:", err)
	}

}
