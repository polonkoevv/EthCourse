package model

import "time"

type Music struct {
	ID         int       `json:"id" db:"music_id"`
	Title      string    `json:"title" db:"title"`
	CID        string    `json:"cid" db:"cid"`
	Link       string    `json:"link" db:"link"`
	OwnerAddr  string    `json:"owner_addr" db:"owner_addr"`
	Signature  string    `json:"signature" db:"signature"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}
