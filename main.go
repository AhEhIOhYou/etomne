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

	routes(router)
	err := router.Run(":8091")
	if err != nil {
		log.Fatal(err)
	}
}

func routes(r *gin.Engine) {

	r.LoadHTMLGlob("app/views/*.html")

	r.Any("/", controllers.Index)

	r.GET("/models", controllers.Models)
	r.GET("/models/:id", controllers.Model)

	r.NoRoute(controllers.NotFound)
}
