package router

import (
	"MedicalLowCode-backend/controller/api"
	"MedicalLowCode-backend/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckToken(c *gin.Context) {
	token := c.Request.Header.Get("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{})
		c.Abort()
	} else {
		token, err := model.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
		} else {
			c.Set("token", token)
			c.Next()
		}
	}
}

func ApiRouterInit(router *gin.Engine) {
	apiRouter := router.Group("/api")
	userRouterInit(apiRouter)
	projectManageRouterInit(apiRouter)
	canvasManageRouterInit(apiRouter)
	moduleManageRouterInit(apiRouter)
}

func userRouterInit(router *gin.RouterGroup) {
	userRouter := router.Group("/user")
	userRouter.POST("/getRegisterCode", api.UserController{}.GetRegisterCode)
	userRouter.POST("/register", api.UserController{}.Register)
	userRouter.POST("/login", api.UserController{}.Login)
	userRouter.POST("/getResetCode", api.UserController{}.GetResetCode)
	userRouter.POST("/resetPassword", api.UserController{}.ResetPassword)
	userRouter.POST("/getUserInfo", api.UserController{}.GetUserInfo)
}

func projectManageRouterInit(router *gin.RouterGroup) {
	projectManageRouter := router.Group("/projectManage")
	projectManageRouter.Use(CheckToken)
	projectManageRouter.POST("/newProject", api.ProjectManageController{}.NewProject)
	projectManageRouter.POST("/deleteProject", api.ProjectManageController{}.DeleteProject)
	projectManageRouter.POST("/editProject", api.ProjectManageController{}.EditProject)
	projectManageRouter.POST("/getProjectInfo", api.ProjectManageController{}.GetProjectInfo)
	projectManageRouter.POST("/getProjectList", api.ProjectManageController{}.GetProjectList)
}

func canvasManageRouterInit(router *gin.RouterGroup) {
	canvasManageRouter := router.Group("/canvasManage")
	canvasManageRouter.Use(CheckToken)
	canvasManageRouter.POST("/editCanvas", api.CanvasManageController{}.EditCanvas)
	canvasManageRouter.POST("/getCanvasInfo", api.CanvasManageController{}.GetCanvasInfo)
}

func moduleManageRouterInit(router *gin.RouterGroup) {
	moduleManageRouter := router.Group("/moduleManage")
	moduleManageRouter.Use(CheckToken)
	moduleManageRouter.POST("/getPersonalModules", api.ModuleManageController{}.GetPersonalModules)
	moduleManageRouter.POST("/addPersonalModule", api.ModuleManageController{}.AddPersonalModule)
	moduleManageRouter.POST("/deletePersonalModule", api.ModuleManageController{}.DeletePersonalModule)
	moduleManageRouter.POST("/editPersonalModule", api.ModuleManageController{}.EditPersonalModule)
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
