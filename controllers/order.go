package controllers

import (
	"assignment-2/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (c *Controller) CreateOrder(ctx *gin.Context) {
	var (
		newOrder models.Order
	)

	if err := ctx.ShouldBindJSON(&newOrder); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": fmt.Sprintf("Error Creating Order Data: %v", err),
		})
		return
	}

	newOrder.OrderedAt = time.Now()

	if err := c.DB.Create(&newOrder).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": fmt.Sprintf("Error Creating Order Data: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"result": newOrder,
	})
}

func (c *Controller) GetOrders(ctx *gin.Context) {
	var (
		orders []models.Order
	)

	if err := c.DB.Model(&models.Order{}).Preload("Items").Order("order_id asc").Find(&orders).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": fmt.Sprintf("Error Getting Order Data: %v", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": orders,
	})
}

func (c *Controller) UpdateOrder(ctx *gin.Context) {
	var (
		orderId      string
		updatedOrder models.Order
		countItem    int
	)

	orderId = ctx.Param("orderId")

	if err := ctx.ShouldBindJSON(&updatedOrder); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"result": fmt.Sprintf("Order with id %v Not Found", orderId),
		})
		return
	}

	if err := c.DB.First(&models.Order{}, "order_id = ?", orderId).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"result": fmt.Sprintf("Order with id %v Not Found", orderId),
		})
		return
	}

	if err := c.DB.Model(&models.Order{}).Where("order_id = ?", orderId).Updates(updatedOrder).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": fmt.Sprintf("Error Updating Order Data: %v", err.Error()),
		})
		return
	}

	countItem = 0
	for _, item := range updatedOrder.Items {
		if err := c.DB.First(&models.Item{}, "order_id = ? and item_id = ?", orderId, item.ItemId).Error; err == nil {
			c.DB.Model(&models.Item{}).Where("item_id = ?", item.ItemId).Updates(item)
			countItem += 1
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprintf("Order with id %v Has Been Successfully Updated (%v Items Updated)", orderId, countItem),
	})
}

func (c *Controller) DeleteOrder(ctx *gin.Context) {
	var (
		orderId string
	)

	orderId = ctx.Param("orderId")

	if err := c.DB.First(&models.Order{}, "order_id = ?", orderId).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"result": fmt.Sprintf("Order with id %v Not Found", orderId),
		})
		return
	}

	if err := c.DB.Where("order_id = ?", orderId).Delete(&models.Order{}).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": fmt.Sprintf("Error Deleting Order: %v", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprintf("Order with id %v Has Been Successfully Deleted", orderId),
	})
}
