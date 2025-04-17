package model

import "time"

type BlockchainTransaction struct {
	Hash            string
	BlockNumber     uint64
	From            string
	To              string
	Value           string
	Gas             uint64
	GasPrice        string
	Input           string
	Timestamp       time.Time
	TransactionType string // "incoming" или "outgoing"
}
