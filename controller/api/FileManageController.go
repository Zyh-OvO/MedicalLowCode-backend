package api

import (
	"MedicalLowCode-backend/model"
	"MedicalLowCode-backend/util"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"strconv"
)

type FileManageController struct{}

type newDirectoryJson struct {
	ParentDirId string `json:"parentDirId" binding:"required"`
	DirName     string `json:"dirName" binding:"required"`
}

type deleteDirectoryJson struct {
	DirId string `json:"dirId" binding:"required"`
}

type renameDirectoryJson struct {
	DirId   string `json:"dirId" binding:"required"`
	DirName string `json:"dirName" binding:"required"`
}

type uploadFileForm struct {
	DirId string                `form:"dirId" binding:"required"`
	File  *multipart.FileHeader `form:"file" binding:"required"`
}

type deleteFileJson struct {
	FileId string `json:"fileId" binding:"required"`
}

type renameFileJson struct {
	FileId   string `json:"fileId" binding:"required"`
	FileName string `json:"fileName" binding:"required"`
}

func (f FileManageController) GetFileTree(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	fileTree := model.GetFileTree(token.UserId)
	c.JSON(http.StatusOK, fileTree)
}

func (f FileManageController) UploadFile(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var form uploadFileForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dirId, _ := strconv.Atoi(form.DirId)
	file, err := form.File.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()
	model.UploadFile(token.UserId, dirId, file, form.File.Filename)
	c.JSON(http.StatusOK, gin.H{})
}

func (f FileManageController) DeleteFile(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json deleteFileJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fileId, _ := strconv.Atoi(json.FileId)
	model.DeleteFile(token.UserId, fileId)
	c.JSON(http.StatusOK, gin.H{})
}

func (f FileManageController) RenameFile(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json renameFileJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fileId, _ := strconv.Atoi(json.FileId)
	model.RenameFile(token.UserId, fileId, json.FileName)
	c.JSON(http.StatusOK, gin.H{})
}

func (f FileManageController) NewDirectory(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json newDirectoryJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	parentDirId, _ := strconv.Atoi(json.ParentDirId)
	dir := model.NewDirectory(token.UserId, parentDirId, json.DirName)
	c.JSON(http.StatusOK, dir)
}

func (f FileManageController) DeleteDirectory(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json deleteDirectoryJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dirId, _ := strconv.Atoi(json.DirId)
	model.DeleteDirectory(token.UserId, dirId)
	c.JSON(http.StatusOK, gin.H{})
}

func (f FileManageController) RenameDirectory(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json renameDirectoryJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dirId, _ := strconv.Atoi(json.DirId)
	model.RenameDirectory(token.UserId, dirId, json.DirName)
	c.JSON(http.StatusOK, gin.H{})
}
