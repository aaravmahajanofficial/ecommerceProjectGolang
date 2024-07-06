package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/aaravmahajanofficial/ecommerce-project/database"
	"github.com/aaravmahajanofficial/ecommerce-project/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

	return func(ctx *gin.Context) {

		userIDFromQuery := ctx.Query("id")

		if userIDFromQuery == "" {
			log.Println("User ID is empty")
			ctx.Header("Content-Type", "application/json")
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "UserID is empty."})
			ctx.Abort()
			return
		}

		userID, err := primitive.ObjectIDFromHex(userIDFromQuery)

		if err != nil {
			log.Println(err)
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": "Internal Server Error"})
			return
		}

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var userModel models.User

		err = UserCollection.FindOne(context, bson.D{primitive.E{Key: "_id", Value: userID}}).Decode(&userModel)

		// get the user documents

		filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userID}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, primitive.E{Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

		cursor, err := UserCollection.Aggregate(context, mongo.Pipeline{filter, unwind, group})

		if err != nil {
			log.Println("Error in aggregation:", err)
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error calculating cart total"})
			return
		}

		defer cursor.Close(context)

		var results []bson.M

		if err = cursor.All(context, &results); err != nil {

			log.Println("Error fetching aggregation results:", err)
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error fetching cart details"})
			return

		}

		if len(results) == 0 {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "no items in cart"})
			return
		}

		// aggregation framework always returns an array of documents

		// eg :
		/*
			[
				{
					"_id": ObjectId("60d5f483f8658e6b9c8e4d65"),
					"total": 3.7
				}
			]
		*/

		total := results[0]["total"]

		response := gin.H{
			"total":    total,
			"userCart": userModel.UserCart,
		}

		ctx.IndentedJSON(http.StatusOK, response)

	}

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
