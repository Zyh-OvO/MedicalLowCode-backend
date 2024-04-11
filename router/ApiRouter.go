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
	if tk == "" || tk == "null" {
		tk = c.Param("token")
		fmt.Println("param_token:", tk)
	}
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
				fmt.Println("token acquired1")
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
		fmt.Println("token acquired2")
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
	fileManageRouterInit(apiRouter)
	defaultDataManageRouterInit(apiRouter)
	websocketRouterInit(apiRouter)
	dataprocessRouterInit(apiRouter)
}

func dataprocessRouterInit(router *gin.RouterGroup) {
	dataprocessRouter := router.Group("/dataprocess")
	dataprocessRouter.POST("/chiSquareTest", api.DataprocessController{}.ChisquarHandler)
	dataprocessRouter.POST("/kMeans", api.DataprocessController{}.K_means_func)
}

func websocketRouterInit(router *gin.RouterGroup) {
	websocketRouter := router.Group("/ws")
	websocketRouter.GET("/inference", api.DefaultModelController{}.WebsocketHandler)
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
	projectDevelopRouter.POST("/submitTrainingTask", api.ProjectDevelopController{}.SubmitTask)
	projectDevelopRouter.POST("/getTaskList", api.ProjectDevelopController{}.GetTaskList)
	projectDevelopRouter.POST("/stopTask", api.ProjectDevelopController{}.StopTask)
	//projectDevelopRouter.POST("/submitReasoningTask", api.ProjectDevelopController{}.SubmitReasoningTask)
}

func defaultModuleManageRouterInit(router *gin.RouterGroup) {
	defaultModuleManageRouter := router.Group("/defaultModule")
	defaultModuleManageRouter.Use(CheckCornerStoneToken)
	//defaultModuleManageRouter.Use(cors.Default())
	defaultModuleManageRouter.POST("/imageTest", api.DefaultModelController{}.ImageTest)
	defaultModuleManageRouter.POST("/niiTest", api.DefaultModelController{}.UploadNiiGzFile)
	defaultModuleManageRouter.POST("/getImages", api.DefaultModelController{}.ReturnMultipleImages)
	defaultModuleManageRouter.GET("/returnNiiGzFile/:token/:id", api.DefaultModelController{}.ReturnNiiGzFile)
	defaultModuleManageRouter.GET("/returnSegData/:token/:id", api.DefaultModelController{}.GetNonZeroLocation)
	defaultModuleManageRouter.POST("/postModelInfo", api.NnunetModelController{}.SetModelInfo)
	defaultModuleManageRouter.GET("/getModelInfoList/:token", api.NnunetModelController{}.GetModelList)
}

func defaultDataManageRouterInit(router *gin.RouterGroup) {
	defaultDataManageRouter := router.Group("/defaultData")
	defaultDataManageRouter.POST("/getAllDataSet", api.DefaultDataController{}.GetAllDataSet)
	defaultDataManageRouter.POST("/getOneDataSet", api.DefaultDataController{}.GetOneDataSet)
}

func fileManageRouterInit(router *gin.RouterGroup) {
	fileManageRouter := router.Group("/fileManage")
	fileManageRouter.Use(CheckToken)
	fileManageRouter.POST("/getDirContent", api.FileManageController{}.GetDirContent)
	fileManageRouter.POST("/uploadFile", api.FileManageController{}.UploadFile)
	fileManageRouter.POST("/deleteFile", api.FileManageController{}.DeleteFile)
	fileManageRouter.POST("/renameFile", api.FileManageController{}.RenameFile)
	//fileManageRouter.POST("/copyFile", api.FileManageController{}.CopyFile)
	fileManageRouter.POST("/newDirectory", api.FileManageController{}.NewDirectory)
	fileManageRouter.POST("/deleteDirectory", api.FileManageController{}.DeleteDirectory)
	fileManageRouter.POST("/renameDirectory", api.FileManageController{}.RenameDirectory)
	fileManageRouter.POST("/getRootDir", api.FileManageController{}.GetRootDir)
}
