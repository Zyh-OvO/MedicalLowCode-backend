package router

import (
	"MedicalLowCode-backend/controller/api"
	"MedicalLowCode-backend/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckToken(c *gin.Context) {
	token := c.Request.Header.Get("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{})
		c.Abort()
	} else {
		token, err := util.ParseToken(token)
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

func CheckCornerStoneToken(c *gin.Context) {
	tk := c.Query("token")
	if tk == "" || tk == "null" || tk == "111" {
		tk = c.Request.Header.Get("token")
		if tk == "" {
			//	c.JSON(http.StatusUnauthorized, gin.H{})
			//	c.Abort()
			//	return
			// TODO: implement real token
			token, err := util.GiveStaticToken()
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				c.Abort()
			} else {
				c.Set("token", token)
				c.Next()
				fmt.Println("token acquired")
			}
		}
	}
	token, err := util.ParseToken(tk)
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

func ApiRouterInit(router *gin.Engine) {
	apiRouter := router.Group("/api")
	userRouterInit(apiRouter)
	projectManageRouterInit(apiRouter)
	canvasManageRouterInit(apiRouter)
	moduleManageRouterInit(apiRouter)
	projectDevelopRouterInit(apiRouter)
	defaultModuleManageRouterInit(apiRouter)
	defaultDataManageRouterInit(apiRouter)
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
}

func projectDevelopRouterInit(router *gin.RouterGroup) {
	projectDevelopRouter := router.Group("/projectDevelop")
	projectDevelopRouter.Use(CheckToken)
	projectDevelopRouter.POST("/exportCode", api.ProjectDevelopController{}.ExportCode)
	projectDevelopRouter.POST("/submitTrainingTask", api.ProjectDevelopController{}.SubmitTrainingTask)
	projectDevelopRouter.POST("/submitReasoningTask", api.ProjectDevelopController{}.SubmitReasoningTask)
}

func defaultModuleManageRouterInit(router *gin.RouterGroup) {
	defaultModuleManageRouter := router.Group("/defaultModule")
	defaultModuleManageRouter.Use(CheckCornerStoneToken)
	defaultModuleManageRouter.POST("/imageTest", api.DefaultModelController{}.ImageTest)
	defaultModuleManageRouter.POST("/niiTest", api.DefaultModelController{}.NiiTest)
	defaultModuleManageRouter.POST("/getImages", api.DefaultModelController{}.ReturnMultipleImages)
	defaultModuleManageRouter.GET("/returnNiiGzFile", api.DefaultModelController{}.ReturnNiiGzFile)
	defaultModuleManageRouter.GET("/returnSegFile", api.DefaultModelController{}.ReturnSegFile)
	defaultModuleManageRouter.GET("/returnSegData", api.DefaultModelController{}.GetNoneZeroLocation)
	defaultModuleManageRouter.GET("/getDim", api.DefaultModelController{}.DimTest)
}

func defaultDataManageRouterInit(router *gin.RouterGroup) {
	defaultDataManageRouter := router.Group("/defaultData")
	defaultDataManageRouter.POST("/getAllDataSet", api.DefaultDataController{}.GetAllDataSet)
	defaultDataManageRouter.POST("/getOneDataSet", api.DefaultDataController{}.GetOneDataSet)
}
