package application

import (
	"os"
	"strconv"
)

type Config struct {
	RedisAddress string
	ServerPort   uint64
}

func LoadConfig() Config {
	cfg := Config{
		RedisAddress: "localhost:6379",
		ServerPort:   3000,
	}

	if redisArr, exists := os.LookupEnv("REDIS_ADDR"); exists {
		cfg.RedisAddress = redisArr
	}
	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = port
		}
	}
	return cfg
}
