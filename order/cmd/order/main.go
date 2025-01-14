package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"github.com/valkyraycho/go-microservices/order"
)

type Config struct {
	DatabaseURL       string `envconfig:"DATABASE_URL"`
	AccountServiceURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogServiceURL string `envconfig:"CATALOG_SERVICE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r order.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		r, err = order.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Fatal(err)
			return err
		}
		return nil
	})
	defer r.Close()

	log.Println("Listening on port 8080...")
	s := order.NewService(r)
	log.Fatal(order.ListenGRPC(s, cfg.AccountServiceURL, cfg.CatalogServiceURL, 8080))
}
