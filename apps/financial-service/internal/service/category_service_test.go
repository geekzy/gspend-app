package service

import (
	"context"
	"testing"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockCategoryRepository for testing
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) ListByUserID(ctx context.Context, userID primitive.ObjectID, categoryType domain.CategoryType) ([]*domain.Category, error) {
	args := m.Called(ctx, userID, categoryType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) ListSystem(ctx context.Context, categoryType domain.CategoryType) ([]*domain.Category, error) {
	args := m.Called(ctx, categoryType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCategoryService_InitializeSystemCategories(t *testing.T) {
	categoryRepo := new(MockCategoryRepository)
	svc := NewCategoryService(categoryRepo)

	t.Run("Success - First Time Initialization", func(t *testing.T) {
		// Mock that no system categories exist
		categoryRepo.On("ListSystem", mock.Anything, domain.CategoryType("")).Return([]*domain.Category{}, nil).Once()
		
		// Mock successful creation of all categories
		familyCategories := svc.getFamilyCategories()
		for range familyCategories {
			categoryRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil).Once()
		}

		err := svc.InitializeSystemCategories(context.Background())

		assert.NoError(t, err)
		categoryRepo.AssertExpectations(t)
	})

	t.Run("Success - Skip When Categories Already Exist", func(t *testing.T) {
		existingCategories := []*domain.Category{
			{
				ID:       primitive.NewObjectID(),
				Name:     "Groceries",
				Type:     domain.CategoryTypeExpense,
				IsSystem: true,
			},
		}

		// Mock that system categories already exist
		categoryRepo.On("ListSystem", mock.Anything, domain.CategoryType("")).Return(existingCategories, nil).Once()
		// No Create calls should be made

		err := svc.InitializeSystemCategories(context.Background())

		assert.NoError(t, err)
		categoryRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Failure", func(t *testing.T) {
		// Mock repository error
		categoryRepo.On("ListSystem", mock.Anything, domain.CategoryType("")).Return(nil, assert.AnError).Once()

		err := svc.InitializeSystemCategories(context.Background())

		assert.Error(t, err)
		categoryRepo.AssertExpectations(t)
	})
}

func TestCategoryService_ValidateSystemCategoryProtection(t *testing.T) {
	categoryRepo := new(MockCategoryRepository)
	svc := NewCategoryService(categoryRepo)

	systemCategoryID := primitive.NewObjectID()
	userCategoryID := primitive.NewObjectID()

	t.Run("Error - System Category Protected", func(t *testing.T) {
		systemCategory := &domain.Category{
			ID:       systemCategoryID,
			Name:     "Groceries",
			Type:     domain.CategoryTypeExpense,
			IsSystem: true,
		}

		categoryRepo.On("GetByID", mock.Anything, systemCategoryID).Return(systemCategory, nil).Once()

		err := svc.ValidateSystemCategoryProtection(context.Background(), systemCategoryID.Hex())

		assert.Error(t, err)
		assert.Equal(t, domain.ErrSystemCategoryProtected, err)
		categoryRepo.AssertExpectations(t)
	})

	t.Run("Success - User Category Not Protected", func(t *testing.T) {
		userCategory := &domain.Category{
			ID:       userCategoryID,
			Name:     "Custom Category",
			Type:     domain.CategoryTypeExpense,
			IsSystem: false,
		}

		categoryRepo.On("GetByID", mock.Anything, userCategoryID).Return(userCategory, nil).Once()

		err := svc.ValidateSystemCategoryProtection(context.Background(), userCategoryID.Hex())

		assert.NoError(t, err)
		categoryRepo.AssertExpectations(t)
	})

	t.Run("Success - Category Not Found", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()

		categoryRepo.On("GetByID", mock.Anything, nonExistentID).Return(nil, nil).Once()

		err := svc.ValidateSystemCategoryProtection(context.Background(), nonExistentID.Hex())

		assert.NoError(t, err)
		categoryRepo.AssertExpectations(t)
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		err := svc.ValidateSystemCategoryProtection(context.Background(), "invalid-id")

		assert.Error(t, err)
		// No repository calls should be made for invalid ID
	})
}

func TestCategoryService_GetFamilyCategories(t *testing.T) {
	categoryRepo := new(MockCategoryRepository)
	svc := NewCategoryService(categoryRepo)

	familyCategories := svc.getFamilyCategories()

	t.Run("Verify Family Categories Structure", func(t *testing.T) {
		assert.Greater(t, len(familyCategories), 30, "Should have a comprehensive set of family categories")

		// Count categories by type
		expenseCount := 0
		incomeCount := 0
		for _, cat := range familyCategories {
			if cat.Type == domain.CategoryTypeExpense {
				expenseCount++
			} else if cat.Type == domain.CategoryTypeIncome {
				incomeCount++
			}
		}

		assert.Greater(t, expenseCount, 20, "Should have many expense categories")
		assert.Greater(t, incomeCount, 3, "Should have several income categories")

		// Verify essential family categories exist
		categoryNames := make(map[string]bool)
		for _, cat := range familyCategories {
			categoryNames[cat.Name] = true
		}

		// Essential family categories
		essentialCategories := []string{
			"Groceries", "Childcare", "School Expenses", "Medical", 
			"Rent/Mortgage", "Utilities", "Car Payment", "Salary",
		}

		for _, essential := range essentialCategories {
			assert.True(t, categoryNames[essential], "Should include essential category: %s", essential)
		}
	})

	t.Run("Verify Category Properties", func(t *testing.T) {
		for _, cat := range familyCategories {
			assert.NotEmpty(t, cat.Name, "Category name should not be empty")
			assert.NotEmpty(t, cat.Icon, "Category icon should not be empty")
			assert.NotEmpty(t, cat.Color, "Category color should not be empty")
			assert.True(t, cat.Type == domain.CategoryTypeExpense || cat.Type == domain.CategoryTypeIncome, "Category type should be valid")
			assert.Greater(t, cat.SortOrder, 0, "Sort order should be positive")
		}
	})

	t.Run("Verify Sort Order Organization", func(t *testing.T) {
		// Verify categories are organized by sort order ranges
		housingCategories := 0
		foodCategories := 0
		childrenCategories := 0
		incomeCategories := 0

		for _, cat := range familyCategories {
			switch {
			case cat.SortOrder >= 1 && cat.SortOrder <= 10:
				housingCategories++
			case cat.SortOrder >= 11 && cat.SortOrder <= 20:
				foodCategories++
			case cat.SortOrder >= 21 && cat.SortOrder <= 40:
				childrenCategories++
			case cat.SortOrder >= 91 && cat.SortOrder <= 100:
				incomeCategories++
			}
		}

		assert.Greater(t, housingCategories, 0, "Should have housing categories")
		assert.Greater(t, foodCategories, 0, "Should have food categories")
		assert.Greater(t, childrenCategories, 0, "Should have children categories")
		assert.Greater(t, incomeCategories, 0, "Should have income categories")
	})
}