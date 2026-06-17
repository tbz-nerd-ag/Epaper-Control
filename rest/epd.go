package rest

import (
	"Control/types"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// @Summary      Get Device by ID
// @Description  Returns the device name for a given EPD ID.
// @Tags         epd
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "EPD ID"
// @Success      200  {object}  DeviceResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /get_device/{id} [get]]

func REST_GetDeviceViaID(c *gin.Context) {
	id := c.Param("id")

	// epd.json einlesen
	data, err := os.ReadFile("epd.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Fehler beim Lesen der Konfiguration",
		})
		return
	}

	var config struct {
		EPD []types.Epd `json:"epd"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Fehler beim Parsen der Konfiguration",
		})
		return
	}

	// ID suchen
	for _, entry := range config.EPD {
		if entry.ID == id {
			c.JSON(http.StatusOK, gin.H{
				"device": entry.Device,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "Gerät nicht gefunden",
	})
}
