// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ink-yht/im/internal/repository/dao/user_dao"
	"github.com/ink-yht/im/internal/repository/user_repo"
	"github.com/ink-yht/im/internal/service/user_service"
	"github.com/ink-yht/im/internal/web/user_web"
	"github.com/ink-yht/im/ioc"
)

import (
	_ "github.com/spf13/viper/remote"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	logger := ioc.InitLogger()
	v := ioc.InitMiddleWares(logger)
	db := ioc.InitDB(logger)
	userDao := user_dao.NewUserDAO(db)
	userRepository := user_repo.NewUserRepository(userDao)
	userService := user_service.NewUserService(userRepository)
	userHandler := user_web.NewUserHandler(userService)
	engine := ioc.InitWebServer(v, userHandler)
	return engine
}
