package order

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type DeliveryInfo struct {
    FullName   string `bson:"fullName" json:"fullName"`
    Phone      string `bson:"phone" json:"phone"`
    Address    string `bson:"address" json:"address"`
    City       string `bson:"city" json:"city"`
    State      string `bson:"state" json:"state"`
    Country    string `bson:"country" json:"country"`
}

type OrderItem struct {
    ProductID   string  `bson:"productId" json:"productId"`
    Name        string  `bson:"name" json:"name"`
    Price       float64 `bson:"price" json:"price"`
    Quantity    int     `bson:"quantity" json:"quantity"`
    Subtotal    float64 `bson:"subtotal" json:"subtotal"` 
}

type Order struct {
    ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID      string        `bson:"userId" json:"userId"`
    Items       []OrderItem   `bson:"items" json:"items"`
    Total       float64       `bson:"total" json:"total"`
    Status      string        `bson:"status" json:"status"`
     Delivery    DeliveryInfo  `bson:"delivery" json:"delivery"` 
    PaymentRef  string        `bson:"paymentRef" json:"-"` 
    CreatedAt   time.Time     `bson:"createdAt" json:"createdAt"`
    UpdatedAt   time.Time     `bson:"updatedAt" json:"updatedAt"`
}