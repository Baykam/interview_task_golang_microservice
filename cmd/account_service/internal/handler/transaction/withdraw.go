package transActionHandler

import (
	"encoding/json"
	"interview_task_golang_microservices/models"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"
	"net/http"
)

func (h *handler) Withdraw(w http.ResponseWriter, r *http.Request) {
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

	body, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.publisher.PublishMessage(r.Context(), rabbitmq.Message{
		QueueName: h.cfg.TransactionService.Queues.WithdrawQueue,
		Body:      body,
	})

	if err != nil {
		h.logger.Error("Failed to publish withdraw event: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to process transaction")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "transaction queued"})
}
