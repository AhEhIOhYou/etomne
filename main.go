package main

import (
	"etomne/app/controllers"
	"github.com/gin-gonic/gin"
	"log"
)

import "github.com/swaggo/gin-swagger" // gin-swagger middleware
import "github.com/swaggo/files"       // swagger embed files

import "etomne/docs"

func main() {
	router := gin.Default()

	router.StaticFile("favicon.ico", "assets/images/favicon.ico")
	router.Static("assets/js/", "assets/js/")
	router.Static("assets/css/", "assets/css/")
	router.Static("assets/fonts/", "assets/fonts/")
	router.Static("upload/", "upload/")

	router.LoadHTMLGlob("app/views/*")

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Models API"
	docs.SwaggerInfo.Description = "This is documentation"
	docs.SwaggerInfo.Version = "0.2"
	docs.SwaggerInfo.Host = "localhost:8092"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	routes(router)
	err := router.Run(":8092")
	if err != nil {
		log.Fatal(err)
	}
}

func routes(r *gin.Engine) {

	r.NoRoute(controllers.NotFound)

	ApiModels := r.Group("/api/models")
	{
		ApiModels.GET("/", controllers.GetModels)
		ApiModels.GET("/:id", controllers.GetModel)
		ApiModels.POST("/create", controllers.CreateModel)
		ApiModels.DELETE("/:id", controllers.DeleteModel)
		ApiModels.PUT(":id", controllers.EditModel)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.GET("/models", controllers.Models)
	//r.GET("/models/:id", controllers.Model)

	//r.GET("/upload", controllers.Upload)
	//r.POST("/upload", controllers.UploadModel)

	//r.GET("/delete", controllers.Delete)

	r.Any("/login", controllers.UserLogin)
	r.Any("/reg", controllers.UserReg)
	r.GET("/logout", controllers.UserLogout)

}
