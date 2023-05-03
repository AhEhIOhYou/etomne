package main

// gin-swagger middleware
// swagger embed files

import (
	"log"
	"net/http"
	"os"

	"github.com/AhEhIOhYou/etomne/backend/infrastructure/auth"
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/logger"
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/persistence"
	"github.com/AhEhIOhYou/etomne/backend/interfaces"
	"github.com/AhEhIOhYou/etomne/backend/interfaces/filemanager"
	"github.com/AhEhIOhYou/etomne/backend/interfaces/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/AhEhIOhYou/etomne/backend/docs"
)

func init() {
	if err := godotenv.Load("configs/.env"); err != nil {
		log.Print("No .env file found")
	}
}

/*
to swagger.yaml in components add
securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
*/


// @title ETOMNE project
// @version 1.0
// @description This is a 3d model viewer app REST API

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host localhost:8093
// @BasePath /
// @query.collection.format multi
func main() {

	dbHost := os.Getenv("DB_HOST")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	services, err := persistence.NewRepo(dbUser, dbPassword, dbPort, dbHost, dbName)
	if err != nil {
		log.Fatal(err)
	}

	redisService, err := auth.NewRedisDb(redisHost, redisPort, redisPassword)
	if err != nil {
		log.Fatal(err)
	}

	tk := auth.NewToken()
	fm := filemanager.NewFileUpload()

	users := interfaces.NewUsers(services.User, fm, redisService.Auth, tk)
	models := interfaces.NewModel(services.Model, services.User, services.File, fm, redisService.Auth, tk)
	files := interfaces.NewFile(services.Model, services.User, services.File, fm, redisService.Auth, tk)
	authenticate := interfaces.NewAuthenticate(services.User, redisService.Auth, tk)
	index := interfaces.Index

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.FrontStaticMiddleware())
	r.Use(middleware.UploadStaticMiddleware())

	r.Any("/", index)

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})

	u := r.Group("api/users")
	{
		u.POST("", users.SaveUser)
		u.GET("", users.GetUsers)
		u.POST("/login", authenticate.Login)
		u.POST("/logout", middleware.AuthMiddleware(), authenticate.Logout)
		u.POST("/refresh", authenticate.Refresh)
		u.POST("/addfile", middleware.AuthMiddleware(), users.SaveUserPhoto)
	}

	m := r.Group("api/model")
	{
		m.POST("", middleware.AuthMiddleware(), models.SaveModel)
		m.PUT("/:model_id", middleware.AuthMiddleware(), models.UpdateModel)
		m.GET("/:model_id", models.GetModel)
		m.DELETE("/:model_id", middleware.AuthMiddleware(), models.DeleteModel)
		m.GET("", models.GetAllModel)
		m.POST("/addfile", middleware.AuthMiddleware(), models.SaveModelFile)
	}

	f := r.Group("api/file")
	{
		f.POST("", files.SaveFile)
		f.DELETE("/:file_id", files.RemoveFile)
	}

	logger.WriteLog(logger.Info, "THE SERVER HAS BEEN SUCCESSFULLY STARTED")

	// docs route
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.LoadHTMLFiles("frontend/dist/index.html")

	log.Fatal(r.Run(":" + "8093"))
}
