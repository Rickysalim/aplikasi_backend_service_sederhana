package handlers

import (
	"fmt"
	"net/http"
	"product_auth/dto"
	"product_auth/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService services.AuthService
}

func (h *AuthHandler) Login(response *gin.Context) {
    var loginRequest dto.LoginRequest 
	if err := response.ShouldBindJSON(&loginRequest); err != nil {
		response.JSON(http.StatusBadRequest, dto.WriteResponse(http.StatusBadRequest, err.Error(),nil, err))
	} else {
		if loginResponse, err := h.AuthService.Login(response, &loginRequest); err != nil {
			response.JSON(err.Code, dto.WriteResponse(err.Code, err.Message, nil, err.Error))
			return
		} else {
			response.JSON(http.StatusOK,dto.WriteResponse(http.StatusOK,"OK",loginResponse,nil))
		}	
	}
}

func (h *AuthHandler) Verify(response *gin.Context) {
	urlParams := make(map[string]string) 

	for k := range response.Request.URL.Query() {
		urlParams[k] = response.Request.URL.Query().Get(k);
	}

	fmt.Println(urlParams)

	if urlParams["token"] != "" {
		appErr := h.AuthService.Verify(response, urlParams)
		if appErr != nil {
			response.JSON(appErr.Code, dto.WriteResponse(appErr.Code, appErr.Message, nil, appErr.Error))
			return
		} else {
			response.JSON(http.StatusOK, authorizedResponse())
			return
		}
	}
}

func (h AuthHandler) Refresh(response *gin.Context) {
    var refreshRequest dto.RefreshTokenRequest 
	if err := response.ShouldBindJSON(&refreshRequest); err != nil {
		response.JSON(http.StatusBadRequest, dto.WriteResponse(http.StatusBadRequest, err.Error(),nil, err))
	} else {
		if  token, err := h.AuthService.Refresh(response, refreshRequest); err != nil {
			response.JSON(http.StatusInternalServerError, dto.WriteResponse(err.Code, err.Message, nil, err.Error))
			return
		} else {
			response.JSON(http.StatusOK,dto.WriteResponse(http.StatusOK,"OK",*token,nil))
		}	
	}
}

func authorizedResponse() map[string]bool {
	return map[string]bool{"is_authorized": true}
}