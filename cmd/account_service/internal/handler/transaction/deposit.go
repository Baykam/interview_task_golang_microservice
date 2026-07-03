package transActionHandler

import (
	"encoding/json"
	"interview_task_golang_microservices/dto"
	"interview_task_golang_microservices/models"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"
	"net/http"
)

func (h *handler) Deposit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, "Account ID is required")
		return
	}

	var req models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body := dto.TransActionToByte(req)
	if body == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err := h.publisher.PublishMessage(r.Context(), rabbitmq.Message{
		QueueName: h.cfg.TransactionService.Queues.WithdrawQueue,
		Body:      body,
	})
	if err != nil {
		h.logger.Error("Failed to publish deposit event: %v", err)
		h.respondWithError(w, 404, "Failed to publish deposit")
		return
	}

	w.WriteHeader(http.StatusAccepted) // 202 Accepted
	json.NewEncoder(w).Encode(map[string]string{"status": "transaction queued"})
}
