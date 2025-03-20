package models

type Song struct {
	ID          int    `db:"id" json:"id"`
	Group       string `db:"group" json:"group"`
	Song        string `db:"song" json:"song"`
	ReleaseDate string `db:"release_date" json:"release_date"`
	Text        string `db:"text" json:"text"`
	Link        string `db:"link" json:"link"`
}

type SongResponse struct {
	ID          int    `json:"id"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func ToSongResponse(song Song) SongResponse {
	return SongResponse{
		ID:          song.ID,
		Group:       song.Group,
		Song:        song.Song,
		ReleaseDate: song.ReleaseDate,
		Text:        song.Text,
		Link:        song.Link,
	}
}

func ToSongResponseList(songs []Song) []SongResponse {
	var response []SongResponse

	for _, song := range songs {
		response = append(response, ToSongResponse(song))
	}

	return response
}

type SongListResponse struct {
	Songs []SongResponse `json:"songs"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
