package order

import (
	"context"
	"errors"
	"fmt"
	"go-auth/internal/cart"
	"go-auth/internal/product"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Service struct {
	repo        *Repo
	cartRepo    *cart.Repo
	productRepo *product.Repo
}

func NewService(repo *Repo, cartRepo *cart.Repo, productRepo *product.Repo) *Service {
	return &Service{
		repo:        repo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (svc *Service) CreateOrder(ctx context.Context, userId string) (Order, error) {

	// Get cart from redis

	cartItem, err := svc.cartRepo.GetCart(ctx, userId)

	if len(cartItem.Items) == 0 {
		return Order{}, errors.New("cart is empty")
	}

	if err != nil {
		return Order{}, fmt.Errorf("Error getting cart items: %v", err)

	}

	var orderItems []OrderItem
	var total float64

	for _, item := range cartItem.Items {

		// validate each product still exists in database

		productId, err := bson.ObjectIDFromHex(item.ProductID)

		if err != nil {
			return Order{}, fmt.Errorf("Unable to convert id: %v", err)
		}

		product, err := svc.productRepo.GetById(ctx, productId)

		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return Order{}, errors.New("product not found")
			}
			return Order{}, fmt.Errorf("unable to get existing product: %v", err)
		}

		// if product is not instock
		if product.InStock == 0 {
			return Order{}, fmt.Errorf("Product not instock: %v", err)
		}

		// if instock then create order

		// insert the cart items into order items
		subtotal := item.Price * float64(item.Quantity)

		orderItems = append(orderItems, OrderItem{
			ProductID: item.ProductID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			Subtotal:  subtotal,
		})

		total = total + subtotal

	}

	o := Order{
		UserID:     userId,
		Items:      orderItems,
		Total:      total,
		Status:     "pending",
		PaymentRef: "nil",
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	order, err := svc.repo.Create(ctx, o)

	if err != nil {
		return Order{}, fmt.Errorf("Error creating order: %v", err)
	}

	// clear cart after order is created

	if err := svc.cartRepo.DeleteAll(ctx, userId); err != nil {
		return Order{}, fmt.Errorf("Error clearing cart")
	}

	return order, nil

}

func (svc *Service) DeleteOrder(ctx context.Context, orderId string, userId string) error {

	// verify the order belongs to this user before deleting
	order, err := svc.repo.GetById(ctx, orderId)

	if err != nil {
		return fmt.Errorf("order not found: %v", err)
	}

	//  if order.UserID != userId {
	//     return errors.New("unauthorized to delete this order")
	// }

	// verify the order belongs to this user before deleting
	if order.UserID != userId {
		return errors.New("Not Authorized")
	}

	id, err := bson.ObjectIDFromHex(orderId)
	if err != nil {
		return fmt.Errorf("invalid order id: %v", err)
	}

	return svc.repo.Delete(ctx, id)
}


func (svc *Service) UpdateOrder(ctx context.Context, orderId string, userId string, status string, paymentRef string) error {


	order, err := svc.repo.GetById(ctx, orderId)

   if err != nil {
		return fmt.Errorf("order not found: %v", err)
	}

    if order.UserID != userId {
		return errors.New("Not Authorized")
	}

	if err := svc.repo.UpdateStatus(ctx, orderId, status, paymentRef); err != nil {
		return fmt.Errorf("Error updating status: %v", err)
	}

	return nil

}

