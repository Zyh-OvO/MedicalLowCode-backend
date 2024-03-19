package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Project struct {
	ProjectId          int `gorm:"primaryKey;autoIncrement"`
	UserId             int
	ProjectName        string
	ProjectDescription string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt
	Canvas             Canvas `gorm:"foreignKey:ProjectId;references:ProjectId"`
}

func (p *Project) TableName() string {
	return "project"
}

func (p *Project) BeforeDelete(tx *gorm.DB) (err error) {
	canvas := Canvas{}
	fmt.Printf("%+v\n", p)
	result := DB.Where("project_id = ?", p.ProjectId).Delete(&canvas)
	err = result.Error
	return
}

type Canvas struct {
	CanvasId      int `gorm:"primaryKey;autoIncrement"`
	ProjectId     int
	CanvasName    string
	CanvasContent string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func (c *Canvas) TableName() string {
	return "canvas"
}

type Module struct {
	ModuleId      int `gorm:"primaryKey;autoIncrement"`
	UserId        int
	ModuleName    string
	ModuleContent string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func (m Module) TableName() string {
	return "module"
}

func NewProject(userId int, projectName string, projectDescription string) Project {
	project := Project{
		UserId:             userId,
		ProjectName:        projectName,
		ProjectDescription: projectDescription,
		Canvas:             Canvas{CanvasName: "默认画布", CanvasContent: ""},
	}
	if err := DB.Create(&project).Error; err != nil {
		panic(err)
	}
	return project
}

func DeleteProject(userId int, projectId int) {
	project := Project{
		ProjectId: projectId,
		UserId:    userId,
	}
	result := DB.Where("user_id = ? and project_id = ?", userId, projectId).Delete(&project)
	if err := result.Error; err != nil {
		panic(err)
	}
}

func EditProject(userId int, projectId int, projectName *string, projectDescription *string) *Project {
	project := Project{}
	var selectFields []string
	if projectName != nil {
		project.ProjectName = *projectName
		selectFields = append(selectFields, "project_name")
	}
	if projectDescription != nil {
		project.ProjectDescription = *projectDescription
		selectFields = append(selectFields, "project_description")
	}
	result := DB.Where("user_id = ? and project_id = ?", userId, projectId).Select(selectFields).Updates(&project)
	if err := result.Error; err != nil {
		panic(err)
	}
	return QueryProject(userId, projectId)
}
func QueryProject(userId int, projectId int) *Project {
	project := Project{}
	result := DB.Preload("Canvas").Where("user_id = ? and project_id = ?", userId, projectId).First(&project)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		} else {
			panic(err)
		}
	}
	return &project
}

func QueryProjectList(userId int) []Project {
	var projectList []Project
	result := DB.Where("user_id = ?", userId).Find(&projectList)
	if err := result.Error; err != nil {
		panic(err)
	}
	return projectList
}

func EditCanvas(userId int, projectId int, canvasContent string) *Canvas {
	canvas := QueryCanvas(userId, projectId)
	if canvas == nil {
		return nil
	}
	canvas.ProjectId = projectId
	canvas.CanvasContent = canvasContent
	result := DB.Where("project_id = ?", projectId).Select("canvas_content").Updates(canvas)
	if err := result.Error; err != nil {
		panic(err)
	}
	return canvas
}

func QueryCanvas(userId int, projectId int) *Canvas {
	project := QueryProject(userId, projectId)
	if project == nil {
		return nil
	}
	return &project.Canvas
}

func QueryPersonalModule(userId int, moduleId int) *Module {
	module := Module{}
	result := DB.Where("user_id = ? and module_id = ?", userId, moduleId).First(&module)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		} else {
			panic(err)
		}
	}
	return &module
}

func QueryPersonalModules(userId int) []Module {
	var modules []Module
	result := DB.Where("user_id = ?", userId).Find(&modules)
	if err := result.Error; err != nil {
		panic(err)
	}
	return modules
}

func NewModule(userId int, moduleName string, moduleContent string) Module {
	module := Module{
		UserId:        userId,
		ModuleName:    moduleName,
		ModuleContent: moduleContent,
	}
	if err := DB.Create(&module).Error; err != nil {
		panic(err)
	}
	return module
}

func DeleteModule(userId int, moduleId int) {
	module := Module{
		ModuleId: moduleId,
		UserId:   userId,
	}
	result := DB.Where("user_id = ? and module_id = ?", userId, moduleId).Delete(&module)
	if err := result.Error; err != nil {
		panic(err)
	}
}

func EditModule(userId int, moduleId int, moduleName *string, moduleContent *string) *Module {
	module := Module{}
	var selectFields []string
	if moduleName != nil {
		module.ModuleName = *moduleName
		selectFields = append(selectFields, "module_name")
	}
	if moduleContent != nil {
		module.ModuleContent = *moduleContent
		selectFields = append(selectFields, "module_content")
	}
	result := DB.Where("user_id = ? and module_id = ?", userId, moduleId).Select(selectFields).Updates(&module)
	if err := result.Error; err != nil {
		panic(err)
	}
	return QueryPersonalModule(userId, moduleId)
}
