package models

type Song struct {
	ID          int    `db:"id" json:"id"`
	Group       string `db:"group" json:"group"`
	Song        string `db:"song" json:"song"`
	ReleaseDate string `db:"release_date" json:"release_date"`
	Text        string `db:"text" json:"text"`
	Link        string `db:"link" json:"link"`
}
