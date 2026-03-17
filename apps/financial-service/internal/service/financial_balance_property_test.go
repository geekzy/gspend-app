package service

import (
    "math/rand"
    "testing"
    "time"

"math"

    "github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// computeBalance calculates total balance from a slice of transactions.
func computeBalance(ts []domain.Transaction) float64 {
    bal := 0.0
    for _, t := range ts {
        switch t.Type {
        case domain.TransactionTypeIncome:
            bal += t.Amount
        case domain.TransactionTypeExpense:
            bal -= t.Amount
        }
    }
    return bal
}

// TestProperty_FinancialBalanceCalculation runs a random‑data property check.
func TestProperty_FinancialBalanceCalculation(t *testing.T) {
    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 1000; i++ {
        n := rnd.Intn(20) + 1
        ts := make([]domain.Transaction, n)
        for j := 0; j < n; j++ {
            amt := rnd.Float64() * 1000
            var ttype domain.TransactionType
            if rnd.Intn(2) == 0 {
                ttype = domain.TransactionTypeIncome
            } else {
                ttype = domain.TransactionTypeExpense
            }
            ts[j] = domain.Transaction{
                ID:             primitive.NewObjectID(),
                UserID:         primitive.NewObjectID(),
                CategoryID:     primitive.NewObjectID(),
                BudgetID:       nil,
                Type:           ttype,
                Amount:         amt,
                Description:    "",
                TransactionDate: time.Now(),
                PaymentMethod: "",
                Notes:          "",
                Metadata:        domain.TransactionMetadata{},
                CreatedAt:      time.Now(),
                UpdatedAt:      time.Now(),
            }
        }
        inc, exp := 0.0, 0.0
        for _, tx := range ts {
            if tx.Type == domain.TransactionTypeIncome {
                inc += tx.Amount
            } else {
                exp += tx.Amount
            }
        }
        if math.Abs(computeBalance(ts)- (inc-exp)) > 1e-9 {
            t.Fatalf("balance mismatch at iteration %d: expected %f, got %f", i, inc-exp, computeBalance(ts))
        }
    }
}
