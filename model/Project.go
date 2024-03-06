package model

import (
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
	CanvasList         []Canvas `gorm:"foreignKey:ProjectId;references:ProjectId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (p Project) TableName() string {
	return "project"
}

type Canvas struct {
	CanvasId      int `gorm:"primaryKey;autoIncrement"`
	ProjectId     int
	CanvasName    string
	CanvasContent string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (c Canvas) TableName() string {
	return "canvas"
}

type Module struct {
	ModuleId      int
	UserId        int
	ModuleName    string
	ModuleContent string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (m Module) TableName() string {
	return "module"
}

func NewProject(userId int, projectName string, projectDescription string) Project {
	project := Project{
		UserId:             userId,
		ProjectName:        projectName,
		ProjectDescription: projectDescription,
		CanvasList: []Canvas{
			{CanvasName: "默认画布", CanvasContent: ""},
		},
	}
	if err := DB.Create(&project).Error; err != nil {
		panic(err)
	}
	return project
}

func CheckProjectPermission(userId int, projectId int) bool {
	project := Project{}
	result := DB.Where("user_id = ? and project_id = ?", userId, projectId).First(&project)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		} else {
			panic(err)
		}
	}
	return true
}

func DeleteProject(userId int, projectId int) {
	project := Project{}
	result := DB.Where("user_id = ? and project_id = ?", userId, projectId).Delete(&project)
	if err := result.Error; err != nil {
		panic(err)
	}
}
