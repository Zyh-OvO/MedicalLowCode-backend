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

type submitTaskJson struct {
	ProjectId string `json:"projectId" binding:"required"`
	TaskName  string `json:"taskName" binding:"required"`
}

type submitReasoningTaskJson struct {
	ProjectId string `json:"projectId" binding:"required"`
}

type getTaskListJson struct {
	ProjectId string `json:"projectId" binding:"required"`
}

type stopTaskJson struct {
	TaskId string `json:"taskId" binding:"required"`
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
	var json submitTaskJson
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
	task := model.NewTask(token.UserId, projectId, json.TaskName)
	code := exportCode.ExportCode(canvas.CanvasContent)
	//本地执行
	go RunTask(token.UserId, task.TaskId, code)
	c.JSON(http.StatusOK, gin.H{})
}

func RunTask(userId int, taskId int, code string) {
	cmd1 := "echo '" + code + "' > ./taskCode/task_" + strconv.Itoa(taskId) + ".py"
	cmd2 := "python ./taskCode/task_" + strconv.Itoa(taskId) + ".py > ./taskLog/task_" + strconv.Itoa(taskId) + ".log"
	allCmd := cmd1 + " && " + cmd2
	if err := exec.Command("bash", "-c", allCmd).Run(); err != nil {
		model.SetTaskStatus(userId, taskId, false)
	}
	model.SetTaskStatus(userId, taskId, true)
}

func (p ProjectDevelopController) GetTaskList(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json getTaskListJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectId, _ := strconv.Atoi(json.ProjectId)
	tasks := model.QueryTaskList(token.UserId, projectId)
	var taskList []gin.H
	for _, task := range tasks {
		var endTime int64
		if task.EndTime == nil {
			endTime = -1
		} else {
			endTime = task.EndTime.Unix()
		}
		taskList = append(taskList, gin.H{
			"taskId":       strconv.Itoa(task.TaskId),
			"taskName":     task.TaskName,
			"submitTime":   task.SubmitTime.Unix(),
			"endTime":      endTime,
			"isSuccessful": task.IsSuccessful,
			"logFilePath":  task.LogFilePath,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"taskList": taskList,
	})
}

func (p ProjectDevelopController) StopTask(c *gin.Context) {
	token := c.MustGet("token").(*util.Token)
	var json stopTaskJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taskId, _ := strconv.Atoi(json.TaskId)
	task := model.SetTaskStatus(token.UserId, taskId, false)
	if task == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "任务不存在"})
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
	code := ""
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
