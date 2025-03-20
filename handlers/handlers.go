package handlers

import (
	"case/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"case/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SongHandler struct {
	service *services.SongService
	log     *logrus.Logger
}

func NewSongHandler(service *services.SongService, log *logrus.Logger) *SongHandler {
	return &SongHandler{service: service, log: log}
}

// GetSongs
// @Summary Get a list of songs
// @Description Get a list of songs with optional filtering and pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Filter by group"
// @Param song query string false "Filter by song name"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} models.SongListResponse "List of songs"
// @Failure 500 {object} models.ErrorResponse "Failed to get songs"
// @Router /songs [get]
func (h *SongHandler) GetSongs(c *gin.Context) {
	filter := make(map[string]string)
	if group := c.Query("group"); group != "" {
		filter["group"] = group
	}
	if song := c.Query("song"); song != "" {
		filter["song"] = song
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: fmt.Sprintf("page must be a positive integer. Got %v", page)})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid limit number"})
		return
	}

	songs, err := h.service.GetSongs(filter, page, limit)
	if err != nil {
		h.log.Errorf("Failed to get songs: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get songs"})
		return
	}

	h.log.WithFields(logrus.Fields{
		"count": len(songs),
	}).Info("Songs retrieved successfully")

	c.JSON(http.StatusOK, models.SongListResponse{Songs: models.ToSongResponseList(songs)})
}

// GetSongLyrics
// @Summary Get a lyrics of song
// @Description Get a lyrics of song by its song's id
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} models.MessageResponse "Lyrics of song"
// @Failure 500 {object} models.ErrorResponse "Failed to get lyrics of the song"
// @Router /songs/{id}/lyrics [get]
func (h *SongHandler) GetSongLyrics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid id"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: fmt.Sprintf("Invalid limit %d", limit)})
		return
	}

	lyrics, err := h.service.GetSongLyrics(id, page, limit)
	if err != nil {
		h.log.Errorf("Failed to get lyrics: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get lyrics"})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: lyrics})
}

// DeleteSong
// @Summary Delete a song
// @Description Delete a song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} models.MessageResponse "Song deleted"
// @Failure 400 {object} models.ErrorResponse "Invalid ID"
// @Failure 500 {object} models.ErrorResponse "Failed to delete song"
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid id"})
		return
	}

	err = h.service.DeleteSong(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete song"})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{Message: "song deleted"})
}

// UpdateSong
// @Summary Update a song
// @Description Update an existing song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Updated song data"
// @Success 200 {object} models.SongResponse "Song updated"
// @Failure 400 {object} models.ErrorResponse "Invalid request body or ID"
// @Failure 500 {object} models.ErrorResponse "Failed to update song"
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid id"})
		return
	}

	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid song"})
		return
	}
	song.ID = id

	if song.ReleaseDate == "" {
		song.ReleaseDate = time.Now().Format("02.01.2006")
	}

	if err := h.service.UpdateSong(&song); err != nil {
		h.log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to update song")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update song"})
		return
	}

	songResponse := models.ToSongResponse(song)

	c.JSON(http.StatusOK, models.SongResponse{
		ID:          songResponse.ID,
		Group:       songResponse.Group,
		Song:        songResponse.Song,
		ReleaseDate: songResponse.ReleaseDate,
		Text:        songResponse.Text,
		Link:        songResponse.Link,
	})
}

// AddSong
// @Summary Add a new song
// @Description Add a new song to the database
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Song data"
// @Success 200 {object} models.SongResponse "Song added"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 500 {object} models.ErrorResponse "Failed to add song"
// @Router /songs [post]
func (h *SongHandler) AddSong(c *gin.Context) {
	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid song"})
		return
	}

	if song.Group == "" || song.Song == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid song"})
		return
	}

	if song.ReleaseDate == "" {
		song.ReleaseDate = time.Now().Format("02.01.2006")
	}
	if song.Text == "" {
		song.Text = ""
	}
	if song.Link == "" {
		song.Link = ""
	}

	if err := h.service.AddSong(&song); err != nil {
		h.log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to add song")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to add song"})
		return
	}

	songResponse := models.ToSongResponse(song)

	c.JSON(http.StatusOK, models.SongResponse{
		ID:          songResponse.ID,
		Group:       songResponse.Group,
		Song:        songResponse.Song,
		ReleaseDate: songResponse.ReleaseDate,
		Text:        songResponse.Text,
		Link:        songResponse.Link,
	})
}
