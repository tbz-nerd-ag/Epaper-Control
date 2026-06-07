package rest

import (
	"Control/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary      Maintenance Status
// @Description  Returns a boolean indicating whether infrastructure maintenance is active.
// @Tags         config
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  WartungResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /get_wartung [get]
func REST_GetWartung(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"wartung": types.Config.Wartung,
	})
}

// @Summary      Maintenance Sleep Time
// @Description  Returns a int indicating how long a display stays in sleep mode during maintenance.
// @Tags         config
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  WartungSleepResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /get_wartung_sleep [get]
func REST_GetWartungSleepTime(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"wartung": types.Config.Wartung_sleep_time,
	})
}

func REST_SetWartung() {}

type WartungResponse struct {
	Wartung bool `json:"wartung" example:"false"`
}

type WartungSleepResponse struct {
	WartungSleepTime int `json:"wartung_sleep_time" example:"20"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Ungültiger Token"`
}
