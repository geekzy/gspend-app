package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/dto"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestIncomeHandler_Create(t *testing.T) {
	e := echo.New()
	incomeRepo := new(MockIncomeRepository)
	svc := service.NewIncomeService(incomeRepo)
	h := NewIncomeHandler(svc)

	userID := primitive.NewObjectID()

	t.Run("Success", func(t *testing.T) {
		reqBody, _ := json.Marshal(dto.CreateIncomeRequest{
			Source:        "Salary",
			Amount:        5000.0,
			Frequency:     "monthly",
			EffectiveDate: time.Now(),
		})

		incomeRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/incomes", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		if assert.NoError(t, h.Create(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})
}

func TestIncomeHandler_List(t *testing.T) {
	e := echo.New()
	incomeRepo := new(MockIncomeRepository)
	svc := service.NewIncomeService(incomeRepo)
	h := NewIncomeHandler(svc)

	userID := primitive.NewObjectID()

	t.Run("Success", func(t *testing.T) {
		incomeRepo.On("ListByUserID", mock.Anything, userID).Return([]*domain.Income{}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/incomes", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		if assert.NoError(t, h.List(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}

func TestIncomeHandler_Update(t *testing.T) {
	e := echo.New()
	incomeRepo := new(MockIncomeRepository)
	svc := service.NewIncomeService(incomeRepo)
	h := NewIncomeHandler(svc)

	incomeID := primitive.NewObjectID()

	t.Run("Success", func(t *testing.T) {
		reqBody, _ := json.Marshal(dto.CreateIncomeRequest{
			Source:        "Bonus",
			Amount:        2000.0,
			Frequency:     "one-time",
			EffectiveDate: time.Now(),
		})

		incomeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/incomes/"+incomeID.Hex(), strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(incomeID.Hex())

		if assert.NoError(t, h.Update(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}

func TestIncomeHandler_Delete(t *testing.T) {
	e := echo.New()
	incomeRepo := new(MockIncomeRepository)
	svc := service.NewIncomeService(incomeRepo)
	h := NewIncomeHandler(svc)

	incomeID := primitive.NewObjectID()

	t.Run("Success", func(t *testing.T) {
		incomeRepo.On("Delete", mock.Anything, incomeID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/incomes/"+incomeID.Hex(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(incomeID.Hex())

		if assert.NoError(t, h.Delete(c)) {
			assert.Equal(t, http.StatusNoContent, rec.Code)
		}
	})
}
