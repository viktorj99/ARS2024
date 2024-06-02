package configuration

import (
	"log"
	"os"
)

type Configuration struct {
	JaegerAddress string
	ServerAddress string
}

func GetConfiguration() Configuration {
	jaegerAddress := os.Getenv("JAEGER_ADDRESS")
	if jaegerAddress == "" {
		log.Fatal("JAEGER_ADDRESS environment variable not set")
	}

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "0.0.0.0:8080"
	}

	return Configuration{
		JaegerAddress: jaegerAddress,
		ServerAddress: serverAddress,
	}
}
