package model

import (
	"gorm.io/gorm"
	"time"
)

type InferenceFile struct {
	Id         int
	UserId     int
	ModelId    int
	CreateTime time.Time
	FinishTime *time.Time
	Name       string
	Address    string
	Info       string
	Share      int
}

func (i InferenceFile) TableName() string {
	return "nnunet_inference"
}

func AddNnunetInferenceFile(userId int, modelId int, fileName string, address string) InferenceFile {
	file := InferenceFile{
		UserId:     userId,
		ModelId:    modelId,
		Name:       fileName,
		Address:    address,
		CreateTime: time.Now(),
	}
	if err := DB.Create(&file).Error; err != nil {
		panic(err)
	}
	return file
}

func QueryNnunetInferenceFile(userId int, modelId int, fileName string) *InferenceFile {
	file := InferenceFile{}
	result := DB.Where("user_id = ? and model_id = ? and name = ?", userId, modelId, fileName).Last(&file)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		} else {
			panic(err)
		}
	}
	return &file
}
