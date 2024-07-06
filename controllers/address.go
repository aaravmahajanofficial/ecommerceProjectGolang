package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/aaravmahajanofficial/ecommerce-project/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAddress() gin.HandlerFunc {

}

func EditHomeAddress() gin.HandlerFunc {

}

func EditWorkAddress() gin.HandlerFunc {

}

func DeleteAddress() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		// get the user id from the query

		userIDFromQuery := ctx.Query("id")

		if userIDFromQuery == "" {
			log.Println("User ID is empty")
			ctx.Header("Content-Type", "application/json")
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "UserID is empty."})
			ctx.Abort()
			return
		}

		userId, err := primitive.ObjectIDFromHex(userIDFromQuery)

		if err != nil {
			log.Println(err)
			ctx.AbortWithError(http.StatusInternalServerError, errors.New("Internal Server Error"))
			return
		}

		addresses := make([]models.Address, 0)

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// get this user from the document
		filter := bson.D{primitive.E{Key: "_id", Value: userId}}
		updatedValue := bson.D{primitive.E{Key: "$set", Value: primitive.E{Key: "address", Value: addresses}}}

		_, err = UserCollection.UpdateOne(context, filter, updatedValue)

		if err != nil {
			log.Println(err)
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"ERROR": "Something Went Wrong!"})
			return

		}

		ctx.IndentedJSON(http.StatusOK, gin.H{"Message": "Successfully Deleted!"})

	}

}
