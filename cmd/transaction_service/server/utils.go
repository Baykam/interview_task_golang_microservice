package server

func (s *server) getConsumerGroupTopics() []string {
	return []string{
		s.cfg.TransactionService.Queues.DepositQueue,
		s.cfg.TransactionService.Queues.WithdrawQueue,
		s.cfg.TransactionService.Queues.TransferQueue,
	}
}
