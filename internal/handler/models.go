package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mfb27/luban/internal/model"
	"github.com/mfb27/luban/internal/response"
)

func (a *App) getModels(c *gin.Context) {
	var models []model.Model
	if err := a.db.Where("status = ?", "active").Order("updated_at desc").Find(&models).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, err.Error())
		return
	}
	response.NewResponseHelper(c).Success(models)
}

