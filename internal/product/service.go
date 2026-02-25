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
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
	InStock int `form:"instock"`
	Files     []*multipart.FileHeader `form:"files"`
}

type UpdateProduct struct {
	Name        string    `form:"name"`
	Description string    `form:"description"`
	Price       float64   `form:"price"`
	Discount    float64  `form:"discount"`
	InStock int `form:"instock"` 
	Category    string    `form:"category"`
	Files     []*multipart.FileHeader `form:"files"`
	ImageURL    []string                
    PublicIDs   []string
}


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
	var publicIds []string

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
	publicIds = append(publicIds, uploadRes.PublicID)

	}
	
	
	now := time.Now().UTC()

	p := Product{
		Name: prodName,
		Description: prodDescription,
		Price: input.Price,
		Discount: input.Discount,
		Category: prodCategory,
		ImageURL: imageUrls,
		PublicIDs: publicIds,
		InStock: input.InStock,
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


func (s *Service) UpdateItems(ctx context.Context, input UpdateProduct, id string) (Product, error) {

	prodName := strings.TrimSpace(input.Name)
	prodDescription := strings.TrimSpace(input.Description)
	prodCategory := strings.TrimSpace(input.Category)

	if prodName == "" || prodDescription == "" || input.Price == 0 || prodCategory == ""{
		return Product{}, errors.New("Required Field")
	}

	 objectId, err := bson.ObjectIDFromHex(id)

	// get existing public ids

	existingProd, err := s.repo.GetById(ctx, objectId)

	if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return Product{}, errors.New("product not found")
        }
        return Product{}, fmt.Errorf("unable to get existing product: %v", err)
    }

	cloud, err := utils.NewCloudinary()
	if err != nil {
		return Product{}, fmt.Errorf("Unable to get cloud")
	}

	if input.Files != nil {
		for _, pubID := range existingProd.PublicIDs{
			_, err := cloud.Upload.Destroy(ctx, uploader.DestroyParams{
					PublicID: pubID,
			})

			if err != nil {
				return Product{}, fmt.Errorf("failed to delete old image: %v", err)
			}
		}


		// upload new images

		var imageUrls []string
		var publicIds []string

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
	publicIds = append(publicIds, uploadRes.PublicID)

	}

	input.ImageURL = imageUrls
	input.PublicIDs = publicIds


	}else{

		input.ImageURL = existingProd.ImageURL
		input.PublicIDs = existingProd.PublicIDs

	}

	prod, err := s.repo.UpdateById(ctx, objectId, input)

	if err != nil {
		return Product{}, fmt.Errorf("Error with the update: %v", err)
	}

	return prod, nil

}


func (s *Service) DeleteProduct(ctx context.Context, id string) error {

	idHex, err := bson.ObjectIDFromHex(id)

	if err != nil {
		return fmt.Errorf("Error converting ID: %v", err)
	}

	existingProd, err := s.repo.GetById(ctx, idHex)

	if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return  errors.New("product not found")
        }
        return fmt.Errorf("unable to get existing product: %v", err)
    }

	// delete images from cloud

	cloud, err := utils.NewCloudinary()

	if err != nil {
		return fmt.Errorf("Error occured inj the cloudinary %v", err)
	}

	for _, pubID := range existingProd.PublicIDs{
		_, err := cloud.Upload.Destroy(ctx, uploader.DestroyParams{
					PublicID: pubID,
			})

			if err != nil {
				return fmt.Errorf("failed to delete old image: %v", err)
			}

	}


	if err := s.repo.DeleteItem(ctx, idHex); err !=nil  {
		return fmt.Errorf("Error occured in the delete item repo: %v", err)
	}

	return nil


}

