package dto

import (
	"encoding/json"
	"interview_task_golang_microservices/models"
)

func ByteToAccount(data []byte) (*models.Account, error) {
	var account models.Account
	err := json.Unmarshal(data, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func AccountToByte(account *models.Account) []byte {
	bb, err := json.Marshal(account)
	if err != nil {
		return nil
	}
	return bb
}
