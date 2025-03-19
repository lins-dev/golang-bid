package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lins-dev/golang-bid.git/internal/store/pgstore"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicatedEmailOrUsername = errors.New("invalid username or email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

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
				return pgstore.User{}, ErrDuplicatedEmailOrUsername
			}
			return pgstore.User{}, err
		}

		return user, nil
	}

func (us *UserService) AuthUser(ctx context.Context, email string, password string) (pgstore.User, error) {
	user, err := us.queries.FindUserByEmail(ctx, email)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return pgstore.User{}, ErrInvalidCredentials
		}
		return pgstore.User{}, err
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return pgstore.User{}, ErrInvalidCredentials
		}
		return pgstore.User{}, err
	}

	return user, nil
}

func (us *UserService) FindUserByUuid(ctx context.Context, user_uuid uuid.UUID) (pgstore.User, error) {
	user, err := us.queries.GetUserByUuid(ctx, user_uuid)
	if err != nil {
		return pgstore.User{}, err
	}

	return user, nil
}
