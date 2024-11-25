package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ink-yht/im/internal/repository/dao"
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
	hdl := user_web.NewUserHandler()
	hdl.RegisterRoutes(server)

	err = server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
