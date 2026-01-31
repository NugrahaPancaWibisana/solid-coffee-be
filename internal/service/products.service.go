package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/repository"
	"github.com/redis/go-redis/v9"
)

type ProductService struct {
	productRepository *repository.ProductRepository
	redis             *redis.Client
}

func NewProductService(productRepository *repository.ProductRepository, rdb *redis.Client) *ProductService {
	return &ProductService{
		productRepository: productRepository,
		redis:             rdb,
	}
}

func (p ProductService) GetAllProducts(ctx context.Context, page int) ([]dto.Products, error) {
	rkey := "solid:products"
	rsc := p.redis.Get(ctx, rkey)
	if rsc.Err() == nil {
		var result []dto.Products
		cache, err := rsc.Bytes()
		if err != nil {
			log.Println(err)
		} else {
			if err := json.Unmarshal(cache, &result); err != nil {
				log.Println(err.Error())
			} else {
				return result, nil
			}
		}
	}

	if rsc.Err() == redis.Nil {
		log.Println("movies cache miss")
	}

	var response []dto.Products

	data, err := p.productRepository.GetAllProduct(ctx, page)
	if err != nil {
		return []dto.Products{}, err
	}
	for _, v := range data {
		response = append(response, dto.Products{
			Id:          v.Id,
			Name:        v.Name,
			Images_Name: v.Images_Name,
			Price:       v.Price,
			Discount:    v.Discount,
			Rating:      v.Rating,
		})
	}

	cacheStr, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		log.Println("failed to marshal")
	}

	rdsStatus := p.redis.Set(ctx, rkey, string(cacheStr), time.Minute*10)
	if rdsStatus.Err() != nil {
		log.Println("caching failed")
		log.Println(rdsStatus.Err().Error())
	}

	return response, nil
}

func (p ProductService) GetTotalPage(ctx context.Context) (int, error) {
	data, err := p.productRepository.GetTotalPage(ctx)
	if err != nil {
		return 0, err
	}

	return data, nil
}
