package dto

import "github.com/gin-gonic/gin"

func WriteResponse(code int, message string, result interface{}, err error) interface{} {
	if err != nil {
		data := map[string]interface{}{
			"error": gin.H{
				"status": code,
				"message": message,
			},
		}
		return data
	}
	data := map[string]interface{}{
		"success": gin.H{
			"status": code,
			"message": message,
			"data": result,
		},
	}
	return data
}