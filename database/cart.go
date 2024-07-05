package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("unable to find the specified product")
	ErrCantDecodeProducts = errors.New("unable to decode product information")
	ErrUserIDIsNotValid   = errors.New("user ID is not valid")
	ErrCantUpdateUser     = errors.New("unable to update user information")
	ErrCantRemoveItem     = errors.New("unable to remove item from cart")
	ErrCantGetItem        = errors.New("unable to retrieve item from cart")
	ErrCantBuyCartItem    = errors.New("unable to process the purchase of cart item")
)

func AddProductToCart() {

}

func RemoveCartItem() {

}

func BuyItemFromCart() {

}

func InstantBuy() {

}
