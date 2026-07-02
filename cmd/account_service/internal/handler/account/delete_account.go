package accountHandler

import (
	"net/http"
)

func (h *handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, "Account ID is required")
		return
	}

	// Servis katmanında Delete metodu olduğunu varsayarak çağırıyoruz
	// (Hizmet katmanında hem DB'den soft delete yapmalı hem de Cache'den DEL etmeli)
	if err := h.service.DeleteAccount(r.Context(), id); err != nil {
		h.logger.Error("Handler error deleting account %s: %v", id, err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to delete account")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
