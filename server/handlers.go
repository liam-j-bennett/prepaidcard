package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"prepaidcard/models"
)

type CardRequest struct {
	CardNumber	string	`json:"card_number,omitempty"`
	MerchantId	string	`json:"merchant_id,omitempty"`
	Amount		int64	`json:"amount"`
}

func handleError(err error, c *gin.Context) {
	if e, ok := err.(models.Error); ok {
		c.AbortWithStatusJSON(e.Code(), gin.H{
			"details": fmt.Sprintln(e),
		})
		return
	} else {
		log.Error(err)
		c.AbortWithStatusJSON(500, gin.H{
			"details": fmt.Sprintln(err),
		})
	}
}

func (s *Server) createCard(c *gin.Context) {
	newCard, err := s.store.CreateCard()
	if err != nil {
		handleError(err, c)
		return
	}
	c.JSON(200, newCard)
}

func (s *Server) getCard(c *gin.Context) {
	cardId := c.Param("cardId")
	card, err := s.store.GetCard(cardId)
	if err != nil {
		handleError(err, c)
		return
	}
	c.JSON(200, card)
}

func (s *Server) listSpending(c *gin.Context) {
	cardId := c.Param("cardId")
	transactionList, err := s.store.TransactionList(cardId)
	if err != nil {
		handleError(err, c)
		return
	}
	c.JSON(200, transactionList)
}

func (s *Server) loadCard(c *gin.Context) {
	cardId := c.Param("cardId")
	var request CardRequest
	if err := c.BindJSON(&request); err != nil {
		return
	}
	if request.Amount <= 0 {
		err := models.InvalidAmount
		handleError(err, c)
		return
	}
	card, err := s.store.LoadCard(cardId, request.Amount)
	if err != nil {
		handleError(err, c)
		return
	}
	c.JSON(200, card)
}

func (s *Server) authRequest(c *gin.Context) {
	var request CardRequest
	if err := c.BindJSON(&request); err != nil {
		return
	}
	if request.Amount <= 0 {
		err := models.InvalidAmount
		handleError(err, c)
		return
	}
	merchant, err := s.store.GetMerchant(request.MerchantId)
	if err != nil {
		handleError(err, c)
		return
	}
	card, err := s.store.GetCard(request.CardNumber)
	if err != nil {
		handleError(err, c)
		return
	}
	transaction, err := s.store.Auth(card, merchant, request.Amount)
	if err != nil {
		handleError(err, c)
		return
	}
	c.JSON(200, transaction)
}

func (s *Server) captureTransaction(c *gin.Context) {
	transactionId := c.Param("transactionId")
	transaction, err := s.store.GetTransaction(transactionId)
	if err != nil {
		handleError(err, c)
		return
	}
	var request CardRequest
	if err := c.BindJSON(&request); err != nil {
		return
	}
	if request.Amount <= 0 || request.Amount > transaction.AuthorizedAmount {
		err := models.InvalidAmount
		handleError(err, c)
		return
	}
	if err = s.store.Capture(transaction, request.Amount); err != nil {
		handleError(err, c)
		return
	}
	c.JSON(200, transaction)
}

func (s *Server) reverseTransaction(c *gin.Context) {
	transactionId := c.Param("transactionId")
	transaction, err := s.store.GetTransaction(transactionId)
	if err != nil {
		handleError(err, c)
		return
	}
	var request CardRequest
	if err := c.BindJSON(&request); err != nil {
		return
	}
	if request.Amount <= 0 || request.Amount > transaction.AuthorizedAmount {
		err := models.InvalidAmount
		handleError(err, c)
		return
	}
	if err = s.store.Reverse(transaction, request.Amount); err != nil {
		handleError(err, c)
		return
	}
	c.JSON(200, transaction)
}

func (s *Server) refundCapture(c *gin.Context) {
	transactionId := c.Param("transactionId")
	transaction, err := s.store.GetTransaction(transactionId)
	if err != nil {
		handleError(err, c)
		return
	}
	var request CardRequest
	if err := c.BindJSON(&request); err != nil {
		return
	}
	if request.Amount < 0 || request.Amount > transaction.CapturedAmount {
		err := models.InvalidAmount
		handleError(err, c)
		return
	}
	err = s.store.Refund(transaction, request.Amount)
	if err != nil {
		handleError(err, c)
		return
	}
	c.JSON(200, transaction)
}
