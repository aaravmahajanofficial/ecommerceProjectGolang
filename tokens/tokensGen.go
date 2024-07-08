package tokens

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aaravmahajanofficial/ecommerce-project/database"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var SECRET_KEY = os.Getenv("SECRET_KEY")

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	UID        string
	jwt.StandardClaims
}

func TokenGenerator(email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshToken string, err error) {

	claims := SignedDetails{

		Email:      email,
		First_Name: firstName,
		Last_Name:  lastName,
		UID:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err

}

func VerifyToken(signedToken string) (claims *SignedDetails, msg string) {

	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {

		return []byte(SECRET_KEY), nil

	})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		msg = "The Token is not valid."
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Toke is expired"
		return
	}

	return claims, msg

}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userID string) {

	context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{primitive.E{Key: "_id", Value: userID}}
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "token", Value: signedToken}, primitive.E{Key: "refresh_token", Value: signedRefreshToken}, primitive.E{Key: "updatedat", Value: updated_at}}}}

	_, err := UserCollection.UpdateOne(context, filter, update, options.Update().SetUpsert(true))

	if err != nil {
		log.Panic(err)
		return
	}

}
