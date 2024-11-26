package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ink-yht/im/internal/service/user_service"
	"log"
	"net/http"
	"strings"
	"time"
)

// LoginJWTMiddlewareBuilder JWT 登录校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		// 我现在用 JWT 来校验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//segs := strings.SplitN(tokenHeader, " ", 2)
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			// 没登录，有人瞎搞
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &user_service.UserClaims{}
		// ParseWithClaims 里面，一定要传入指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return user_service.JWTKey, nil
		})
		if err != nil {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//claims.ExpiresAt.Time.Before(time.Now()) {
		//	// 过期了
		//}
		// err 为 nil，token 不为 nil
		if token == nil || !token.Valid || claims.Id == 0 {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			// 严重的安全问题
			// 要监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.ExpiresAt.Time.Before(time.Now()) {
			// 如果 token 过期了
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Hour*60 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 3))
			tokenStr, err = token.SignedString(user_service.JWTKey)
			if err != nil {
				// 记录日志
				log.Println("jwt 续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
		ctx.Set("claims", claims)
	}
}
