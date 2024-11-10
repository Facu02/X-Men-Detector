package router

import (
	"x-menDetector/pkg/controller"
	"x-menDetector/pkg/repository"
	"x-menDetector/pkg/services"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func RegisterMutantRouter(router *mux.Router, redis *redis.Client) {
	mutantRepository := repository.NewMutanRepository(redis)
	mutantServices := services.NewMutantServices(mutantRepository)
	mutantController := controller.NewMutantController(mutantServices)
	router.HandleFunc("/mutant", mutantController.HandleMutant).Methods("POST")
	router.HandleFunc("/stats", mutantController.HandleStats).Methods("GET")
}
