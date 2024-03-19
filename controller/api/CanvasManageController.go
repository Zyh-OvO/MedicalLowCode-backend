package api

import (
	"MedicalLowCode-backend/model"
	"MedicalLowCode-backend/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CanvasManageController struct{}

type editCanvasJson struct {
	ProjectId     string `json:"projectId" binding:"required"`
	CanvasContent string `json:"canvasContent" binding:"required"`
}

type getCanvasInfoJson struct {
	ProjectId string `json:"projectId" binding:"required"`
}

func (CanvasManageController) EditCanvas(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json editCanvasJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectId, _ := strconv.Atoi(json.ProjectId)
	canvas := model.EditCanvas(token.UserId, projectId, json.CanvasContent)
	if canvas == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "项目不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"projectId":     strconv.Itoa(canvas.ProjectId),
		"canvasContent": canvas.CanvasContent,
		"createdAt":     canvas.CreatedAt.Unix(),
		"updatedAt":     canvas.UpdatedAt.Unix(),
	})
}

func (CanvasManageController) GetCanvasInfo(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json getCanvasInfoJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectId, _ := strconv.Atoi(json.ProjectId)
	canvas := model.QueryCanvas(token.UserId, projectId)
	if canvas == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "项目不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"projectId":     strconv.Itoa(canvas.ProjectId),
		"canvasContent": canvas.CanvasContent,
		"createdAt":     canvas.CreatedAt.Unix(),
		"updatedAt":     canvas.UpdatedAt.Unix(),
	})
}
