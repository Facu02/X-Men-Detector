package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"x-menDetector/internal/models"
	"x-menDetector/internal/utils"
)

type MutantController struct {
	mutantServices models.MutanServices
}

func NewMutantController(mutantServices models.MutanServices) *MutantController {
	return &MutantController{mutantServices: mutantServices}
}

func (mc *MutantController) HandleMutant(w http.ResponseWriter, r *http.Request) {
	var dnaRequest struct {
		DNA []string `json:"dna"`
	}

	if err := json.NewDecoder(r.Body).Decode(&dnaRequest); err != nil {
		http.Error(w, "tipo de dato invalido. Asegurese de enviar un JSON con una matriz de dna.", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateDNA(dnaRequest.DNA); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isMutant, err := mc.mutantServices.IsMutant(dnaRequest.DNA)
	if err != nil {
		http.Error(w, fmt.Sprintf("error al obtener validar el dna: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if !isMutant {
		http.Error(w, "acceso no autorizado", http.StatusForbidden)
		return
	}

	response := map[string]bool{"is_mutant": isMutant}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("error al procesar la respuesta: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (mc *MutantController) HandleStats(w http.ResponseWriter, r *http.Request) {
	stats, err := mc.mutantServices.GetMutantStats()
	if err != nil {
		http.Error(w, fmt.Sprintf("error al obtener las estadisticas: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, fmt.Sprintf("error al enviar la respuesta: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
