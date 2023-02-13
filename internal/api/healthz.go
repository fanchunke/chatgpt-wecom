package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *router) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "healthz"})
}
