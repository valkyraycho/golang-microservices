package order

import (
	"context"
	"time"

	pb "github.com/valkyraycho/go-microservices/order/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{conn, pb.NewOrderServiceClient(conn)}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	pbProducts := []*pb.PostOrderRequest_OrderProduct{}

	for _, p := range products {
		pbProducts = append(pbProducts, &pb.PostOrderRequest_OrderProduct{
			ProductId: p.ID,
			Quantity:  p.Quantity,
		})
	}
	res, err := c.service.PostOrder(ctx, &pb.PostOrderRequest{
		AccountId: accountID,
		Products:  pbProducts,
	})

	if err != nil {
		return nil, err
	}
	newOrderCreatedAt := time.Time{}
	newOrderCreatedAt.UnmarshalBinary(res.Order.CreatedAt)

	return &Order{
		ID:         res.Order.Id,
		TotalPrice: res.Order.TotalPrice,
		AccountID:  res.Order.AccountId,
		Products:   products,
		CreatedAt:  newOrderCreatedAt,
	}, nil
}

func (c *Client) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	res, err := c.service.GetOrdersForAccount(ctx, &pb.GetOrdersForAccountRequest{AccountId: accountID})
	if err != nil {
		return nil, err
	}

	orders := []Order{}

	for _, pbOrder := range res.Orders {
		newOrder := Order{
			ID:         pbOrder.Id,
			AccountID:  pbOrder.AccountId,
			TotalPrice: pbOrder.TotalPrice,
		}
		newOrder.CreatedAt = time.Time{}
		newOrder.CreatedAt.UnmarshalBinary(pbOrder.CreatedAt)
		products := []OrderedProduct{}
		for _, p := range pbOrder.Products {
			products = append(products, OrderedProduct{
				ID:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    p.Quantity,
			})
		}
		newOrder.Products = products
		orders = append(orders, newOrder)
	}
	return orders, err
}
