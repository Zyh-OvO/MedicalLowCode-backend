package api

import (
	"MedicalLowCode-backend/model"
	"MedicalLowCode-backend/util"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/KyungWonPark/nifti"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const viewNiiGz = "data/view_nii_gz/"
const fileEnding = ".nii.gz"
const outputLog = "output.log"
const device = util.Device
const env = util.Env
const maxLabelValue = 5

type DefaultModelController struct {
}

type InitMessage struct {
	FileId  int    `json:"fileId"`
	Token   string `json:"token"`
	ModelId int    `json:"modelId"`
}

type InferenceProgress struct {
	Progress      float32 `json:"progress"`
	PassedTime    string  `json:"passedTime"`
	RemainingTime string  `json:"remainingTime"`
	Status        string  `json:"status"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (u DefaultModelController) WebsocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	flag := false
	var fileId int
	var userId int
	var modelId int

	for {
		// 读取客户端发送的消息
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		// 处理消息
		fmt.Printf("Received message: %s\n", p)
		if !flag {
			// 解析JSON字符串到map
			var message InitMessage
			if err := json.Unmarshal([]byte(p), &message); err != nil {
				fmt.Println("解析JSON失败:", err)
				return
			}

			// 获得解析后的字段值
			fileId = message.FileId
			// TODO:token判断不是很优雅
			token, err := util.ParseToken(message.Token)
			if err != nil {
				token, _ = util.GiveStaticToken()
			}
			userId = token.UserId
			modelId = message.ModelId

			InferenceFile := model.QueryNnunetInferenceFile(fileId)

			// 判断文件推理状态
			if InferenceFile == nil || InferenceFile.FinishTime != nil {
				var progress InferenceProgress
				if InferenceFile == nil {
					progress.Status = "no such file"
				} else {
					progress.Status = "finished"
				}
				jsonData, err := json.Marshal(progress)
				if err != nil {
					fmt.Println("Error marshalling JSON:", err)
					break
				}
				if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
					fmt.Println("Error writing message:", err)
					break
				}
				return
			}

			go WatchInferenceProgress(fileId, userId, modelId, conn)

			// 将 Message 结构体实例序列化为 JSON 字符串
			jsonData, err := json.Marshal(message)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				break
			}

			// 发送 JSON 字符串给客户端
			if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				fmt.Println("Error writing message:", err)
				break
			}

			flag = true
		}
	}
}

func (u DefaultModelController) Test(c *gin.Context) {
	username := c.Query("username")
	age := c.Query("age")

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"age":      age,
	})
}

func (u DefaultModelController) UploadNiiGzFile(c *gin.Context) {
	modelId, _ := strconv.Atoi(c.Request.FormValue("modelId"))
	// 检查
	if model.QueryNnunetModelReady(modelId) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Model not ready"})
		return
	}
	file, handler, err := c.Request.FormFile("nifti")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	name := handler.Filename
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot close file stream"})
			return
		}
	}(file)
	token := c.MustGet("token").(*util.Token)
	filePath := viewNiiGz + strconv.Itoa(token.UserId) + "/" + strconv.Itoa(modelId) + "/"
	if util.CreateFolderIfNotExists(filePath) != true {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
		return
	}

	// 加入数据库操作
	// TODO:文件名重复？现在是当新文件处理
	share := 1
	inferenceFile := model.AddNnunetInferenceFile(token.UserId, modelId, name, filePath+name, share)

	// 在本地创建一个同名的文件
	out, err := os.Create(inferenceFile.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}
	}(out)

	// 将上传的 nifti 文件复制到本地文件中
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	inputFolder := filePath + strconv.Itoa(inferenceFile.Id) + "_input/"
	util.CreateFolderIfNotExists(inputFolder)

	// 进行预处理，分开channels
	IdSlice := []int{inferenceFile.Id}
	NameSlice := []string{inferenceFile.Name}

	go PreprocessAndInferenceNiiGzFile(filePath, inputFolder, modelId, IdSlice, NameSlice)

	c.JSON(http.StatusOK, gin.H{"message": "NII file uploaded successfully",
		"file_id": inferenceFile.Id})
}

func (u DefaultModelController) ImageTest(c *gin.Context) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot close file stream"})
			return
		}
	}(file)
	// 创建一个名为 uploaded.png 的文件
	out, err := os.Create("uploaded.png")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}
	}(out)

	// 将上传的 PNG 文件复制到本地文件中
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
}

func (u DefaultModelController) ReturnMultipleImages(c *gin.Context) {
	// 从本地文件系统中读取多个PNG图片
	imagePaths := []string{
		"C:\\code\\go\\fengru-backend\\output_slice_20.png",
		"C:\\code\\go\\fengru-backend\\output_slice_21.png",
		"C:\\code\\go\\fengru-backend\\output_slice_22.png"}

	var images []string
	for _, path := range imagePaths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image file"})
			return
		}
		encodedImage := base64.StdEncoding.EncodeToString(data)
		images = append(images, encodedImage)
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}

func (u DefaultModelController) ReturnNiiGzFile(c *gin.Context) {
	fmt.Println("请求的URL是：", c.Request.URL.String())
	token := c.MustGet("token").(*util.Token)
	fmt.Println(token)
	fileId, _ := strconv.Atoi(c.Param("id"))
	inferenceFile := model.QueryNnunetInferenceFile(fileId)
	filePath := inferenceFile.Address
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read nii.gz file"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+inferenceFile.Name)

	c.Data(http.StatusOK, "application/octet-stream", data)

	// 删除已经查看的临时文件
	//os.Remove(filePath)
}

func (u DefaultModelController) GetNonZeroLocation(c *gin.Context) {
	fmt.Println("请求的URL是：", c.Request.URL.String())
	token := c.MustGet("token").(*util.Token)
	fmt.Println(token)
	fileId, _ := strconv.Atoi(c.Param("id"))
	inferenceFile := model.QueryNnunetInferenceFile(fileId)
	if inferenceFile.FinishTime == nil {
		c.JSON(http.StatusOK, gin.H{"status": "inferencing"})
		return
	}
	segmentationFilePath := inferenceFile.OutputFolder + inferenceFile.Name
	nifti1Image := nifti.Nifti1Image{}
	nifti1Image.LoadImage(segmentationFilePath, true)
	dims := nifti1Image.GetDims()
	var nonZero [][]int
	for i := 0; i <= maxLabelValue; i++ { // 0-5 共六个2维数组
		nonZero = append(nonZero, []int{})
	}

	maxLabel := 0
	index := 0
	for z := 0; z < dims[2]; z++ {
		for y := 0; y < dims[1]; y++ {
			for x := 0; x < dims[0]; x++ {
				label := int(nifti1Image.GetAt(uint32(x), uint32(y), uint32(z), 0))
				if label > maxLabelValue {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "max label value exceeded"})
					return
				}
				if label > maxLabel {
					maxLabel = label
				}
				nonZero[label] = append(nonZero[label], index)
				index++
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"nonZero1": nonZero[1], "nonZero2": nonZero[2], "nonZero3": nonZero[3], "nonZero4": nonZero[4], "nonZero5": nonZero[5], "status": "finished", "maxLabel": maxLabel})
}

// 将整数数组转换为字符串
func intArrayToString(arr []int) string {
	var strArr []string
	for _, v := range arr {
		strArr = append(strArr, strconv.Itoa(v))
	}
	return strings.Join(strArr, ", ")
}

func PreprocessNiiGzFile(filePath string, inputFolder string, modelId int, IdSlice []int, NameSlice []string) {
	channel := model.QueryNnunetModelChannel(modelId)
	outputFolder := filePath + "/output/"
	util.CreateFolderIfNotExists(outputFolder)
	if channel == 1 {
		//	如果文件channel数为1则只对文件重命名
		for i := 0; i < len(NameSlice); i++ {
			srcFile, err := os.Open(filePath + NameSlice[i])
			if err != nil {
				fmt.Println("Error opening source file:", err)
				return
			}
			defer srcFile.Close()
			inputFileName := strings.ReplaceAll(NameSlice[i], fileEnding, "") + "_0000" + fileEnding
			fmt.Println("inputFileName:" + inputFileName)
			destFile, err := os.Create(inputFolder + inputFileName)
			if err != nil {
				fmt.Println("Error creating destination file:", err)
				return
			}
			defer destFile.Close()
			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				fmt.Println("Error copying file:", err)
				return
			}
		}
	} else {
		//	TODO:多channel分开
		for i := 0; i < len(NameSlice); i++ {
			srcFile := filePath + NameSlice[i]
			SplitNiiGFile(srcFile, NameSlice[i], channel, inputFolder)
		}
	}

	model.SetNnunetInferenceFilePreprocess(IdSlice)

	fmt.Println("File copied and renamed successfully.")
}

func SplitNiiGFile(srcFile string, fileName string, modelChannel int, inputFolder string) {
	nifti1Image := nifti.Nifti1Image{}
	nifti1Image.LoadImage(srcFile, true)
	dims := nifti1Image.GetDims()
	channels := nifti1Image.GetDims()[3]
	if channels == modelChannel {
		for i := 0; i < modelChannel; i++ {
			img := nifti.NewImg(dims[0], dims[1], dims[2], 1)
			img.SetNewHeader(nifti1Image.GetHeader())
			img.SetHeaderDim2(dims[0], dims[1], dims[2], 1)
			for z := 0; z < dims[2]; z++ {
				for y := 0; y < dims[1]; y++ {
					for x := 0; x < dims[0]; x++ {
						value := nifti1Image.GetAt(uint32(x), uint32(y), uint32(z), uint32(i))
						img.SetAt(uint32(x), uint32(y), uint32(z), 0, value)

					}
				}
			}
			fileName, _ = strings.CutSuffix(fileName, ".nii.gz")
			inputFileName := fileName + fmt.Sprintf("_%04d", i) + ".nii"
			img.Save(inputFolder + inputFileName)
		}
	} else {
		panic("Wrong srcFile dims!")
		//	TODO: file与model维数不匹配报错
	}
}

func InferenceNiiGzFile(filePath string, inputFolder string, modelId int, idSlice []int, nameSlice []string) {
	outputFolder := filePath + "output/"
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && conda activate nnUNet && nnUNetv2_predict -i %s -o %s -c 3d_fullres -d %d -device %s -f all -chk checkpoint_best.pth --disable_progress_bar > %s", env, inputFolder, outputFolder, modelId, device, inputFolder+outputLog))
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("命令执行失败:", err)
		return
	}
	// 输出命令执行结果
	fmt.Println("命令输出:", string(output))
	model.SetNnunetInferenceFileProcessed(idSlice, outputFolder)
}

func PreprocessAndInferenceNiiGzFile(filePath string, inputFolder string, modelId int, idSlice []int, nameSlice []string) {
	// 预处理
	PreprocessNiiGzFile(filePath, inputFolder, modelId, idSlice, nameSlice)
	// 进行推断

	InferenceNiiGzFile(filePath, inputFolder, modelId, idSlice, nameSlice)
}

func WatchInferenceProgress(id int, userId int, modelId int, conn *websocket.Conn) {
	inputFolder := viewNiiGz + strconv.Itoa(userId) + "/" + strconv.Itoa(modelId) + "/" + strconv.Itoa(id) + "_input/"
	filePath := inputFolder + outputLog
	t, err := tail.TailFile(filePath, tail.Config{
		Follow:    true,                                           // 实时监听文件变化
		ReOpen:    true,                                           // 当文件被删除或者移动后重新打开
		MustExist: false,                                          // 如果文件不存在不报错
		Poll:      true,                                           // 使用轮询的方式监听文件变化
		Location:  &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}, // 从文件末尾开始读取
	})
	if err != nil {
		log.Fatal(err)
	}
	defer t.Cleanup()

	totalFiles := 0
	//var nowFile int

	// 循环处理文件内容
	for line := range t.Lines {
		// 输出文件内容
		log.Println(line.Text)
		str := line.Text
		if strings.HasPrefix(str, "startSGSGSG:") { // 推断开始
			// 提取后续子字符串
			startIndex := len("startSGSGSG:")
			totalFiles, _ = strconv.Atoi(str[startIndex:])
		} else if str == "finishedSGSGSG" {
			var progress InferenceProgress
			progress.Status = "finished"
			jsonData, err := json.Marshal(progress)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				break
			}
			if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				fmt.Println("Error writing message:", err)
				break
			}
		} else if totalFiles == 1 && strings.HasPrefix(str, "infile:") { // 单个文件的推断任务
			progress := ParseProgressLine(str, "infile:")
			jsonData, err := json.Marshal(progress)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				break
			}
			if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				fmt.Println("Error writing message:", err)
				break
			}
		} else if totalFiles > 1 && strings.HasPrefix(str, "inlist:") { // 多个文件的推断任务
			progress := ParseProgressLine(str, "inlist:")
			jsonData, err := json.Marshal(progress)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				break
			}
			if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				fmt.Println("Error writing message:", err)
				break
			}
		}
	}
}

func ParseProgressLine(str string, prefix string) InferenceProgress {
	var proress InferenceProgress
	proress.Status = "inferencing"
	_, s1, _ := strings.Cut(str, prefix)
	progressSlice := strings.Split(s1, "/")
	f, _ := strconv.ParseFloat(progressSlice[0], 32)
	proress.Progress = float32(f) * 100
	secondsF, _ := strconv.ParseFloat(progressSlice[1], 32)
	seconds := int(secondsF + 0.5)
	hours := seconds / 3600          // 小时数
	minutes := (seconds % 3600) / 60 // 分钟数
	seconds = seconds % 60           // 剩余的秒数
	proress.PassedTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	secondsF, _ = strconv.ParseFloat(progressSlice[2], 32)
	seconds = int(secondsF + 0.5)
	hours = seconds / 3600          // 小时数
	minutes = (seconds % 3600) / 60 // 分钟数
	seconds = seconds % 60          // 剩余的秒数
	proress.RemainingTime = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	fmt.Println(proress)
	return proress
}
