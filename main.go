package main

import (
	controllers2 "etomne/old/app/controllers"
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

	r.NoRoute(controllers2.NotFound)

	ApiModels := r.Group("/api/models")
	{
		ApiModels.GET("/", controllers2.GetModels)
		ApiModels.GET("/:id", controllers2.GetModel)
		ApiModels.POST("/create", controllers2.CreateModel)
		ApiModels.DELETE("/:id", controllers2.DeleteModel)
		ApiModels.PUT("/:id", controllers2.EditModel)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.GET("/models", controllers.Models)
	//r.GET("/models/:id", controllers.Model)

	//r.GET("/upload", controllers.Upload)
	//r.POST("/upload", controllers.UploadModel)

	//r.GET("/delete", controllers.Delete)

	r.Any("/login", controllers2.UserLogin)
	r.Any("/reg", controllers2.UserReg)
	r.GET("/logout", controllers2.UserLogout)

}
