package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"product_auth/domain"
	"product_auth/dto/errs"
	"product_auth/dto/logger"

	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	FindByUsernameAndPassword(ctx context.Context, username string, password string) (*domain.Login, *errs.AppError)
    GenerateAndSaveRefreshTokenToStore(ctx context.Context, authToken domain.AuthToken, customer_id string) (string, *errs.AppError)
	RefreshTokenExists(ctx context.Context, refreshToken string) *errs.AppError 
}

type AuthRepositoryDb struct {
	client *sqlx.DB
}

func (d AuthRepositoryDb) FindByEmailAndPassword(ctx context.Context, email string, password string) (*domain.Login, *errs.AppError) {
	var login domain.Login
	sqlVerify := `SELECT customers_id, username, role, age, address, gender FROM customers WHERE email = $1 AND password = $2`
	err := d.client.GetContext(ctx, &login, sqlVerify, email, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewAuthenticationError("invalid credentials", err)
		} else {
			logger.Error("Error while verifying login request from database: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error", err)
		}
	}
	fmt.Println(login)
	return &login, nil
}

func (d AuthRepositoryDb) RefreshTokenExists(ctx context.Context, refreshToken string) *errs.AppError {
	sqlSelect := "select refresh_token from refresh_token_store where refresh_token = $1"
	var token string 
	err := d.client.GetContext(ctx, &token, sqlSelect, refreshToken)
    if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewAuthenticationError("refresh token not registered in the store", err)
		} else {
			logger.Error("Unexpected database error: "+ err.Error())
			return errs.NewUnexpectedError("Unexpected database error", err)
		}
	}
	return nil
}  

<<<<<<< HEAD
func (d AuthRepositoryDb) GenerateAndSaveRefreshTokenToStore(ctx context.Context, authToken domain.AuthToken, customer_id string) (string, *errs.AppError) {
=======
func (d AuthRepositoryDb) GenerateAndSaveRefreshTokenToStore(ctx context.Context, authToken domain.AuthToken, claims domain.AccessTokenClaims) (string, *errs.AppError) {
>>>>>>> 4e6c5d0f9f51aad147517bd831c7187f2a4b6ed0
	var appErr *errs.AppError 
	var refreshToken string 
	if refreshToken, appErr = authToken.NewRefreshToken(); appErr != nil {
		return "", appErr 
	}
	fmt.Println(refreshToken)
	sqlInsert := "INSERT INTO refresh_token_store (refresh_token, customers_id) VALUES ($1, $2)"
<<<<<<< HEAD
	_, err := d.client.ExecContext(ctx, sqlInsert, refreshToken, customer_id)
=======
	_, err := d.client.ExecContext(ctx, sqlInsert, refreshToken, claims.CustomersId)
>>>>>>> 4e6c5d0f9f51aad147517bd831c7187f2a4b6ed0
    if err != nil {
		fmt.Println(err)
		logger.Error("unexpected database error: " + err.Error())
		return "", errs.NewUnexpectedError("unexpected database error", err)
	}
	return refreshToken, nil 
}

func NewAuthRepository(client *sqlx.DB) AuthRepositoryDb {
	return AuthRepositoryDb{client}
}
