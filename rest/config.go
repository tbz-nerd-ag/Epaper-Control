package rest

import (
	"Control/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func REST_GetWartung(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"wartung": types.Config.Wartung,
	})
}

func REST_SetWartung() {}
