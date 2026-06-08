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

// @Summary      Set Maintenance Status
// @Description  Sets whether infrastructure maintenance is active.
// @Tags         config
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      WartungResponse  true  "Wartung Status"
// @Success      200   {object}  WartungResponse
// @Failure      400   {object}  badrequest
// @Failure      401   {object}  ErrorResponse
// @Router       /set_wartung [post]
func REST_PostWartung(c *gin.Context) {
	var req WartungResponse
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige Anfrage"})
		return
	}

	types.Config.Wartung = req.Wartung

	if err := types.SaveConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Speichern Fehler"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wartung": types.Config.Wartung,
	})
}

// @Summary      Set Maintenance Sleeptime
// @Description  Sets the sleep time during maintenance.
// @Tags         config
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      WartungSleepResponse  20  "Sleeptime"
// @Success      200   {object}  WartungSleepResponse
// @Failure      400   {object}  badrequest
// @Failure      401   {object}  ErrorResponse
// @Router       /set_wartung [post]
func REST_PostWartungSleep(c *gin.Context) {
	var req WartungSleepResponse
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige Anfrage"})
		return
	}
	types.Config.Wartung_sleep_time = req.WartungSleepTime

	if err := types.SaveConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Speicher Fehler"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"wartung_sleep_time": types.Config.Sleep_time,
	})
}

type badrequest struct {
	Error string `json:"error" example:"Ungültige Anfrage"`
}

type WartungResponse struct {
	Wartung bool `json:"wartung" example:"false"`
}

type WartungSleepResponse struct {
	WartungSleepTime int `json:"wartung_sleep_time" example:"20"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Ungültiger Token"`
}
