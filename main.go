package main

// gin-swagger middleware
// swagger embed files

import (
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/auth"
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/persistence"
	"github.com/AhEhIOhYou/etomne/backend/interfaces"
	"github.com/AhEhIOhYou/etomne/backend/interfaces/filemanager"
	"github.com/AhEhIOhYou/etomne/backend/interfaces/middleware"
	"github.com/AhEhIOhYou/etomne/docs"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
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

openapi: 3.0.1
info:
  title: "aaaa"
  version: 1.0.0
servers:
- url: https://modelshowtime.serdcebolit.ru/api
*/
func main() {

	//Main DB vars
	dbHost := os.Getenv("DB_HOST")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	//Redis vars
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	services, err := persistence.NewRepo(dbUser, dbPassword, dbPort, dbHost, dbName)
	if err != nil {
		log.Fatal(err)
	}

	//services.Migrate()

	redisService, err := auth.NewRedisDb(redisHost, redisPort, redisPassword)
	if err != nil {
		log.Fatal(err)
	}

	tk := auth.NewToken()
	fm := filemanager.NewFileUpload()

	users := interfaces.NewUsers(services.User, redisService.Auth, tk)
	models := interfaces.NewModel(services.Model, services.User, services.File, fm, redisService.Auth, tk)
	files := interfaces.NewFile(services.Model, services.User, services.File, fm, redisService.Auth, tk)
	authenticate := interfaces.NewAuthenticate(services.User, redisService.Auth, tk)
	comments := interfaces.NewComment(services.Model, services.User, services.Comment, redisService.Auth, tk)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	u := r.Group("api/users")
	{
		u.POST("", users.SaveUser)
		u.GET("", users.GetUsers)
		u.GET("/:user_id", users.GetUser)
		u.POST("/login", authenticate.Login)
		u.POST("/logout", middleware.AuthMiddleware(), authenticate.Logout)
		u.POST("/refresh", authenticate.Refresh)
	}

	m := r.Group("api/model")
	{
		m.POST("", middleware.AuthMiddleware(), models.SaveModel)
		m.PUT("/:model_id", middleware.AuthMiddleware(), models.UpdateModel)
		m.GET("/:model_id", models.GetModel)
		m.DELETE("/:model_id", middleware.AuthMiddleware(), models.DeleteModel)
		m.GET("", models.GetAllModel)
	}

	f := r.Group("api/file")
	{
		f.POST("", files.SaveFile)
		f.DELETE("/:file_id", files.RemoveFile)
	}

	c := r.Group("api/comment")
	{
		c.GET("", comments.GetComments)
		c.POST("", comments.SaveComment)
	}

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Models API"
	docs.SwaggerInfo.Description = "This is documentation"
	docs.SwaggerInfo.Version = "0.3"
	docs.SwaggerInfo.Host = "localhost:8093/api"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(r.Run(":" + "8093"))
}
