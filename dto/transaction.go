package dto

import (
	"encoding/json"
	"interview_task_golang_microservices/models"
)

func ByteToTransAction(data []byte) *models.Transaction {
	var transaction models.Transaction
	err := json.Unmarshal(data, &transaction)
	if err != nil {
		return nil
	}
	return &transaction
}

func TransActionToByte(transAction models.Transaction) []byte {
	bb, err := json.Marshal(transAction)
	if err != nil {
		return nil
	}
	return bb
}
