package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ink-yht/im/internal/web/user_web"
	"github.com/ink-yht/im/pkg/logger"
	"net/http"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc,

	userHdl *user_web.UserHandler,
) *gin.Engine {

	server := gin.Default()
	server.StaticFS("uploads", http.Dir("uploads"))
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	return server
}

func InitMiddleWares(l logger.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{

		corsHdl(),

		//log.NewMiddlewaresLoggerBuilder(func(ctx context.Context, al *log.AccessLog) {
		//	l.Debug("HTTP请求", logger.Field{Key: "al", Value: al})
		//}).AllowReqBody().AllowRespBody().Build(),

		//middlelware.NewLoginJWTMiddlewareBuilder().
		//	IgnorePaths("/users/signup").
		//	IgnorePaths("/users/login").Build(),
		//
		//ratelimit.NewBuilder(redisClient, time.Minute, 100).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins:     []string{"https://foo.com"},
		//AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
