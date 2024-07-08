package main

import (
	"log"
	"os"

	"github.com/aaravmahajanofficial/ecommerce-project/controllers"
	"github.com/aaravmahajanofficial/ecommerce-project/database"
	"github.com/aaravmahajanofficial/ecommerce-project/middleware"
	"github.com/aaravmahajanofficial/ecommerce-project/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)

	router.Use(middleware.Authorization())
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/listcart", controllers.GetItemFromCart())
	router.POST("/addaddress", controllers.AddAddress())
	router.PUT("/edithomeaddress", controllers.EditHomeAddress())
	router.PUT("/editworkaddress", controllers.EditWorkAddress())
	router.GET("/deleteaddresses", controllers.DeleteAddress())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())
	log.Fatal(router.Run(":" + port))

}
