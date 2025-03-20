package services

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lins-dev/golang-bid.git/internal/store/pgstore"
)

type ProductService struct{
	pool *pgxpool.Pool
	queries *pgstore.Queries
}

var ErrProductNotFound = errors.New("product not found")

func NewProductService(pool *pgxpool.Pool) ProductService {
	return ProductService{
		pool: pool,
		queries: pgstore.New(pool),
	}
}

func (ps *ProductService) CreateProduct(
	ctx context.Context,
	sellerId int32,
	productName string,
	description string,
	price int32,
	auction_end time.Time,
) (pgstore.Product, error) {
	args := pgstore.CreateProductParams{
		SellerID: sellerId,
		ProductName: productName,
		Description: description,
		Price: price,
		AuctionEnd: auction_end,
	}
	product, err := ps.queries.CreateProduct(ctx, args)
	if err != nil {
		return pgstore.Product{}, err
	}

	return product, nil
}