package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
	"novaro-server/service"
)

type EventsApi struct {
	service *service.EventsService
}

func NewEventApi() *EventsApi {
	return &EventsApi{
		service: service.NewEventsService(),
	}
}

func (api *EventsApi) CreateEvents(c *gin.Context) {
	var event model.Events
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	if err := api.service.Create(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Successfully added event"})
}

func (api *EventsApi) DeleteEvents(c *gin.Context) {
	id := c.Query("id")
	if id != "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}
	if err := api.service.Delete(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Successfully deleted event"})
}

func (api *EventsApi) UpdateEvents(c *gin.Context) {
	var event model.Events
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	if err := api.service.Update(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Successfully updated event"})
}

func (api *EventsApi) GetId(c *gin.Context) {
	value := c.Query(("id"))
	if value == "" {
		c.JSON(400, gin.H{"error": "id is required"})
	}

	resp, err := api.service.GetById(value)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": resp,
	})
}

func (api *EventsApi) GetList(c *gin.Context) {
	var event model.Events
	c.ShouldBindJSON(&event)

	resp, err := api.service.GetList(&event)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": resp,
	})
}
