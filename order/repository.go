package order

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	CreateOrder(ctx context.Context, o Order) error
	GetOrdersForAccount(ctx context.Context, accoundID string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgresRepository) CreateOrder(ctx context.Context, o Order) error {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES($1, $2, $3, $4)",
		o.ID,
		o.CreatedAt,
		o.AccountID,
		o.TotalPrice,
	)

	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))

	if err != nil {
		return err
	}

	for _, p := range o.Products {
		_, err = stmt.ExecContext(ctx, o.ID, p.ID, p.Quantity)
		if err != nil {
			return err
		}
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}
	stmt.Close()
	return nil
}

func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, accoundID string) ([]Order, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT 
		o.id, 
		o.created_at, 
		o.account_id, 
		o.total_price::money::numeric::float8,	   
		op.product_id, 
		op.quantity
		FROM orders o
		JOIN order_products op ON o.id = op.order_id
		WHERE o.account_id = $1
		ORDER BY o.id`,
		accoundID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []Order{}
	order := &Order{}
	lastOrder := &Order{}
	orderedProduct := &OrderedProduct{}
	products := []OrderedProduct{}

	for rows.Next() {
		if err := rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.AccountID,
			&order.TotalPrice,
			&orderedProduct.ID,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}

		if lastOrder.ID != "" && order.ID != lastOrder.ID {
			orders = append(orders, Order{
				ID:         lastOrder.ID,
				CreatedAt:  lastOrder.CreatedAt,
				TotalPrice: lastOrder.TotalPrice,
				AccountID:  lastOrder.AccountID,
				Products:   products,
			})
			products = []OrderedProduct{}
		}

		products = append(products, *orderedProduct)
		*lastOrder = *order
	}

	if lastOrder.ID != "" {
		orders = append(orders, Order{
			ID:         lastOrder.ID,
			CreatedAt:  lastOrder.CreatedAt,
			TotalPrice: lastOrder.TotalPrice,
			AccountID:  lastOrder.AccountID,
			Products:   products,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
