package main

import (
	"context"
	"errors"
	"time"

	orderServ "github.com/valkyraycho/go-microservices/order"
)

type mutationResolver struct {
	server *Server
}

var ErrInvalidParameter = errors.New("invalid parameter")

func (r *mutationResolver) CreateAccount(ctx context.Context, account AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	a, err := r.server.accountClient.PostAccount(ctx, account.Name)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   a.ID,
		Name: a.Name,
	}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, product ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	p, err := r.server.catalogClient.PostProduct(ctx, product.Name, product.Description, product.Price)
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, order OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	orderedProducts := []orderServ.OrderedProduct{}

	for _, p := range order.Products {
		if p.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}

		orderedProducts = append(orderedProducts, orderServ.OrderedProduct{
			ID:       p.ID,
			Quantity: uint32(p.Quantity),
		})
	}

	o, err := r.server.orderClient.PostOrder(ctx, order.AccountID, orderedProducts)
	if err != nil {
		return nil, err
	}

	return &Order{
		ID:         o.ID,
		CreatedAt:  o.CreatedAt,
		TotalPrice: o.TotalPrice,
	}, nil
}
