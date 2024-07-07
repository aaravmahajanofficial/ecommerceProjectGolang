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
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {

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
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}

		var newAddress models.Address
		newAddress.Address_id = primitive.NewObjectID()

		if err = ctx.BindJSON(&newAddress); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"error": "Invalid JSON data"})
			return
		}

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// aggregate the records, of the user

		filter := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: userID}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_value"}, primitive.E{Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		// eg :
		/*
			1. get the user record with id -> userId
			{
			  "$match": {
			    "_id": ObjectId("60e91db17f20d4d4f899fb56")
			  }
			}

			2. after unwinding the documents ->
				{
					"_id": ObjectId("60e91db17f20d4d4f899fb56"),
					"name": "John Doe",
					"address": { "address_id": ObjectId("60e91db17f20d4d4f899fb57"), "city": "New York" }
				},
				{
					"_id": ObjectId("60e91db17f20d4d4f899fb56"),
					"name": "John Doe",
					"address": { "address_id": ObjectId("60e91db17f20d4d4f899fb58"), "city": "San Francisco" }
				}

			3. Groups by address id

				[
					{ "_id": ObjectId("60e91db17f20d4d4f899fb57"), "count": 1 },
					{ "_id": ObjectId("60e91db17f20d4d4f899fb58"), "count": 1 }
				]

		*/

		cursor, err := UserCollection.Aggregate(context, mongo.Pipeline{filter, unwind, group})

		if err != nil {
			log.Println(err)
			ctx.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		defer cursor.Close(context)

		var currentAddresses []bson.M

		if err = cursor.All(context, &currentAddresses); err != nil {
			log.Println(err)
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error fetching adddress details"})
			return
		}

		var totalAddresses int32

		for _, value := range currentAddresses {
			totalAddresses += value["count"].(int32)
		}

		if totalAddresses < 2 {

			// add the address
			filter := bson.D{primitive.E{Key: "_id", Value: userID}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: newAddress}}}}
			_, err := UserCollection.UpdateOne(context, filter, update)

			if err != nil {
				ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update data"})
				return
			}

			ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Address added successfully"})

		} else {

			ctx.IndentedJSON(http.StatusBadRequest, "Not Allowed")

		}

	}

}

func EditHomeAddress() gin.HandlerFunc {

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
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}

		var newAddress models.Address

		if err = ctx.BindJSON(&newAddress); err != nil {

			ctx.JSON(http.StatusNotAcceptable, gin.H{"error": "Invalid JSON data"})
			return

		}

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userID}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: newAddress.House}, primitive.E{Key: "address.0.street_name", Value: newAddress.Street}, primitive.E{Key: "address.0.city_name", Value: newAddress.City}, primitive.E{Key: "address.0.pin_code", Value: newAddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(context, filter, update)

		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "Something Went Wrong")
			return
		}

		ctx.IndentedJSON(http.StatusOK, "Successfully updated home address")

	}

}

func EditWorkAddress() gin.HandlerFunc {

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
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}

		var newAddress models.Address

		if err = ctx.BindJSON(&newAddress); err != nil {

			ctx.JSON(http.StatusNotAcceptable, gin.H{"error": "Invalid JSON data"})
			return

		}

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userID}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: newAddress.House}, primitive.E{Key: "address.1.street_name", Value: newAddress.Street}, primitive.E{Key: "address.1.city_name", Value: newAddress.City}, primitive.E{Key: "address.1.pin_code", Value: newAddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(context, filter, update)

		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "Something Went Wrong")
			return
		}

		ctx.IndentedJSON(http.StatusOK, "Successfully updated home address")

	}

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
