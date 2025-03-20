package repositories

import (
	"case/models"
	"database/sql"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type SongRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewSongRepository(db *sql.DB, log *logrus.Logger) *SongRepository {
	return &SongRepository{db: db, log: log}
}

func (r *SongRepository) GetSongs(filter map[string]string, page, limit int) ([]models.Song, error) {
	query := `SELECT id, "group", song, release_date, text, link FROM songs`
	var args []interface{}

	if group, ok := filter["group"]; ok {
		query += " AND \"group\" = $" + strconv.Itoa(len(args)+1)
		args = append(args, group)
	}

	if song, ok := filter["song"]; ok {
		query += " AND song = $" + strconv.Itoa(len(args)+1)
		args = append(args, song)
	}

	query += " LIMIT $" + strconv.Itoa(len(args)+1)
	args = append(args, limit)
	query += " OFFSET $" + strconv.Itoa(len(args)+1)
	args = append(args, (page-1)*limit)

	r.log.WithFields(logrus.Fields{
		"query": query,
		"group": filter["group"],
		"song":  filter["song"],
	}).Debug("Executing SQL query")

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.log.WithFields(logrus.Fields{
				"error": err,
			}).Error("Error closing rows")
		}
	}(rows)

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func (r *SongRepository) GetSongLyrics(id, page, limit int) (string, error) {
	var text string

	query := "SELECT text FROM songs WHERE id=$1"

	r.log.WithFields(logrus.Fields{
		"query": query,
	}).Debug("Executing SQL query")

	err := r.db.QueryRow(query, id).Scan(&text)
	if err != nil {
		return "", err
	}

	verses := strings.Split(text, "\n\n")
	start := (page - 1) * limit
	end := start + limit
	if start >= len(verses) {
		return "", nil
	}
	if end > len(verses) {
		end = len(verses)
	}

	return strings.Join(verses[start:end], "\n\n"), nil
}

func (r *SongRepository) DeleteSong(id int) error {
	query := "DELETE FROM songs WHERE id = $1"
	r.log.WithFields(logrus.Fields{
		"query": query,
	}).Debug("Executing SQL query")
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SongRepository) UpdateSong(song *models.Song) error {
	query := `UPDATE songs
SET "group" = $1, song = $2, release_date = $3, text = $4, link = $5
WHERE id = $6
`
	r.log.WithFields(logrus.Fields{
		"query": query,
	}).Debug("Executing SQL query")
	_, err := r.db.Exec(query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link, song.ID)
	return err
}

func (r *SongRepository) AddSong(song *models.Song) error {
	query := `
        INSERT INTO songs ("group", song, release_date, text, link)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	r.log.WithFields(logrus.Fields{
		"query": query,
	}).Debug("Executing SQL query")

	return r.db.QueryRow(query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link).Scan(&song.ID)
}
