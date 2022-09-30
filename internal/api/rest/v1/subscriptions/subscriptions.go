package subscribrions

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ArtemVoronov/indefinite-studies-subscriptions-service/internal/services"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/api"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/api/validation"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/services/kafka"
	"github.com/gin-gonic/gin"
)

type AddEventDTO struct {
	EventType string `json:"EventType" binding:"required"`
	EventBody string `json:"EventBody" binding:"required"`
}

func AddEvent(c *gin.Context) {
	var dto AddEventDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		validation.SendError(c, err)
		return
	}

	err := services.Instance().KafkaProducer().CreateMessage(dto.EventType, dto.EventBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unable to add event")
		log.Printf("Unable to add event: %s", err)
		return
	}

	c.IndentedJSON(http.StatusOK, api.DONE)
}

func AddSendEmailEvent(c *gin.Context) {
	var dto kafka.SendEmailEvent
	if err := c.ShouldBindJSON(&dto); err != nil {
		validation.SendError(c, err)
		return
	}

	data, err := json.Marshal(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unable to add SEND_EMAIL event")
		log.Printf("Unable to add SEND_EMAIL event: %s", err)
		return
	}

	err = services.Instance().KafkaProducer().CreateMessage(kafka.EVENT_TYPE_SEND_EMAIL, string(data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unable to add SEND_EMAIL event")
		log.Printf("Unable to add SEND_EMAIL event: %s", err)
		return
	}

	c.IndentedJSON(http.StatusOK, api.DONE)
}
