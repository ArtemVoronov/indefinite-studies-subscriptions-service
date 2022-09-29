package subscribrions

import (
	"log"
	"net/http"

	"github.com/ArtemVoronov/indefinite-studies-subscriptions-service/internal/services"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/api"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/api/validation"
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
