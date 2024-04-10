package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/stat"
	"math"
	"net/http"
)

type DataprocessController struct {
}

type item struct {
	item_name string `json:"item_name"`
	item_type int    `json:"item_type"`
}

type ChiSquareTest struct {
	VarAValue []int `json:"varAValue"`
	VarBValue []int `json:"varBValue"`
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
	for i := 0; i < len(chisquaredata.VarAValue); i++ {
		sum += chisquaredata.VarAValue[i]
	}
	table[0][1] = float64(sum)
	table[0][0] = float64(len(chisquaredata.VarAValue) - sum)

	sum = 0
	for i := 0; i < len(chisquaredata.VarBValue); i++ {
		sum += chisquaredata.VarBValue[i]
	}
	table[1][1] = float64(sum)
	table[1][0] = float64(len(chisquaredata.VarBValue) - sum)

	var ChiSquare float64

	c_a := table[0][0]
	c_b := table[0][1]
	c_c := table[1][0]
	c_d := table[1][1]
	N := c_a + c_b + c_c + c_d

	if table[0][0] > 5 && table[0][1] > 5 && table[1][0] > 5 && table[1][1] > 5 && table[0][0]+table[0][1]+table[1][0]+table[1][1] >= 40 {
		ChiSquare = stat.ChiSquare(table[0][:], table[1][:])
	} else if table[0][0] >= 1 && table[0][1] >= 1 && table[1][0] >= 1 && table[1][1] >= 1 && table[0][0]+table[0][1]+table[1][0]+table[1][1] >= 40 {
		ChiSquare = (math.Abs(c_a*c_d-c_b*c_c) - N*N/4) * (math.Abs(c_a*c_d-c_b*c_c) - N*N/4) / (c_a + c_b) / (c_c + c_d) / (c_a + c_c) / (c_b + c_d)
	} else {
		ChiSquare = combination(c_a+c_b, c_a) * combination(c_c+c_d, c_c) / combination(N, c_a+c_c)
	}

	c.JSON(http.StatusOK, gin.H{"message": "数据已处理完毕", "chiSquare": ChiSquare})
}

func combination(n, k float64) float64 {
	if k == 0 || k == n {
		return 1
	}
	// 计算阶乘函数 math.factorial() 已废弃，使用 math.Gamma() 来近似计算阶乘
	// 组合数公式：C(n, k) = n! / (k! * (n-k)!)
	// 其中 Gamma 函数的近似值为 (n+1)!
	return math.Gamma(float64(n+1)) / (math.Gamma(float64(k+1)) * math.Gamma(float64(n-k+1)))
}
