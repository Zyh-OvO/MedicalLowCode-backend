package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

type DefaultDataController struct {
}

type getOneDataSetJson struct {
	Location string `json:"location" binding:"required"`
}

type AllDataset struct {
	Location        string            `json:"location"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Reference       string            `json:"reference"`
	Licence         string            `json:"licence"`
	Release         string            `json:"release"`
	TensorImageSize string            `json:"tensorImageSize"`
	Modality        map[string]string `json:"modality"`
	Labels          map[string]string `json:"labels"`
	NumTraining     int               `json:"numTraining"`
	NumTest         int               `json:"numTest"`
}

type OneDataset struct {
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Reference       string            `json:"reference"`
	License         string            `json:"licence"`
	Release         string            `json:"release"`
	TensorImageSize string            `json:"tensorImageSize"`
	Modality        map[string]string `json:"modality"`
	Labels          map[string]string `json:"labels"`
	NumTraining     int               `json:"numTraining"`
	NumTest         int               `json:"numTest"`
	Training        []struct {
		Image string `json:"image"`
		Label string `json:"label"`
	} `json:"training"`
	Test []string `json:"test"`
}

func (u DefaultDataController) GetAllDataSet(c *gin.Context) {
	file, err := os.Open("C:\\code\\go\\fengru-backend\\config\\DefaultData.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var datasets []AllDataset
	err = json.NewDecoder(file).Decode(&datasets)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, datasets)
}

func (u DefaultDataController) GetOneDataSet(c *gin.Context) {
	var dataSetJson getOneDataSetJson
	if err := c.ShouldBindJSON(&dataSetJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loc := dataSetJson.Location + "\\dataset.json"
	file, err := os.Open(loc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(file)
	defer file.Close()

	var datasets OneDataset
	err = json.NewDecoder(file).Decode(&datasets)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, datasets)
}
