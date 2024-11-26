package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ink-yht/im/internal/repository/dao"
	"github.com/ink-yht/im/internal/repository/dao/user_dao"
	"github.com/ink-yht/im/internal/repository/user_repo"
	"github.com/ink-yht/im/internal/service/user_service"
	"github.com/ink-yht/im/internal/web/user_web"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13326)/im?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	server := gin.Default()
	userDao := user_dao.NewUserDAO(db)
	userRepo := user_repo.NewUserRepository(userDao)
	svc := user_service.NewUserServiceImpl(userRepo)
	hdl := user_web.NewUserHandler(svc)
	hdl.RegisterRoutes(server)

	err = server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
