package services

import (
	"case/models"
	"case/repositories"
)

type SongService struct {
	repo   *repositories.SongRepository
	apiURL string
}

func NewSongService(repo *repositories.SongRepository, apiURL string) *SongService {
	return &SongService{repo: repo, apiURL: apiURL}
}

func (s *SongService) GetSongs(filter map[string]string, page, limit int) ([]models.Song, error) {
	return s.repo.GetSongs(filter, page, limit)
}

func (s *SongService) GetSongLyrics(id, page, limit int) (string, error) {
	return s.repo.GetSongLyrics(id, page, limit)
}

func (s *SongService) DeleteSong(id int) error {
	return s.repo.DeleteSong(id)
}

func (s *SongService) UpdateSong(song *models.Song) error {
	return s.repo.UpdateSong(song)
}

func (s *SongService) AddSong(song *models.Song) error {
	return s.repo.AddSong(song)
}

func (s *SongService) GetApiURL() string {
	return s.apiURL
}
