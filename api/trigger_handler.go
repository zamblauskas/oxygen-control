package api

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zamblauskas/oxygen-control/models"
	"github.com/zamblauskas/oxygen-control/service"
)

func NewTriggerHandler(triggerService *service.TriggerService) *TriggerHandler {
	return &TriggerHandler{triggerService: triggerService}
}

type TriggerHandler struct {
	triggerService *service.TriggerService
}

func (h *TriggerHandler) AddTrigger(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trigger, err := models.ParseTriggerFromJson(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.triggerService.AddTrigger(trigger)
}

func (h *TriggerHandler) GetTriggers(c *gin.Context) {
	triggers, err := h.triggerService.GetTriggers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, triggers)
}
