package api

import (
	"MedicalLowCode-backend/model"
	"MedicalLowCode-backend/util"
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
	ProjectId string `json:"projectId" binding:"required"`
}

type editProjectJson struct {
	ProjectId          string  `json:"projectId" binding:"required"`
	ProjectName        *string `json:"projectName"`
	ProjectDescription *string `json:"projectDescription"`
}

type getProjectInfoJson struct {
	ProjectId string `json:"projectId" binding:"required"`
}

func (p ProjectManageController) NewProject(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
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
	token := c.MustGet("token").(*util.Token)
	var json deleteProjectJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectId, _ := strconv.Atoi(json.ProjectId)
	if model.QueryProject(token.UserId, projectId) == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "项目不存在"})
		return
	}
	model.DeleteProject(token.UserId, projectId)
	c.JSON(http.StatusOK, gin.H{})
}

func (p ProjectManageController) EditProject(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json editProjectJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectId, _ := strconv.Atoi(json.ProjectId)
	project := model.QueryProject(token.UserId, projectId)
	if project == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "项目不存在"})
		return
	}
	if json.ProjectName == nil && json.ProjectDescription == nil {
		c.JSON(http.StatusOK, gin.H{
			"projectId":          strconv.Itoa(project.ProjectId),
			"projectName":        project.ProjectName,
			"projectDescription": project.ProjectDescription,
			"createdAt":          project.CreatedAt.Unix(),
			"updatedAt":          project.UpdatedAt.Unix(),
		})
	} else {
		editedProject := model.EditProject(token.UserId, projectId, json.ProjectName, json.ProjectDescription)
		c.JSON(http.StatusOK, gin.H{
			"projectId":          strconv.Itoa(editedProject.ProjectId),
			"projectName":        editedProject.ProjectName,
			"projectDescription": editedProject.ProjectDescription,
			"createdAt":          editedProject.CreatedAt.Unix(),
			"updatedAt":          editedProject.UpdatedAt.Unix(),
		})
	}
}

func (p ProjectManageController) GetProjectInfo(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json getProjectInfoJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectId, _ := strconv.Atoi(json.ProjectId)
	project := model.QueryProject(token.UserId, projectId)
	if project == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "项目不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"projectId":          strconv.Itoa(project.ProjectId),
		"projectName":        project.ProjectName,
		"projectDescription": project.ProjectDescription,
		"createdAt":          project.CreatedAt.Unix(),
		"updatedAt":          project.UpdatedAt.Unix(),
	})
}

func (p ProjectManageController) GetProjectList(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	projectList := model.QueryProjectList(token.UserId)
	var projectInfoList []gin.H
	for _, project := range projectList {
		projectInfoList = append(projectInfoList, gin.H{
			"projectId":          strconv.Itoa(project.ProjectId),
			"projectName":        project.ProjectName,
			"projectDescription": project.ProjectDescription,
			"createdAt":          project.CreatedAt.Unix(),
			"updatedAt":          project.UpdatedAt.Unix(),
		})
	}
	response := gin.H{
		"projectList": projectInfoList,
	}
	c.JSON(http.StatusOK, response)
}
