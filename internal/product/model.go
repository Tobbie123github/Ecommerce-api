package product

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Product struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
	Price       float64       `bson:"price" json:"price"`
	Discount    float64       `bson:"discount" json:"discount"`
	Size        string        `bson:"size" json:"size"`
	Featured    bool          `bson:"featured" json:"featured"`
	InStock     int           `bson:"instock" json:"instock"`
	Category    string        `bson:"category" json:"category"`
	ImageURL    []string      `bson:"imageUrl" json:"imageUrl"`
	PublicIDs   []string      `bson:"publicIds" json:"-"`
	CreatedAt   time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt" json:"updatedAt"`
}
