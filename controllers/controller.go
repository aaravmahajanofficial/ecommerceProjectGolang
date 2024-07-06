package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aaravmahajanofficial/ecommerce-project/database"
	"github.com/aaravmahajanofficial/ecommerce-project/models"
	generate "github.com/aaravmahajanofficial/ecommerce-project/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection = database.UserData(database.Client, "Users")
var ProductsCollection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

func HashPassword(password string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}

	return string(bytes)

}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(givenPassword))
	valid := true

	msg := ""

	if err != nil {
		msg = "Incorrect or Password is Incorrect"
		valid = false
	}

	return valid, msg

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

		// we want to create the new User model for insertion to the mongoDb collection

		password := HashPassword(*user.Password)
		user.Password = &password

		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshToken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshToken
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, insertError := UserCollection.InsertOne(context, user)

		if insertError != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"ERROR": "USER NOT CREATED"})
		}

		defer cancel()

		ctx.JSON(http.StatusCreated, "Successfully Signed Up!")

	}

}

func Login() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		// initialize the context

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var userDataFromDB models.User

		// bind the JSON

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"ERROR": err})
			return
		}

		err := UserCollection.FindOne(context, bson.M{"email": user.Email}).Decode(&userDataFromDB)
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"ERROR": "Incorrect email or password!"})
			return
		}

		PasswordIsValid, msg := VerifyPassword(*user.Password, *userDataFromDB.Password)
		defer cancel()

		if !PasswordIsValid {
			ctx.JSON(http.StatusInternalServerError, gin.H{"ERROR": msg})
			return
		}

		token, refreshToken, _ := generate.TokenGenerator(*userDataFromDB.Email, *userDataFromDB.First_Name, *userDataFromDB.Last_Name, userDataFromDB.User_ID)
		defer cancel()

		generate.UpdateAllTokens(token, refreshToken, userDataFromDB.User_ID)
		ctx.JSON(http.StatusFound, userDataFromDB)
	}

}

func ProductViewerAdmin() gin.HandlerFunc {

}

func SearchProduct() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		var productsList []models.Product

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductsCollection.Find(context, bson.D{{}})

		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "Someting Went Wrong Please Try After Some Time")
			return
		}

		err = cursor.All(ctx, &productsList)

		// If there is a problem while fetching all the documents at once (e.g., a network failure or invalid data that cannot be unmarshalled into the productsList slice), the error will be caught here.
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(context)

		// After attempting to read all documents, if there were any errors encountered during the cursor’s lifetime (not just during the All method call), they will be caught here.

		if err := cursor.Err(); err != nil {
			log.Println(err)
			ctx.IndentedJSON(400, "Invalid")
			return
		}

		ctx.IndentedJSON(200, productsList)

	}

}

func SearchProductByQuery() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		var searchResults []models.Product

		searchQuery := ctx.Query("name")

		if searchQuery == "" {

			log.Println("Query is Empty")
			ctx.Header("Content-Type", "application/json")
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"ERROR": "Invalid Search Index"})
			ctx.Abort()
			return

		}

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductsCollection.Find(context, bson.M{"$regex": searchQuery})

		if err != nil {
			log.Println(err)
			ctx.IndentedJSON(404, "Something went wrong while fetching the items.")
			return
		}

		err = cursor.All(context, &searchResults)

		if err != nil {
			log.Print(err)
			ctx.IndentedJSON(400, "Invalid Result")
			return
		}

		defer cursor.Close(context)

		if err := cursor.Err(); err != nil {
			log.Println(err)
			ctx.IndentedJSON(400, "Invalid Request")
			return
		}

		ctx.IndentedJSON(200, searchResults)

	}

}
