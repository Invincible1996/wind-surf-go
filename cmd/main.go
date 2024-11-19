package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"wind-surf-go/internal/model"
	"wind-surf-go/internal/router"
)

func main() {

	// 读取配置文件 ../config/config.yaml
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Printf("Current working directory: %s", dir)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/Users/kevin/Documents/go-project/wind-surf-go/internal/config")
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// Initialize database
	db, err := gorm.Open(sqlite.Open("wind-surf.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic("failed to migrate database")
	}

	// Setup router
	r := router.SetupRouter(db)

	// 创建 http.Server
	srv := &http.Server{
		Addr:    ":" + viper.GetString("server.port"),
		Handler: r,
	}

	// 在独立的 goroutine 中启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 设置关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
