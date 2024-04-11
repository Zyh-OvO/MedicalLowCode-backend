package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/stat"
	"math"
	"math/rand"
	"net/http"
)

type DataprocessController struct {
}

type ChiSquareTest struct {
	VarAValue []int `json:"varAValue"`
	VarBValue []int `json:"varBValue"`
	BTypeNum  int   `json:"BTypeNum"`
}

type K_means_data struct {
	K         int         `json:"k"`
	Data      [][]float64 `json:"data"`
	ItorTimes int         `json:"itorTimes"`
}

func (u DataprocessController) ChisquarHandler(c *gin.Context) {
	// 读取传输的数据
	info := c.Request.FormValue("data")
	fmt.Println(info)
	// 解析 JSON 数据到结构体
	var chisquaredata ChiSquareTest
	if err := json.Unmarshal([]byte(info), &chisquaredata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON info"})
		return
	}
	var table [2][]float64

	table[0] = make([]float64, chisquaredata.BTypeNum)
	table[1] = make([]float64, chisquaredata.BTypeNum)

	for i := 0; i < len(chisquaredata.VarAValue); i++ {
		table[chisquaredata.VarAValue[i]][chisquaredata.VarBValue[i]]++
	}

	var ChiSquare float64

	c_a := table[0][0]
	c_b := table[0][1]
	c_c := table[1][0]
	c_d := table[1][1]
	N := c_a + c_b + c_c + c_d

	if chisquaredata.BTypeNum > 2 || table[0][0] > 5 && table[0][1] > 5 && table[1][0] > 5 && table[1][1] > 5 && table[0][0]+table[0][1]+table[1][0]+table[1][1] >= 40 {
		for i := 0; i < 2; i++ {
			for j := 0; j < chisquaredata.BTypeNum; j++ {
				if table[i][j] == 0 {
					table[i][j] = 0.01
				}
			}
		}
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

func (u DataprocessController) k_means_func(c *gin.Context) {
	info := c.Request.FormValue("data")
	fmt.Println(info)
	// 解析 JSON 数据到结构体
	var k_means_data K_means_data
	if err := json.Unmarshal([]byte(info), &k_means_data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON info"})
		return
	}
	var types []int // 取值是0-K-1
	types = make([]int, len(k_means_data.Data[0]))

	var points_location [][]float64
	points_location = make([][]float64, k_means_data.K)

	for i := 0; i < len(points_location); i++ {
		points_location[i] = make([]float64, len(k_means_data.Data))
	}

	// 初始化k个点
	for i := 0; i < len(points_location); i++ {
		for j := 0; j < len(points_location[0]); j++ {
			points_location[i][j] = rand.Float64()
		}
	}

	for i := 0; i < k_means_data.K; i++ {
		var points_location_temp [][]float64
		points_location_temp = make([][]float64, len(k_means_data.Data[0]))
		for i := 0; i < len(points_location); i++ {
			points_location_temp[i] = make([]float64, len(k_means_data.Data))
		}

		num := make([]int, len(points_location))
		//首先计算属于每个中心点的点
		for j := 0; j < len(k_means_data.Data[0]); j++ { // 对于每一行
			var this_min float64 // 当前值
			var dis_min float64  //最小值
			var min_point int    // 最小值对应的哪一行(代表一个点)
			dis_min = 2
			for k := 0; k < k_means_data.K; k++ { //对于每一个聚点
				//计算第j个点对于第k个聚点的距离
				this_min = 0
				for l := 0; l < len(k_means_data.Data); l++ { // 对于第l维，也就是第k列
					this_min += (k_means_data.Data[l][j] - points_location[k][j]) * (k_means_data.Data[l][j] - points_location[k][j])
				}
				if this_min < dis_min {
					dis_min = this_min
					min_point = k
				}
			}
			num[min_point]++
			types[j] = min_point
			for k := 0; k < len(k_means_data.Data); k++ {
				points_location_temp[min_point][k] += k_means_data.Data[k][j]
			}
		}
		//然后移动中心点
		for j := 0; j < len(points_location); j++ {
			// 对于第j个聚点
			if num[j] == 0 {
				for k := 0; k < len(points_location[0]); k++ {
					points_location[j][k] = rand.Float64()
				}
			} else {
				for k := 0; k < len(points_location[0]); k++ {
					points_location[j][k] = points_location_temp[j][k] / float64(num[j])
				}
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "数据已处理完毕", "chiSquare": types})
}
