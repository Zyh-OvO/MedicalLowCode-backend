package api

import (
	"MedicalLowCode-backend/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ModuleManageController struct{}

type addPersonalModuleJson struct {
	ModuleName    string `json:"moduleName" binding:"required"`
	ModuleContent string `json:"moduleContent" binding:"required"`
}

type deletePersonalModuleJson struct {
	ModuleId string `json:"moduleId" binding:"required"`
}

type editPersonalModuleJson struct {
	ModuleId      string  `json:"moduleId" binding:"required"`
	ModuleName    *string `json:"moduleName"`
	ModuleContent *string `json:"moduleContent"`
}

func (m ModuleManageController) GetPersonalModules(c *gin.Context) {
	token := c.MustGet("token").(*model.Token)
	modules := model.QueryPersonalModules(token.UserId)
	var moduleList []gin.H
	for _, module := range modules {
		moduleList = append(moduleList, gin.H{
			"moduleId":      strconv.Itoa(module.ModuleId),
			"moduleName":    module.ModuleName,
			"moduleContent": module.ModuleContent,
			"createdAt":     module.CreatedAt.Unix(),
			"updatedAt":     module.UpdatedAt.Unix(),
		})
	}
	response := gin.H{
		"modules": moduleList,
	}
	c.JSON(200, response)
}

func (m ModuleManageController) AddPersonalModule(c *gin.Context) {
	token := c.MustGet("token").(*model.Token)
	var json addPersonalModuleJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newModule := model.NewModule(token.UserId, json.ModuleName, json.ModuleContent)
	c.JSON(http.StatusOK, gin.H{
		"moduleId":      strconv.Itoa(newModule.ModuleId),
		"moduleName":    newModule.ModuleName,
		"moduleContent": newModule.ModuleContent,
		"createdAt":     newModule.CreatedAt.Unix(),
		"updatedAt":     newModule.UpdatedAt.Unix(),
	})
}

func (m ModuleManageController) DeletePersonalModule(c *gin.Context) {
	token := c.MustGet("token").(*model.Token)
	var json deletePersonalModuleJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	moduleId, _ := strconv.Atoi(json.ModuleId)
	if model.QueryPersonalModule(token.UserId, moduleId) == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "组件不存在"})
		return
	}
	model.DeleteModule(token.UserId, moduleId)
	c.JSON(http.StatusOK, gin.H{})
}

func (m ModuleManageController) EditPersonalModule(c *gin.Context) {
	token := c.MustGet("token").(*model.Token)
	var json editPersonalModuleJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	moduleId, _ := strconv.Atoi(json.ModuleId)
	module := model.QueryPersonalModule(token.UserId, moduleId)
	if module == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "组件不存在"})
		return
	}
	if json.ModuleName == nil && json.ModuleContent == nil {
		c.JSON(http.StatusOK, gin.H{
			"moduleId":      strconv.Itoa(module.ModuleId),
			"moduleName":    module.ModuleName,
			"moduleContent": module.ModuleContent,
			"createdAt":     module.CreatedAt.Unix(),
			"updatedAt":     module.UpdatedAt.Unix(),
		})
	} else {
		editedModule := model.EditModule(token.UserId, moduleId, json.ModuleName, json.ModuleContent)
		c.JSON(http.StatusOK, gin.H{
			"moduleId":      strconv.Itoa(editedModule.ModuleId),
			"moduleName":    editedModule.ModuleName,
			"moduleContent": editedModule.ModuleContent,
			"createdAt":     editedModule.CreatedAt.Unix(),
			"updatedAt":     editedModule.UpdatedAt.Unix(),
		})
	}
}
