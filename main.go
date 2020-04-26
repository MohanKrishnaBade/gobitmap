package main

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gobitmap/routes"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os/exec"
)

func init() {

	//react build setup.
	if err := copyNewBuild(); err != nil {
		panic(err)
	}
	//database setup.
	routes.CreateSchema()

	//environment variables setup.
	config()

	//redis Connection
}

func main() {

	// Set the router as the default one shipped with Gin
	router := gin.Default()
	router.NoRoute(customNoRouteHandler)

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./build", true)))

	//Setup route group for the API
	user := router.Group("/user")
	{
		user.GET("/all", routes.FindAll)
		user.GET("/fromCookie", routes.GetUserFromCookie)
		user.POST("/create", routes.Create)
		user.PATCH("/update/:id", routes.Update)
		user.POST("/login", routes.Login)
		user.DELETE("/delete", routes.DeleteAll)
	}

	news := router.Group("/news")
	{
		news.GET("/all", routes.News)
		news.GET("/search", routes.News)
	}

	// Start and run the server
	panic(router.Run(":8000"))
}

//redirect to home page when route not found.
func customNoRouteHandler(c *gin.Context) {
	c.Redirect(301, "/")
}
func copyNewBuild() error {
	err := exec.Command("/bin/sh", "./buildHelper.sh").Start()
	return err
}
