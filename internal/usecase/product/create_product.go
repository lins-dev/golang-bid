package product

import (
	"context"
	"time"

	"github.com/lins-dev/golang-bid.git/internal/validator"
)

type CreateProductReq struct {
	SellerID    int32     `json:"seller_id"`
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	Price       float32     `json:"price"`
	AuctionEnd  time.Time `json:"auction_end"`
}

const minAuctionDuration = 2 * time.Hour

func (req CreateProductReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator
	
	eval.CheckField(req.Price > 0, "price", "this field must be greater than 0")

	eval.CheckField(validator.NotBlank(req.ProductName), "product_name", "this field cannot be empty")
	eval.CheckField(validator.MinChars(req.ProductName, 3) && validator.MaxChars(req.ProductName, 200), "product_name", "this field must have length between 3 and 200 characters")

	eval.CheckField(validator.NotBlank(req.Description), "description", "this field cannot be empty")
	eval.CheckField(validator.MinChars(req.Description, 5) && validator.MaxChars(req.Description, 300), "description", "this field must have betlenghtween 3 and 300 caharacters")

	eval.CheckField(req.AuctionEnd.Sub(time.Now()) >= minAuctionDuration, "auction_end", "this filed must be at least 2 hours duration")

	return eval
}