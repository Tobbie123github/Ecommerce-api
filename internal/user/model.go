package user

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {

	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"passwordHash" json:"-"`
	Bio          string             `bson:"bio" json:"bio"`
	Role         string             `bson:"role" json:"role"`
	Phone        string             `bson:"phone" json:"phone"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type PublicUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Bio       string    `json:"bio"`
	Role      string    `json:"role"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToPublic(u User) PublicUser {
	return PublicUser{
		ID:        u.ID.Hex(),
		Name:      u.Name,
		Email:     u.Email,
		Bio:       u.Bio,
		Role:      u.Role,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt.UTC(),
		UpdatedAt: u.CreatedAt.UTC(),
	}
}
