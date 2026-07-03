package transActionHandler

import (
	"encoding/json"
	accountProto "interview_task_golang_microservices/protos"
	"net/http"
)

func (h *handler) GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.grpcClient.GetTransactionHistory(r.Context(), &accountProto.GetTransactionHistoryRequest{
		AccountId: id,
	})

	if err != nil {
		h.logger.Error("gRPC error fetching history: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch transaction history")
		return
	}

	json.NewEncoder(w).Encode(resp)
}
