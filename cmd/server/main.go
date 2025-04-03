package main

import (
	"TzTrood/internal/server/config"
	"TzTrood/internal/server/delivery/httpapi/handlers"
	"TzTrood/internal/server/repository/redisdb"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func main() {
	l := logrus.New()
	log := logrus.NewEntry(l)

	cfg := config.NewMock()

	db := redisdb.NewRedis(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	h := handlers.NewHandlers(cfg, db, log)

	router := gin.Default()

	router.POST("/message", func(ctx *gin.Context) {
		h.PostUserQuestion(ctx)
	})

	log.Info("Сервер запустился")
	router.Run(cfg.HttpPort)
}
