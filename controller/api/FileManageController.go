package api

import (
	"MedicalLowCode-backend/model"
	"MedicalLowCode-backend/util"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"path/filepath"
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

type getDirContentJson struct {
	DirId string `json:"dirId" binding:"required"`
}

type retFile struct {
	FileId   string `json:"fileId"`
	FileName string `json:"fileName"`
	Path     string `json:"path"`
}

type retDir struct {
	DirId   string `json:"dirId"`
	DirName string `json:"dirName"`
	Path    string `json:"path"`
}

func (f FileManageController) GetDirContent(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json getDirContentJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dirId, _ := strconv.Atoi(json.DirId)
	parentDirPath, err := model.QueryDirPath(token.UserId, dirId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dirs := model.GetDirsUnderDir(token.UserId, dirId)
	files := model.GetFilesUnderDir(token.UserId, dirId)
	dirContent := make(map[string]interface{})
	retDirs := make([]retDir, 0)
	retFiles := make([]retFile, 0)
	for _, dir := range dirs {
		dirPath := filepath.Join(parentDirPath, dir.DirName)
		retDirs = append(retDirs, retDir{
			DirId:   strconv.Itoa(dir.DirId),
			DirName: dir.DirName,
			Path:    dirPath,
		})
	}
	for _, file := range files {
		filePath := filepath.Join(parentDirPath, file.FileName)
		retFiles = append(retFiles, retFile{
			FileId:   strconv.Itoa(file.FileId),
			FileName: file.FileName,
			Path:     filePath,
		})
	}
	dirContent["directories"] = retDirs
	dirContent["files"] = retFiles
	c.JSON(http.StatusOK, dirContent)
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
	if dir == nil {
	}
	c.JSON(http.StatusOK, gin.H{})
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

func (f FileManageController) GetRootDir(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	dir := model.GetRootDir(token.UserId)
	if dir == nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"dirId": strconv.Itoa(dir.DirId),
	})
}
