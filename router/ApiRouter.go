package router

import (
	"MedicalLowCode-backend/controller/api"
	"github.com/gin-gonic/gin"
)

func ApiRouterInit(router *gin.Engine) {
	apiRouter := router.Group("/api")
	userRouterInit(apiRouter)
}

func userRouterInit(router *gin.RouterGroup) {
	router.POST("/user/getRegisterCode", api.UserController{}.GetRegisterCode)
	router.POST("/user/register", api.UserController{}.Register)
	router.POST("/user/login", api.UserController{}.Login)
	router.POST("/user/getResetCode", api.UserController{}.GetResetCode)
	router.POST("/user/resetPassword", api.UserController{}.ResetPassword)
	router.POST("/user/getUserInfo", api.UserController{}.GetUserInfo)
	router.POST("/user/test", api.CtModelController{}.Test)
	router.POST("/user/imageTest", api.CtModelController{}.ImageTest)
	router.POST("/user/niiTest", api.CtModelController{}.NiiTest)
	router.POST("/user/getImages", api.CtModelController{}.ReturnMultipleImages)
	router.GET("/user/returnNiiGzFile", api.CtModelController{}.ReturnNiiGzFile)
	router.GET("/user/returnSegFile", api.CtModelController{}.ReturnSegFile)
	router.GET("/user/returnSegData", api.CtModelController{}.GetNoneZeroLocation)
	router.GET("/user/getDim", api.CtModelController{}.DimTest)
}
