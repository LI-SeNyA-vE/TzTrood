package main

import (
	"TzTrood/internal/server/config"
	"TzTrood/internal/server/repository/nlp"
	"TzTrood/internal/server/repository/redisdb"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type RequestMessage struct {
	Text string `json:"text"`
}

func main() {
	l := logrus.New()
	log := logrus.NewEntry(l)

	cfg := config.NewMock()

	db := redisdb.NewRedis(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	router := gin.Default()

	router.POST("/message", func(ctx *gin.Context) {
		var mess RequestMessage

		if err := ctx.ShouldBindJSON(&mess); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid request body",
				"details": err.Error(),
			})

			log.Errorf("ошибка маршлинга запроса пользователя в структуру %w", err)
			return
		}

		if mess.Text == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "nil request body",
			})

			log.Error("пользователь оставил пустой запрос")
			return
		}

		// отправляем в сторонний NPL
		key, err := nlp.DetectIntentHTTP(cfg.NLPAddress, mess.Text)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "failed to detect intent",
			})
			log.Errorf("ошибка NLP: %v", err)
			return
		}

		response, err := db.KeyResponse.Search(ctx, key)
		if err != nil || response == "" {
			resp, err := http.Post(cfg.ServiceHumanAgent, "text/plain", strings.NewReader(mess.Text))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "server error",
				})

				log.Errorf("ошибка отправки запроса на сервис агент-человек. %w", err)
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "server error",
				})

				log.Errorf("ошибка чтения ответа от агента-человека: %w", err)
				return
			}
			defer resp.Body.Close()

			response = string(body)
			err = db.KeyResponse.Add(ctx, key, response)
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"response": response,
				})

				log.Errorf("ошибка сохранения в БЗ пары ключ/ответ от агента-человека. %w", err)
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"response": response,
			})

			log.Infof("на вопрос пользователя: \"%s\" | был отправлен ответ: \"%s\"", mess.Text, response)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"response": response,
		})
		log.Infof("на вопрос пользователя: \"%s\" | был отправлен ответ: \"%s\"", mess.Text, response)
		return
	})

	log.Info("Сервер запустился")
	router.Run(cfg.HttpPort)
}
