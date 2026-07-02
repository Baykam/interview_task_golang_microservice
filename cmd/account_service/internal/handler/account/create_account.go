package accountHandler

import (
	"encoding/json"
	"interview_task_golang_microservices/models"
	"net/http"
)

func (h *handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var account models.Account
	// İstemciden gelen JSON gövdesini (body) parse ediyoruz
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		h.logger.Info("Failed to decode create account request: %v", err)
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Servis katmanını çağırıyoruz (O arkada DB ve Cache işlemlerini hallediyor)
	if err := h.service.CreateAccount(r.Context(), &account); err != nil {
		h.logger.Error("Handler error creating account: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}
