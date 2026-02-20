package product

import (
	"context"
	"errors"
	"fmt"
	"go-auth/internal/utils"
	"mime/multipart"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

type ProductInput struct {
	Name        string    `form:"name"`
	Description string    `form:"description"`
	Price       float64   `form:"price"`
	Discount    float64  `form:"discount"`
	Category    string    `form:"category"`
	Files     []*multipart.FileHeader `form:"files"`
}

// type AuthResult struct {
// 	Product   `json:"product"`
// }

func (s *Service) UploadItems(ctx context.Context, input ProductInput) (Product, error) {

	prodName := strings.TrimSpace(input.Name)
	prodDescription := strings.TrimSpace(input.Description)
	prodCategory := strings.TrimSpace(input.Category)

	if prodName == "" || prodDescription == "" || input.Price == 0 || prodCategory == ""{
		return Product{}, errors.New("Required Field")
	}

	// handle upload to cloudinary

	if input.Files == nil {
		return Product{}, errors.New("File should not be empty")
	}

	 cloud, err := utils.NewCloudinary()

	 if err != nil {
		return Product{}, fmt.Errorf("Issue with the cloud utils: %v", err)
	}

	var imageUrls []string

	for _,FileHeader := range input.Files {
		
		file, err := FileHeader.Open()

		if err != nil {
			return Product{}, fmt.Errorf("failed to get file header: %v", err)
		}

		defer file.Close()

		uploadRes , err := cloud.Upload.Upload(ctx, file, uploader.UploadParams{
				Folder:   "product-images",
				UniqueFilename: api.Bool(true),
		})

	if err != nil {
       return Product{}, fmt.Errorf("Error getting result: %v", err)
    }

	imageUrls = append(imageUrls, uploadRes.SecureURL)

	}
	
	
	now := time.Now().UTC()

	p := Product{
		Name: prodName,
		Description: prodDescription,
		Price: input.Price,
		Discount: input.Discount,
		Category: prodCategory,
		ImageURL: imageUrls,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Insert into db

	createdProduct, err := s.repo.Add(ctx, p)

	if err != nil {
		return Product{}, fmt.Errorf("Error from the Add repo : %v" , err)
	}

	return createdProduct, nil


}

