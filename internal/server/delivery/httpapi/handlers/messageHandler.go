package handlers

import (
	"TzTrood/internal/server/config"
	"TzTrood/internal/server/repository"
	"TzTrood/internal/server/repository/nlp"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type RequestMessage struct {
	Text string `json:"text"`
}

type handlersGIN struct {
	cfg *config.Server
	db  *repository.DataBase
	log *logrus.Entry
}

func (h handlersGIN) PostUserQuestion(ctx *gin.Context) {
	var mess RequestMessage

	if err := ctx.ShouldBindJSON(&mess); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})

		h.log.Errorf("ошибка маршлинга запроса пользователя в структуру %w", err)
		return
	}

	if mess.Text == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "nil request body",
		})

		h.log.Error("пользователь оставил пустой запрос")
		return
	}

	// отправляем в сторонний NPL
	key, err := nlp.DetectIntentHTTP(h.cfg.NLPAddress, mess.Text)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to detect intent",
		})
		h.log.Errorf("ошибка NLP: %v", err)
		return
	}

	response, err := h.db.KeyResponse.Search(ctx, key)
	if err != nil || response == "" {
		resp, err := http.Post(h.cfg.ServiceHumanAgent, "text/plain", strings.NewReader(mess.Text))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "server error",
			})

			h.log.Errorf("ошибка отправки запроса на сервис агент-человек. %w", err)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "server error",
			})

			h.log.Errorf("ошибка чтения ответа от агента-человека: %w", err)
			return
		}
		defer resp.Body.Close()

		response = string(body)
		err = h.db.KeyResponse.Add(ctx, key, response)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"response": response,
			})

			h.log.Errorf("ошибка сохранения в БЗ пары ключ/ответ от агента-человека. %w", err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"response": response,
		})

		h.log.Infof("на вопрос пользователя: \"%s\" | был отправлен ответ: \"%s\"", mess.Text, response)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"response": response,
	})
	h.log.Infof("на вопрос пользователя: \"%s\" | был отправлен ответ: \"%s\"", mess.Text, response)
	return
}

func NewHandlers(cfg *config.Server, db *repository.DataBase, log *logrus.Entry) *handlersGIN {
	return &handlersGIN{
		cfg: cfg,
		db:  db,
		log: log,
	}
}
