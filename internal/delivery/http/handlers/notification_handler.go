package handlers

import (
	"net/http"
	"strconv"

	"github.com/corechain/notification-service/internal/application/services"
	"github.com/corechain/notification-service/internal/delivery/http/response"
	"github.com/corechain/notification-service/internal/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// GetUserNotifications godoc
// @Summary Get notifications by user ID
// @Description Get all notifications for a specific user with pagination
// @Tags notifications
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/notifications/{userId} [get]
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		response.Error(c, http.StatusBadRequest, "User ID is required")
		return
	}

	// Parse pagination parameters
	limit := 20
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 100 {
				limit = 100 // Cap at 100
			}
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	notifications, err := h.notificationService.GetUserNotifications(c.Request.Context(), userID, limit, offset)
	if err != nil {
		logger.Error("Failed to get user notifications",
			zap.Error(err),
			zap.String("user_id", userID),
		)
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve notifications")
		return
	}

	response.JSON(c, http.StatusOK, gin.H{
		"notifications": notifications,
		"count":         len(notifications),
		"limit":         limit,
		"offset":        offset,
	})
}

// GetNotificationDetail godoc
// @Summary Get notification by ID
// @Description Get a single notification by its ID
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/notifications/detail/{id} [get]
func (h *NotificationHandler) GetNotificationDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Notification ID is required")
		return
	}

	notification, err := h.notificationService.GetNotificationByID(c.Request.Context(), id)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "notification not found" {
			response.Error(c, http.StatusNotFound, "Notification not found")
			return
		}

		logger.Error("Failed to get notification detail",
			zap.Error(err),
			zap.String("notification_id", id),
		)
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve notification")
		return
	}

	response.JSON(c, http.StatusOK, notification)
}
