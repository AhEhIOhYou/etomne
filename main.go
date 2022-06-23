package main

import (
	"etomne/app/controllers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()

	router.Static("assets/js/", "assets/js/")
	router.Static("assets/css/", "assets/css/")
	router.StaticFile("favicon.ico", "assets/images/favicon.ico")

	routes(router)
	err := router.Run(":8091")
	if err != nil {
		log.Fatal(err)
	}
}

func routes(r *gin.Engine) {

	r.LoadHTMLGlob("app/views/*")

	r.Any("/", controllers.Index)
	r.NoRoute(controllers.NotFound)

	r.GET("/models", controllers.Models)
	r.GET("/models/:id", controllers.Model)

	r.GET("/upload", controllers.Upload)
	r.POST("/upload", controllers.UploadModel)

	r.GET("/delete", controllers.Delete)

	r.Any("/login", controllers.UserLogin)
	r.Any("/reg", controllers.UserReg)
	r.GET("/logout", controllers.UserLogout)

}
