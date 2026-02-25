package cart

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-auth/internal/product"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Repo struct {
	rdb         *redis.Client
	productRepo *product.Repo
}

func NewRepo(rdb *redis.Client, productRepo *product.Repo) *Repo {
	return &Repo{
		rdb:         rdb,
		productRepo: productRepo,
	}
}

//helper key func

func helperKey(userId string) string {
	return "cart:" + userId
}

func (pr Repo) AddCart(ctx context.Context, productId string, userId string, quantity int) (Cart, error) {

	if productId == "" || userId == "" || quantity == 0 {
		return Cart{}, errors.New("Invalid body request")
	}



	id, err := bson.ObjectIDFromHex(productId)

	if err != nil {
		return Cart{}, fmt.Errorf("Error getting id: %v", err)
	}

	product, err := pr.productRepo.GetById(ctx, id)
	if err != nil {
		return Cart{}, fmt.Errorf("unable to get product id: %v", err)
	}

	if quantity > product.InStock {
		return Cart{}, fmt.Errorf("Your quantity exceed peoduct in stock %v, %v", product.InStock, err)
	}
	

	cartItem := CartItem{
		ProductID: product.ID.Hex(),
		Name:      product.Name,
		Price:     product.Price,
		ImageURL:  product.ImageURL,
		Quantity:  quantity,
	}

	//convert to struct

	cartData, err := json.Marshal(cartItem)

	if err != nil {
		return Cart{}, fmt.Errorf("Error converting to string: %v", err)
	}

	// store redis using userId

	key := helperKey(userId)

	// redisClient, err := db.Redis(ctx)

	err = pr.rdb.HSet(ctx, key, productId, cartData).Err()

	if err != nil {
		return Cart{}, fmt.Errorf("error storing cart in redis: %v", err)
	}

	return Cart{
		UserId: userId,
		Items:  []CartItem{cartItem},
		Total:  cartItem.Price * float64(cartItem.Quantity),
	}, nil

}

func (pr Repo) GetCart(ctx context.Context, userId string) (Cart, error) {

	key := helperKey(userId)
	res, err := pr.rdb.HGetAll(ctx, key).Result()

	if err == redis.Nil {
		return Cart{}, fmt.Errorf("Cart Not found: %v", err)
	} else if err != nil {
		return Cart{}, fmt.Errorf("Failed to fetch cart")

	}

	if len(res) == 0 {
		return Cart{
			UserId: userId,
			Items:  []CartItem{},
			Total:  0,
		}, nil
	}

	var items []CartItem
	var total float64

	for _, val := range res {
		var item CartItem
		if err := json.Unmarshal([]byte(val), &item); err != nil {
			return Cart{}, fmt.Errorf("Failed to parse cart data: %v", err)
		}

		items = append(items, item)
		total += item.Price * float64(item.Quantity)

	}

	return Cart{
		UserId: userId,
		Items:  items,
		Total:  total,
	}, nil

}


func (pr Repo) DeleteCartItem(ctx context.Context, productId string, userId string) error {

	key := helperKey(userId)

	_, err := pr.rdb.HDel(ctx, key, productId).Result()

	if err != nil {
		return errors.New("Error deleting product from cart")
	}

	return nil
}

func (pr Repo) DeleteAll(ctx context.Context, userId string) error {
	key := helperKey(userId)

	_, err := pr.rdb.Del(ctx, key).Result() 

	if err != nil {
		return errors.New("Error deleting all product from cart")
	}

	return nil
}
