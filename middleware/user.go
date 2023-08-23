package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"whitelist/service/admins"
	"whitelist/utils"
)

// UserAuthorized 校验user中间件
func UserAuthorized() gin.HandlerFunc {
	return func(c *gin.Context){
		tokenString 	:= c.GetHeader("token")
		if tokenString == "" {
			utils.CreateResponse(c).Unauthorized()
			c.Abort()
		}
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(admins.AdminPasswordKey), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("admin",claims)
		} else {
			utils.CreateResponse(c).Unauthorized()
			c.Abort()
		}
	}
}
