package database

import (
	"context"
	"errors"
	"log"
	"time"

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

	filter := bson.D{primitive.E{Key: "_id", Value: productID}}
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

func BuyItemFromCart(context context.Context, usersCollection *mongo.Collection, userID string) error {

	// convert userID to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}

	// initialize order model
	orderModel := models.Order{
		Order_ID:       primitive.NewObjectID(),
		Orderered_At:   time.Now(),
		Order_Cart:     make([]models.ProductUser, 0),
		Payment_Method: models.Payment{COD: true},
	}

	// find the documents of the user, and get the total of the items in the cart
	filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userObjectID}}}}
	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, primitive.E{Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "usercart.price"}}}}}}

	cursor, err := usersCollection.Aggregate(context, mongo.Pipeline{filter, unwind, group})

	if err != nil {
		log.Panic(err)
		return err
	}

	// extract total price from aggregation result
	var totalAmount int32
	var totalResults []bson.M

	if err = cursor.All(context, &totalResults); err != nil {
		log.Println(err)
		return err
	}

	for _, result := range totalResults {
		totalAmount = result["total"].(int32)
	}

	orderModel.Price = int(totalAmount)

	// push orderModel to user's orders array
	filter = bson.D{primitive.E{Key: "_id", Value: userObjectID}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderModel}}}}
	if _, err := usersCollection.UpdateOne(context, filter, update); err != nil {
		log.Println(err)
		return err
	}

	// retrieve user document to access usercart
	var userModel models.User
	if err := usersCollection.FindOne(context, filter).Decode(&userModel); err != nil {
		log.Println(err)
		return err
	}

	// add ordered items to each order's order_list array
	update = bson.D{{Key: "$push", Value: bson.D{{Key: "orders.$[].order_list", Value: bson.D{{Key: "$each", Value: userModel.UserCart}}}}}}
	if _, err := usersCollection.UpdateOne(context, filter, update); err != nil {
		log.Println(err)
		return err
	}

	// Clear user's cart after successful purchase
	emptyCart := make([]models.ProductUser, 0)
	emptyCartUpdate := bson.D{{Key: "$set", Value: bson.D{{Key: "usercart", Value: emptyCart}}}}
	if _, err := usersCollection.UpdateOne(context, filter, emptyCartUpdate); err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	return nil

}

func InstantBuy() {

}
