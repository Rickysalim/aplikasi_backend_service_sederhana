package domain

import (
	"errors"
	"product_auth/dto/errs"
	"product_auth/dto/logger"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const HMAC_SAMPLE_SECRET = "@SECRETKEY123@"
const ACCESS_TOKEN_DURATION = time.Hour
const REFRESH_TOKEN_DURATION = time.Hour * 24 * 30

var now = time.Now().Add(ACCESS_TOKEN_DURATION)
var refresh = time.Now().Add(REFRESH_TOKEN_DURATION)
var newJwtTime = jwt.NewNumericDate(now)
var newJwtTimeRefresj = jwt.NewNumericDate(refresh)


type AuthToken struct {
	token *jwt.Token
}

type Login struct {
	CustomersId string `db:"customers_id"`
	Username    string `db:"username"`
	Role        string `db:"role"`
	Age         int    `db:"age"`
	Address     string `db:"address"`
	Gender      string `db:"gender"`
}

type AccessTokenClaims struct {
	CustomersId string `json:"customers_id"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	Age         int    `json:"age"`
	Address     string `json:"address"`
	Gender      string `json:"gender"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	TokenType   string `json:"token_type"`
	CustomersId string `json:"customers_id"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	Age         int    `json:"age"`
	Address     string `json:"address"`
	Gender      string `json:"gender"`
	jwt.RegisteredClaims
}

func (c AccessTokenClaims) IsSameRole(urlParams map[string]string) error {
	if c.Role != urlParams["role"] {
		return errors.New("role not same")
	}
	return nil
}

func (c AccessTokenClaims) IsRequestVerifiedWithTokenClaims(urlParams map[string]string) error {
	if c.CustomersId != urlParams["customer_id"] {
		return errors.New("customer_id doesn't match")
	}
	return nil
}

func (l Login) ClaimsForAccessToken() AccessTokenClaims {
	return AccessTokenClaims{
		CustomersId: l.CustomersId,
		Username:    l.Username,
		Role:        l.Role,
		Age:         l.Age,
		Address:     l.Address,
		Gender:      l.Gender,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: newJwtTime,
		},
	}
}

func (c RefreshTokenClaims) AccessTokenClaims() AccessTokenClaims {
	return AccessTokenClaims{
		CustomersId: c.CustomersId,
		Username: c.Username,
		Role: c.Role,
		Age: c.Age,
		Address: c.Address,
		Gender: c.Gender,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: newJwtTime,
		},
	}
}

func (c AccessTokenClaims) RefreshTokenClaims() RefreshTokenClaims {
	return RefreshTokenClaims{
		TokenType: "refresh_token",
		CustomersId: c.CustomersId,
		Username: c.Username,
		Role: c.Role,
		Age: c.Age,
		Address: c.Address,
		Gender: c.Gender,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: newJwtTime,
		},
	}
}

func NewAuthToken(claims AccessTokenClaims) AuthToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return AuthToken{token: token}
}

func (t AuthToken) NewAccessToken() (string, *errs.AppError) {
	signedString, err := t.token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		return "", errs.NewUnexpectedError("cannot generate access token", err)
	}
	return signedString, nil
}

func (t AuthToken) NewRefreshToken() (string, *errs.AppError) {
	c := t.token.Claims.(AccessTokenClaims)
	refreshClaims := c.RefreshTokenClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    signedString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
    if err != nil {
		logger.Error("Failed while signing refresh token: " + err.Error())
	}
	return signedString, nil
}

func NewAccessTokenFromRefreshToken(refreshToken string) (string, *errs.AppError) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		return "", errs.NewAuthenticationError("invalid or expired refresh token",err)
	}
	r := token.Claims.(*RefreshTokenClaims)
	accessTokenClaims := r.AccessTokenClaims()
	authToken := NewAuthToken(accessTokenClaims)
	return authToken.NewAccessToken()
}
