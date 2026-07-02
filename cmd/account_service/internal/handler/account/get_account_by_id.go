package accountHandler

import (
	"encoding/json"
	"net/http"
)

func (h *handler) GetAccountById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, "Account ID is required")
		return
	}

	account, err := h.service.GetAccount(r.Context(), id)
	if err != nil {
		if err.Error() == "account not found" {
			h.respondWithError(w, http.StatusNotFound, "Account not found")
			return
		}
		h.logger.Error("Handler error fetching account %s: %v", id, err)
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}
