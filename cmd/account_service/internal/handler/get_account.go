package accountHandler

import (
	"encoding/json"
	rabbitmq "interview_task_golang_microservices/pkgs/rabbit_mq"
	"net/http"
)

func (h *handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	ff, err := json.Marshal("something")
	if err != nil {
		return
	}

	h.publisher.PublishMessage(r.Context(), rabbitmq.Message{
		QueueName: h.cfg.AccountService.Queues.DepositQueue,
		Body:      ff,
	})

	// data, err := h.grpcClient.CheckAccountExists(r.Context(), &accountProto.CheckAccountExistsRequest{})
	// if err != nil {
	// 	return
	// }

	// print(data)
}
