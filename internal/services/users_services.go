package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lins-dev/golang-bid.git/internal/store/pgstore"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicatedEmailOrPassword = errors.New("invalid username or email already exists")

type UserService struct{
	pool *pgxpool.Pool
	queries *pgstore.Queries
}

func NewUserService (pool *pgxpool.Pool) UserService {
	return UserService{
		pool: pool,
		queries: pgstore.New(pool),
	}
}

func (us *UserService) CreateUser(
	ctx context.Context,
	userName string,
	email string,
	password string,
	bio string) (
	pgstore.User,
	error) {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			return pgstore.User{}, err
		}

		args := pgstore.CreateUserParams{
			UserName: userName,
			Email: email,
			PasswordHash: hash,
			Bio: bio,
		}
		user, err := us.queries.CreateUser(ctx, args)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return pgstore.User{}, ErrDuplicatedEmailOrPassword
			}
			return pgstore.User{}, err
		}

		return user, nil
	}
