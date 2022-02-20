package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/peienxie/go-bank/db/sqlc"
)

func (s *Server) initTransferRoutes() {
	s.router.POST("/transfers", s.createTransfer)
}

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,oneof=USD TWD"`
}

func (s *Server) createTransfer(c *gin.Context) {
	var req createTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if ok := s.validAccount(c, req.FromAccountID, req.Currency); !ok {
		return
	}

	if ok := s.validAccount(c, req.ToAccountID, req.Currency); !ok {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	result, err := s.store.TransferTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) validAccount(c *gin.Context, id int64, currency string) bool {
	account, err := s.store.GetAccount(c, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account %d currency expect %s, but got %s",
			id, currency, account.Currency)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
