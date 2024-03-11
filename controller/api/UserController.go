package api

import (
	"MedicalLowCode-backend/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct{}

type getRegisterCodeJson struct {
	UserEmail string `json:"userEmail"`
}

type registerJson struct {
	RegisterCode string `json:"registerCode"`
	UserEmail    string `json:"userEmail"`
	UserName     string `json:"userName"`
	UserPassword string `json:"userPassword"`
}

type loginJson struct {
	UserEmail    string `json:"userEmail"`
	UserPassword string `json:"userPassword"`
}

type getResetCodeJson struct {
	UserEmail string `json:"userEmail"`
}

type resetPassword struct {
	NewPassword string `json:"newPassword"`
	ResetCode   string `json:"resetCode"`
	UserEmail   string `json:"userEmail"`
}

func (u UserController) GetRegisterCode(c *gin.Context) {
	var json getRegisterCodeJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//校验邮箱
	if !model.CheckEmail(json.UserEmail) {
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
	code := model.GenRandomString(6)

	//操作数据库
	model.AddRegisterCode(json.UserEmail, code)

	//发送邮件
	toEmail := json.UserEmail
	subject := "欢迎注册医学低代码平台"
	body := fmt.Sprintf("欢迎注册医学低代码平台，您的验证码为：%s，验证码有效期为10分钟，如非您个人操作，请忽略", code)
	model.SendEmail(toEmail, subject, body)
	fmt.Println("123123123")
	c.JSON(http.StatusOK, gin.H{})
}

func (u UserController) Register(c *gin.Context) {
	var json registerJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//校验邮箱
	if !model.CheckEmail(json.UserEmail) {
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
	if !model.CheckUserName(json.UserName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名格式不正确"})
		return
	}
	if !model.CheckPassword(json.UserPassword) {
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
	if !model.CheckEmail(json.UserEmail) {
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
	if !model.CheckBcryptPassword(json.UserPassword, user.UserPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码错误"})
		return
	}
	token, err := model.SignToken(user.UserId, user.UserName, user.UserEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (u UserController) GetResetCode(c *gin.Context) {
	fmt.Println(c)
	fmt.Println(u)
	var json getResetCodeJson
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//校验邮箱
	if !model.CheckEmail(json.UserEmail) {
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
	code := model.GenRandomString(6)

	//操作数据库
	model.AddResetCode(json.UserEmail, code)

	//发送邮件
	toEmail := json.UserEmail
	subject := "医学低代码平台密码重置"
	body := fmt.Sprintf("您正在重置密码，验证码为：%s，验证码有效期为10分钟，如非您个人操作，请忽略", code)
	model.SendEmail(toEmail, subject, body)

	c.JSON(http.StatusOK, gin.H{})
}

func (u UserController) ResetPassword(c *gin.Context) {
	var json resetPassword
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//校验邮箱
	if !model.CheckEmail(json.UserEmail) {
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
	if !model.CheckPassword(json.NewPassword) {
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
	token, err := model.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"userId":    token.UserId,
		"userName":  token.UserName,
		"userEmail": token.UserEmail,
	})
}
