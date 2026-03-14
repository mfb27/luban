package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mfb27/luban/internal/model"
)

func (a *App) getModels(c *gin.Context) {
	var models []model.Model
	if err := a.db.Order("updated_at desc").Find(&models).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models)
}

