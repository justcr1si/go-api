package handlers

import (
	"case/models"
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

// GetSongs @Summary Get a list of songs
// @Description Get a list of songs with optional filtering and pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Filter by group"
// @Param song query string false "Filter by song name"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} gin.H "List of songs"
// @Failure 500 {object} gin.H "Failed to get songs"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	songs, err := h.service.GetSongs(filter, page, limit)
	if err != nil {
		h.log.Errorf("Failed to get songs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get songs"})
		return
	}

	h.log.WithFields(logrus.Fields{
		"count": len(songs),
	}).Info("Songs retrieved successfully")

	c.JSON(http.StatusOK, gin.H{"songs": songs})
}

// GetSongLyrics @Summary Get a lyrics of song
// @Description Get a lyrics of song by its song's id
// @Tags songs
// @Accept json
// @Produce json
// @Param song query string false "Filter by song name"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} gin.H "Lyrics of song"
// @Failure 500 {object} gin.H "Failed to get lyrics of the song"
// @Router /songs [get]
func (h *SongHandler) GetSongLyrics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	lyrics, err := h.service.GetSongLyrics(id, page, limit)
	if err != nil {
		h.log.Errorf("Failed to get lyrics: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get lyrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lyrics": lyrics})
}

// DeleteSong @Summary Delete a song
// @Description Delete a song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} gin.H "Song deleted"
// @Failure 400 {object} gin.H "Invalid ID"
// @Failure 500 {object} gin.H "Failed to delete song"
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.service.DeleteSong(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted"})
}

// UpdateSong @Summary Update a song
// @Description Update an existing song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Updated song data"
// @Success 200 {object} gin.H "Song updated"
// @Failure 400 {object} gin.H "Invalid request body or ID"
// @Failure 500 {object} gin.H "Failed to update song"
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song updated", "song": song})
}

// AddSong @Summary Add a new song
// @Description Add a new song to the database
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Song data"
// @Success 200 {object} gin.H "Song added"
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 500 {object} gin.H "Failed to add song"
// @Router /songs [post]
func (h *SongHandler) AddSong(c *gin.Context) {
	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if song.Group == "" || song.Song == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields 'group' and 'song' are required"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song added", "song": song})
}
