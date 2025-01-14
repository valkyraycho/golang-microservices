package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
)

type Order struct {
	ID         string
	CreatedAt  time.Time
	TotalPrice float64
	AccountID  string
	Products   []OrderedProduct
}

type OrderedProduct struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    uint32
}

type Service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type orderService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &orderService{r}
}

func (s *orderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	totalPrice := 0.0

	for _, p := range products {
		totalPrice += p.Price * float64(p.Quantity)
	}

	order := Order{
		ID:         ksuid.New().String(),
		CreatedAt:  time.Now().UTC(),
		AccountID:  accountID,
		Products:   products,
		TotalPrice: totalPrice,
	}

	if err := s.repository.CreateOrder(ctx, order); err != nil {
		return nil, err
	}
	return &order, nil
}
func (s *orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	return s.repository.GetOrdersForAccount(ctx, accountID)
}
