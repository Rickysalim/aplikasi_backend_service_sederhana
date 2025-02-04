package services

import (
	"context"
	"errors"
	"fmt"
	"product_auth/domain"
	"product_auth/dto"
	"product_auth/dto/errs"
	"product_auth/repositories"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Verify(ctx context.Context, urlParams map[string]string) *errs.AppError
    Refresh(ctx context.Context, request dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError)

}

type DefaultAuthService struct {
	repo            repositories.AuthRepository
	rolePermissions domain.RolePermissions
}

func (s DefaultAuthService) Refresh(ctx context.Context, request dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError) {
	if vErr := request.IsAccessTokenValid(); vErr != nil {
		if vErr.Errors == jwt.ValidationErrorExpired {
			var appErr *errs.AppError
			if appErr = s.repo.RefreshTokenExists(ctx, request.RefreshToken); appErr != nil {
			   return nil, appErr	
			}
			var accessToken string 
			if accessToken, appErr = domain.NewAccessTokenFromRefreshToken(request.RefreshToken); appErr != nil {
				return nil, appErr
			}
			return &dto.LoginResponse{AccessToken: accessToken}, nil
		}
		return nil, errs.NewAuthenticationError("invalid token", vErr)
	}
	return nil, errs.NewAuthenticationError("cannot generate a new access token until the current one expires", errors.New("cannot generate a new access token until the current one expires"))
}

func (s DefaultAuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	var err *errs.AppError
	var login *domain.Login

	if login, err = s.repo.FindByEmailAndPassword(ctx, req.Email, req.Password); err != nil {
		return nil, err
	}

	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims)

	var accessToken , refreshToken string

	if accessToken, err = authToken.NewAccessToken(); err != nil {
		return nil, err
	}

	if refreshToken, err = s.repo.GenerateAndSaveRefreshTokenToStore(ctx, authToken, claims.CustomersId); err != nil {
		return nil, err 
	}

	return &dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil

}

func (s DefaultAuthService) Verify(ctx context.Context, urlParams map[string]string) *errs.AppError {
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return errs.NewAuthenticationError(err.Error(), err)
	} else {
		if jwtToken.Valid {
			claims := jwtToken.Claims.(*domain.AccessTokenClaims)
			if err := claims.IsRequestVerifiedWithTokenClaims(urlParams); err != nil {
				return errs.NewAuthorizationError(err.Error(), err)
			}
			if err := claims.IsSameRole(urlParams); err != nil {
				return errs.NewAuthorizationError(err.Error(), err)
		    }
			if err := s.rolePermissions.IsAuthorizedFor(urlParams["role"], urlParams["route_name"]); err != nil {
				return errs.NewAuthorizationError(fmt.Sprintf("%s role is not authorized", claims.Username), err)
			}
			return nil
		} else {
			return errs.NewAuthorizationError("Invalid token", err)
		}
	}
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func NewAuthService(repo repositories.AuthRepository, permissions domain.RolePermissions) DefaultAuthService {
	return DefaultAuthService{repo, permissions}
}
