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
	router.Static("upload/", "upload/")

	err := router.Run(":8091")
	if err != nil {
		log.Fatal(err)
	}
}

func Routes(r *gin.Engine) {

	r.LoadHTMLGlob("app/views/*")

	r.GET("/", controllers.Index)
	r.NoRoute(controllers.NotFound)

	r.GET("/models", controllers.Models)
	r.GET("/models/:id", controllers.Model)

	r.GET("/upload", controllers.Upload)
	r.POST("/upload", controllers.UploadModel)

	r.GET("/delete", controllers.Delete)

	r.GET("/login", controllers.UserLogin)
	r.POST("/login", controllers.UserLogin)

	r.GET("/reg", controllers.UserReg)
	r.POST("/reg", controllers.UserReg)

	r.GET("/exit", controllers.UserLogout)
}
