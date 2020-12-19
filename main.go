package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-boot/conf"
	"gin-boot/router"

	"github.com/gin-gonic/gin"
)

// @title gin-boot API
// @version 1.0
// @description gin-boot server.

// @contact.email jhq0113@163.com

// @host 127.0.0.1
// @BasePath /
func main() {
	Bootstrap()

	server := &http.Server{
		Addr:    conf.Conf.App.Addr,
		Handler: router.Load(),
	}

	gin.SetMode(conf.Conf.App.Mode)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	//优雅退出
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-ch
	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(cxt)
	if err != nil {
		fmt.Println(err)
	}
}
