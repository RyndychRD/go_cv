package application

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerPort uint64
}

func LoadConfig() Config {
	cfg := Config{
		ServerPort: 3000,
	}

	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = port
		} else {
			log.Println("bad SERVER_PORT, using default ", serverPort)
		}
	} else {
		log.Println("bad SERVER_PORT, using default ", serverPort)
	}
	return cfg
}
