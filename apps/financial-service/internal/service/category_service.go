package service

import (
	"context"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryService struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryService(categoryRepo domain.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) error {
	return s.categoryRepo.Create(ctx, category)
}

func (s *CategoryService) GetCategory(ctx context.Context, id string) (*domain.Category, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.categoryRepo.GetByID(ctx, objectID)
}

func (s *CategoryService) ListUserCategories(ctx context.Context, userID string, categoryType domain.CategoryType) ([]*domain.Category, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	return s.categoryRepo.ListByUserID(ctx, objectID, categoryType)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, category *domain.Category) error {
	return s.categoryRepo.Update(ctx, category)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.categoryRepo.Delete(ctx, objectID)
}
