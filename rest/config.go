package rest

import (
	"Control/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary		Maintenance Status
// @Description	Returns a boolean indicating whether infrastructure maintenance is active.
// @Tags         config
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /get_wartung [get]
func REST_GetWartung(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"wartung": types.Config.Wartung,
	})
}

func REST_SetWartung() {}
