package user

import (
	"context"

	"github.com/lins-dev/golang-bid.git/internal/validator"
)

type CreateUserReq struct {
	UserName     string             `json:"user_name"`
	Email        string             `json:"email"`
	Password 	 string             `json:"password"`
	Bio          string             `json:"bio"`
}

func (req CreateUserReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.UserName), "user_name", "this field cannot be empty")
	eval.CheckField(validator.NotBlank(req.Bio), "bio", "this field cannot be empty")
	eval.CheckField(validator.MinChars(req.Bio, 5), "bio", "this field must be at least size 5")
	eval.CheckField(validator.MaxChars(req.Bio, 255), "bio", "this size field must be less than 255")
	eval.CheckField(validator.MinChars(req.Password, 3), "password", "this field must be at least size 3")
	eval.CheckField(validator.MaxChars(req.Password, 255), "bio", "this size field must be less than 255")
	eval.CheckField(validator.NotBlank(req.Email), "email", "this field cannot be empty")
	eval.CheckField(validator.Matches(req.Email, validator.EmailRx), "email", "this field must be a valid email")

	return eval
}