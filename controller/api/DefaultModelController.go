package api

import (
	"MedicalLowCode-backend/model"
	"MedicalLowCode-backend/util"
	"encoding/base64"
	"fmt"
	"github.com/KyungWonPark/nifti"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const viewNiiGz = "data/view_nii_gz/"

type DefaultModelController struct {
}

func (u DefaultModelController) Test(c *gin.Context) {
	username := c.Query("username")
	age := c.Query("age")

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"age":      age,
	})
}

func (u DefaultModelController) NiiTest(c *gin.Context) {
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
	modelId, _ := strconv.Atoi(c.Request.FormValue("modelId"))
	token := c.MustGet("token").(*util.Token)
	filePath := viewNiiGz + strconv.Itoa(token.UserId) + "/"
	if createFolderIfNotExists(filePath) != true {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
		return
	}

	// 加入数据库操作
	// TODO:文件名重复？
	inferenceFile := model.AddNnunetInferenceFile(token.UserId, modelId, name, filePath+name)

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
	//name := c.Query("file")
	//modelId, err := strconv.Atoi(c.Query("model_id"))
	//fmt.Println(modelId)
	token := c.MustGet("token").(*util.Token)
	fmt.Println(token)
	fileId, _ := strconv.Atoi(c.Param("id"))
	//filePath := "/Users/qhy/Desktop/lth/冯如杯/hepaticvessel_001.nii.gz"
	inferenceFile := model.QueryNnunetInferenceFile(fileId)
	filePath := inferenceFile.Address
	//filePath := "C:\\BUAA\\3rd\\FengRu\\MICCAI-LITS2017\\Task06_Lung\\Task06_Lung\\labelsTr\\lung_001.nii.gz"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read nii.gz file"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+inferenceFile.Name)
	c.Header("file_id", strconv.Itoa(inferenceFile.Id))

	c.Data(http.StatusOK, "application/octet-stream", data)

	// 删除已经查看的临时文件
	//os.Remove(filePath)
}

func (u DefaultModelController) ReturnSegFile(c *gin.Context) {
	filePath := "C:\\BUAA\\3rd\\FengRu\\MICCAI-LITS2017\\Task06_Lung\\Task06_Lung\\labelsTr\\lung_010.nii.gz"
	//filePath := "C:\\BUAA\\3rd\\FengRu\\MICCAI-LITS2017\\Task06_Lung\\Task06_Lung\\labelsTr\\lung_001.nii.gz"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read nii.gz file"})
		return
	}

	filename := "seg.nii.gz"
	c.Header("Content-Disposition", "attachment; filename="+filename)

	c.Data(http.StatusOK, "application/octet-stream", data)
}

func (u DefaultModelController) GetNoneZeroLocation(c *gin.Context) {
	nifti1Image := nifti.Nifti1Image{}
	nifti1Image.LoadImage("C:\\BUAA\\3rd\\FengRu\\MICCAI-LITS2017\\Task06_Lung\\Task06_Lung\\labelsTr\\lung_010.nii.gz", true)
	dims := nifti1Image.GetDims()
	var nonZero []int
	index := 0
	for z := 0; z < dims[2]; z++ {
		for y := 0; y < dims[1]; y++ {
			for x := 0; x < dims[0]; x++ {
				if nifti1Image.GetAt(uint32(x), uint32(y), uint32(z), 0) != 0 {
					nonZero = append(nonZero, index)
				}
				index++
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"nonZero": nonZero})
}

func (u DefaultModelController) DimTest(c *gin.Context) {
	nifti1Image := nifti.Nifti1Image{}
	nifti1Image.LoadImage("C:\\BUAA\\3rd\\FengRu\\MICCAI-LITS2017\\Task06_Lung\\Task06_Lung\\imagesTr\\lung_023.nii.gz", true)
	dims := nifti1Image.GetDims()
	fmt.Println(dims)
	fmt.Println(1)
}

// 将整数数组转换为字符串
func intArrayToString(arr []int) string {
	var strArr []string
	for _, v := range arr {
		strArr = append(strArr, strconv.Itoa(v))
	}
	return strings.Join(strArr, ", ")
}

func createFolderIfNotExists(folderPath string) bool {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return false
	}
	folderPath = filepath.Join(currentDir, folderPath)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, 0755)
		if err != nil {
			fmt.Println("Error creating folder:", err)
			return false
		} else {
			fmt.Println("Folder created successfully:", folderPath)
			return true
		}
	} else {
		fmt.Println("Folder already exists, no action taken:", folderPath)
		return true
	}
}
