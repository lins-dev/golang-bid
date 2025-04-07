package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lins-dev/golang-bid.git/internal/store/pgstore"
)

type BidsService struct {
	pool *pgxpool.Pool
	queries *pgstore.Queries
}

var ErrNewBidLowerThanPrevious = errors.New("current bid is lower than the previous one")
var ErrNewBidLowerThanBasePrice = errors.New("current bid is lower than the base price")

func NewBidService(pool *pgxpool.Pool) BidsService{
	return BidsService {
		pool: pool,
		queries: pgstore.New(pool),
	}
}

func (bs *BidsService) CreateBid(
	ctx context.Context, 
	productId int32, 
	bidderId int32, 
	bidAmount int32,
) (pgstore.Bid, error) {

	args := pgstore.CreateBidParams{
		BidderID: bidderId,
		ProductID: productId,
		BidAmount: bidAmount,
	}

	bid, err := bs.queries.CreateBid(ctx, args)
	if err != nil {
		return pgstore.Bid{}, err
	}

	return bid, nil

}

func (bs *BidsService) PlaceBid(
	ctx context.Context,
	productId int32,
	bidderId int32,
	bidAmount int32,
) (pgstore.Bid, error) {

	args := pgstore.CreateBidParams{
		ProductID: productId,
		BidderID: bidderId,
		BidAmount: bidAmount,
	}

	product, err := bs.queries.GetProduct(ctx, args.ProductID)
	if err != nil {
		return pgstore.Bid{}, ErrProductNotFound
	}

	mostValueBid, err := bs.queries.GetHightestBidByProductId(ctx, productId)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
	}
	slog.Info("Log of bids", "args.BidAmount", args.BidAmount, "product.Price", product.Price, "mostValueBid.BidAmount", mostValueBid.BidAmount)
	if args.BidAmount < product.Price {
		return pgstore.Bid{}, ErrNewBidLowerThanBasePrice
	}

	if args.BidAmount < mostValueBid.BidAmount {
		return pgstore.Bid{}, ErrNewBidLowerThanPrevious
	}

	bid, err := bs.queries.CreateBid(ctx, args)
	if err != nil {
		return pgstore.Bid{}, err
	}

	return bid, nil
	
}