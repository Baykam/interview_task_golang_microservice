package accountHandler

import (
	"encoding/json"
	"net/http"
)

func (h *handler) GetAccountsList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	accounts, err := h.service.GetAccountsList(r.Context())
	if err != nil {
		h.logger.Error("Handler error fetching accounts list: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch accounts")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accounts)
}
