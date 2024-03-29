package router

import (
	"MedicalLowCode-backend/controller/api"
	"MedicalLowCode-backend/util"
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

func ApiRouterInit(router *gin.Engine) {
	apiRouter := router.Group("/api")
	userRouterInit(apiRouter)
	projectManageRouterInit(apiRouter)
	canvasManageRouterInit(apiRouter)
	moduleManageRouterInit(apiRouter)
	projectDevelopRouterInit(apiRouter)
	defaultModuleManageRouterInit(apiRouter)
	fileManageRouterInit(apiRouter)
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
	defaultModuleManageRouter.POST("/imageTest", api.CtModelController{}.ImageTest)
	defaultModuleManageRouter.POST("/niiTest", api.CtModelController{}.NiiTest)
	defaultModuleManageRouter.POST("/getImages", api.CtModelController{}.ReturnMultipleImages)
	defaultModuleManageRouter.GET("/returnNiiGzFile", api.CtModelController{}.ReturnNiiGzFile)
	defaultModuleManageRouter.GET("/returnSegFile", api.CtModelController{}.ReturnSegFile)
	defaultModuleManageRouter.GET("/returnSegData", api.CtModelController{}.GetNoneZeroLocation)
	defaultModuleManageRouter.GET("/getDim", api.CtModelController{}.DimTest)
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
