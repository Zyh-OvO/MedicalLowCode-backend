package api

import (
	"MedicalLowCode-backend/exportCode"
	"MedicalLowCode-backend/model"
	"MedicalLowCode-backend/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

type ProjectDevelopController struct{}

type exportCodeJson struct {
	ProjectId string `json:"projectId" binding:"required"`
}

type submitTrainingTaskJson struct {
	ProjectId string `json:"projectId" binding:"required"`
}

type submitReasoningTaskJson struct {
	ProjectId string `json:"projectId" binding:"required"`
}

func (p ProjectDevelopController) ExportCode(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json exportCodeJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectId, _ := strconv.Atoi(json.ProjectId)
	canvas := model.QueryCanvas(token.UserId, projectId)
	if canvas == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "项目不存在"})
		return
	}
	code := exportCode.ExportCode(canvas.CanvasContent)
	fmt.Println(code)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
	})
}

func (p ProjectDevelopController) SubmitTask(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json submitTrainingTaskJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectId, _ := strconv.Atoi(json.ProjectId)
	canvas := model.QueryCanvas(token.UserId, projectId)
	if canvas == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "项目不存在"})
		return
	}
	//todo
	code := exportCode.ExportCode(canvas.CanvasContent)
	//ssh
	vpnCmd1 := exec.Command("sh", "-c", "~/lowcode/actvpn.sh")
	err := vpnCmd1.Run()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sshConfig, err := util.GetSSHConfig()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client, err := ssh.Dial("tcp", "192.168.5.201:22", sshConfig)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer client.Close()
	session1, err := client.NewSession()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer session1.Close()
	userDirPath := fmt.Sprintf("user_%d/project_%d/%s/", token.UserId, projectId, util.UnixToDate(time.Now().Unix()))
	cmd1 := "mkdir -p ~/lowcode/" + userDirPath
	cmd2 := "cd ~/lowcode/" + userDirPath
	cmd3 := "echo '" + code + "' > train.py"
	cmd4 := "cp ~/lowcode/myscript ./myscript"
	cmd5 := "sbatch --gres=gpu:V100:1 ./myscript"
	allCmd := cmd1 + " && " + cmd2 + " && " + cmd3 + " && " + cmd4 + " && " + cmd5
	if err := session1.Run(allCmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vpnCmd2 := exec.Command("sh", "-c", "~/lowcode/stopactvpn.sh")
	err = vpnCmd2.Run()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (p ProjectDevelopController) SubmitReasoningTask(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json submitReasoningTaskJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectId, _ := strconv.Atoi(json.ProjectId)
	canvas := model.QueryCanvas(token.UserId, projectId)
	if canvas == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "项目不存在"})
		return
	}
	//todo

	c.JSON(http.StatusOK, gin.H{})
}
