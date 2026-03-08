package payment

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Payment struct {
	ID                bson.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID           string        `bson:"orderId" json:"orderId"`
	UserID            string        `bson:"userId" json:"userId"`
	Amount            float64       `bson:"amount" json:"amount"`
	Currency          string        `bson:"currency" json:"currency"`
	Status            string        `bson:"status" json:"status"`
	PaystackReference string        `bson:"paystackReference" json:"-"`
	AuthorizationURL  string        `bson:"authorizationUrl" json:"-"`
	CreatedAt         time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time     `bson:"updatedAt" json:"updatedAt"`
}
