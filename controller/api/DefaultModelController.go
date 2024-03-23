package api

import (
	"encoding/base64"
	"fmt"
	"github.com/KyungWonPark/nifti"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

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
	file, _, err := c.Request.FormFile("nifti")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded11"})
			return
		}
	}(file)
	// 创建一个名为 uploaded 的文件
	out, err := os.Create("uploaded.nii.gz")
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

	c.JSON(http.StatusOK, gin.H{"message": "NII file uploaded successfully"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
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
	filePath := "C:\\BUAA\\3rd\\FengRu\\MICCAI-LITS2017\\Task06_Lung\\Task06_Lung\\imagesTr\\lung_010.nii.gz"
	//filePath := "C:\\BUAA\\3rd\\FengRu\\MICCAI-LITS2017\\Task06_Lung\\Task06_Lung\\labelsTr\\lung_001.nii.gz"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read nii.gz file"})
		return
	}

	filename := "file.nii.gz"
	c.Header("Content-Disposition", "attachment; filename="+filename)

	c.Data(http.StatusOK, "application/octet-stream", data)
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
