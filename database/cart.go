package database

import (
	"context"
	"errors"
	"log"

	"github.com/aaravmahajanofficial/ecommerce-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("unable to find the specified product")
	ErrCantDecodeProducts = errors.New("unable to decode product information")
	ErrUserIDIsNotValid   = errors.New("user ID is not valid")
	ErrCantUpdateUser     = errors.New("unable to update user information")
	ErrCantRemoveItem     = errors.New("unable to remove item from cart")
	ErrCantGetItem        = errors.New("unable to retrieve item from cart")
	ErrCantBuyCartItem    = errors.New("unable to process the purchase of cart item")
)

func AddProductToCart(context context.Context, productsCollection *mongo.Collection, usersCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {

	filter := bson.D{primitive.E{Key: "_id", Value: userID}}
	cursor, err := productsCollection.Find(context, filter)

	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var itemsToBeAdded []models.ProductUser

	if err = cursor.All(context, &itemsToBeAdded); err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	objectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}

	// need to find the document of the user, to insert this item in the user cart
	filter = bson.D{primitive.E{Key: "_id", Value: objectID}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{primitive.E{Key: "$each", Value: itemsToBeAdded}}}}}}
	_, err = usersCollection.UpdateOne(context, filter, update)

	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}

	return nil

}

func RemoveCartItem(context context.Context, productsCollection *mongo.Collection, usersCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {

	// first convert the userID to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: objectID}}
	update := bson.D{{Key: "$pull", Value: bson.D{primitive.E{Key: "usercart", Value: productID}}}}
	_, err = usersCollection.UpdateMany(context, filter, update)

	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}

	return nil

}

func BuyItemFromCart() {

}

func InstantBuy() {

}
