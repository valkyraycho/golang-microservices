package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/valkyraycho/go-microservices/account"
	"github.com/valkyraycho/go-microservices/catalog"
	pb "github.com/valkyraycho/go-microservices/order/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type orderServer struct {
	pb.OrderServiceServer
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
}

func ListenGRPC(s Service, accountServiceURL string, catalogServiceURL string, port int) error {
	accountClient, err := account.NewClient(accountServiceURL)
	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogServiceURL)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}

	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, &orderServer{service: s, accountClient: accountClient, catalogClient: catalogClient})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *orderServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("Error while finding account: ", err)
		return nil, errors.New("account not found")
	}

	productIDs := make([]string, len(r.Products))
	quantities := make(map[string]uint32, len(r.Products))

	for i, p := range r.Products {
		productIDs[i] = p.ProductId
		quantities[p.ProductId] = p.Quantity
	}

	products, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("Error while getting products: ", err)
		return nil, errors.New("products not found")
	}

	orderedProducts := make([]OrderedProduct, 0, len(products))
	for _, p := range products {
		if quantity := quantities[p.ID]; quantity > 0 {
			orderedProducts = append(orderedProducts, OrderedProduct{
				ID:          p.ID,
				Quantity:    quantity,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	order, err := s.service.PostOrder(ctx, r.AccountId, orderedProducts)
	if err != nil {
		log.Println("Error posting order: ", err)
		return nil, errors.New("could not post order")
	}

	pbOrderedProducts := make([]*pb.Order_OrderProduct, len(order.Products))
	for i, p := range order.Products {
		pbOrderedProducts[i] = &pb.Order_OrderProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		}
	}
	createdAt, err := order.CreatedAt.MarshalBinary()
	if err != nil {
		log.Printf("Error marshaling timestamp: %v", err)
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}
	return &pb.PostOrderResponse{
		Order: &pb.Order{
			Id:         order.ID,
			AccountId:  order.AccountID,
			TotalPrice: order.TotalPrice,
			CreatedAt:  createdAt,
			Products:   pbOrderedProducts,
		},
	}, nil
}
