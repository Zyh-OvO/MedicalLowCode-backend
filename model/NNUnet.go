package model

import (
	"gorm.io/gorm"
	"time"
)

type InferenceFile struct {
	Id           int
	UserId       int
	ModelId      int
	CreateTime   time.Time
	FinishTime   *time.Time
	Name         string
	Address      string
	Info         string
	Share        int
	Preprocessed int
	OutputFolder string
}

type Model struct {
	Id          int
	UserId      int
	Name        string
	Description string
	Cover       string
	Share       int
	Channel     int
	Ready       int
}

func (i InferenceFile) TableName() string {
	return "nnunet_inference"
}
func (m Model) TableName() string {
	return "nnunet_models"
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

func QueryNnunetInferenceFile(id int) *InferenceFile {
	file := InferenceFile{}
	result := DB.Where("id = ?", id).Last(&file)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		} else {
			panic(err)
		}
	}
	return &file
}

func SetNnunetInferenceFilePreprocess(idSlice []int) {
	if err := DB.Model(&InferenceFile{}).Where("id IN ?", idSlice).Update("preprocessed", 1).Error; err != nil {
		panic(err)
	}
}

func SetNnunetInferenceFileProcessed(idSlice []int, outputFolder string) {
	if err := DB.Model(&InferenceFile{}).Where("id IN ?", idSlice).Updates(map[string]interface{}{
		"finish_time":   time.Now(),
		"output_folder": outputFolder,
	}).Error; err != nil {
		panic(err)
	}
}

func QueryNnunetModelChannel(modelId int) int {
	var channel int
	if err := DB.Model(&Model{}).Where("id = ?", modelId).Pluck("channel", &channel).Error; err != nil {
		panic("failed to query database")
	}
	return channel
}
