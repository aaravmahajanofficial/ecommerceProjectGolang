package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aaravmahajanofficial/ecommerce-project/database"
	"github.com/aaravmahajanofficial/ecommerce-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

var UserCollection = database.UserData(database.Client, "Users")
var ProductsCollection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

func HashPassword(password string) string {

}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {

}

func SignUp() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var user models.User

		// convert the incoming json in this type
		if err := ctx.BindJSON(&user); err != nil {

			// 400 Error
			ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return

		}

		// used for validate checks applied in the User Model

		validationError := Validate.Struct(user)

		if validationError != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"ERROR": validationError})
		}

		count, err := UserCollection.CountDocuments(context, bson.M{"email": user.Email})
		defer cancel()

		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"ERROR": err})
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"ERROR": "User already exists"})
		}

		count, err = UserCollection.CountDocuments(context, bson.M{"phone": user.Phone})
		defer cancel()

		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"ERROR": err})
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"ERROR": "Phone number already in use"})
			return
		}
	}

}

func Login() gin.HandlerFunc {

}

func ProductViewerAdmin() gin.HandlerFunc {

}

func SearchProduct() gin.HandlerFunc {

}

func SearchProductByQuery() gin.HandlerFunc {

}
