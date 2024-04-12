package api

import (
	"MedicalLowCode-backend/model"
	"MedicalLowCode-backend/util"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type NnunetModelController struct {
}

type NnunetModelListElement struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	Ready       int    `json:"ready"`
}

const toTrain = "data/to_train/" // 存储模型训练原数据

func (u NnunetModelController) SetModelInfo(c *gin.Context) {

	fmt.Println("请求的URL是：", c.Request.URL.String())

	token := c.MustGet("token").(*util.Token)
	fmt.Println(token)
	userId := token.UserId

	// 读取传输的数据
	info := c.Request.FormValue("modelInfo")
	fmt.Println(info)
	// 解析 JSON 数据到结构体
	var modelInfo model.ModelInfo
	if err := json.Unmarshal([]byte(info), &modelInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON info"})
		return
	}

	if len(modelInfo.ChannelNames) > 1 {
		modelInfo.TensorImageSize = "4D"
	} else {
		modelInfo.TensorImageSize = "3D"
	}

	// TODO: 支持其他文件格式
	modelInfo.FileEnding = ".nii.gz"

	nnUnetModel := model.AddNnunetModel(modelInfo, userId)
	fmt.Println(nnUnetModel)

	// 获取表单中的文件信息
	data, handler, err := c.Request.FormFile("modelData")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	fileName := handler.Filename

	// 在本地创建一个文件,并把传入的zip文件数据拷入其中
	util.CreateFolderIfNotExists(toTrain)
	fileName = strings.ReplaceAll(fileName, "_", "-") // 将_替换为-，以免nnunet不支持
	fileNameInSystemWithNoExtension, _ := strings.CutSuffix(fileName, ".zip")
	fileName = fmt.Sprintf("Task%02d_%s", nnUnetModel.Id, fileName)
	out, err := os.Create(toTrain + fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	_, err = io.Copy(out, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer func(data multipart.File) {
		err := data.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot close file stream"})
			return
		}
	}(data)

	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}
	}(out)

	destDir, _ := strings.CutSuffix(fileName, ".zip")
	destDir = toTrain + destDir + "/"
	fmt.Println(toTrain + fileName)
	err = util.Unzip(toTrain+fileName, destDir)
	if err != nil {
		fmt.Println("解压 ZIP 文件时出错:", err)
		return
	}

	fmt.Println("ZIP 文件已成功解压到目录:", destDir)

	newDestDir := util.FindFileNumFolder(destDir)
	if newDestDir != destDir {
		fileNameInSystemWithNoExtension, _ = strings.CutPrefix(newDestDir, destDir)
		fileNameInSystemWithNoExtension, _ = strings.CutSuffix(fileNameInSystemWithNoExtension, "/")
		list := strings.Split(fileNameInSystemWithNoExtension, "_")
		fileNameInSystemWithNoExtension = list[len(list)]
	}

	trNums := util.FindFileNumInFolder(destDir+"imagesTr", modelInfo.FileEnding)
	tsNums := util.FindFileNumInFolder(destDir+"imagesTs", modelInfo.FileEnding)

	if trNums == -1 || tsNums == -1 {
		fmt.Println("遍历文件时出错:", err)
		return
	}

	modelInfo.NumTraining = trNums
	modelInfo.NumTest = tsNums

	nnUnetModel = model.UpdateNnunetModel(modelInfo, nnUnetModel.Id, 0, fileNameInSystemWithNoExtension)

	// 调用 os.Stat 获取文件信息
	datasetJson := destDir + "dataset.json"
	_, err = os.Stat(datasetJson)

	// 将结构体转换为 JSON
	jsonData, err := json.MarshalIndent(modelInfo, "", "    ")
	if err != nil {
		fmt.Println("转换为 JSON 时出错:", err)
		return
	}

	// 判断dataset.json文件是否存在，若存在则保留原来的dataset.json文件
	// TODO：根据原dataset.json修改modelInfo
	_, err = os.Stat(destDir + "dataset.json")
	if os.IsNotExist(err) {
		fmt.Printf("文件 %s 不存在\n", datasetJson)
		// 将json文件输出
		err = ioutil.WriteFile(destDir+"dataset.json", []byte(jsonData), 0644)
		if err != nil {
			fmt.Println("写入文件时出错:", err)
			return
		}

		fmt.Println("JSON 数据已成功写入文件 dataset.json")
	} else if err != nil {
		fmt.Println("获取文件信息时出错:", err)
	} else {
		fmt.Printf("文件 %s 存在\n", datasetJson)
	}

	go TrainNnunetModel(destDir, nnUnetModel.Id)

	c.JSON(http.StatusOK, gin.H{"message": "数据已收到", "modelId": nnUnetModel.Id, "model": nnUnetModel})

}

func (u NnunetModelController) GetModelList(c *gin.Context) {
	fmt.Println("请求的URL是：", c.Request.URL.String())
	token := c.MustGet("token").(*util.Token)
	fmt.Println(token)
	modelList := model.QueryUserNnunetModelList(token.UserId) // 自己的模型或公开的模型 // 目前应该是默认所有模型都公开

	var jsonDataList []interface{}

	for i := 0; i < len(modelList); i++ {
		var element NnunetModelListElement
		element.Name = modelList[i].Name
		element.Id = modelList[i].Id
		element.UserId = modelList[i].UserId
		element.Description = modelList[i].Description
		element.Cover = modelList[i].Cover
		element.Ready = modelList[i].Ready
		jsonDataList = append(jsonDataList, element)
	}

	c.JSON(http.StatusOK, gin.H{"modelList": jsonDataList})

}

func (u NnunetModelController) GetModelInfoInference(c *gin.Context) {
	fmt.Println("请求的URL是：", c.Request.URL.String())
	token := c.MustGet("token").(*util.Token)
	fmt.Println(token)
	historyFileList := model.QueryUserInferenceFileList(token.UserId)
	modelId, _ := strconv.Atoi(c.Param("modelId"))
	nnunetModel := model.QueryNnunetModel(modelId)
	if nnunetModel.Ready != 1 { //模型还在inference TODO: 模型不存在
		c.JSON(http.StatusOK, gin.H{"status": "inferencing"})
		return
	}
	var labels map[string]string
	if err := json.Unmarshal([]byte(nnunetModel.LabelNames), &labels); err != nil {
		panic(err)
	}
	modelName := nnunetModel.Name
	modelDescription := nnunetModel.Description

	// dice
	nnUNetResultsDir := os.Getenv("nnUNet_results")
	fmt.Println("env:\n" + nnUNetResultsDir)
	datasetFolder := fmt.Sprintf("Dataset%03d_%s", nnunetModel.Id, nnunetModel.FileNameInSystemWithNoExtension)
	resultDir := nnUNetResultsDir + "/" + datasetFolder + "/" + util.TrainPlanFolder + "/"
	resultFoldDir := resultDir + util.TrainFoldFolder + "/"
	meanValidationDice := GetMeanValidationDice(resultFoldDir)
	datasetInfo := GetDatasetInfoJson(resultDir)

	fmt.Println(resultDir)
	fmt.Println(resultFoldDir)
	fmt.Println(meanValidationDice)
	fmt.Println(datasetInfo)
	// 返回json文件

	c.JSON(http.StatusOK, gin.H{"history_file_list": historyFileList, "label_names": labels, "model_name": modelName, "description": modelDescription, "status": "finished", "mean_validation_dice": meanValidationDice, "dataset_info": datasetInfo})

}

func TrainNnunetModel(destDir string, modelId int) {
	ConvertNnunetModelData(destDir, modelId)
	PreprocessNnunetModel(modelId)
	TrainAllFold(modelId)
}

func ConvertNnunetModelData(destDir string, modelId int) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && conda activate nnUNet && nnUNetv2_convert_MSD_dataset -i %s -overwrite_id %d", env, destDir, modelId))
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("命令执行失败:", err)
		return
	}
	// 输出命令执行结果
	fmt.Println("命令输出:", string(output))
}

func PreprocessNnunetModel(modelId int) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && conda activate nnUNet && nnUNetv2_plan_and_preprocess -d %d -c 3d_fullres --verify_dataset_integrity", env, modelId))
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("命令执行失败:", err)
		return
	}
	// 输出命令执行结果
	fmt.Println("命令输出:", string(output))
}

func TrainAllFold(modelId int) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && conda activate nnUNet && nnUNetv2_train %d 3d_fullres all --npz --c -device %s", env, modelId, device))
	command := fmt.Sprintf("source %s && conda activate nnUNet && nnUNetv2_train %d 3d_fullres all --npz --c -device %s", env, modelId, device)
	fmt.Println(command)
	//go WatchTrainingProgress(modelId) TODO:实时监测训练进度
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("命令执行失败:", err)
		return
	}
	// 输出命令执行结果
	fmt.Println("命令输出:", string(output))
	model.SetNnunetModelReady(modelId, 1)
}

func GetMeanValidationDice(folder string) float64 {
	filePath := SearchFileInFolder(folder, "training_log")
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("无法打开文件：", file.Name())
		panic(err)
	}
	defer file.Close()

	// 创建一个 Scanner 来扫描文件内容
	scanner := bufio.NewScanner(file)
	iconStr := "Mean Validation Dice:  "
	var lastRow string

	// 逐行读取文件内容
	for scanner.Scan() {
		line := scanner.Text()
		// 查找子串的索引
		index := strings.Index(line, iconStr)

		// 如果找到子串
		if index != -1 {
			// 输出子串后的内容
			lastRow = fmt.Sprintf(line[index+len(iconStr):])
		}
	}

	// 检查扫描过程中是否出现错误
	if err := scanner.Err(); err != nil {
		fmt.Println("扫描文件时出错：", err)
	}
	lastRow = strings.TrimSpace(lastRow)
	dice, _ := strconv.ParseFloat(lastRow, 64)
	return dice
}

func GetDatasetInfoJson(folder string) string {
	filePath := SearchFileInFolder(folder, "plans.json")
	jsonFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("命令执行失败:", err)
	}

	// 将 JSON 文件内容转换为字符串
	jsonString := string(jsonFile)
	return jsonString
}

func SearchFileInFolder(absoluteFolder string, prefix string) string {
	var filePath string
	// 遍历指定文件夹下的所有文件和文件夹
	err := filepath.Walk(absoluteFolder, func(path string, info os.FileInfo, err error) error {
		// 忽略目录
		if info.IsDir() {
			return nil
		}
		// 检查文件名是否以 "training_log" 开头
		if strings.HasPrefix(filepath.Base(path), prefix) {
			// 如果是，则打印文件路径
			fmt.Println("找到匹配的文件：", path)
			filePath = path
		}
		return nil
	})

	if err != nil {
		fmt.Println("搜索文件时出错：", err)
	}
	return filePath
}

//func WatchTrainingProgress(modelId int) {
//}
