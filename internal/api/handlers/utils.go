package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type ContextData struct {
	FamilyID int `json:"family_id"`
	UserID   int `json:"user_id"`
	RoleID   int `json:"role_id"`
}

// extractClaims help to extract the JWT claims
func extractClaims(c *gin.Context) jwt.MapClaims {

	if _, exists := c.Get("JWT_PAYLOAD"); !exists {
		emptyClaims := make(jwt.MapClaims)
		return emptyClaims
	}

	jwtClaims, _ := c.Get("JWT_PAYLOAD")

	return jwtClaims.(jwt.MapClaims)
}

func getClaimsFromContext(c *gin.Context) *ContextData {
	claims := extractClaims(c)
	return &ContextData{
		UserID:   int(claims["user_id"].(float64)),
		FamilyID: int(claims["family_id"].(float64)),
		RoleID:   int(claims["role_id"].(float64)),
	}
}
