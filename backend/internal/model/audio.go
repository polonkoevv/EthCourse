package model

import (
	"time"
)

type Audio struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Artist     string    `json:"artist"`
	IPFSCID    string    `json:"ipfsCid"`
	OwnerAddr  string    `json:"ownerAddress"`
	Signature  string    `json:"signature"`
	UploadedAt time.Time `json:"uploadedAt"`
}
