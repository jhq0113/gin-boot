package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-boot/global"

	"github.com/gin-gonic/gin"
)

// @title gin-boot API
// @version 1.0
// @description gin-boot server.

// @contact.email jhq0113@163.com

// @host 127.0.0.1
// @BasePath /
func main() {
	global.Bootstrap()

	r := gin.Default()
	server := &http.Server{
		Addr:           "",
		Handler:        r,
		ReadTimeout:    time.Millisecond * 200,
		WriteTimeout:   time.Millisecond * 200,
		MaxHeaderBytes: 1024*1024,
	}

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
