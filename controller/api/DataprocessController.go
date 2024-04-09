package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/stat"
	"net/http"
)

type DataprocessController struct {
}

type item struct {
	item_name string `json:"item_name"`
	item_type int    `json:"item_type"`
}

type ChiSquareTest struct {
	varAValue []int `json:"varAValue"`
	varBValue []int `json:"varBValue"`
}

type AnalysisResult struct {
	ChiSquare              float64   `json:"chi_square"`
	RegressionCoefficients []float64 `json:"regression_coefficients"`
}

func (u DataprocessController) MedicalDataAnalysisHandler(c *gin.Context) {
	// 读取传输的数据
	info := c.Request.FormValue("data")
	fmt.Println(info)
	// 解析 JSON 数据到结构体
	var chisquaredata ChiSquareTest
	if err := json.Unmarshal([]byte(info), &chisquaredata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON info"})
		return
	}
	var table [2][2]float64
	sum := 0
	for i := 0; i < len(chisquaredata.varAValue); i++ {
		sum += chisquaredata.varAValue[i]
	}
	table[0][1] = float64(sum)
	table[0][0] = float64(len(chisquaredata.varAValue) - sum)

	sum = 0
	for i := 0; i < len(chisquaredata.varBValue); i++ {
		sum += chisquaredata.varBValue[i]
	}
	table[1][1] = float64(sum)
	table[1][0] = float64(len(chisquaredata.varBValue) - sum)

	ChiSquare := stat.ChiSquare(table[0][:], table[1][:])

	c.JSON(http.StatusOK, gin.H{"message": "数据已处理完毕", "chiSquare": ChiSquare})
}
