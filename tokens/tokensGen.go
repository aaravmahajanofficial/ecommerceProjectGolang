package tokens

import (
	"log"
	"os"
	"time"

	"github.com/aaravmahajanofficial/ecommerce-project/database"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo"
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
