package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/corechain/notification-service/internal/delivery/http/handlers"
	"github.com/corechain/notification-service/internal/delivery/http/middleware"
	"github.com/corechain/notification-service/internal/utils/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	router              *gin.Engine
	httpServer          *http.Server
	notificationHandler *handlers.NotificationHandler
}

type ServerConfig struct {
	Port                int
	NotificationHandler *handlers.NotificationHandler
}

func NewServer(config ServerConfig) *Server {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Add middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestLogger())

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	server := &Server{
		router:              router,
		notificationHandler: config.NotificationHandler,
	}

	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "notification-service",
		})
	})

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Notification routes
		notifications := v1.Group("/notifications")
		{
			notifications.GET("/:userId", s.notificationHandler.GetUserNotifications)
			notifications.GET("/detail/:id", s.notificationHandler.GetNotificationDetail)
		}
	}
}

func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	logger.Info("HTTP server starting", zap.String("address", addr))

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down HTTP server...")
	
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown HTTP server: %w", err)
		}
	}

	logger.Info("HTTP server stopped")
	return nil
}
