package api

import (
	"MedicalLowCode-backend/model"
	"MedicalLowCode-backend/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserController struct{}

type getRegisterCodeJson struct {
	UserEmail string `json:"userEmail" binding:"required"`
}

type registerJson struct {
	RegisterCode string `json:"registerCode" binding:"required"`
	UserEmail    string `json:"userEmail" binding:"required"`
	UserName     string `json:"userName" binding:"required"`
	UserPassword string `json:"userPassword" binding:"required"`
}

type loginJson struct {
	UserEmail    string `json:"userEmail" binding:"required"`
	UserPassword string `json:"userPassword" binding:"required"`
}

type getResetCodeJson struct {
	UserEmail string `json:"userEmail" binding:"required"`
}

type resetPassword struct {
	NewPassword string `json:"newPassword" binding:"required"`
	ResetCode   string `json:"resetCode" binding:"required"`
	UserEmail   string `json:"userEmail" binding:"required"`
}

func (u UserController) GetRegisterCode(c *gin.Context) {
	var json getRegisterCodeJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//校验邮箱
	if !util.CheckEmail(json.UserEmail) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱格式不正确"})
		return
	}

	//检查用户是否存在
	user := model.QueryUser(json.UserEmail)
	if user != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户已存在"})
		return
	}

	//生成6位验证码
	code := util.GenRandomString(6)

	//操作数据库
	model.AddRegisterCode(json.UserEmail, code)

	//发送邮件
	toEmail := json.UserEmail
	subject := "欢迎注册医学低代码平台"
	body := fmt.Sprintf("欢迎注册医学低代码平台，您的验证码为：%s，验证码有效期为10分钟，如非您个人操作，请忽略", code)
	util.SendEmail(toEmail, subject, body)
	c.JSON(http.StatusOK, gin.H{})
}

func (u UserController) Register(c *gin.Context) {
	var json registerJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//校验邮箱
	if !util.CheckEmail(json.UserEmail) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱格式不正确"})
		return
	}
	//检查用户是否存在
	user := model.QueryUser(json.UserEmail)
	if user != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户已存在"})
		return
	}
	//校验名字和密码
	if !util.CheckUserName(json.UserName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名格式不正确"})
		return
	}
	if !util.CheckPassword(json.UserPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码格式不正确"})
		return
	}
	//校验验证码
	dbCode := model.QueryRegisterCode(json.UserEmail)
	if dbCode != json.RegisterCode {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	model.RegisterNewUser(json.UserEmail, json.UserName, json.UserPassword)
	c.JSON(http.StatusOK, gin.H{})
}

func (u UserController) Login(c *gin.Context) {
	var json loginJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//校验邮箱
	if !util.CheckEmail(json.UserEmail) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱格式不正确"})
		return
	}
	//检验用户是否存在
	user := model.QueryUser(json.UserEmail)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
		return
	}
	//校验密码
	if !util.CheckBcryptPassword(json.UserPassword, user.UserPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码错误"})
		return
	}
	token, err := util.SignToken(user.UserId, user.UserName, user.UserEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (u UserController) GetResetCode(c *gin.Context) {
	var json getResetCodeJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//校验邮箱
	if !util.CheckEmail(json.UserEmail) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱格式不正确"})
		return
	}

	//检查用户是否存在
	user := model.QueryUser(json.UserEmail)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
		return
	}

	//生成6位验证码
	code := util.GenRandomString(6)

	//操作数据库
	model.AddResetCode(json.UserEmail, code)

	//发送邮件
	toEmail := json.UserEmail
	subject := "医学低代码平台密码重置"
	body := fmt.Sprintf("您正在重置密码，验证码为：%s，验证码有效期为10分钟，如非您个人操作，请忽略", code)
	util.SendEmail(toEmail, subject, body)

	c.JSON(http.StatusOK, gin.H{})
}

func (u UserController) ResetPassword(c *gin.Context) {
	var json resetPassword
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//校验邮箱
	if !util.CheckEmail(json.UserEmail) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱格式不正确"})
		return
	}
	//检查用户是否存在
	user := model.QueryUser(json.UserEmail)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
		return
	}
	//校验密码
	if !util.CheckPassword(json.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码格式不正确"})
		return
	}
	//校验验证码
	dbCode := model.QueryResetCode(json.UserEmail)
	if dbCode != json.ResetCode {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}
	model.ResetPassword(json.UserEmail, json.NewPassword)
	c.JSON(http.StatusOK, gin.H{})
}

func (u UserController) GetUserInfo(c *gin.Context) {
	tokenString := c.Request.Header.Get("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	token, err := util.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"userId":    strconv.Itoa(token.UserId),
		"userName":  token.UserName,
		"userEmail": token.UserEmail,
	})
}
