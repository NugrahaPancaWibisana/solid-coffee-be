package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ProductService struct {
	productRepository *repository.ProductRepository
	redis             *redis.Client
	db                *pgxpool.Pool
}

func NewProductService(productRepository *repository.ProductRepository, db *pgxpool.Pool, rdb *redis.Client) *ProductService {
	return &ProductService{
		productRepository: productRepository,
		redis:             rdb,
		db:                db,
	}
}

func (ps *ProductService) invalidateProductsCache(ctx context.Context) error {
	pattern := fmt.Sprintf("%s:products:*", os.Getenv("RDB_KEY"))
	
	iter := ps.redis.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		err := ps.redis.Del(ctx, iter.Val()).Err()
		if err != nil {
			log.Printf("failed to delete cache key %s: %v", iter.Val(), err)
		}
	}
	
	if err := iter.Err(); err != nil {
		log.Printf("error during cache invalidation: %v", err)
		return err
	}
	
	return nil
}

func (ps *ProductService) invalidateCache(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	
	err := ps.redis.Del(ctx, keys...).Err()
	if err != nil {
		log.Printf("failed to invalidate cache: %v", err)
		return err
	}
	
	return nil
}

func (ps ProductService) GetAllProducts(ctx context.Context, req dto.ProductQueries) ([]dto.Products, int, error) {
	rkey := fmt.Sprintf("%s:products:page=%s:title=%s:min=%s:max=%s:categories=%s",
		os.Getenv("RDB_KEY"), req.Page, req.Title, req.Min, req.Max, strings.Join(req.Category, ","))

	rsc := ps.redis.Get(ctx, rkey)
	if rsc.Err() == nil {
		var result struct {
			Products  []dto.Products `json:"products"`
			TotalPage int            `json:"total_page"`
		}
		cache, err := rsc.Bytes()
		if err != nil {
			log.Println(err)
		} else {
			if err := json.Unmarshal(cache, &result); err != nil {
				log.Println(err.Error())
			} else {
				return result.Products, result.TotalPage, nil
			}
		}
	}

	if rsc.Err() == redis.Nil {
		log.Println("products cache miss")
	}

	totalPage, err := ps.productRepository.GetTotalPage(ctx, ps.db, req)
	if err != nil {
		return []dto.Products{}, 0, err
	}

	data, err := ps.productRepository.GetProducts(ctx, ps.db, req)
	if err != nil {
		return []dto.Products{}, 0, err
	}

	var response []dto.Products
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

	cacheData := struct {
		Products  []dto.Products `json:"products"`
		TotalPage int            `json:"total_page"`
	}{
		Products:  response,
		TotalPage: totalPage,
	}

	cacheStr, err := json.Marshal(cacheData)
	if err != nil {
		log.Println(err)
		log.Println("failed to marshal")
	}

	rdsStatus := ps.redis.Set(ctx, rkey, string(cacheStr), time.Minute*10)
	if rdsStatus.Err() != nil {
		log.Println("caching failed")
		log.Println(rdsStatus.Err().Error())
	}

	return response, totalPage, nil
}

func (ps ProductService) PostProduct(ctx context.Context, post dto.PostProductsRequest, images dto.PostImagesRequest) (dto.PostProductResponse, error) {
	tx, err := ps.db.Begin(ctx)
	if err != nil {
		log.Println(err)
		return dto.PostProductResponse{}, err
	}

	data, err := ps.productRepository.PostProduct(ctx, tx, post)
	if err != nil {
		return dto.PostProductResponse{}, err
	}
	defer tx.Rollback(ctx)

	for i := range len(images.Images_Name) {
		_, err := ps.productRepository.PostImages(ctx, tx, data.Id, images.Images_Name[i])
		if err != nil {
			return dto.PostProductResponse{}, err
		}
	}

	if e := tx.Commit(ctx); e != nil {
		log.Println("failed to commit", e.Error())
		return dto.PostProductResponse{}, e
	}

	ps.invalidateProductsCache(ctx)
	ps.invalidateCache(ctx, 
		fmt.Sprintf("%s:product_admin", os.Getenv("RDB_KEY")),
	)

	response := dto.PostProductResponse{
		Id: data.Id,
	}

	return response, nil
}

func (ps ProductService) UpdateProduct(ctx context.Context, update dto.UpdateProductsRequest, images dto.PostImagesRequest, idProduct int) error {
	tx, err := ps.db.Begin(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	if update.Description != "" || update.Price != 0 || update.ProductName != "" {
		cmd, err := ps.productRepository.UpdateProduct(ctx, tx, update, idProduct)
		if err != nil {
			return err
		}
		if cmd.RowsAffected() == 0 {
			return errors.New("no data updated")
		}
	}

	defer tx.Rollback(ctx)

	for i := range len(images.Images_Name) {
		cmd, err := ps.productRepository.PostImages(ctx, tx, idProduct, images.Images_Name[i])
		if err != nil {
			return err
		}
		if cmd.RowsAffected() == 0 {
			return errors.New("no data inserted")
		}
	}

	if e := tx.Commit(ctx); e != nil {
		log.Println("failed to commit", e.Error())
		return e
	}

	ps.invalidateProductsCache(ctx)
	ps.invalidateCache(ctx, 
		fmt.Sprintf("%s:product_admin", os.Getenv("RDB_KEY")),
	)

	return nil
}

func (ps ProductService) DeleteProductById(ctx context.Context, idProduct int) error {
	tx, err := ps.db.Begin(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	cmd, err := ps.productRepository.DeleteProductById(ctx, tx, idProduct)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("no data deleted")
	}

	defer tx.Rollback(ctx)

	cmdDel, errDel := ps.productRepository.DeleteProductImage(ctx, tx, idProduct)
	if errDel != nil {
		return err
	}
	if cmdDel.RowsAffected() == 0 {
		return errors.New("no data deleted")
	}

	if e := tx.Commit(ctx); e != nil {
		log.Println("failed to commit", e.Error())
		return e
	}

	ps.invalidateProductsCache(ctx)
	ps.invalidateCache(ctx, 
		fmt.Sprintf("%s:product_admin", os.Getenv("RDB_KEY")),
	)

	return nil
}

func (ps ProductService) DeleteProductImageById(ctx context.Context, idImages int) error {
	cmd, err := ps.productRepository.DeleteProductImageById(ctx, ps.db, idImages)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("no data deleted")
	}
	
	ps.invalidateProductsCache(ctx)
	ps.invalidateCache(ctx, 
		fmt.Sprintf("%s:product_admin", os.Getenv("RDB_KEY")),
	)
	
	return nil
}

func (ps ProductService) GetProductById(ctx context.Context, idProduct int) (dto.DetailProduct, error) {
	var response dto.DetailProduct

	data, err := ps.productRepository.GetProductById(ctx, ps.db, idProduct)
	if err != nil {
		return dto.DetailProduct{}, err
	}

	response = dto.DetailProduct{
		IdProduct:   data.IdProduct,
		ProductName: data.ProductName,
		Description: data.Description,
		Price:       data.Price,
		IdImages:    data.IdImages,
		Images:      data.Images,
	}

	return response, nil
}

func (ps ProductService) GetDetailProductByUserWithId(ctx context.Context, idMenu int) (dto.DetailProductUser, error) {
	var response dto.DetailProductUser

	data, err := ps.productRepository.GetDetailProductByUserWithId(ctx, ps.db, idMenu)
	if err != nil {
		return dto.DetailProductUser{}, err
	}

	response = dto.DetailProductUser{
		IdProduct:    data.IdProduct,
		ProductName:  data.ProductName,
		Images:       data.Images,
		Price:        data.Price,
		Description:  data.Description,
		Discount:     data.Discount,
		Rating:       data.Rating,
		Total_Review: data.Total_Review,
	}
	return response, nil
}

func (ps *ProductService) GetAllProductType(ctx context.Context) ([]dto.ProductType, error) {
	rkey := fmt.Sprintf("%s:product_type", os.Getenv("RDB_KEY"))

	rsc := ps.redis.Get(ctx, rkey)
	if rsc.Err() == nil {
		var result []dto.ProductType
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
		log.Println("product_type cache miss")
	}

	data, err := ps.productRepository.GetAllProductType(ctx, ps.db)
	if err != nil {
		return []dto.ProductType{}, err
	}

	var response []dto.ProductType
	for _, v := range data {
		response = append(response, dto.ProductType{
			Id:    v.Id,
			Name:  v.Name,
			Price: v.Price,
		})
	}

	cacheStr, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		log.Println("failed to marshal")
	}

	rdsStatus := ps.redis.Set(ctx, rkey, string(cacheStr), time.Minute*10)
	if rdsStatus.Err() != nil {
		log.Println("caching failed")
		log.Println(rdsStatus.Err().Error())
	}

	return response, nil
}

func (ps *ProductService) GetAllProductSize(ctx context.Context) ([]dto.ProductSize, error) {
	rkey := fmt.Sprintf("%s:product_size", os.Getenv("RDB_KEY"))

	rsc := ps.redis.Get(ctx, rkey)
	if rsc.Err() == nil {
		var result []dto.ProductSize
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
		log.Println("product_size cache miss")
	}

	data, err := ps.productRepository.GetAllProductSize(ctx, ps.db)
	if err != nil {
		return []dto.ProductSize{}, err
	}

	var response []dto.ProductSize
	for _, v := range data {
		response = append(response, dto.ProductSize{
			Id:    v.Id,
			Name:  v.Name,
			Price: v.Price,
		})
	}

	cacheStr, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		log.Println("failed to marshal")
	}

	rdsStatus := ps.redis.Set(ctx, rkey, string(cacheStr), time.Minute*10)
	if rdsStatus.Err() != nil {
		log.Println("caching failed")
		log.Println(rdsStatus.Err().Error())
	}

	return response, nil
}

func (ps *ProductService) GetAllProductByAdmin(ctx context.Context) ([]dto.ProductSize, error) {
	rkey := fmt.Sprintf("%s:product_admin", os.Getenv("RDB_KEY"))

	rsc := ps.redis.Get(ctx, rkey)
	if rsc.Err() == nil {
		var result []dto.ProductSize
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
		log.Println("product_admin cache miss")
	}

	data, err := ps.productRepository.GetAllProductSize(ctx, ps.db)
	if err != nil {
		return []dto.ProductSize{}, err
	}

	var response []dto.ProductSize
	for _, v := range data {
		response = append(response, dto.ProductSize{
			Id:    v.Id,
			Name:  v.Name,
			Price: v.Price,
		})
	}

	cacheStr, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		log.Println("failed to marshal")
	}

	rdsStatus := ps.redis.Set(ctx, rkey, string(cacheStr), time.Minute*10)
	if rdsStatus.Err() != nil {
		log.Println("caching failed")
		log.Println(rdsStatus.Err().Error())
	}

	return response, nil
}