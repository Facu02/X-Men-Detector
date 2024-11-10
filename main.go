package main

import (
	"log"
	"x-menDetector/pkg/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}
	server.Init()
}
