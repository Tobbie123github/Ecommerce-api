package user

import (
	"context"
	"errors"
	"fmt"
	"go-auth/internal/auth"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *Repo
	jwtSecret string
}

func NewService(repo *Repo, jwtSecret string) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

type RegisterInput struct {
	Name string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResult struct {
	Token string     `json:"token"`
	User  PublicUser `json:"user"`
}

// jwt, authorization

func (s *Service) Register(ctx context.Context, input RegisterInput) (AuthResult, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := strings.ToLower(strings.TrimSpace(input.Password))
	name := strings.ToLower(strings.TrimSpace(input.Name))

	if email == "" || password == "" || name == ""{
		return AuthResult{}, errors.New("Name, Email and password must not be empty")
	}

	if len(password) < 6 {
		return AuthResult{}, errors.New("Password must be of lenght 6")
	}

	// check if the user already exist from repo.go

	_, err := s.repo.FindByEmail(ctx, email)

	if err == nil {
		return AuthResult{}, fmt.Errorf("Email already exist: %v", err)
	}

	// 	if !errors.Is(err, mongo.ErrNoDocuments) {
	//     return AuthResult{}, fmt.Errorf("database error: %v", err)
	//    }

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return AuthResult{}, err
	}

	// HashPassword

	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return AuthResult{}, fmt.Errorf("Unable to Hash paswword:%v", err)

	}

	now := time.Now().UTC()

	u := User{
		Name: name,
		Email:        email,
		PasswordHash: string(hashPass),
		Role:         "user",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	createdUser, err := s.repo.Create(ctx, u)

	if err != nil {
		return AuthResult{}, fmt.Errorf("Unable to Create User:%v", err)

	}

	// create token from jwt.go
	token, err := auth.CreateToken(s.jwtSecret, createdUser.ID.Hex(), createdUser.Role)

	if err != nil {
		return AuthResult{}, fmt.Errorf("Err to Create Token: %v", err)
	}

	return AuthResult{
		Token: token,
		User:  ToPublic(createdUser),
	}, nil

}

func checkPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("Error macthing pw %v", err)
	}

	return nil
}

// Login

func (s *Service) Login(ctx context.Context, input LoginInput) (AuthResult, error) {

	email := strings.ToLower(strings.TrimSpace(input.Email))
	pass := strings.ToLower(strings.TrimSpace(input.Password))

	if email == "" || pass == "" {
		return AuthResult{}, errors.New("Email and password must not be empty")
	}

	// check if user is in db
	u, err := s.repo.FindByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return AuthResult{}, mongo.ErrNoDocuments
		}

		return AuthResult{}, fmt.Errorf("Find by email failed: %s", err)
	}

	// if user is there, then compare the user password with the hashed in the db

	if err := checkPassword(u.PasswordHash, pass); err != nil {
		return AuthResult{}, fmt.Errorf("Incorrect Password %v", err)
	}

	token, err := auth.CreateToken(s.jwtSecret, u.ID.Hex(), u.Role)

	if err != nil {
		return AuthResult{}, fmt.Errorf("Err to Create Token: %v", err)
	}

	return AuthResult{
		Token: token,
		User:  ToPublic(u),
	}, nil

}
