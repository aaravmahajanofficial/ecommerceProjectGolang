package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/aaravmahajanofficial/ecommerce-project/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	productsCollection *mongo.Collection
	usersCollection    *mongo.Collection
}

func NewApplication(productCollection *mongo.Collection, userCollection *mongo.Collection) *Application {

	return &Application{
		productsCollection: productCollection,
		usersCollection:    userCollection,
	}
}

// use "AbortWithError", when dealing with critical functions, like validation errors, database queries, authorization errors and "JSON" or "IndentedJSON" only when simple logging like success code etc.

func (app *Application) AddToCart() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		// get the product id from the query parameter
		productQueryID := ctx.Query("id")

		if productQueryID == "" {
			log.Println("Product ID is empty")
			ctx.AbortWithError(http.StatusBadRequest, errors.New("ProductID is empty."))
			return
		}

		userQueryID := ctx.Query("userID")

		if userQueryID == "" {
			log.Println("User ID is empty")
			ctx.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty."))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return

		}

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = database.AddProductToCart(context, app.productsCollection, app.usersCollection, productID, userQueryID)

		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
		}

		ctx.IndentedJSON(200, "Successfully Added to Cart")

	}

}

func (app *Application) RemoveItem() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		productQueryID := ctx.Query("id")
		if productQueryID == "" {
			log.Println("Product ID is empty")
			ctx.AbortWithError(http.StatusBadRequest, errors.New("ProductID is empty."))
			return
		}

		userQueryID := ctx.Query("userId")
		if userQueryID == "" {
			log.Println("User ID is empty")
			ctx.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty."))
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return

		}

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.RemoveCartItem(context, app.productsCollection, app.usersCollection, productId, userQueryID)

		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
		}

		ctx.IndentedJSON(200, "Successfully Removed from Cart")

	}

}

func GetItemFromCart() gin.HandlerFunc {

}

func (app *Application) BuyFromCart() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		userQueryId := ctx.Query("userId")

		if userQueryId == "" {
			log.Print("User ID is empty")

			ctx.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty."))
			return
		}

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(context, app.usersCollection, userQueryId)

		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
		}

		ctx.IndentedJSON(200, "Successfully Placed the Order")
	}

}

func (app *Application) InstantBuy() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		productQueryID := ctx.Query("id")
		if productQueryID == "" {
			log.Println("Product ID is empty")
			ctx.AbortWithError(http.StatusBadRequest, errors.New("ProductID is empty."))
			return
		}

		userQueryID := ctx.Query("userId")
		if userQueryID == "" {
			log.Println("User ID is empty")
			ctx.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty."))
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return

		}

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.InstantBuy(context, app.productsCollection, app.usersCollection, productId, userQueryID)

		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
		}

		ctx.IndentedJSON(200, "Successfully Placed the Order")
	}

}
