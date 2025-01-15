package main

import (
	"context"
	"log"
	"time"
)

type queryResolver struct {
	server *Server
}

func (p *PaginationInput) bounds() (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(100)
	if p.Skip != nil {
		skipValue = uint64(*p.Skip)
	}
	if p.Take != nil {
		takeValue = uint64(*p.Take)
	}
	return skipValue, takeValue
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if id != nil {
		a, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Account{
			{
				ID:   a.ID,
				Name: a.Name,
			},
		}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}
	res, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}

	for _, a := range res {
		accounts = append(accounts, &Account{
			ID:   a.ID,
			Name: a.Name,
		})
	}
	return accounts, nil
}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if id != nil {
		p, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*Product{
			{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			},
		}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	q := ""
	if query != nil {
		q = *query
	}

	productList, err := r.server.catalogClient.GetProducts(ctx, skip, take, nil, q)
	if err != nil {
		return nil, err
	}

	products := []*Product{}
	for _, p := range productList {
		products = append(products, &Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	return products, nil
}
