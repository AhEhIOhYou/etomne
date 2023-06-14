package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AhEhIOhYou/etomne/pkg/server/constants"
	"github.com/AhEhIOhYou/etomne/pkg/server/infrastructure/auth"
	"github.com/AhEhIOhYou/etomne/pkg/server/infrastructure/logger"
	"github.com/AhEhIOhYou/etomne/pkg/server/infrastructure/persistence"
	"github.com/AhEhIOhYou/etomne/pkg/server/interfaces"
	"github.com/AhEhIOhYou/etomne/pkg/server/interfaces/filemanager"
	"github.com/AhEhIOhYou/etomne/pkg/server/interfaces/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/AhEhIOhYou/etomne/docs"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		logger.WriteLog(logger.Error, constants.ServerNotFoundEnvFile)
	}
}

//	@title			ETOMNE project
//	@version		2.0
//	@description	This is the REST API of the application for viewing 3d models

//	@contact.name	API Support
//	@contact.email	email@man.you

//	@host						localhost:8095
//	@BasePath					/api
//	@query.collection.format	multi

//	securityDefinitions.apikey	BearerAuth
//	in							header
//	name						Authorization
func Start() {

	dbHost := os.Getenv("DB_HOST")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASS")

	logger.WriteLog(logger.Info, constants.DatabaseConnectionStart)

	services, err := persistence.NewRepo(dbUser, dbPassword, dbPort, dbHost, dbName)
	if err != nil {
		logger.WriteLog(logger.Error, fmt.Sprintf(constants.DatabaseConnectionError, err))
		return
	}

	redisService, err := auth.NewRedisDb(redisHost, redisPort, redisPassword)
	if err != nil {
		logger.WriteLog(logger.Error, fmt.Sprintf(constants.DatabaseConnectionError, err))
		return
	}

	logger.WriteLog(logger.Info, constants.DatabaseConnectionSuccess)

	tk := auth.NewToken()
	fm := filemanager.NewFileUpload()

	users := interfaces.NewUsers(services.User, fm, redisService.Auth, tk)
	models := interfaces.NewModel(services.Model, services.User, services.File, fm, redisService.Auth, tk)
	files := interfaces.NewFile(services.Model, services.User, services.File, fm, redisService.Auth, tk)
	authenticate := interfaces.NewAuthenticate(services.User, redisService.Auth, tk)
	index := interfaces.Index

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.FrontStaticMiddleware())
	router.Use(middleware.UploadStaticMiddleware())

	router.Any("/", index)

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	u := router.Group("api/users")
	{
		u.POST("", users.SaveUser)
		u.GET("/:user_id", users.GetUserByID)
		u.POST("/login", authenticate.Login)
		u.GET("/logout", middleware.AuthMiddleware(), authenticate.Logout)
		u.POST("/update/:user_id", middleware.AuthMiddleware(), users.UpdateUser)
		u.POST("/update/admin/:user_id", middleware.AuthMiddleware(), users.UpdateUserAdminRights)
		u.DELETE("/:user_id", middleware.AuthMiddleware(), users.DeleteUser)
		u.POST("/refresh", authenticate.Refresh)
	}

	m := router.Group("api/model")
	{
		m.POST("", middleware.AuthMiddleware(), models.SaveModel)
		m.POST("/update/:model_id", middleware.AuthMiddleware(), models.UpdateModel)
		m.GET("/:model_id", models.GetModel)
		m.DELETE("/:model_id", middleware.AuthMiddleware(), models.DeleteModel)
		m.GET("", models.GetModelList)
		m.POST("/addfile/:model_id", middleware.AuthMiddleware(), models.SaveModelFile)
	}

	f := router.Group("api/file")
	{
		f.POST("", files.SaveFile)
		f.DELETE("/:file_id", files.RemoveFile)
	}

	logger.WriteLog(logger.Info, constants.ServerStartSuccess)

	logger.WriteLog(logger.Error,
		fmt.Sprintf(
			constants.ServerStartErr,
			router.Run(":8080")),
	)
}
