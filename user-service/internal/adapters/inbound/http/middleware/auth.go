package middleware

import (
	"net/http"
	"strconv"
	"strings"

	httpcommon "user-service/internal/adapters/inbound/http/common"

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
		tokenString := bearerToken(c.GetHeader("Authorization"))
		if tokenString == "" {
			httpcommon.RespondError(c, http.StatusUnauthorized, httpcommon.CodeUnauthorized, "missing authorization header", nil)
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenUnverifiable
			}

			return []byte(m.jwtSecret), nil
		})
		if err != nil || !token.Valid {
			httpcommon.RespondError(c, http.StatusUnauthorized, httpcommon.CodeUnauthorized, "invalid or expired token", nil)
			c.Abort()
			return
		}

		if m.jwtIssuer != "" && claims["iss"] != m.jwtIssuer {
			httpcommon.RespondError(c, http.StatusUnauthorized, httpcommon.CodeUnauthorized, "invalid token issuer", nil)
			c.Abort()
			return
		}

		userID, ok := subjectAsInt64(claims["sub"])
		if !ok {
			httpcommon.RespondError(c, http.StatusUnauthorized, httpcommon.CodeUnauthorized, "invalid token subject", nil)
			c.Abort()
			return
		}

		c.Set("auth.user_id", userID)
		c.Set("auth.email", claims["email"])
		c.Set("auth.user_type", claims["user_type"])
		c.Set("auth.role", claims["role"])
		c.Next()
	}
}

func bearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

func subjectAsInt64(value any) (int64, bool) {
	switch sub := value.(type) {
	case float64:
		return int64(sub), true
	case string:
		id, err := strconv.ParseInt(sub, 10, 64)
		return id, err == nil
	default:
		return 0, false
	}
}
