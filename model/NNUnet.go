package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type InferenceFile struct {
	Id           int        `json:"id"`
	UserId       int        `json:"user_id"`
	ModelId      int        `json:"model_id"`
	CreateTime   time.Time  `json:"create_time"`
	FinishTime   *time.Time `json:"finish_time"`
	Name         string     `json:"name"`
	Address      string     `json:"address"`
	Info         string     `json:"info"`
	Share        int        `json:"share"`
	Preprocessed int        `json:"preprocessed"`
	OutputFolder string     `json:"output_folder"`
}

type NnunetModel struct {
	Id                              int     `json:"id"`
	UserId                          int     `json:"user_id"`
	Name                            string  `json:"name"`
	Description                     string  `json:"description"`
	Cover                           string  `json:"cover"`
	Share                           int     `json:"share"`
	Channel                         int     `json:"channel"`
	Ready                           int     `json:"ready"`
	Reference                       string  `json:"reference"`
	License                         string  `json:"license"`
	Release                         string  `json:"release"`
	TensorImageSize                 string  `json:"tensor_image_size"`
	Label                           int     `json:"label"`
	LabelNames                      string  `json:"label_names"`
	NumTraining                     int     `json:"num_training"`
	NumTest                         int     `json:"num_test"`
	FileEnding                      string  `json:"file_ending"`
	ChannelNames                    string  `json:"channel_names"`
	FileNameInSystemWithNoExtension string  `json:"file_name_in_system_with_no_extension"`
	MeanValidationDice              float64 `json:"mean_validation_dice"`
	DatasetInfo                     string  `json:"dataset_info"`
	NetStructure                    string  `json:"net_structure"`
}

type ModelInfo struct {
	Name            string            `json:"name"` //
	Description     string            `json:"description"`
	Reference       string            `json:"reference"`
	TensorImageSize string            `json:"tensorImageSize"`
	Labels          map[string]string `json:"labels"` //
	NumTraining     int               `json:"numTraining"`
	NumTest         int               `json:"numTest"`
	FileEnding      string            `json:"file_ending"`
	ChannelNames    map[string]string `json:"channel_names"` //
	License         string            `json:"license"`
	Release         string            `json:"release"`
	// 以上信息可以生成dataset.json
	Share int    `json:"share"`
	Cover string `json:"cover"`
}

func (i InferenceFile) TableName() string {
	return "nnunet_inference"
}
func (m NnunetModel) TableName() string {
	return "nnunet_models"
}

func AddNnunetInferenceFile(userId int, modelId int, fileName string, address string, share int) InferenceFile {
	file := InferenceFile{
		UserId:     userId,
		ModelId:    modelId,
		Name:       fileName,
		Address:    address,
		CreateTime: time.Now(),
		Share:      share,
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

func QueryUserInferenceFileList(userId int) []InferenceFile {
	var fileList []InferenceFile
	DB.Where("user_id = ?", userId).Find(&fileList)
	return fileList
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

func AddNnunetModel(info ModelInfo, userId int) NnunetModel {
	channel := len(info.ChannelNames)
	label := len(info.Labels)
	labels, _ := json.Marshal(info.Labels)
	channels, _ := json.Marshal(info.ChannelNames)
	model := NnunetModel{
		UserId:          userId,
		Name:            info.Name,
		Description:     info.Description,
		Cover:           info.Cover,
		Share:           info.Share,
		Channel:         channel,
		Ready:           0, // 1为训练好
		Reference:       info.Reference,
		License:         info.License,
		Release:         info.Release,
		TensorImageSize: info.TensorImageSize,
		Label:           label,
		LabelNames:      string(labels),
		NumTraining:     info.NumTraining,
		NumTest:         info.NumTest,
		FileEnding:      info.FileEnding,
		ChannelNames:    string(channels),
		DatasetInfo:     string("{}"),
		NetStructure:    string("{}"),
	}
	if err := DB.Create(&model).Error; err != nil {
		panic(err)
	}
	return model
}

func UpdateNnunetModel(info ModelInfo, modelId int, ready int, fileNameInSystemWithNoExtension string) NnunetModel {
	channel := len(info.ChannelNames)
	label := len(info.Labels)
	labels, _ := json.Marshal(info.Labels)
	channels, _ := json.Marshal(info.ChannelNames)
	if err := DB.Model(&NnunetModel{}).Where("id = ? ", modelId).Updates(map[string]interface{}{
		//UserId:          userId,
		"Name":                            info.Name,
		"Description":                     info.Description,
		"Cover":                           info.Cover,
		"Share":                           info.Share,
		"Channel":                         channel,
		"Ready":                           ready, // 1为训练好
		"Reference":                       info.Reference,
		"License":                         info.License,
		"Release":                         info.Release,
		"TensorImageSize":                 info.TensorImageSize,
		"Label":                           label,
		"LabelNames":                      string(labels),
		"NumTraining":                     info.NumTraining,
		"NumTest":                         info.NumTest,
		"FileEnding":                      info.FileEnding,
		"ChannelNames":                    string(channels),
		"FileNameInSystemWithNoExtension": fileNameInSystemWithNoExtension,
	}).Error; err != nil {
		panic(err)
	}
	nnunetModel := NnunetModel{}
	DB.Where("id = ?", modelId).Last(&nnunetModel)
	return nnunetModel
}

func SetNnunetModelReady(modelId int, ready int) {
	if err := DB.Model(&NnunetModel{}).Where("id = ? ", modelId).Updates(map[string]interface{}{
		"Ready": ready, // 1为训练好
	}).Error; err != nil {
		panic(err)
	}
}

func QueryNnunetModelChannel(modelId int) int {
	var channel int
	if err := DB.Model(&NnunetModel{}).Where("id = ?", modelId).Pluck("channel", &channel).Error; err != nil {
		panic("failed to query database")
	}
	return channel
}

func QueryNnunetModel(modelId int) NnunetModel {

	nnunetModel := NnunetModel{}
	DB.Where("id = ?", modelId).Last(&nnunetModel)

	return nnunetModel
}

func QueryNnunetModelReady(modelId int) int {
	ready := 0
	if err := DB.Model(&NnunetModel{}).Where("id = ?", modelId).Pluck("ready", &ready).Error; err != nil {
		panic("failed to query database")
	}
	return ready
}

func QueryUserNnunetModelList(userId int) []NnunetModel {
	var modelList []NnunetModel
	DB.Where("user_id = ? OR share = ?", userId, 1).Find(&modelList)
	return modelList
}
