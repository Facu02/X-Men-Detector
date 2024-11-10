package server

import (
	"log"
	"net/http"
	"os"
	"x-menDetector/pkg/config"
	"x-menDetector/pkg/router"

	"github.com/gorilla/mux"
)

func createServer(port string) error {
	log.Printf("servidor iniciado en el puerto %s", port)
	return http.ListenAndServe(":"+port, nil)
}

func Init() {
	registerRouters()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	createServer(port)
}

func registerRouters() {
	redis, err := config.InitRedis()
	if err != nil {
		log.Fatalf("error al iniciar la base de datos: %v", err)
	}

	r := mux.NewRouter()
	router.RegisterMutantRouter(r, redis)
	http.Handle("/", r)
}
