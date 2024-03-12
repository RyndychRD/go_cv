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
	modelPath, mExists := os.LookupEnv("MODEL_FULL_PATH")
	if mExists && modelPath != "" {
		recognizer.ModelsDir = modelPath
	}
	storagePath, sExists := os.LookupEnv("STORAGE_FULL_PATH")
	if sExists && storagePath != "" {
		recognizer.ImagesDir = storagePath
	}
	isTestThreshold, tExists := os.LookupEnv("IS_TEST_THRESHOLD")
	if tExists && isTestThreshold == "1" {
		recognizer.IsTestThreshold = true
	}
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
