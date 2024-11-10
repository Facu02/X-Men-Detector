package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	mocks "x-menDetector/internal/mocks/services"
	"x-menDetector/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestHandleMutant(t *testing.T) {
	mockService := new(mocks.MockMutanServices)
	controller := NewMutantController(mockService)

	t.Run("Valid DNA - is mutant", func(t *testing.T) {
		dna := []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"}
		mockService.On("IsMutant", dna).Return(true, nil)

		body, _ := json.Marshal(map[string]interface{}{"dna": dna})
		req := httptest.NewRequest("POST", "/mutant", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.HandleMutant(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]bool
		json.NewDecoder(rec.Body).Decode(&response)
		assert.True(t, response["is_mutant"])

		mockService.AssertExpectations(t)
	})

	t.Run("Valid DNA - not a mutant", func(t *testing.T) {
		dna := []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCTTA", "TCACTG"}
		mockService.On("IsMutant", dna).Return(false, nil)

		body, _ := json.Marshal(map[string]interface{}{"dna": dna})
		req := httptest.NewRequest("POST", "/mutant", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		controller.HandleMutant(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]bool
		json.NewDecoder(rec.Body).Decode(&response)
		assert.False(t, response["is_mutant"])

		mockService.AssertExpectations(t)
	})

}

func TestHandleMutantError(t *testing.T) {
	mockService := new(mocks.MockMutanServices)
	controller := NewMutantController(mockService)

	t.Run("Invalid DNA format", func(t *testing.T) {
		body := []byte(`{"dna": "invalid_dna"}`)
		req := httptest.NewRequest("POST", "/mutant", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		controller.HandleMutant(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Tipo de dato inválido")
	})

	t.Run("Error_in_IsMutant", func(t *testing.T) {
		dna := []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"}

		mockService.On("IsMutant", dna).Return(false, assert.AnError)

		reqBody := map[string][]string{"dna": dna}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/mutant", bytes.NewBuffer(reqJSON))
		rr := httptest.NewRecorder()

		controller.HandleMutant(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		assert.Contains(t, rr.Body.String(), "error al obtener validar el dna")
	})
}

func TestHandleStats(t *testing.T) {
	mockService := new(mocks.MockMutanServices)
	controller := NewMutantController(mockService)

	t.Run("error_in_GetMutantStats", func(t *testing.T) {
		mockService.On("GetMutantStats").Return(models.MutantStats{}, assert.AnError) 

		req := httptest.NewRequest("GET", "/stats", nil)
		rr := httptest.NewRecorder()

		controller.HandleStats(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		assert.Contains(t, rr.Body.String(), "error al obtener las estadisticas")
	})
}
