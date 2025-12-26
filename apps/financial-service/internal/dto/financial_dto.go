package dto

import "time"

// Category DTOs
type CategoryResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Icon      string `json:"icon"`
	Color     string `json:"color"`
	IsSystem  bool   `json:"isSystem"`
	SortOrder int    `json:"sortOrder"`
}

type CreateCategoryRequest struct {
	Name      string `json:"name" validate:"required"`
	Type      string `json:"type" validate:"required,oneof=income expense"`
	Icon      string `json:"icon" validate:"required"`
	Color     string `json:"color" validate:"required"`
	SortOrder int    `json:"sortOrder"`
}

type UpdateCategoryRequest struct {
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Color     string `json:"color"`
	SortOrder int    `json:"sortOrder"`
}

// Income DTOs
type IncomeResponse struct {
	ID            string    `json:"id"`
	Source        string    `json:"source"`
	Amount        float64   `json:"amount"`
	Frequency     string    `json:"frequency"`
	EffectiveDate time.Time `json:"effectiveDate"`
}

type CreateIncomeRequest struct {
	Source        string    `json:"source" validate:"required"`
	Amount        float64   `json:"amount" validate:"required,gt=0"`
	Frequency     string    `json:"frequency" validate:"required,oneof=one-time daily weekly monthly yearly"`
	EffectiveDate time.Time `json:"effectiveDate" validate:"required"`
}

// Budget DTOs
type BudgetItemDTO struct {
	ID            string  `json:"id"`
	CategoryID    string  `json:"categoryId"`
	CategoryName  string  `json:"categoryName"`
	PlannedAmount float64 `json:"plannedAmount"`
	SpentAmount   float64 `json:"spentAmount"`
	Notes         string  `json:"notes"`
}

type BudgetResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	PeriodType  string          `json:"periodType"`
	StartDate   time.Time       `json:"startDate"`
	EndDate     time.Time       `json:"endDate"`
	TotalAmount float64         `json:"totalAmount"`
	Items       []BudgetItemDTO `json:"items"`
}

type CreateBudgetRequest struct {
	Name        string          `json:"name" validate:"required"`
	PeriodType  string          `json:"periodType" validate:"required,oneof=monthly quarterly yearly"`
	StartDate   time.Time       `json:"startDate" validate:"required"`
	EndDate     time.Time       `json:"endDate" validate:"required"`
	TotalAmount float64         `json:"totalAmount" validate:"required,gt=0"`
	Items       []CreateBudgetItemRequest `json:"items" validate:"required,dive"`
}

type CreateBudgetItemRequest struct {
	CategoryID    string  `json:"categoryId" validate:"required"`
	CategoryName  string  `json:"categoryName" validate:"required"`
	PlannedAmount float64 `json:"plannedAmount" validate:"required,gt=0"`
	Notes         string  `json:"notes"`
}

type UpdateBudgetItemRequest struct {
	CategoryID    string  `json:"categoryId" validate:"required"`
	CategoryName  string  `json:"categoryName" validate:"required"`
	PlannedAmount float64 `json:"plannedAmount" validate:"required,gt=0"`
	Notes         string  `json:"notes"`
}

// Transaction DTOs
type TransactionResponse struct {
	ID              string    `json:"id"`
	CategoryID      string    `json:"categoryId"`
	CategoryName    string    `json:"categoryName"`
	BudgetID        string    `json:"budgetId,omitempty"`
	BudgetName      string    `json:"budgetName,omitempty"`
	Type            string    `json:"type"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transactionDate"`
	PaymentMethod   string    `json:"paymentMethod"`
	Notes           string    `json:"notes"`
}

type CreateTransactionRequest struct {
	CategoryID      string    `json:"categoryId" validate:"required"`
	BudgetID        string    `json:"budgetId"`
	Type            string    `json:"type" validate:"required,oneof=income expense"`
	Amount          float64   `json:"amount" validate:"required,gt=0"`
	Description     string    `json:"description" validate:"required"`
	TransactionDate time.Time `json:"transactionDate" validate:"required"`
	PaymentMethod   string    `json:"paymentMethod" validate:"required"`
	Notes           string    `json:"notes"`
}
