package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"prepaidcard/models"
)

type Server struct {
	Router *gin.Engine
	store models.CardStore
}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"hello": "world!"})
}

func (s *Server) bindHandlers() {
	router := s.Router
	router.GET("/", handlePing)
	router.POST("/cards", s.createCard)
	router.GET("/cards/:cardId", s.getCard)
	router.GET("cards/:cardId/spending", s.listSpending)
	router.POST("/cards/:cardId", s.loadCard)
	router.POST("/transactions", s.authRequest)
	router.PATCH("/transactions/:transactionId/capture", s.captureTransaction)
	router.PATCH("/transactions/:transactionId/reverse", s.reverseTransaction)
	router.PATCH("/transactions/:transactionId/refund", s.refundCapture)
}

func InitServer(store models.CardStore) *Server {
	router := gin.Default()
	server := Server{
		Router: router,
		store: store,
	}
	server.bindHandlers()
	return &server
}