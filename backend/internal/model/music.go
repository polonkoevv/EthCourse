package model

type Music struct {
	ID    int    `json:"id" db:"music_id"`
	Title string `json:"title" db:"title"`
	CID   string `json:"cid" db:"cid"`
	Link  string `json:"link" db:"link"`
}
