package api

import (
	"MedicalLowCode-backend/model"
	"MedicalLowCode-backend/util"
	"MedicalLowCode-backend/util/exportCode"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProjectDevelopController struct{}

type exportCodeJson struct {
	ProjectId string `json:"projectId" binding:"required"`
}

type submitTrainingTaskJson struct {
	ProjectId string `json:"projectId" binding:"required"`
}

type submitReasoningTaskJson struct {
	ProjectId string `json:"projectId" binding:"required"`
}

func (p ProjectDevelopController) ExportCode(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json exportCodeJson
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
	code := exportCode.ExportCode(canvas.CanvasContent)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
	})
}

func (p ProjectDevelopController) SubmitTrainingTask(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json submitTrainingTaskJson
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
	//todo

	c.JSON(http.StatusOK, gin.H{})
}

func (p ProjectDevelopController) SubmitReasoningTask(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json submitReasoningTaskJson
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
	//todo

	c.JSON(http.StatusOK, gin.H{})
}
