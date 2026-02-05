package service

import (
	"context"
	"log"

	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/cache"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/dto"
	"github.com/NugrahaPancaWibisana/solid-coffee-be/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type MenuService struct {
	menuRepository *repository.MenuRepository
	redis          *redis.Client
	db             *pgxpool.Pool
}

func NewMenuService(menuRepository *repository.MenuRepository, rdb *redis.Client, db *pgxpool.Pool) *MenuService {
	return &MenuService{menuRepository: menuRepository, redis: rdb, db: db}
}

func (ms *MenuService) CreateMenu(ctx context.Context, req dto.MenuRequest, userID int, token string) error {
	if err := cache.CheckToken(ctx, ms.redis, userID, token); err != nil {
		return err
	}

	if err := ms.menuRepository.CreateMenu(ctx, ms.db, req); err != nil {
		return err
	}

	return nil
}

func (ms *MenuService) GetMenu(ctx context.Context, userID, menuID int, token string) (dto.Menu, error) {
	if err := cache.CheckToken(ctx, ms.redis, userID, token); err != nil {
		return dto.Menu{}, err
	}

	data, err := ms.menuRepository.GetMenu(ctx, ms.db, menuID)
	if err != nil {
		return dto.Menu{}, err
	}

	res := dto.Menu{
		ID:          data.ID,
		Discount:    data.Discount,
		ProductID:   data.ProductID,
		ProductName: data.ProductName,
		Stock:       data.Stock,
	}

	return res, nil
}

func (ms *MenuService) GetMenus(ctx context.Context, req dto.MenuParams, userID, menuID int, token string) ([]dto.Menu, int, error) {
	if err := cache.CheckToken(ctx, ms.redis, userID, token); err != nil {
		return nil, 0, err
	}

	tx, err := ms.db.Begin(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil, 0, err
	}

	totalPage, err := ms.menuRepository.GetTotalPage(ctx, tx, req)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	data, err := ms.menuRepository.GetMenus(ctx, tx, req)
	if err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("failed to commit", err.Error())
		return nil, 0, err
	}

	var response []dto.Menu
	for _, v := range data {
		response = append(response, dto.Menu{
			ID:          v.ID,
			Discount:    v.Discount,
			ProductID:   v.ProductID,
			ProductName: v.ProductName,
			Stock:       v.Stock,
		})
	}

	return response, totalPage, nil
}

func (ms *MenuService) UpdateMenu(ctx context.Context, req dto.UpdateMenuRequest, userID, menuID int, token string) error {
	if err := cache.CheckToken(ctx, ms.redis, userID, token); err != nil {
		return err
	}

	if err := ms.menuRepository.UpdateMenu(ctx, ms.db, req, menuID); err != nil {
		return err
	}

	return nil
}

func (ms *MenuService) DeleteMenu(ctx context.Context, userID, menuID int, token string) error {
	if err := cache.CheckToken(ctx, ms.redis, userID, token); err != nil {
		return err
	}

	if err := ms.menuRepository.DeleteMenu(ctx, ms.db, menuID); err != nil {
		return err
	}

	return nil
}
