package api

import (
	"MedicalLowCode-backend/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProjectManageController struct{}

type newProjectJson struct {
	ProjectName        string `json:"projectName" binding:"required"`
	ProjectDescription string `json:"projectDescription" binding:"required"`
}

type deleteProjectJson struct {
	ProjectId string `json:"projectId"`
}

type editProjectJson struct {
	ProjectId          string `json:"projectId"`
	ProjectName        string `json:"projectName"`
	ProjectDescription string `json:"projectDescription"`
}

type getProjectInfoJson struct {
	ProjectId string `json:"projectId"`
}

func (p ProjectManageController) NewProject(c *gin.Context) {
	token := c.MustGet("token").(*model.Token)
	var json newProjectJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newProject := model.NewProject(token.UserId, json.ProjectName, json.ProjectDescription)
	c.JSON(http.StatusOK, gin.H{
		"projectId":          strconv.Itoa(newProject.ProjectId),
		"projectName":        newProject.ProjectName,
		"projectDescription": newProject.ProjectDescription,
		"createdAt":          newProject.CreatedAt.Unix(),
		"updatedAt":          newProject.UpdatedAt.Unix(),
	})
}

func (p ProjectManageController) DeleteProject(c *gin.Context) {
	token := c.MustGet("token").(*model.Token)
	var json deleteProjectJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectId, _ := strconv.Atoi(json.ProjectId)
	if !model.CheckProjectPermission(token.UserId, projectId) {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "无权限"})
		return
	}
	model.DeleteProject(token.UserId, projectId)
	c.JSON(http.StatusOK, gin.H{})
}

func (p ProjectManageController) EditProject(c *gin.Context) {

}

func (p ProjectManageController) GetProjectInfo(c *gin.Context) {

}

func (p ProjectManageController) GetProjectList(c *gin.Context) {

}
