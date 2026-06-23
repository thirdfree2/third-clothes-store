package middleware

import (
	"net/http"
	"strings"

	httpcommon "clothing-service/internal/adapters/inbound/http/common"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwtSecret string
	jwtIssuer string
}

func NewAuthMiddleware(jwtSecret string, jwtIssuer string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
		jwtIssuer: jwtIssuer,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			httpcommon.RespondError(
				c,
				http.StatusUnauthorized,
				httpcommon.CodeUnauthorized,
				"missing authorization header",
				nil,
			)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			httpcommon.RespondError(
				c,
				http.StatusUnauthorized,
				httpcommon.CodeUnauthorized,
				"invalid authorization header",
				nil,
			)
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])

		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenUnverifiable
			}

			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			httpcommon.RespondError(
				c,
				http.StatusUnauthorized,
				httpcommon.CodeUnauthorized,
				"invalid or expired token",
				nil,
			)
			c.Abort()
			return
		}

		issuer, _ := claims["iss"].(string)
		if m.jwtIssuer != "" && issuer != m.jwtIssuer {
			httpcommon.RespondError(
				c,
				http.StatusUnauthorized,
				httpcommon.CodeUnauthorized,
				"invalid token issuer",
				nil,
			)
			c.Abort()
			return
		}

		tokenType, _ := claims["type"].(string)
		if tokenType != "admin_access" {
			httpcommon.RespondError(
				c,
				http.StatusUnauthorized,
				httpcommon.CodeUnauthorized,
				"invalid token type",
				nil,
			)
			c.Abort()
			return
		}

		c.Set("auth.email", claims["email"])
		c.Set("auth.role", claims["role"])
		c.Set("auth.permissions", extractPermissions(claims["permissions"]))

		c.Next()
	}
}

func (m *AuthMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawPermissions, exists := c.Get("auth.permissions")
		if !exists {
			httpcommon.RespondError(
				c,
				http.StatusForbidden,
				httpcommon.CodeForbidden,
				"permission denied",
				nil,
			)
			c.Abort()
			return
		}

		permissions, ok := rawPermissions.([]string)
		if !ok {
			httpcommon.RespondError(
				c,
				http.StatusForbidden,
				httpcommon.CodeForbidden,
				"permission denied",
				nil,
			)
			c.Abort()
			return
		}

		for _, p := range permissions {
			if p == permission {
				c.Next()
				return
			}
		}

		httpcommon.RespondError(
			c,
			http.StatusForbidden,
			httpcommon.CodeForbidden,
			"permission denied",
			nil,
		)
		c.Abort()
	}
}

func extractPermissions(value any) []string {
	result := []string{}

	switch permissions := value.(type) {
	case []any:
		for _, item := range permissions {
			if permission, ok := item.(string); ok {
				result = append(result, permission)
			}
		}

	case []string:
		result = append(result, permissions...)
	}

	return result
}
