package configuration

import (
	"log"
	"os"
)

type Configuration struct {
	JaegerEndpoint string
	ServerAddress  string
}

func GetConfiguration() Configuration {
	jaegerEndpoint := os.Getenv("JAEGER_ENDPOINT")
	if jaegerEndpoint == "" {
		log.Fatal("JAEGER_ENDPOINT environment variable not set")
	}

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = ":8080"
	}

	return Configuration{
		JaegerEndpoint: jaegerEndpoint,
		ServerAddress:  serverAddress,
	}
}
