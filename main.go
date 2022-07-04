package main

import (
	"etomne/app/controllers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()

	router.StaticFile("favicon.ico", "assets/images/favicon.ico")
	router.Static("assets/js/", "assets/js/")
	router.Static("assets/css/", "assets/css/")
	router.Static("assets/fonts/", "assets/fonts/")
	router.Static("upload/", "upload/")

	router.LoadHTMLGlob("app/views/*")

	routes(router)
	err := router.Run(":8092")
	if err != nil {
		log.Fatal(err)
	}
}

func routes(r *gin.Engine) {

	r.NoRoute(controllers.NotFound)

	api := r.Group("/api")
	{
		api.GET("/models", controllers.GetModels)
		api.GET("/model/:id", controllers.GetModel)
		api.POST("/model/create", controllers.CreateModel)
		api.DELETE("/models/:id", controllers.DeleteModel)
	}
	//r.GET("/models", controllers.Models)
	//r.GET("/models/:id", controllers.Model)

	r.GET("/upload", controllers.Upload)
	//r.POST("/upload", controllers.UploadModel)

	//r.GET("/delete", controllers.Delete)

	r.Any("/login", controllers.UserLogin)
	r.Any("/reg", controllers.UserReg)
	r.GET("/logout", controllers.UserLogout)

}
