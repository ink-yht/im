//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/ink-yht/im/internal/repository/dao/user_dao"
	"github.com/ink-yht/im/internal/repository/user_repo"
	"github.com/ink-yht/im/internal/service/user_service"
	"github.com/ink-yht/im/internal/web/user_web"
	"github.com/ink-yht/im/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitLogger,

		// DAO 部分
		user_dao.NewUserDAO,

		// cache 部分

		// repository 部分
		user_repo.NewUserRepository,

		// service 部分
		user_service.NewUserService,

		// Handler 部分
		user_web.NewUserHandler,

		// 中间件
		ioc.InitWebServer,
		ioc.InitMiddleWares,
	)
	return new(gin.Engine)
}
