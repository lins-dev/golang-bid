package user

import (
	"context"

	"github.com/lins-dev/golang-bid.git/internal/validator"
)

type LoginUserReq struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (req LoginUserReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.Email), "email", "this field cannot be empty")
	eval.CheckField(validator.Matches(req.Email, validator.EmailRx), "email", "this field must be a valid email")
	eval.CheckField(validator.NotBlank(req.Password), "password", "this field cannot be empty")

	return eval
}