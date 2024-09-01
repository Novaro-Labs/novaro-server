package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
)

type EventsApi struct {
}

var events = model.Events{}

func (EventsApi) CreateEvents(c *gin.Context) {
	var event model.Events
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	if err := events.Create(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Successfully added event"})
}

func (EventsApi) DeleteEvents(c *gin.Context) {
	id := c.Query("id")
	if id != "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}
	if err := events.Delete(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Successfully deleted event"})
}

func (EventsApi) UpdateEvents(c *gin.Context) {
	var event model.Events
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	if err := events.Update(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Successfully updated event"})
}

func (EventsApi) GetId(c *gin.Context) {
	value := c.Query(("id"))
	if value == "" {
		c.JSON(400, gin.H{"error": "id is required"})
	}

	resp, err := events.Get(value)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": resp,
	})
}

func (EventsApi) GetList(c *gin.Context) {
	var event model.Events
	c.ShouldBindJSON(&event)

	resp, err := events.GetList(&event)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": resp,
	})
}
