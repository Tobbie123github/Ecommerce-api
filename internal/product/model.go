package product

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Product struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Price float64 `bson:"price" json:"price"`
	Discount float64 `bson:"discount" json:"discount"`
	Category string `bson:"category" json:"category"`
	ImageURL []string `bson:"imageUrl" json:"imageUrl"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

// type PublicProduct struct {
// 	ID        string    `json:"id"`
// 	Name      string    `json:"name"`
// 	Description     string    `json:"email"`
// 	Price       string    `json:"bio"`
// 	Discount      string    `json:"role"`
// 	Category     string    `json:"phone"`
// 	ImageURL string `json:"imageUrl"`
// 	CreatedAt time.Time `json:"createdAt"`
// 	UpdatedAt time.Time `json:"updatedAt"`
// }

// func ToPublicProduct(p Product) *PublicProduct {
// 	return &PublicProduct{
// 		ID:        p.ID.Hex(),
// 		Name:      p.Name,
// 		Description:     p.Description,
// 		Price:       p.Price,
// 		Discount:      p.Discount,
// 		Category:     p.Category,
// 		ImageURL: p.ImageURL,
// 		CreatedAt: p.CreatedAt.UTC(),
// 		UpdatedAt: p.CreatedAt.UTC(),
// 	}
// }